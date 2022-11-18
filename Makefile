PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.50.1

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml


TESTS_WD = $(PROJECT_DIR)/tests

# === Test ===
.PHONY: test
test:
	mkdir -p $(TESTS_WD)
	go test -v --timeout=5m --covermode=count --coverprofile=$(TESTS_WD)/profile.cov_tmp ./...
	cat $(TESTS_WD)/profile.cov_tmp > $(TESTS_WD)/profile.cov

.PHONY: test-coverage
test-coverage: test
	go tool cover --func=$(TESTS_WD)/profile.cov 

.PHONY: test-coverage-html
test-coverage-html: test
	go tool cover --html=$(TESTS_WD)/profile.cov 


# === Build ===
.PHONY: build
build:
	go build -o $(PROJECT_BIN)/dp-cli ./cmd/
