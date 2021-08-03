// Copyright Contributors to the Open Cluster Management project

package clusterpoolhost

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	printclusterpoolv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ClusterPoolHostsDir = ".kube"
const ServiceAccountNameSpace = "default"

const (
	ClusterPoolHostContextPrefix string = "clusterpoolhost"
)

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

type ErrorType string

const (
	errorNotFound ErrorType = "not found"
)

//GetContextName returns the context name for a given clusterpoolhost
func (c *ClusterPoolHost) GetContextName() string {
	u, err := url.Parse(c.APIServer)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s/%s/%s", ClusterPoolHostContextPrefix, c.Namespace, u.Hostname(), c.Name)
}

//GetClusterPoolHosts returns all clusterpoolhosts
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

//IsClusterPoolHost checks if the provided context name is a clusterpoolhost context
func IsClusterPoolHost(contextName string) (bool, error) {
	cphs, err := GetClusterPoolHosts()
	if err != nil {
		return false, err
	}
	_, ok := cphs.ClusterPoolHosts[contextName]
	return ok, nil
}

//GetClusterPoolHost returns the clusterpoolhost corresponding to the provided name
func GetClusterPoolHost(clusterPoolHostName string) (*ClusterPoolHost, error) {
	cphs, err := GetClusterPoolHosts()
	if err != nil {
		return nil, err
	}
	if c, ok := cphs.ClusterPoolHosts[clusterPoolHostName]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("%s %s", clusterPoolHostName, errorNotFound)
}

//ApplyClusterPoolHosts saves the list of clusterpoolhost
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

//OpenClusterPoolHost returns the clusterpoolhost
func OpenClusterPoolHost(consoleUrl string) error {
	return helpers.Openbrowser(consoleUrl)
}

//GetClusterPoolHost opens a browzer on the clusterpoolhost console
func (cs *ClusterPoolHosts) GetClusterPoolHost(name string) (*ClusterPoolHost, error) {
	if c, ok := cs.ClusterPoolHosts[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("cluster pool host %s not found", name)
}

var (
	ClusterPoolHostsColumns string = "custom-columns=ACTIVE:.spec.active,CLUSTER_POOL_HOST:.metadata.name,NAMESPACE:.spec.namespace,API_SERVER:.spec.apiServer,CONSOLE:.spec.console"
)

func ConvertToPrintClusterPoolHostList(cphs *ClusterPoolHosts) *printclusterpoolv1alpha1.PrintClusterPoolHostList {
	pcps := &printclusterpoolv1alpha1.PrintClusterPoolHostList{}
	for i := range cphs.ClusterPoolHosts {
		pcp := printclusterpoolv1alpha1.PrintClusterPoolHost{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cphs.ClusterPoolHosts[i].Name,
				Namespace: cphs.ClusterPoolHosts[i].Namespace,
			},
			Spec: printclusterpoolv1alpha1.PrintClusterPoolHostSpec{
				Active:    cphs.ClusterPoolHosts[i].Active,
				Name:      cphs.ClusterPoolHosts[i].Name,
				Namespace: cphs.ClusterPoolHosts[i].Namespace,
				APIServer: cphs.ClusterPoolHosts[i].APIServer,
				Console:   cphs.ClusterPoolHosts[i].Console,
				Group:     cphs.ClusterPoolHosts[i].Group,
			},
		}
		pcps.Items = append(pcps.Items, pcp)
	}
	return pcps
}

//UnActiveAll unactives all clusterpoolhosts
func (cs *ClusterPoolHosts) UnActiveAll() error {
	for _, c := range cs.ClusterPoolHosts {
		c.Active = false
	}
	return cs.ApplyClusterPoolHosts()
}

//SetActive actives a specific clusterpoolhost
func (cs *ClusterPoolHosts) SetActive(c *ClusterPoolHost) error {
	if err := cs.UnActiveAll(); err != nil {
		return err
	}
	c.Active = true
	return cs.ApplyClusterPoolHosts()
}

//GetCurrentClusterPoolHost gets the current clusterpoolhost
func (cs *ClusterPoolHosts) GetCurrentClusterPoolHost() (*ClusterPoolHost, error) {
	for _, c := range cs.ClusterPoolHosts {
		if c.IsActive() {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no active cluster pool host found")
}

//AddClusterPoolHost adds a clusterpoolhost
func (c *ClusterPoolHost) AddClusterPoolHost() error {
	cs, err := GetClusterPoolHosts()
	if err != nil {
		return err
	}
	cs.ClusterPoolHosts[c.Name] = c
	return cs.ApplyClusterPoolHosts()
}

//DeleteClusterPoolHost deletes a clusterpoolhost
func (c *ClusterPoolHost) DeleteClusterPoolHost() error {
	cs, err := GetClusterPoolHosts()
	if err != nil {
		return err
	}
	delete(cs.ClusterPoolHosts, c.Name)
	return cs.ApplyClusterPoolHosts()
}

//IsActive checks if clusterpoolhost is active
func (c *ClusterPoolHost) IsActive() bool {
	return c.Active
}

//GetCurrentClusterPoolHost gets the current active clusterpoolhost
func GetCurrentClusterPoolHost() (*ClusterPoolHost, error) {
	cs, err := GetClusterPoolHosts()
	if err != nil {
		return nil, err
	}
	return cs.GetCurrentClusterPoolHost()
}
