#!/bin/bash

set -e

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."
export os_list=linux && \
export arch_list=amd64 && \
scripts/build.local.sh
