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

var contextBackup, globalContextBackup string

//GetGlobalConfig gets the config from the global file
func GetGlobalConfig() (*clientcmdapi.Config, error) {
	return getConfig(true)
}

//GetConfig gets the config from the file specified by the env var if set otherwise the global file
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

//SetGlobalCurrentContext sets the current context in the global file
func SetGlobalCurrentContext(contextName string) error {
	return setCurrentContext(contextName, true)
}

//SetCurrentContext sets the current context in the file specified by the env var if set otherwise in the global file.
func SetCurrentContext(contextName string) error {
	return setCurrentContext(contextName, false)
}

func setCurrentContext(contextName string, globalKubeConfig bool) error {
	if len(contextName) == 0 {
		return fmt.Errorf("context name is empty")
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	if globalKubeConfig && os.Getenv(pathOptions.EnvVar) != "" {
		fmt.Printf(KubeConfigIgnoredMessage, contextName)
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
	if config.CurrentContext != contextName {
		config.CurrentContext = contextName
		fmt.Printf(KubeConfigSwitchMessage, contextName)
		return clientcmd.ModifyConfig(pathOptions, *config, true)
	}
	return nil

}

//MoveContextToDefault Move the context from its current location to the global file.
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

//CreateContextFronConfigAPI creates a new context in the global file
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

//GetCurrentRestConfig gest the *rest.Config of the current context in the file specified by the env var if set.
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

//SetCPHContext sets the clusterpoolhost context as current
func SetCPHContext(contextName string) error {
	if strings.HasPrefix(contextName, ClusterPoolHostContextPrefix) {
		IsGlobalContext, err := IsGlobalContext(contextName)
		if err != nil {
			return err
		}
		if IsGlobalContext {
			if err := SetGlobalCurrentContext(contextName); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("%s is not cph context", contextName)
}

func findConfigAPIByAPIServer(contextName, apiServer string) (bool, error) {
	configAPI, isEnvVarSet, err := GetConfigAPI()
	if err != nil {
		return false, err
	}
	isGlobal, foundContxt, err := findConfigAPIByAPIServerForConfig(configAPI, contextName, apiServer)
	if err == nil {
		if err := setCurrentContext(foundContxt, isGlobal); err != nil {
			return false, err
		}
		return isGlobal, err
	}
	if isEnvVarSet {
		configAPI, _, err = GetGlobalConfigAPI()
		if err != nil {
			return isGlobal, err
		}
		isGlobal, foundContxt, err = findConfigAPIByAPIServerForConfig(configAPI, contextName, apiServer)
		if err == nil {
			if err := setCurrentContext(foundContxt, isGlobal); err != nil {
				return false, err
			}
			return isGlobal, err
		}
	}
	return isGlobal, err
}

func findConfigAPIByAPIServerForConfig(configAPI *clientcmdapi.Config, contextName, apiServer string) (bool, string, error) {
	if len(configAPI.CurrentContext) == 0 {
		return false, "", fmt.Errorf("no current context")
	}

	//Search for the cluster in kubeconfig.clusters
	var foundCluster string
	for clusterName, cluster := range configAPI.Clusters {
		if cluster.Server == apiServer {
			foundCluster = clusterName
			break
		}
	}
	if len(foundCluster) == 0 {
		return false, "", fmt.Errorf("not found %s as current context", apiServer)
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
		return false, "", fmt.Errorf("not found %s as current context", apiServer)
	}
	isGlobal, err := IsGlobalContext(configAPI.CurrentContext)
	return isGlobal, foundContext, err
}

//BackupCurrentContexts backups the names for the current contexts, the context in the
//global file and the one defined in the env var if set.
func BackupCurrentContexts() (err error) {
	configAPI, isEnvVarSet, err := GetConfigAPI()
	if err != nil {
		return
	}
	if isEnvVarSet {
		contextBackup = configAPI.CurrentContext
		configAPI, _, err = GetGlobalConfigAPI()
		if err != nil {
			return
		}
	}
	globalContextBackup = configAPI.CurrentContext
	return
}

//RestoreCurrentContexts restores the backuped contexts.
func RestoreCurrentContexts() error {
	if len(contextBackup) != 0 {
		if err := setCurrentContext(contextBackup, false); err != nil {
			return err
		}
	}
	if len(globalContextBackup) != 0 {
		if err := setCurrentContext(globalContextBackup, true); err != nil {
			return err
		}
	}
	return nil
}
