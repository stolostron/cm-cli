// Copyright Contributors to the Open Cluster Management project
package policies

import (
	"context"
	"time"

	printpoliciesv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	policyv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policy/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// cmdutil.CheckErr(o.Complete(cmFlags.KubectlFactory, cmd, args))
// cmdutil.CheckErr(o.Validate(cmd))
// cmdutil.CheckErr(o.Run(cmFlags.KubectlFactory, cmd, args))

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
	policyRestConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(policyRestConfig)
	if err != nil {
		return err
	}

	list, err := dynamicClient.Resource(helpers.GvrPol).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
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
