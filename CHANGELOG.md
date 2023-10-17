# Change Log

## v0.4.1

> Change log since v0.4.0

### Upgrade Notice

> No, really, you must read this before you upgrade.

### Fixed
1. Fixed path error of examples/demo-pipeline.yaml.

### Changes
1. Renamed the attribute of additions of the multi-tenant such as ProductResourcePathPipeline, ProductResourceRevision, and SyncResourceTypes.
```yaml
  componentsList:
     multiTenant:
        name: hnc
        namespace: hnc-system
        additions:
           productResourceKustomizeFileFolder: templates/pipelines
           productResourceRevision: main
           syncResourceTypes: tekton.dev/Pipeline
```

### New Feature
1. By default, the runtime creates an account with the name of the runtime. You can also specify an account or not.
   What does mean the account for runtime which is the deployment runtime, or project pipeline runtime? It's a ServiceAccount for Kubernetes and a Role for Vault.
- DeploymentRuntime
```yaml
apiVersion: nautes.resource.nautes.io/v1alpha1
kind: DeploymentRuntime
spec:
  name: dr-demo
  account: dr-demo-account
```
- ProjectPipelineRuntime
```yaml
apiVersion: nautes.resource.nautes.io/v1alpha1
kind: ProjectPipelineRuntime
spec:
  name: pr-demo
  account: pr-demo-account
```

### How to use
```shell
nautes get ppr -p $PRODUCT-NAME -t $TOKEN -s $API-SERVER
```

## v0.4.0

> Change log since v0.3.9

### Upgrade Notice

> No, really, you must read this before you upgrade.

### Changes
1. Deleted the argocdHost tektonHost and traefik attributes when adding a cluster.

### New Feature
1. Supported custom components for cluster resource. It would be best to use the componentsList attribute when adding a cluster.
   The componentsList includes three properties which are name and namespace and additions which are additional properties. The Key is the component attribute, the Value is value of the component attribute.
   eg: If the traefik as gateway, it can be set attributes of traefik by the additions attribute.
```yaml
  componentsList:
    gateway:
      name: traefik
      namespace: traefik
      additions:
        httpNodePort: "30080"
        httpsNodePort: "30443"
```

eg: At least be used multiTenant component and gateway of the cluster when adding a pipeline runtime cluster.
```yaml
  componentsList:
    multiTenant:
      name: hnc
      namespace: hnc-system
      additions:
        ProductResourcePathPipeline: templates/pipelines
        ProductResourceRevision: main
        SyncResourceTypes: tekton.dev/Pipeline
    gateway:
      name: traefik
      namespace: traefik
      additions:
        httpNodePort: "30080"
        httpsNodePort: "30443"
```

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