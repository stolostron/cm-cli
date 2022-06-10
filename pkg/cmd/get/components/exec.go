// Copyright Contributors to the Open Cluster Management project
package components

import (
	"context"

	"github.com/spf13/cobra"
	printv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run(streams genericclioptions.IOStreams) (err error) {

	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	mceu, err := dynamicClient.Resource(helpers.GvrMCEV1alpha1).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		mceu, err = dynamicClient.Resource(helpers.GvrMCEV1).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
	}
	if err != nil {
		return err
	}
	componentsMap := make(map[string]bool, 0)
	components, _, err := unstructured.NestedSlice(mceu.Object, "spec", "overrides", "components")
	if err != nil {
		return err
	}
	mchs, err := dynamicClient.Resource(helpers.GvrMCH).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(mchs.Items) != 0 {
		mchComponents, _, err := unstructured.NestedSlice(mchs.Items[0].Object, "spec", "overrides", "components")
		if err != nil {
			return err
		}
		components = append(components, mchComponents...)
	}
	for _, imceComponents := range components {
		mceComponents := imceComponents.(map[string]interface{})
		componentsMap[mceComponents["name"].(string)] = mceComponents["enabled"].(bool)
	}
	printComponentList := &printv1alpha1.PrintComponentList{}
	printComponentList.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printv1alpha1.GroupName,
				Kind:    "PrintComponent",
				Version: printv1alpha1.GroupVersion.Version})
	for k, v := range componentsMap {
		printComponent := printv1alpha1.PrintComponent{
			ObjectMeta: metav1.ObjectMeta{
				Name: k,
			},
			Spec: printv1alpha1.PrintComponentSpec{
				Enabled: v,
			},
		}
		printComponentList.Items = append(printComponentList.Items, printComponent)
	}
	return helpers.Print(printComponentList, o.GetOptions.PrintFlags)
}
