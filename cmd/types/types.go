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

package types

import (
	"reflect"
)

const _PRODUCT_PATH_TEMPLATE = "/api/v1/products/%s"
const _ENV_PATH_TEMPLATE = "/api/v1/products/%s/environments/%s"
const _PROJECT_PATH_TEMPLATE = "/api/v1/products/%s/projects/%s"
const _CODEREPO_PATH_TEMPLATE = "/api/v1/products/%s/coderepos/%s"
const _CODEREPO_BINDING_PATH_TEMPLATE = "/api/v1/products/%s/coderepobindings/%s"
const _DEPLOYMENTRUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/deploymentruntimes/%s"
const _PROJECTPIPELINERUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/projectpipelineruntimes/%s"
const _CLUSTER_PATH_TEMPLATE = "/api/v1/clusters/%s"

const (
	ResourceKind = "Kind"
	ApplyOrder   = "applyOrder"
	RemoveOrder  = "removeOrder"
	Column       = "column"
	MergeTo      = "mergeTo"
)

type ResourceFunc func(apiServer string, token string, skipCheck bool, resource string, resourceHandler ResourceHandler) error

type ResourcesType struct {
	ResourceType     reflect.Type
	ResponseItemType reflect.Type
}

type ResourceHandler interface {
	GetKind() string
	GetPathTemplate() string
	GetPathVarNames() []string
}

type Base struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
}

type Cluster struct {
	APIVersion string      `yaml:"apiVersion" json:"api_version"`
	Kind       string      `yaml:"kind" json:"kind" commands:"cls" applyOrder:"0" removeOrder:"7"`
	Spec       ClusterSpec `yaml:"spec" json:"spec"`
}

type ClusterSpec struct {
	Name           string         `yaml:"name" json:"name"`
	ApiServer      string         `yaml:"apiServer" json:"apiServer"`
	ClusterKind    string         `yaml:"clusterKind" json:"clusterKind"`
	ClusterType    string         `yaml:"clusterType" json:"clusterType"`
	Usage          string         `yaml:"usage" json:"usage"`
	WorkerType     string         `yaml:"workerType" json:"workerType"`
	HostCluster    string         `yaml:"hostCluster" json:"hostCluster"`
	PrimaryDomain  string         `yaml:"primaryDomain" json:"primaryDomain"`
	TektonHost     string         `yaml:"tektonHost" json:"tektonHost"`
	ArgoCDHost     string         `yaml:"argocdHost" json:"argocdHost"`
	Kubeconfig     string         `yaml:"kubeconfig" json:"kubeconfig"`
	Traefik        Traefik        `yaml:"traefik" json:"traefik"`
	VCluster       VCluster       `yaml:"vcluster" json:"vcluster"`
	ComponentsList ComponentsList `yaml:"componentsList" json:"componentsList"`
	// ReservedNamespacesAllowedProducts key is namespace name, value is the product name list witch can use namespace.
	ReservedNamespacesAllowedProducts map[string][]string `yaml:"reservedNamespacesAllowedProducts" json:"reservedNamespacesAllowedProducts"`
	// +optional
	// ReservedNamespacesAllowedProducts key is product name, value is the list of cluster resources.
	ProductAllowedClusterResources map[string][]ClusterResourceInfo `yaml:"productAllowedClusterResources" json:"productAllowedClusterResources"`
}

type ClusterResponse struct {
	Items []*ClusterResponseItem `json:"items"`
}

type ClusterResponseItem struct {
	Name           string                 `json:"name" column:"name"`
	ApiServer      string                 `json:"api_server" column:"ApiServer"`
	ClusterKind    string                 `json:"cluster_kind"`
	ClusterType    string                 `json:"cluster_type" column:"ClusterType"`
	Usage          string                 `json:"usage" column:"usage"`
	WorkerType     string                 `json:"worker_type" column:"WorkerType"`
	HostCluster    string                 `json:"host_cluster"`
	PrimaryDomain  string                 `json:"primary_domain" column:"PrimaryDomain"`
	TektonHost     string                 `json:"tekton_host"`
	ArgoCDHost     string                 `json:"argocd_host"`
	Kubeconfig     string                 `json:"kubeconfig"`
	Traefik        TraefikResponse        `json:"traefik"`
	VCluster       VClusterResponse       `json:"vcluster"`
	ComponentsList ComponentsListResponse `json:"components_list"`
	// ReservedNamespacesAllowedProducts key is namespace name, value is the product name list witch can use namespace.
	ReservedNamespacesAllowedProducts map[string][]string `json:"reserved_namespaces_allowed_products"`
	// +optional
	// ReservedNamespacesAllowedProducts key is product name, value is the list of cluster resources.
	ProductAllowedClusterResources map[string][]ClusterResourceInfo `json:"product_allowed_cluster_resources"`
}

