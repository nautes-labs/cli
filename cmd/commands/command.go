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
	"bufio"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/nautes-labs/cli/cmd/printers"
	"github.com/nautes-labs/cli/cmd/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

const (
	OutputYaml = "yaml"
	OutputJson = "json"
	OutputWide = "wide"
)

// CheckError logs a fatal message and exits with error code if err is not nil
func CheckError(err error) {
	if err != nil {
		Fatal(20, err)
	}
}

// Fatal is a wrapper for logrus.Fatal() to exit with custom code
func Fatal(exitcode int, args ...interface{}) {
	exitfunc := func() {
		os.Exit(exitcode)
	}
	log.RegisterExitHandler(exitfunc)
	log.Fatal(args...)
}

// NewResourceCommand creates and returns a set of Cobra commands for a resource type based on reflection and provided options.
// It takes client options, resource type, response item type, and a subCommandFunc responsible for creating subcommands.
// The generated commands include those for the resource itself, its plural form, and any short commands specified in tags.
// The subCommandFunc is called to create subcommands for each of these names.
func NewResourceCommand(clientOptions *types.ClientOptions, resourceType reflect.Type, responseItemType reflect.Type,
	subCommandFunc func(clientOptions *types.ClientOptions, resourceHandler types.ResourceHandler, resourceName string,
		resourceType, responseItemType reflect.Type) *cobra.Command) (ccCommands []*cobra.Command) {
	// Instantiate a ResourceHandler of the specified type
	resourceHandler := reflect.New(resourceType).Interface().(types.ResourceHandler)

	// Get short commands from tags
	var shortCommands []string
	if field, ok := resourceType.FieldByName(types.ResourceKind); ok {
		tag := field.Tag
		shortCommands = strings.Split(tag.Get("commands"), ",")
	}

	// Set the 'kind' value in the ResourceHandler
	resourcePtr := reflect.ValueOf(resourceHandler).Elem()
	rType := reflect.TypeOf(resourceHandler).String()
	resourcePtr.FieldByName(types.ResourceKind).SetString(strings.TrimPrefix(rType, "*types."))

	// Generate resource name and related command names
	var resourceName = strings.ToLower(resourceHandler.GetKind())
	var allCmds []string
	allCmds = append(allCmds, resourceName, fmt.Sprintf("%ss", resourceName))
	allCmds = append(allCmds, shortCommands...)

	// Create commands using the subCommandFunc for each name
	for _, cmd := range allCmds {
		command := subCommandFunc(clientOptions, resourceHandler, cmd, resourceType, responseItemType)
		ccCommands = append(ccCommands, command)
	}

	// Return the generated Cobra commands
	return ccCommands
}

