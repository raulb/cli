package cmd

/*
Copyright © 2020 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create meroxa pipeline components",
	Long: `Use the create command to create various Meroxa pipeline components
including Connectors and Functions.

If you need to add a resource, try:

$ meroxa add resource <resource-type> [name]`,
}

var createConnectorCmd = &cobra.Command{
	Use:   "connector",
	Short: "create connector",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resource Name
		resName := args[0]

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		cfgString, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		cfg := &Config{}
		if cfgString != "" {
			err = json.Unmarshal([]byte(cfgString), cfg)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		// Process metadata
		metadataString, err := cmd.Flags().GetString("metadata")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		metadata := map[string]string{}
		if metadataString != "" {
			err = json.Unmarshal([]byte(metadataString), &metadata)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		// merge in input
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		err = cfg.Set("input", input)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if !flagRootOutputJSON {
			fmt.Println("Creating connector...")
		}

		con, err := createConnector(name, resName, cfg, metadata, input)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flagRootOutputJSON {
			jsonPrint(con)
		} else {
			fmt.Println("Connector successfully created!")
			prettyPrint("connector", con)
		}
	},
}

var createFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "create function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create function called - Not Implemented")
	},
}

var createPipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "create pipeline",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		p := &meroxa.Pipeline{
			Name: pipelineName,
		}

		// Process metadata
		metadataString, err := cmd.Flags().GetString("metadata")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if metadataString != "" {
			var metadata map[string]string
			err = json.Unmarshal([]byte(metadataString), &metadata)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			p.Metadata = metadata
		}

		if !flagRootOutputJSON {
			fmt.Println("Creating Pipeline...")
		}

		res, err := c.CreatePipeline(ctx, p)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flagRootOutputJSON {
			jsonPrint(res)
		} else {
			fmt.Println("Pipeline successfully created!")
			prettyPrint("pipeline", res)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createConnectorCmd)
	createConnectorCmd.Flags().StringP("name", "n", "", "connector name")
	createConnectorCmd.Flags().StringP("config", "c", "", "connector configuration")
	createConnectorCmd.Flags().StringP("metadata", "m", "", "connector metadata")
	createConnectorCmd.Flags().String("input", "", "command delimeted list of input streams")
	createConnectorCmd.MarkFlagRequired("input")

	createCmd.AddCommand(createFunctionCmd)

	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringP("metadata", "m", "", "pipeline metadata")
}

func createConnector(connectorName string, resourceName string, config *Config, metadata map[string]string, input string) (*meroxa.Connector, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}

	// get resource ID from name
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.GetResourceByName(ctx, resourceName)
	if err != nil {
		return nil, err
	}

	// create connector
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cfg := Config{}
	if config != nil {
		cfg.Merge(*config)
	}

	// merge in input
	if input != "" {
		err = cfg.Set("input", input)
		if err != nil {
			return nil, err
		}
	}

	con, err := c.CreateConnector(ctx, connectorName, res.ID, cfg, metadata)
	if err != nil {
		return nil, err
	}

	return con, nil
}
