# Copyright Contributors to the Open Cluster Management project

BEFORE_SCRIPT := $(shell build/before-make.sh)

SCRIPTS_PATH ?= build

# Install software dependencies
INSTALL_DEPENDENCIES ?= ${SCRIPTS_PATH}/install-dependencies.sh

GOPATH := ${shell go env GOPATH}
GOBIN ?= ${GOPATH}/bin

CRD_OPTIONS ?= "crd:crdVersions=v1"

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

.PHONY: 
build-bin: doc-help
	tar -czf docs/help.tar.gz -C docs/help/ .
	zip -q docs/help.zip -j docs/help/*
	@rm -rf bin
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_darwin_amd64.tar.gz -C bin/ cm
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_amd64.tar.gz -C bin/ cm 
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_arm64.tar.gz -C bin/ cm 
	GOOS=linux GOARCH=ppc64le go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_ppc64le.tar.gz -C bin/ cm 
	GOOS=linux GOARCH=s390x go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm ./cmd/cm.go && tar -czf bin/cm_linux_s390x.tar.gz -C bin/ cm
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -gcflags=-trimpath=x/y  -o bin/cm.exe ./cmd/cm.go && zip -q bin/cm_windows_amd64.zip -j bin/cm.exe

,PHONY: doc-help
doc-help:
	@echo "Generate help markdown in docs/help"
	@go build -o docs/tools/cm docs/tools/cm.go && cd docs/tools && cm && rm cm
	@echo "Markdown generated"

.PHONY: install
install: build

.PHONY: plugin
plugin: build
	cp ${GOPATH}/bin/cm ${GOPATH}/bin/oc-cm
	cp ${GOPATH}/bin/cm ${GOPATH}/bin/kubectl-cm

.PHONY: check
## Runs a set of required checks
check: check-copyright doc-help

.PHONY: check-copyright
check-copyright:
	@build/check-copyright.sh

.PHONY: test
test: manifests
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
