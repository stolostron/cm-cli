## cm get

get a resource

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
      --add-dir-header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files
      --as string                        Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray             Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                    UID to impersonate for the operation.
      --beta                             If set commands or functionalities in beta version will be available
      --cache-dir string                 Default cache directory (default "${HOME}/.kube/cache")
      --certificate-authority string     Path to a cert file for the certificate authority
      --client-certificate string        Path to a client certificate file for TLS
      --client-key string                Path to a client key file for TLS
      --cluster string                   The name of the kubeconfig cluster to use
      --context string                   The name of the kubeconfig context to use
      --dry-run                          If set the generated resources will be displayed but not applied
      --insecure-skip-tls-verify         If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string                Path to the kubeconfig file to use for CLI requests.
      --log-backtrace-at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                   If non-empty, write log files in this directory
      --log-file string                  If non-empty, use this log file
      --log-file-max-size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --match-server-version             Require server version to match client version
  -n, --namespace string                 If present, the namespace scope for this CLI request
      --one-output                       If true, only write logs to their native severity level (vs also writing to each lower severity level)
      --password string                  Password for basic authentication to the API server
      --request-timeout string           The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                    The address and port of the Kubernetes API server
      --server-namespace string          The namespace where the server (RHACM/MCE) is installed
      --skip-headers                     If true, avoid header prefixes in the log messages
      --skip-log-headers                 If true, avoid headers when opening log files
      --skip-server-check                If set commands will not check the installed server (RHACM/MCE) target
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --tls-server-name string           Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                     Bearer token for authentication to the API server
      --user string                      The name of the kubeconfig user to use
      --username string                  Username for basic authentication to the API server
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [cm](cm.md)	 - CLI for Red Hat Advanced Cluster Management
* [cm get addon](cm_get_addon.md)	 - get enabled addon on specified managed cluster
* [cm get addon](cm_get_addon.md)	 - disable specified addon on specified managed clusters
* [cm get clusterclaims](cm_get_clusterclaims.md)	 - Display clusterclaims
* [cm get clusterpoolhosts](cm_get_clusterpoolhosts.md)	 - list the clusterpoolhosts
* [cm get clusterpools](cm_get_clusterpools.md)	 - Get clusterpool
* [cm get clusters](cm_get_clusters.md)	 - Display the attached clusters
* [cm get clustersets](cm_get_clustersets.md)	 - get clustersets
* [cm get components](cm_get_components.md)	 - Get the list of available components
* [cm get config](cm_get_config.md)	 - get the config of a resource
* [cm get contexts](cm_get_contexts.md)	 - Get the managedcluster's contexts of a hub
* [cm get credentials](cm_get_credentials.md)	 - list the credentials of cloud providers
* [cm get hub-info](cm_get_hub-info.md)	 - get hub-info
* [cm get machinepools](cm_get_machinepools.md)	 - list the machinepools for a give cluster
* [cm get policies](cm_get_policies.md)	 - Display policies
* [cm get works](cm_get_works.md)	 - get manifestwork on a specified managed cluster

