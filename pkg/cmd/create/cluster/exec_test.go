// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	fakeapiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekubernetes "k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubectl/pkg/scheme"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned/fake"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned/fake"
)

var testDir = filepath.Join("test", "unit")

func TestOptions_complete(t *testing.T) {
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		cloud       string
		valuesPath  string
		values      map[string]interface{}
		outputFile  string
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
			name: "Failed, empty values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-empty.yaml"),
			},
			wantErr: true,
		},
		{
			name: "Sucess, with values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-fake-aws.yaml"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				clusterName: tt.fields.clusterName,
				cloud:       tt.fields.cloud,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
				outputFile:  tt.fields.outputFile,
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
		cloud       string
		valuesPath  string
		values      map[string]interface{}
		outputFile  string
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
						"name":  "test",
						"cloud": "aws",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success Azure all info in values",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "test",
						"cloud": "azure",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success GCP all info in values",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "test",
						"cloud": "gcp",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success VSphere all info in values",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "test",
						"cloud": "vsphere",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed, bad valuesPath",
			fields: fields{
				valuesPath: "bad-values-path.yaml",
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
			name: "Failed name missing",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"cloud": "vsphere",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed name empty",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "",
						"cloud": "vsphere",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed cloud missing",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name": "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Failed cloud enpty",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name":  "test",
						"cloud": "",
					},
				},
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
				cloud:       tt.fields.cloud,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
				outputFile:  tt.fields.outputFile,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "Success replace clusterName" {
				if imc, ok := o.values["managedCluster"]; ok {
					mc := imc.(map[string]interface{})
					if icn, ok := mc["name"]; ok {
						cm := icn.(string)
						if cm != "test2" {
							t.Errorf("got %s and expected %s", tt.fields.clusterName, cm)
						}
					} else {
						t.Error("name not found")
					}
				} else {
					t.Error("managedCluster not found")
				}
			}
		})
	}
}

func TestOptions_runWithClient(t *testing.T) {
	pullSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pull-secret",
			Namespace: "openshift-config",
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte("crds: mycrds"),
		},
	}
	values, err := helpers.ConvertValuesFileToValuesMap(filepath.Join(testDir, "values-fake-aws.yaml"), "")
	if err != nil {
		t.Fatal(err)
	}
	apiextensionsClient := fakeapiextensionsclient.NewSimpleClientset()
	s := scheme.Scheme
	kubeClient := fakekubernetes.NewSimpleClientset(pullSecret.DeepCopyObject())
	kubeClientNoPullSecret := fakekubernetes.NewSimpleClientset()
	dynamicClient := fakedynamic.NewSimpleDynamicClient(s)
	discoveryClient := kubeClient.Discovery()
	discoveryClient.(*fakediscovery.FakeDiscovery).Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{
					Name:       "secrets",
					Namespaced: true,
					Kind:       "Secret",
				},
			},
		},
		{
			GroupVersion: "cluster.open-cluster-management.io/v1",
			APIResources: []metav1.APIResource{
				{
					Name:       "managedclusters",
					Namespaced: false,
					Kind:       "ManagedCluster",
				},
			},
		},
		{
			GroupVersion: "agent.open-cluster-management.io/v1",
			APIResources: []metav1.APIResource{
				{
					Name:       "klusteraddonconfigs",
					Namespaced: false,
					Kind:       "KlusterletAddonConfig",
				},
			},
		},
		{
			GroupVersion: "hive.openshift.io/v1",
			APIResources: []metav1.APIResource{
				{
					Name:       "machinepools",
					Namespaced: false,
					Kind:       "MachinePool",
				},
				{
					Name:       "clusterimagesets",
					Namespaced: false,
					Kind:       "ClusterImageSet",
				},
				{
					Name:       "clusterdeployments",
					Namespaced: false,
					Kind:       "ClusterDeployment",
				},
			},
		},
	}
	clusterClient := clusterclientset.NewSimpleClientset()
	workClient := workclientset.NewSimpleClientset()
	type fields struct {
		CMFlags     *genericclioptionscm.CMFlags
		clusterName string
		cloud       string
		valuesPath  string
		values      map[string]interface{}
		outputFile  string
	}
	type args struct {
		kubeClient          kubernetes.Interface
		dynamicClient       dynamic.Interface
		apiextensionsClient apiextensionsclient.Interface
		discoveryClient     discovery.DiscoveryInterface
		clusterClient       *clusterclientset.Clientset
		workClient          *workclientset.Clientset
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
				CMFlags: genericclioptionscm.NewCMFlags(nil),
				values:  values,
				cloud:   "aws",
			},
			args: args{
				kubeClient:          kubeClient,
				discoveryClient:     discoveryClient,
				apiextensionsClient: apiextensionsClient,
				dynamicClient:       dynamicClient,
				clusterClient:       clusterClient,
				workClient:          workClient,
			},
			wantErr: false,
		},
		{
			name: "Failed no pullsecret",
			fields: fields{
				CMFlags: genericclioptionscm.NewCMFlags(nil),
				values:  values,
				cloud:   "aws",
			},
			args: args{
				kubeClient:          kubeClientNoPullSecret,
				discoveryClient:     discoveryClient,
				apiextensionsClient: apiextensionsClient,
				dynamicClient:       dynamicClient,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				clusterName: tt.fields.clusterName,
				cloud:       tt.fields.cloud,
				valuesPath:  tt.fields.valuesPath,
				values:      tt.fields.values,
				outputFile:  tt.fields.outputFile,
			}
			if err := o.runWithClient(tt.args.kubeClient,
				tt.args.apiextensionsClient,
				tt.args.dynamicClient,
				tt.args.clusterClient,
				tt.args.workClient); (err != nil) != tt.wantErr {
				t.Errorf("Options.runWithClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
