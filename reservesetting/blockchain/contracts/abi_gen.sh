#!/usr/bin/env bash

set -euo pipefail

readonly src_dir=/tmp/reserve-stats-abi-gen/src/github.com/ethereum
export GOPATH=/tmp/reserve-stats-abi-gen
export PATH=$GOPATH/bin:$PATH

mkdir -p "$src_dir"
cd "$src_dir"
[[ -d go-ethereum ]] || git clone https://github.com/ethereum/go-ethereum.git
go install -v github.com/ethereum/go-ethereum/cmd/abigen


abigen -abi "$OLDPWD"/wrapper.abi -pkg contracts -type Wrapper -out "$OLDPWD"/wrapper_abi.go

