// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"os"

	"k8s.io/kubectl/pkg/cmd/get"

	"k8s.io/apimachinery/pkg/runtime"
)

const (
	YamlFormat            string = "yaml"
	JsonFormat            string = "json"
	JsonPathFormat        string = "jsonpath="
	CustomColumnsFormat   string = "columns="
	ColumnsSeparator      string = ","
	SupportedOutputFormat string = YamlFormat + "|" + JsonFormat + "|" + JsonPathFormat + "|" + CustomColumnsFormat + "..."
)

func Print(obj runtime.Object, printFlags *get.PrintFlags) error {
	pf := printFlags.Copy()
	printer, err := pf.ToPrinter()
	if err != nil {
		return err
	}
	return printer.PrintObj(obj, os.Stdout)
}
