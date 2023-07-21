// Copyright 2023 Nautes Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

const _PRODUCT_PATH_TEMPLATE = "/api/v1/products/%s"
const _ENV_PATH_TEMPLATE = "/api/v1/products/%s/environments/%s"
const _PROJECT_PATH_TEMPLATE = "/api/v1/products/%s/projects/%s"
const _CODEREPO_PATH_TEMPLATE = "/api/v1/products/%s/coderepos/%s"
const _CODEREPO_BINDING_PATH_TEMPLATE = "/api/v1/products/%s/coderepobindings/%s"
const _DEPLOYMENTRUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/deploymentruntimes/%s"
const _PROJECTPIPELINERUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/projectpipelineruntimes/%s"
const _CLUSTER_PATH_TEMPLATE = "/api/v1/clusters/%s"

type resourceFunc func(apiServer string, token string, skipCheck bool, resource string, resourceHandler ResourceHandler) error

type ResourceHandler interface {
	getKind() string
	getPathTemplate() string
	getPathVarNames() []string
}

type Base struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
}

type Product struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name string `json:"name"`
		Git  *struct {
			Gitlab *struct {
				Name        string `json:"name"`
				Path        string `json:"path"`
				Visibility  string `json:"visibility"`
				Description string `json:"description"`
				ParentID    int    `yaml:"parentID" json:"parent_id"`
			} `json:"gitlab"`
		} `json:"git"`
	} `json:"spec"`
}

func (p Product) getKind() string {
	return p.Kind
}

func (p Product) getPathTemplate() string {
	return _PRODUCT_PATH_TEMPLATE
}

func (p Product) getPathVarNames() []string {
	return []string{"Name"}
}

type Environment struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name    string `json:"name"`
		Cluster string `json:"cluster"`
		EnvType string `yaml:"envType" json:"env_type"`
		Product string `json:"product"`
	} `json:"spec"`
}

func (e Environment) getKind() string {
	return e.Kind
}

func (e Environment) getPathTemplate() string {
	return _ENV_PATH_TEMPLATE
}

func (e Environment) getPathVarNames() []string {
	return []string{"Product", "Name"}
}

type Project struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name     string `json:"name"`
		Language string `json:"language"`
		Product  string `json:"product"`
	} `json:"spec"`
}

func (p Project) getKind() string {
	return p.Kind
}

func (p Project) getPathTemplate() string {
	return _PROJECT_PATH_TEMPLATE
}

func (p Project) getPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepo struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Product                string `json:"product"`
		Name                   string `json:"name"`
		Project                string `json:"project"`
		DeploymentRuntime      bool   `yaml:"deploymentRuntime" json:"deployment_runtime"`
		ProjectPipelineRuntime bool   `yaml:"pipelineRuntime" json:"pipeline_runtime"`
		Webhook                *struct {
			Events *[]string `json:"events"`
		} `json:"webhook"`
		Git *struct {
			Gitlab *struct {
				Name        string `json:"name"`
				Path        string `json:"path"`
				Visibility  string `json:"visibility"`
				Description string `json:"description"`
			} `json:"gitlab"`
		} `json:"git"`
	} `json:"spec"`
}

func (c CodeRepo) getKind() string {
	return c.Kind
}

func (c CodeRepo) getPathTemplate() string {
	return _CODEREPO_PATH_TEMPLATE
}

func (c CodeRepo) getPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepoBinding struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		ProductName string    `yaml:"productName" json:"product_name"`
		Name        string    `json:"name"`
		CodeRepo    string    `json:"coderepo"`
		Product     string    `json:"product"`
		Projects    *[]string `json:"projects"`
		Permissions string    `json:"permissions"`
	} `json:"spec"`
}

func (c CodeRepoBinding) getKind() string {
	return c.Kind
}

func (c CodeRepoBinding) getPathTemplate() string {
	return _CODEREPO_BINDING_PATH_TEMPLATE
}

func (c CodeRepoBinding) getPathVarNames() []string {
	return []string{"ProductName", "Name"}
}

type DeploymentRuntime struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name        string    `json:"name"`
		Product     string    `json:"product"`
		ProjectsRef *[]string `yaml:"projectsRef" json:"projects_ref"`
		Destination struct {
			Environment string   `yaml:"environment" json:"environment"`
			Namespaces  []string `yaml:"namespaces" json:"namespaces"`
		} `yaml:"destination" json:"destination"`
		Manifestsource *struct {
			CodeRepo       string `yaml:"codeRepo" json:"code_repo"`
			TargetRevision string `yaml:"targetRevision" json:"target_revision"`
			Path           string `json:"path"`
		} `yaml:"manifestsource" json:"manifest_source"`
	} `json:"spec"`
}

func (d DeploymentRuntime) getKind() string {
	return d.Kind
}

func (d DeploymentRuntime) getPathTemplate() string {
	return _DEPLOYMENTRUNTIME_PATH_TEMPLATE
}

func (d DeploymentRuntime) getPathVarNames() []string {
	return []string{"Product", "Name"}
}

