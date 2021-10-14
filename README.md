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

## Contributing

See our [Contributing Document](CONTRIBUTING.md) for more information.  

## Commands

[command help](docs/help/cm.md)

[general commands](docs/general.md)

[cluster commands](docs/cluster.md)

[clusterpool commands](docs/clusterpool.md)

[policies commands](docs/policies.md)
