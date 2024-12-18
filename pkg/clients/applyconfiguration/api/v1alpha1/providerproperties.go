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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
)

// ProviderPropertiesApplyConfiguration represents a declarative configuration of the ProviderProperties type for use
// with apply.
type ProviderPropertiesApplyConfiguration struct {
	NetworkType     *v1alpha1.ProviderNetworkType `json:"networkType,omitempty"`
	PhysicalNetwork *v1alpha1.PhysicalNetwork     `json:"physicalNetwork,omitempty"`
	SegmentationID  *int32                        `json:"segmentationID,omitempty"`
}

// ProviderPropertiesApplyConfiguration constructs a declarative configuration of the ProviderProperties type for use with
// apply.
func ProviderProperties() *ProviderPropertiesApplyConfiguration {
	return &ProviderPropertiesApplyConfiguration{}
}

// WithNetworkType sets the NetworkType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NetworkType field is set to the value of the last call.
func (b *ProviderPropertiesApplyConfiguration) WithNetworkType(value v1alpha1.ProviderNetworkType) *ProviderPropertiesApplyConfiguration {
	b.NetworkType = &value
	return b
}

// WithPhysicalNetwork sets the PhysicalNetwork field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PhysicalNetwork field is set to the value of the last call.
func (b *ProviderPropertiesApplyConfiguration) WithPhysicalNetwork(value v1alpha1.PhysicalNetwork) *ProviderPropertiesApplyConfiguration {
	b.PhysicalNetwork = &value
	return b
}

// WithSegmentationID sets the SegmentationID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SegmentationID field is set to the value of the last call.
func (b *ProviderPropertiesApplyConfiguration) WithSegmentationID(value int32) *ProviderPropertiesApplyConfiguration {
	b.SegmentationID = &value
	return b
}
