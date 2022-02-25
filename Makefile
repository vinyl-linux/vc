PREFIX ?= ""
BINDIR ?= "$(PREFIX)/usr/bin"

BUILT_ON ?= $(shell date --rfc-3339=seconds | sed 's/ /T/')
BUILT_BY := $(shell whoami)
BUILD_REF := $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

.PHONY: default install
default: vc

vc: pkg = "github.com/vinyl-linux/vc/bin/cmd"
vc: bin/*.go bin/**/*.go go.mod go.sum
	(cd bin && CGO_ENABLED=0 go build -ldflags="-s -w -X $(pkg).Ref=$(BUILD_REF) -X $(pkg).BuildUser=$(BUILT_BY) -X $(pkg).BuiltOn=$(BUILT_ON)" -trimpath -o ../$@)

install: "$(BINDIR)/vc"
$(BINDIR)/vc: vc
	install -m 0755 -o root $< $@
