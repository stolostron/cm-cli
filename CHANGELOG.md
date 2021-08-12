[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- [Create a cmd which generate a config to create a clusterpool based on an existing clusterpool #66](https://github.com/open-cluster-management/cm-cli/issues/66)
- Create help.tar.gz and help.zip files contains command-line markdown help

## Breacking changes

- As the project leverages the [printers](https://github.com/kubernetes/cli-runtime/blob/master/pkg/printers/interface.go) the output format might change. 
## Bug fixes

- [Increase QPS as lots of throttling message #68](https://github.com/open-cluster-management/cm-cli/issues/68)
- [cm console cc doesn't take into account the --cph #95](https://github.com/open-cluster-management/cm-cli/issues/95)