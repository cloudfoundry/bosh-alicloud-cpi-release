#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh

: ${access_key:?}
: ${secret_key:?}
: ${region:?}
: ${env_name:=""}
: ${public_key:?}
: ${delete_on_failure:= true}
: ${generate_random_name:= false}
: ${action:=""}
: ${terraform_role_arn:?}
: ${terraform_source:?}
: ${output_module:="metadata"}
# Remote state parameters
#: ${remote_state_access_key:=${access_key}}
#: ${remote_state_secret_key:=${secret_key}}
: ${remote_state_region:=${region}}
: ${remote_state_bucket:?}
: ${remote_state_file_path:="terraform-state"}
: ${remote_state_file_name:=""}

if [[ ${env_name} == "" ]]; then
    if [[ ${generate_random_name} = true ]]; then
        env_name="bosh-concourse-$(echo $RANDOM)"
    fi
fi

if [[ ${remote_state_file_name} == "" ]]; then
    remote_state_file_name="${env_name}.tfstate"
fi

CURRENT_PATH=${PWD}
ls -l ./
echo -e "ls -l CURRENT_PATH"
ls -l ${CURRENT_PATH}

SOURCE_PATH=$CURRENT_PATH/bosh-alicloud-cpi-release
TERRAFORM_METADATA=$CURRENT_PATH/environment
output_path=${CURRENT_PATH}/environment
TERRAFORM_VERSION=1.3.3

wget -qN https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
apt-get install unzip
unzip -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin

# 调用 AssumeRole API 获取临时凭证
response=$(aliyun sts AssumeRole --RoleArn ${terraform_role_arn} --RoleSessionName "bosh-cpi-e2e-test" --access-key-id ${access_key} --access-key-secret ${secret_key})

# 解析返回结果中的 AccessKeyId, AccessKeySecret 和 SecurityToken
ACCESS_KEY_ID=$(echo $response | jq -r '.Credentials.AccessKeyId')
ACCESS_KEY_SECRET=$(echo $response | jq -r '.Credentials.AccessKeySecret')
SECURITY_TOKEN=$(echo $response | jq -r '.Credentials.SecurityToken')

# 检查是否成功获取到凭证
if [ -z "$ACCESS_KEY_ID" ] || [ -z "$ACCESS_KEY_SECRET" ]; then
  echo "Failed to get credentials."
  exit 1
fi

pushd ${terraform_source}
    terraform init \
        -backend-config="access_key=${ACCESS_KEY_ID}" \
        -backend-config="secret_key=${ACCESS_KEY_SECRET}" \
        -backend-config="security_token=${SECURITY_TOKEN}" \
        -backend-config="region=${remote_state_region}" \
        -backend-config="bucket=${remote_state_bucket}" \
        -backend-config="prefix=${remote_state_file_path}" \
        -backend-config="key=${remote_state_file_name}"

    set +e

    echo -e "******* Build terraform environment ******** \n"

    if [[ ${action} == "destroy" ]]; then
        echo -e "******** Try to delete environment ********\n"
        terraform apply -destroy -auto-approve -var access_key=${ACCESS_KEY_ID} -var secret_key=${ACCESS_KEY_SECRET} -var security_token=${SECURITY_TOKEN} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
    else
        echo -e "******** Try to build environment ********\n"
        terraform apply --auto-approve -var access_key=${ACCESS_KEY_ID} -var secret_key=${ACCESS_KEY_SECRET} -var security_token=${SECURITY_TOKEN} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
        if [[ $? -eq 0 ]]; then
            echo -e "******** Build terraform environment successfully ******** \n"
            ls -al
            echo "{" > ${output_path}/${output_module}
            terraform output >> ${output_path}/${output_module}
            sed -i '2,$s/^/"/g; 2, $s/$/,/g; $s/,//g; 2,$s/ = /": /g' ${output_path}/${output_module}
            echo "}" >> ${output_path}/${output_module}
        elif [[ ${delete_on_failure} = true ]]; then
            echo -e "******** Destroy terraform environment... ******** \n"
            terraform apply -destroy -auto-approve -var access_key=${ACCESS_KEY_ID} -var secret_key=${ACCESS_KEY_SECRET} -var security_token=${SECURITY_TOKEN} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
        fi
    fi

    set -e

    echo -e "******** Write metadata successfully ********\n"

popd
