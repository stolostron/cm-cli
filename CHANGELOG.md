[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Release Content
## Additions

- [Clarify the auth error and login prompts #162](https://github.com/stolostron/cm-cli/issues/162)
- `--skip-server-check` skips now the platform checks and not the only the platform version check.
- Upgrade a number of packages. The main ones are kubebuilder to v0.8.0, clusteradm to v0.1.1-0.20220128120402-ba85108480ae, k8s packages to 0.23.3, cobra to 1.3.0
- Add check on `cm get policies` if the platform is RHACM
- Add `Standby` column in the `cm get clusterpool`
- Display an error if the clusterclaim has no namespace,  this can happen when new clusterclaim with the same name is created just of the deletion of a clusterclaim with the same name. [When clusterclaim not yet running, cm use/run generates an error](https://github.com/stolostron/cm-cli/issues/167)
- Check if Identity Configuration Management is installed when running `cm create authrealm` 
- [Create a cmd to install RHACM #94](https://github.com/stolostron/cm-cli/issues/94)
- [Create a cmd to install MCE #180](https://github.com/stolostron/cm-cli/issues/180)
## Breaking changes

## Bug fixes

- [When clusterclaim not yet running, cm use/run generates an error](https://github.com/stolostron/cm-cli/issues/167)
- Fix running message in `cm run cc <clusterclaim_name>`
- [cm get cc <cc_name> should not wait for the cc to be running #177](https://github.com/stolostron/cm-cli/issues/177)
- [Scaling of clusterpool default size can cause problems with unexpectedly scaling down the pool #175](https://github.com/stolostron/cm-cli/issues/175)
- [Display which cc is currently active when running cm get cc #178](https://github.com/stolostron/cm-cli/issues/178)

