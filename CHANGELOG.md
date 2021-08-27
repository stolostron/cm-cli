[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add `get policies` command [#98](https://github.com/open-cluster-management/cm-cli/pull/98)
- Add "Lifetime" column to `get clusterclaims` [#115](https://github.com/open-cluster-management/cm-cli/pull/115)

## Breaking changes

- When cm get cc <cc_name> provides the credentials only when flag `--creds` is set [#107](https://github.com/open-cluster-management/cm-cli/issues/107). Now the end-user must explacitly set the `--creds` flag to get the credentials displayed.

## Bug fixes

- When `-o wide` was specified for resources using a printer CRD only Name and Age columns display
- When running cm get, the element are listed in random order [#105](https://github.com/open-cluster-management/cm-cli/issues/105)
- The --cph should not change the current active cph [#102](https://github.com/open-cluster-management/cm-cli/issues/102)
- No active cluster pool hosts [#89](https://github.com/open-cluster-management/cm-cli/issues/89)
- cm create fails when KUBECONFIG is set [#108](https://github.com/open-cluster-management/cm-cli/issues/108)
