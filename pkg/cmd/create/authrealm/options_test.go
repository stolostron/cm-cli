// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	"reflect"
	"testing"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Test_newOptions(t *testing.T) {
	cmFlags := genericclioptionscm.NewCMFlags(nil)
	type args struct {
		cmFlags *genericclioptionscm.CMFlags
		streams genericclioptions.IOStreams
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "success",
			args: args{
				cmFlags: cmFlags,
			},
			want: &Options{
				CMFlags: cmFlags,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptions(tt.args.cmFlags, tt.args.streams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
