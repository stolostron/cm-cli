[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Cluster

The CLI has commands to manage clusters.

```bash
cm <verb> cluster <options...>
```
## Help

```bash
cm <verb> cluster -h
```

## Verbs

### Get clusters

Get the list of attached clusters

```bash
cm get clusters
```
### Attach Cluster

```bash
cm attach cluster --values <values_yaml_path>
```
or
```bash
cm attach cluster --values <values_yaml_path> --import-file <import_file>
```
or
```bash
cm attach cluster --cluster <cluster_name> [--cluster-kubeconfig <managed_cluster_kubeconfig_path>]
```

The `attach` verb provides the capability to attach a cluster to a hub.
The `attach` can be done on different ways. 
1. Manually:
    By adding the parameter `--import-file` to the `attach` command, an import yaml files will be generated then run the command `kubectl apply -f <import_file>` on the cluster you would like to import. This will install the agent on the cluster and it will connect to the hub.

2. Automatically:
    a) By providing the kubeconfig in the [values.yaml](../pkg/cmd/attach/cluster/scenario/attach/values-template.yaml), then a secret will be created on the hub cluster and the system will use it to install the agent. The secret is deleled if the `attach` failed or succeed and so the credentials are not kept on the hub.
    b) By providing the pair server/token in the [values.yaml](../pkg/cmd/attach/cluster/scenario/attach/values-template.yaml) and again a secret will be created on the hub and the system will use it to install the agent. The secret is deleled if the `attach` failed or succeed and so the credentials are not kept on the hub. 
    c) When the cluster was provisionned with hive. If the cluster was provisionned with hive, a clusterdeployemnt custom resource exists which contain a secret to access the remote cluster and thus if you `attach` a hive cluster, you don't have to provide any credential to access the cluster. The system will find out the credentials and attach the cluster.

    The `attach` command also takes `--cluster` and `--cluster-kubeconfig` instead of the `--values`, in that case the default [values.yaml](../pkg/attach/cluster/scenario/attach/values-default.yaml) will be used.

5. Attaching the hub: by default the hub is attached to itself but if you detached it and want to reattach it you just have to provide a [values.yaml](../pkg/cmd/attach/cluster/scenario/attach/values-template.yaml) with a cluster name `local-cluster`. The system will recognized that name and use the cluster credentials to do the attach.

### Detach Cluster

```bash
cm detach cluster --values <values_yaml_path>
```
or
```bash
cm detach cluster --cluster <cluster_name>
```

The `detach` verb will detach from the hub an already managed cluster.

### Create Cluster

```bash
cm create cluster --values <values_yaml_path>
```

The `create` will create a new managed cluster and attach it to the hub. Cloud provider credentials must be given in the values.yaml.

### Delete Cluster


```bash
cm delete cluster --values <values_yaml_path>
```
or
```bash
cm delete cluster --cluster <cluster_name>
```

The `delete` detaches and deletes an existing managed cluster.

### Scale a cluster:

This is valid only for cluster deployed with hive.

First find the machinepool you would like to alterate
```bash
cm get machinepools --cluster <cluster_name>
```
then change the number of replicas for the machinepool.
```bash
cm scale cluster --cluster <cluster_name> --machine-pool <machine_pool_name> --replicas <nb_replicas>
```