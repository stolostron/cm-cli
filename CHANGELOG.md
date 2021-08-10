[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Create a cmd which generate a config to create a clusterpool based on an existing clusterpool [#66](https://github.com/open-cluster-management/cm-cli/issues/66)
- Add `get policies` command [#98](https://github.com/open-cluster-management/cm-cli/pull/98)

## Breaking changes

- As the project leverages the [printers](https://github.com/kubernetes/cli-runtime/blob/master/pkg/printers/interface.go) the output format might change.
- Users without clusterwide access to policies will not be able to list policies. [#100](https://github.com/open-cluster-management/cm-cli/issues/100)

## Bug fixes

- [Increase QPS as lots of throttling message #68](https://github.com/open-cluster-management/cm-cli/issues/68)
- When `-o wide` is specified for resources using a printer CRD only Name and Age columns display