commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)
repo ?= github.com/awslabs/k8s-eniconfig-controller

.PHONY: build
build:
	go build -o eniconfig-controller -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/eniconfig-controller

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: dev-release
dev-release:
	goreleaser --rm-dist --snapshot --skip-publish

.PHONY: test
test:
	go test -v -cover -race $(repo)/...

.PHONY: tag
tag:
	git tag -a ${VERSION} -s
	git push origin --tags

