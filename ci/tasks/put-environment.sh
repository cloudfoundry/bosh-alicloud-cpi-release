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
: ${output_statefile:=""}
: ${output_module:="metadata"}
# OSS bucket request args
: ${access_key_id:=${access_key}}
: ${secret_access_key:=${secret_key}}
: ${region_name:=${region}}
: ${bucket:?}
: ${bucket_path:="terraform-state"}

source bosh-cpi-src/ci/tasks/utils.sh

if [[ ${env_name} == "" ]]; then
    if [[ ${generate_random_name} = true ]]; then
        env_name="bosh-concourse-$(echo $RANDOM)"
    fi
fi

if [[ ${output_statefile} == "" ]]; then
    output_statefile="${env_name}.tfstate"
fi

CURRENT_PATH=${PWD}
ls -l ./
echo -e "ls -l CURRENT_PATH"
ls -l ${CURRENT_PATH}

SOURCE_PATH=$CURRENT_PATH/bosh-alicloud-cpi-release
TERRAFORM_METADATA=$CURRENT_PATH/environment
output_path=${CURRENT_PATH}/environment
TERRAFORM_VERSION=0.11.4

wget -qN https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
unzip -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin

pushd ${terraform_source}
    terraform init
    set +e

    aliyun oss cp oss://${bucket}/${output_statefile} ${output_statefile} \
        --access-key-id ${access_key_id} \
        --access-key-secret ${secret_access_key} \
        --region ${region_name} \
        --force

    echo -e "******* Build terraform environment ******** \n"

    if [[ ${action} == "destroy" ]]; then
        echo -e "******** Try to delete environment ********\n"
        terraform destroy -force -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}" -state=${output_statefile}
    else
        echo -e "******** Try to build environment ********\n"
        terraform apply --auto-approve -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var "public_key=${public_key}" -state=${output_statefile}
        if [[ $? -eq 0 ]]; then
            echo -e "******** Build terraform environment successfully ******** \n"
            ls -al
            echo "{" > ${output_path}/${output_module}
            terraform output -state=${output_statefile} >> ${output_path}/${output_module}
            sed -i '2,$s/^/"/g; 2, $s/$/",/g; $s/,//g; 2,$s/ = /": "/g' ${output_path}/${output_module}
            echo "}" >> ${output_path}/${output_module}
        elif [[ ${delete_on_failure} = true ]]; then
            echo -e "******** Destroy terraform environment... ******** \n"
            terraform destroy -force -var access_key=${access_key} -var secret_key=${secret_key} -var region=${region} -var env_name=${env_name} -var 'public_key=${public_key}' -state=${output_statefile}
        fi
    fi

    set -e

    aliyun oss cp ${output_statefile} oss://${bucket}/${output_statefile} \
        --access-key-id ${access_key_id} \
        --access-key-secret ${secret_access_key} \
        --region ${region_name} \
        --force

    echo -e "******** Write metadata successfully ********\n"

popd
