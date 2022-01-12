// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"

	"github.com/stolostron/cm-cli/pkg/cmd/create/authrealm/scenario"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	scenarioDirectory = "create"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")

var example = `
# Create a authrealm
%[1]s create authrealm myauthrealm --values values.yaml

# Create a authrealm with routeSubDomain overwrite
%[1]s create authrealm myauthrealm --routeSubDomain mysso --values values.yaml
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "authrealm",
		Short:        "Create a authrealm",
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
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
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
	cmd.Flags().StringVar(&o.name, "name", "", "The name of the authrealm")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "", "The name of the authrealm")
	cmd.Flags().StringVar(&o.typeName, "type", "", "type of proxy (dex)")
	cmd.Flags().StringVar(&o.routeSubDomain, "route-sub-domain", "", "the route sub domain")
	cmd.Flags().StringVar(&o.placement, "placement", "", "The name of the placement")
	cmd.Flags().StringVar(&o.managedClusterSet, "cluster-set", "", "The name of the managed cluster set")
	cmd.Flags().StringVar(&o.managedClusterSetBinding, "cluster-set-binding", "", "The of the cluster set binding")
	cmd.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")

	return cmd
}
