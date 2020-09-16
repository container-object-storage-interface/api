package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	// objectstorage
	v1alpha1 "github.com/container-object-storage-interface/api/apis/objectstorage.k8s.io/v1alpha1"
	bucketclientset "github.com/container-object-storage-interface/api/clientset"

	// k8s api
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"

	// k8s client
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"

	// logging
	"github.com/golang/glog"

	// config
	"github.com/spf13/viper"
)

type controller struct {
	LeaseDuration time.Duration
	RenewDeadline time.Duration
	RetryPeriod   time.Duration

	// Controller
	ResyncPeriod time.Duration

	// Listeners
	BucketListener              BucketListener
	BucketClassListener         BucketClassListener
	BucketRequestListener       BucketRequestListener
	BucketAccessListener        BucketAccessListener
	BucketAccessClassListener   BucketAccessClassListener
	BucketAccessRequestListener BucketAccessRequestListener

	// leader election
	leaderLock string
	identity   string

	// internal
	initialized  bool
	bucketClient bucketclientset.Interface
	kubeClient   kubeclientset.Interface
}

func New(identity string, leaderLockName string) (*controller, error) {
	cfg, err := func() (*rest.Config, error) {
		kubeConfig := viper.GetString("kube-config")

		if kubeConfig != "" {
			return clientcmd.BuildConfigFromFlags("", kubeConfig)
		}
		return rest.InClusterConfig()
	}()
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubeclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	bucketClient, err := bucketclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	id := identity
	if id == "" {
		id, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	return &controller{
		identity:     id,
		kubeClient:   kubeClient,
		bucketClient: bucketClient,
		initialized:  false,
		leaderLock:   leaderLockName,

		ResyncPeriod:  30 * time.Second,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   5 * time.Second,
	}, nil
}

// Run - runs the controller. Note that ctx must be cancellable i.e. ctx.Done() should not return nil
func (c *controller) Run(ctx context.Context) error {
	if !c.initialized {
		fmt.Errorf("Uninitialized controller. Atleast 1 listener should be added")
	}

	ns := func() string {
		if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
			return ns
		}

		if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
			if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
				return ns
			}
		}
		return "default"
	}()

	leader := func() string {
		re := regexp.MustCompile("[^a-zA-Z0-9-]")
		name := strings.ToLower(re.ReplaceAllString(fmt.Sprintf("%s/%s", c.leaderLock, c.identity), "-"))
		if name[len(name)-1] == '-' {
			// name must not end with '-'
			name = name + "X"
		}
		return name
	}()

	recorder := record.NewBroadcaster()
	recorder.StartRecordingToSink(&corev1.EventSinkImpl{Interface: c.kubeClient.CoreV1().Events(ns)})
	eRecorder := recorder.NewRecorder(scheme.Scheme, v1.EventSource{Component: leader})

	rlConfig := resourcelock.ResourceLockConfig{
		Identity:      c.identity,
		EventRecorder: eRecorder,
	}

	l, err := resourcelock.New(resourcelock.LeasesResourceLock, ns, leader, c.kubeClient.CoreV1(), c.kubeClient.CoordinationV1(), rlConfig)
	if err != nil {
		return err
	}

	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock:          l,
		LeaseDuration: c.LeaseDuration,
		RenewDeadline: c.RenewDeadline,
		RetryPeriod:   c.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				glog.V(2).Info("became leader, starting")
				c.runController(ctx)
			},
			OnStoppedLeading: func() {
				glog.Fatal("stopped leading")
			},
			OnNewLeader: func(identity string) {
				glog.V(3).Infof("new leader detected, current leader: %s", identity)
			},
		},
	}

	leaderelection.RunOrDie(ctx, leaderConfig)
	return nil // should never reach here
}

