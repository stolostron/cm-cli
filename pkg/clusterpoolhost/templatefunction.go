// Copyright Contributors to the Open Cluster Management project

package clusterpoolhost

import (
	"context"
	"regexp"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	userv1typedclient "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

//ApplierFuncMap adds the function map
func FuncMap() template.FuncMap {
	return template.FuncMap(GenericFuncMap())
}

// GenericFuncMap returns a copy of the basic function map as a map[string]interface{}.
func GenericFuncMap() map[string]interface{} {
	gfm := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}
	return gfm
}

var genericMap = map[string]interface{}{
	"nomalizeName": NormalizeName,
	"getUser":      GetUser,
}

func NormalizeName(name string) string {
	normalizedName := strings.ToLower(name)
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(normalizedName, "")
}

func GetUser(f cmdutil.Factory) (string, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return "", err
	}
	userInterface, err := userv1typedclient.NewForConfig(clientConfig)
	if err != nil {
		return "", err
	}
	me, err := userInterface.Users().Get(context.TODO(), "~", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	b, err := yaml.Marshal(me)
	if err != nil {
		return "", err
	}
	return string(b), nil
	// return NormalizeName(me.FullName), nil
}
