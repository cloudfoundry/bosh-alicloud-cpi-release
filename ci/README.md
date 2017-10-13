## 概览
![image](https://yqfile.alicdn.com/62bee4cf860524c63de9590a4a5d604ae77baadc.png)

## Codepipeline 简介
1. Concourse搭建在阿里云容器服务上
 
   + URL: http://47.90.16.199:8080/teams/main/pipelines/ali-cli
   + 用户名密码：
   
2. Codepipeline是一个任务编排引擎，按照一定的规则执行SHELL脚本。总共包含5个Job，每个Job下面包含N个Task。

## 项目详情

1. 代码位置

   + codepipeline的代码，分布在三个仓库
     - bosh-alicloud-cpi-release工程下的ci目录，放置pipeline主要代码
     - bosh-cpi-certification，放置通用的Task脚本
     - bosh-deployment, 放置用于部署的manifest
     
2. pipeline outline

   + pipeline模板，由三部分组成
     - jobs: 定义所有的工作
     - resources: 定义每个Job的输入、输出资源
     - shared: 定义可共用的Task, 不同的Job可能会执行相同的Task
     
3. Jobs简介

   + build-candidate, 打CPI release包
     - 需要三个资源：
       - bosh-cpi-src, bosh-alicloud-cpi-release源码
       - ruby-cpi-blobs, 创建CPI release的bundler Dependency
     - 包含一个输出资源
       - ruby-cpi-blobs, 就是CPI release的构建物.
       - ps: 不支持OSS的存储，临时放到github上
     - 包含一个任务build
       - 运行时用boshcpi/aws-cpi-release这个镜像起一个container
       - 会运行UT
       - 安装CPI bundler依赖
       - 打CPI release包
       - 把release包放到输出目录
       - todo: 打阿里云docker镜像, 把bosh-cli, cpi依赖也打进去 @箫竹
       - todo: 完善ut @埃兰
       - todo: bosh2打出来的release不能用 @思皓
   + integration cpi-release集成测试
      - 需要两个资源:
        - 第1步的CPI release
        - bosh-alicloud-cpi-release源码
      - 引用了三个公共Task: 创建IaaS资源和销毁IaaS资源
      - 包含一个任务integration
        - 获取IaaS层的metadata信息
        - 执行integration测试
        - todo: 跟Terrafom联调 @箫竹
        - todo: 完善integration脚本 @埃兰
   + bats, 创建Bosh, bats测试
     - 需要几个资源：
        - 第1步的CPI release
        - bosh-release, bosh-alicloud-cpi-release源码, bosh-stemcell
        - bosh-deployment源码
        - bosh-cpi-certification源码
     - 引用了三个公共Task: 创建IaaS资源和销毁IaaS资源
     - 包含三个Task
       - prepare-director
         - 获取第1步的metadata,生成bosh的director.yml文件
         - todo: 需要跟@箫竹联调
       - deploy-director
         - bosh create-env创建bosh
         - todo: 支持ssl_ca @厚泼
         - todo: 支持jumpbox @厚泼
         - todo: 联调 @箫竹
       - run-bats
         - 执行bats脚本
         - todo: bats脚本@埃兰
   + end-2-end，创建Bosh, 部署CF
     - 需要几个资源：
        - 第1步的CPI release
        - bosh-release, bosh-alicloud-cpi-release源码, bosh-stemcell
        - bosh-deployment源码, heavy-stemcell
        - bosh-cpi-certification源码
        - todo: heavy-stemcell测试需要 @思皓评估一下
     - 引用了三个公共Task: 创建IaaS资源和销毁IaaS资源
     - 包含三个Task
       - prepare-director
         - 获取第1步的metadata,生成bosh的director.yml文件
       - deploy-director
         - bosh create-env创建bosh
       - run-e2e
         - 部署CF
         - 执行errand测试
         - todo errand测试@厚泼



## Concourse 运维
   
2. 常用命令

   + 安装fly CLI, 在控制台下载对应版本的安装包
   + 登陆Concourse
   
   ```
   fly -t ali login -c http://REPLACE_WITH_YOUR_IP:8080
   ```
   + 修改密码
   
   ```
   fly -t ali set-team --team-name main \
       --basic-auth-username concourse \
       --basic-auth-password REPLACE_WITH_YOUR_PWD
   ```
   + 更新pipeline
   
   ```
   fly -t ali set-pipeline -p cpi-release -c pipeline-develop.yml --load-vars-from vars-pipeline-develop.yml
   ```
   + 删除pipeline
   
   ```
   fly -t ali destroy-pipeline -p cpi-release

   ```