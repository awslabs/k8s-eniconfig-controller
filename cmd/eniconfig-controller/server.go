package main

import (
	"time"

	clientset "github.com/aws/amazon-vpc-cni-k8s/pkg/client/clientset/versioned"
	informers "github.com/aws/amazon-vpc-cni-k8s/pkg/client/informers/externalversions"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/christopherhein/eniconfig-controller/pkg/config"
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

		eniconfigClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building eniconfig clientset: %s", err.Error())
		}

		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
		eniconfigInformerFactory := informers.NewSharedInformerFactory(eniconfigClient, time.Second*30)

		awsSession, err := session.NewSession()
		if err != nil {
			glog.Fatalf("Error getting creating aws session: %s", err.Error())
		}

		metadata := ec2metadata.New(awsSession)
		if region == "" {
			region, err = metadata.Region()
			if err != nil {
				glog.Fatalf("Error getting ec2 region: %s", err.Error())
			}
		}

		awsSession, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			glog.Fatalf("Error getting creating aws session: %s", err.Error())
		}

		conf := config.New(automaticENIConfig, eniconfigName, eniconfigTagName, awsSession)

		ctrl := controller.New(kubeClient,
			kubeInformerFactory.Core().V1().Nodes(),
			eniconfigInformerFactory.Crd().V1alpha1().ENIConfigs(),
			conf)

		// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
		// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
		kubeInformerFactory.Start(stopCh)
		eniconfigInformerFactory.Start(stopCh)

		if err = ctrl.Run(2, stopCh); err != nil {
			glog.Fatalf("Error running controller: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