// SubGetCommand creates a Cobra command for the "get" subcommand of a resource.
// It retrieves information about a specific resource or a list of resources based on the provided arguments.
// The command supports various output formats such as json, yaml, or a wide table format.
// The "product" flag allows filtering resources by product name.
func SubGetCommand(clientOptions *types.ClientOptions, resourceHandler types.ResourceHandler, resourceName string, _, responseItemType reflect.Type) *cobra.Command {
	var (
		output  string
		product string
	)

	// Reflect on the resource handler and initialize some variables
	resourceValue := reflect.ValueOf(resourceHandler).Elem()
	resourceKind := resourceHandler.GetKind()
	var resourceNameUpper = strings.ToUpper(resourceKind)

	// Create the "get" command
	var command = &cobra.Command{
		Use:   fmt.Sprintf("%s name", resourceName),
		Short: fmt.Sprintf("Get %s information", resourceNameUpper),
		Example: fmt.Sprintf(`nautes get %s example-name

nautes get %s name-101 name-102`, resourceName, resourceName),

		Run: func(c *cobra.Command, args []string) {
			var responseValue reflect.Value

			// Process the "product" flag to filter resources by product name
			if product != "" {
				if resourceKind != IgnoreProductOfCluster && resourceKind != IgnoreProductOfProduct {
					if resourceKind == CodeRepoBinding {
						resourceValue.FieldByName("Spec").FieldByName("ProductName").SetString(product)
					} else {
						resourceValue.FieldByName("Spec").FieldByName("Product").SetString(product)
					}
				}
			}

			var outputFlag bool
			resourceResponseList := make([]interface{}, 0)
			resourceResponseListValue := make([]reflect.Value, 0)

			if len(args) == 0 {
				// Retrieve a list of resources
				// Dynamic creation of a struct for storing resource items
				sliceType := reflect.SliceOf(responseItemType)
				sliceValue := reflect.MakeSlice(sliceType, 5, 10)
				slice := sliceValue.Interface()
				fields := []reflect.StructField{
					{
						Name: "Items",
						Type: reflect.TypeOf(slice),
					},
				}
				customStructType := reflect.StructOf(fields)
				responseValue = reflect.New(customStructType)

				// Build and retrieve resources from the server
				resBytes, err := buildResourceAndDo(MethodGet, clientOptions.ServerAddr, clientOptions.Token, clientOptions.SkipCheck, resourceHandler)
				if err != nil {
					CheckError(err)
				}
				err = json.Unmarshal(resBytes, responseValue.Interface())
				if err != nil {
					CheckError(err)
				}

				// Extract individual resource items from the response
				instance := responseValue.Elem().FieldByName("Items")
				instanceLen := instance.Len()
				for i := 0; i < instanceLen; i++ {
					item := instance.Index(i)
					resourceResponseList = append(resourceResponseList, item.Interface())
					resourceResponseListValue = append(resourceResponseListValue, item)
				}
			} else {
				// Retrieve specific resources by name
				for _, argsSelector := range args {
					resourceValue.FieldByName("Spec").FieldByName("Name").SetString(argsSelector)
					resBytes, err := buildResourceAndDo(MethodGet, clientOptions.ServerAddr, clientOptions.Token, clientOptions.SkipCheck, resourceHandler)
					if err != nil {
						CheckError(err)
					}
					responseValue = reflect.New(responseItemType)
					err = json.Unmarshal(resBytes, responseValue.Interface())
					if err != nil {
						CheckError(err)
					}
					resourceResponseList = append(resourceResponseList, responseValue.Interface())
					resourceResponseListValue = append(resourceResponseListValue, responseValue)
				}
				if len(args) == 1 {
					outputFlag = true
				}
			}

			// Output formatting based on the specified format
			switch output {
			case OutputYaml, OutputJson:
				err := PrintResourceResponseList(resourceResponseList, output, outputFlag)
				CheckError(err)
			case OutputWide, "":
				table, err := printers.GenerateTable(resourceResponseListValue, responseItemType)
				CheckError(err)
				err = printers.PrintTable(table, os.Stdout)
				CheckError(err)
			default:
				CheckError(fmt.Errorf("unknown output format: %s", output))
			}
		},
	}

	// Add flags to the command
	command.Flags().StringVarP(&output, "output", "o", "wide", "Output format. One of: json|yaml|wide")
	if resourceKind != IgnoreProductOfCluster && resourceKind != IgnoreProductOfProduct {
		command.Flags().StringVarP(&product, "product", "p", "", "List resource by product name")
		err := command.MarkFlagRequired("product")
		if err != nil {
			CheckError(err)
		}

		// Set the "product" flag from the environment variable if present
		if os.Getenv("PRODUCT") != "" {
			err = command.Flags().Set("product", os.Getenv("PRODUCT"))
			if err != nil {
				CheckError(err)
			}
		}
	}

	return command
}

