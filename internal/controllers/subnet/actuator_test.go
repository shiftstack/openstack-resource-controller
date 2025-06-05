package subnet

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
)

func TestHandleNameUpdate(t *testing.T) {
	testCases := []struct {
		name         string
		new          v1alpha1.OpenStackName
		existing     string
		expectChange bool
	}{
		{name: "Identical", new: v1alpha1.OpenStackName("name"), existing: "name", expectChange: false},
		{name: "Identical to object name", new: "", existing: "object-name", expectChange: false},
		{name: "Different", new: "new-name", existing: "name", expectChange: true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resource := &orcv1alpha1.Subnet{}
			resource.Name = "object-name"
			resource.Spec = orcv1alpha1.SubnetSpec{
				Resource: &orcv1alpha1.SubnetResourceSpec{Name: &tt.new},
			}
			osResource := &subnets.Subnet{Name: tt.existing}

			updateOpts := subnets.UpdateOpts{}
			handleNameUpdate(&updateOpts, resource, osResource)

			got := !reflect.ValueOf(updateOpts).IsZero()
			if got != tt.expectChange {
				t.Errorf("Expected change: %v, got %v", tt.expectChange, got)
			}
		})

	}
}

func TestHandleDNSNameserversUpdate(t *testing.T) {
	testCases := []struct {
		name         string
		new          []v1alpha1.IPvAny
		existing     []string
		expectChange bool
	}{
		{name: "Duplicate", new: []v1alpha1.IPvAny{"one", "two", "two"}, existing: []string{"one", "two", "three"}, expectChange: true},
		{name: "Different order", new: []v1alpha1.IPvAny{"one", "two"}, existing: []string{"two", "one"}, expectChange: true},
		{name: "Identical", new: []v1alpha1.IPvAny{"one", "two", "three"}, existing: []string{"one", "two", "three"}, expectChange: false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resource := &orcv1alpha1.SubnetResourceSpec{DNSNameservers: tt.new}
			osResource := &subnets.Subnet{DNSNameservers: tt.existing}

			updateOpts := subnets.UpdateOpts{}
			handleDNSNameserversUpdate(&updateOpts, resource, osResource)

			got := !reflect.ValueOf(updateOpts).IsZero()
			if got != tt.expectChange {
				t.Errorf("Expected change: %v, got %v", tt.expectChange, got)
			}
		})

	}
}
