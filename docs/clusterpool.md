[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# ClusterPools

The CLI has commands to manage clusterpools. It allows to manage multiple clusters that are holding the clusterpool management system. 

These clusters are called here `clusterpoolhost`

## ClusterPoolHosts
### Create ClusterPoolHosts

A clusterpoolhost defines where clusterclaim can be created using the cm-cli.

```bash
cm create clusterpoolhost <clusterpoolhost_name> --api-server <api_server_url> --console <console_url> --namespace <my_namespace> [--group <my_user_group>]
```

The namespace is where the clusterclaim will be created on that clusterpoolhost.
The group is user group that will be bind to roles in order to retreive the cluster credentials

When you create a clusterpoolhosts it becomes the active one and all other commands will be done toward that clusterpoolhosts.

Example:
```
cm create clusterpoolhost my_cluster_pool_host_name --api-server https://api.mycluster.my.domain:6443 --console https://console-openshift-console.apps.mycluster.my.domain --namespace my_namespace --group my_user_group
```

### Get ClusterPoolHosts

```bash
cm get clusterpoolhosts <options>
```

The list of clusterpoolhosts is maintained in the `~/.kube/known-cphs`.

### Delete ClusterPoolHosts

```bash
cm delete clusterpoolhost <clusterpoolhost_name>
```
### Set a clusterpoolhost as active

```bash
cm set clusterpoolhost <clusterpoolhost_name>
```
Setting as active means the CLI will used it to find a cluster claims attached to that clusterpoolhost.

Once the clusterpoolhost is created, the `~/.kube/config` is updated with a context pointing to that cluster and the clusterpoolhost is set as the current one.

### Use a clusterpoolhost

If you want to run `oc` or `kubectl` toward a clusterpoolhost

```bash
cm use cph <clusterpoolhost_name>
```

It updates the current-context in `~/.kube/config` to point to that cluster.

## ClusterPools
### Create a clusterpool

```bash
cm create clusterpool [<clusterpool_name>] --values <values_yaml_path>
```

if the clusterpool_name is specified then it overwrites the one in the values yaml.

The template can be retreived by running 

```bash
cm create cp -h
```

it supports clusterpools for AWS, Azure and Google


### Get cluserpools

```bash
cm get clusterpool [--cph <clusterpoolhost>|-A]
```
### Scale a clusterpool

```bash
cm scale <clusterpool> --size <size> [--cph <clusterpoolhost>] 
```
### Delete a clusterpool

```bash
cm delete clusterpool <clusterpool_name> [--cph <clusterpoolhost_name>]
```

## ClusterClaims
### Creeate clusterclaims

To create clusterclaims on the active clusterpool, the command `cm create clusterclaim|cc <clusterpool_name> <clusterclaim_name>` can be executed. Multiple clusterclaims can be created simultaneously by providing a list (comma-separated) of clusterclaim name.

For example:

```bash
cm create clusterclaim myclusterpool_name clusterclaim1,clusterclaim2
```

NB: The comma-separated list must not contain space, if it does it should be surrounded by double-quotes.

### Use a cluster claim managed by a clusterpoolhost

```bash
cm use clusterclaim <clusterclaim_name> [--cph <clusterpoolhost_name>]
```

It updates the kubeconfig with a context toward that cluster. If the KUBECONFIG environment variable is set, the file specifed in the environment variable is updated with the context.

### Get the list of clusterclaims

```bash
cm get clusterclaim [--cph <clusterpoolhost_name>| -A]
```

### Get the credential for a clusterclaim
```bash
cm get clusterclaim <clusterclaim_name> [--cph <clusterpoolhost_name>]
```
### Hibernate clusterclaims

```bash
cm hibernate clusterclaim <clusterclaim>[,<clusterlcaim>...] [--skip-schedule]
```
The option `--skip-schedule` will opt-out the clusterclaim from the cronjob hibernation.

### Run clusterclaims

```bash
cm run clusterclaim <clusterclaim>[,<clusterlcaim>...] [--skip-schedule]
```
The option `--skip-schedule` will opt-out the clusterclaim from the cronjob hibernation.

### Attach clusterclaims

```bash
cm attach clusterclaim <clusterclaim_name>
```

### Detach a clusterclaim

```bash
cm detach cluster <clusterclaim_name>
```
### Delete clusterclaims

```bash
cm delete clusterclaim <clusterclaim>[,<clusterclaim>...] [--cph <clusterpoolhost_name>]
```
