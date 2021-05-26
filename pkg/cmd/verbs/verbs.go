// Copyright Contributors to the Open Cluster Management project
package verbs

import (
	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	attachcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster"
	createcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster"
	deletecluster "github.com/open-cluster-management/cm-cli/pkg/cmd/delete/cluster"
	detachcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/detach/cluster"
	scalecluster "github.com/open-cluster-management/cm-cli/pkg/cmd/scale/cluster"
	getclusters "github.com/open-cluster-management/cm-cli/pkg/cmd/get/clusters"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/version"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewVerbCreate(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use: parent,
	}
	cmd.AddCommand(
		createcluster.NewCmd(streams),
	)

	return cmd
}

func NewVerbGet(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := get.NewCmdGet("cm", f, streams)

	cmd.AddCommand(
		getclusters.NewCmd(f, streams),
	)
	return cmd
}

func NewVerbUpdate(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   parent,
		Short: "Not yet implemented",
	}

	return cmd
}

func NewVerbDelete(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use: parent,
	}
	cmd.AddCommand(
		deletecluster.NewCmd(streams),
	)

	return cmd
}

func NewVerbList(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   parent,
		Short: "Not yet implemented",
	}

	return cmd
}

func NewVerbApplier(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := appliercmd.NewCmd(streams)

	return cmd
}

func NewVerbAttach(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   parent,
		Short: "Attach cluster to hub",
	}

	cmd.AddCommand(attachcluster.NewCmd(streams))

	return cmd
}

func NewVerbDetach(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   parent,
		Short: "Detatch a cluster from the hub",
	}

	cmd.AddCommand(detachcluster.NewCmd(streams))

	return cmd
}

func newVerbScale(verb string, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use: verb,
		Short: "Scale a worker pool on a managed hub",
	}
	cmd.AddCommand(
		scalecluster.NewCmd(streams),
	)

func NewVerbVersion(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := version.NewCmd(streams)

	return cmd
}
