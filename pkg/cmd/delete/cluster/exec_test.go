// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/kubectl/pkg/scheme"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterclientsetfake "open-cluster-management.io/api/client/cluster/clientset/versioned/fake"
	cluster "open-cluster-management.io/api/cluster/v1"
)

var testDir = filepath.Join("test", "unit")

func TestOptions_complete(t *testing.T) {
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		valuesPath  string
		values      map[string]interface{}
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
			name: "Failed, bad valuesPath",
			fields: fields{
				valuesPath: "bad-values-path.yaml",
			},
			wantErr: true,
		},
		{
			name: "Failed, empty values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-empty.yaml"),
			},
			wantErr: true,
		},
		{
			name:    "Success, no values.yaml, no name",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Success, no values.yaml",
			fields: fields{
				clusterName: "myCluster",
			},
			wantErr: false,
		},
		{
			name: "Success, with values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-fake.yaml"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				clusterName: tt.fields.clusterName,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
			}
			if err := o.complete(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Options.complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_validate(t *testing.T) {
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		valuesPath  string
		values      map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success AWS all info in values",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name": "test",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed name missing",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed name empty",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name": "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed managedCluster missing",
			fields: fields{
				values: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "Success replace clusterName",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "test",
						"cloud": "aws",
					},
				},
				clusterName: "test2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				clusterName: tt.fields.clusterName,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_runWithClient(t *testing.T) {
	mc := &cluster.ManagedCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	clusterClient := clusterclientsetfake.NewSimpleClientset(mc)
	clusterClientNoManagedCluster := clusterclientsetfake.NewSimpleClientset()
	s := scheme.Scheme
	dynamicClient := fakedynamic.NewSimpleDynamicClient(s)
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		valuesPath  string
		values      map[string]interface{}
	}
	type args struct {
		clusterClient clusterclientset.Interface
		dynamicClient dynamic.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				CMFlags:     genericclioptionscm.NewCMFlags(nil),
				clusterName: "test",
			},
			args: args{
				clusterClient: clusterClient,
				dynamicClient: dynamicClient,
			},
			wantErr: false,
		},
		{
			name: "Success no managedcluster",
			fields: fields{
				CMFlags:     genericclioptionscm.NewCMFlags(nil),
				clusterName: "test",
			},
			args: args{
				clusterClient: clusterClientNoManagedCluster,
				dynamicClient: dynamicClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				clusterName: tt.fields.clusterName,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
			}
			if err := o.runWithClient(tt.args.clusterClient, tt.args.dynamicClient); (err != nil) != tt.wantErr {
				t.Errorf("Options.runWithClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
