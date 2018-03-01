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
package v1

import (
	v1 "github.com/sample-crd/pkg/apis/samplecrd/v1"
	scheme "github.com/sample-crd/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SampleCRDsGetter has a method to return a SampleCRDInterface.
// A group's client should implement this interface.
type SampleCRDsGetter interface {
	SampleCRDs(namespace string) SampleCRDInterface
}

// SampleCRDInterface has methods to work with SampleCRD resources.
type SampleCRDInterface interface {
	Create(*v1.SampleCRD) (*v1.SampleCRD, error)
	Update(*v1.SampleCRD) (*v1.SampleCRD, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.SampleCRD, error)
	List(opts meta_v1.ListOptions) (*v1.SampleCRDList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.SampleCRD, err error)
	SampleCRDExpansion
}

// sampleCRDs implements SampleCRDInterface
type sampleCRDs struct {
	client rest.Interface
	ns     string
}

// newSampleCRDs returns a SampleCRDs
func newSampleCRDs(c *SamplecrdV1Client, namespace string) *sampleCRDs {
	return &sampleCRDs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sampleCRD, and returns the corresponding sampleCRD object, and an error if there is any.
func (c *sampleCRDs) Get(name string, options meta_v1.GetOptions) (result *v1.SampleCRD, err error) {
	result = &v1.SampleCRD{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("samplecrds").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SampleCRDs that match those selectors.
func (c *sampleCRDs) List(opts meta_v1.ListOptions) (result *v1.SampleCRDList, err error) {
	result = &v1.SampleCRDList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("samplecrds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sampleCRDs.
func (c *sampleCRDs) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("samplecrds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a sampleCRD and creates it.  Returns the server's representation of the sampleCRD, and an error, if there is any.
func (c *sampleCRDs) Create(sampleCRD *v1.SampleCRD) (result *v1.SampleCRD, err error) {
	result = &v1.SampleCRD{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("samplecrds").
		Body(sampleCRD).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sampleCRD and updates it. Returns the server's representation of the sampleCRD, and an error, if there is any.
func (c *sampleCRDs) Update(sampleCRD *v1.SampleCRD) (result *v1.SampleCRD, err error) {
	result = &v1.SampleCRD{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("samplecrds").
		Name(sampleCRD.Name).
		Body(sampleCRD).
		Do().
		Into(result)
	return
}

// Delete takes name of the sampleCRD and deletes it. Returns an error if one occurs.
func (c *sampleCRDs) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("samplecrds").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sampleCRDs) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("samplecrds").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sampleCRD.
func (c *sampleCRDs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.SampleCRD, err error) {
	result = &v1.SampleCRD{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("samplecrds").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
