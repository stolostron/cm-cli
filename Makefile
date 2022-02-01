# Copyright Contributors to the Open Cluster Management project

BEFORE_SCRIPT := $(shell build/before-make.sh)

SCRIPTS_PATH ?= build

# Install software dependencies
INSTALL_DEPENDENCIES ?= ${SCRIPTS_PATH}/install-dependencies.sh

GOPATH := ${shell go env GOPATH}
GOBIN ?= ${GOPATH}/bin
GOOS := ${shell go env GOOS}
GOARCH := ${shell go env GOARCH}

CRD_OPTIONS ?= "crd:crdVersions=v1"

export KREW_DIR=$(shell mktemp -d)

export PROJECT_DIR            = $(shell 'pwd')
export PROJECT_NAME			  = $(shell basename ${PROJECT_DIR})

export GOPACKAGES   = $(shell go list ./... | grep -v /vendor | grep -v /build | grep -v /test | grep -v /scenario )

.PHONY: clean
clean: clean-test
	kind delete cluster --name ${PROJECT_NAME}-functional-test
	
.PHONY: deps
deps:
	@$(INSTALL_DEPENDENCIES)

.PHONY: build
build: 
	rm -f ${GOPATH}/bin/cm
	go install ./cmd/cm.go

.PHONY: build-bin
build-bin: doc-help
	tar -czf docs/help.tar.gz -C docs/help/ .
	zip -q docs/help.zip -j docs/help/*
	@rm -rf bin
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_darwin_amd64.tar.gz LICENSE -C bin/ cm
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_amd64.tar.gz LICENSE -C bin/ cm 
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_arm64.tar.gz LICENSE -C bin/ cm 
	GOOS=linux GOARCH=ppc64le go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_ppc64le.tar.gz LICENSE -C bin/ cm 
	GOOS=linux GOARCH=s390x go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_s390x.tar.gz LICENSE -C bin/ cm
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm.exe ./cmd/cm.go && zip -q bin/cm_windows_amd64.zip LICENSE -j bin/cm.exe

.PHONY: release
release: 
	@if [[ -z "${VERSION}" ]]; then VERSION=`cat VERSION.txt`; echo $$VERSION; fi; \
	git tag v$$VERSION && git push upstream --tags

.PHONY: build-krew
build-krew: krew-tools
	@if [[ -z "${VERSION}" ]]; then VERSION=`cat VERSION.txt`; echo $$VERSION; fi; \
	docker run -v ${PROJECT_DIR}/.krew.yaml:/tmp/template-file.yaml rajatjindal/krew-release-bot:v0.0.40 \
	krew-release-bot template --tag v$$VERSION --template-file /tmp/template-file.yaml > cm.yaml; 
	KREW=/tmp/krew-${GOOS}\_$(GOARCH) && \
	KREW_ROOT=`mktemp -d` KREW_OS=darwin KREW_ARCH=amd64 $$KREW install --manifest=cm.yaml && \
	KREW_ROOT=`mktemp -d` KREW_OS=linux KREW_ARCH=amd64 $$KREW install --manifest=cm.yaml && \
	KREW_ROOT=`mktemp -d` KREW_OS=linux KREW_ARCH=arm64 $$KREW install --manifest=cm.yaml && \
	KREW_ROOT=`mktemp -d` KREW_OS=windows KREW_ARCH=amd64 $$KREW install --manifest=cm.yaml;

.PHONY: krew-tools
krew-tools:
ifeq (, $(shell which /tmp/krew-$(GOOS)\_$(GOARCH)))
	@( \
		set -x; cd /tmp && \
		KREW=krew-$(GOOS)\_$(GOARCH); \
		curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/$$KREW.tar.gz" && \
		tar zxvf $$KREW.tar.gz \
	) 
endif

.PHONY: doc-help
doc-help:
	@echo "Generate help markdown in docs/help"
	go build -o docs/tools/cm docs/tools/cm.go && PATH=docs/tools cm && rm docs/tools/cm
	@echo "Markdown generated"
	@build/clean-docs.sh

.PHONY: install
install: build

.PHONY: plugin
plugin: build
	cp ${GOPATH}/bin/cm ${GOPATH}/bin/oc-cm
	cp ${GOPATH}/bin/cm ${GOPATH}/bin/kubectl-cm

.PHONY: check
## Runs a set of required checks
check: check-copyright

.PHONY: check-copyright
check-copyright:
	@build/check-copyright.sh

.PHONY: test
test: controller-gen manifests
	@build/run-unit-tests.sh

.PHONY: clean-test
clean-test: 
	-rm -r ./test/unit/coverage
	-rm -r ./test/unit/tmp
	-rm -r ./test/functional/tmp
	-rm -r ./test/out

.PHONY: functional-test-full
functional-test-full: deps install
	@build/run-functional-tests.sh

.PHONY: functional-test-full-clean
functional-test-full-clean:
	@build/run-functional-tests-clean.sh

.PHONY: manifests
manifests:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=api/cm-cli/v1alpha1/crd

.PHONY: generate
generate: controller-gen manifests
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	@hack/update-codegen.sh

# find or download controller-gen
# download controller-gen if necessary
.PHONY: controller-gen
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: kubebuilder-tools
## Find or download kubebuilder
kubebuilder-tools:
ifeq (, $(shell which kubebuilder))
	@( \
		set -ex ;\
		KUBEBUILDER_TMP_DIR=$$(mktemp -d) ;\
		cd $$KUBEBUILDER_TMP_DIR ;\
		curl -L -o $$KUBEBUILDER_TMP_DIR/kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/3.1.0/$$(go env GOOS)/$$(go env GOARCH) ;\
		chmod +x $$KUBEBUILDER_TMP_DIR/kubebuilder && mv $$KUBEBUILDER_TMP_DIR/kubebuilder /usr/local/bin/ ;\
	)
endif

# See https://book.kubebuilder.io/reference/envtest.html.
#    kubebuilder 2.3.x contained kubebuilder and etc in a tgz
#    kubebuilder 3.x only had the kubebuilder, not etcd, so we had to download a different way
# After running this make target, you will need to either:
# - export KUBEBUILDER_ASSETS=$HOME/kubebuilder/bin
# OR
# - sudo mv $HOME/kubebuilder /usr/local
#
# This will allow you to run `make test`
.PHONY: envtest-tools
## Install envtest tools to allow you to run `make test`
envtest-tools:
ifeq (, $(shell which etcd))
		@{ \
			set -ex ;\
			ENVTEST_TMP_DIR=$$(mktemp -d) ;\
			cd $$ENVTEST_TMP_DIR ;\
			K8S_VERSION=1.19.2 ;\
			curl -sSLo envtest-bins.tar.gz https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-$$K8S_VERSION-$$(go env GOOS)-$$(go env GOARCH).tar.gz ;\
			tar xf envtest-bins.tar.gz ;\
			mv $$ENVTEST_TMP_DIR/kubebuilder $$HOME ;\
			rm -rf $$ENVTEST_TMP_DIR ;\
		}
else
   
endif
