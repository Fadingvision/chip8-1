build:
	GOOS=js GOARCH=wasm go build -o chip8.wasm ./cmd/main.go