CHS_ENV_HOME ?= $(HOME)/.chs_env
TESTS        ?= ./...

bin          := emergency-auth-code-api 
chs_envs     := $(CHS_ENV_HOME)/global_env $(CHS_ENV_HOME)/emergency-auth-code-api/env
source_env   := for chs_env in $(chs_envs); do test -f $$chs_env && . $$chs_env; done
xunit_output := test.xml
lint_output  := lint.txt

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

.PHONY: all
all: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt $(bin)

$(bin):
	CGO_ENABLED=0 go build -o ./$(bin)

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test $(TESTS) -run 'Unit' -coverprofile=coverage.out

.PHONY: test-integration
test-integration:
	$(source_env); go test $(TESTS) -run 'Integration'

.PHONY: clean
clean:
	go mod tidy
	rm -f ./$(bin) ./$(bin)-*.zip $(test_path) build.log

.PHONY: package
package:
ifndef version
	$(error No version given. Aborting)
endif
	$(info Packaging version: $(version))
	$(eval tmpdir := $(shell mktemp -d build-XXXXXXXXXX))
	cp ./$(bin) $(tmpdir)
	cp ./routes.yaml $(tmpdir)
	cp ./start.sh $(tmpdir)
	cp -r ./assets  $(tmpdir)/assets
	cd $(tmpdir) && zip -r ../$(bin)-$(version).zip $(bin) start.sh routes.yaml assets
	rm -rf $(tmpdir)

.PHONY: dist
dist: clean build package

.PHONY: xunit-tests
xunit-tests: GO111MODULE = off
xunit-tests:
	go get github.com/tebeka/go2xunit
	@set -a; go test -v $(TESTS) -run 'Unit' | go2xunit -output $(xunit_output)

.PHONY: lint
lint: GO111MODULE = off
lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter ./... > $(lint_output); true
