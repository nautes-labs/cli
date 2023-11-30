# Nautes CLI

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![golang](https://img.shields.io/badge/golang-v1.20-brightgreen)](https://go.dev/doc/install)
[![version](https://img.shields.io/badge/version-v0.4.2-green)]()

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

## CLI 扩展子命令

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
| nautes get product                   | prod,prods      | product                | name  |       | nautes get prod product-name                   |
| nautes delete product                | prod,prods      | product                | name  |       | nautes delete prod product-name                |
| nautes get cluster                   | cls             | cluster                | name  |       | nautes get cls cluster-name                    |
| nautes delete cluster                | cls             | cluster                | name  |       | nautes delete cls cluster-name                 |
| nautes get environment               | env,envs        | environment            | name  | -p    | nautes get env env-name -p product-name        |
| nautes delete environment            | env,envs        | environment            | name  | -p    | nautes delete env env-name -p product-name     |
| nautes get project                   | pro,pros,proj   | project                | name  | -p    | nautes get pro project-name -p product-name    |
| nautes delete project                | pro,pros,proj   | project                | name  | -p    | nautes delete pro project-name -p product-name |
| nautes get coderepo                  | cr,crs          | coderepo               | name  | -p    | nautes get cr cr-name -p product-name          |
| nautes delete coderepo               | cr,crs          | coderepo               | name  | -p    | nautes delete cr cr-name -p product-name       |
| nautes get coderepobinding           | crb,crbs        | coderepobinding        | name  | -p    | nautes get crb crb-name -p product-name        |
| nautes delete coderepobinding        | crb,crbs        | coderepobinding        | name  | -p    | nautes delete crb crb-name -p product-name     |
| nautes get deploymentruntime         | dr,drs          | deploymentruntime      | name  | -p    | nautes get dr dr-name -p product-name          |
| nautes delete deploymentruntime      | dr,drs          | deploymentruntime      | name  | -p    | nautes delete dr dr-name -p product-name       |
| nautes get projectpipelineruntime    | ppr,pprs        | projectpipelineruntime | name  | -p    | nautes get ppr ppr-name -p product-name        |
| nautes delete projectpipelineruntime | ppr,pprs        | projectpipelineruntime | name  | -p    | nautes delete ppr ppr-name -p product-name     |


CLI 的具体的使用方法请参见[用户手册](https://nautes.io/guide/user-guide/deploy-an-application.html)

## 快速添加资源

### 添加资源

举例，比如添加制品库资源，在 cmd/types/types.go 中添加制品库资源的定义，注意标签 YAML 是小驼峰，JSON 是下划线:

```yaml
type ArtifactRepo struct {
  APIVersion string              `yaml:"apiVersion" json:"api_version"`
  Kind       string              `yaml:"kind" json:"kind" commands:"ar,ars" applyOrder:"8" removeOrder:"0"`
  Spec       ArtifactRepoResponseItem `yaml:"spec" json:"spec"`
}

type ArtifactRepoResponse struct {
  Items []*ArtifactRepoResponseItem `json:"items"`
}

type ArtifactRepoResponseItem struct {
  ArtifactRepoProvider string `json:"artifact_repo_provider" yaml:"artifactRepoProvider"`
  Product              string `json:"product" yaml:"product"  column:"product"`
  Projects []string  `json:"projects" yaml:"projects"`
  RepoName string    `json:"repo_name" yaml:"repoName" column:"repoName"`
  RepoType string    `json:"repo_type" yaml:"repoType"  column:"repoType" mergeTo:"repoName"`
  PackageType string `json:"package_type" yaml:"packageType"`
}
```

### 新资源的请求路径

需要在 cmd/types/types.go 定义常量，声明请求 api-server 的路径，在执行创建或删除时会替换路径中的参数

```go
const ArtifactRepoPathTemplate = "/api/v1/products/%s/artifactrepos/%s"

```

### 新加资源要实现接口中的三个方法

需要实现 cmd/types/types.go 中的 ResourceHandler 接口

- GetKind() 返回资源的类型，即资源中的 Kind 属性.
- GetPathTemplate() 返回请求 api-server 的接口路径.
- GetPathVarNames() 返回接口路径中需要填充的参数.

### 设置扩展标签

在 Nautes 的资源中，Cluster 和 Product 属于一等公民，而 Project, Environment, CodeRepo, CodeRepoBinding, ProjectPipelineRuntime, DeploymentRuntime 属于二等公民。

> 二等公民需要依赖一等公民的创建，比如：创建 Project 时需要先创建 Product； 

> 二等公民之间也有依赖关系，比如 创建 CodeRepo 时需要依赖 Project，创建 DeploymentRuntime 时需要依赖 Project, Environment, CodeRepo。

当在一个 yaml 文件中定义了所有资源且无序时，cli 客户端在执行的时候会按照内部自动排序，在添加资源时按升序添加，而在删除资源时按降序删除，避免因资源的依赖关系产生错误。

- commands: 出现在 Kind 属性中，用标签 key:value 的形式表示客户端简写的命令，命令可以有多个，比如示例中的：commands:"ar,ars"，在 cli 执行时可以使用如：nautes get ar 这样的命令。
- applyOrder: 添加资源时的顺序，按升序添加，数字越小，优先级越高，优先创建
- removeOrder: 删除资源时的顺序，按降序删除，数字越大，优先级越高，优先删除
- column: 要打印显示的列，用标签  key:value 的形式表示要显示的列
- mergeTo: 如果一行要显示多列，可以用合并列的方式，把一列添加到目标列上来显示

### 把新加的资源类型添加到 cmd/main.go 的 resourcesTypeArr 数组中

- 定义资源类型
- 定义资源响应返回值的类型
```struct
{
  ResourceType:     reflect.TypeOf(types.ArtifactRepo{}),
  ResponseItemType: reflect.TypeOf(types.ArtifactRepoResponseItem{}),
}
```

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