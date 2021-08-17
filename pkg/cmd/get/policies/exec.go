// Copyright Contributors to the Open Cluster Management project
package policies

import (
	"context"
	"fmt"
	"time"

	printpoliciesv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	policyv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policy/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.PolicyName = args[0]
	}
	return nil
}

func (o *Options) validate() (err error) {
	return nil
}

func (o *Options) run(f cmdutil.Factory) (err error) {
	// Retrieve default namespace from config
	defaultNamespace, _, err := o.CMFlags.KubectlFactory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		fmt.Println(err.Error())
	}
	// Create dynamic client to retrieve resources
	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return err
	}
	// Retrieve policies
	list := &unstructured.UnstructuredList{}
	// Retrieve individual policy (use namespace from config if not provided)
	if o.PolicyName != "" {
		var policyu *unstructured.Unstructured
		if o.GetOptions.Namespace != "" {
			policyu, err = dynamicClient.Resource(helpers.GvrPol).Namespace(o.GetOptions.Namespace).Get(context.TODO(), o.PolicyName, metav1.GetOptions{})
		} else {
			policyu, err = dynamicClient.Resource(helpers.GvrPol).Namespace(defaultNamespace).Get(context.TODO(), o.PolicyName, metav1.GetOptions{})
		}
		if err != nil {
			return err
		}
		list.Items = append(list.Items, *policyu)
	} else {
		// Retrieve list of policies (first, try all namespaces (or given namespace) and then try namespace from config)
		list, err = dynamicClient.Resource(helpers.GvrPol).Namespace(o.GetOptions.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil && o.GetOptions.Namespace == "" {
			fmt.Println(err.Error())
			fmt.Printf("Failed to retrieve policies clusterwide. Attempting to retrieve policies from namespace '%s':\n", defaultNamespace)
			list, err = dynamicClient.Resource(helpers.GvrPol).Namespace(defaultNamespace).List(context.TODO(), metav1.ListOptions{})
		}
		if err != nil {
			return err
		}
	}
	// Construct and print policies using PrintPolicies CRD
	policy := &policyv1.Policy{}
	printPoliciesList := &printpoliciesv1alpha1.PrintPoliciesList{}
	printPoliciesList.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printpoliciesv1alpha1.GroupName,
				Kind:    "PrintPolicies",
				Version: printpoliciesv1alpha1.GroupVersion.Version})
	for _, policyu := range list.Items {
		if runtime.DefaultUnstructuredConverter.FromUnstructured(policyu.UnstructuredContent(), policy); err != nil {
			return err
		}
		printpol := printpoliciesv1alpha1.PrintPolicies{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policy.Name,
				Namespace: policy.Namespace,
			},
			Spec: printpoliciesv1alpha1.PrintPoliciesSpec{
				Policy: *policy,
				Age:    helpers.TimeDiff(policy.CreationTimestamp.Time, time.Second),
			},
		}
		printPoliciesList.Items = append(printPoliciesList.Items, printpol)
	}
	helpers.Print(printPoliciesList, o.GetOptions.PrintFlags)
	return err
}
