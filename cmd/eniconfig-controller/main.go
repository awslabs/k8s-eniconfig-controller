package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfig, region, masterURL, eniconfigName, eniconfigTagName string
	automaticENIConfig                                             bool

	rootCmd = &cobra.Command{
		Use:   "eniconfig-controller",
		Short: `ENIConfig Controller will listen for new nodes being added and automatically annotate them with the proper ENIConfig CR name`,
		Long:  `When using the new Amazon ENI CNI for Kubernetes you are able to setup a secondary CIDR for the pods to run on using a CRD, this requires you to annotate your nodes with the name of the ENIConfig CR it should use, this controller will automatically annotate them.`,
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-west-2", "AWS Region which the nodes are deployed into.")
	rootCmd.PersistentFlags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().BoolVar(&automaticENIConfig, "automatic-eniconfig", true, "Automatic ENIConfig will configure the controller to load the ENIConfig name from an EC2 tag.")
	rootCmd.PersistentFlags().StringVar(&eniconfigName, "eniconfig-name", "default-eniconfig", "The name of the ENIConfig resource to annotate the nodes with. Ignored if --automatic-eniconfig set.")
	rootCmd.PersistentFlags().StringVar(&eniconfigTagName, "eniconfig-tag-name", "k8s.amazonaws.com/eniConfig", "The name of the EC2 tag name to get the ENIConfig Name from.")
}
