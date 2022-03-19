GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name  *.proto)
API_PROTO_FILES=$(shell cd ../../ && find api/$(APP_NAME) -name [^error]*.proto)
ERROR_PROTO_FILES=$(shell cd ../../ && find api/$(APP_NAME) -name error*.proto)
APP_NAME=demo

.PHONY: init
# init env
init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/Paradox315/protoc-gen-go-xhttp
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2

	go get -u github.com/gogo/protobuf/proto
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/gogo/protobuf/protoc-gen-gogoslick

	go get -u github.com/google/wire/cmd/wire
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/envoyproxy/protoc-gen-validate@latest



.PHONY: errors
# generate errors code
errors:
	cd ../../ && protoc --proto_path=. \
           --proto_path=./third_party \
           --gogofaster_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(ERROR_PROTO_FILES)

.PHONY: config
# generate config proto code
config:
	protoc --proto_path=. \
	       --proto_path=../../third_party \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

.PHONY: validate
# generate validate code
validate:
	cd ../../ && protoc --proto_path=. \
                 		--proto_path=./third_party \
                 		--gofast_out=paths=source_relative:. \
                 		--validate_out=paths=source_relative,lang=go:. \
                 $(API_PROTO_FILES)


.PHONY: swagger
# generate swagger
swagger:
	cd ../../ && protoc --proto_path=. \
	        --proto_path=./third_party \
	        --openapiv2_out . \
	        --openapiv2_opt logtostderr=true \
           $(API_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	 cd ../../ && protoc  --proto_path=. \
					      --proto_path=./third_party \
					  	  --gofast_out=paths=source_relative:. \
					  	  --go-xhttp_out=paths=source_relative:. \
					 	  --go-grpc_out=paths=source_relative:. \
	              $(API_PROTO_FILES)

.PHONY: service
# generate service proto
service:
	 cd ../../ && $(foreach proto_file,$(API_PROTO_FILES), \
    				echo $(proto_file); \
            		kratos proto server $(proto_file) -t app/$(APP_NAME)/internal/service; \
    				)\
     echo 'generate API proto file done!'

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...

.PHONY: all
# generate all
all:
	make api;
	make errors;
	make config;
	make validate;
	make swagger;
	make service;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
