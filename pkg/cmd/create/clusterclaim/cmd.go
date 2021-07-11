// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"path/filepath"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster/scenario"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	scenarioDirectory = "create"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")

var example = `
# Create clusterclaims in the current clusterpoolhost
%[1]s create clusterclaim|cc <clusterpool> <clusterclaim_name>[,<clusterclaim_name>...] <options>

# Create clusterclaims on a given clusterpoolhost
%[1]s create clusterclaim|cc <clusterpool> <clusterclaim_name>[,<clusterclaim_name>...] <clusterpoolhost> <options>
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "clusterclaim",
		Aliases:      []string{"cc", "clusterclaims"},
		Short:        "Create clusterclaims",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return clusterpoolhost.BackupCurrentContexts()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return clusterpoolhost.RestoreCurrentContexts()
		},
	}

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")

	return cmd
}
