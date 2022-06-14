// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"

	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func SetComponentEnable(cmFlags *genericclioptions.CMFlags, componentName string, enable bool) error {
	dynamicClient, err := cmFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	//Update MCE
	foundMCE := false
	version := GvrMCEV1alpha1
	mceu, err := dynamicClient.Resource(version).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		version = GvrMCEV1
		mceu, err = dynamicClient.Resource(version).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
	}
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err == nil {
		components, _, err := unstructured.NestedSlice(mceu.Object, "spec", "overrides", "components")
		if err != nil {
			return err
		}
		for i := range components {
			component := components[i].(map[string]interface{})
			if component["name"].(string) == componentName {
				component["enabled"] = enable
				foundMCE = true
				break
			}
		}
		if foundMCE {
			err = unstructured.SetNestedSlice(mceu.Object, components, "spec", "overrides", "components")
			if err != nil {
				return err
			}
			_, err = dynamicClient.Resource(version).Update(context.TODO(), mceu, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	//Update MCH

	foundMCH := false
	mchs, err := dynamicClient.Resource(GvrMCH).List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		if len(mchs.Items) != 0 {
			mch := mchs.Items[0].DeepCopy()
			components, _, err := unstructured.NestedSlice(mch.Object, "spec", "overrides", "components")
			if err != nil {
				return err
			}
			for i := range components {
				component := components[i].(map[string]interface{})
				if component["name"].(string) == componentName {
					component["enabled"] = enable
					foundMCH = true
					break
				}
			}
			if foundMCH {
				err = unstructured.SetNestedSlice(mch.Object, components, "spec", "overrides", "components")
				if err != nil {
					return err
				}
				_, err = dynamicClient.Resource(GvrMCH).Namespace(mch.GetNamespace()).Update(context.TODO(), mch, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
			}
		}
	}
	if !foundMCE && !foundMCH {
		return fmt.Errorf("component %s not found", componentName)
	}
	return nil
}
