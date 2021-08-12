// Copyright Contributors to the Open Cluster Management project

package main

import (
	"os"

	"k8s.io/klog/v2"

	"github.com/open-cluster-management/cm-cli/pkg/cmd"
)

func main() {
	root := cmd.NewCMCommand()
	err := root.Execute()
	if err != nil {
		klog.V(1).ErrorS(err, "Error:")
	}
	klog.Flush()
	if err != nil {
		os.Exit(1)
	}
}
