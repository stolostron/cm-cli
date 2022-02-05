// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Attach a cluster
%[1]s attach cluster --values values.yaml

# Attach a cluster with overwritting the cluster name
%[1]s attach cluster --cluster mycluster --values values.yaml

# Attach a cluster with overwritting the cluster name as arg
%[1]s attach cluster mycluster --values values.yaml
`

const (
	scenarioDirectory = "attach"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")
var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "cluster",
		Short:        "Import a cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			isSupported, err := helpers.IsSupported(o.CMFlags)
			if err != nil {
				return err
			}
			if !isSupported {
				return fmt.Errorf("this command '%s %s' is only available on %s or %s",
					helpers.GetExampleHeader(),
					strings.Join(os.Args[1:], " "),
					helpers.RHACM,
					helpers.MCE)
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.clusterName, "cluster", "", "Name of the cluster")
	cmd.Flags().StringVar(&o.clusterServer, "cluster-server", "", "cluster server url of the cluster to import")
	cmd.Flags().StringVar(&o.clusterToken, "cluster-token", "", "token to access the cluster to import")
	cmd.Flags().StringVar(&o.clusterKubeConfig, "cluster-kubeconfig", "", "path to the kubeconfig the cluster to import")
	cmd.Flags().StringVar(&o.clusterKubeConfigContent, "cluster-kubeconfig-content", "", "content of the kubeconfig the cluster to import")
	cmd.Flags().StringVar(&o.importFile, "import-file", "", "the file path and prefix which will contain the import yaml files for manual import")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.waitAgent, "wait", false, "Wait until the klusterlet agent is installed")
	//Not implemented as it requires to import all addon packages
	// cmd.Flags().BoolVar(&o.waitAddOns, "wait-addons", false, "Wait until the klusterlet agent and the addons are is installed")
	cmd.Flags().IntVar(&o.timeout, "timeout", 180, "Timeout to get the klusterlet agent or addons ready in seconds")
	return cmd
}
