// Copyright Contributors to the Open Cluster Management project

package clusterpoolhost

import (
	"fmt"
	"os"

	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	KubeConfigIgnoredMessage = "WARNING: KUBECONFIG is set and is being ignored by context name=%s\n"
	KubeConfigSwitchMessage  = "Switching to context %s\n"
)

//IsGlobalContext checks if the context is in the global file
func IsGlobalContext(contextName string) (bool, error) {
	return isContextExists(contextName, true)
}

//IsContext check if the context is in the file specified by the env var if set otherwise the global file
func IsContext(contextName string) (bool, error) {
	return isContextExists(contextName, false)
}

func isContextExists(contextName string, globalKubeConfig bool) (bool, error) {
	if len(contextName) == 0 {
		return false, fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	if globalKubeConfig {
		pathOptions.EnvVar = ""
	}
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return false, err
	}
	_, ok := config.Contexts[contextName]
	return ok, nil
}

//MoveContextToDefault Move the context from its current location to the global file.
func MoveContextToDefault(currentContextName, clusterPoolContextName, defaultNamespace, user, token string) error {
	if len(currentContextName) == 0 {
		return fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	//The cph context is already in the config, no move needed
	if _, ok := config.Contexts[clusterPoolContextName]; ok {
		config.CurrentContext = clusterPoolContextName
		//refesh token
		authInfo := config.AuthInfos[clusterPoolContextName]
		authInfo.Token = token

		file := pathOptions.GetDefaultFilename()
		return clientcmd.WriteToFile(*config, file)
	}

	//Search the context used for `oc login` while creating the cph
	context, ok := config.Contexts[currentContextName]
	if !ok {
		//Search in Globalfile
		pathOptions.EnvVar = ""
		if config, err = pathOptions.GetStartingConfig(); err != nil {
			return err
		}
		if context, ok = config.Contexts[currentContextName]; !ok {
			return fmt.Errorf("context name %s not found", currentContextName)
		}
	}
	cluster, ok := config.Clusters[context.Cluster]
	if !ok {
		return fmt.Errorf("cluster not found for context %s", currentContextName)
	}
	authInfo, ok := config.AuthInfos[context.AuthInfo]
	if !ok {
		return fmt.Errorf("authInfo not found for context %s", currentContextName)
	}

	pathOptions = clientcmd.NewDefaultPathOptions()
	pathOptions.EnvVar = ""

	if _, err := os.Stat(pathOptions.GetDefaultFilename()); os.IsNotExist(err) {
		_, err := os.Create(pathOptions.GetDefaultFilename())
		if err != nil {
			return err
		}
	}

	config, err = pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	config.Clusters[clusterPoolContextName] = cluster
	context.AuthInfo = clusterPoolContextName
	context.Namespace = defaultNamespace
	context.Cluster = clusterPoolContextName
	config.CurrentContext = clusterPoolContextName
	config.Contexts[clusterPoolContextName] = context
	authInfo.Token = token
	config.AuthInfos[clusterPoolContextName] = authInfo

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       config.Clusters,
		Contexts:       config.Contexts,
		CurrentContext: config.CurrentContext,
		AuthInfos:      config.AuthInfos,
	}
	file := pathOptions.GetDefaultFilename()
	return clientcmd.WriteToFile(clientConfig, file)
}

//CreateContextFronConfigAPI creates a new context in the global file
func CreateContextFronConfigAPI(configAPI *clientcmdapi.Config, token, contextName, defaultNamespace, user string, setAsCurrent bool) error {
	if len(contextName) == 0 {
		return fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	if _, err := os.Stat(pathOptions.GetDefaultFilename()); os.IsNotExist(err) {
		_, err := os.Create(pathOptions.GetDefaultFilename())
		if err != nil {
			return err
		}
	}
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}
	configAPI.AuthInfos["admin"].Token = token
	configAPI.AuthInfos["admin"].ClientKeyData = nil
	configAPI.AuthInfos["admin"].ClientCertificateData = nil

	contextConfigAPI := configAPI.Contexts["admin"]
	config.AuthInfos[contextName] = configAPI.AuthInfos[contextConfigAPI.AuthInfo]
	config.Clusters[contextName] = configAPI.Clusters[contextConfigAPI.Cluster]
	configAPI.Contexts["admin"].Namespace = contextConfigAPI.Namespace
	config.Contexts[contextName] = configAPI.Contexts["admin"]
	config.Contexts[contextName].AuthInfo = contextName
	config.Contexts[contextName].Cluster = contextName
	if setAsCurrent {
		config.CurrentContext = contextName
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       config.Clusters,
		Contexts:       config.Contexts,
		CurrentContext: config.CurrentContext,
		AuthInfos:      config.AuthInfos,
	}
	file := pathOptions.GetDefaultFilename()
	return clientcmd.WriteToFile(clientConfig, file)
}

//GetGlobalConfigAPI returns the Global ConfigAPI and if the KUBECONFIG was set.
func GetGlobalConfigAPI() (*clientcmdapi.Config, bool, error) {
	return getConfigAPI(true)
}

//GetConfigAPI returns the ConfigAPI and if the KUBECONFIG was set.
func GetConfigAPI() (*clientcmdapi.Config, bool, error) {
	return getConfigAPI(false)
}

func getConfigAPI(globalKubeConfig bool) (*clientcmdapi.Config, bool, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()
	isEnvVarSet := os.Getenv(pathOptions.EnvVar) != ""
	if globalKubeConfig {
		pathOptions.EnvVar = ""
	}
	configapi, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, isEnvVarSet, err
	}
	return configapi, isEnvVarSet, nil
}

//GetGlobalCurrentRestConfig gets the *rest.Config of the current context in the global file.
func GetGlobalCurrentRestConfig() (*rest.Config, error) {
	return getCurrentRestConfig(true)
}

//GetCurrentRestConfig gets the *rest.Config of the current context in the file specified by the env var if set.
func GetCurrentRestConfig() (*rest.Config, error) {
	return getCurrentRestConfig(false)
}

func getCurrentRestConfig(globalKubeConfig bool) (*rest.Config, error) {
	configapi, _, err := getConfigAPI(globalKubeConfig)
	if err != nil {
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*configapi, nil)
	return clientConfig.ClientConfig()
}

//GetGlobalRestConfig gets the *rest.Config of the current context in the global file.
func (cph *ClusterPoolHost) GetGlobalRestConfig() (*rest.Config, error) {
	return cph.getRestConfig(true)
}

//GetRestConfig gest the *rest.Config of the current context in the file specified by the env var if set.
func (cph *ClusterPoolHost) GetRestConfig() (*rest.Config, error) {
	return cph.getRestConfig(false)
}

func (cph *ClusterPoolHost) getRestConfig(globalKubeConfig bool) (*rest.Config, error) {
	configapi, _, err := getConfigAPI(globalKubeConfig)
	if err != nil {
		return nil, err
	}
	configapi.CurrentContext = cph.GetContextName()
	clientConfig := clientcmd.NewDefaultClientConfig(*configapi, nil)

	// clientConfig = clientcmd.NewDefaultClientConfig(rawConfig, &clientcmd.ConfigOverrides{})
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = helpers.QPS
	config.Burst = helpers.Burst

	return config, nil

}
