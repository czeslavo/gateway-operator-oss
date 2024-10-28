package konnect

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/gateway-operator/controller/konnect/constraints"

	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	konnectv1alpha1 "github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

type handleKeySetRefTestCase[T constraints.SupportedKonnectEntityType, TEnt constraints.EntityType[T]] struct {
	name                string
	ent                 TEnt
	objects             []client.Object
	expectResult        ctrl.Result
	expectErrorContains string
	// Returns true if the updated entity satisfy the assertion.
	// Returns false and error message if entity fails to satisfy it.
	updatedEntAssertions []func(TEnt) (ok bool, message string)
}

func TestHandleKeySetRef(t *testing.T) {
	// Test objects definitions.
	var (
		commonKeyMeta = metav1.ObjectMeta{
			Name:      "key-1",
			Namespace: "ns",
		}

		cp = &konnectv1alpha1.KonnectGatewayControlPlane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cp-1",
				Namespace: "ns",
			},
		}
		cpRef = &configurationv1alpha1.ControlPlaneRef{
			Type: configurationv1alpha1.ControlPlaneRefKonnectNamespacedRef,
			KonnectNamespacedRef: &configurationv1alpha1.KonnectNamespacedRef{
				Name: "cp-1",
			},
		}

		notProgrammedKeySet = &configurationv1alpha1.KongKeySet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "key-set-1",
				Namespace: "ns",
			},
			Spec: configurationv1alpha1.KongKeySetSpec{
				ControlPlaneRef: cpRef,
				KongKeySetAPISpec: configurationv1alpha1.KongKeySetAPISpec{
					Name: "key-set-1",
				},
			},
		}
		programmedKeySet = &configurationv1alpha1.KongKeySet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "key-set-2",
				Namespace: "ns",
			},
			Spec: configurationv1alpha1.KongKeySetSpec{
				ControlPlaneRef: cpRef,
				KongKeySetAPISpec: configurationv1alpha1.KongKeySetAPISpec{
					Name: "key-set-2",
				},
			},
			Status: configurationv1alpha1.KongKeySetStatus{
				Konnect: &konnectv1alpha1.KonnectEntityStatusWithControlPlaneRef{
					KonnectEntityStatus: konnectv1alpha1.KonnectEntityStatus{
						ID: "key-set-id",
					},
					ControlPlaneID: "cp-id",
				},
				Conditions: []metav1.Condition{
					{
						Type:               konnectv1alpha1.KonnectEntityProgrammedConditionType,
						Status:             metav1.ConditionTrue,
						ObservedGeneration: 1,
						LastTransitionTime: metav1.Now(),
					},
				},
			},
		}
		keySetDuringDeletion = &configurationv1alpha1.KongKeySet{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "key-set-3",
				Namespace:         "ns",
				DeletionTimestamp: lo.ToPtr(metav1.Now()),
				Finalizers: []string{
					KonnectCleanupFinalizer,
				},
			},
		}
	)

	// Common assertions.
	var (
		keySetIDIsEmpty = func(key *configurationv1alpha1.KongKey) (bool, string) {
			if key.Status.Konnect != nil && key.Status.Konnect.KeySetID != "" {
				return false, "KeySetID should be empty"
			}
			return true, ""
		}
		keySetIDIs = func(expectedID string) func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
			return func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
				if key.Status.Konnect == nil || key.Status.Konnect.KeySetID != expectedID {
					return false, fmt.Sprintf("KeySetID should be %s", expectedID)
				}
				return true, ""
			}
		}
		keySetRefConditionIs = func(expectedStatus metav1.ConditionStatus) func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
			return func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
				containsCondition := lo.ContainsBy(key.Status.Conditions, func(condition metav1.Condition) bool {
					return condition.Type == konnectv1alpha1.KeySetRefValidConditionType &&
						condition.Status == expectedStatus
				})
				if !containsCondition {
					return false, fmt.Sprintf("KeySetRefValid condition should be %s", expectedStatus)
				}
				return true, ""
			}
		}
		hasExactlyOneOwnerReferenceTo = func(obj client.Object) func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
			return func(key *configurationv1alpha1.KongKey) (ok bool, message string) {
				if len(key.GetOwnerReferences()) != 1 {
					return false, "KongKey should have exactly one owner reference"
				}

				hasOwnerRef := lo.ContainsBy(key.GetOwnerReferences(), func(owner metav1.OwnerReference) bool {
					return owner.Name == obj.GetName() &&
						reflect.TypeOf(obj).Elem().Name() == owner.Kind &&
						owner.BlockOwnerDeletion != nil && *owner.BlockOwnerDeletion
				})
				if !hasOwnerRef {
					return false, fmt.Sprintf("KongKey should have owner reference to %s", client.ObjectKeyFromObject(obj))
				}
				return true, ""
			}
		}
	)

	testCases := []handleKeySetRefTestCase[configurationv1alpha1.KongKey, *configurationv1alpha1.KongKey]{
		{
			name: "key set ref is nil",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef:       nil,
				},
			},
			expectResult: ctrl.Result{},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIsEmpty,
			},
		},
		{
			name: "key set ref is nil but Konnect ID in status is set",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef:       nil,
				},
				Status: configurationv1alpha1.KongKeyStatus{
					Konnect: &konnectv1alpha1.KonnectEntityStatusWithControlPlaneAndKeySetRef{
						ControlPlaneID: "cp-id",
					},
				},
			},
			expectResult: ctrl.Result{},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIsEmpty,
			},
		},
		{
			name: "key set ref points to non-existing key set",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef: &configurationv1alpha1.KeySetRef{
						Type: configurationv1alpha1.KeySetRefNamespacedRef,
						NamespacedRef: &configurationv1alpha1.KeySetNamespacedRef{
							Name: "key-set-1",
						},
					},
				},
			},
			expectResult:        ctrl.Result{},
			expectErrorContains: "keysets.configuration.konghq.com \"key-set-1\" not found",
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetRefConditionIs(metav1.ConditionFalse),
				keySetIDIsEmpty,
			},
		},
		{
			name: "key set ref points to a key set during deletion",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef: &configurationv1alpha1.KeySetRef{
						Type: configurationv1alpha1.KeySetRefNamespacedRef,
						NamespacedRef: &configurationv1alpha1.KeySetNamespacedRef{
							Name: keySetDuringDeletion.Name,
						},
					},
				},
			},
			objects:             []client.Object{keySetDuringDeletion},
			expectResult:        ctrl.Result{},
			expectErrorContains: "referenced KongKeySet ns/key-set-3 is being deleted",
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIsEmpty,
			},
		},
		{
			name: "key set ref points to a key set that is not programmed yet",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef: &configurationv1alpha1.KeySetRef{
						Type: configurationv1alpha1.KeySetRefNamespacedRef,
						NamespacedRef: &configurationv1alpha1.KeySetNamespacedRef{
							Name: notProgrammedKeySet.Name,
						},
					},
				},
			},
			objects:      []client.Object{notProgrammedKeySet},
			expectResult: ctrl.Result{Requeue: true},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIsEmpty,
				keySetRefConditionIs(metav1.ConditionFalse),
			},
		},
		{
			name: "key set ref points to a programmed key set",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef: &configurationv1alpha1.KeySetRef{
						Type: configurationv1alpha1.KeySetRefNamespacedRef,
						NamespacedRef: &configurationv1alpha1.KeySetNamespacedRef{
							Name: programmedKeySet.Name,
						},
					},
				},
			},
			objects:      []client.Object{programmedKeySet},
			expectResult: ctrl.Result{},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetRefConditionIs(metav1.ConditionTrue),
				keySetIDIs(programmedKeySet.Status.Konnect.ID),
				hasExactlyOneOwnerReferenceTo(programmedKeySet),
			},
		},
		{
			name: "key set ref in spec changed to nil after resolving ref",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: commonKeyMeta,
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef:       nil,
				},
				Status: configurationv1alpha1.KongKeyStatus{
					Konnect: &konnectv1alpha1.KonnectEntityStatusWithControlPlaneAndKeySetRef{
						ControlPlaneID: "cp-id",
						KeySetID:       "key-set-id",
					},
				},
			},
			expectResult: ctrl.Result{},
			objects:      []client.Object{cp},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIsEmpty,
				keySetRefConditionIs(metav1.ConditionTrue),
				hasExactlyOneOwnerReferenceTo(&konnectv1alpha1.KonnectGatewayControlPlane{
					ObjectMeta: metav1.ObjectMeta{
						Name: cpRef.KonnectNamespacedRef.Name,
					},
				}),
			},
		},
		{
			name: "when entity has owning reference to control plane and refers key set, the ownership should be transferred to key set",
			ent: &configurationv1alpha1.KongKey{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "key-1",
					Namespace: "ns",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: konnectv1alpha1.GroupVersion.String(),
							Kind:       "KonnectGatewayControlPlane",
							Name:       cpRef.KonnectNamespacedRef.Name,
						},
					},
				},
				Spec: configurationv1alpha1.KongKeySpec{
					ControlPlaneRef: cpRef,
					KeySetRef: &configurationv1alpha1.KeySetRef{
						Type: configurationv1alpha1.KeySetRefNamespacedRef,
						NamespacedRef: &configurationv1alpha1.KeySetNamespacedRef{
							Name: programmedKeySet.Name,
						},
					},
				},
			},
			objects:      []client.Object{programmedKeySet},
			expectResult: ctrl.Result{},
			updatedEntAssertions: []func(*configurationv1alpha1.KongKey) (ok bool, message string){
				keySetIDIs(programmedKeySet.Status.Konnect.ID),
				keySetRefConditionIs(metav1.ConditionTrue),
				hasExactlyOneOwnerReferenceTo(programmedKeySet),
			},
		},
	}
	testHandleKeySetRef(t, testCases)
}

func testHandleKeySetRef[T constraints.SupportedKonnectEntityType, TEnt constraints.EntityType[T]](
	t *testing.T, testCases []handleKeySetRefTestCase[T, TEnt],
) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, configurationv1alpha1.AddToScheme(scheme))
			require.NoError(t, konnectv1alpha1.AddToScheme(scheme))
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.ent).
				WithObjects(tc.objects...).
				// WithStatusSubresource is required for updating status of handled entity.
				WithStatusSubresource(tc.ent).
				Build()
			require.NoError(t, fakeClient.Status().Update(context.Background(), tc.ent))

			res, err := handleKongKeySetRef(context.Background(), fakeClient, tc.ent)

			updatedEnt := tc.ent.DeepCopyObject().(TEnt)
			require.NoError(t, fakeClient.Get(context.Background(), client.ObjectKeyFromObject(tc.ent), updatedEnt))
			for _, assertion := range tc.updatedEntAssertions {
				ok, msg := assertion(updatedEnt)
				require.True(t, ok, msg)
			}

			if len(tc.expectErrorContains) > 0 {
				require.ErrorContains(t, err, tc.expectErrorContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectResult, res)
		})
	}
}