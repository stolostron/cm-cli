// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
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
