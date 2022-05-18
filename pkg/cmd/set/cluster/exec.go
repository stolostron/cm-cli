// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"strings"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/spf13/cobra"
)

var scheduleSkip string

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("cluster names are missing")
	}
	o.Clusters = args[0]
	if cmd.Flags().Lookup("hibernate-schedule-on").Changed {
		scheduleSkip = "true"
	}
	if cmd.Flags().Lookup("hibernate-schedule-off").Changed {
		scheduleSkip = "skip"
	}

	return nil
}

func (o *Options) validate(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("hibernate-schedule-on").Changed &&
		cmd.Flags().Lookup("hibernate-schedule-off").Changed {
		return fmt.Errorf("flags hibernate-schedule-on and hibernate-schedule-off are mutually exclusif")
	}
	if cmd.Flags().Lookup("hibernate-schedule-on").Changed {
		dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
		if err != nil {
			return err
		}
		cjus, err := dynamicClient.Resource(helpers.GvrCD).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		found := false
		for _, cju := range cjus.Items {
			if cju.GetName() == "hibernation-cronjob" {
				found = true
			}
		}
		if !found && !o.HibernateScheduleForce {
			return fmt.Errorf(`
hibernation cronjob not installed,
please visit https://github.com/stolostron/hibernate-cronjob#hibernate-your-hive-provisioned-clusters\n
you can use --hibernate-schedule-force to force the setting`)
		}
	}
	return nil
}

func (o *Options) run() (err error) {

	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	for _, ccn := range strings.Split(o.Clusters, ",") {
		ccn := strings.TrimSpace(ccn)
		cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(ccn).Get(context.TODO(), ccn, metav1.GetOptions{})
		if err != nil {
			return err
		}
		cd := &hivev1.ClusterDeployment{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd)
		if err != nil {
			return err
		}

		cd.Labels["hibernate"] = scheduleSkip

		if !o.CMFlags.DryRun {
			cdu.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cd)
			if err != nil {
				return err
			}
			_, err = dynamicClient.Resource(helpers.GvrCD).Namespace(cd.Namespace).Update(context.TODO(), cdu, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
