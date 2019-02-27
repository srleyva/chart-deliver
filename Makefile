BINARY = chart
VERSION=0.0.1
GITHUB_USERNAME=srleyva
GITHUB_REPO=${GITHUB_USERNAME}/${BINARY}
DOCKER_REPO=sleyva97/$(BINARY)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION}"
pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)
gobuild = go build ${LDFLAGS} -o ${BINARY} main.go

all: format build-linux-server

format:
	go fmt $(pkgs)

test:
	go test `go list ./pkg/...`
	go test `go list ./cmd/...`
	rm -rf pkg/helpers/tester/

build:
	$(gobuild)

build-linux-server:
	CGO_ENABLED=0 GOOS=linux $(gobuild)

docker-build:
	docker build . -t ${DOCKER_REPO}:${VERSION}

docker-push: docker-build
	docker push ${DOCKER_REPO}:${VERSION}
