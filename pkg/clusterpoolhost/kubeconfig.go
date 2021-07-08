// Copyright Contributors to the Open Cluster Management project

package clusterpoolhost

import (
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	KubeConfigIgnoredMessage = "WARNING: KUBECONFIG is set and is being ignored by context name=%s\n"
	KubeConfigSwitchMessage  = "Switching to context %s\n"
)

func GetGlobalConfig() (*clientcmdapi.Config, error) {
	return getConfig(true)
}

func GetConfig() (*clientcmdapi.Config, error) {
	return getConfig(false)
}

func getConfig(globalKubeConfig bool) (*clientcmdapi.Config, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if globalKubeConfig {
		pathOptions.EnvVar = ""
	}
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func IsGlobalContext(contextName string) (bool, error) {
	return isContextExists(contextName, true)
}

func IsContext(contextName string) (bool, error) {
	return isContextExists(contextName, false)
}

func isContextExists(contextName string, globalKubeConfig bool) (bool, error) {
	if len(contextName) == 0 {
		return false, fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	if globalKubeConfig {
		fmt.Printf(KubeConfigIgnoredMessage, contextName)
		pathOptions.EnvVar = ""
	}
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return false, err
	}
	_, ok := config.Contexts[contextName]
	return ok, nil
}

func SetGlobalCurrentContext(contextName string) error {
	return setCurrentContext(contextName, true)
}

func SetCurrentContext(contextName string) error {
	return setCurrentContext(contextName, false)
}

func setCurrentContext(contextName string, globalKubeConfig bool) error {
	if len(contextName) == 0 {
		return fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	fmt.Printf(KubeConfigIgnoredMessage, contextName)
	if globalKubeConfig {
		pathOptions.EnvVar = ""
	}
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	_, ok := config.Contexts[contextName]
	if !ok {
		return fmt.Errorf("context name %s not found", contextName)
	}
	config.CurrentContext = contextName
	fmt.Printf(KubeConfigSwitchMessage, contextName)
	return clientcmd.ModifyConfig(pathOptions, *config, true)

}

func MoveContextToDefault(contextName, clusterPoolContextName, defaultNamespace, user, token string) error {
	if len(contextName) == 0 {
		return fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	context, ok := config.Contexts[contextName]
	if !ok {
		//Search in Globalfile
		pathOptions.EnvVar = ""
		if config, err = pathOptions.GetStartingConfig(); err != nil {
			return err
		}
		if context, ok = config.Contexts[contextName]; !ok {
			return fmt.Errorf("context name %s not found", contextName)
		}
	}
	cluster, ok := config.Clusters[context.Cluster]
	if !ok {
		return fmt.Errorf("cluster not found for context %s", contextName)
	}
	authInfo, ok := config.AuthInfos[context.AuthInfo]
	if !ok {
		return fmt.Errorf("authInfo not found for context %s", contextName)
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

	delete(config.Clusters, context.Cluster)
	delete(config.Contexts, contextName)
	delete(config.AuthInfos, context.AuthInfo)

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

func CreateContextFronConfigAPI(configAPI *clientcmdapi.Config, token, contextName, defaultNamespace, user string) error {
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
	config.CurrentContext = contextName

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

func GetGlobalConfigAPI() (*clientcmdapi.Config, error) {
	return getConfigAPI(true)
}

func GetConfigAPI() (*clientcmdapi.Config, error) {
	return getConfigAPI(false)
}

func getConfigAPI(globalKubeConfig bool) (*clientcmdapi.Config, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if globalKubeConfig {
		pathOptions.EnvVar = ""
	}
	configapi, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	return configapi, nil
}

func GetGlobalCurrentRestConfig() (*rest.Config, error) {
	return getCurrentRestConfig(true)
}

func GetCurrentRestConfig() (*rest.Config, error) {
	return getCurrentRestConfig(false)
}

func getCurrentRestConfig(globalKubeConfig bool) (*rest.Config, error) {
	configapi, err := getConfigAPI(globalKubeConfig)
	if err != nil {
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*configapi, nil)
	return clientConfig.ClientConfig()
}

func FindConfigAPIByAPIServer(contextName, apiServer string) (bool, error) {
	configAPI, err := GetConfigAPI()
	if err != nil {
		return false, err
	}
	if len(configAPI.CurrentContext) == 0 {
		return false, fmt.Errorf("no current context")
	}

	if _, ok := configAPI.Contexts[contextName]; ok {
		if err := SetCurrentContext(contextName); err != nil {
			return false, err
		}
		return IsGlobalContext(configAPI.CurrentContext)
	}

	//Search for the cluster in kubeconfig.clusters
	var foundCluster string
	for clusterName, cluster := range configAPI.Clusters {
		if cluster.Server == apiServer && !strings.HasPrefix(clusterName, ClusterPoolHostContextPrefix) {
			foundCluster = clusterName
			break
		}
	}
	if len(foundCluster) == 0 {
		return false, fmt.Errorf("not found %s as current context", apiServer)
	}

	//Search for the found cluster in the kubeconfig.context.
	var foundContext string
	for contextName, context := range configAPI.Contexts {
		if context.Cluster == foundCluster {
			foundContext = contextName
			break
		}
	}
	if len(foundContext) == 0 {
		return false, fmt.Errorf("not found %s as current context", apiServer)
	}
	err = SetCurrentContext(foundContext)
	if err != nil {
		return false, err
	}
	return IsGlobalContext(configAPI.CurrentContext)
}
