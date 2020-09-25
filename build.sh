#!/bin/bash
set -e
mkdir -p build
cp ./frontend/* ./build/
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./build/script.js

export GOPATH=$(go env GOPATH)
$GOPATH/bin/statik -src=./include/ -dest=./wasm/

GOOS=js GOARCH=wasm go build -o build/frontend.wasm ./wasm/
