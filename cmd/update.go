package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

var (
	state string // connector state
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a component",
	Long:  `Update a component of the Meroxa platform, including connectors`,
}

var updateConnectorCmd = &cobra.Command{
	Use:     "connector <name> --state <pause|resume|restart>",
	Aliases: []string{"connectors"},
	Short:   "Update connector state",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires connector name\n\nUsage:\n  meroxa update connector <name> --state <pause|resume|restart>")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Connector Name
		conName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// call meroxa-go to update connector status with name
		if !flagRootOutputJSON {
			fmt.Printf("Updating %s connector...\n", conName)
		}

		con, err := c.UpdateConnectorStatus(ctx, conName, state)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(con)
		} else {
			fmt.Printf("Connector %s successfully updated!\n", con.Name)
		}

		return nil
	},
}

var updateResourceCmd = &cobra.Command{
	Use:     "resource <resource-name>",
	Short:   "Update a resource",
	Long:    `Use the update command to update various Meroxa resources.`,
	Aliases: []string{"resources"},
	// TODO: Change the design so a new name for the resource could be set
	// meroxa update resource <old-resource-name> --name <new-resource-name>
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || (resURL == "" && resMetadata == "" && resCredentials == "") {
			return errors.New("requires a resource name and either `--metadata`, `--url` or `--credentials` to update the resource \n\nUsage:\n  meroxa update resource <resource-name> [--url <url> | --metadata <metadata> | --credentials <credentials>]")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Resource Name
		resName = args[0]
		c, err := client()

		if err != nil {
			return err
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var res meroxa.UpdateResourceInput

		// If url was provided, update it
		if resURL != "" {
			res.URL = resURL
		}

		// TODO: Figure out best way to handle creds and metadata
		// Get credentials (expect a JSON string)
		if resCredentials != "" {
			var creds meroxa.Credentials
			err = json.Unmarshal([]byte(resCredentials), &creds)
			if err != nil {
				return err
			}

			res.Credentials = &creds
		}

		// If metadata was provided, update it
		if resMetadata != "" {
			var metadata map[string]string
			err = json.Unmarshal([]byte(resMetadata), &metadata)
			if err != nil {
				return err
			}

			res.Metadata = metadata
		}

		// call meroxa-go to update resource
		if !flagRootOutputJSON {
			fmt.Printf("Updating %s resource...\n", resName)
		}

		resource, err := c.UpdateResource(ctx, resName, res)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(resource)
		} else {
			fmt.Printf("Resource %s successfully updated!\n", resName)
		}

		return nil
	},
}

var updatePipelineCmd = &cobra.Command{
	Use:     "pipeline <name> --state <pause|resume|restart>",
	Aliases: []string{"pipelines"},
	Short:   "Update pipeline state",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires pipeline name\n\nUsage:\n  meroxa update pipeline <name> --state <pause|resume|restart>")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Pipeline Name
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get pipeline id from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		pipeline, err := c.GetPipelineByName(ctx, pipelineName)
		if err != nil {
			return err
		}

		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		// call meroxa-go to update pipeline status with name
		if !flagRootOutputJSON {
			fmt.Printf("Updating %s pipeline...\n", pipelineName)
		}

		p, err := c.UpdatePipelineStatus(ctx, pipeline.ID, state)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(p)
		} else {
			fmt.Printf("Pipeline %s successfully updated!\n", p.Name)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Subcommands
	updateCmd.AddCommand(updateConnectorCmd)
	updateConnectorCmd.Flags().StringVarP(&state, "state", "", "", "connector state")
	updateConnectorCmd.MarkFlagRequired("state")

	updateCmd.AddCommand(updatePipelineCmd)
	updatePipelineCmd.Flags().StringVarP(&state, "state", "", "", "pipeline state")
	updatePipelineCmd.MarkFlagRequired("state")

	updateCmd.AddCommand(updateResourceCmd)
	RootCmd.AddCommand(updateCmd)

	updateResourceCmd.Flags().StringVarP(&resURL, "url", "u", "", "resource url")
	updateResourceCmd.Flags().StringVarP(&resMetadata, "metadata", "m", "", "resource metadata")
	updateResourceCmd.Flags().StringVarP(&resCredentials, "credentials", "", "", "resource credentials")
}
