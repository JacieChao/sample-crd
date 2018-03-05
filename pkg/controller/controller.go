package controller

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/kubernetes/scheme"
	typecorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/golang/glog"
	sampleClientset "github.com/sample-crd/pkg/client/clientset/versioned"
	sampleLister "github.com/sample-crd/pkg/client/listers/samplecrd/v1"
	sampleInformer "github.com/sample-crd/pkg/client/informers/externalversions"
	sampleschema "github.com/sample-crd/pkg/client/clientset/versioned/scheme"
	samplev1 "github.com/sample-crd/pkg/apis/samplecrd/v1"
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

	deployInformer := kubeInformerFactory.Apps().V1().Deployments()
	sampleInformer := sampleInformerFactory.Samplecrd().V1().SampleCRDs()
	sampleschema.AddToScheme(scheme.Scheme)

	glog.Info("--------- generate controller -----------")
	bct := record.NewBroadcaster()
	bct.StartLogging(glog.Infof)
	bct.StartRecordingToSink(&typecorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	recorder := bct.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "sample-crd"})

	c := &SampleCRDController{
		kubeClient: kubeClient,
		sampleClient: sampleClient,
		dplLister: deployInformer.Lister(),
		dplSynced: deployInformer.Informer().HasSynced(),
		sampleLister: sampleInformer.Lister(),
		sampleSynced: sampleInformer.Informer().HasSynced(),
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "SampleCRDs"),
		recorder: recorder,
	}

	glog.Info("--------- set handlers -----------")
	sampleInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.sample,
		UpdateFunc: func(old, new interface{}) {
			c.sample(new)
		},
	})

	deployInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.deploy,
		UpdateFunc: func(old, new interface{}) {
			newDpl := new.(*appsv1.Deployment)
			oldDpl := old.(*appsv1.Deployment)
			if newDpl.ResourceVersion == oldDpl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			c.deploy(new)
		},
	})

	return c
}

// sample handler
func (c *SampleCRDController) sample(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}


// deployment handler
func (c *SampleCRDController) deploy(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "SampleCRD" {
			return
		}

		sample, err := c.sampleLister.SampleCRDs(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.Infof("ignoring orphaned object '%s' of samplecrd '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.sample(sample)
		return
	}
}

// controller run
func (c *SampleCRDController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("--------------------- Starting SampleCRD controller ---------------------------")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.dplSynced, c.sampleSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("----------------------- Starting workers -------------------------------")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("---------------------- workers has been started ------------------------")
	<-stopCh
	glog.Info("---------------------- Shutting down workers ----------------------")

	return nil
}

// process workqueue
func (c *SampleCRDController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// read a item from workqueue and handle it
func (c *SampleCRDController) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// sync compare desire and actual
func (c *SampleCRDController) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the resource with this namespace/name
	sample, err := c.sampleLister.SampleCRDs(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("SampleCRD '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	deploymentName := sample.Spec.DeploymentName
	if deploymentName == "" {
		runtime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment name
	deployment, err := c.dplLister.Deployments(sample.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeClient.AppsV1().Deployments(sample.Namespace).Create(newDeployment(sample))
	}

	if err != nil {
		return err
	}

	// If the deployment has already exits, return an error
	if !metav1.IsControlledBy(deployment, sample) {
		msg := fmt.Sprintf("%s has already exits", deployment.Name)
		c.recorder.Event(sample, corev1.EventTypeWarning, "create deployment fail", msg)
		return fmt.Errorf(msg)
	}

	//update deployment replicas
	if sample.Spec.Replicas != nil && *sample.Spec.Replicas != *deployment.Spec.Replicas {
		glog.Infof("update %s deployments replicas to %d", sample.Spec.DeploymentName, sample.Spec.Replicas)
		deployment, err = c.kubeClient.AppsV1().Deployments(sample.Namespace).Update(newDeployment(sample))
	}

	if err != nil {
		return err
	}

	err = c.updateStatus(sample, deployment)
	if err != nil {
		return err
	}

	// get pods which belongs to deployment
	labels := map[string]string{
		"app": sample.Name,
		"controller": sample.Name,
	}
	podList, err := c.kubeClient.CoreV1().Pods(sample.Namespace).List(labels)
	if err != nil {
		return err
	}
	sample.Spec.Pods = podList

	c.recorder.Event(sample, corev1.EventTypeNormal, "Synced", "SampleCRD synced!")
	return nil
}

func (c *SampleCRDController) updateStatus(sample *samplev1.SampleCRD, deployment *appsv1.Deployment) error {
	sampleCopy := sample.DeepCopy()
	sampleCopy.Status.CurrentReplicas = deployment.Status.ReadyReplicas

	_, err := c.sampleClient.SamplecrdV1().SampleCRDs(sample.Namespace).Update(sampleCopy)
	return err
}

func newDeployment(sample *samplev1.SampleCRD) *appsv1.Deployment {
	labels := map[string]string{
		"app": sample.Name,
		"controller": sample.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sample.Spec.DeploymentName,
			Namespace: sample.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(sample, schema.GroupVersionKind{
					Group:   samplev1.SchemeGroupVersion.Group,
					Version: samplev1.SchemeGroupVersion.Version,
					Kind:    "SampleCRD",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: sample.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}