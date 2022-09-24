.PHONY: clean
clean:
	rm -rf bin/

.PHONY: build
build:
	go build -o bin/asciidonut

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/asciidonut.wasm

.PHONY: copy
copy:
	cp bin/asciidonut.wasm web/assets/wasm/

.PHONY: serve
serve:
	cd docs/ && python3 -m http.server 8080
