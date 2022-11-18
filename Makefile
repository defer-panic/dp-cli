PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

TOOLKIT = $(PROJECT_BIN)/toolkit

.PHONY: toolkit
toolkit: 
	go build -o $(TOOLKIT) .

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
