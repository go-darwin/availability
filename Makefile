GO_TAGS := osusergo,netgo
GO_LDFLAGS := -s -w "-extldflags=-static"
GO_INSTALLSUFFIX := netgo
GO_FLAGS = -trimpath -tags='${GO_TAGS}' -ldflags='${GO_LDFLAGS}' -installsuffix='${GO_INSTALLSUFFIX}'

GO_TEST ?= go test -v
GO_TEST_FLAGS ?= -race -count=1
GO_TEST_FUNC ?= .
GO_TEST_PACKAGE ?= ./...

JOBS := $(shell getconf _NPROCESSORS_CONF)
ifeq ($(CIRCLECI),true)
ifeq (${GO_OS},linux)
	# https://circleci.com/changelog#container-cgroup-limits-now-visible-inside-the-docker-executor
	JOBS := $(shell echo $$(($$(cat /sys/fs/cgroup/cpu/cpu.shares) / 1024)))
	GO_TEST_FLAGS+=-p=${JOBS} -cpu=${JOBS}
endif
endif

##@ gen

define gen
@rm -f $2
printf '// Copyright 2021 The Go Darwin Authors\n// SPDX-License-Identifier: BSD-3-Clause\n\n//go:build darwin\n// +build darwin\n\n' > $2
@go tool cgo -godefs $1 $3 | gofmt -s | tee -a $2 > /dev/null 2>&1
@sed -i 's|${CURDIR}/||' $2
@rm -rf _obj
endef

ztypes_darwin.go: types_darwin.go
	$(call gen,$^,$@)

##@ test

define go_test
${GO_TEST} $(strip ${GO_FLAGS}) ${GO_TEST_FLAGS} -run=${GO_TEST_FUNC} ${GO_TEST_PACKAGE}
endef

.PHONY: test
test:  ## Run test.
	$(call go_test)

.PHONY: coverage
coverage: GO_TEST_FLAGS+=-covermode=atomic -coverpkg=./... -coverprofile=coverage.out
coverage: tools/gotestsum  ## Run test and collect coverages.
	$(call go_test)

##@ clean

.PHONY: clean
clean:  ## Cleanups binaries and extra files in the package.
	$(call target)
	@$(RM) -rf *.out *.test *.prof *.txt **/_obj ${TOOLS_BIN}

##@ help

.PHONY: help
help:  ## Show make target help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: env/% env
env:  ## Print the value of MAKEFILE_VARIABLE. Use `make env/MAKEFILE_VARIABLE`.
env/%:
	@echo $($*)
