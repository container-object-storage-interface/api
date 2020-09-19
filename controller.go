package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/container-object-storage-interface/api/apis/objectstorage.k8s.io/v1alpha1"
	bucketclientset "github.com/container-object-storage-interface/api/clientset"
	bucketcontroller "github.com/container-object-storage-interface/api/controller"
	kubeclientset "k8s.io/client-go/kubernetes"

	"github.com/golang/glog"
)

var ctx context.Context
var cmd = &cobra.Command{
	Use:           os.Args[0],
	Short:         "Sample controller for listening and responding to bucket* and bucketAccess* API objects",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(c *cobra.Command, args []string) error {
		return run(args)
	},
	DisableFlagsInUseLine: true,
}

var kubeConfig string

func init() {
	viper.AutomaticEnv()

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.Set("logtostderr", "true")

	strFlag := func(c *cobra.Command, ptr *string, name string, short string, dfault string, desc string) {
		c.PersistentFlags().
			StringVarP(ptr, name, short, dfault, desc)
	}
	strFlag(cmd, &kubeConfig, "kube-config", "", kubeConfig, "path to kubeconfig file")

	hideFlag := func(name string) {
		cmd.PersistentFlags().MarkHidden(name)
	}
	hideFlag("alsologtostderr")
	hideFlag("log_backtrace_at")
	hideFlag("log_dir")
	hideFlag("logtostderr")
	hideFlag("master")
	hideFlag("stderrthreshold")
	hideFlag("vmodule")

	// suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})
	viper.BindPFlags(cmd.PersistentFlags())

	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	go func() {
		s := <-sigs
		cancel()
		panic(fmt.Sprintf("%s %s", s.String(), "Signal received. Exiting"))
	}()
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

type bucketListener struct {
	kubeClient   kubeclientset.Interface
	bucketClient bucketclientset.Interface
}

func (b *bucketListener) InitializeKubeClient(k kubeclientset.Interface) {
	b.kubeClient = k
}

func (b *bucketListener) InitializeBucketClient(bc bucketclientset.Interface) {
	b.bucketClient = bc
}

func (b *bucketListener) Add(ctx context.Context, obj *v1alpha1.Bucket) error {
	glog.V(1).Infof("add called for bucket %s", obj.Name)
	return nil
}

func (b *bucketListener) Update(ctx context.Context, old, new *v1alpha1.Bucket) error {
	return nil
}

func (b *bucketListener) Delete(ctx context.Context, obj *v1alpha1.Bucket) error {
	return nil
}

type bucketAccessListener struct {
	kubeClient   kubeclientset.Interface
	bucketClient bucketclientset.Interface

	count int
}

func (b *bucketAccessListener) InitializeKubeClient(k kubeclientset.Interface) {
	b.kubeClient = k
}

func (b *bucketAccessListener) InitializeBucketClient(bc bucketclientset.Interface) {
	b.bucketClient = bc
}

func (b *bucketAccessListener) Add(ctx context.Context, obj *v1alpha1.BucketAccess) error {
	b.count += 1
	if b.count < 5 {
		return errors.New("failing")
	}
	glog.V(1).Infof("add called for bucketAccess %s", obj.Name)
	return nil
}

func (b *bucketAccessListener) Update(ctx context.Context, old, new *v1alpha1.BucketAccess) error {
	return nil
}

func (b *bucketAccessListener) Delete(ctx context.Context, obj *v1alpha1.BucketAccess) error {
	return nil
}

func run(args []string) error {
	ctrl, err := bucketcontroller.NewDefaultObjectStorageController("sample-controller", "leader-lock", 40)
	if err != nil {
		glog.Error(err)
		return err
	}
	ctrl.AddBucketListener(&bucketListener{})
	ctrl.AddBucketAccessListener(&bucketAccessListener{})
	return ctrl.Run(ctx)
}
