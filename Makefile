TARGET 			= amcert
SERVER   = ./cmd/${TARGET}

GO          = go
GOBUILD     = $(GO) build

.PHONY: wire
# generate wire
wire:
	wire ./service/

.PHONY: amd64
# build amd64
amd64: wire
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-s -w"  -o $(TARGET) $(SERVER)

.PHONY: arm64
# build arm64
arm64: wire
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "-s -w" -o $(TARGET) $(SERVER)
