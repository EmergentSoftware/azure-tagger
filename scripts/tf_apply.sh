#!/bin/bash

set -e

# Set variables
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."
source vars.sh
cd tf
terraform apply tf.plan
