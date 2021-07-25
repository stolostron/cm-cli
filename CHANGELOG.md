[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content

- [Add examples in the clusterpool.md or in the help of the command for cm create cph](https://github.com/open-cluster-management/cm-cli/issues/59)
- Add check on RHACM version when running `cm attach clusterclaim` as auto-import is only for >= 2.3
- [open the clusterpoolhost console when adding a cph and not logged in](https://github.com/open-cluster-management/cm-cli/issues/55)
- `cm get credentials` [Create a command to list the provider-connections](https://github.com/open-cluster-management/cm-cli/issues/42)
- `cm console cph <cph_name>` `cm console cc <cc_name> [--cph <cph_name>]` [create a cmd to open the console of a cc or cph](https://github.com/open-cluster-management/cm-cli/issues/56)
- `cm get config cluster <clsuter_name> --output-file <file_name>` [Create a cmd which generate a config to create a cluster based on an existing cluster](https://github.com/open-cluster-management/cm-cli/issues/67)
### Additions

### Breacking changes
### Bug fixes
- [Columns are not well aligned when running cm get cc and cp](https://github.com/open-cluster-management/cm-cli/issues/58)
- [cm get cp -A incomplete when one cphs is hibernating](https://github.com/open-cluster-management/cm-cli/issues/57)
