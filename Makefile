PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

TOOLKIT = $(PROJECT_BIN)/toolkit

.PHONY: toolkit
toolkit: 
	go build -o $(TOOLKIT) .
