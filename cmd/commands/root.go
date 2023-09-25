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

package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nautes-labs/cli/cmd/types"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	MethodGet              = "GET"
	MethodDelete           = "DELETE"
	MethodPost             = "POST"
	IgnoreProductOfCluster = "Cluster"
	IgnoreProductOfProduct = "Product"
	CodeRepoBinding        = "CodeRepoBinding"
)

func Execute(apiServer string, token string, filePath string, skipCheck bool,
	resourceTypeArr []types.ResourcesType, resourceFunc types.ResourceFunc) error {
	apiServer = formatAPIServer(apiServer)
	fmt.Printf("API server: %s\n", apiServer)

	resourcesMap, err := loadResourcesMap(filePath)
	if err != nil {
		return fmt.Errorf("failed to load resource file: %w", err)
	}

	// Send requests in the order of the given types.
	for _, value := range resourceTypeArr {
		typeName := value.ResourceType.Name()
		for _, resource := range resourcesMap[typeName] {
			resourceObj := reflect.New(value.ResourceType).Interface().(types.ResourceHandler)
			if err = resourceFunc(apiServer, token, skipCheck, resource, resourceObj); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadResourcesMap(filePath string) (map[string][]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	resources := strings.Split(string(content), "---")

	fmt.Printf("%d resources found\n\n", len(resources))

	resourcesMap := make(map[string][]string)

	for _, resource := range resources {
		var crMetadata types.Base
		err := yaml.Unmarshal([]byte(resource), &crMetadata)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
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

func DeleteResource(apiServer string, token string, skipCheck bool, resource string, resourceHandler types.ResourceHandler) error {
	if err := manageResource(MethodDelete, apiServer, token, skipCheck, resource, resourceHandler); err != nil {
		return err
	}
	fmt.Printf("%s deleted successfully.\n", resourceHandler.GetKind())
	return nil
}

func SaveResource(apiServer string, token string, skipCheck bool, resource string, resourceHandler types.ResourceHandler) error {
	if err := manageResource(MethodPost, apiServer, token, skipCheck, resource, resourceHandler); err != nil {
		return err
	}
	fmt.Printf("%s saved successfully.\n", resourceHandler.GetKind())
	return nil
}

func manageResource(method string, apiServer string, token string, skipCheck bool, resource string, resourceHandler types.ResourceHandler) error {
	err := yaml.Unmarshal([]byte(resource), resourceHandler)
	if err != nil {
		fmt.Println("error unmarshaling YAML: %w", err)
		return err
	}
	_, err = buildResourceAndDo(method, apiServer, token, skipCheck, resourceHandler)
	if err != nil {
		return err
	}

	return nil
}

func buildResourceAndDo(method string, apiServer string, token string, skipCheck bool, resourceHandler types.ResourceHandler) ([]byte, error) {
	requestURL, requestBody, err := buildRequestURLAndBodys(apiServer, resourceHandler)
	if err != nil {
		return nil, err
	}
	if skipCheck {
		requestURL = fmt.Sprintf("%s?insecure_skip_check=%t", requestURL, skipCheck)
	}
	var resBytes []byte
	// An exception occurred in the delete request, continue execution.
	if resBytes, err = buildAndSendRequest(resourceHandler.GetKind(), method, requestURL, requestBody, token); err != nil {
		return nil, err
	}
	return resBytes, nil
}

func formatAPIServer(apiServer string) string {
	if strings.HasSuffix(apiServer, "/") {
		length := len(apiServer)
		return apiServer[0 : length-1]
	}
	return apiServer
}

func buildRequestURLAndBodys(apiServer string, resourceHandler types.ResourceHandler) (string, []byte, error) {
	specValue := reflect.ValueOf(resourceHandler).Elem().FieldByName("Spec")
	pathVarValues, err := getPathVarValues(specValue.Interface(), resourceHandler.GetPathVarNames())
	if err != nil {
		return "", nil, err
	}
	requestURL := apiServer + buildURLByParameters(resourceHandler.GetPathTemplate(), pathVarValues)

	requestBodyBytes, err := json.Marshal(specValue.Interface())
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return requestURL, requestBodyBytes, nil
}

func buildURLByParameters(template string, pathVarValues []string) string {
	for _, pathVarValue := range pathVarValues {
		template = replaceFirstPlaceholder(template, "%s", pathVarValue)
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
			return nil, fmt.Errorf("type Assertion Failure: %+v", pathVarValue.Interface())
		}
		pathVarValues = append(pathVarValues, pathVarValueStr)
	}

	return pathVarValues, nil
}

func buildAndSendRequest(kind string, method string, requestURL string, requestBody []byte, token string) ([]byte, error) {
	var req *http.Request
	var err error

	fmt.Printf("Request[%s] URL: %s\n", method, requestURL)

	if requestBody != nil && method != MethodDelete {
		fmt.Printf("Request body: %s\n\n", string(requestBody))
		req, err = http.NewRequest(method, requestURL, bytes.NewBuffer(requestBody))
	} else {
		req, err = http.NewRequest(method, requestURL, http.NoBody)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	//sTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()
	//eTime := time.Now()
	//fmt.Printf("sTime: %s\neTime: %s\nsub: %fs\n", sTime, eTime, eTime.Sub(sTime).Seconds())
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		return bodyBytes, nil
	}
	//return nil, fmt.Errorf("failed to operate %s: status code: %d", kind, resp.StatusCode)
	return nil, fmt.Errorf("failed to operate %s:\n%s", kind, bodyBytes)
}
