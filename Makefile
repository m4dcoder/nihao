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

DEBUG ?= false
BUILDPATH=$(CURDIR)
VERSION ?= unset

GO=$(shell which go)
GOINSTALL=$(GO) install
GOBUILD=$(GO) build
GOLIST=$(GO) list
GOTOOL=$(GO) tool
GOVET=$(GO) vet
GOFMT=$(shell which gofmt)
GOLINT=$(shell which golint)
GOSEC=$(shell which gosec)
GOEXCLUDES=-e generated -e swagger -e vendor -e pkg/dep/sources

GOTEST=$(GO) test
ifeq ($(DEBUG),true)
    GOTEST=$(GO) test -v
endif

unexport GOPATH
export GO111MODULE := on
export GOBIN := $(BUILDPATH)/bin

debug:
	@echo "DEBUG=$(DEBUG)"
	@echo "BUILDPATH=$(BUILDPATH)"
	@echo "GOTEST=$(GOTEST)"
	@echo ""

clean:
	@echo "Start cleaning..."
	@rm -rf $(BUILDPATH)/cmd/nihao/nihao
	@echo "Completed cleaning."

fmt:
	@echo "Run $(GOFMT)..."
	@$(GOFMT) -l $(BUILDPATH) | grep -v $(GOEXCLUDES) | xargs -r $(GOFMT) -l -w 2>&1
	@echo "Completed $(GOFMT)."

lint:
	@echo "Run $(GOLINT)..."
	@cd $(BUILDPATH) && $(GOLINT) -set_exit_status `$(GOLIST) ./... | grep -v $(GOEXCLUDES)` 2>&1
	@echo "Completed $(GOLINT)."

sec:
	@echo "Run $(GOSEC)..."
	@cd $(BUILDPATH) && $(GOSEC) -quiet -exclude=G402,G505 ./...
	@echo "Completed $(GOSEC)."

vet:
	@echo "Run $(GOVET)..."
	@cd $(BUILDPATH) && $(GOVET) `$(GOLIST) ./... | grep -v $(GOEXCLUDES)` 2>&1
	@echo "Completed $(GOVET)."

unit-tests:
	@echo "Run unit tests..."
	@cd $(BUILDPATH) && $(GOTEST) ./... -tags=unit -covermode=count -coverprofile=coverage.out 2>&1
	@echo "Completed unit tests."

integration-tests:
	@echo "Run unit tests..."
	@cd $(BUILDPATH) && $(GOTEST) ./... -tags=integration -covermode=count -coverprofile=coverage.out 2>&1
	@echo "Completed unit tests."

tests:
	@echo "Run unit and integration tests..."
	@cd $(BUILDPATH) && $(GOTEST) ./... -tags=unit,integration -covermode=count -coverprofile=coverage.out 2>&1
	@echo "Completed unit tests."

cover:
	@cd $(BUILDPATH) && $(GOTOOL) cover -func=coverage.out 2>&1

build:
	@echo "Building binaries..."
	@cd $(BUILDPATH)/cmd/nihao && $(GOBUILD) -ldflags="-X 'main.Version=$(VERSION)'" -o $(BUILDPATH)/cmd/nihao/nihao
	@echo "Completed building binaries."

all: debug clean fmt vet lint sec tests
