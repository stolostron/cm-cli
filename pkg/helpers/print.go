// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/stolostron/applier/pkg/asset"
	printclusterpoolv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/kubectl/pkg/cmd/get"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
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
	if pf.OutputFormat == nil || len(*pf.OutputFormat) == 0 || *pf.OutputFormat == "wide" {
		reader := printclusterpoolv1alpha1.GetScenarioResourcesReader()
		crd, err := searchCRD(reader, obj.GetObjectKind().GroupVersionKind().Kind)
		if err != nil {
			return err
		}
		printerColumns, err := searchPrinterColumns(crd)
		if err != nil {
			return err
		}
		columns := make([]string, 0)
		for _, pc := range printerColumns {
			if *pf.OutputFormat == "wide" && pc.Priority > 0 || pc.Priority == 0 {
				columns = append(columns, strings.ToUpper(pc.Name)+":"+pc.JSONPath)
			}
		}
		cc := "custom-columns=" + strings.Join(columns, ",")
		pf.OutputFormat = &cc
	}
	printer, err := pf.ToPrinter()
	if err != nil {
		return err
	}
	return printer.PrintObj(obj, os.Stdout)
}

func searchCRD(reader *asset.ScenarioResourcesReader, kind string) (*apiextensionsv1.CustomResourceDefinition, error) {
	crdFileNames, err := reader.AssetNames(nil)
	if err != nil {
		return nil, err
	}
	for _, crdFileName := range crdFileNames {
		crdData, err := reader.Asset(crdFileName)
		if err != nil {
			return nil, err
		}
		crd := &apiextensionsv1.CustomResourceDefinition{}
		err = yaml.Unmarshal(crdData, crd)
		if err != nil {
			continue
		}
		if crd.Spec.Names.Kind != kind {
			continue
		}
		return crd, nil
	}
	return nil, fmt.Errorf("crd %s not found", kind)
}

func searchPrinterColumns(crd *apiextensionsv1.CustomResourceDefinition) ([]apiextensionsv1.CustomResourceColumnDefinition, error) {
	for _, v := range crd.Spec.Versions {
		if v.Name == printclusterpoolv1alpha1.GroupVersion.Version {
			return v.AdditionalPrinterColumns, nil
		}
	}
	return nil, fmt.Errorf("column definition not found for version %s in crd %s", printclusterpoolv1alpha1.GroupVersion.Version, crd.GetName())
}
