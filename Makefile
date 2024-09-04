ifeq ($(GOPATH),)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

.PHONY: build
build:
	tinygo build -o plugin.wasm -scheduler=none -target=wasi ./main.go

.PHONY: run
run:
	envoy -c ./envoy.yaml --concurrency 2 --log-format '%v'