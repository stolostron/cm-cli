[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- Add "--current" in `cm get cc` to display the current clusterclaim in use.
- Create binary for darwin/arm64.

## Breaking changes

- Remove commands "init, join, accept, get token, delete token" as they are pure OCM and should not be used on ACM/MCE. Please use commands like `cm attach cluster`.
## Bug fixes

- Use "--server-namespace" when checking if it is an MCE hub.