package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const _PRODUCT_PATH_TEMPLATE = "/api/v1/products/%s"
const _ENV_PATH_TEMPLATE = "/api/v1/products/%s/environments/%s"
const _PROJECT_PATH_TEMPLATE = "/api/v1/products/%s/projects/%s"
const _CODEREPO_PATH_TEMPLATE = "/api/v1/products/%s/coderepos/%s"
const _DEPLOYMENTRUNTIME_PATH_TEMPLATE = "/api/v1/products/%s/deploymentruntimes/%s"

const _KIND_PRODUCT = "Product"
const _KIND_ENVIRONMENT = "Environment"
const _KIND_PROJECT = "Project"
const _KIND_CODEREPO = "CodeRepo"
const _KIND_DEPLOYMENTRUNTIME = "DeploymentRuntime"

func main() {
	var filePath, token, apiServer string

	var rootCmd = &cobra.Command{
		Use: "nautes",
	}

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply resources",
		Run: func(cmd *cobra.Command, args []string) {

			apiServer = formatApiServer(apiServer)
			fmt.Printf("API server: %s\n", apiServer)

			resourcesMap, err := loadResourcesMap(filePath)
			if err != nil {
				fmt.Println("Failed to load resource file: %w", err)
				os.Exit(1)
			}

			for _, resource := range resourcesMap[_KIND_PRODUCT] {
				var product Product
				err := yaml.Unmarshal([]byte(resource), &product)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = saveResource(_KIND_PRODUCT, token, product, apiServer, _PRODUCT_PATH_TEMPLATE, []string{"Name"})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			for _, resource := range resourcesMap[_KIND_ENVIRONMENT] {
				var env Environment
				err := yaml.Unmarshal([]byte(resource), &env)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = saveResource(_KIND_ENVIRONMENT, token, env, apiServer, _ENV_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			for _, resource := range resourcesMap[_KIND_PROJECT] {
				var project Project
				err := yaml.Unmarshal([]byte(resource), &project)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = saveResource(_KIND_PROJECT, token, project, apiServer, _PROJECT_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			for _, resource := range resourcesMap[_KIND_CODEREPO] {
				var coderepo CodeRepo
				err := yaml.Unmarshal([]byte(resource), &coderepo)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = saveResource(_KIND_CODEREPO, token, coderepo, apiServer, _CODEREPO_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			for _, resource := range resourcesMap[_KIND_DEPLOYMENTRUNTIME] {
				var dr DeploymentRuntime
				err := yaml.Unmarshal([]byte(resource), &dr)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = saveResource(_KIND_DEPLOYMENTRUNTIME, token, dr, apiServer, _DEPLOYMENTRUNTIME_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove resources",
		Run: func(cmd *cobra.Command, args []string) {
			apiServer = formatApiServer(apiServer)
			fmt.Printf("API server: %s\n", apiServer)

			resourcesMap, err := loadResourcesMap(filePath)
			if err != nil {
				fmt.Println("Failed to load resource file: %w", err)
				os.Exit(1)
			}

			for _, resource := range resourcesMap[_KIND_DEPLOYMENTRUNTIME] {
				var dr DeploymentRuntime
				err := yaml.Unmarshal([]byte(resource), &dr)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = deleteResource(_KIND_DEPLOYMENTRUNTIME, token, dr, apiServer, _DEPLOYMENTRUNTIME_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, resource := range resourcesMap[_KIND_CODEREPO] {
				var coderepo CodeRepo
				err := yaml.Unmarshal([]byte(resource), &coderepo)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = deleteResource(_KIND_CODEREPO, token, coderepo, apiServer, _CODEREPO_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, resource := range resourcesMap[_KIND_PROJECT] {
				var project Project
				err := yaml.Unmarshal([]byte(resource), &project)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = deleteResource(_KIND_PROJECT, token, project, apiServer, _PROJECT_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, resource := range resourcesMap[_KIND_ENVIRONMENT] {
				var env Environment
				err := yaml.Unmarshal([]byte(resource), &env)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = deleteResource(_KIND_ENVIRONMENT, token, env, apiServer, _ENV_PATH_TEMPLATE, []string{"Product", "Name"})
				if err != nil {
					fmt.Println(err)
				}
			}

			for _, resource := range resourcesMap[_KIND_PRODUCT] {
				var product Product
				err := yaml.Unmarshal([]byte(resource), &product)
				if err != nil {
					fmt.Println("Error unmarshaling YAML: %w", err)
					os.Exit(1)
				}

				err = deleteResource(_KIND_PRODUCT, token, product, apiServer, _PRODUCT_PATH_TEMPLATE, []string{"Name"})
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	}

	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	applyCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(applyCmd)

	removeCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	removeCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(removeCmd)

	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Authentication token (required)")
	rootCmd.MarkPersistentFlagRequired("token")

	rootCmd.PersistentFlags().StringVarP(&apiServer, "api-server", "s", "", "URL to API server (required)")
	rootCmd.MarkPersistentFlagRequired("api-server")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadResourcesMap(filePath string) (map[string][]string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	resources := strings.Split(string(content), "---")

	fmt.Printf("%d resources found\n\n", len(resources))

	resourcesMap := make(map[string][]string)

	for _, resource := range resources {
		var crMetadata Base
		err := yaml.Unmarshal([]byte(resource), &crMetadata)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling YAML: %w", err)
		}
		resourceList := resourcesMap[crMetadata.Kind]
		if resourceList == nil {
			resourceList = []string{resource}
		} else {
			resourceList = append(resourceList, resource)
		}
		resourcesMap[crMetadata.Kind] = resourceList
	}

	return resourcesMap, nil
}

func deleteResource(kind string, token string, resourceObj interface{},
	apiServer string, pathTemplate string, pathVarNames []string) error {
	requestUrl, _, err := buildRequestURLAndBodys(resourceObj, apiServer, pathTemplate, pathVarNames)
	if err != nil {
		return err
	}

	if err := buildAndSendRequest(kind, "DELETE", requestUrl, nil, token); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s deleted successfully.\n", kind))
	return nil
}

func saveResource(kind string, token string, resourceObj interface{},
	apiServer string, pathTemplate string, pathVarNames []string) error {
	requestUrl, requestBody, err := buildRequestURLAndBodys(resourceObj, apiServer, pathTemplate, pathVarNames)
	if err != nil {
		return err
	}

	if err := buildAndSendRequest(kind, "POST", requestUrl, requestBody, token); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s saved successfully.\n", kind))
	return nil
}

func formatApiServer(apiServer string) string {
	if strings.HasSuffix(apiServer, "/") {
		length := len(apiServer)
		return apiServer[0 : length-1]
	}
	return apiServer
}

func buildRequestURLAndBodys(resourceObj interface{},
	apiServer string, pathTemplate string, pathVarNames []string) (string, []byte, error) {

	_, ok := reflect.TypeOf(resourceObj).FieldByName("Spec")

	if !ok {
		return "", nil, fmt.Errorf("Spec field not found: %+v", resourceObj)
	}

	specValue := reflect.ValueOf(resourceObj).FieldByName("Spec")

	pathVarValues, err := getPathVarValues(specValue.Interface(), pathVarNames)
	if err != nil {
		return "", nil, err
	}
	requestURL := apiServer + buildURLByParameters(pathTemplate, pathVarValues)

	requestBodyBytes, err := json.Marshal(specValue.Interface())

	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return requestURL, requestBodyBytes, nil
}

func buildURLByParameters(template string, pathVarValues []string) string {
	if pathVarValues != nil {
		for _, pathVarValue := range pathVarValues {
			template = replaceFirstPlaceholder(template, "%s", pathVarValue)
		}
	}

	return template
}

func replaceFirstPlaceholder(s, placeholder, replacement string) string {
	index := strings.Index(s, placeholder)
	if index == -1 {
		return s
	}
	return s[:index] + replacement + s[index+len(placeholder):]
}

func getPathVarValues(specObj interface{}, pathVarNames []string) ([]string, error) {
	pathVarValues := make([]string, 0, len(pathVarNames))
	for _, pathVarName := range pathVarNames {
		_, ok := reflect.TypeOf(specObj).FieldByName(pathVarName)

		if !ok {
			return nil, fmt.Errorf("%s field not found: %+v", pathVarName, specObj)
		}

		pathVarValue := reflect.ValueOf(specObj).FieldByName(pathVarName)
		pathVarValueStr, ok := pathVarValue.Interface().(string)

		if !ok {
			return nil, fmt.Errorf("Type Assertion Failure: %+v", pathVarValue.Interface())
		}
		pathVarValues = append(pathVarValues, pathVarValueStr)
	}

	return pathVarValues, nil
}

func buildAndSendRequest(kind string, method string, requestURL string, requestBody []byte, token string) error {
	var req *http.Request
	var err error

	fmt.Printf("Request[%s] URL: %s\n", method, requestURL)

	if requestBody != nil {
		fmt.Printf("Request body: %s\n", string(requestBody))
		req, err = http.NewRequest(method, requestURL, bytes.NewBuffer(requestBody))
	} else {
		req, err = http.NewRequest(method, requestURL, http.NoBody)
	}

	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("failed to create %s: status code %d, response body: %s", kind, resp.StatusCode, string(bodyBytes))
	}

	return nil
}