type ClusterResourceInfo struct {
	Kind  string `yaml:"kind" json:"kind"`
	Group string `yaml:"group" json:"group"`
}

type Traefik struct {
	HTTPNodePort  string `yaml:"httpNodePort" json:"httpNodePort"`
	HTTPSNodePort string `yaml:"httpsNodePort" json:"httpsNodePort"`
}

type TraefikResponse struct {
	HTTPNodePort  string `yaml:"httpNodePort" json:"http_node_port"`
	HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
}

type VCluster struct {
	HTTPSNodePort string `yaml:"httpsNodePort" json:"httpsNodePort"`
}

type VClusterResponse struct {
	HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
}

// ComponentsList declares the specific components used by the cluster
type ComponentsList struct {
	// +optional
	CertMgt *Component `yaml:"certMgt" json:"certMgt"`
	// +optional
	Deployment *Component `yaml:"deployment" json:"deployment"`
	// +optional
	EventListener *Component `yaml:"eventListener" json:"eventListener"`
	// +optional
	IngressController *Component `yaml:"ingressController" json:"ingressController"`
	// +optional
	MultiTenant *Component `yaml:"multiTenant" json:"multiTenant"`
	// +optional
	Pipeline *Component `yaml:"pipeline" json:"pipeline"`
	// +optional
	ProgressiveDelivery *Component `yaml:"progressiveDelivery" json:"progressiveDelivery"`
	// +optional
	SecretMgt *Component `yaml:"secretMgt" json:"secretMgt"`
	// +optional
	SecretSync *Component `yaml:"secretSync" json:"secretSync"`
}

type ComponentsListResponse struct {
	// +optional
	CertMgt *Component `json:"cert_mgt"`
	// +optional
	Deployment *Component `json:"deployment"`
	// +optional
	EventListener *Component `json:"event_listener"`
	// +optional
	IngressController *Component `json:"ingress_controller"`
	// +optional
	MultiTenant *Component `json:"multi_tenant"`
	// +optional
	Pipeline *Component `json:"pipeline"`
	// +optional
	ProgressiveDelivery *Component `json:"progressive_delivery"`
	// +optional
	SecretMgt *Component `json:"secret_mgt"`
	// +optional
	SecretSync *Component `json:"secret_sync"`
}

type Component struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

func (c *Cluster) GetKind() string {
	return c.Kind
}

func (c *Cluster) GetPathTemplate() string {
	return _CLUSTER_PATH_TEMPLATE
}

func (c *Cluster) GetPathVarNames() []string {
	return []string{"Name"}
}

type Product struct {
	APIVersion string      `yaml:"apiVersion" json:"api_version"`
	Kind       string      `json:"kind" commands:"prod,prods" applyOrder:"1" removeOrder:"6"`
	Spec       ProductSpec `json:"spec"`
}

type ProductSpec struct {
	Name string `json:"name" column:"name"`
	Git  *struct {
		Gitlab *ProductGitRepo `json:"gitlab"`
		//GitHub *ProductGitRepo `json:"github"`
	} `json:"git"`
}

type ProductGitRepo struct {
	Name        string `json:"name"`
	Path        string `json:"path" column:"path"`
	Visibility  string `json:"visibility" column:"visibility"`
	Description string `json:"description" column:"description"`
	ParentID    int    `yaml:"parentID" json:"parent_id"`
}

type ProductResponse struct {
	Items []*ProductResponseItem `json:"items"`
}

type ProductResponseItem struct {
	ProductSpec
}

func (p *Product) GetKind() string {
	return p.Kind
}

func (p *Product) GetPathTemplate() string {
	return _PRODUCT_PATH_TEMPLATE
}

func (p *Product) GetPathVarNames() []string {
	return []string{"Name"}
}

type Environment struct {
	APIVersion string          `yaml:"apiVersion" json:"api_version"`
	Kind       string          `json:"kind" commands:"env,envs" applyOrder:"2" removeOrder:"5"`
	Spec       EnvironmentSpec `json:"spec"`
}

type EnvironmentSpec struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
	EnvType string `yaml:"envType"`
	Product string `json:"product"`
}

type EnvironmentResponse struct {
	Items []*EnvironmentResponseItem
}

type EnvironmentResponseItem struct {
	Name    string `json:"name" column:"name"`
	Product string `json:"product" column:"product"`
	Cluster string `json:"cluster" column:"cluster"`
	EnvType string `json:"env_type" column:"env_type"`
}

func (e *Environment) GetKind() string {
	return e.Kind
}

func (e *Environment) GetPathTemplate() string {
	return _ENV_PATH_TEMPLATE
}

func (e *Environment) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type Project struct {
	APIVersion string      `yaml:"apiVersion" json:"api_version"`
	Kind       string      `json:"kind" commands:"pro,proj,pros" applyOrder:"3" removeOrder:"4"`
	Spec       ProjectSpec `json:"spec"`
}

