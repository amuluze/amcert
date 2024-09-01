TARGET 			= amcert
SERVER   = ./cmd/${TARGET}

GO          = go
GOBUILD     = $(GO) build

.PHONY: amd64
# build amd64
amd64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET) $(SERVER)
	#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-gnu-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(TARGET) $(SERVER)

.PHONY: arm64
# build arm64
arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(TARGET) $(SERVER)
	#CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(TARGET) $(SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/

