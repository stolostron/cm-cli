[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content

### Additions
### Breacking changes
- Fix [Manual import failed because klusterlet crd not ready yet](https://github.com/open-cluster-management/cm-cli/issues/30). 
Two files are generated instead of one and the command to apply to manual import is like `kubectl apply -f import_crd.yaml;sleep 10;kubectl apply -f import_yaml.yaml`
with `import` being the `--import-file` parameter value.

### Bug fixes