// SubDeleteCommand creates a Cobra command for the "delete" subcommand of a resource.
// It removes one or more resources based on the provided arguments and options.
// The command supports confirmation prompts and the option to bypass prompts using the "--yes" flag.
// The "product" flag allows filtering resources by product name.
func SubDeleteCommand(clientOptions *types.ClientOptions, resourceHandler types.ResourceHandler, resourceName string, _, _ reflect.Type) *cobra.Command {
	var noPrompt bool
	var product string
	resourceValue := reflect.ValueOf(resourceHandler).Elem()
	resourceKind := resourceHandler.GetKind()
	var resourceNameUpper = strings.ToUpper(resourceKind)
	var command = &cobra.Command{
		Use:   fmt.Sprintf("%s name", resourceName),
		Short: fmt.Sprintf("Remove %s credentials", resourceNameUpper),
		Example: fmt.Sprintf(`nautes delete %s example-name

nautes delete %s name-101 name-102`, resourceName, resourceName),

		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()(c, args)
				os.Exit(1)
			}
			if product != "" {
				if resourceKind != IgnoreProductOfCluster && resourceKind != IgnoreProductOfProduct {
					if resourceKind == CodeRepoBinding {
						resourceValue.FieldByName("Spec").FieldByName("ProductName").SetString(product)
					} else {
						resourceValue.FieldByName("Spec").FieldByName("Product").SetString(product)
					}
				}
			}
			var isConfirmAll bool
			for _, argsSelector := range args {
				var lowercaseAnswer string
				if !noPrompt {
					if !isConfirmAll {
						lowercaseAnswer = AskToProceedS("Are you sure you want to remove '" + argsSelector +
							"'? [y/n/A] where 'A' is to remove all specified resources without prompting. ")
						if lowercaseAnswer == "a" {
							lowercaseAnswer = "y"
							isConfirmAll = true
						}
					} else {
						lowercaseAnswer = "y"
					}
				} else {
					lowercaseAnswer = "y"
				}

				if lowercaseAnswer == "y" {
					resourceValue.FieldByName("Spec").FieldByName("Name").SetString(argsSelector)
					_, err := buildResourceAndDo(MethodDelete, clientOptions.ServerAddr, clientOptions.Token, clientOptions.SkipCheck, resourceHandler)
					if err != nil {
						CheckError(err)
					}
					fmt.Printf("%s '%s' removed\n", resourceHandler.GetKind(), argsSelector)
				} else {
					fmt.Println("The command to remove '" + argsSelector + "' was canceled.")
				}
			}
		},
	}
	command.Flags().BoolVarP(&noPrompt, "yes", "y", false, "Turn off prompting to confirm remove of resources")
	if resourceKind != IgnoreProductOfCluster && resourceKind != IgnoreProductOfProduct {
		command.Flags().StringVarP(&product, "product", "p", "", "List resource by product name")
		err := command.MarkFlagRequired("product")
		if err != nil {
			CheckError(err)
		}

		if os.Getenv("PRODUCT") != "" {
			err = command.Flags().Set("product", os.Getenv("PRODUCT"))
			if err != nil {
				CheckError(err)
			}
		}
	}
	return command
}

// AskToProceedS prompts the user with a message (typically a yes, no or all question) and returns string
// "a", "y" or "n".
func AskToProceedS(message string) string {
	for {
		fmt.Print(message)
		reader := bufio.NewReader(os.Stdin)
		proceedRaw, err := reader.ReadString('\n')
		CheckError(err)
		switch strings.ToLower(strings.TrimSpace(proceedRaw)) {
		case "y", "yes":
			return "y"
		case "n", "no":
			return "n"
		case "a", "all":
			return "a"
		}
	}
}

// PrintResource prints a single resource in YAML or JSON format to stdout according to the output format
func PrintResource(resource interface{}, output string) error {
	switch output {
	case OutputJson:
		jsonBytes, err := jsoniter.MarshalIndent(resource, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal resource to json: %w", err)
		}
		fmt.Println(string(jsonBytes))
	case OutputYaml:
		yamlBytes, err := yaml.Marshal(resource)
		if err != nil {
			return fmt.Errorf("unable to marshal resource to yaml: %w", err)
		}
		fmt.Print(string(yamlBytes))
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
	return nil
}

// PrintResourceResponseList marshals & prints a list of resources to stdout according to the output format
func PrintResourceResponseList(resources interface{}, output string, single bool) error {
	kt := reflect.ValueOf(resources)
	//Sometimes, we want to marshal the first resource of a slice or array as single item
	if kt.Kind() == reflect.Slice || kt.Kind() == reflect.Array {
		if single && kt.Len() == 1 {
			return PrintResource(kt.Index(0).Interface(), output)
		}

		// If we have a zero len list, prevent printing "null"
		if kt.Len() == 0 {
			return PrintResource([]string{}, output)
		}
	}

	switch output {
	case OutputJson:
		jsonBytes, err := jsoniter.MarshalIndent(resources, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal resources to json: %w", err)
		}
		_, err = fmt.Println(string(jsonBytes))
		if err != nil {
			return err
		}
	case OutputYaml:
		yamlBytes, err := yaml.Marshal(resources)
		if err != nil {
			return fmt.Errorf("unable to marshal resources to yaml: %w", err)
		}
		_, err = fmt.Print(string(yamlBytes))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
	return nil
}
