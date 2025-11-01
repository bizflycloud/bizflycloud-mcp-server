package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// safeString returns empty string if s is nil, otherwise returns the string
func safeString(s string) string {
	return s
}

// RegisterDatabaseTools registers all database-related tools with the MCP server
func RegisterDatabaseTools(s *server.MCPServer, client *gobizfly.Client) {
	// List databases tool
	listDatabasesTool := mcp.NewTool("bizflycloud_list_databases",
		mcp.WithDescription("List all Bizfly Cloud databases"),
	)
	s.AddTool(listDatabasesTool, func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
		// Panic recovery - return error result on panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Recovered from panic in listDatabasesTool: %v", r)
				errMsg := fmt.Sprintf("Panic occurred while listing databases: %v", r)
				result = mcp.NewToolResultError(errMsg)
				err = nil // Don't propagate panic error
			}
		}()

		log.Printf("[DEBUG] Database List tool called")
		
		// Check if client is nil
		if client == nil {
			log.Printf("[ERROR] Client is nil")
			return mcp.NewToolResultText("Available databases:\n\n(Database service is not available - client is nil)"), nil
		}

		// Check if CloudDatabase service exists
		if client.CloudDatabase == nil {
			log.Printf("[ERROR] CloudDatabase is nil")
			return mcp.NewToolResultText("Available databases:\n\n(Database service is not available)"), nil
		}

		// Check if Instances() is available
		if client.CloudDatabase.Instances() == nil {
			log.Printf("[ERROR] CloudDatabase.Instances() is nil")
			return mcp.NewToolResultText("Available databases:\n\n(Database service is not available)"), nil
		}

		// Get databases - call List() with empty struct to avoid nil pointer dereference in AddParamsListOption
		databases, err := client.CloudDatabase.Instances().List(ctx, &gobizfly.CloudDatabaseListOption{})
		
		if err != nil {
			log.Printf("[ERROR] Failed to list databases: %v", err)
			// Check if error is 404 or service not available
			errStr := strings.ToLower(err.Error())
			if strings.Contains(errStr, "404") ||
				strings.Contains(errStr, "not found") ||
				strings.Contains(errStr, "resource not found") ||
				strings.Contains(errStr, "<svg") ||
				strings.Contains(errStr, "<html") {
				return mcp.NewToolResultText("Available databases:\n\n(No databases found or Database service is not enabled)"), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list databases: %v", err)), nil
		}

		// Check if databases is nil
		if databases == nil {
			log.Printf("[WARN] databases is nil after List() call")
			return mcp.NewToolResultText("Available databases:\n\n(No databases found)"), nil
		}

		resultText := "Available databases:\n\n"
		if len(databases) == 0 {
			resultText += "(No databases found)\n"
		} else {
			for _, db := range databases {
				if db == nil {
					log.Printf("[WARN] Found nil database in list")
					continue
				}
				
				// Safely access db fields
				resultText += fmt.Sprintf("Database: %s\n", safeString(db.Name))
				resultText += fmt.Sprintf("  ID: %s\n", safeString(db.ID))
				
				if db.Datastore.Type != "" {
					resultText += fmt.Sprintf("  DataStore Type: %s\n", db.Datastore.Type)
				}
				if db.Datastore.ID != "" {
					resultText += fmt.Sprintf("  DataStore ID: %s\n", db.Datastore.ID)
				}
				
				resultText += fmt.Sprintf("  Status: %s\n", safeString(db.Status))
				resultText += fmt.Sprintf("  Created At: %s\n", safeString(db.CreatedAt))
				resultText += "\n"
			}
		}
		return mcp.NewToolResultText(resultText), nil
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

	// Get database tool
	getDatabaseTool := mcp.NewTool("bizflycloud_get_database",
		mcp.WithDescription("Get details of a Bizfly Cloud database instance"),
		mcp.WithString("database_id",
			mcp.Required(),
			mcp.Description("ID of the database to get details for"),
		),
	)
	s.AddTool(getDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok {
			return nil, errors.New("database_id must be a string")
		}
		db, err := client.CloudDatabase.Instances().Get(ctx, databaseID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get database: %v", err)), nil
		}

		result := fmt.Sprintf("Database Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", db.Name)
		result += fmt.Sprintf("ID: %s\n", db.ID)
		result += fmt.Sprintf("Status: %s\n", db.Status)
		result += fmt.Sprintf("DataStore Type: %s\n", db.Datastore.Type)
		result += fmt.Sprintf("DataStore ID: %s\n", db.Datastore.ID)
		result += fmt.Sprintf("Instance Type: %s\n", db.InstanceType)
		result += fmt.Sprintf("Volume Size: %d GB\n", db.Volume.Size)
		result += fmt.Sprintf("Volume Used: %.2f GB\n", db.Volume.Used)
		result += fmt.Sprintf("Public Access: %v\n", db.PublicAccess)
		result += fmt.Sprintf("Enable Failover: %v\n", db.EnableFailover)
		result += fmt.Sprintf("Nodes Count: %d\n", len(db.Nodes))
		
		// Get detailed nodes information using ListNodes
		nodes, nodesErr := client.CloudDatabase.Instances().ListNodes(ctx, databaseID, &gobizfly.CloudDatabaseListOption{})
		if nodesErr == nil && len(nodes) > 0 {
			result += fmt.Sprintf("\nNodes Details:\n")
			for i, node := range nodes {
				result += fmt.Sprintf("  Node %d:\n", i+1)
				result += fmt.Sprintf("    ID: %s\n", node.ID)
				result += fmt.Sprintf("    Name: %s\n", node.Name)
				result += fmt.Sprintf("    Status: %s\n", node.Status)
				result += fmt.Sprintf("    Operating Status: %s\n", node.OperatingStatus)
				result += fmt.Sprintf("    Node Type: %s\n", node.NodeType)
				result += fmt.Sprintf("    Role: %s\n", node.Role)
				result += fmt.Sprintf("    Flavor: %s\n", node.Flavor)
				result += fmt.Sprintf("    Availability Zone: %s\n", node.AvailabilityZone)
				result += fmt.Sprintf("    Enable Failover: %v\n", node.EnableFailover)
				if node.ReplicaOf != "" {
					result += fmt.Sprintf("    Replica Of: %s\n", node.ReplicaOf)
				}
				if len(node.Addresses.Private) > 0 {
					result += fmt.Sprintf("    Private Addresses:\n")
					for _, addr := range node.Addresses.Private {
						result += fmt.Sprintf("      - IP: %s, Port: %d, Network: %s\n", addr.IPAddress, addr.Port, addr.Network)
					}
				}
				if len(node.Addresses.Public) > 0 {
					result += fmt.Sprintf("    Public Addresses:\n")
					for _, addr := range node.Addresses.Public {
						result += fmt.Sprintf("      - IP: %s, Port: %d, Network: %s\n", addr.IPAddress, addr.Port, addr.Network)
					}
				}
				if node.DNS.Private != "" {
					result += fmt.Sprintf("    Private DNS: %s\n", node.DNS.Private)
				}
				if node.DNS.Public != "" {
					result += fmt.Sprintf("    Public DNS: %s\n", node.DNS.Public)
				}
				if len(node.Replicas) > 0 {
					result += fmt.Sprintf("    Replicas: %d\n", len(node.Replicas))
				}
				result += fmt.Sprintf("    Created At: %s\n", node.CreatedAt)
				result += "\n"
			}
		} else if len(db.Nodes) > 0 {
			// Fallback to basic node info from Get() response
			result += fmt.Sprintf("\nNodes:\n")
			for _, node := range db.Nodes {
				result += fmt.Sprintf("  - Node ID: %s\n", node.ID)
			}
		}
		
		result += fmt.Sprintf("Created At: %s\n", db.CreatedAt)
		return mcp.NewToolResultText(result), nil
	})
	
	// List database nodes tool
	listNodesTool := mcp.NewTool("bizflycloud_list_database_nodes",
		mcp.WithDescription("List all nodes in a Bizfly Cloud database instance"),
		mcp.WithString("database_id",
			mcp.Required(),
			mcp.Description("ID of the database instance"),
		),
	)
	s.AddTool(listNodesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok {
			return nil, errors.New("database_id must be a string")
		}
		
		nodes, err := client.CloudDatabase.Instances().ListNodes(ctx, databaseID, &gobizfly.CloudDatabaseListOption{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list database nodes: %v", err)), nil
		}
		
		result := fmt.Sprintf("Database Nodes for Instance %s:\n\n", databaseID)
		if len(nodes) == 0 {
			result += "(No nodes found)\n"
		} else {
			for i, node := range nodes {
				result += fmt.Sprintf("Node %d:\n", i+1)
				result += fmt.Sprintf("  ID: %s\n", node.ID)
				result += fmt.Sprintf("  Name: %s\n", node.Name)
				result += fmt.Sprintf("  Status: %s\n", node.Status)
				result += fmt.Sprintf("  Operating Status: %s\n", node.OperatingStatus)
				result += fmt.Sprintf("  Node Type: %s\n", node.NodeType)
				result += fmt.Sprintf("  Role: %s\n", node.Role)
				result += fmt.Sprintf("  Flavor: %s\n", node.Flavor)
				result += fmt.Sprintf("  Availability Zone: %s\n", node.AvailabilityZone)
				result += fmt.Sprintf("  Enable Failover: %v\n", node.EnableFailover)
				if node.ReplicaOf != "" {
					result += fmt.Sprintf("  Replica Of: %s\n", node.ReplicaOf)
				}
				if len(node.Addresses.Private) > 0 {
					result += fmt.Sprintf("  Private Addresses:\n")
					for _, addr := range node.Addresses.Private {
						result += fmt.Sprintf("    - IP: %s, Port: %d, Network: %s\n", addr.IPAddress, addr.Port, addr.Network)
					}
				}
				if len(node.Addresses.Public) > 0 {
					result += fmt.Sprintf("  Public Addresses:\n")
					for _, addr := range node.Addresses.Public {
						result += fmt.Sprintf("    - IP: %s, Port: %d, Network: %s\n", addr.IPAddress, addr.Port, addr.Network)
					}
				}
				if node.DNS.Private != "" {
					result += fmt.Sprintf("  Private DNS: %s\n", node.DNS.Private)
				}
				if node.DNS.Public != "" {
					result += fmt.Sprintf("  Public DNS: %s\n", node.DNS.Public)
				}
				if node.DNS.SRV != "" {
					result += fmt.Sprintf("  SRV DNS: %s\n", node.DNS.SRV)
				}
				if len(node.Replicas) > 0 {
					result += fmt.Sprintf("  Replicas Count: %d\n", len(node.Replicas))
					for j, replica := range node.Replicas {
						result += fmt.Sprintf("    Replica %d: %s (%s)\n", j+1, replica.ID, replica.Name)
					}
				}
				result += fmt.Sprintf("  Created At: %s\n", node.CreatedAt)
				if node.Message != "" {
					result += fmt.Sprintf("  Message: %s\n", node.Message)
				}
				result += "\n"
			}
		}
		return mcp.NewToolResultText(result), nil
	})

	// List backups tool
	listBackupsTool := mcp.NewTool("bizflycloud_list_database_backups",
		mcp.WithDescription("List backups for a Bizfly Cloud database instance"),
		mcp.WithString("database_id",
			mcp.Required(),
			mcp.Description("ID of the database instance"),
		),
	)
	s.AddTool(listBackupsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok {
			return nil, errors.New("database_id must be a string")
		}

		resource := &gobizfly.CloudDatabaseBackupResource{
			ResourceID:   databaseID,
			ResourceType: "instance",
		}
		backups, err := client.CloudDatabase.Backups().List(ctx, resource, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list backups: %v", err)), nil
		}

		result := "Available backups:\n\n"
		for _, backup := range backups {
			result += fmt.Sprintf("Backup: %s\n", backup.Name)
			result += fmt.Sprintf("  ID: %s\n", backup.ID)
			result += fmt.Sprintf("  Status: %s\n", backup.Status)
			result += fmt.Sprintf("  Type: %s\n", backup.Type)
			result += fmt.Sprintf("  Size: %.2f GB\n", backup.Size)
			result += fmt.Sprintf("  Created At: %s\n", backup.Created)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create backup tool
	createBackupTool := mcp.NewTool("bizflycloud_create_database_backup",
		mcp.WithDescription("Create a backup for a Bizfly Cloud database instance"),
		mcp.WithString("database_id",
			mcp.Required(),
			mcp.Description("ID of the database instance"),
		),
		mcp.WithString("backup_name",
			mcp.Required(),
			mcp.Description("Name of the backup"),
		),
	)
	s.AddTool(createBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok {
			return nil, errors.New("database_id must be a string")
		}
		backupName, ok := request.Params.Arguments["backup_name"].(string)
		if !ok {
			return nil, errors.New("backup_name must be a string")
		}

		backup, err := client.CloudDatabase.Backups().Create(ctx, "instance", databaseID, &gobizfly.CloudDatabaseBackupCreate{
			Name: backupName,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create backup: %v", err)), nil
		}

		result := fmt.Sprintf("Backup created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", backup.Name)
		result += fmt.Sprintf("  ID: %s\n", backup.ID)
		result += fmt.Sprintf("  Status: %s\n", backup.Status)
		result += fmt.Sprintf("  Type: %s\n", backup.Type)
		return mcp.NewToolResultText(result), nil
	})
} 