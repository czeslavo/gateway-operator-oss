package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/gateway-operator/api/v1alpha1"
	"github.com/kong/gateway-operator/api/v1beta1"
	"github.com/kong/gateway-operator/pkg/consts"
	gatewayutils "github.com/kong/gateway-operator/pkg/utils/gateway"
	testutils "github.com/kong/gateway-operator/pkg/utils/test"
	"github.com/kong/gateway-operator/test/helpers"
)

func TestAIGatewayCreation(t *testing.T) {
	t.Parallel()

	namespace, cleaner := helpers.SetupTestEnv(t, GetCtx(), GetEnv())

	t.Log("deploying a GatewayConfiguration resource")
	gatewayConfiguration := &v1beta1.GatewayConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      uuid.New().String(),
			Namespace: namespace.Name,
		},
		Spec: v1beta1.GatewayConfigurationSpec{
			DataPlaneOptions: &v1beta1.GatewayConfigDataPlaneOptions{
				Deployment: v1beta1.DataPlaneDeploymentOptions{
					DeploymentOptions: v1beta1.DeploymentOptions{
						PodTemplateSpec: &corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Name:            consts.DataPlaneProxyContainerName,
									Image:           helpers.GetDefaultDataPlaneImage(),
									ImagePullPolicy: corev1.PullAlways,
									ReadinessProbe: &corev1.Probe{
										InitialDelaySeconds: 1,
										PeriodSeconds:       1,
									},
									Env: []corev1.EnvVar{
										{
											Name:  "KONG_ADMIN_GUI_LISTEN",
											Value: "0.0.0.0:8002",
										},
										{
											Name:  "KONG_ADMIN_LISTEN",
											Value: "0.0.0.0:8001, 0.0.0.0:8444 ssl reuseport backlog=16384",
										},
									},
								}},
							},
						},
					},
				},
			},
			ControlPlaneOptions: &v1beta1.ControlPlaneOptions{
				Deployment: v1beta1.ControlPlaneDeploymentOptions{
					PodTemplateSpec: &corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  consts.ControlPlaneControllerContainerName,
								Image: consts.DefaultControlPlaneImage,
								Env: []corev1.EnvVar{{
									Name:  "CONTROLLER_LOG_LEVEL",
									Value: "debug",
								}},
							}},
						},
					},
				},
			},
		},
	}
	gatewayConfiguration, err := GetClients().OperatorClient.ApisV1beta1().GatewayConfigurations(namespace.Name).Create(GetCtx(), gatewayConfiguration, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(gatewayConfiguration)

	t.Log("deploying a GatewayClass resource, [", &gatewayConfiguration.Name, "]")
	gatewayClass := testutils.GenerateGatewayClass()
	gatewayClass.Spec.ParametersRef = &gatewayv1.ParametersReference{
		Group:     gatewayv1.Group("gateway-operator.konghq.com"),
		Kind:      gatewayv1.Kind("GatewayConfiguration"),
		Name:      gatewayConfiguration.Name,
		Namespace: (*gatewayv1.Namespace)(&namespace.Name),
	}
	gatewayClass, err = GetClients().GatewayClient.GatewayV1().GatewayClasses().Create(GetCtx(), gatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(gatewayClass)

	credSecretName := uuid.New().String()
	t.Log("creating null secret containing the required credentials, [", credSecretName, "]")
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credSecretName,
			Namespace: namespace.Name,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			// TODO get real credentials from a vault...
			//
			// See: https://github.com/Kong/gateway-operator/issues/1368
			"openai": []byte("openai-key"),
			"cohere": []byte("cohere-key"),
		},
	}
	credSecret, err = GetClients().K8sClient.CoreV1().Secrets(namespace.Name).Create(GetCtx(), credSecret, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(credSecret)

	modelOpenAI := "gpt-3.5-turbo-instruct"
	modelCohere := "command"
	promptTypeCompletions := "completions"
	promptTypeChat := "chat"
	maxTokens := 256
	identifierOpenAI := "devteam-chatgpt"
	identifierCohere := "cohere-command"

	aigatewayName := "aigateway-test"
	t.Log("deploying an AIGateway, [", aigatewayName, "]")
	aigateway := &v1alpha1.AIGateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace.Name,
			Name:      aigatewayName,
		},
		Spec: v1alpha1.AIGatewaySpec{
			GatewayClassName: gatewayClass.Name,
			LargeLanguageModels: &v1alpha1.LargeLanguageModels{
				CloudHosted: []v1alpha1.CloudHostedLargeLanguageModel{
					{
						Identifier: identifierOpenAI,
						Model:      &modelOpenAI,
						PromptType: (*v1alpha1.LLMPromptType)(&promptTypeChat),
						AICloudProvider: v1alpha1.AICloudProvider{
							Name: v1alpha1.AICloudProviderOpenAI,
						},
						DefaultPromptParams: &v1alpha1.LLMPromptParams{
							MaxTokens: &maxTokens,
						},
					},
					{
						Identifier: identifierCohere,
						Model:      &modelCohere,
						PromptType: (*v1alpha1.LLMPromptType)(&promptTypeCompletions),
						AICloudProvider: v1alpha1.AICloudProvider{
							Name: v1alpha1.AICloudProviderCohere,
						},
						// Deliberately no Default Params, to test nil checks
					},
				},
			},
			CloudProviderCredentials: &v1alpha1.AICloudProviderAPITokenRef{
				Name:      credSecretName,
				Namespace: &namespace.Name,
			},
		},
	}
	aigateway, err = GetClients().OperatorClient.ApisV1alpha1().AIGateways(namespace.Name).Create(GetCtx(), aigateway, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(aigateway)

	t.Log("checking for the Gateway that should have been created for the AIGateway")
	gateway := eventuallyDetermineGatewayForAIGateway(t, aigateway, GetClients())

	t.Log("verifying Gateway gets marked as Scheduled")
	gatewayExpectedNN := types.NamespacedName{Name: gateway.Name, Namespace: gateway.Namespace}
	require.Eventually(t, testutils.GatewayIsScheduled(t, GetCtx(), gatewayExpectedNN, clients), testutils.GatewaySchedulingTimeLimit, time.Second)

	t.Log("verifying Gateway gets marked as Programmed")
	require.Eventually(t, testutils.GatewayIsProgrammed(t, GetCtx(), gatewayExpectedNN, clients), testutils.GatewayReadyTimeLimit, time.Second)
	require.Eventually(t, testutils.GatewayListenersAreProgrammed(t, GetCtx(), gatewayExpectedNN, clients), testutils.GatewayReadyTimeLimit, time.Second)

	t.Log("verifying Gateway gets an IP address")
	require.Eventually(t, testutils.GatewayIPAddressExist(t, GetCtx(), gatewayExpectedNN, clients), testutils.SubresourceReadinessWait, time.Second)
	gateway = testutils.MustGetGateway(t, GetCtx(), gatewayExpectedNN, clients)
	gatewayIPAddress := gateway.Status.Addresses[0].Value

	t.Log("verifying that the DataPlane becomes Ready")
	require.Eventually(t, testutils.GatewayDataPlaneIsReady(t, GetCtx(), gateway, clients), testutils.SubresourceReadinessWait, time.Second)
	dataplanes := testutils.MustListDataPlanesForGateway(t, GetCtx(), gateway, clients)
	require.Len(t, dataplanes, 1)
	dataplane := dataplanes[0]

	t.Log("verifying that the ControlPlane becomes provisioned")
	require.Eventually(t, testutils.GatewayControlPlaneIsProvisioned(t, GetCtx(), gateway, clients), testutils.SubresourceReadinessWait, time.Second)
	controlplanes := testutils.MustListControlPlanesForGateway(t, GetCtx(), gateway, clients)
	require.Len(t, controlplanes, 1)
	controlplane := controlplanes[0]

	t.Log("verifying networkpolicies are created")
	require.Eventually(t, testutils.GatewayNetworkPoliciesExist(t, GetCtx(), gateway, clients), testutils.SubresourceReadinessWait, time.Second)

	t.Log("verifying connectivity to the Gateway")
	require.Eventually(t, Expect404WithNoRouteFunc(t, GetCtx(), "http://"+gatewayIPAddress), testutils.SubresourceReadinessWait, time.Second)

	dataplaneNN := types.NamespacedName{Namespace: namespace.Name, Name: dataplane.Name}
	controlplaneNN := types.NamespacedName{Namespace: namespace.Name, Name: controlplane.Name}

	t.Log("verifying that dataplane has 1 ready replica")
	require.Eventually(t, testutils.DataPlaneHasNReadyPods(t, GetCtx(), dataplaneNN, clients, 1), time.Minute, time.Second)

	t.Log("verifying that controlplane has 1 ready replica")
	require.Eventually(t, testutils.ControlPlaneHasNReadyPods(t, GetCtx(), controlplaneNN, clients, 1), time.Minute, time.Second)

	t.Log("verifying that the HTTPRoute is now available for these LLMs")
	require.Eventually(t, func() bool {
		gateway, err = GetClients().GatewayClient.GatewayV1().Gateways(namespace.Name).Get(GetCtx(), gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return gatewayutils.IsScheduled(gateway)
	}, testutils.GatewaySchedulingTimeLimit, time.Second)

	// TODO - for now we don't have AI Cloud provider credentials in CI,
	// this is something we're considering adding for later but it has
	// cost implications we need to work through. Rather than have this
	// test hit the real cloud provider, we test manually for now after
	// at least verifying that all the resources are in place.
	//
	// See:  https://github.com/Kong/gateway-operator/issues/1368
}

// This is temporary while we work on statuses for AIGateways.
//
// See: https://github.com/Kong/gateway-operator/issues/1368
func eventuallyDetermineGatewayForAIGateway(
	t *testing.T,
	aigateway *v1alpha1.AIGateway,
	clients testutils.K8sClients,
) (gateway *gatewayv1.Gateway) {
	require.Eventually(t, func() bool {
		gateways, err := clients.GatewayClient.GatewayV1().Gateways(aigateway.Namespace).List(GetCtx(), metav1.ListOptions{})
		require.NoError(t, err)

		for _, item := range gateways.Items {
			gw := item
			for _, ownerRef := range gw.OwnerReferences {
				if ownerRef.UID == aigateway.UID {
					gateway = &gw
					return true
				}
			}
		}

		return false
	}, time.Minute, time.Second)
	return
}
