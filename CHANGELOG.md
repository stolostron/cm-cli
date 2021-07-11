[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content

# Technical Preview
- `create cph`, `get cphs`, `delete cph`, `use-cph`, `set-cph` : [Create a command to register to clusterpool hosts](https://github.com/open-cluster-management/cm-cli/issues/43)
- `use <clusterName>` : [Create a ck use like command](https://github.com/open-cluster-management/cm-cli/issues/32)
- `create <clusterpool> <clusterclaim>[,<clusterclaim>...]`: [Create a command to claim a clusters](https://github.com/open-cluster-management/cm-cli/issues/33)
- `delete <clusterclaim>[,<clusterclaim>...]`: [Create a command to delete a clusters](https://github.com/open-cluster-management/cm-cli/issues/34)
- `get cc [clusterpoolhost|-A]` to get clusterclaims for a given clusterpoolhost or for all of them.
### Additions
- Add the command `delete token` from [clusteradm](https://github.com/open-cluster-management-io/clusteradm)
### Breacking changes
- Fix [Manual import failed because klusterlet crd not ready yet](https://github.com/open-cluster-management/cm-cli/issues/30). 
Two files are generated instead of one and the command to apply to manual import is like `kubectl apply -f import_crd.yaml;sleep 10;kubectl apply -f import_yaml.yaml`
with `import` being the `--import-file` parameter value.

### Bug fixes

