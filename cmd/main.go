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
	"fmt"
	"github.com/nautes-labs/cli/cmd/commands"
	"github.com/nautes-labs/cli/cmd/types"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"sort"
	"strconv"
)

func main() {
	var filePath string
	var clientOpts types.ClientOptions
	var resourcesTypeArr = []types.ResourcesType{
		{
			ResourceType:     reflect.TypeOf(types.Cluster{}),
			ResponseItemType: reflect.TypeOf(types.ClusterResponseItem{}), //response type
		},
		{
			ResourceType:     reflect.TypeOf(types.Product{}),
			ResponseItemType: reflect.TypeOf(types.ProductResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.Environment{}),
			ResponseItemType: reflect.TypeOf(types.EnvironmentResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.Project{}),
			ResponseItemType: reflect.TypeOf(types.ProjectResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.CodeRepo{}),
			ResponseItemType: reflect.TypeOf(types.CodeRepoResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.CodeRepoBinding{}),
			ResponseItemType: reflect.TypeOf(types.CodeRepoBindingResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.ProjectPipelineRuntime{}),
			ResponseItemType: reflect.TypeOf(types.ProjectPipelineRuntimeResponseItem{}),
		},
		{
			ResourceType:     reflect.TypeOf(types.DeploymentRuntime{}),
			ResponseItemType: reflect.TypeOf(types.DeploymentRuntimeResponseItem{}),
		},
		//{
		//	ResourceType:     reflect.TypeOf(types.ArtifactRepo{}),
		//	ResponseItemType: reflect.TypeOf(types.ArtifactRepoResponseItem{}),
		//},
	}

	// create slices to store resources types based on apply and remove orders
	var applyResourceTypes = make([]types.ResourcesType, len(resourcesTypeArr))
	var removeResourceTypes = make([]types.ResourcesType, len(resourcesTypeArr))

	// structure to store order information
	type order struct {
		applyOrderIndex  int
		removeOrderIndex int
		index            int
	}

	var orderArr []order

	// iterate over resource types and extract apply and remove orders from tags
	for idx, value := range resourcesTypeArr {
		// check if the "ResourceKind" field exists in the struct
		if field, ok := value.ResourceType.FieldByName(types.ResourceKind); ok {
			tag := field.Tag
			applyOrder := tag.Get(types.ApplyOrder)
			removeOrder := tag.Get(types.RemoveOrder)

			// parse apply and remove orders as integers
			applyOrderIndex, err := strconv.ParseInt(applyOrder, 10, 32)
			if err != nil {
				commands.CheckError(err)
			}
			removeOrderIndex, err := strconv.ParseInt(removeOrder, 10, 32)
			if err != nil {
				commands.CheckError(err)
			}

			// store order information in the array
			orderArr = append(orderArr, order{
				applyOrderIndex:  int(applyOrderIndex),
				removeOrderIndex: int(removeOrderIndex),
				index:            idx,
			})
		}
	}

	// sort the order array based on apply order index
	sort.Slice(orderArr, func(i, j int) bool {
		return orderArr[i].applyOrderIndex < orderArr[j].applyOrderIndex
	})

	// populate applyResourceTypes with resource types sorted by apply order
	for i := 0; i < len(orderArr); i++ {
		applyResourceTypes[i] = resourcesTypeArr[orderArr[i].index]
		//fmt.Println("apply", i, resourcesTypeArr[orderArr[i].index])
	}

	// sort the order array based on remove order index
	sort.Slice(orderArr, func(i, j int) bool {
		return orderArr[i].removeOrderIndex > orderArr[j].removeOrderIndex
	})

	// populate removeResourceTypes with resource types sorted by remove order
	for i := 0; i < len(orderArr); i++ {
		removeResourceTypes[i] = resourcesTypeArr[orderArr[i].index]
		//fmt.Println("delete", i, resourcesTypeArr[orderArr[i].index])
	}

	var rootCmd = &cobra.Command{
		Use:   "nautes",
		Short: "nautes controls a Nautes API server",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply resources",
		Run: func(cmd *cobra.Command, args []string) {
			if err := commands.Execute(clientOpts.ServerAddr, clientOpts.Token, filePath, clientOpts.SkipCheck, applyResourceTypes, commands.SaveResource); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove resources",
		Run: func(cmd *cobra.Command, args []string) {
			if err := commands.Execute(clientOpts.ServerAddr, clientOpts.Token, filePath, clientOpts.SkipCheck, removeResourceTypes, commands.DeleteResource); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	applyCmd.Flags().BoolVarP(&clientOpts.SkipCheck, "insecure", "i", false, "Skipping the compliance check (optional)")
	err := applyCmd.MarkFlagRequired("file")
	if err != nil {
		commands.CheckError(err)
	}
	rootCmd.AddCommand(applyCmd)

	removeCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the input file (required)")
	removeCmd.Flags().BoolVarP(&clientOpts.SkipCheck, "insecure", "i", false, "Skipping the compliance check (optional)")
	err = removeCmd.MarkFlagRequired("file")
	if err != nil {
		commands.CheckError(err)
	}
	rootCmd.AddCommand(removeCmd)

	rootCmd.PersistentFlags().StringVarP(&clientOpts.Token, "token", "t", "", "Authentication token (required)")
	err = rootCmd.MarkPersistentFlagRequired("token")
	if err != nil {
		commands.CheckError(err)
	}

	rootCmd.PersistentFlags().StringVarP(&clientOpts.ServerAddr, "api-server", "s", "", "URL to API server (required)")
	err = rootCmd.MarkPersistentFlagRequired("api-server")
	if err != nil {
		commands.CheckError(err)
	}

	if os.Getenv("API_SERVER") != "" {
		clientOpts.ServerAddr = os.Getenv("API_SERVER")
		err = rootCmd.PersistentFlags().Set("api-server", clientOpts.ServerAddr)
		if err != nil {
			commands.CheckError(err)
		}
	}

	if os.Getenv("GIT_TOKEN") != "" {
		clientOpts.Token = os.Getenv("GIT_TOKEN")
		err = rootCmd.PersistentFlags().Set("token", clientOpts.Token)
		if err != nil {
			commands.CheckError(err)
		}
	}

	// add get command for resource
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get resources",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()(c, args)
				os.Exit(1)
			}
		},
	}
	for _, rc := range applyResourceTypes {
		getCmd.AddCommand(commands.NewResourceCommand(&clientOpts, rc.ResourceType, rc.ResponseItemType, commands.SubGetCommand)...)
	}
	rootCmd.AddCommand(getCmd)

	// add delete command for resource
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()(c, args)
				os.Exit(1)
			}
		},
	}
	for _, rc := range applyResourceTypes {
		deleteCmd.AddCommand(commands.NewResourceCommand(&clientOpts, rc.ResourceType, rc.ResponseItemType, commands.SubDeleteCommand)...)
	}
	rootCmd.AddCommand(deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
