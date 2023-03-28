package main

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

type Project struct {
	APIVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `json:"kind"`
	Spec       struct {
		Name     string `json:"name"`
		Language string `json:"language"`
		Product  string `json:"product"`
	} `json:"spec"`
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
