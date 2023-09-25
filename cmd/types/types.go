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

const ProductPathTemplate = "/api/v1/products/%s"
const EnvPathTemplate = "/api/v1/products/%s/environments/%s"
const ProjectPathTemplate = "/api/v1/products/%s/projects/%s"
const CodeRepoPathTemplate = "/api/v1/products/%s/coderepos/%s"
const CodeRepoBindingPathTemplate = "/api/v1/products/%s/coderepobindings/%s"
const DeploymentRuntimePathTemplate = "/api/v1/products/%s/deploymentruntimes/%s"
const ProjectPipelineRuntimePathTemplate = "/api/v1/products/%s/projectpipelineruntimes/%s"
const ClusterPathTemplate = "/api/v1/clusters/%s"

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
	APIVersion string              `yaml:"apiVersion" json:"api_version"`
	Kind       string              `yaml:"kind" json:"kind" commands:"cls" applyOrder:"0" removeOrder:"7"`
	Spec       ClusterResponseItem `yaml:"spec" json:"spec"`
}

type ClusterResponse struct {
	Items []*ClusterResponseItem `json:"items"`
}

type ClusterResponseItem struct {
	Name          string   `yaml:"name" json:"name" column:"name"`
	ApiServer     string   `yaml:"apiServer" json:"api_server" column:"ApiServer"`
	ClusterKind   string   `yaml:"clusterKind" json:"cluster_kind"`
	Usage         string   `yaml:"usage" json:"usage" column:"Usage" mergeTo:"ApiServer"`
	ClusterType   string   `yaml:"clusterType" json:"cluster_type" column:"CT"  mergeTo:"ApiServer"`
	WorkerType    string   `yaml:"workerType" json:"worker_type" column:"WT" mergeTo:"ApiServer"`
	HostCluster   string   `yaml:"hostCluster" json:"host_cluster"`
	PrimaryDomain string   `yaml:"primaryDomain" json:"primary_domain" column:"PrimaryDomain"`
	Kubeconfig    string   `yaml:"kubeconfig" json:"kubeconfig"`
	VCluster      VCluster `yaml:"vcluster" json:"vcluster"`
	// ReservedNamespacesAllowedProducts key is namespace name, value is the product name list witch can use namespace.
	ReservedNamespacesAllowedProducts map[string][]string `yaml:"reservedNamespacesAllowedProducts" json:"reserved_namespaces_allowed_products"`
	// +optional
	// ReservedNamespacesAllowedProducts key is product name, value is the list of cluster resources.
	ProductAllowedClusterResources map[string][]ClusterResourceInfo `yaml:"productAllowedClusterResources" json:"product_allowed_cluster_resources"`
	ComponentsList                 ComponentsList                   `yaml:"componentsList" json:"components_list" column:"ComponentsList:Name"`
}

type ClusterResourceInfo struct {
	Kind  string `yaml:"kind" json:"kind"`
	Group string `yaml:"group" json:"group"`
}

type VCluster struct {
	HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
}

// ComponentsList declares the specific components used by the cluster
type ComponentsList struct {
	// +optional
	Deployment *Component `yaml:"deployment" json:"deployment"`
	// +optional
	EventListener *Component `yaml:"eventListener" json:"event_listener"`
	// +optional
	Gateway *Component `yaml:"gateway" json:"gateway"`
	// +optional
	MultiTenant *Component `yaml:"multiTenant" json:"multi_tenant"`
	// +optional
	Pipeline *Component `yaml:"pipeline" json:"pipeline"`
	// +optional
	ProgressiveDelivery *Component `yaml:"progressiveDelivery" json:"progressive_delivery"`
	// +optional
	SecretSync *Component `yaml:"secretSync" json:"secret_sync"`
}

type Component struct {
	Name      string            `yaml:"name" json:"name"`
	Namespace string            `yaml:"namespace" json:"namespace"`
	Additions map[string]string `yaml:"additions" json:"additions"`
}

func (c *Cluster) GetKind() string {
	return c.Kind
}

func (c *Cluster) GetPathTemplate() string {
	return ClusterPathTemplate
}

func (c *Cluster) GetPathVarNames() []string {
	return []string{"Name"}
}

