// Copyright Contributors to the Open Cluster Management project
package contexts

import (
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {

	// cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	// if err != nil {
	// 	fmt.Println("no clusterpoolhost found, will only get the contexts of hive generated clusters")
	// }

	// return err
	return nil

}
