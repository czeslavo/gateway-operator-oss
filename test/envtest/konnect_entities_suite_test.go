package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/gateway-operator/controller/konnect"
	"github.com/kong/gateway-operator/controller/konnect/constraints"
	"github.com/kong/gateway-operator/controller/konnect/ops"
	"github.com/kong/gateway-operator/modules/manager/scheme"

	konnectv1alpha1 "github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

// TestKonnectEntityReconcilers tests Konnect entity reconcilers. The test cases are run against a real Kubernetes API
// server provided by the envtest package and a mock Konnect SDK.
func TestKonnectEntityReconcilers(t *testing.T) {
	cfg, _ := Setup(t, context.Background(), scheme.Get())

	testNewKonnectEntityReconciler(t, cfg, konnectv1alpha1.KonnectGatewayControlPlane{}, konnectGatewayControlPlaneTestCases)
}

type konnectEntityReconcilerTestCase struct {
	name                string
	objectOps           func(ctx context.Context, t *testing.T, cl client.Client, ns *corev1.Namespace)
	mockExpectations    func(t *testing.T, sdk *ops.MockSDKWrapper, cl client.Client, ns *corev1.Namespace)
	eventuallyPredicate func(ctx context.Context, t *assert.CollectT, cl client.Client, ns *corev1.Namespace)
}

// testNewKonnectEntityReconciler is a helper function to test Konnect entity reconcilers.
// It creates a new namespace for each test case and starts a new controller manager.
// The provided rest.Config designates the Kubernetes API server to use for the tests.
func testNewKonnectEntityReconciler[
	T constraints.SupportedKonnectEntityType,
	TEnt constraints.EntityType[T],
](
	t *testing.T,
	cfg *rest.Config,
	ent T,
	testCases []konnectEntityReconcilerTestCase,
) {
	t.Helper()

	t.Run(ent.GetTypeName(), func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mgr, logs := NewManager(t, ctx, cfg, scheme.Get())

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: NameFromT(t),
			},
		}
		require.NoError(t, mgr.GetClient().Create(ctx, ns))

		cl := client.NewNamespacedClient(mgr.GetClient(), ns.Name)
		factory := ops.NewMockSDKFactory(t)
		sdk := factory.SDK

		StartReconcilers(ctx, t, mgr, logs, konnect.NewKonnectEntityReconciler[T, TEnt](factory, false, cl))

		const (
			wait = time.Second
			tick = 20 * time.Millisecond
		)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockExpectations(t, sdk, cl, ns)
				tc.objectOps(ctx, t, cl, ns)
				require.EventuallyWithT(t, func(collect *assert.CollectT) {
					tc.eventuallyPredicate(ctx, collect, cl, ns)
				}, wait, tick)
			})
		}
	})
}
