# Nautes CLI

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![golang](https://img.shields.io/badge/golang-v1.20-brightgreen)](https://go.dev/doc/install)
[![version](https://img.shields.io/badge/version-v0.3.9-green)]()

CLI 项目通过封装 API Server 的 REST API 提供了一个简单的命令行工具，用于简化用户使用 API 的操作。

## 功能简介

CLI 目前支持的 API 包括：集群、产品、环境、项目、代码库、代码库权限、流水线运行时、和部署运行时。

CLI 的执行文件为 nautes，包含以下子命令：

- apply：通过 `-f` 接收一个文件参数，新增或修改文件中声明的所有实体，操作顺序为：集群、产品、环境、项目、代码库、代码库权限、流水线运行时、部署运行时。
- remove：通过 `-f` 接收一个文件参数，删除文件中声明的所有实体，操作顺序为：部署运行时、流水线运行时、代码库权限、代码库、项目、环境、产品、集群。

以上两个子命令可以通过添加 `-i` 参数，跳过 API 的合规性校验，强制执行请求。

CLI 还包含以下参数标志：

- -t, --token：认证所需 token，目前只支持 GitLab AccessToken。
- -s, --api-server：API Server URL，您可以在[安装程序](https://nautes.io/guide/user-guide/installation.html#%E6%9F%A5%E7%9C%8B%E5%AE%89%E8%A3%85%E7%BB%93%E6%9E%9C)的输出目录中找到这个地址。

您可以在 `examples` 目录下找到参数文件的模板，其中：

- demo-cluster-host.yaml：宿主集群模板
- demo-cluster-physical-worker-pipeline.yaml：物理流水线运行时集群模板
- demo-cluster-physical-worker-deployment.yaml：物理部署运行时集群模板
- demo-cluster-virtual-worker-pipeline.yaml：虚拟流水线运行时集群模板
- demo-cluster-virtual-worker-deployment.yaml：虚拟部署运行时集群模板
- demo-product.yaml：可被两种运行时公用的基础实体的模板
- demo-pipeline.yaml：流水线运行时相关实体的模板
- demo-deployment.yaml：部署运行时相关实体的模板

### CLI 扩展子命令

比如：nautes get cr -p demo-101 -t xxxxxx -s xxxxxx 执行结果是列出当前产品下所有的coderepo资源，可以再简约一下

> 设置环境变量，提升访问效率，API_SERVER 是请求的地址，GIT_TOKEN是 Git 仓库的 access token，PRODUCT 是产品名

- export API_SERVER=http://127.0.0.1:8000

- export GIT_TOKEN=glpat-yYTvmC9Vnzom5k2NuzUU

- export PRODUCT=demo-101

> 一次查询多个 CodeRepo 资源：nautes get cr coderepo-name-101 coderepo-name-102 coderepo-name-103 -o json

> 查询 DeploymentRuntime 资源列表：nautes get dr

> 一次删除多个 Project 资源：nautes delete pro project-101 project-102

| command                              | short command   | resource               | args  | flags | example                                        |
|--------------------------------------|-----------------|------------------------|-------|-------|------------------------------------------------|
| nautes product get                   | prod,prods      | product                | name  |       | nautes prod get product-name                   |
| nautes product delete                | prod,prods      | product                | name  |       | nautes prod delete product-name                |
| nautes cluster get                   | cls             | cluster                | name  |       | nautes cls get cluster-name                    |
| nautes cluster delete                | cls             | cluster                | name  |       | nautes cls delete cluster-name                 |
| nautes environment get               | env,envs        | environment            | name  | -p    | nautes env get env-name -p product-name        |
| nautes environment delete            | env,envs        | environment            | name  | -p    | nautes env delete env-name -p product-name     |
| nautes project get                   | pro,pros,proj   | project                | name  | -p    | nautes pro get project-name -p product-name    |
| nautes project delete                | pro,pros,proj   | project                | name  | -p    | nautes pro delete project-name -p product-name |
| nautes coderepo get                  | cr,crs          | coderepo               | name  | -p    | nautes cr get cr-name -p product-name          |
| nautes coderepo delete               | cr,crs          | coderepo               | name  | -p    | nautes cr delete cr-name -p product-name       |
| nautes coderepobinding get           | crb,crbs        | coderepobinding        | name  | -p    | nautes crb get crb-name -p product-name        |
| nautes coderepobinding delete        | crb,crbs        | coderepobinding        | name  | -p    | nautes crb delete crb-name -p product-name     |
| nautes deploymentruntime get         | dr,drs          | deploymentruntime      | name  | -p    | nautes dr get dr-name -p product-name          |
| nautes deploymentruntime delete      | dr,drs          | deploymentruntime      | name  | -p    | nautes dr delete dr-name -p product-name       |
| nautes projectpipelineruntime get    | ppr,pprs        | projectpipelineruntime | name  | -p    | nautes ppr get ppr-name -p product-name        |
| nautes projectpipelineruntime delete | ppr,pprs        | projectpipelineruntime | name  | -p    | nautes ppr delete ppr-name -p product-name     |


CLI 的具体的使用方法请参见[用户手册](https://nautes.io/guide/user-guide/deploy-an-application.html)

## 快速开始

### 准备

安装 [go](https://golang.org/dl/)

### 构建

```
go build -o nautes
```

### 运行

```bash
nautes apply -f $FILE -t $TOKEN -s $API-SERVER
```