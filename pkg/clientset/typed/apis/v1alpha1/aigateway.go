/*
Copyright 2022 Kong Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	apisv1alpha1 "github.com/kong/gateway-operator/api/v1alpha1"
	scheme "github.com/kong/gateway-operator/pkg/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// AIGatewaysGetter has a method to return a AIGatewayInterface.
// A group's client should implement this interface.
type AIGatewaysGetter interface {
	AIGateways(namespace string) AIGatewayInterface
}

// AIGatewayInterface has methods to work with AIGateway resources.
type AIGatewayInterface interface {
	Create(ctx context.Context, aIGateway *apisv1alpha1.AIGateway, opts v1.CreateOptions) (*apisv1alpha1.AIGateway, error)
	Update(ctx context.Context, aIGateway *apisv1alpha1.AIGateway, opts v1.UpdateOptions) (*apisv1alpha1.AIGateway, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, aIGateway *apisv1alpha1.AIGateway, opts v1.UpdateOptions) (*apisv1alpha1.AIGateway, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*apisv1alpha1.AIGateway, error)
	List(ctx context.Context, opts v1.ListOptions) (*apisv1alpha1.AIGatewayList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *apisv1alpha1.AIGateway, err error)
	AIGatewayExpansion
}

// aIGateways implements AIGatewayInterface
type aIGateways struct {
	*gentype.ClientWithList[*apisv1alpha1.AIGateway, *apisv1alpha1.AIGatewayList]
}

// newAIGateways returns a AIGateways
func newAIGateways(c *ApisV1alpha1Client, namespace string) *aIGateways {
	return &aIGateways{
		gentype.NewClientWithList[*apisv1alpha1.AIGateway, *apisv1alpha1.AIGatewayList](
			"aigateways",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *apisv1alpha1.AIGateway { return &apisv1alpha1.AIGateway{} },
			func() *apisv1alpha1.AIGatewayList { return &apisv1alpha1.AIGatewayList{} },
		),
	}
}
