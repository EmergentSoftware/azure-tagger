#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

docker run --rm \
  --volume=$PWD:/root \
  --workdir=/root \
  golang \
  bash -c 'apt update && apt install -y zip && git config --global --add safe.directory /root && ./scripts/build.local.sh' "$@"
