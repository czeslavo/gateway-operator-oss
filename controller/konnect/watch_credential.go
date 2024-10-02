package konnect

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kong/gateway-operator/controller/konnect/constraints"
	operatorerrors "github.com/kong/gateway-operator/internal/errors"

	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
)

// kongCredentialRefersToKonnectGatewayControlPlane returns a predicate function that
// returns true if the KongCredential refers to a KongConsumer which uses
// KonnectGatewayControlPlane reference.
func kongCredentialRefersToKonnectGatewayControlPlane[
	T interface {
		*configurationv1alpha1.KongCredentialACL |
			*configurationv1alpha1.KongCredentialAPIKey |
			*configurationv1alpha1.KongCredentialBasicAuth |
			*configurationv1alpha1.KongCredentialJWT
		// TODO add support for HMAC Auth https://github.com/Kong/gateway-operator/issues/621

		GetTypeName() string
		GetNamespace() string
	},
](cl client.Client) func(obj client.Object) bool {
	return func(obj client.Object) bool {
		credential, ok := obj.(T)
		if !ok {
			ctrllog.FromContext(context.Background()).Error(
				operatorerrors.ErrUnexpectedObject,
				"failed to run predicate function",
				"expected", constraints.EntityTypeName[T](), "found", reflect.TypeOf(obj),
			)
			return false
		}

		var consumerRefName string
		switch credential := any(credential).(type) {
		case *configurationv1alpha1.KongCredentialACL:
			consumerRefName = credential.Spec.ConsumerRef.Name
		case *configurationv1alpha1.KongCredentialAPIKey:
			consumerRefName = credential.Spec.ConsumerRef.Name
		case *configurationv1alpha1.KongCredentialBasicAuth:
			consumerRefName = credential.Spec.ConsumerRef.Name
		case *configurationv1alpha1.KongCredentialJWT:
			consumerRefName = credential.Spec.ConsumerRef.Name
		}

		nn := types.NamespacedName{
			Namespace: credential.GetNamespace(),
			Name:      consumerRefName,
		}
		var consumer configurationv1.KongConsumer
		if err := cl.Get(context.Background(), nn, &consumer); client.IgnoreNotFound(err) != nil {
			return true
		}

		return objHasControlPlaneRefKonnectNamespacedRef(&consumer)
	}
}
