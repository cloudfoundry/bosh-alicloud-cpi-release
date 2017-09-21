BINDIR := $(CURDIR)/bin
MAINDIR := bosh-alicloud-cpi
MAINFILE := $(CURDIR)/src/$(MAINDIR)/main/main.go
EXECUTABLE := $(BINDIR)/alicloud_cpi

GOPATH := $(CURDIR)

GO_OPTION ?=
ifeq ($(VERBOSE), 1)
GO_OPTIONS += -v
endif

BUILD_OPTIONS = -a

all: clean build

clean:
	rm -f $(BINDIR)/*

build:
	mkdir -p $(BINDIR)
	go build $(GO_OPTIONS) $(BUILD_OPTIONS) -o ${EXECUTABLE} $(MAINFILE)

test:
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/action)
	go test -v $(shell find $(CURDIR) -name *_test.go | grep $(MAINDIR)/alicloud)
