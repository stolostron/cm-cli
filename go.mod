module github.com/open-cluster-management/cm-cli

go 1.16

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
)

require (
	github.com/Masterminds/semver v1.5.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/open-cluster-management/governance-policy-propagator v0.0.0-20210810132451-6fcf70131732
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142
	github.com/openshift/hive/apis v0.0.0-20210707015124-49b5837aa081
	github.com/openshift/library-go v0.0.0-20210603104821-259346e2fd4c // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/cli-runtime v0.21.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/code-generator v0.21.1
	k8s.io/component-base v0.21.1
	k8s.io/klog/v2 v2.9.0
	k8s.io/kubectl v0.21.1
	open-cluster-management.io/api v0.0.0-20210607023841-cd164385e2bb
	open-cluster-management.io/clusteradm v0.1.0-alpha.4.0.20210709205037-2347693f34cd
)
