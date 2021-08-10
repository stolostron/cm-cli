// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
)

func ExecuteWithContext(context string, args []string, dryRun bool, streams genericclioptions.IOStreams, outputFile string) error {
	var newArgs []string
	index := -1
	for i, a := range args {
		if a == "--" {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("sub-command separator '--' missing in commannd %s", strings.Join(args, ""))
	}
	if index+1 == len(args) {
		return fmt.Errorf("no command provided after '--'")
	}
	newArgs = args[index+1:]
	for _, a := range newArgs {
		if strings.HasPrefix(a, "--context") {
			fmt.Printf("--context is overwritten with --context=%s", context)
			break
		}
	}
	newArgs = append(newArgs, "--context", context)
	klog.V(5).Infof("Args: %s, newCMD: %s\n", strings.Join(args, " "), strings.Join(newArgs, " "))
	cmd := exec.Command(newArgs[0], newArgs[1:]...)
	cmd.Stdout = streams.Out
	cmd.Stderr = streams.ErrOut
	return cmd.Run()
}