func (c *controller) runController(ctx context.Context) {
	type addFunc func(ctx context.Context, obj interface{}) error
	type updateFunc func(ctx context.Context, old, new interface{}) error
	type deleteFunc func(ctx context.Context, obj interface{}) error

	controllerFor := func(name string, objType runtime.Object, add addFunc, update updateFunc, delete deleteFunc) {
		indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
		resyncPeriod := c.ResyncPeriod

		lw := cache.NewListWatchFromClient(c.bucketClient.ObjectstorageV1alpha1().RESTClient(), name, "", fields.Everything())
		cfg := &cache.Config{
			Queue: cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
				KnownObjects:          indexer,
				EmitDeltaTypeReplaced: true,
			}),
			ListerWatcher:    lw,
			ObjectType:       objType,
			FullResyncPeriod: resyncPeriod,
			RetryOnError:     true,
			Process: func(obj interface{}) error {
				for _, d := range obj.(cache.Deltas) {
					switch d.Type {
					case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
						if old, exists, err := indexer.Get(d.Object); err == nil && exists {
							if err := indexer.Update(d.Object); err != nil {
								return err
							}
							return update(ctx, old, d.Object)
						} else {
							if err := indexer.Add(d.Object); err != nil {
								return err
							}
							return add(ctx, d.Object)
						}
					case cache.Deleted:
						if err := indexer.Delete(d.Object); err != nil {
							return err
						}
						return delete(ctx, d.Object)
					}
				}
				return nil
			},
		}
		ctrlr := cache.New(cfg)
		ctrlr.Run(ctx.Done())
	}

	if c.BucketListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketListener.Add(ctx, obj.(*v1alpha1.Bucket))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketListener.Update(ctx, old.(*v1alpha1.Bucket), new.(*v1alpha1.Bucket))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketListener.Delete(ctx, obj.(*v1alpha1.Bucket))
		}
		go controllerFor("Buckets", &v1alpha1.Bucket{}, addFunc, updateFunc, deleteFunc)
	}
	if c.BucketRequestListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketRequestListener.Add(ctx, obj.(*v1alpha1.BucketRequest))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketRequestListener.Update(ctx, old.(*v1alpha1.BucketRequest), new.(*v1alpha1.BucketRequest))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketRequestListener.Delete(ctx, obj.(*v1alpha1.BucketRequest))
		}
		go controllerFor("BucketRequests", &v1alpha1.BucketRequest{}, addFunc, updateFunc, deleteFunc)
	}
	if c.BucketClassListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketClassListener.Add(ctx, obj.(*v1alpha1.BucketClass))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketClassListener.Update(ctx, old.(*v1alpha1.BucketClass), new.(*v1alpha1.BucketClass))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketClassListener.Delete(ctx, obj.(*v1alpha1.BucketClass))
		}
		go controllerFor("BucketClasses", &v1alpha1.BucketClass{}, addFunc, updateFunc, deleteFunc)
	}

	if c.BucketAccessListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessListener.Add(ctx, obj.(*v1alpha1.BucketAccess))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketAccessListener.Update(ctx, old.(*v1alpha1.BucketAccess), new.(*v1alpha1.BucketAccess))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessListener.Delete(ctx, obj.(*v1alpha1.BucketAccess))
		}
		go controllerFor("BucketAccesses", &v1alpha1.BucketAccess{}, addFunc, updateFunc, deleteFunc)
	}
	if c.BucketAccessRequestListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessRequestListener.Add(ctx, obj.(*v1alpha1.BucketAccessRequest))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketAccessRequestListener.Update(ctx, old.(*v1alpha1.BucketAccessRequest), new.(*v1alpha1.BucketAccessRequest))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessRequestListener.Delete(ctx, obj.(*v1alpha1.BucketAccessRequest))
		}
		go controllerFor("BucketAccessRequests", &v1alpha1.BucketAccessRequest{}, addFunc, updateFunc, deleteFunc)
	}
	if c.BucketAccessClassListener != nil {
		addFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessClassListener.Add(ctx, obj.(*v1alpha1.BucketAccessClass))
		}
		updateFunc := func(ctx context.Context, old interface{}, new interface{}) error {
			return c.BucketAccessClassListener.Update(ctx, old.(*v1alpha1.BucketAccessClass), new.(*v1alpha1.BucketAccessClass))
		}
		deleteFunc := func(ctx context.Context, obj interface{}) error {
			return c.BucketAccessClassListener.Delete(ctx, obj.(*v1alpha1.BucketAccessClass))
		}
		go controllerFor("BucketAccessClasses", &v1alpha1.BucketAccessClass{}, addFunc, updateFunc, deleteFunc)
	}

	<-ctx.Done()
}
