$(shell PATH=$PATH:$GOPATH/bin)
BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
IMAGE_NAME := registry.gitlab.com/isaiahwong/cluster/api/accounts
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

set-image:
	kubectl set image deployments/accounts-deployment accounts=$(IMAGE_NAME)

set-image-latest:
	kubectl set image deployments/accounts-deployment accounts=$(IMAGE_NAME):latest

set-image-sha:
	kubectl set image deployments/accounts-deployment accounts=$(IMAGE_NAME):${BUILD_ID}

gen-manifest-release:
	./tools/gen-manifest.sh gen-cert --release true

genproto:
	if [ ! -d "api" ]; then \
			mkdir api; \
	fi

	protoc -I./proto/accounts-proto/api -I./proto/third_party/googleapis --go_out=plugins=grpc:./api ./proto/accounts-proto/api/accounts/v1/*.proto

genmocks:
	mockery -name=DataStore -dir=./internal/store -recursive=true -output=./tests/mocks       
	mockery -name=Repo -dir=./internal/store/repo/accounts -recursive=true -output=./tests/mocks        

compose-token:
	docker-compose -f docker-compose.yml exec hydra \
    hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id auth-code-client-2 \
    --secret secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5555/callback

compose-client:
	docker-compose -f docker-compose.yml exec hydra \
	hydra token user \
	--client-id auth-code-client-2 \
	--client-secret secret \
	--endpoint http://127.0.0.1:4444/ \
	--port 5555 \
	--scope openid,offline