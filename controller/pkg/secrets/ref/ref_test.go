package ref

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/gateway-operator/api/v1alpha1"
	gwtypes "github.com/kong/gateway-operator/internal/types"
)

func TestIsSecretCrossReferenceGranted(t *testing.T) {
	customizeReferenceGrant := func(rg gatewayv1beta1.ReferenceGrant, opts ...func(rg *gatewayv1beta1.ReferenceGrant)) gatewayv1beta1.ReferenceGrant {
		rg = *rg.DeepCopy()
		for _, opt := range opts {
			opt(&rg)
		}
		return rg
	}

	const (
		badSecretName   = gwtypes.ObjectName("wrong-secret")
		emptySecretName = gwtypes.ObjectName("")
		goodSecretName  = gwtypes.ObjectName("good-secret")
	)
	referenceGrantForObj := func(obj client.Object) gatewayv1beta1.ReferenceGrant {
		return gatewayv1beta1.ReferenceGrant{
			Spec: gatewayv1beta1.ReferenceGrantSpec{
				From: []gatewayv1beta1.ReferenceGrantFrom{
					{
						Group:     gatewayv1.Group(obj.GetObjectKind().GroupVersionKind().Group),
						Kind:      gatewayv1.Kind(obj.GetObjectKind().GroupVersionKind().Kind),
						Namespace: gatewayv1.Namespace(obj.GetNamespace()),
					},
				},
				To: []gatewayv1beta1.ReferenceGrantTo{
					{
						Group: "",
						Kind:  "Secret",
						Name:  lo.ToPtr(goodSecretName),
					},
				},
			},
		}
	}
	var (
		objGateway = &gatewayv1.Gateway{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Gateway",
				APIVersion: gatewayv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "goodNamespace",
			},
		}
		objKPI = &v1alpha1.KongPluginInstallation{
			TypeMeta: metav1.TypeMeta{
				Kind:       "KongPluginInstallation",
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "goodNamespace",
			},
		}
	)

	testCases := []struct {
		name            string
		forObj          client.Object
		referenceGrants []gatewayv1beta1.ReferenceGrant
		isGranted       bool
	}{
		{
			name:      "no referenceGrants",
			forObj:    objGateway,
			isGranted: false,
		},
		{
			name:   "granted for Gateway",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				referenceGrantForObj(objGateway),
			},
			isGranted: true,
		},
		{
			name:   "granted for KPI",
			forObj: objKPI,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				referenceGrantForObj(objKPI),
			},
			isGranted: true,
		},
		{
			name:   "not granted, wrong Kind",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				referenceGrantForObj(objKPI),
			},
			isGranted: false,
		},
		{
			name:   "not granted, bad 'from' group",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objGateway), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.From[0].Group = "wrong-group"
				}),
			},
			isGranted: false,
		},
		{
			name:   "granted, for one Kind Gateway bad 'from' group, but it expects KPI that is properly granted",
			forObj: objKPI,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objGateway), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.From[0].Group = "wrong-group"
				}),
				referenceGrantForObj(objKPI),
			},
			isGranted: true,
		},
		{
			name:   "not granted, bad 'to' group",
			forObj: objKPI,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objKPI), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.To[0].Group = "wrong-group"
				}),
			},
			isGranted: false,
		},
		{
			name:   "not granted, bad 'from' kind",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objGateway), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.From[0].Kind = "wrong-kind"
				}),
			},
			isGranted: false,
		},
		{
			name:   "not granted, bad 'to' kind",
			forObj: objKPI,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objKPI), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.To[0].Kind = "wrong-kind"
				}),
			},
			isGranted: false,
		},
		{
			name:   "not granted, bad 'from' namespace",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objGateway), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.From[0].Namespace = "bad-namespace"
				}),
			},
			isGranted: false,
		},
		{
			name:   "not granted, empty 'to' secret name",
			forObj: objKPI,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objKPI), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.To[0].Name = lo.ToPtr(emptySecretName)
				}),
			},
			isGranted: false,
		},
		{
			name:   "not granted, bad 'to' secret name",
			forObj: objGateway,
			referenceGrants: []gatewayv1beta1.ReferenceGrant{
				customizeReferenceGrant(referenceGrantForObj(objGateway), func(rg *gatewayv1beta1.ReferenceGrant) {
					rg.Spec.To[0].Name = lo.ToPtr(badSecretName)
				}),
			},
			isGranted: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isGranted, isSecretCrossReferenceGranted(tc.forObj, goodSecretName, tc.referenceGrants))
		})
	}
}