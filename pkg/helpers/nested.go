// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func NestedString(obj map[string]interface{}, dotedPath string) (string, error) {
	fields := strings.Split(dotedPath, ".")
	s, ok, err := unstructured.NestedString(obj, fields...)
	if err != nil {
		return s, err
	}
	if !ok {
		return s, fmt.Errorf("%s is missing", dotedPath)
	}
	if len(s) == 0 {
		return s, fmt.Errorf("%s not specified", dotedPath)
	}
	return s, nil
}

func SetNestedField(obj map[string]interface{}, value interface{}, dotedPath string) error {
	fields := strings.Split(dotedPath, ".")
	return unstructured.SetNestedField(obj, value, fields...)
}

//NestedExists returns true if the nested field exists and an error if unable to traverse obj.
func NestedExists(obj map[string]interface{}, dotedPath string) (bool, error) {
	fields := strings.Split(dotedPath, ".")
	_, ok, err := unstructured.NestedFieldNoCopy(obj, fields...)
	return ok, err
}
