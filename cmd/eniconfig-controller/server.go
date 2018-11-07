package main

import (
	"time"

	"github.com/christopherhein/eniconfig-controller/pkg/controller"
	"github.com/christopherhein/eniconfig-controller/pkg/signals"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the controller",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// set up signals so we handle the first shutdown signal gracefully
		stopCh := signals.SetupSignalHandler()

		cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			glog.Fatalf("Error building kubeconfig: %s", err.Error())
		}

		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		}

		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

		ctrl := controller.New(kubeClient,
			kubeInformerFactory.Core().V1().Nodes(),
			eniconfigName)

		// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
		// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
		kubeInformerFactory.Start(stopCh)

		if err = ctrl.Run(2, stopCh); err != nil {
			glog.Fatalf("Error running controller: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
