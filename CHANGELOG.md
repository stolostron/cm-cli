[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add `--current` as parameter of `cm get contexts` to generate a kubeconfig containing only the current used context.
- Generate secret token due to automatic token generation deprecation on kubernetes 1.24 [Kubernetes 1.24 doesn't automatically generate token anymore. #233](https://github.com/stolostron/cm-cli/issues/233)

## Breaking changes

- `cm create hd` template changed to a more structured template.

## Bug fixes

