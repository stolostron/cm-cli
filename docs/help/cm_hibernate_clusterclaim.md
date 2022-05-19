## cm hibernate clusterclaim

hibernate clusterclaims

```
cm hibernate clusterclaim [flags]
```

### Examples

```

# Hibernate clusterclaims
cm hibernate cc <clusterclaim_name>[,<clusterclaim_name>...] <options>

# run clusterclaims on a given clusterpoolhost
cm hibernate cc <clusterclaim_name>[,<clusterclaim_name>...] --cph <clusterpoolhost> <options>

```

### Options

```
      --cph string               The clusterpoolhost to use
  -h, --help                     help for clusterclaim
      --hibernate-schedule-off   Set the hibernation schedule to off
      --hibernate-schedule-on    Set the hibernation schedule to on
      --output-file string       The generated resources will be copied in the specified file
      --skip-schedule            Set the hibernation schedule to skip (deprecated)
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

* [cm hibernate](cm_hibernate.md)	 - hibernate a resource

