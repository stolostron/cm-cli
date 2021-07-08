// Copyright Contributors to the Open Cluster Management project

package clusterpoolhost

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/rest"

	"github.com/ghodss/yaml"
	userv1 "github.com/openshift/api/user/v1"
	userv1typedclient "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ClusterPoolHostsDir = ".kube"
const ServiceAccountNameSpace = "default"

var (
	clusterPoolClustersFile = filepath.Join(ClusterPoolHostsDir, "known-cphs")
)

type ClusterPoolHosts struct {
	ClusterPoolHosts map[string]*ClusterPoolHost `json:"clusters"`
}

type ClusterPoolHost struct {
	// Name of the cluster pool
	Name string `json:"name"`
	// true if this cluster pool is the Active one
	Active bool `json:"active"`
	// The API address of the cluster where your `ClusterPools` are defined. Also referred to as the "ClusterPool host"
	APIServer string `json:"apiServer"`
	// The URL of the OpenShift console for the ClusterPool host
	Console string `json:"console"`
	// Namespace where `ClusterPools` are defined
	Namespace string `json:"namespace"`
	// Name of a `Group` (`user.openshift.io/v1`) that should be added to each `ClusterClaim` for team access
	Group string `json:"group"`
}

type ClusterPoolHostError struct {
	Name string
	Err  error
}

type ErrorType string

const (
	errorNotFound ErrorType = "not found"
)

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	switch v := err.(type) {
	case *ClusterPoolHostError:
		return strings.Contains(v.Err.Error(), string(errorNotFound))
	}
	return false
}

func (err *ClusterPoolHostError) Error() string {
	return err.Name + ":" + err.Err.Error()
}

func newError(name string, err error) *ClusterPoolHostError {
	return &ClusterPoolHostError{
		Name: name,
		Err:  err,
	}
}

func (c *ClusterPoolHost) WhoAmI(restConfig *rest.Config) (*userv1.User, error) {
	userInterface, err := userv1typedclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	me, err := userInterface.Users().Get(context.TODO(), "~", metav1.GetOptions{})
	if err == nil {
		return me, err
	}

	return me, err
}

func (c *ClusterPoolHost) GetContextName() string {
	u, err := url.Parse(c.APIServer)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s/%s/%s", ClusterPoolHostContextPrefix, c.Namespace, u.Hostname(), c.Name)
}

func GetClusterPoolHosts() (*ClusterPoolHosts, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	fileName := filepath.Clean(filepath.Join(home, clusterPoolClustersFile))
	cpc := &ClusterPoolHosts{}
	cpc.ClusterPoolHosts = make(map[string]*ClusterPoolHost)
	if _, err := os.Stat(fileName); err != nil {
		return cpc, nil
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, cpc)
	if err != nil {
		return nil, err
	}
	return cpc, nil
}

func IsClusterPoolHost(contextName string) (bool, error) {
	cphs, err := GetClusterPoolHosts()
	if err != nil {
		return false, err
	}
	_, ok := cphs.ClusterPoolHosts[contextName]
	return ok, nil
}

func GetClusterPoolHost(clusterPoolHostName string) (*ClusterPoolHost, error) {
	cphs, err := GetClusterPoolHosts()
	if err != nil {
		return nil, err
	}
	if c, ok := cphs.ClusterPoolHosts[clusterPoolHostName]; ok {
		return c, nil
	}
	return nil, newError(clusterPoolHostName, fmt.Errorf("%s", errorNotFound))
}

func (cs *ClusterPoolHosts) ApplyClusterPoolHosts() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fileName := filepath.Clean(filepath.Join(home, clusterPoolClustersFile))
	b, err := yaml.Marshal(cs)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(fileName), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, b, 0600)
}

func (cs *ClusterPoolHosts) GetClusterPoolHost(name string) (*ClusterPoolHost, *ClusterPoolHostError) {
	if c, ok := cs.ClusterPoolHosts[name]; ok {
		return c, nil
	}
	return nil, newError(name, fmt.Errorf("cluster pool host not found"))
}

func (cs *ClusterPoolHosts) RawPrint() error {
	b, err := yaml.Marshal(cs)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (cs *ClusterPoolHosts) Print() {
	for _, c := range cs.ClusterPoolHosts {
		star := " "
		if c.IsActive() {
			star = "*"
		}
		fmt.Printf("%s%s\t%s\n", star, c.Name, c.APIServer)
	}
}

func (cs *ClusterPoolHosts) UnActiveAll() error {
	for _, c := range cs.ClusterPoolHosts {
		c.Active = false
	}
	return cs.ApplyClusterPoolHosts()
}

func (cs *ClusterPoolHosts) SetActive(c *ClusterPoolHost) error {
	if err := cs.UnActiveAll(); err != nil {
		return err
	}
	c.Active = true
	return cs.ApplyClusterPoolHosts()
}

func (c *ClusterPoolHost) AddClusterPoolHost(current bool) error {
	cs, err := GetClusterPoolHosts()
	if err != nil {
		return err
	}
	cs.ClusterPoolHosts[c.Name] = c
	return cs.SetActive(c)
}

func (c *ClusterPoolHost) IsActive() bool {
	return c.Active
}

func GetCurrentClusterPoolHost() (*ClusterPoolHost, error) {
	cs, err := GetClusterPoolHosts()
	if err != nil {
		return nil, err
	}
	for _, c := range cs.ClusterPoolHosts {
		if c.IsActive() {
			return c, nil
		}
	}
	return nil, newError("", fmt.Errorf("active cluster pool host not found"))
}
