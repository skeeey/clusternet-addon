all: build
.PHONY: all

IMAGE_REGISTRY ?= quay.io/skeeey
IMAGE_TAG ?= 0.2
IMAGE_NAME ?= $(IMAGE_REGISTRY)/clusternet-addon:$(IMAGE_TAG)

# build
vendor: 
	go mod tidy
	go mod vendor
.PHONY: vendor

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o clusternet cmd/clusternet/main.go
.PHONY: build

image:
	docker build -f Dockerfile -t $(IMAGE_NAME) .
.PHONY: image
