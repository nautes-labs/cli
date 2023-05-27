# Nautes CLI

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![golang](https://img.shields.io/badge/golang-v1.19-brightgreen)](https://go.dev/doc/install)
[![version](https://img.shields.io/badge/version-v0.3.0-green)]()

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