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

const METHOD_DELETE = "DELETE"
const METHOD_POST = "POST"

func main() {
	var filePath, token, apiServer string
	var skipCheck bool

	var resourceTypeArr4Save = []reflect.Type{reflect.TypeOf(Cluster{}),
		reflect.TypeOf(Product{}), reflect.TypeOf(Environment{}), reflect.TypeOf(Project{}),
		reflect.TypeOf(CodeRepo{}), reflect.TypeOf(CodeRepoBinding{}), reflect.TypeOf(ProjectPipelineRuntime{}),
		reflect.TypeOf(DeploymentRuntime{})}

	var resourceTypeArr4Remove = []reflect.Type{reflect.TypeOf(DeploymentRuntime{}), reflect.TypeOf(ProjectPipelineRuntime{}),
		reflect.TypeOf(CodeRepoBinding{}), reflect.TypeOf(CodeRepo{}), reflect.TypeOf(Project{}), reflect.TypeOf(Environment{}),
		reflect.TypeOf(Product{}), reflect.TypeOf(Cluster{})}

	var rootCmd = &cobra.Command{
		Use: "nautes",
	}

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply resources",
		Run: func(cmd *cobra.Command, args []string) {
			if err := execute(apiServer, token, filePath, skipCheck, resourceTypeArr4Save, saveResource); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove resources",
		Run: func(cmd *cobra.Command, args []string) {
			if err := execute(apiServer, token, filePath, skipCheck, resourceTypeArr4Remove, deleteResource); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	applyCmd.Flags().BoolVarP(&skipCheck, "insecure", "i", false, "Skipping the compliance check (optional)")
	applyCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(applyCmd)

	removeCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	removeCmd.Flags().BoolVarP(&skipCheck, "insecure", "i", false, "Skipping the compliance check (optional)")
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

func execute(apiServer string, token string, filePath string, skipCheck bool,
	resourceTypeArr []reflect.Type, resourceFunc resourceFunc) error {
	apiServer = formatApiServer(apiServer)
	fmt.Printf("API server: %s\n", apiServer)

	resourcesMap, err := loadResourcesMap(filePath)
	if err != nil {
		return fmt.Errorf("Failed to load resource file: %w", err)
	}

	// Send requests in the order of the given types.
	for _, resourceType := range resourceTypeArr {
		typeName := resourceType.Name()
		for _, resource := range resourcesMap[typeName] {
			resourceObj := reflect.New(resourceType).Interface().(ResourceHandler)
			if err = resourceFunc(apiServer, token, skipCheck, resource, resourceObj); err != nil {
				return err
			}
		}
	}

	return nil
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

func deleteResource(apiServer string, token string, skipCheck bool, resource string, resourceHandler ResourceHandler) error {
	if err := manageResource(METHOD_DELETE, apiServer, token, skipCheck, resource, resourceHandler); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%s deleted successfully.\n", resourceHandler.getKind()))
	return nil
}

func saveResource(apiServer string, token string, skipCheck bool, resource string, resourceHandler ResourceHandler) error {
	if err := manageResource(METHOD_POST, apiServer, token, skipCheck, resource, resourceHandler); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%s saved successfully.\n", resourceHandler.getKind()))
	return nil
}

func manageResource(method string, apiServer string, token string, skipCheck bool, resource string, resourceHandler ResourceHandler) error {
	err := yaml.Unmarshal([]byte(resource), resourceHandler)
	if err != nil {
		fmt.Println("Error unmarshaling YAML: %w", err)
		os.Exit(1)
	}

	requestUrl, requestBody, err := buildRequestURLAndBodys(apiServer, skipCheck, resourceHandler)
	if err != nil {
		return err
	}

	// An exception occurred in the delete request, continue execution.
	if err := buildAndSendRequest(resourceHandler.getKind(), method, requestUrl, requestBody, token); err != nil && method != METHOD_DELETE {
		return err
	} else if err != nil {
		fmt.Println(err)
	}

	return nil
}

func formatApiServer(apiServer string) string {
	if strings.HasSuffix(apiServer, "/") {
		length := len(apiServer)
		return apiServer[0 : length-1]
	}
	return apiServer
}

func buildRequestURLAndBodys(apiServer string, skipCheck bool, resourceHandler ResourceHandler) (string, []byte, error) {
	specValue := reflect.ValueOf(resourceHandler).Elem().FieldByName("Spec")
	pathVarValues, err := getPathVarValues(specValue.Interface(), resourceHandler.getPathVarNames())
	if err != nil {
		return "", nil, err
	}

	requestURL := apiServer + buildURLByParameters(resourceHandler.getPathTemplate(), pathVarValues)
	if skipCheck {
		requestURL = fmt.Sprintf("%s?insecure_skip_check=%t", requestURL, skipCheck)
	}

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

	if requestBody != nil && method != METHOD_DELETE {
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
		return fmt.Errorf("failed to operate %s: status code %d, response body: %s", kind, resp.StatusCode, string(bodyBytes))
	}

	return nil
}
