GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=v0.0.1
SERVER_NAME=sidecar
AUTHOR=frambos
IMG=$(AUTHOR)/$(SERVER_NAME):$(VERSION)
BUILDER_NAME=$(SERVER_NAME)-multi-platform-buildx
PLATFORMS = linux/amd64 linux/arm64

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest


.PHONY: api
api:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/sidecar

.PHONY: docker-build
docker-build:
	docker build -t $(IMG) -f Dockerfile .

.PHONY: docker-run
docker-run:
	docker run -it -d --rm -p 8000:8000 -p 9000:9000 -v ./configs/:/data/conf --name $(SERVER_NAME)-$(VERSION) $(IMG) 

.PHONY: docker-stop
docker-stop:
	docker stop $(SERVER_NAME)-$(VERSION)

.PHONY: docker-push
docker-push:
	docker push $(IMG)

.PHONY: docker-dev-build
docker-dev-build:
	docker build -f Dockerfile.dev -t sidecar-dev .

.PHONY: docker-dev
docker-dev:
	docker run -it -d --rm --name sidecar-dev -p 8001:8001 -p 9001:9001 -v ./:/go/src/sidecar sidecar-dev

.PHONY: run
run:
	go run ./cmd/sidecar -conf ./configs/

.PHONY: generate
generate:
	go mod tidy
	@cd cmd/sidecar && wire && cd -

.PHONY: multi-platform-build
multi-platform-build:
	@for platform in $(PLATFORMS); do \
		image_name=$$platform/$(IMG); \
		docker build -f Dockerfile.multi_platform_build --platform=$$platform -t $$image_name --load . ; \
		platform_formated=$$(echo $$platform | tr '[:upper:]' '[:lower:]' | tr '/' '-'); \
		container_name=$$platform_formated-$(SERVER_NAME)-$(VERSION); \
		docker run -it -d --rm --name $$container_name $$image_name; \
		sleep 5; \
		docker cp $$container_name:/app.tar.gz ./built/$$container_name.tar.gz; \
		cd ./built && sha256sum $$container_name.tar.gz > $$container_name.tar.gz.sha256sum && cd - ; \
		sleep 5; \
	done

.PHONY: all
all:
	make api;
	make generate;

help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