type ProjectPipelineRuntime struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name           string `json:"name"`
		Product        string `json:"product"`
		Project        string `json:"project"`
		PipelineSource string `yaml:"pipelineSource" json:"pipeline_source"`
		Pipelines      *[]struct {
			Name  string `json:"name"`
			Label string `json:"label"`
			Path  string `json:"path"`
		} `json:"pipelines"`
		Destination  string `yaml:"destination" json:"destination"`
		EventSources *[]struct {
			Name   string `json:"name"`
			Gitlab *struct {
				RepoName string   `yaml:"repoName" json:"repo_name"`
				Revision string   `json:"revision"`
				Events   []string `json:"events"`
			} `json:"gitlab"`
			Calendar *struct {
				Schedule       string   `json:"schedule"`
				Interval       string   `json:"interval"`
				ExclusionDates []string `yaml:"exclusionDates" json:"exclusion_dates"`
				Timezone       string   `json:"timezone"`
			} `json:"calendar"`
		} `yaml:"eventSources" json:"event_sources"`
		Isolation        string `json:"isolation"`
		PipelineTriggers *[]struct {
			EventSource string `yaml:"eventSource" json:"event_source"`
			Pipeline    string `json:"pipeline"`
			Revision    string `json:"revision"`
		} `yaml:"pipelineTriggers" json:"pipeline_triggers"`
	}
}

func (p ProjectPipelineRuntime) getKind() string {
	return p.Kind
}

func (p ProjectPipelineRuntime) getPathTemplate() string {
	return _PROJECTPIPELINERUNTIME_PATH_TEMPLATE
}

func (p ProjectPipelineRuntime) getPathVarNames() []string {
	return []string{"Product", "Name"}
}

type Cluster struct {
	APIVersion string      `yaml:"apiVersion" json:"api_version"`
	Kind       string      `yaml:"kind" json:"kind"`
	Spec       ClusterSpec `yaml:"spec" json:"spec"`
}

type ClusterSpec struct {
	Name           string         `yaml:"name" json:"name"`
	ApiServer      string         `yaml:"apiServer" json:"api_server"`
	ClusterKind    string         `yaml:"clusterKind" json:"cluster_kind"`
	ClusterType    string         `yaml:"clusterType" json:"cluster_type"`
	Usage          string         `yaml:"usage" json:"usage"`
	WorkerType     string         `yaml:"workerType" json:"worker_type"`
	HostCluster    string         `yaml:"hostCluster" json:"host_cluster"`
	PrimaryDomain  string         `yaml:"primaryDomain" json:"primary_domain"`
	TektonHost     string         `yaml:"tektonHost" json:"tekton_host"`
	ArgoCDHost     string         `yaml:"argocdHost" json:"argocd_host"`
	Kubeconfig     string         `yaml:"kubeconfig" json:"kubeconfig"`
	Traefik        Traefik        `yaml:"traefik" json:"traefik"`
	VCluster       VCluster       `yaml:"vcluster" json:"vcluster"`
	ComponentsList ComponentsList `yaml:"componentsList" json:"components_list"`
	// ReservedNamespacesAllowedProducts key is namespace name, value is the product name list witch can use namespace.
	ReservedNamespacesAllowedProducts map[string][]string `yaml:"reservedNamespacesAllowedProducts" json:"reserved_namespaces_allowed_products"`
	// +optional
	// ReservedNamespacesAllowedProducts key is product name, value is the list of cluster resources.
	ProductAllowedClusterResources map[string][]ClusterResourceInfo `yaml:"productAllowedClusterResources" json:"product_allowed_cluster_resources"`
}

type ClusterResourceInfo struct {
	Kind  string `yaml:"kind" json:"kind"`
	Group string `yaml:"group" json:"group"`
}

type Traefik struct {
	HTTPNodePort  string `yaml:"httpNodePort" json:"http_node_port"`
	HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
}

type VCluster struct {
	HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
}

// ComponentsList declares the specific components used by the cluster
type ComponentsList struct {
	// +optional
	CertMgt *Component `yaml:"certMgt" json:"cert_mgt"`
	// +optional
	Deployment *Component `yaml:"deployment" json:"deployment"`
	// +optional
	EventListener *Component `yaml:"eventListener" json:"event_listener"`
	// +optional
	IngressController *Component `yaml:"ingressController" json:"ingress_controller"`
	// +optional
	MultiTenant *Component `yaml:"multiTenant" json:"multi_tenant"`
	// +optional
	Pipeline *Component `yaml:"pipeline" json:"pipeline"`
	// +optional
	ProgressiveDelivery *Component `yaml:"progressiveDelivery" json:"progressive_delivery"`
	// +optional
	SecretMgt *Component `yaml:"secretMgt" json:"secret_mgt"`
	// +optional
	SecretSync *Component `yaml:"secretSync" json:"secret_sync"`
}

type Component struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

func (c Cluster) getKind() string {
	return c.Kind
}

func (c Cluster) getPathTemplate() string {
	return _CLUSTER_PATH_TEMPLATE
}

func (c Cluster) getPathVarNames() []string {
	return []string{"Name"}
}
