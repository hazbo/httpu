GOCMD         := go
GOBUILD       := $(GOCMD) build
GOTEST        := $(GOCMD) test
GOGET         := $(GOCMD) get
GOINSTALL     := $(GOCMD) install
BUILD_WD      := cmd/httpu
BIN_PATH      := build
BIN_NAME      := httpu
BIN_OUT       := $(BIN_PATH)/$(BIN_NAME)
BIN_OUT_LINUX := $(BIN_OUT)-linux
GIT_COMMIT    := `git rev-parse HEAD`
GIT_BRANCH    := `git rev-parse --abbrev-ref HEAD`
META_PACKAGE  := github.com/hazbo/httpu/meta.Commit

# Running tests on build for now as it's quick and helpful during development.
all: deps build

build: build/httpu

build/httpu:
	cd $(BUILD_WD) && \
		$(GOBUILD) -v \
		-o ../../$(BIN_OUT) \
		-ldflags "-X $(META_PACKAGE)=$(GIT_COMMIT)"

test:
	$(GOTEST) -cover -v \
		./ \
		./resource \
		./resource/request \
		./utils/varparser

run: $(BIN_OUT)
	./$(BIN_OUT) version

clean:
	rm -f $(BIN_OUT)
	rm -f $(BIN_OUT_LINUX)

deps:
	$(GOGET) -v ./...

install: build
	cp $(BIN_OUT) $(GOPATH)/bin/$(BIN_NAME)

.PHONY: all build test clean deps install