type Product struct {
	APIVersion string              `yaml:"apiVersion" json:"api_version"`
	Kind       string              `yaml:"kind" json:"kind" commands:"prod,prods" applyOrder:"1" removeOrder:"6"`
	Spec       ProductResponseItem `yaml:"spec" json:"spec"`
}

type ProductResponse struct {
	Items []*ProductResponseItem `yaml:"items" json:"items"`
}

type ProductResponseItem struct {
	Name string          `yaml:"name" json:"name" column:"name"`
	Git  *ProductSpecGit `yaml:"git" json:"git"`
}

type ProductSpecGit struct {
	Gitlab *ProductGitRepo `yaml:"gitlab" json:"gitlab"`
}

type ProductGitRepo struct {
	Name        string `yaml:"name" json:"name"`
	Path        string `yaml:"path" json:"path" column:"path"`
	Visibility  string `yaml:"visibility" json:"visibility" column:"visibility"`
	Description string `yaml:"description" json:"description" column:"description"`
	ParentID    int    `yaml:"parentID" json:"parent_id"`
}

func (p *Product) GetKind() string {
	return p.Kind
}

func (p *Product) GetPathTemplate() string {
	return ProductPathTemplate
}

func (p *Product) GetPathVarNames() []string {
	return []string{"Name"}
}

type Environment struct {
	APIVersion string                  `yaml:"apiVersion" json:"api_version"`
	Kind       string                  `yaml:"kind" json:"kind" commands:"env,envs" applyOrder:"2" removeOrder:"5"`
	Spec       EnvironmentResponseItem `yaml:"spec" json:"spec"`
}

type EnvironmentResponse struct {
	Items []*EnvironmentResponseItem
}

type EnvironmentResponseItem struct {
	Name    string `yaml:"name" json:"name" column:"name"`
	Product string `yaml:"product" json:"product" column:"product"`
	Cluster string `yaml:"cluster" json:"cluster" column:"cluster"`
	EnvType string `yaml:"envType" json:"env_type" column:"env_type"`
}

func (e *Environment) GetKind() string {
	return e.Kind
}

func (e *Environment) GetPathTemplate() string {
	return EnvPathTemplate
}

func (e *Environment) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type Project struct {
	APIVersion string              `yaml:"apiVersion" json:"api_version"`
	Kind       string              `yaml:"kind" json:"kind" commands:"pro,proj,pros" applyOrder:"3" removeOrder:"4"`
	Spec       ProjectResponseItem `yaml:"spec" json:"spec"`
}

type ProjectResponse struct {
	Items []*ProjectResponseItem `yaml:"items" json:"items"`
}

type ProjectResponseItem struct {
	Name     string `yaml:"name" json:"name" column:"name"`
	Product  string `yaml:"product" json:"product" column:"product"`
	Language string `yaml:"language" json:"language" column:"language"`
}

func (p *Project) GetKind() string {
	return p.Kind
}

func (p *Project) GetPathTemplate() string {
	return ProjectPathTemplate
}

func (p *Project) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepo struct {
	APIVersion string               `yaml:"apiVersion" json:"api_version"`
	Kind       string               `yaml:"kind" json:"kind" commands:"cr,crs" applyOrder:"4" removeOrder:"3"`
	Spec       CodeRepoResponseItem `yaml:"spec" json:"spec"`
}

type CodeRepoResponse struct {
	Items []*CodeRepoResponseItem `yaml:"items" json:"items"`
}

type CodeRepoResponseItem struct {
	Name                   string                       `yaml:"name" json:"name" column:"name"`
	Product                string                       `yaml:"product" json:"product" column:"product"`
	Project                string                       `yaml:"project" json:"project" column:"project" mergeTo:"product"`
	Git                    *CodeRepoResponseItemGit     `yaml:"git" json:"git"`
	Webhook                *CodeRepoResponseItemWebhook `yaml:"webhook" json:"webhook"`
	DeploymentRuntime      bool                         `yaml:"deploymentRuntime" json:"deployment_runtime"`
	ProjectPipelineRuntime bool                         `yaml:"projectPipelineRuntime" json:"pipeline_runtime"`
}

type CodeRepoResponseItemGit struct {
	Gitlab *CodeRepoGitRepoDetails `yaml:"gitlab" json:"gitlab"`
}

type CodeRepoResponseItemWebhook struct {
	Events *[]string `yaml:"events" json:"events"`
}

