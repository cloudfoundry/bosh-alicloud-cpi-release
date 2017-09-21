BINDIR := $(CURDIR)/bin
MAINDIR := bosh-alicloud-cpi
MAINFILE := $(CURDIR)/src/$(MAINDIR)/main/main.go
EXECUTABLE := $(BINDIR)/alicloud_cpi

GOPATH := $(CURDIR)

COMMIT = $(shell git rev-parse --short HEAD)

GO_OPTION ?=
ifeq ($(VERBOSE), 1)
GO_OPTIONS += -v
endif

# TODO add local link invocation
BUILD_OPTIONS = -a -ldflags "-X main.GitCommit=\"$(COMMIT)\""

all: clean deps build create-release deploy-scp

clean:
	rm -f $(BINDIR)/*

deps:
	go get -v github.com/cppforlife/bosh-cpi-go/...
	go get -v github.com/denverdino/aliyungo/...
	go get -v github.com/cloudfoundry/bosh-utils/logger
	go get -v github.com/cloudfoundry/bosh-utils/uuid
	go get -v github.com/cloudfoundry/bosh-utils/system
	go get -v github.com/onsi/ginkgo/ginkgo
	go get -v github.com/onsi/gomega

build: deps
	mkdir -p $(BINDIR)
	go build $(GO_OPTIONS) $(BUILD_OPTIONS) -o ${EXECUTABLE} $(MAINFILE)

test: build
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/action)
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/alicloud)

create-release: build
    git commit -c 'commit for test deploy release'
    bosh create-release --force --tarball=bin/bosh-alicloud-cpi.tgz

deploy-scp: create-release
    scp bin/bosh-alicloud-cpi.tgz root@${DEPLOY_SERVER_IP}:/root
