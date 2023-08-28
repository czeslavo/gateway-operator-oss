package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorv1beta1 "github.com/kong/gateway-operator/apis/v1beta1"
	"github.com/kong/gateway-operator/internal/consts"
	k8sutils "github.com/kong/gateway-operator/internal/utils/kubernetes"
	"github.com/kong/gateway-operator/internal/versions"
)

// -----------------------------------------------------------------------------
// DataPlane - Private Functions - Generators
// -----------------------------------------------------------------------------

func generateDataPlaneImage(dataplane *operatorv1beta1.DataPlane, validators ...versions.VersionValidationOption) (string, error) {
	if dataplane.Spec.DataPlaneOptions.Deployment.PodTemplateSpec == nil {
		return consts.DefaultDataPlaneImage, nil // TODO: https://github.com/Kong/gateway-operator/issues/20
	}

	container := k8sutils.GetPodContainerByName(&dataplane.Spec.DataPlaneOptions.Deployment.PodTemplateSpec.Spec, consts.DataPlaneProxyContainerName)
	if container != nil && container.Image != "" {
		for _, v := range validators {
			supported, err := v(container.Image)
			if err != nil {
				return "", err
			}
			if !supported {
				return "", fmt.Errorf("unsupported DataPlane image %s", container.Image)
			}
		}
		return container.Image, nil
	}

	if relatedKongImage := os.Getenv("RELATED_IMAGE_KONG"); relatedKongImage != "" {
		// RELATED_IMAGE_KONG is set by the operator-sdk when building the operator bundle.
		// https://github.com/Kong/gateway-operator/issues/261
		return relatedKongImage, nil
	}

	return consts.DefaultDataPlaneImage, nil // TODO: https://github.com/Kong/gateway-operator/issues/20
}

// -----------------------------------------------------------------------------
// DataPlane - Private Functions - Kubernetes Object Labels and Annotations
// -----------------------------------------------------------------------------

func addAnnotationsForDataplaneIngressService(obj client.Object, dataplane operatorv1beta1.DataPlane) {
	specAnnotations := extractDataPlaneIngressServiceAnnotations(&dataplane)
	if specAnnotations == nil {
		return
	}
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	for k, v := range specAnnotations {
		annotations[k] = v
	}
	encodedSpecAnnotations, err := json.Marshal(specAnnotations)
	if err == nil {
		annotations[consts.AnnotationLastAppliedAnnotations] = string(encodedSpecAnnotations)
	}
	obj.SetAnnotations(annotations)
}

func extractDataPlaneIngressServiceAnnotations(dataplane *operatorv1beta1.DataPlane) map[string]string {
	if dataplane.Spec.DataPlaneOptions.Network.Services == nil ||
		dataplane.Spec.DataPlaneOptions.Network.Services.Ingress == nil ||
		dataplane.Spec.DataPlaneOptions.Network.Services.Ingress.Annotations == nil {
		return nil
	}

	anns := dataplane.Spec.DataPlaneOptions.Network.Services.Ingress.Annotations
	return anns
}

// extractOutdatedDataPlaneIngressServiceAnnotations returns the last applied annotations
// of ingress service from `DataPlane` spec but disappeared in current `DataPlane` spec.
func extractOutdatedDataPlaneIngressServiceAnnotations(
	dataplane *operatorv1beta1.DataPlane, existingAnnotations map[string]string,
) (map[string]string, error) {
	if existingAnnotations == nil {
		return nil, nil
	}
	lastAppliedAnnotationsEncoded, ok := existingAnnotations[consts.AnnotationLastAppliedAnnotations]
	if !ok {
		return nil, nil
	}
	outdatedAnnotations := map[string]string{}
	err := json.Unmarshal([]byte(lastAppliedAnnotationsEncoded), &outdatedAnnotations)
	if err != nil {
		return nil, fmt.Errorf("failed to decode last applied annotations: %w", err)
	}
	// If an annotation is present in last applied annotations but not in current spec of annotations,
	// the annotation is outdated and should be removed.
	// So we remove the annotations present in current spec in last applied annotations,
	// the remaining annotations are outdated and should be removed.
	currentSpecifiedAnnotations := extractDataPlaneIngressServiceAnnotations(dataplane)
	for k := range currentSpecifiedAnnotations {
		delete(outdatedAnnotations, k)
	}
	return outdatedAnnotations, nil
}

