.PHONY: build
build:
	go build -o bin/asciidonut

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/asciidonut-wasm
