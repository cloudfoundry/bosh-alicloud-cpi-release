#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_SECRET_ACCESS_KEY:?}
: ${ALICLOUD_DEFAULT_REGION:?}
: ${GIT_USER_EMAIL:?}
: ${GIT_USER_NAME:?}
: ${GIT_USER_ID:?}
: ${GIT_USER_PASSWORD:?}
: ${BOSH_REPO_HOST:?}

CURRENT_PATH=$(pwd)
SOURCE_PATH=$CURRENT_PATH/bosh-alicloud-cpi-release
TERRAFORM_PATH=$CURRENT_PATH/terraform
TERRAFORM_MODULE=$SOURCE_PATH/ci/assets/terraform
TERRAFORM_METADATA=$CURRENT_PATH/terraform-metadata
METADATA=metadata
TERRAFORM_VERSION=0.10.0
TERRAFORM_PROVIDER_VERSION=1.2.6


wget -N https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
wget -N https://github.com/alibaba/terraform-provider/releases/download/V${TERRAFORM_PROVIDER_VERSION}/terraform-provider-alicloud_linux-amd64.tgz

mkdir -p ${TERRAFORM_PATH}

unzip -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d ${TERRAFORM_PATH}
tar -xzvf terraform-provider-alicloud_linux-amd64.tgz
mv -f bin/terraform* ${TERRAFORM_PATH}
rm -rf ./bin
export PATH="${TERRAFORM_PATH}:$PATH"

#echo "******** clone terraform template ********"
#
#git clone ${BOSH_REPO_HOST}
#cd ${TERRAFORM_MODULE}
#git checkout -b ${BOSH_REPO_BRANCH} origin/${BOSH_REPO_BRANCH}

cd $CURRENT_PATH/bosh-cpi-src

echo "******** tell docker who am I ********"
git config --global user.email ${GIT_USER_EMAIL}
git config --global user.name ${GIT_USER_NAME}
git config --local -l

ls -la ~/.ssh/


echo "******** clone terraform template ********"
cd $CURRENT_PATH
ls -la ~/.ssh/
git clone ${BOSH_REPO_HOST}
ls -l
cd ${TERRAFORM_MODULE}
git checkout -b ${BOSH_REPO_BRANCH} origin/${BOSH_REPO_BRANCH}


#echo "******** git install expect ********"
#sudo apt-get install expect -y
#
#echo "******** git pull by https ********"
#echo "#!/usr/bin/expect" > git_install.sh
#echo "spawn git fetch https://${GIT_USER_ID}@${BOSH_REPO_HOST} ${BOSH_REPO_BRANCH}" >> git_install.sh
#echo "expect \"Password for 'https://${GIT_USER_ID}@github.com': \"" >> git_install.sh
#echo "send \"${GIT_USER_PASSWORD}\r\"" >> git_install.sh
#echo "expect eof" >> git_install.sh
#echo exit >> git_install.sh
#chmod +x git_install.sh
#./git_install.sh
#rm -rf ./git_install.sh
#
#echo $'\n'
#echo "****** git merge ******"
#git merge FETCH_HEAD

echo $'\n'
echo "Destroy terraform environment......"
terraform init
echo terraform destroy -var alicloud_access_key=${ALICLOUD_ACCESS_KEY_ID} -var alicloud_secret_key=${ALICLOUD_SECRET_ACCESS_KEY} -var alicloud_region=${ALICLOUD_DEFAULT_REGION}  \<\< EOF > terraform_destroy.sh
echo yes >> terraform_destroy.sh
echo EOF >> terraform_destroy.sh
chmod +x terraform_destroy.sh
./terraform_destroy.sh
echo "Destroy terraform environment successfully."
rm -rf ./terraform_destroy.sh

if [ -e ${METADATA} ];
then
    echo "" > $METADATA
fi


function copyToOutput(){

    cp -rf $1/. $2

    cd $2
    ls -la
    echo "******** show git repo info ********"
    git remote -v
    git branch

    echo "******** show git repo info ********"
    git remote -v
    git branch

    git status | sed -n 'p' |while read LINE
    do
        echo "echo LINE: $LINE"
        if [[ $LINE == HEAD*detached* ]];
        then
            echo "****** fix detached branch ******"
            read -r -a Words <<< $LINE

            git status | sed -n 'p' |while read LI
            do
                echo "echo LI: $LI"
                if [[ $LI == Changes*not*staged*for*commit* ]];
                then
                    git add .
                    git commit -m 'destroy environment commit on detached'
                    git branch temp
                    git checkout ${BOSH_REPO_BRANCH}
                    git merge temp
                    git branch
                    git branch -d temp
                fi
            done
            break
        fi
    done

    echo "******* git status ******"
    git status

    git status | sed -n '$p' |while read LINE
    do
        echo $LINE
        if [[ $LINE != nothing*clean ]];
        then
            echo $LINE
            git add .
            git commit -m 'destroy environment commit -a'
            return 0
        fi
    done

    git status
    return 0
}

echo "Copy to output ......"
copyToOutput ${SOURCE_PATH} ${TERRAFORM_METADATA}