package controller

import (
	"fmt"
	"time"

	eniv1alpha1 "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/crd.k8s.amazonaws.com/v1alpha1"
	clientset "github.com/aws/amazon-vpc-cni-k8s/pkg/client/clientset/versioned"
	informers "github.com/aws/amazon-vpc-cni-k8s/pkg/client/informers/externalversions/crd.k8s.amazonaws.com/v1alpha1"
	listers "github.com/aws/amazon-vpc-cni-k8s/pkg/client/listers/crd.k8s.amazonaws.com/v1alpha1"
	"github.com/christopherhein/eniconfig-controller/pkg/config"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const (
	controllerAgentName   = "eniconfig-controller"
	annotationName        = "k8s.amazonaws.com/eniConfig"
	SuccessSynced         = "Synced"
	MessageResourceSynced = "Node and ENIConfig synced successfully"
)

// controller is the controller implementation
type controller struct {
	kubeclientset      kubernetes.Interface
	eniconfigclientset clientset.Interface

	nodesLister      corelisters.NodeLister
	nodesSynced      cache.InformerSynced
	eniconfigsLister listers.ENIConfigLister
	eniconfigsSynced cache.InformerSynced

	conf config.Config

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

// New returns a new eniconfig controller
func New(
	kubeclientset kubernetes.Interface,
	nodeInformer coreinformers.NodeInformer,
	eniconfigInformer informers.ENIConfigInformer,
	conf config.Config) *controller {

	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	ctrl := &controller{
		kubeclientset:    kubeclientset,
		nodesLister:      nodeInformer.Lister(),
		nodesSynced:      nodeInformer.Informer().HasSynced,
		eniconfigsLister: eniconfigInformer.Lister(),
		eniconfigsSynced: eniconfigInformer.Informer().HasSynced,
		conf:             conf,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ENIConfigs"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	eniconfigInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.handleENIConfig,
		UpdateFunc: func(old, updated interface{}) {
			newENI := updated.(*eniv1alpha1.ENIConfig)
			oldENI := old.(*eniv1alpha1.ENIConfig)
			if newENI.ResourceVersion == oldENI.ResourceVersion {
				return
			}
			ctrl.handleENIConfig(updated)
		},
		DeleteFunc: ctrl.handleENIConfig,
	})

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.enqueueNode,
		UpdateFunc: func(old, updated interface{}) {
			updatedNode := updated.(*corev1.Node)
			oldNode := old.(*corev1.Node)
			if updatedNode.ResourceVersion == oldNode.ResourceVersion {
				return
			}
			ctrl.enqueueNode(updated)
		},
	})

	return ctrl
}

// Run will configure the event handlers, sync the caches and will block until
// stopCh is closed, then it will shutdown and wait for work items to process
func (c *controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.nodesSynced, c.eniconfigsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("started workers")
	<-stopCh
	glog.Info("shutting down workers")

	return nil
}

// runWorker runs forever and processes work items using processNextWorkItem
func (c *controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will pop off the queue and process using the syncHandler
func (c *controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Infof("successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHander will reconcile the desired/actual state
func (c *controller) syncHandler(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	node, err := c.nodesLister.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("node '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	eniconfigName, err := c.conf.GetName(node.Spec.ProviderID)
	if err != nil {
		return err
	}

	if node.Annotations[annotationName] == eniconfigName {
		return nil
	}

	eniconfig, err := c.eniconfigsLister.Get(eniconfigName)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("eniconfig '%s' not found", eniconfigName))
			return nil
		}

		return err
	}

	instanceAZ, err := c.conf.GetInstanceAZ(node.Spec.ProviderID)
	if err != nil {
		return err
	}

	subnetAZ, err := c.conf.GetSubnetAZ(eniconfig.Spec.Subnet)
	if err != nil {
		return err
	}

	if instanceAZ != subnetAZ {
		runtime.HandleError(fmt.Errorf("instance AZ doesn't match ENIConfig subnet AZ '%s' != '%s", instanceAZ, subnetAZ))
		return nil
	}

	err = c.updateNodeAnnotations(node, eniconfigName)
	if err != nil {
		return err
	}

	c.recorder.Event(node, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// updateNodeAnnotations will update the annotation to equal the proper ENIConfig
func (c *controller) updateNodeAnnotations(node *corev1.Node, eniconfigName string) error {
	nodeCopy := node.DeepCopy()
	nodeCopy.Annotations[annotationName] = eniconfigName
	_, err := c.kubeclientset.Core().Nodes().Update(nodeCopy)
	return err
}

// enqueueNode will translate the node into a workable item
func (c *controller) enqueueNode(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *controller) handleENIConfig(obj interface{}) {
	nodes, err := c.nodesLister.List(labels.Everything())
	if err != nil {
		runtime.HandleError(fmt.Errorf("Error listing nodes."))
		return
	}
	for _, node := range nodes {
		c.enqueueNode(node)
	}

	return
}
