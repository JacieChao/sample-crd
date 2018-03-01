/*
Copyright 2018 The sample-crd Authors.

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
package fake

import (
	samplecrd_v1 "github.com/sample-crd/pkg/apis/samplecrd/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSampleCRDs implements SampleCRDInterface
type FakeSampleCRDs struct {
	Fake *FakeSamplecrdV1
	ns   string
}

var samplecrdsResource = schema.GroupVersionResource{Group: "samplecrd.k8s.io", Version: "v1", Resource: "samplecrds"}

var samplecrdsKind = schema.GroupVersionKind{Group: "samplecrd.k8s.io", Version: "v1", Kind: "SampleCRD"}

// Get takes name of the sampleCRD, and returns the corresponding sampleCRD object, and an error if there is any.
func (c *FakeSampleCRDs) Get(name string, options v1.GetOptions) (result *samplecrd_v1.SampleCRD, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(samplecrdsResource, c.ns, name), &samplecrd_v1.SampleCRD{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrd_v1.SampleCRD), err
}

// List takes label and field selectors, and returns the list of SampleCRDs that match those selectors.
func (c *FakeSampleCRDs) List(opts v1.ListOptions) (result *samplecrd_v1.SampleCRDList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(samplecrdsResource, samplecrdsKind, c.ns, opts), &samplecrd_v1.SampleCRDList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &samplecrd_v1.SampleCRDList{}
	for _, item := range obj.(*samplecrd_v1.SampleCRDList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sampleCRDs.
func (c *FakeSampleCRDs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(samplecrdsResource, c.ns, opts))

}

// Create takes the representation of a sampleCRD and creates it.  Returns the server's representation of the sampleCRD, and an error, if there is any.
func (c *FakeSampleCRDs) Create(sampleCRD *samplecrd_v1.SampleCRD) (result *samplecrd_v1.SampleCRD, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(samplecrdsResource, c.ns, sampleCRD), &samplecrd_v1.SampleCRD{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrd_v1.SampleCRD), err
}

// Update takes the representation of a sampleCRD and updates it. Returns the server's representation of the sampleCRD, and an error, if there is any.
func (c *FakeSampleCRDs) Update(sampleCRD *samplecrd_v1.SampleCRD) (result *samplecrd_v1.SampleCRD, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(samplecrdsResource, c.ns, sampleCRD), &samplecrd_v1.SampleCRD{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrd_v1.SampleCRD), err
}

// Delete takes name of the sampleCRD and deletes it. Returns an error if one occurs.
func (c *FakeSampleCRDs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(samplecrdsResource, c.ns, name), &samplecrd_v1.SampleCRD{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSampleCRDs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(samplecrdsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &samplecrd_v1.SampleCRDList{})
	return err
}

// Patch applies the patch and returns the patched sampleCRD.
func (c *FakeSampleCRDs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *samplecrd_v1.SampleCRD, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(samplecrdsResource, c.ns, name, data, subresources...), &samplecrd_v1.SampleCRD{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrd_v1.SampleCRD), err
}
