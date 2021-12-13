# Copyright 2021 Linka Cloud  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODULE = go.linka.cloud/protoautoindex


PROTO_BASE_PATH = $(PWD)

THIRD_PARTY = .deps

TMP_ROOT := .tmp
TMP_GOOGLE_API := $(TMP_ROOT)/google/api

INCLUDE_PROTO_PATH = -I$(PROTO_BASE_PATH) -I $(THIRD_PARTY) \
	-I $(shell go list -m -f {{.Dir}} go.linka.cloud/protoc-gen-defaults) \
	-I $(shell go list -m -f {{.Dir}} go.linka.cloud/protofilters) \
	-I $(shell go list -m -f {{.Dir}} github.com/envoyproxy/protoc-gen-validate) \
	-I $(shell go list -m -f {{.Dir}} github.com/alta/protopatch) \
	-I $(shell go list -m -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway/v2)

PROTO_OPTS = paths=source_relative

PROTO_FILES = $(shell find $(PROTO_BASE_PATH)/api -name '*.proto' -type f)

$(shell mkdir -p .bin .deps)

export GOBIN=$(PWD)/.bin

export PATH := $(GOBIN):$(PATH)

bin:
	@go install github.com/golang/protobuf/protoc-gen-go
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install go.linka.cloud/protoc-gen-defaults
	@go install go.linka.cloud/protoc-gen-go-fields
	@go install github.com/envoyproxy/protoc-gen-validate
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
	@go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto

clean:
	@rm -rf .bin $(TMP_ROOT) $(THIRD_PARTY)
	@find $(PROTO_BASE_PATH) -name '*.pb*.go' -type f -exec rm {} \;

proto-deps:
	@test -f $(THIRD_PARTY)/google/api/annotations.proto || rm -rf $(TMP_ROOT) && \
		mkdir -p $(TMP_ROOT) && \
		git clone https://fuchsia.googlesource.com/third_party/googleapis $(TMP_ROOT) --depth 1 &> /dev/null && \
		mkdir -p $(THIRD_PARTY)/google/api && \
		mv $(TMP_GOOGLE_API)/annotations.proto \
			$(TMP_GOOGLE_API)/http.proto \
			$(TMP_GOOGLE_API)/httpbody.proto \
			$(THIRD_PARTY)/google/api && \
			rm -rf $(TMP_ROOT)


.PHONY: proto
proto: proto-deps gen-proto lint


.PHONY: gen-proto
gen-proto: bin
	@protoc $(INCLUDE_PROTO_PATH) \
		--go-patch_out=plugin=go,$(PROTO_OPTS):. \
		--go-patch_out=plugin=go-grpc,$(PROTO_OPTS):. \
		--go-patch_out=plugin=defaults,$(PROTO_OPTS):. \
		--go-patch_out=plugin=go-fields,$(PROTO_OPTS):. \
		--go-patch_out=plugin=go-vtproto,features=marshal+unmarshal+size,$(PROTO_OPTS):. \
		--go-patch_out=plugin=grpc-gateway,$(PROTO_OPTS):. \
		--go-patch_out=plugin=openapiv2:. \
		--go-patch_out=plugin=validate,lang=go,$(PROTO_OPTS):. $(PROTO_FILES)

.PHONY: lint
lint:
	@goimports -w -local $(MODULE) $(PWD)
	@gofmt -w $(PWD)

.PHONY: tests
tests: proto
	@go test -v ./...
