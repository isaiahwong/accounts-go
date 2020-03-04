$(shell PATH=$PATH:$GOPATH/bin)
BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
IMAGE_NAME := registry.gitlab.com/isaiahwong/api/auth
VERSION := 0.0.1

PROTO_DIR := ../../pb

.DEFAULT_GOAL := help

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:
	docker build -t $(IMAGE_NAME):latest . --rm=true

build-sha: 
	docker build -t $(IMAGE_NAME):$(BUILD_ID) . --rm=true

push: 
	docker push $(IMAGE_NAME):latest

push-sha:
	docker push $(IMAGE_NAME):$(BUILD_ID)

build-all:
	docker build -t $(IMAGE_NAME):latest -t $(IMAGE_NAME):$(BUILD_ID) . --rm=true

push-all:
	make push
	make push-sha

build-push:
	make build-all
	make push-all

clean:
	docker rmi $( docker images | grep '<none>') --force 2>/dev/null

gen-manifest-release:
	./tools/gen-manifest.sh gen-cert --release true

genproto:
	if [ ! -d "api" ]; then \
			mkdir api; \
	fi

	protoc -I./proto/api -I./proto/third_party/googleapis --go_out=plugins=grpc:./api ./proto/api/accounts/v1/*.proto

genmocks:
	mockery -name=DataStore -dir=./internal/store -recursive=true -output=./tests/mocks       
	mockery -name=Repo -dir=./internal/store/repo/user -recursive=true -output=./tests/mocks        
