BINARY_SUFFIX ?= ""

restore:
	go get -C ./cmd/mdrop-client
	go get -C ./cmd/mdrop-tunnel-tools
	go get -C ./internal

build-client: restore
	go build -C ./cmd/mdrop-client -ldflags="-linkmode external -extldflags -static -w -s" -o "../../mdrop${BINARY_SUFFIX}"

build-tunnel: restore
	go build -C ./cmd/mdrop-tunnel-tools -ldflags="-linkmode external -extldflags -static -w -s" -o "../../mdrop-tunnel"

build-client-general: restore
	go build -C ./cmd/mdrop-client -ldflags="-extldflags -static -w -s" -o "../../mdrop${BINARY_SUFFIX}"

build-tunnel-general: restore
	go build -C ./cmd/mdrop-tunnel-tools -ldflags="-extldflags -static -w -s" -o "../../mdrop-tunnel"
