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
const _DEPLOYMENTRUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/deploymentruntimes/%s"
const _CLUSTER_PATH_TEMPLATE = "/api/v1/clusters/%s"

type resourceFunc func(apiServer string, token string, resource string, resourceHandler ResourceHandler) error

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
		Git  struct {
			Gitlab struct {
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
		Product           string `json:"product"`
		Name              string `json:"name"`
		Project           string `json:"project"`
		DeploymentRuntime bool   `yaml:"deploymentRuntime" json:"deployment_runtime"`
		PipelineRuntime   bool   `yaml:"pipelineRuntime" json:"pipeline_runtime"`
		Webhook           struct {
			Events    []string `json:"events"`
			Isolation string   `json:"isolation"`
		} `json:"webhook"`
		Git struct {
			Gitlab struct {
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

type DeploymentRuntime struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name           string   `json:"name"`
		Product        string   `json:"product"`
		ProjectsRef    []string `yaml:"projectsRef" json:"projects_ref"`
		Destination    string   `json:"destination"`
		Manifestsource struct {
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

type Cluster struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Spec       struct {
		Name        string `yaml:"name" json:"name"`
		ApiServer   string `yaml:"apiServer" json:"api_server"`
		ClusterKind string `yaml:"clusterKind" json:"cluster_kind"`
		ClusterType string `yaml:"clusterType" json:"cluster_type"`
		Usage       string `yaml:"usage" json:"usage"`
		HostCluster string `yaml:"hostCluster" json:"host_cluster"`
		ArgoCDHost  string `yaml:"argocdHost" json:"argocd_host"`
		Traefik     struct {
			HTTPNodePort  string `yaml:"httpNodePort" json:"http_node_port"`
			HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
		} `yaml:"traefik" json:"traefik"`
		VCluster struct {
			HTTPSNodePort string `yaml:"httpsNodePort" json:"https_node_port"`
		} `yaml:"vcluster" json:"vcluster"`
		Kubeconfig string `yaml:"kubeconfig" json:"kubeconfig"`
	} `yaml:"spec" json:"spec"`
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
