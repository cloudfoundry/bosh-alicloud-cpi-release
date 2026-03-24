BINDIR := $(CURDIR)/bin
MAINDIR := bosh-alicloud-cpi
MAINFILE := ./main/main.go
EXECUTABLE := $(BINDIR)/alicloud_cpi

GO_OPTION ?=
ifeq ($(VERBOSE), 1)
GO_OPTIONS += -v
endif

BUILD_OPTIONS = -a -mod=vendor

all: clean build

clean:
	rm -f $(BINDIR)/*

build:
	mkdir -p $(BINDIR)
	cd $(MAINDIR) && go build $(GO_OPTIONS) $(BUILD_OPTIONS) -o $(EXECUTABLE) $(MAINFILE)

test:
	cd $(MAINDIR) && go test -mod=vendor ./...

update-deps:
	cd $(MAINDIR) && go get -u ./... && go mod tidy

update-deps-patch:
	cd $(MAINDIR) && go get -u=patch ./... && go mod tidy