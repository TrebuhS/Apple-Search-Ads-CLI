BINARY_NAME=asa-cli
BUILD_DIR=.
GOBIN=$(shell go env GOPATH)/bin
INSTALL_DIR=$(GOBIN)

.PHONY: build install clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
