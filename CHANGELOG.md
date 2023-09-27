# Change Log

## v0.4.0

> Change log since v0.3.9

### Upgrade Notice

> No, really, you must read this before you upgrade.

- Support custom components for cluster resource.

### How to use
```shell
nautes get ppr -p $PRODUCT-NAME -t $TOKEN -s $API-SERVER
```

## v0.3.9

> Change log since v0.3.8

### Upgrade Notice

> No, really, you must read this before you upgrade.

- When you use the apply command to create a ProjectPipelineRuntime resource by yaml file, you must change the Destination attribute of ProjectPipelineRuntime resource to object which the result includes Environment and Namespace. The type of Namespace is a string. For example:
```yaml
  destination:
    environment: env-test-demo
    namespaces: dr-demo
```

### How to use
```shell
nautes get ppr -p $PRODUCT-NAME -t $TOKEN -s $API-SERVER
```

### Changes
1. Changed the destination attribute of the ProjectPipelineRuntime resource.

### New Feature
1. Added optional additional resources for ProjectPipelineRuntime.
   ([#3](https://github.com/nautes-labs/cli/pull/3), [@rubinus](https://github.com/rubinus))

## v0.3.8

> Change log since v0.3.0

### Upgrade Notice

> No, really, you must read this before you upgrade.

- When you use the apply command to create a DeploymentRuntime resource by yaml file, you must change the Destination attribute of the DeploymentRuntime resource to an object the result includes Environment and Namespaces. The type of Namespaces is an array. For example:
```yaml
  destination:
    environment: env-test-demo
    namespaces:
      - dr-demo
```

### How to use
```shell
nautes get cr -p $PRODUCT-NAME -t $TOKEN -s $API-SERVER
```
**There is a simple method to use that is set environments for example.**

- export API_SERVER=http://127.0.0.1:8000

- export GIT_TOKEN=glpat-yYTvmC9Vnzom5k2NuzUU

- export PRODUCT=demo-101

```shell
nautes get cr coderepo-name-101 coderepo-name-102 coderepo-name-103 -o json

nautes get dr

nautes delete pro project-101 project-102
```

### Changes
1. Changed the destination attribute of DeploymentRuntime resource.

### New Feature
1. Added get and delete commands to CLI.
   ([#2](https://github.com/nautes-labs/cli/pull/2), [@rubinus](https://github.com/rubinus))

2. Added cluster attributes such as components、reservedNamespaces and clusterResources

## v0.3.0

### How to use

```shell
nautes apply -f $FILE -t $TOKEN -s $API-SERVER
```

### Changes

The following APIs have been added:

1. CodeRepoBinding: save_coderepobinding, delete_coderepobinding.
2. ProjectPipelineRuntime: save_projectpipelineruntime, delete_projectpipelineruntime.

The input parameters of the following APIs have been upgraded:

1. Cluster: save_cluster.
2. CodeRepo: save_coderepo.

## v0.2.0

### How to use
```shell
nautes apply -f $FILE -t $TOKEN -s $API-SERVER
```
### Changes
This is the first release, list of supported APIs:
1. Cluster: save_cluster, delete_cluster.
2. Product: save_product, delete_product.
3. Environment: save_environment, delete_environment.
4. Project: save_project, delete_project.
5. CodeRepo: save_coderepo, delete_coderepo.
6. DeploymentRuntime: save_deploymentruntime, delete_deploymentruntime.