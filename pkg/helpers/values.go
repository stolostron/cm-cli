// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"open-cluster-management.io/clusteradm/pkg/helpers/asset"
)

func ConvertValuesFileToValuesMap(path, prefix string) (values map[string]interface{}, err error) {
	var b []byte
	if path != "" {
		b, err = ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		b = append(b, '\n')
		pdata, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		b = append(b, pdata...)
	}

	valuesc := make(map[string]interface{})
	err = yaml.Unmarshal(b, &valuesc)
	if err != nil {
		if path != "" {
			fmt.Printf("Error while unmarshaling stdin or values file %s\n", path)
		} else {
			fmt.Printf("Error while unmarshaling stdin:\n%s\n", string(b))
		}
		return nil, err
	}

	values = make(map[string]interface{})
	if prefix != "" {
		values[prefix] = valuesc
	} else {
		values = valuesc
	}

	return values, nil
}

func ConvertReaderFileToValuesMap(path string,
	reader *asset.ScenarioResourcesReader) (values map[string]interface{}, err error) {
	values = make(map[string]interface{})
	b, err := reader.Asset(path)
	if err != nil {
		return values, err
	}
	err = yaml.Unmarshal(b, &values)
	if err != nil {
		return values, err
	}
	return values, nil
}
