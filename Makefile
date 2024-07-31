os-archs=darwin:amd64 darwin:arm64 freebsd:amd64 linux:amd64 linux:arm:7 linux:arm:5 linux:arm64 windows:amd64 windows:arm64 linux:mips64 linux:mips64le linux:mips:softfloat linux:mipsle:softfloat linux:riscv64 android:arm64

restore:
	go get -C ./cmd
	go get -C ./internal

build: restore
	go build -C ./cmd -ldflags="-linkmode external -extldflags -static -w -s" -o "../mdrop"

