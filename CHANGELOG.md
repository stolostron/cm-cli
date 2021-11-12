[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions
- Display message "clusterclaim xx is running" only when it wasn't running before.
- Add support for listing a single clusterpool `cm get cp <cluster_poolname>`
- Update authrealm documentation help.
## Breaking changes

## Bug fixes
- [Klusteraddonconfig error when attach cluster on MCE #143](https://github.com/open-cluster-management/cm-cli/issues/143)
- Fix set timeout 0 when running `cm get clusterclaim`
- timeout n > 0 return in n minutes rather than n+2
