#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_SECRET_ACCESS_KEY:?}
: ${ALICLOUD_DEFAULT_REGION:?}
: ${DESTROY_ENVIRONMENT:?}

CURRENT_PATH=$(pwd)
SOURCE_PATH=$CURRENT_PATH/bosh-cpi-src
TERRAFORM_PATH=$CURRENT_PATH/terraform
TERRAFORM_MODULE=$SOURCE_PATH/ci/assets/terraform
TERRAFORM_METADATA=$CURRENT_PATH/terraform-metadata
METADATA=metadata
TERRAFORM_VERSION=0.10.0
TERRAFORM_PROVIDER_VERSION=1.2.4

wget -N -r -np http://releases.hashicorp.com/terraform/$TERRAFORM_VERSION/terraform_$TERRAFORM_VERSION_linux_amd64.zip
wget -N -r -np http://github.com/alibaba/terraform-provider/releases/download/V$TERRAFORM_PROVIDER_VERSION/terraform-provider-alicloud_linux-amd64.tgz
mkdir -p ${TERRAFORM_PATH}
unzip -o terraform_$TERRAFORM_VERSION_linux_amd64.zip -d ${TERRAFORM_PATH}
tar -xzvf terraform-provider-alicloud_linux-amd64.tgz
mv -f bin/terraform* ${TERRAFORM_PATH}
rm -rf ./bin
export PATH="${TERRAFORM_PATH}:$PATH"

cd ${TERRAFORM_MODULE}

rm -rf ${METADATA}
touch ${METADATA}

terraform init && terraform apply -var alicloud_access_key=${ALICLOUD_ACCESS_KEY_ID} -var alicloud_secret_key=${ALICLOUD_SECRET_ACCESS_KEY} -var alicloud_region=${ALICLOUD_DEFAULT_REGION}

if [ ${DESTROY_ENVIRONMENT} == "true" ]
then

   echo "Destroy terraform environment......"
   terraform init
   echo terraform destroy -var alicloud_access_key=${ALICLOUD_ACCESS_KEY_ID} -var alicloud_secret_key=${ALICLOUD_SECRET_ACCESS_KEY} -var alicloud_region=${ALICLOUD_DEFAULT_REGION}  \<\< EOF > terraform_destroy.sh
   echo yes >> terraform_destroy.sh
   echo EOF >> terraform_destroy.sh
   chmod +x terraform_destroy.sh
   ./terraform_destroy.sh

   rm -rf ./terraform_destroy.sh
fi

function copyToOutput(){

    cp -rf $1/. $2

    cd $2
    git config --global user.email "guimin.hgm@alibaba-inc.com "
    git config --global user.name "He Guimin"
    git config --local -l

    git status | sed -n '$p' | while read LINE
    do
        echo $LINE
        if [[ $LINE != nothing* ]];
        then
            git add .
            git commit -m 'commit metadata'
        fi
    done
    return 0
}

if [ ! -e "./terraform.tfstate" ];
then
    copyToOutput ${SOURCE_PATH} ${TERRAFORM_METADATA}
    echo "./terraform.tfstate is not exist and then quit."
    exit 0
fi

terraform state list > all_state

cat all_state | while read LINE
do
    if [ $LINE == "alicloud_vswitch.default" ];
    then
        terraform state show $LINE >> $METADATA
        cat $METADATA | while read line
        do
          echo $line
          if [[ $line == id* ]];
          then
              echo vswitch_$line >> $METADATA
          fi
        done
        sed -i '/^id/d' $METADATA
    fi
    if [ $LINE == "alicloud_security_group.sg" ];
    then
        terraform state show $LINE >> $METADATA
        cat $METADATA | while read line
        do
          echo $line
          if [[ $line == id* ]];
          then
              echo security_group_$line >> $METADATA
          fi
        done
        sed -i '/^id/d' $METADATA
    fi
done

sed -i 's/=/:/g' $METADATA

copyToOutput ${SOURCE_PATH} ${TERRAFORM_METADATA}
