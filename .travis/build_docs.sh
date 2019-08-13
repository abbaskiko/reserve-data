#!/bin/bash
# -*- firestarter: "shfmt -i 4 -ci -w %p" -*-

set -euxo pipefail

pushd ./apidocs
docker-compose up
popd
