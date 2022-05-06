[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add "--current" in `cm get cc` to display the current clusterclaim in use.
- Create binary for darwin/arm64.
- [Support import of a cluster that is the non-active kubecontext in a kubeconfig #197](https://github.com/stolostron/cm-cli/issues/197)
- Add "--kubeconfig-creds" in `cm get cc` to display the kubeconfig
- [Open openshift console of an hypershift hosted #201](https://github.com/stolostron/cm-cli/issues/201)
- [Add cmd cm console cluster to open the console of a managedcluster either created by cc or hypershift #200](https://github.com/stolostron/cm-cli/issues/200)
## Breaking changes

- Remove commands "init, join, accept, get token, delete token" as they are pure OCM and should not be used on ACM/MCE. Please use commands like `cm attach cluster`.
## Bug fixes

- Use "--server-namespace" when checking if it is an MCE hub.