type ProjectSpec struct {
	Name     string `json:"name" column:"name"`
	Product  string `json:"product" column:"product"`
	Language string `json:"language" column:"language"`
}

type ProjectResponse struct {
	Items []*ProjectResponseItem `json:"items"`
}

type ProjectResponseItem struct {
	ProjectSpec
}

func (p *Project) GetKind() string {
	return p.Kind
}

func (p *Project) GetPathTemplate() string {
	return _PROJECT_PATH_TEMPLATE
}

func (p *Project) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepo struct {
	APIVersion string       `yaml:"apiVersion" json:"api_version"`
	Kind       string       `yaml:"kind" commands:"cr,crs" applyOrder:"4" removeOrder:"3"`
	Spec       CodeRepoSpec `yaml:"spec"`
}

type CodeRepoSpec struct {
	DeploymentRuntime      bool `yaml:"deploymentRuntime"`
	ProjectPipelineRuntime bool `yaml:"pipelineRuntime" json:"pipelineRuntime"`
	CodeRepoCommon
}

type CodeRepoResponse struct {
	Items []*CodeRepoResponseItem `json:"items"`
}

type CodeRepoResponseItem struct {
	CodeRepoCommon
	DeploymentRuntime      bool `json:"deployment_runtime"`
	ProjectPipelineRuntime bool `json:"pipeline_runtime"`
}

type CodeRepoCommon struct {
	Name    string `json:"name" column:"name"`
	Product string `json:"product" column:"product"`
	Project string `json:"project" column:"project" mergeTo:"product"`
	Git     *struct {
		Gitlab *CodeRepoGitRepoDetails `json:"gitlab"`
		//Github *CodeRepoGitRepoDetails `json:"github"`
	} `json:"git"`
	Webhook *struct {
		Events *[]string `json:"events"`
	} `json:"webhook"`
}

type CodeRepoGitRepoDetails struct {
	Name          string `json:"name"`
	Path          string `json:"path" column:"path"`
	Visibility    string `json:"visibility" column:"visibility" mergeTo:"path"`
	Description   string `json:"description"`
	SshUrlToRepo  string `json:"ssh_url_to_repo" column:"ssh_url_to_repo"`
	HttpUrlToRepo string `json:"http_url_to_repo" column:"http_url_to_repo" mergeTo:"ssh_url_to_repo"`
}

func (c *CodeRepo) GetKind() string {
	return c.Kind
}

func (c *CodeRepo) GetPathTemplate() string {
	return _CODEREPO_PATH_TEMPLATE
}

func (c *CodeRepo) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepoBinding struct {
	APIVersion string              `yaml:"apiVersion" json:"api_version"`
	Kind       string              `json:"kind" commands:"crb,crbs" applyOrder:"5" removeOrder:"2"`
	Spec       CodeRepoBindingSpec `json:"spec"`
}

type CodeRepoBindingSpec struct {
	ProductName string `yaml:"productName"`
	CodeRepoBindingCommon
}

type CodeRepoBindingResponse struct {
	Items []*CodeRepoBindingResponseItem `json:"items"`
}

type CodeRepoBindingResponseItem struct {
	CodeRepoBindingCommon
}

type CodeRepoBindingCommon struct {
	Name        string   `json:"name" column:"name"`
	Product     string   `json:"product" column:"product"`
	CodeRepo    string   `json:"coderepo" column:"coderepo"`
	Permissions string   `json:"permissions" column:"permissions"`
	Projects    []string `json:"projects" column:"projects"`
}

func (c *CodeRepoBinding) GetKind() string {
	return c.Kind
}

func (c *CodeRepoBinding) GetPathTemplate() string {
	return _CODEREPO_BINDING_PATH_TEMPLATE
}

func (c *CodeRepoBinding) GetPathVarNames() []string {
	return []string{"ProductName", "Name"}
}

type ProjectPipelineRuntime struct {
	APIVersion string                     `yaml:"apiVersion" json:"api_version"`
	Kind       string                     `yaml:"kind" commands:"ppr,pprs" applyOrder:"6" removeOrder:"1"`
	Spec       ProjectPipelineRuntimeSpec `yaml:"spec"`
}

type ProjectPipelineRuntimeSpec struct {
	ProjectPipelineRuntimeCommon
	Product          string                                    `yaml:"product"`
	PipelineSource   string                                    `yaml:"pipelineSource" json:"pipelineSource"`
	EventSources     *[]ProjectPipelineRuntimeSpecEventSources `yaml:"eventSources"`
	PipelineTriggers *[]struct {
		EventSource string `yaml:"eventSource"`
		Pipeline    string `yaml:"pipeline"`
		Revision    string `yaml:"revision"`
	} `yaml:"pipelineTriggers"`
}

