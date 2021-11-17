## cm create clusterclaim

Create clusterclaims

```
cm create clusterclaim [flags]
```

### Examples

```

# Create clusterclaims in the current clusterpoolhost
cm create cc <clusterpool> <clusterclaim_name>[,<clusterclaim_name>...] <options>

# Create clusterclaims on a given clusterpoolhost
cm create cc <clusterpool> <clusterclaim_name>[,<clusterclaim_name>...] --cph <clusterpoolhost> <options>

```

### Options

```
      --allow-missing-template-keys   If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      --cph string                    The clusterpoolhost to use
      --creds                         If set the credentials will be displayed
  -h, --help                          help for clusterclaim
      --import                        If set the clusterclaim will be imported
  -L, --label-columns strings         Accepts a comma separated list of labels that are going to be presented as columns. Names are case-sensitive. You can also use multiple flag options like -L label1 -L label2...
      --no-headers                    When using the default or custom-column output format, don't print headers (default print headers).
  -o, --output string                 Output format. One of: json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-as-json|jsonpath-file|custom-columns-file|custom-columns|wide See custom columns [https://kubernetes.io/docs/reference/kubectl/overview/#custom-columns], golang template [http://golang.org/pkg/text/template/#pkg-overview] and jsonpath template [https://kubernetes.io/docs/reference/kubectl/jsonpath/].
      --output-file string            The generated resources will be copied in the specified file
      --show-kind                     If present, list the resource type for the requested object(s).
      --show-labels                   When printing, show all labels as the last column (default hide labels column)
      --show-managed-fields           If true, keep the managedFields when printing objects in JSON or YAML format.
      --skip-schedule                 Set the hibernation schedule to skip
      --sort-by string                If non-empty, sort list types using this field specification.  The field specification is expressed as a JSONPath expression (e.g. '{.metadata.name}'). The field in the API resource specified by this JSONPath expression must be an integer or a string.
      --template string               Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --timeout int                   Timeout to get the cluster claim running in minutes (default 60)
```

### Options inherited from parent commands

```
      --add-dir-header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files
      --as string                        Username to impersonate for the operation
      --as-group stringArray             Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
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
      --skip-headers                     If true, avoid header prefixes in the log messages
      --skip-log-headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --tls-server-name string           Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                     Bearer token for authentication to the API server
      --user string                      The name of the kubeconfig user to use
      --username string                  Username for basic authentication to the API server
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [cm create](cm_create.md)	 - create a resource

