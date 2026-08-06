#!/usr/bin/env bash

set -e

: ${region:?}
: ${env_name:=""}
: ${public_key:?}
: ${delete_on_failure:= true}
: ${generate_random_name:= false}
: ${action:=""}
: ${terraform_role_arn:?}
: ${terraform_source:?}
: ${output_module:="metadata"}
: ${director_role_name:=""}
# Remote state parameters
: ${remote_state_region:=${region}}
: ${remote_state_bucket:?}
: ${remote_state_file_path:="terraform-state"}
: ${remote_state_file_name:=""}

source bosh-cpi-src/ci/tasks/credentials.sh

# Auto-detect Concourse worker's outbound IP for security group rules
CONCOURSE_WORKER_IP=$(curl -s --max-time 5 ifconfig.me || echo "")
if [[ -n "${CONCOURSE_WORKER_IP}" ]]; then
    echo "Detected Concourse worker outbound IP: ${CONCOURSE_WORKER_IP}"
fi

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

# The CLI comes from the aliyun-cli input, which pins a version in a bucket the
# operator controls, rather than an unpinned download unpacked into /usr/bin.
ensure_aliyun_cli

# Start from the worker's instance role and assume the provisioning role. Both
# create and destroy do this independently, so a destroy still authenticates
# after the task that created the environment failed.
assume_pipeline_role "${terraform_role_arn}" "tf"

pushd ${terraform_source}
    # The provider and the OSS backend read the assumed credentials from the
    # environment. Passing them as -var/-backend-config would put them in the
    # process arguments and in Terraform's debug output.
    terraform init \
        -backend-config="region=${remote_state_region}" \
        -backend-config="bucket=${remote_state_bucket}" \
        -backend-config="prefix=${remote_state_file_path}" \
        -backend-config="key=${remote_state_file_name}"

    terraform_vars=(
        -var "region=${region}"
        -var "env_name=${env_name}"
        -var "public_key=${public_key}"
        -var "concourse_worker_ip=${CONCOURSE_WORKER_IP}"
        -var "director_role_name=${director_role_name}"
    )

    # Terraform failures must not be swallowed. `set +e` below lets a failed
    # apply fall through to the delete-on-failure path, so every exit code is
    # captured and the task fails at the end if any of them was non-zero.
    # Reporting success after a failed destroy is what lets an environment
    # survive unnoticed and collide with the next build.
    terraform_status=0
    set +e

    echo -e "******* Build terraform environment ******** \n"

    if [[ ${action} == "destroy" ]]; then
        echo -e "******** Try to delete environment ********\n"
        terraform apply -destroy -auto-approve "${terraform_vars[@]}"
        terraform_status=$?
    else
        echo -e "******** Try to build environment ********\n"
        terraform apply --auto-approve "${terraform_vars[@]}"
        apply_status=$?
        if [[ ${apply_status} -eq 0 ]]; then
            echo -e "******** Build terraform environment successfully ******** \n"
            ls -al
            # Later jobs locate every resource through this file, and cleanup
            # cannot find them if it is empty or truncated, so a failure here is
            # a failure of the task.
            {
                echo "{"
                terraform output
            } > ${output_path}/${output_module}
            output_status=$?
            if [[ ${output_status} -ne 0 ]]; then
                terraform_status=${output_status}
                echo "terraform output failed; the metadata file is unusable" >&2
            else
                sed -i '2,$s/^/"/g; 2, $s/$/,/g; $s/,//g; 2,$s/ = /": /g' ${output_path}/${output_module}
                sed_status=$?
                echo "}" >> ${output_path}/${output_module}
                if [[ ${sed_status} -ne 0 ]]; then
                    terraform_status=${sed_status}
                    echo "rewriting the metadata file failed" >&2
                fi
            fi
        else
            terraform_status=${apply_status}
            if [[ ${delete_on_failure} = true ]]; then
                echo -e "******** Destroy terraform environment... ******** \n"
                terraform apply -destroy -auto-approve "${terraform_vars[@]}"
                destroy_status=$?
                if [[ ${destroy_status} -ne 0 ]]; then
                    echo "Destroy after a failed apply also failed; the environment may still exist" >&2
                fi
            fi
        fi
    fi

    set -e

popd

# The metadata file is published as a task output and consumed by later jobs, so
# it must not carry credentials.
assert_no_credentials_in_file "${output_path}/${output_module}"

if [[ ${terraform_status} -ne 0 ]]; then
    echo "terraform ${action:-apply} exited ${terraform_status}" >&2
    exit ${terraform_status}
fi

echo -e "******** Write metadata successfully ********\n"
