#!/usr/bin/env bash

set -e

TERRAFORM_PATH=$(pwd)
echo TERRAFORM_PATH
export PATH="${TERRAFORM_PATH}:$PATH"
terraform init && terraform plan