GIT_COMMIT=$(shell git rev-parse --short HEAD)

GOTEST=go test
GOCOVER=go tool cover

ARCHES=amd64 arm64
PLATFORMS=darwin linux windows

BUILDARCH?=$(shell uname -m)
BUILDPLATFORM?=$(shell uname -s)
ifeq ($(BUILDARCH),aarch64)
  BUILDARCH=arm64
endif
ifeq ($(BUILDARCH),x86_64)
  BUILDARCH=amd64
endif
ifeq ($(BUILDPLATFORM),Darwin)
  BUILDPLATFORM=darwin
endif
ifeq ($(BUILDPLATFORM),Linux)
  BUILDPLATFORM=linux
endif
ifeq ($(BUILDPLATFORM),Win)
  BUILDPLATFORM=windows
endif

# unless otherwise set, I am building for my own architecture, i.e. not cross-compiling
ARCH ?= $(BUILDARCH)
PLATFORM ?= $(BUILDPLATFORM)

# canonicalized names for target architecture
ifeq ($(ARCH),aarch64)
  override ARCH=arm64
endif
ifeq ($(ARCH),x86_64)
  override ARCH=amd64
endif
ifeq ($(PLATFORM),Darwin)
  override PLATFORM=darwin
endif
ifeq ($(PLATFORM),Linux)
  override PLATFORM=linux
endif
ifeq ($(PLATFORM),Win)
  override PLATFORM=windows
endif

VERSION ?= $(GIT_COMMIT)
DEFAULTIMAGE ?= dibi/kube-resource-explorer:$(VERSION)

.PHONY: all

all: clean test cover build install

docker:
	docker build . -t sadrian99/microservice && docker push sadrian99/microservice
my:
	go build -o ./out/memreq ./cmd/memreq
my-run:
	go run ./cmd/memreq
test:
	$(shell mkdir TestResults)
	$(GOTEST) -v -coverprofile=TestResults/coverage.out ./...

cover: test
	$(GOCOVER) -func=TestResults/coverage.out
	$(GOCOVER) -html=TestResults/coverage.out

build:
	CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) GO111MODULE=on\
		go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -a -installsuffix cgo -o ./out/kube-resource-explorer-$(PLATFORM)-$(ARCH) ./cmd/kube-resource-explorer

install:
	CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) GO111MODULE=on\
		go install -ldflags "-X main.GitCommit=${GIT_COMMIT}" -a -installsuffix cgo ./cmd/kube-resource-explorer

run:
	CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) GO111MODULE=on\
		go run -ldflags "-X main.GitCommit=${GIT_COMMIT}" -a -installsuffix cgo ./cmd/kube-resource-explorer

package:
	DOCKER_BUILDKIT=1 docker build -t $(DEFAULTIMAGE) .

clean:
	rm -rf ./TestResults/* ./out/*
	docker rmi $(DEFAULTIMAGE) || true
