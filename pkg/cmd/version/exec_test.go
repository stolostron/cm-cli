// Copyright Contributors to the Open Cluster Management project
package version

import (
	"testing"

	"github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekubernetes "k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubectl/pkg/scheme"
)

func TestOptions_complete(t *testing.T) {
	type fields struct {
		CMFlags *genericclioptions.CMFlags
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags: tt.fields.CMFlags,
			}
			if err := o.complete(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Options.complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_validate(t *testing.T) {
	type fields struct {
		CMFlags *genericclioptions.CMFlags
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags: tt.fields.CMFlags,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_runWithClient(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myName",
			Namespace: "myNamespace",
			Labels: map[string]string{
				"ocm-configmap-type":  "image-manifest",
				"ocm-release-version": "2.3.0",
			},
		},
		Data: map[string]string{},
	}
	kubeClient := fakekubernetes.NewSimpleClientset(cm)
	s := scheme.Scheme
	dynamicClient := fakedynamic.NewSimpleDynamicClient(s)
	type fields struct {
		CMFlags *genericclioptions.CMFlags
	}
	type args struct {
		kubeClient    kubernetes.Interface
		dynamicClient dynamic.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				kubeClient:    kubeClient,
				dynamicClient: dynamicClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags: tt.fields.CMFlags,
			}
			if err := o.runWithClient(tt.args.kubeClient, tt.args.dynamicClient); (err != nil) != tt.wantErr {
				t.Errorf("Options.runWithClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
