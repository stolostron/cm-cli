[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add `--current` as parameter of `cm get contexts` to generate a kubeconfig containing only the current used context.
- Generate secret token due to automatic token generation deprecation on kubernetes 1.24

## Breaking changes

- `cm create hd` template changed to a more structured template.

## Bug fixes

