# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOFMT = gofmt

# Build target
build:
	$(GOBUILD) -o nuwa-terminal ./cmd

# Clean target
clean:
	$(GOCLEAN)
	rm -f nuwa-terminal

# Test target
test:
	$(GOTEST) -v ./...

# Get dependencies
deps:
	$(GOGET) -v ./...

fmt:
	$(GOFMT) -w .


# Default target
default: build
