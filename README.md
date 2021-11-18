[comment]: # ( Copyright Contributors to the Open Cluster Management project )
# Open Cluster Management CLI & CLI Plugin

A CLI and kubernetes CLI plugin that allows you to interact with OCM/ACM/MCE to provision and managed your Hybrid Cloud presence from the command-line.

## Requirements

Go 1.16 is required in order to build or contribute on this project as it leverage the `go:embed` tip.

## Installation

### Using releases

You can download the binary from [https://github.com/open-cluster-management/cm-cli/releases](https://github.com/open-cluster-management/cm-cli/releases)

### Using Krew

1. Install krew [https://krew.sigs.k8s.io/docs/user-guide/setup/install/](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)
2. Plugins can be installed with the `kubectl krew install cm`

### CLI
The binary will be installed in GOPATH/bin

```bash
git clone https://github.com/open-cluster-management/cm-cli.git
cd cm-cli
make build
cm
```

### Plugin
The binary will be installed in GOPATH/bin

This will create a binary `oc-cm` and `kubectl-cm` in the `$GOPATH/go/bin` allowing you to call `oc cm` or `kubectl cm`
```bash
git clone https://github.com/open-cluster-management/cm-cli.git
cd cm-cli
make plugin
kubectl cm
oc cm
```

## Disclaimer

This CLI (and plugin) is still in development, but aims to expose OCM/ACM's functional through a useful and lightweight CLI and kubectl/oc CLI plugin.  Some features may not be present, fully implemented, and it might be buggy!  

## Getting Started

### Setting up a ClusterPoolHost

In order to work with clusters, you need set up `cm` with your hub cluster(s) - `cm` refers to these hubs as "clusterpoolhost"(s) or "cph"(s) for short!  

To set up your first ClusterPoolHost:
1. `oc login` to your ClusterPoolHost running [Red Hat Advanced Cluster Management](https://access.redhat.com/products/red-hat-advanced-cluster-management-for-kubernetes), [Multicluster Engine for Kubernetes](https://github.com/open-cluster-management/backplane-operator), or [Open Cluster Management](http://github.com/open-cluster-management-io).  **Your user must be able to create ServiceAccounts in the target namespace, given that `create cph` creates a ServiceAccount.  Also ensure that ServiceAccounts in that namespace have the relevant access such as create/delete ClusterClaims, ClusterPools, etc.   
2. Run `cm create cph --api-server=<api-url> --console=<console-url> --group=<rbac-group> --namespace=<namespace-containing-clusters> <name-of-cph>` and run `cm create cph --help` to view all options
3. Run `cm get cph` to verify that your active clusterpoolhost is correct, `cm set cph <clusterpoolhost-name>` to swap clusterpoolhosts, and `cm use cph <clusterpoolhost-name>` to switch to that clusterpoolhost's context.  

### Working with Clusters

#### Using ClusterPools

ClusterPools maintain a configurable and scalable number of OpenShift clusters in a hibernating state.  `cm` exposes the capability to create, view, consume clusters from, and destroy ClusterPools.  

To create a ClusterPool, see `cm create <clusterpool/cp> --help` for all options.  You can view your created ClusterPools with `cm get <clusterpool/cp>`.

Once you have a ClusterPool, you can claim clusters from the pool for use using `cm create <clusterclaim/cc> <clusterpool-name> <clusterclaim-name>`, see `cm create <clusterclaim/cc> --help` for more options.  You can also view ClusterClaims and details with `cm get <clusterclaim/cc>` and delete ClusterClaims with `cm delete <clusterclaim/cc>`.  

Finally, you can delete a ClusterPool with `cm delete <clusterpool/cp>`.  

#### Creating Clusters

To create individual clusters with specific configurations, you can use `cm create cluster`, see `cm create cluster --help` for more options, including a `values.yaml` template.  

#### Nagivating with Clusters

`cm` also allows you to easily change cluster contexts without losing visibility to your Hub/ClusterPoolHost.  

To view the available clusters, use `cm get <clusterclaim/clusterpoolhost>`.  

Once you've identified a cluster, you can use `cm use <clusterclaim/clusterpoolhost> <cluster-name>` to switch to that cluster's context.  You can use these commands again to list and change contexts without losing your Hub/ClusterPoolHost context.  

`cm console <clusterclaim/clusterpoolhost>` allows you to quickly open the console of a claimed cluster or ClusterPoolHost.  

#### Hibernating/Waking Clusters

`cm` also allows you to hibernate clusters via `cm hibernate <clusterclaim/cc> <cluster-name>` and wake clusters using `cm run <clusterclaim/cc> <cluster-name>`.  

#### Bringing Clusters Under Management

If you wish to bring a cluster under the management of a hub cluster, you can use `cm attach <cluster/clusterclaim>`, see `cm attach --help` for all options.  

## Contributing

See our [Contributing Document](CONTRIBUTING.md) for more information.  

## Commands

[command help](docs/help/cm.md)

[general commands](docs/general.md)

[cluster commands](docs/cluster.md)

[clusterpool commands](docs/clusterpool.md)

[policies commands](docs/policies.md)
