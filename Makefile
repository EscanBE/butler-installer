###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "Building Escan Butler binary..."
	@go build -mod=readonly -o build/ebid.exe ./cmd/ebid
	@echo "Builded successfully"
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "Installing Escan Butler binary..."
	@go install -mod=readonly ./cmd/ebid
	@echo "Installed successfully"
.PHONY: install