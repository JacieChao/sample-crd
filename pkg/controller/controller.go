package controller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/tools/record"

	sampleClientset "github.com/sample-crd/pkg/client/clientset/versioned"
	sampleLister "github.com/sample-crd/pkg/client/listers/samplecrd/v1"
	sampleInformer "github.com/sample-crd/pkg/client/informers/externalversions"

)

type SampleCRDController struct {
	kubeClient kubernetes.Interface
	sampleClient sampleClientset.Interface

	dplLister v1.DeploymentLister
	dplSynced cache.InformerSynced
	sampleLister sampleLister.SampleCRDLister
	sampleSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

func NewSampleController(
	kubeClient kubernetes.Interface,
	sampleClient sampleClientset.Interface,
	kubeInformerFactory informers.SharedInformerFactory,
	sampleInformerFactory sampleInformer.SharedInformerFactory) *SampleCRDController {
		c := &SampleCRDController{}

		return c
}