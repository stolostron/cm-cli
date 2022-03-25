[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add `cm delete clusterset`, inhirated from clusteradm
- Add `cm delete work`, inhirated from clusteradm
- Add `cm set clusterset`, inhirated from clusteradm `clusteradm clusterset set`
- Add `cm bind clusterset`, inhirated from clusteradm `clusteradm clusterset bind`
- Add `cm unbind clusterset`, inhirated from clusteradm`clusteradm clusterset unbind`
- update to the latest https://github.com/open-cluster-management-io/clusteradm version.
## Breaking changes

- Remove commands "init, join, accept, get token, delete token" as they are pure OCM and should not be used on ACM/MCE. Please use commands like `cm attach cluster`.
## Bug fixes


