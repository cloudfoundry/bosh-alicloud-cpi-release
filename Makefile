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

all: clean deps build

clean:
	rm -f $(BINDIR)/*

deps:
	go get -v github.com/cppforlife/bosh-cpi-go/...
	go get -v github.com/denverdino/aliyungo/...
	go get -v github.com/cloudfoundry/bosh-utils/logger
	go get -v github.com/cloudfoundry/bosh-utils/uuid
	go get -v github.com/cloudfoundry/bosh-utils/system

build:
	mkdir -p $(BINDIR)
	go build $(GO_OPTIONS) $(BUILD_OPTIONS) -o ${EXECUTABLE} $(MAINFILE)

test:
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/action)
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/alicloud)
