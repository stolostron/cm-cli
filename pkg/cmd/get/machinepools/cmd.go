// Copyright Contributors to the Open Cluster Management project
package machinepools

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get machine pools
	%[1]s machinepools <cluster-name> -oyaml`
)

var clusterName string

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := get.NewGetOptions("cm", streams)
	machinepools := &cobra.Command{
		Use:                   "machinepools [(-o|--output=)json|yaml|wide|custom-columns=...|custom-columns-file=...|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=...] (TYPE[.VERSION][.GROUP] [NAME | -l label] | TYPE[.VERSION][.GROUP]/NAME ...) [flags]",
		Aliases:               []string{"machinepool", "mp"},
		DisableFlagsInUseLine: true,
		Short:                 "list the machinepools for a give cluster",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"machinepools"}, args...)
			cmd.Flags().Set("namespace", clusterName)
			cmdutil.CheckErr(o.Complete(cmFlags.KubectlFactory, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run(cmFlags.KubectlFactory, cmd, args))
		},
		SuggestFor: []string{"list", "ps"},
	}

	o.PrintFlags.AddFlags(machinepools)

	machinepools.Flags().StringVar(&clusterName, "cluster", "", "Name of the cluster")
	machinepools.Flags().BoolVarP(&o.Watch, "watch", "w", o.Watch, "After listing/getting the requested object, watch for changes. Uninitialized objects are excluded if no object name is provided.")
	machinepools.Flags().BoolVar(&o.WatchOnly, "watch-only", o.WatchOnly, "Watch for changes to the requested object(s), without listing/getting first.")
	machinepools.Flags().BoolVar(&o.OutputWatchEvents, "output-watch-events", o.OutputWatchEvents, "Output watch event objects when --watch or --watch-only is used. Existing objects are output as initial ADDED events.")
	machinepools.Flags().Int64Var(&o.ChunkSize, "chunk-size", o.ChunkSize, "Return large lists in chunks rather than all at once. Pass 0 to disable. This flag is beta and may change in the future.")
	machinepools.Flags().BoolVar(&o.IgnoreNotFound, "ignore-not-found", o.IgnoreNotFound, "If the requested object does not exist the command will return exit code 0.")
	machinepools.Flags().StringVarP(&o.LabelSelector, "selector", "l", o.LabelSelector, "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	machinepools.Flags().StringVar(&o.FieldSelector, "field-selector", o.FieldSelector, "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The server only supports a limited number of field queries per type.")
	machinepools.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", o.AllNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	addOpenAPIPrintColumnFlags(machinepools, o)
	addServerPrintColumnFlags(machinepools, o)
	cmdutil.AddFilenameOptionFlags(machinepools, &o.FilenameOptions, "identifying the resource to get from a server.")

	return machinepools
}

const (
	useOpenAPIPrintColumnFlagLabel = "use-openapi-print-columns"
	useServerPrintColumns          = "server-print"
)

func addOpenAPIPrintColumnFlags(cmd *cobra.Command, opt *get.GetOptions) {
	cmd.Flags().BoolVar(&opt.PrintWithOpenAPICols, useOpenAPIPrintColumnFlagLabel, opt.PrintWithOpenAPICols, "If true, use x-kubernetes-print-column metadata (if present) from the OpenAPI schema for displaying a resource.")
	cmd.Flags().MarkDeprecated(useOpenAPIPrintColumnFlagLabel, "deprecated in favor of server-side printing")
}

func addServerPrintColumnFlags(cmd *cobra.Command, opt *get.GetOptions) {
	cmd.Flags().BoolVar(&opt.ServerPrint, useServerPrintColumns, opt.ServerPrint, "If true, have the server return the appropriate table output. Supports extension APIs and CRDs.")
}
