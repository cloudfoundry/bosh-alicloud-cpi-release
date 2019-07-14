#!/usr/bin/env bash

set -e

: ${access_key:?}
: ${secret_key:?}
: ${region:?}
: ${env_name:=""}
: ${public_key:?}
: ${delete_on_failure:= true}
: ${generate_random_name:= false}
: ${action:=""}
: ${terraform_source:?}
: ${output_module:="metadata"}
# Remote state parameters
: ${remote_state_access_key:=${access_key}}
: ${remote_state_secret_key:=${secret_key}}
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
TERRAFORM_VERSION=0.12.2

wget -qN https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
unzip -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin

pushd ${terraform_source}
    terraform init \
        -backend-config="access_key=${remote_state_access_key}" \
        -backend-config="secret_key=${remote_state_secret_key}" \
        -backend-config="region=${remote_state_region}" \
        -backend-config="bucket=${remote_state_bucket}" \
        -backend-config="prefix=${remote_state_file_path}" \
        -backend-config="key=${remote_state_file_name}"

    set +e

    echo -e "******* Build terraform environment ******** \n"

    if [[ ${action} == "destroy" ]]; then
        echo -e "******** Try to delete environment ********\n"
        terraform destroy -force -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
    else
        echo -e "******** Try to build environment ********\n"
        terraform apply --auto-approve -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
        if [[ $? -eq 0 ]]; then
            echo -e "******** Build terraform environment successfully ******** \n"
            ls -al
            echo "{" > ${output_path}/${output_module}
            terraform output >> ${output_path}/${output_module}
            sed -i '2,$s/^/"/g; 2, $s/$/",/g; $s/,//g; 2,$s/ = /": "/g' ${output_path}/${output_module}
            echo "}" >> ${output_path}/${output_module}
        elif [[ ${delete_on_failure} = true ]]; then
            echo -e "******** Destroy terraform environment... ******** \n"
            terraform destroy -force -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}"
        fi
    fi

    set -e

    echo -e "******** Write metadata successfully ********\n"

popd
