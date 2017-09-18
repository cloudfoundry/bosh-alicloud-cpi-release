#!/usr/bin/env bash

set -e

TERRAFORM_PATH=$(pwd)
export PATH="${TERRAFORM_PATH}:$PATH"
cd ../terraform-module

terraform init && terraform plan