package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	cosiapiv1alpha1 "github.com/container-object-storage-interface/api/apis/v1alpha1"
	"github.com/container-object-storage-interface/api/client/cosiv1alpha1"
)

var scheme = runtime.NewScheme()
var kubeconfig *string

func init() {
	cosapiv1alpha1.AddToScheme(sc)

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

type Interface interface {
	kubernetes.Interface

	COSIV1alpha1() cosiv1alpha1.CosiV1alpha1Interface
}

type clientSet struct {
	restClient *rest.Client
}

func (c *clientSet) COSIV1alpha1() BucketV1alpha1Interface {
	return cosiv1alpha1.New(c.restClient)
}

func New() (*clientSet, error) {
	var errInCluster error
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorf("unable to obtain in-cluster-config")
		errInCluster = err
	}

	if errInCluster != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			glog.Errorf("unable to obtain out-of-cluster config", err)
			return nil, fmt.Errorf("unable to obtain-in-cluster: %w and out-of-cluster: %w configs", errInCluster, err)
		}
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("unable to initialize kubernetes client: %v", err)
		return err
	}

	var cl clientSet
	cl.Client = kubeClient
	cl.bucketClient = newBucketClient(config)
	cl.bucketClassClient = newBucketClassClient(config)
	cl.bucketContentClient = newBucketContentClient(config)

	return &cl, nil
}
