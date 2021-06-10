// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/spf13/cobra"
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
				valuesPath: "badpath",
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
			name: "Sucess, with values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-with-data.yaml"),
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
			name: "Success all info in values",
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
				values: map[string]interface{}{},
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
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		valuesPath  string
		values      map[string]interface{}
	}
	type args struct {
		clusterClient clusterclientset.Interface
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
			},
			wantErr: false,
		},
		{
			name: "Failed no managedcluster",
			fields: fields{
				CMFlags:     genericclioptionscm.NewCMFlags(nil),
				clusterName: "test",
			},
			args: args{
				clusterClient: clusterClientNoManagedCluster,
			},
			wantErr: true,
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
			if err := o.runWithClient(tt.args.clusterClient); (err != nil) != tt.wantErr {
				t.Errorf("Options.runWithClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