type ProjectPipelineRuntimeCommon struct {
	Name        string `json:"name" column:"name"`
	Project     string `json:"project" column:"project"`
	Destination string `json:"destination" column:"destination"`
	Isolation   string `json:"isolation"`
	Pipelines   *[]struct {
		Name  string `json:"name"`
		Label string `json:"label"`
		Path  string `json:"path"`
	} `json:"pipelines"`
}

type ProjectPipelineRuntimeResponse struct {
	Items []*ProjectPipelineRuntimeResponseItem `json:"items"`
}

type ProjectPipelineRuntimeResponseItem struct {
	ProjectPipelineRuntimeCommon
	Product        string `json:"product"`
	PipelineSource string `json:"pipeline_source" column:"PipelineSource"`
	EventSources   *[]struct {
		Name   string `json:"name"`
		Gitlab *struct {
			RepoName string   `json:"repo_name"  column:"EventSources-RepoName"`
			Revision string   `json:"revision"`
			Events   []string `json:"events"`
		} `json:"gitlab"`
		Calendar *struct {
			Schedule       string   `json:"schedule"`
			Interval       string   `json:"interval"`
			ExclusionDates []string `json:"exclusion_dates"`
			Timezone       string   `json:"timezone"`
		} `json:"calendar"`
	} `json:"event_sources"`
	PipelineTriggers *[]struct {
		EventSource string `json:"event_source" column:"Triggers-EventSource" mergeTo:"EventSources-RepoName"`
		Pipeline    string `json:"pipeline"`
		Revision    string `json:"revision"`
	} `json:"pipeline_triggers"`
}

type ProjectPipelineRuntimeSpecEventSources struct {
	Name   string `yaml:"name"`
	Gitlab *struct {
		RepoName string   `yaml:"repoName"`
		Revision string   `yaml:"revision"`
		Events   []string `yaml:"events"`
	} `json:"gitlab"`
	Calendar *struct {
		Schedule       string   `yaml:"schedule"`
		Interval       string   `yaml:"interval"`
		ExclusionDates []string `yaml:"exclusionDates"`
		Timezone       string   `yaml:"timezone"`
	} `yaml:"calendar"`
}

func (p *ProjectPipelineRuntime) GetKind() string {
	return p.Kind
}

func (p *ProjectPipelineRuntime) GetPathTemplate() string {
	return _PROJECTPIPELINERUNTIME_PATH_TEMPLATE
}

func (p *ProjectPipelineRuntime) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type DeploymentRuntime struct {
	APIVersion string                `yaml:"apiVersion" json:"api_version"`
	Kind       string                `json:"kind" commands:"dr,drs" applyOrder:"7" removeOrder:"0"`
	Spec       DeploymentRuntimeSpec `yaml:"spec" json:"spec"`
}

type DeploymentRuntimeSpec struct {
	Name           string `json:"name" column:"name"`
	Product        string `json:"product" column:"product"`
	ManifestSource struct {
		CodeRepo       string `yaml:"codeRepo"`
		TargetRevision string `yaml:"targetRevision"`
		Path           string `yaml:"path"`
	} `yaml:"manifestsource"`
	ProjectsRef []string `yaml:"projectsRef"`
	Destination *struct {
		Environment string   `yaml:"environment"`
		Namespaces  []string `yaml:"namespaces"`
	} `yaml:"destination" json:"destination"`
}

type DeploymentRuntimeResponse struct {
	Items []*DeploymentRuntimeResponseItem `json:"items"`
}

type DeploymentRuntimeResponseItem struct {
	Name           string `json:"name" column:"name"`
	Product        string `json:"product" column:"product"`
	ManifestSource *struct {
		CodeRepo       string `yaml:"codeRepo" json:"code_repo" column:"codeRepo"`
		TargetRevision string `yaml:"targetRevision" json:"target_revision" column:"targetRevision" mergeTo:"codeRepo"`
		Path           string `yaml:"path" json:"path" column:"path" mergeTo:"codeRepo"`
	} `yaml:"manifestsource" json:"manifest_source"`
	ProjectsRef []string `yaml:"projectsRef" json:"projects_ref" column:"projectsRef"`
	Destination *struct {
		Environment string   `yaml:"environment" json:"environment" column:"environment"`
		Namespaces  []string `yaml:"namespaces" json:"namespaces" column:"namespaces" mergeTo:"environment"`
	} `yaml:"destination" json:"destination"`
}

func (d *DeploymentRuntime) GetKind() string {
	return d.Kind
}

func (d *DeploymentRuntime) GetPathTemplate() string {
	return _DEPLOYMENTRUNTIME_PATH_TEMPLATE
}

func (d *DeploymentRuntime) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

// ClientOptions hold api server address, token, and other settings for the API client.
type ClientOptions struct {
	ServerAddr string
	Token      string
	SkipCheck  bool
}
