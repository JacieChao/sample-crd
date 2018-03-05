package main

import (
	"flag"
	"time"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	kubeinformers "k8s.io/client-go/informers"

	"github.com/golang/glog"
	sampleClient "github.com/sample-crd/pkg/client/clientset/versioned"
	sampleinformers "github.com/sample-crd/pkg/client/informers/externalversions"
	"github.com/sample-crd/pkg/controller"
)

var (
	masterURL  string
	kubeconfig string
	stopCh chan struct{}

)

func init() {
	stopCh = make(chan struct{})
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server.")
}

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	sampleClient, err := sampleClient.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	sampleInformerFactory := sampleinformers.NewSharedInformerFactory(sampleClient, time.Second*30)

	controller := controller.NewSampleController(kubeClient, sampleClient, kubeInformerFactory, sampleInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go sampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running sample controller: %s", err.Error())
	}
}