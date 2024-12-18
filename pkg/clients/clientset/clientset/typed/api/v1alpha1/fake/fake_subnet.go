/*
Copyright 2024 The ORC Authors.

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

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
	apiv1alpha1 "github.com/k-orc/openstack-resource-controller/pkg/clients/applyconfiguration/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSubnets implements SubnetInterface
type FakeSubnets struct {
	Fake *FakeOpenstackV1alpha1
	ns   string
}

var subnetsResource = v1alpha1.SchemeGroupVersion.WithResource("subnets")

var subnetsKind = v1alpha1.SchemeGroupVersion.WithKind("Subnet")

// Get takes name of the subnet, and returns the corresponding subnet object, and an error if there is any.
func (c *FakeSubnets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Subnet, err error) {
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(subnetsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// List takes label and field selectors, and returns the list of Subnets that match those selectors.
func (c *FakeSubnets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SubnetList, err error) {
	emptyResult := &v1alpha1.SubnetList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(subnetsResource, subnetsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SubnetList{ListMeta: obj.(*v1alpha1.SubnetList).ListMeta}
	for _, item := range obj.(*v1alpha1.SubnetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested subnets.
func (c *FakeSubnets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(subnetsResource, c.ns, opts))

}

// Create takes the representation of a subnet and creates it.  Returns the server's representation of the subnet, and an error, if there is any.
func (c *FakeSubnets) Create(ctx context.Context, subnet *v1alpha1.Subnet, opts v1.CreateOptions) (result *v1alpha1.Subnet, err error) {
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(subnetsResource, c.ns, subnet, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// Update takes the representation of a subnet and updates it. Returns the server's representation of the subnet, and an error, if there is any.
func (c *FakeSubnets) Update(ctx context.Context, subnet *v1alpha1.Subnet, opts v1.UpdateOptions) (result *v1alpha1.Subnet, err error) {
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(subnetsResource, c.ns, subnet, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSubnets) UpdateStatus(ctx context.Context, subnet *v1alpha1.Subnet, opts v1.UpdateOptions) (result *v1alpha1.Subnet, err error) {
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(subnetsResource, "status", c.ns, subnet, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// Delete takes name of the subnet and deletes it. Returns an error if one occurs.
func (c *FakeSubnets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(subnetsResource, c.ns, name, opts), &v1alpha1.Subnet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSubnets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(subnetsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SubnetList{})
	return err
}

// Patch applies the patch and returns the patched subnet.
func (c *FakeSubnets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Subnet, err error) {
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(subnetsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied subnet.
func (c *FakeSubnets) Apply(ctx context.Context, subnet *apiv1alpha1.SubnetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Subnet, err error) {
	if subnet == nil {
		return nil, fmt.Errorf("subnet provided to Apply must not be nil")
	}
	data, err := json.Marshal(subnet)
	if err != nil {
		return nil, err
	}
	name := subnet.Name
	if name == nil {
		return nil, fmt.Errorf("subnet.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(subnetsResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeSubnets) ApplyStatus(ctx context.Context, subnet *apiv1alpha1.SubnetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Subnet, err error) {
	if subnet == nil {
		return nil, fmt.Errorf("subnet provided to Apply must not be nil")
	}
	data, err := json.Marshal(subnet)
	if err != nil {
		return nil, err
	}
	name := subnet.Name
	if name == nil {
		return nil, fmt.Errorf("subnet.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Subnet{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(subnetsResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Subnet), err
}