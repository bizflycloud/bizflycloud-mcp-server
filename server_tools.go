package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterServerTools registers all server-related tools with the MCP server
func RegisterServerTools(s *server.MCPServer, client *gobizfly.Client) {
	// List servers tool
	listServersTool := mcp.NewTool("bizflycloud_list_servers",
		mcp.WithDescription("List all Bizfly Cloud servers"),
	)
	s.AddTool(listServersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		servers, err := client.CloudServer.List(ctx, &gobizfly.ServerListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list servers: %v", err)), nil
		}

		result := "Available servers:\n\n"
		for _, server := range servers {
			result += fmt.Sprintf("Server: %s\n", server.Name)
			result += fmt.Sprintf("  ID: %s\n", server.ID)
			result += fmt.Sprintf("  Status: %s\n", server.Status)
			result += fmt.Sprintf("  Flavor: %s\n", server.FlavorName)
			result += fmt.Sprintf("  Zone: %s\n", server.AvailabilityZone)
			if len(server.IPAddresses.WanV4Addresses) > 0 {
				result += fmt.Sprintf("  WAN IP: %s\n", string(server.IPAddresses.WanV4Addresses[0].Address))
			}
			if len(server.IPAddresses.LanAddresses) > 0 {
				result += fmt.Sprintf("  LAN IP: %s\n", string(server.IPAddresses.LanAddresses[0].Address))
			}
			result += fmt.Sprintf("  Created At: %s\n", server.CreatedAt)
			result += fmt.Sprintf("  Updated At: %s\n", server.UpdatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Reboot server tool
	rebootServerTool := mcp.NewTool("bizflycloud_reboot_server",
		mcp.WithDescription("Reboot a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to reboot"),
		),
	)
	s.AddTool(rebootServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.SoftReboot(ctx, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to reboot server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s rebooted successfully", serverID)), nil
	})

	// Delete server tool
	deleteServerTool := mcp.NewTool("bizflycloud_delete_server",
		mcp.WithDescription("Delete a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to delete"),
		),
	)
	s.AddTool(deleteServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.Delete(ctx, serverID, []string{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s deleted successfully", serverID)), nil
	})

	// Start server tool
	startServerTool := mcp.NewTool("bizflycloud_start_server",
		mcp.WithDescription("Start a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to start"),
		),
	)
	s.AddTool(startServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.Start(ctx, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s started successfully", serverID)), nil
	})

	// Resize server tool
	resizeServerTool := mcp.NewTool("bizflycloud_resize_server",
		mcp.WithDescription("Resize a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to resize"),
		),
		mcp.WithString("flavor_name",
			mcp.Required(),
			mcp.Description("Name of the new flavor to resize to"),
		),
	)
	s.AddTool(resizeServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		flavorName, ok := request.Params.Arguments["flavor_name"].(string)
		if !ok {
			return nil, errors.New("flavor_name must be a string")
		}

		// Get the flavor ID from the name
		flavors, err := client.CloudServer.Flavors().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get flavors: %v", err)), nil
		}

		var flavorID string
		for _, flavor := range flavors {
			if flavor.Name == flavorName {
				flavorID = flavor.ID
				break
			}
		}
		if flavorID == "" {
			return mcp.NewToolResultError(fmt.Sprintf("Flavor '%s' not found", flavorName)), nil
		}

		_, err = client.CloudServer.Resize(ctx, serverID, flavorID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resize server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s resizing to flavor %s successfully", serverID, flavorName)), nil
	})

	// List flavors tool
	listFlavorsTool := mcp.NewTool("bizflycloud_list_flavors",
		mcp.WithDescription("List all available Bizfly Cloud server flavors"),
	)
	s.AddTool(listFlavorsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		flavors, err := client.CloudServer.Flavors().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list flavors: %v", err)), nil
		}

		result := "Available flavors:\n\n"
		for _, flavor := range flavors {
			result += fmt.Sprintf("Flavor: %s\n", flavor.Name)
			result += fmt.Sprintf("  ID: %s\n", flavor.ID)
			result += fmt.Sprintf("  vCPUs: %d\n", flavor.VCPUs)
			result += fmt.Sprintf("  RAM: %d MB\n", flavor.RAM)
			result += fmt.Sprintf("  Disk: %d GB\n", flavor.Disk)
			result += fmt.Sprintf("  Category: %s\n", flavor.Category)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get server tool
	getServerTool := mcp.NewTool("bizflycloud_get_server",
		mcp.WithDescription("Get details of a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to get details for"),
		),
	)
	s.AddTool(getServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		server, err := client.CloudServer.Get(ctx, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get server: %v", err)), nil
		}

		result := fmt.Sprintf("Server Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", server.Name)
		result += fmt.Sprintf("ID: %s\n", server.ID)
		result += fmt.Sprintf("Status: %s\n", server.Status)
		result += fmt.Sprintf("Flavor: %s\n", server.FlavorName)
		result += fmt.Sprintf("Zone: %s\n", server.AvailabilityZone)
		if len(server.IPAddresses.WanV4Addresses) > 0 {
			result += fmt.Sprintf("WAN IPs:\n")
			for _, ip := range server.IPAddresses.WanV4Addresses {
				result += fmt.Sprintf("  - %s\n", string(ip.Address))
			}
		}
		if len(server.IPAddresses.LanAddresses) > 0 {
			result += fmt.Sprintf("LAN IPs:\n")
			for _, ip := range server.IPAddresses.LanAddresses {
				result += fmt.Sprintf("  - %s\n", string(ip.Address))
			}
		}
		result += fmt.Sprintf("Created At: %s\n", server.CreatedAt)
		result += fmt.Sprintf("Updated At: %s\n", server.UpdatedAt)
		return mcp.NewToolResultText(result), nil
	})

	// Stop server tool
	stopServerTool := mcp.NewTool("bizflycloud_stop_server",
		mcp.WithDescription("Stop a Bizfly Cloud server"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to stop"),
		),
	)
	s.AddTool(stopServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.Stop(ctx, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s stopped successfully", serverID)), nil
	})

	// Hard reboot server tool
	hardRebootServerTool := mcp.NewTool("bizflycloud_hard_reboot_server",
		mcp.WithDescription("Hard reboot a Bizfly Cloud server (force reboot)"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to hard reboot"),
		),
	)
	s.AddTool(hardRebootServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.HardReboot(ctx, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to hard reboot server: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Server %s hard rebooted successfully", serverID)), nil
	})
} 