package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDatabaseTools registers all database-related tools with the MCP server
func RegisterDatabaseTools(s *server.MCPServer, client *gobizfly.Client) {
	// List databases tool
	listDatabasesTool := mcp.NewTool("bizflycloud_list_databases",
		mcp.WithDescription("List all Bizfly Cloud databases"),
	)
	s.AddTool(listDatabasesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get available engines first
		engines, err := client.CloudDatabase.Engines().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list database engines: %v", err)), nil
		}

		// Get databases
		databases, err := client.CloudDatabase.Instances().List(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list databases: %v", err)), nil
		}

		result := "Available databases:\n\n"
		for _, db := range databases {
			result += fmt.Sprintf("Database: %s\n", db.Name)
			result += fmt.Sprintf("  ID: %s\n", db.ID)
			result += fmt.Sprintf("  DataStore Type: %s\n", db.Datastore.Type)
			result += fmt.Sprintf("  DataStore ID: %s\n", db.Datastore.ID)
			
			// Find matching engine
			for _, engine := range engines {
				if engine.ID == db.Datastore.ID {
					result += fmt.Sprintf("  Engine Name: %s\n", engine.Name)
					if len(engine.Versions) > 0 {
						result += fmt.Sprintf("  Available Versions: %v\n", engine.Versions)
					}
					break
				}
			}
			
			result += fmt.Sprintf("  Status: %s\n", db.Status)
			result += fmt.Sprintf("  Created At: %s\n", db.CreatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// List datastores tool
	listDatastoresTool := mcp.NewTool("bizflycloud_list_datastores",
		mcp.WithDescription("List all available Bizfly Cloud database engines and versions"),
	)
	s.AddTool(listDatastoresTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		engines, err := client.CloudDatabase.Engines().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list database engines: %v", err)), nil
		}

		result := "Available database engines and versions:\n\n"
		for _, engine := range engines {
			result += fmt.Sprintf("Database Engine: %s\n", engine.Name)
			result += fmt.Sprintf("  ID: %s\n", engine.ID)
			if len(engine.Versions) > 0 {
				result += fmt.Sprintf("  Available Versions:\n")
				for _, version := range engine.Versions {
					result += fmt.Sprintf("    - %s\n", version)
				}
			}
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create database tool
	createDatabaseTool := mcp.NewTool("bizflycloud_create_database",
		mcp.WithDescription("Create a new Bizfly Cloud database"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the database"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of database (mysql, postgresql, mongodb)"),
		),
		mcp.WithString("version",
			mcp.Required(),
			mcp.Description("Version of the database"),
		),
		mcp.WithString("flavor",
			mcp.Required(),
			mcp.Description("Flavor name for the database instance"),
		),
		mcp.WithNumber("volume_size",
			mcp.Required(),
			mcp.Description("Size of the volume in GB"),
		),
		mcp.WithString("availability_zone",
			mcp.Required(),
			mcp.Description("Availability zone for the database"),
		),
	)
	s.AddTool(createDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		dbType, ok := request.Params.Arguments["type"].(string)
		if !ok {
			return nil, errors.New("type must be a string")
		}
		version, ok := request.Params.Arguments["version"].(string)
		if !ok {
			return nil, errors.New("version must be a string")
		}
		flavor, ok := request.Params.Arguments["flavor"].(string)
		if !ok {
			return nil, errors.New("flavor must be a string")
		}
		volumeSize, ok := request.Params.Arguments["volume_size"].(float64)
		if !ok {
			return nil, errors.New("volume_size must be a number")
		}
		availabilityZone, ok := request.Params.Arguments["availability_zone"].(string)
		if !ok {
			return nil, errors.New("availability_zone must be a string")
		}

		database, err := client.CloudDatabase.Instances().Create(ctx, &gobizfly.CloudDatabaseInstanceCreate{
			Name: name,
			Datastore: gobizfly.CloudDatabaseDatastore{
				Type: dbType,
				ID:   version,
			},
			FlavorName:       flavor,
			VolumeSize:       int(volumeSize),
			AvailabilityZone: availabilityZone,
			Networks:         []gobizfly.CloudDatabaseNetworks{{}}, // Default network
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create database: %v", err)), nil
		}

		result := fmt.Sprintf("Database created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", database.Name)
		result += fmt.Sprintf("  ID: %s\n", database.ID)
		result += fmt.Sprintf("  DataStore Type: %s\n", database.Datastore.Type)
		result += fmt.Sprintf("  DataStore ID: %s\n", database.Datastore.ID)
		result += fmt.Sprintf("  Status: %s\n", database.Status)
		result += fmt.Sprintf("  Created At: %s\n", database.CreatedAt)
		return mcp.NewToolResultText(result), nil
	})

	// Delete database tool
	deleteDatabaseTool := mcp.NewTool("bizflycloud_delete_database",
		mcp.WithDescription("Delete a Bizfly Cloud database"),
		mcp.WithString("database_id",
			mcp.Required(),
			mcp.Description("ID of the database to delete"),
		),
	)
	s.AddTool(deleteDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok {
			return nil, errors.New("database_id must be a string")
		}
		_, err := client.CloudDatabase.Instances().Delete(ctx, databaseID, &gobizfly.CloudDatabaseDelete{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete database: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Database %s deleted successfully", databaseID)), nil
	})
} 