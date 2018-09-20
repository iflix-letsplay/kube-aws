package cmd

import (
	"fmt"

	"github.com/kubernetes-incubator/kube-aws/core/root"
	"github.com/spf13/cobra"
)

var (
	cmdUpdate = &cobra.Command{
		Use:          "update",
		Short:        "Update an existing Kubernetes cluster",
		Long:         ``,
		RunE:         runCmdUpdate,
		SilenceUsage: true,
	}

	updateOpts = struct {
		awsDebug, prettyPrint, skipWait bool
	}{}
)

func init() {
	RootCmd.AddCommand(cmdUpdate)
	cmdUpdate.Flags().BoolVar(&updateOpts.awsDebug, "aws-debug", false, "Log debug information from aws-sdk-go library")
	cmdUpdate.Flags().BoolVar(&updateOpts.prettyPrint, "pretty-print", false, "Pretty print the resulting CloudFormation")
	cmdUpdate.Flags().BoolVar(&updateOpts.skipWait, "skip-wait", false, "Don't wait the resources finish")
}

func runCmdUpdate(_ *cobra.Command, _ []string) error {
	opts := root.NewOptions(updateOpts.prettyPrint, updateOpts.skipWait)

	cluster, err := root.ClusterFromFile(configPath, opts, updateOpts.awsDebug)
	if err != nil {
		return fmt.Errorf("Failed to read cluster config: %v", err)
	}

	if _, err := cluster.ValidateStack(); err != nil {
		return err
	}

	report, err := cluster.Update()
	if err != nil {
		return fmt.Errorf("Error updating cluster: %v", err)
	}
	if report != "" {
		fmt.Printf("Update stack: %s\n", report)
	}

	info, err := cluster.Info()
	if err != nil {
		return fmt.Errorf("Failed fetching cluster info: %v", err)
	}

	successMsg :=
		`Success! Your AWS resources are being updated:
%s
`
	fmt.Printf(successMsg, info.String())

	return nil
}
