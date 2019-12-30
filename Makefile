
GO_FILES?=$(shell find . -name '*.go' |grep -v vendor)
TAG?=latest
GIT_COMMIT=$(shell git rev-list -1 HEAD)
VERSION=$(shell git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')

.PHONY: build

install-ci:
	./hack/install-ci.sh
ci:
	./hack/integration-test.sh

.PHONY: static
static:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /template/|grep -v /build/) -cover \
    && CGO_ENABLED=0 go build --ldflags "-s -w \
    -X github.com/Qolzam/ofc-bootstrap/version.GitCommit=${GIT_COMMIT} \
    -X github.com/Qolzam/ofc-bootstrap/version.Version=${VERSION}" \
    -a -installsuffix cgo -o ofc-bootstrap

.PHONY: build
build:
	./build.sh

.PHONY: dist
dist:
	./build_redist.sh