// ensureDataPlaneReadyStatus ensures that the provided DataPlane gets an up to
// date Ready status condition.
// It sets the condition based on the readiness of DataPlane's Deployment and
// its ingress Service receiving an address.
func ensureDataPlaneReadyStatus(
	ctx context.Context,
	cl client.Client,
	log logr.Logger,
	dataplane *operatorv1beta1.DataPlane,
) (ctrl.Result, error) {
	deployments, err := k8sutils.ListDeploymentsForOwner(ctx,
		cl,
		dataplane.Namespace,
		dataplane.UID,
		client.MatchingLabels{
			"app":                                dataplane.Name,
			consts.DataPlaneDeploymentStateLabel: consts.DataPlaneStateLabelValueLive,
		},
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed listing deployments for DataPlane %s/%s: %w", dataplane.Namespace, dataplane.Name, err)
	}
	if len(deployments) != 1 {
		info(log, "expected only 1 Deployment for DataPlane", dataplane)
		return ctrl.Result{Requeue: true}, nil
	}

	deployment := deployments[0]
	// We check if the Deployment is not Ready.
	// This is the case when status has replicas set to 0 or status.availableReplicas
	// in status is less than status.replicas.
	// The second condition takes into account the time when new version (ReplicaSet)
	// is being rolled out by Deployment controller and there might be more available
	// replicas than specified in spec.replicas but we don't consider it fully ready
	// until it stabilized to be equal to status.replicas.
	// If any of those conditions is specified we mark the DataPlane as not ready yet.
	if deployment.Status.Replicas == 0 || deployment.Status.AvailableReplicas < deployment.Status.Replicas {
		debug(log, "Deployment for DataPlane not ready yet", dataplane)

		// Set Ready to false for dataplane as the underlying deployment is not ready.
		k8sutils.SetCondition(
			k8sutils.NewConditionWithGeneration(
				k8sutils.ReadyType,
				metav1.ConditionFalse,
				k8sutils.WaitingToBecomeReadyReason,
				k8sutils.WaitingToBecomeReadyMessage,
				dataplane.Generation,
			),
			dataplane,
		)
		ensureDataPlaneReadinessStatus(dataplane, &deployment)
		if err = patchDataPlaneStatus(ctx, cl, log, dataplane); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed patching status (Deployment not ready) for DataPlane %s/%s: %w", dataplane.Namespace, dataplane.Name, err)
		}
		return ctrl.Result{}, nil
	}

	services, err := k8sutils.ListServicesForOwner(ctx,
		cl,
		dataplane.Namespace,
		dataplane.UID,
		client.MatchingLabels{
			"app":                             dataplane.Name,
			consts.DataPlaneServiceStateLabel: consts.DataPlaneStateLabelValueLive,
			consts.DataPlaneServiceTypeLabel:  string(consts.DataPlaneIngressServiceLabelValue),
		},
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed listing ingress services for DataPlane %s/%s: %w", dataplane.Namespace, dataplane.Name, err)
	}
	if len(services) != 1 {
		info(log, "expected only 1 ingress Service for DataPlane", dataplane)
		return ctrl.Result{Requeue: true}, nil
	}

	ingressService := services[0]
	if !dataPlaneIngressServiceIsReady(dataplane, &ingressService) {
		debug(log, "Ingress Service for DataPlane not ready yet", dataplane)

		// Set Ready to false for dataplane as the underlying deployment is not ready.
		k8sutils.SetCondition(
			k8sutils.NewConditionWithGeneration(
				k8sutils.ReadyType,
				metav1.ConditionFalse,
				k8sutils.WaitingToBecomeReadyReason,
				k8sutils.WaitingToBecomeReadyMessage,
				dataplane.Generation,
			),
			dataplane,
		)
		ensureDataPlaneReadinessStatus(dataplane, &deployment)
		if err = patchDataPlaneStatus(ctx, cl, log, dataplane); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed patching status (ingress Service not ready) for DataPlane %s/%s: %w", dataplane.Namespace, dataplane.Name, err)
		}
		return ctrl.Result{}, nil
	}

	markAsProvisioned(dataplane)
	k8sutils.SetReady(dataplane)
	ensureDataPlaneReadinessStatus(dataplane, &deployment)

	if err = patchDataPlaneStatus(ctx, cl, log, dataplane); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed patching status for DataPlane %s/%s: %w", dataplane.Namespace, dataplane.Name, err)
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// DataPlane - Private Functions - Equality Checks
// -----------------------------------------------------------------------------

func dataplaneSpecDeepEqual(spec1, spec2 *operatorv1beta1.DataPlaneOptions) bool {
	// TODO: Doesn't take .Rollout field into account.
	if !deploymentOptionsDeepEqual(&spec1.Deployment.DeploymentOptions, &spec2.Deployment.DeploymentOptions) ||
		!servicesOptionsDeepEqual(&spec1.Network, &spec2.Network) {
		return false
	}

	return true
}
