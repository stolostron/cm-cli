// Copyright Contributors to the Open Cluster Management project
package addons

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	fakeapiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	fakekubernetes "k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubectl/pkg/scheme"
)

var testDir = filepath.Join("test", "unit")

func TestOptions_complete(t *testing.T) {
	type fields struct {
		valuesPath  string
		clusterName string
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
			name:    "Failed, no values.yaml, no name",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Success, not replacing values",
			fields: fields{
				valuesPath: filepath.Join(testDir, "values-with-data.yaml"),
			},
			wantErr: false,
		},
		{
			name: "Success, no values.yaml",
			fields: fields{
				clusterName: "mycluster",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				valuesPath:  tt.fields.valuesPath,
				clusterName: tt.fields.clusterName,
			}
			if err := o.complete(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Options.complete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				imc, ok := o.values["managedCluster"]
				if !ok || imc == nil {
					t.Errorf("missing managedCluster")
				}
				mc := imc.(map[string]interface{})

				if tt.name == "Success, no values.yaml" {
					if mc["name"] != o.clusterName {
						t.Errorf("Expect %s got %s", o.clusterName, mc["name"])
					}
				}
			}
		})
	}
}

func TestEnableAddonsOptions_ValidateWithClient(t *testing.T) {
	s := scheme.Scheme
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
	dynamicClient := fakedynamic.NewSimpleDynamicClient(s)
	type fields struct {
		values            map[string]interface{}
		clusterName       string
		clusterServer     string
		clusterToken      string
		clusterKubeConfig string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success local-cluster, all info in values",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name": "local-cluster",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed local-cluster, cluster name empty",
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
			name: "Failed local-cluster, cluster name missing",
			fields: fields{
				values: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "Success non-local-cluster, overrite cluster-name with local-cluster",
			fields: fields{
				values: map[string]interface{}{
					"managedCluster": map[string]interface{}{
						"name": "test-cluster",
					},
				},

				clusterName: "local-cluster",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				values:      tt.fields.values,
				clusterName: tt.fields.clusterName,
			}
			if err := o.validateWithClient(kubeClient, dynamicClient); (err != nil) != tt.wantErr {
				t.Errorf("EnableAddonsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_runWithClient(t *testing.T) {
	dir, err := ioutil.TempDir(testDir, "tmp")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)
	importSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-import",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"crds.yaml":   []byte("crds: mycrds"),
			"import.yaml": []byte("import: myimport"),
		},
	}
	values, err := helpers.ConvertValuesFileToValuesMap(filepath.Join(testDir, "values-with-data.yaml"), "")
	if err != nil {
		t.Fatal(err)
	}
	s := scheme.Scheme
	kubeClient := fakekubernetes.NewSimpleClientset(importSecret)
	dynamicClient := fakedynamic.NewSimpleDynamicClient(s)
	apiExtensionsClient := fakeapiextensionsclient.NewSimpleClientset()
	discoveryClient := kubeClient.Discovery()
	discoveryClient.(*fakediscovery.FakeDiscovery).Resources = []*metav1.APIResourceList{
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
	}
	type fields struct {
		CMFlags     *genericclioptions.CMFlags
		values      map[string]interface{}
		clusterName string
	}
	type args struct {
		kubeClient          kubernetes.Interface
		apiExtensionsClient apiextensionsclient.Interface
		dynamicClient       dynamic.Interface
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
				CMFlags:     genericclioptions.NewCMFlags(nil),
				values:      values,
				clusterName: "test",
			},
			args: args{
				kubeClient:          kubeClient,
				apiExtensionsClient: apiExtensionsClient,
				dynamicClient:       dynamicClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				CMFlags:     tt.fields.CMFlags,
				values:      tt.fields.values,
				clusterName: tt.fields.clusterName,
			}
			if err := o.runWithClient(tt.args.kubeClient, tt.args.apiExtensionsClient, tt.args.dynamicClient); (err != nil) != tt.wantErr {
				t.Errorf("Options.runWithClient() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
