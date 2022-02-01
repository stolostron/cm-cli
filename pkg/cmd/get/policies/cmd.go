// Copyright Contributors to the Open Cluster Management project
package policies

import (
	"fmt"
	"os"
	"strings"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get all policies
	%[1]s get policies

	# get a single policy
	%[1]s get policy <policy_name>
	`
)

// NewCmd ...
func NewCmd(f cmdutil.Factory, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	policies := &cobra.Command{
		Use:                   "policies [(-o|--output=)json|yaml|wide|custom-columns=...|custom-columns-file=...|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=...] (TYPE[.VERSION][.GROUP] [NAME | -l label] | TYPE[.VERSION][.GROUP]/NAME ...) [flags]",
		Aliases:               []string{"pol", "policy"},
		DisableFlagsInUseLine: true,
		Short:                 "Display policies",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		PreRunE: func(c *cobra.Command, args []string) error {
			if !helpers.IsRHACM(o.CMFlags) {
				return fmt.Errorf("this command '%s %s' is only available on RHACM", helpers.GetExampleHeader(), strings.Join(os.Args[1:], " "))
			}
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run(f))
		},
		SuggestFor: []string{"list", "ps"},
	}

	o.GetOptions.PrintFlags.AddFlags(policies)

	policies.Flags().BoolVarP(&o.GetOptions.Watch, "watch", "w", o.GetOptions.Watch, "After listing/getting the requested object, watch for changes. Uninitialized objects are excluded if no object name is provided.")
	policies.Flags().BoolVar(&o.GetOptions.WatchOnly, "watch-only", o.GetOptions.WatchOnly, "Watch for changes to the requested object(s), without listing/getting first.")
	policies.Flags().BoolVar(&o.GetOptions.OutputWatchEvents, "output-watch-events", o.GetOptions.OutputWatchEvents, "Output watch event objects when --watch or --watch-only is used. Existing objects are output as initial ADDED events.")
	policies.Flags().Int64Var(&o.GetOptions.ChunkSize, "chunk-size", o.GetOptions.ChunkSize, "Return large lists in chunks rather than all at once. Pass 0 to disable. This flag is beta and may change in the future.")
	policies.Flags().BoolVar(&o.GetOptions.IgnoreNotFound, "ignore-not-found", o.GetOptions.IgnoreNotFound, "If the requested object does not exist the command will return exit code 0.")
	policies.Flags().StringVarP(&o.GetOptions.LabelSelector, "selector", "l", o.GetOptions.LabelSelector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	policies.Flags().StringVar(&o.GetOptions.FieldSelector, "field-selector", o.GetOptions.FieldSelector, "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.")
	policies.Flags().StringVar(&o.GetOptions.Namespace, "cluster", o.GetOptions.Namespace, "List the requested object(s) in the specified cluster namespace.")
	policies.Flags().StringVarP(&o.GetOptions.Namespace, "namespace", "n", o.GetOptions.Namespace, "List the requested object(s) in the specified namespace.")
	addServerPrintColumnFlags(policies, o.GetOptions)
	cmdutil.AddFilenameOptionFlags(policies, &o.GetOptions.FilenameOptions, "identifying the resource to get from a server.")

	return policies
}

const (
	useServerPrintColumns = "server-print"
)

func addServerPrintColumnFlags(cmd *cobra.Command, opt *get.GetOptions) {
	cmd.Flags().BoolVar(&opt.ServerPrint, useServerPrintColumns, opt.ServerPrint, "If true, have the server return the appropriate table output. Supports extension APIs and CRDs.")
}
