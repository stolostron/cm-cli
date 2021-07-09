[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# ClusterPools

The CLI has commands to manage clusterpools. It allows to manage multiple clusters that are holding the clusterpool management system. 

These clusters are called here `clusterpoolhost`

## Manage ClusterPoolHosts

A clusterpoolhost can be created using the command `cm create clusterpoolhost|cph <clusterpoolhost_name> <options>`. 

The list of clusterpoolhosts can be retreived by calling the command `cm get clusterpoolhosts|cphs <options>`.

The list of clusterpoolhosts is maintained in the `~/.kube/known-cphs`.

Delete a clusterpoolhosts can be acheived by runnig `cm delete clusterpoolhost|cph <clusterpoolhost_name>`.

Set a clusterpoolhost active or current with `cm set-cph <clusterpoolhost_name>`. Setting as active means the CLI will used it to find a cluster claims attached to that clusterpoolhost.

Once the clusterpoolhost is created, the `~/.kube/config` is updated with a context pointing to that cluster and the clusterpoolhost is set as the current one.

## Use a clusterpoolhost

`cm use-cph <clusterpoolhost_name>` updates the `~/.kube/config` to point to that cluster.

## Use a cluster claim managed by a clusterpoolhost

First, the clusterpoolhost which manage the cluster claim must be the current one. This can be done by using the command `cm use <clusterpoolhost_name>`

Then use `cm use <cluster_claim_name>` to update the kubeconfig with a context toward that cluster. If the KUBECONFIG environment variable is set, the file specifed in the environment variable is updated with the context.