type CodeRepoGitRepoDetails struct {
	Name          string `yaml:"name" json:"name"`
	Path          string `yaml:"path" json:"path" column:"path"`
	Visibility    string `yaml:"visibility" json:"visibility" column:"visibility" mergeTo:"path"`
	Description   string `yaml:"description" json:"description"`
	SshUrlToRepo  string `yaml:"sshUrlToRepo" json:"ssh_url_to_repo" column:"ssh_url_to_repo"`
	HttpUrlToRepo string `yaml:"httpUrlToRepo" json:"http_url_to_repo" column:"http_url_to_repo" mergeTo:"ssh_url_to_repo"`
}

func (c *CodeRepo) GetKind() string {
	return c.Kind
}

func (c *CodeRepo) GetPathTemplate() string {
	return CodeRepoPathTemplate
}

func (c *CodeRepo) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type CodeRepoBinding struct {
	APIVersion string                      `yaml:"apiVersion" json:"api_version"`
	Kind       string                      `yaml:"kind" json:"kind" commands:"crb,crbs" applyOrder:"5" removeOrder:"2"`
	Spec       CodeRepoBindingResponseItem `yaml:"spec" json:"spec"`
}

type CodeRepoBindingResponse struct {
	Items []*CodeRepoBindingResponseItem `json:"items"`
}

type CodeRepoBindingResponseItem struct {
	Name        string   `yaml:"name" json:"name" column:"name"`
	ProductName string   `yaml:"productName" json:"product_name"`
	Product     string   `yaml:"product" json:"product" column:"product"`
	CodeRepo    string   `yaml:"coderepo" json:"coderepo" column:"coderepo"`
	Permissions string   `yaml:"permissions" json:"permissions" column:"permissions"`
	Projects    []string `yaml:"projects" json:"projects" column:"projects"`
}

func (c *CodeRepoBinding) GetKind() string {
	return c.Kind
}

func (c *CodeRepoBinding) GetPathTemplate() string {
	return CodeRepoBindingPathTemplate
}

func (c *CodeRepoBinding) GetPathVarNames() []string {
	return []string{"ProductName", "Name"}
}

type ProjectPipelineRuntime struct {
	APIVersion string                             `yaml:"apiVersion" json:"api_version"`
	Kind       string                             `yaml:"kind" json:"kind" commands:"ppr,pprs" applyOrder:"6" removeOrder:"1"`
	Spec       ProjectPipelineRuntimeResponseItem `yaml:"spec" json:"spec"`
}

type ProjectPipelineRuntimeCommonPipelines struct {
	Name  string `yaml:"name" json:"name"`
	Label string `yaml:"label" json:"label"`
	Path  string `yaml:"path" json:"path"`
}

type ProjectPipelineRuntimeCommonDestination struct {
	Environment string `yaml:"environment" json:"environment" column:"environment"`
	Namespace   string `yaml:"namespace" json:"namespace" column:"namespace" mergeTo:"environment"`
}

// ProjectPipelineRuntimeAdditionalResources defines the additional resources witch runtime needed
type ProjectPipelineRuntimeAdditionalResources struct {
	// Optional
	Git ProjectPipelineRuntimeAdditionalResourcesGit `yaml:"git" json:"git"`
}

// ProjectPipelineRuntimeAdditionalResourcesGit defines the additional resources if it comes from git
type ProjectPipelineRuntimeAdditionalResourcesGit struct {
	// Optional
	CodeRepo string `yaml:"codeRepo" json:"coderepo"`
	// Optional
	// If git repo is a public repo, use url instead
	URL      string `yaml:"url" json:"url"`
	Revision string `yaml:"revision" json:"revision"`
	Path     string `yaml:"path" json:"path"`
}

type ProjectPipelineRuntimeResponse struct {
	Items []*ProjectPipelineRuntimeResponseItem `json:"items"`
}

