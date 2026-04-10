GO ?= go
OUTPUT ?= build/dae-next
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

date := $(shell git log -1 --format="%cd" --date=short | sed s/-//g)
count := $(shell git rev-list --count HEAD)
commit := $(shell git rev-parse --short HEAD)
VERSION ?= dae-next-$(date).r$(count).$(commit)

.PHONY: dae-next

dae-next:
	mkdir -p $(dir $(OUTPUT))
	GOWORK=off GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -ldflags "-s -w -X main.Version=$(VERSION)" -o $(OUTPUT) ./cmd/dae-next
