package main

import (
	"os"
	"time"

	clientset "github.com/aws/amazon-vpc-cni-k8s/pkg/client/clientset/versioned"
	informers "github.com/aws/amazon-vpc-cni-k8s/pkg/client/informers/externalversions"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/christopherhein/eniconfig-controller/pkg/config"
	"github.com/christopherhein/eniconfig-controller/pkg/controller"
	"github.com/christopherhein/eniconfig-controller/pkg/signals"
	"github.com/kris-nova/logger"
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
			logger.Critical("Error building kubeconfig: %s", err.Error())
			os.Exit(1)
		}

		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			logger.Critical("Error building kubernetes clientset: %s", err.Error())
			os.Exit(1)
		}

		eniconfigClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			logger.Critical("Error building eniconfig clientset: %s", err.Error())
			os.Exit(1)
		}

		eniconfigClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building eniconfig clientset: %s", err.Error())
		}

		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
		eniconfigInformerFactory := informers.NewSharedInformerFactory(eniconfigClient, time.Second*30)

		awsSession, err := session.NewSession()
		if err != nil {
			logger.Critical("Error getting creating aws session: %s", err.Error())
			os.Exit(1)
		}

		metadata := ec2metadata.New(awsSession)
		if region == "" {
			region, err = metadata.Region()
			if err != nil {
				logger.Critical("Error getting ec2 region: %s", err.Error())
				os.Exit(1)
			}
		}

		awsSession, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			logger.Critical("Error getting creating aws session: %s", err.Error())
			os.Exit(1)
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
			logger.Critical("Error running controller: %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
