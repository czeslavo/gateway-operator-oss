package konnect

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	konnectv1alpha1 "github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

// KongKeySetReconciliationWatchOptions returns the watch options for the KongKeySet.
func KongKeySetReconciliationWatchOptions(cl client.Client) []func(*ctrl.Builder) *ctrl.Builder {
	return []func(*ctrl.Builder) *ctrl.Builder{
		func(b *ctrl.Builder) *ctrl.Builder {
			return b.For(&configurationv1alpha1.KongKeySet{},
				builder.WithPredicates(
					predicate.NewPredicateFuncs(objRefersToKonnectGatewayControlPlane[configurationv1alpha1.KongKeySet]),
				),
			)
		},
		func(b *ctrl.Builder) *ctrl.Builder {
			return b.Watches(
				&konnectv1alpha1.KonnectAPIAuthConfiguration{},
				handler.EnqueueRequestsFromMapFunc(
					enqueueObjectForAPIAuthThroughControlPlaneRef[configurationv1alpha1.KongKeySetList](
						cl, IndexFieldKongKeySetOnKonnectGatewayControlPlane,
					),
				),
			)
		},
		func(b *ctrl.Builder) *ctrl.Builder {
			return b.Watches(
				&konnectv1alpha1.KonnectGatewayControlPlane{},
				handler.EnqueueRequestsFromMapFunc(
					enqueueObjectForKonnectGatewayControlPlane[configurationv1alpha1.KongKeySetList](
						cl, IndexFieldKongKeySetOnKonnectGatewayControlPlane,
					),
				),
			)
		},
	}
}
