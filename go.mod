module github.com/open-cluster-management/cm-cli

go 1.16

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
	open-cluster-management.io/clusteradm => /Users/dvernier/acm-tools/clusteradm
)

require (
	github.com/Masterminds/semver v1.5.0
	github.com/ghodss/yaml v1.0.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/openshift/library-go v0.0.0-20210603104821-259346e2fd4c
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/cli-runtime v0.21.1
	k8s.io/client-go v1.5.2
	k8s.io/component-base v0.21.1
	k8s.io/kubectl v0.21.1
	open-cluster-management.io/api v0.0.0-20210607023841-cd164385e2bb
	open-cluster-management.io/clusteradm v0.0.0-20210611012639-e513b4f17927
)
