package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// masterURL, kubeconfig, eniconfigName
	masterURL, kubeconfig, eniconfigName string

	// rootCmd initializes the base cobra command
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
	rootCmd.PersistentFlags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().StringVar(&eniconfigName, "eniconfig-name", "default-eniconfig", "The name of the ENIConfig resource to annotate the nodes with.")

	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("master", rootCmd.PersistentFlags().Lookup("master"))
	viper.BindPFlag("eniconfig-name", rootCmd.PersistentFlags().Lookup("eniconfig-name"))
}