type ProjectPipelineRuntimeResponseItem struct {
	Name        string                                   `yaml:"name" json:"name" column:"name"`
	Project     string                                   `yaml:"project" json:"project" column:"project"`
	Destination *ProjectPipelineRuntimeCommonDestination `yaml:"destination" json:"destination"`
	Isolation   string                                   `yaml:"isolation" json:"isolation"`
	Pipelines   *[]ProjectPipelineRuntimeCommonPipelines `yaml:"pipelines" json:"pipelines"`
	// Optional
	Product             string                                              `yaml:"product" json:"product"`
	PipelineSource      string                                              `yaml:"pipelineSource" json:"pipeline_source" column:"PipelineSource"`
	EventSources        *[]ProjectPipelineRuntimeResponseItemEventSources   `yaml:"eventSources" json:"event_sources"`
	PipelineTriggers    *ProjectPipelineRuntimeResponseItemPipelineTriggers `yaml:"pipelineTriggers" json:"pipeline_triggers"`
	AdditionalResources ProjectPipelineRuntimeAdditionalResources           `yaml:"additionalResources" json:"additional_resources"`
}

type ProjectPipelineRuntimeResponseItemEventSources struct {
	Name     string                                                  `yaml:"name" json:"name"`
	Gitlab   *ProjectPipelineRuntimeResponseItemEventSourcesGitlab   `yaml:"gitlab" json:"gitlab"`
	Calendar *ProjectPipelineRuntimeResponseItemEventSourcesCalendar `yaml:"calendar" json:"calendar"`
}

type ProjectPipelineRuntimeResponseItemEventSourcesGitlab struct {
	RepoName string   `yaml:"repoName" json:"repo_name"  column:"RepoName"`
	Revision string   `yaml:"revision" json:"revision"`
	Events   []string `yaml:"events" json:"events"`
}

type ProjectPipelineRuntimeResponseItemEventSourcesCalendar struct {
	Schedule       string   `yaml:"schedule" json:"schedule"`
	Interval       string   `yaml:"interval" json:"interval"`
	ExclusionDates []string `yaml:"exclusionDates" json:"exclusion_dates"`
	Timezone       string   `yaml:"timezone" json:"timezone"`
}

type ProjectPipelineRuntimeResponseItemPipelineTriggers []struct {
	EventSource string `yaml:"eventSource" json:"event_source" column:"EventSource" mergeTo:"RepoName"`
	Pipeline    string `yaml:"pipeline" json:"pipeline"`
	Revision    string `yaml:"revision" json:"revision"`
}

func (p *ProjectPipelineRuntime) GetKind() string {
	return p.Kind
}

func (p *ProjectPipelineRuntime) GetPathTemplate() string {
	return ProjectPipelineRuntimePathTemplate
}

func (p *ProjectPipelineRuntime) GetPathVarNames() []string {
	return []string{"Product", "Name"}
}

type DeploymentRuntime struct {
	APIVersion string                        `yaml:"apiVersion" json:"api_version"`
	Kind       string                        `yaml:"kind" json:"kind" commands:"dr,drs" applyOrder:"7" removeOrder:"0"`
	Spec       DeploymentRuntimeResponseItem `yaml:"spec" json:"spec"`
}

type DeploymentRuntimeResponse struct {
	Items []*DeploymentRuntimeResponseItem `yaml:"items" json:"items"`
}

type DeploymentRuntimeResponseItem struct {
	Name           string                                       `yaml:"name" json:"name" column:"name"`
	Product        string                                       `yaml:"product" json:"product" column:"product"`
	ManifestSource *DeploymentRuntimeResponseItemManifestSource `yaml:"manifestsource" json:"manifest_source"`
	ProjectsRef    []string                                     `yaml:"projectsRef" json:"projects_ref" column:"projectsRef"`
	Destination    *DeploymentRuntimeResponseItemDestination    `yaml:"destination" json:"destination"`
}

type DeploymentRuntimeResponseItemManifestSource struct {
	CodeRepo       string `yaml:"codeRepo" json:"code_repo" column:"codeRepo"`
	TargetRevision string `yaml:"targetRevision" json:"target_revision" column:"targetRevision" mergeTo:"codeRepo"`
	Path           string `yaml:"path" json:"path" column:"path" mergeTo:"codeRepo"`
}

type DeploymentRuntimeResponseItemDestination struct {
	Environment string   `yaml:"environment" json:"environment" column:"environment"`
	Namespaces  []string `yaml:"namespaces" json:"namespaces" column:"namespaces" mergeTo:"environment"`
}

func (d *DeploymentRuntime) GetKind() string {
	return d.Kind
}

func (d *DeploymentRuntime) GetPathTemplate() string {
	return DeploymentRuntimePathTemplate
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
