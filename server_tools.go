package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	// Create server tool - Create a server with customizable OS, flavor, disk size and volume type
	createServerTool := mcp.NewTool("bizflycloud_create_server",
		mcp.WithDescription("Create a new Bizfly Cloud server"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the server"),
		),
		mcp.WithString("os_type",
			mcp.Description("OS type (ubuntu, centos, etc.) - optional, defaults to ubuntu"),
		),
		mcp.WithString("image_id",
			mcp.Description("ID of the image (optional, will auto-select based on os_type if not provided)"),
		),
		mcp.WithString("flavor_name",
			mcp.Description("Name of the flavor (optional, defaults to nix.1c_1g for smallest config)"),
		),
		mcp.WithNumber("root_disk_size",
			mcp.Description("Root disk size in GB (optional, defaults to 20 GB)"),
		),
		mcp.WithString("volume_type",
			mcp.Description("Volume type for root disk (optional, defaults to SSD - PREMIUM-SSD1)"),
		),
		mcp.WithString("availability_zone",
			mcp.Description("Availability zone (optional, defaults to HN1)"),
		),
		mcp.WithString("use_password",
			mcp.Description("Set to 'true' to use password authentication (optional, defaults to SSH key)"),
		),
	)
	s.AddTool(createServerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}

		// Get OS type (default to ubuntu)
		osType := "ubuntu"
		if ost, ok := request.Params.Arguments["os_type"].(string); ok && ost != "" {
			osType = strings.ToLower(ost)
		}

		flavorName := "nix.1c_1g" // Default to smallest flavor
		if fn, ok := request.Params.Arguments["flavor_name"].(string); ok && fn != "" {
			flavorName = fn
		}

		// Get root disk size (default to 20GB)
		rootDiskSize := 20
		if rds, ok := request.Params.Arguments["root_disk_size"].(float64); ok && rds > 0 {
			rootDiskSize = int(rds)
		}

		availabilityZone := "HN1"
		if zone, ok := request.Params.Arguments["availability_zone"].(string); ok && zone != "" {
			availabilityZone = zone
		}

		// Verify flavor exists
		flavors, err := client.CloudServer.Flavors().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get flavors: %v", err)), nil
		}

		flavorFound := false
		for _, flavor := range flavors {
			if flavor.Name == flavorName {
				flavorFound = true
				break
			}
		}
		if !flavorFound {
			return mcp.NewToolResultError(fmt.Sprintf("Flavor '%s' not found. Use bizflycloud_list_flavors to see available flavors", flavorName)), nil
		}

		// Get image ID
		var imageID string
		if imgID, ok := request.Params.Arguments["image_id"].(string); ok && imgID != "" {
			imageID = imgID
		} else {
			// Try to find image from custom images first
			customImages, err := client.CloudServer.CustomImages().List(ctx)
			if err == nil && len(customImages) > 0 {
				// Look for OS type in custom images
				for _, img := range customImages {
					if strings.Contains(strings.ToLower(img.Name), osType) {
						imageID = img.ID
						break
					}
				}
			}
			
			// If not found in custom images, try OS images
			if imageID == "" {
				images, err := client.CloudServer.OSImages().List(ctx)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to get images: %v. Please provide image_id manually", err)), nil
				}

				// Find image matching OS type
				for _, image := range images {
					if strings.ToLower(image.OSDistribution) == osType {
						// Get first version's ID
						if len(image.Version) > 0 {
							imageID = image.Version[0].ID
							break
						}
					}
				}
			}
			
			if imageID == "" {
				return mcp.NewToolResultError(fmt.Sprintf("%s image not found automatically. Please provide image_id parameter", strings.Title(osType))), nil
			}
		}

		// Determine password authentication
		usePassword := false
		if pwd, ok := request.Params.Arguments["use_password"].(string); ok && pwd == "true" {
			usePassword = true
		}

		// Create server request
		// Determine server type based on flavor category
		serverType := "premium" // Default to premium
		for _, flavor := range flavors {
			if flavor.Name == flavorName {
				// Map flavor category to server type
				switch flavor.Category {
				case "basic":
					serverType = "basic"
				case "enterprise":
					serverType = "enterprise"
				case "dedicated":
					serverType = "dedicated"
				case "premium":
					serverType = "premium"
				case "vps":
					serverType = "premium"
				default:
					serverType = "premium"
				}
				break
			}
		}
		
		// Ensure serverType is not empty
		if serverType == "" {
			serverType = "premium"
		}

		// Get volume type for root disk
		volumeType := ""
		if vt, ok := request.Params.Arguments["volume_type"].(string); ok && vt != "" {
			volumeType = vt
		} else {
			// Try to find SSD volume type from existing volumes
			volumes, err := client.CloudServer.Volumes().List(ctx, &gobizfly.VolumeListOptions{})
			if err == nil {
				// Look for SSD volume types in existing volumes (highest priority)
				for _, vol := range volumes {
					if strings.Contains(strings.ToUpper(vol.VolumeType), "SSD") {
						volumeType = vol.VolumeType
						break
					}
				}
				// If no SSD found, try NVME (also fast storage, second priority)
				if volumeType == "" {
					for _, vol := range volumes {
						if strings.Contains(strings.ToUpper(vol.VolumeType), "NVME") {
							volumeType = vol.VolumeType
							break
						}
					}
				}
			}
			
			// If no SSD/NVME found, use default SSD types in order of preference
			if volumeType == "" {
				possibleSSDTypes := []string{"PREMIUM-SSD1", "PREMIUM-NVME1", "SSD", "BASIC-SSD1"}
				volumeType = possibleSSDTypes[0] // Default to PREMIUM-SSD1 (SSD)
			}
		}
		
		rootDisk := &gobizfly.ServerDisk{
			Size:       rootDiskSize,
			VolumeType: &volumeType,
		}

		// OS Type should be "image" when using an image ID
		createReq := &gobizfly.ServerCreateRequest{
			Name:             name,
			FlavorName:       flavorName,
			Type:             serverType,
			RootDisk:         rootDisk,
			AvailabilityZone: availabilityZone,
			OS: &gobizfly.ServerOS{
				ID:   imageID,
				Type: "image", // OS type: "image", "snapshot", "volume", "volume_source", "prebuild_app"
			},
			Password: usePassword,
		}

		// Create the server
		createResp, err := client.CloudServer.Create(ctx, createReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create server: %v", err)), nil
		}

		result := fmt.Sprintf("Server creation initiated successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", name)
		result += fmt.Sprintf("  Flavor: %s\n", flavorName)
		result += fmt.Sprintf("  OS: %s\n", strings.Title(osType))
		result += fmt.Sprintf("  Root Disk: %d GB (%s)\n", rootDiskSize, volumeType)
		result += fmt.Sprintf("  Zone: %s\n", availabilityZone)
		result += fmt.Sprintf("  Task IDs: %v\n", createResp.Task)
		result += fmt.Sprintf("\nNote: Server is being created. Use bizflycloud_list_servers to check status.\n")
		return mcp.NewToolResultText(result), nil
	})
} 