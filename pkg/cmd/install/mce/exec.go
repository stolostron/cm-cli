// Copyright Contributors to the Open Cluster Management project
package mce

import (
	"context"
	"fmt"
	"time"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/cmd/install/mce/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	kubeClient, apiextensionsClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, apiextensionsClient, dynamicClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	apiextensionsClient apiextensionsclient.Interface,
	dynamicClient dynamic.Interface) (err error) {
	_, err = dynamicClient.Resource(helpers.GvrMCEV1alpha1).Namespace(o.namespace).Get(context.TODO(), "multiclusterhub", metav1.GetOptions{})
	if err == nil {
		return errors.NewUnauthorized("mce already installed")
	}
	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	files := []string{
		"install/namespace.yaml",
	}

	approval := "Automatic"
	if o.manualApproval {
		approval = "Manual"
	}
	values := struct {
		Channel       string
		Namespace     string
		OperatorGroup string
		Approval      string
	}{
		Channel:       o.channel,
		Namespace:     o.namespace,
		OperatorGroup: o.operatorGroup,
		Approval:      approval,
	}

	applierBuilder := clusteradmapply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()
	out, err := applier.ApplyDirectly(reader, values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	files = []string{
		"install/operator-group.yaml",
		"install/subscription.yaml",
	}
	out, err = applier.ApplyCustomResources(reader, values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	//Wait MCE CRD to be created
	if !o.CMFlags.DryRun {
		err = wait.PollImmediate(10*time.Second, time.Duration(3)*time.Minute, func() (bool, error) {
			_, err := apiextensionsClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), "multiclusterengines.multicluster.openshift.io", metav1.GetOptions{})
			if err != nil {
				fmt.Printf("%s, waiting...\n", err)
				return false, nil
			}
			return true, nil
		})

		if err != nil {
			return err
		}
	}

	files = []string{
		"install/multicluster-engine_v1alpha1.yaml",
	}
	out, err = applier.ApplyCustomResources(reader, values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		files = []string{
			"install/multicluster-engine_v1.yaml",
		}
		out, err = applier.ApplyCustomResources(reader, values, o.CMFlags.DryRun, "", files...)
		if err != nil {
			return err
		}
	}
	output = append(output, out...)

	if o.wait {
		i := 0
		wait.PollImmediate(1*time.Minute, time.Duration(o.timeout)*time.Minute, func() (bool, error) {
			mchu, err := dynamicClient.Resource(helpers.GvrMCEV1alpha1).Namespace(o.namespace).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
			if errors.IsNotFound(err) {
				mchu, err = dynamicClient.Resource(helpers.GvrMCEV1).Namespace(o.namespace).Get(context.TODO(), "multiclusterengine", metav1.GetOptions{})
			}
			if err != nil {
				fmt.Printf("%s, waiting...\n", err)
				return false, nil
			}
			i += 1
			if statusu, ok := mchu.Object["status"]; ok {
				status := statusu.(map[string]interface{})
				if phaseu, ok := status["phase"]; ok {
					phase := phaseu.(string)
					if phase != "Running" {
						fmt.Printf("(%d/%d), phase is %s, waiting for Running", i, o.timeout, phase)
						return false, nil
					}
				}
			}
			return true, nil
		})
	}
	return clusteradmapply.WriteOutput(o.outputFile, output)
}
