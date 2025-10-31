package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterVolumeTools registers all volume-related tools with the MCP server
func RegisterVolumeTools(s *server.MCPServer, client *gobizfly.Client) {
	// List volumes tool
	listVolumesTool := mcp.NewTool("bizflycloud_list_volumes",
		mcp.WithDescription("List all Bizfly Cloud volumes"),
	)
	s.AddTool(listVolumesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumes, err := client.CloudServer.Volumes().List(ctx, &gobizfly.VolumeListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list volumes: %v", err)), nil
		}

		result := "Available volumes:\n\n"
		for _, volume := range volumes {
			result += fmt.Sprintf("Volume: %s\n", volume.Name)
			result += fmt.Sprintf("  ID: %s\n", volume.ID)
			result += fmt.Sprintf("  Status: %s\n", volume.Status)
			result += fmt.Sprintf("  Size: %d GB\n", volume.Size)
			result += fmt.Sprintf("  Type: %s\n", volume.VolumeType)
			result += fmt.Sprintf("  Zone: %s\n", volume.AvailabilityZone)
			result += fmt.Sprintf("  Created At: %s\n", volume.CreatedAt)
			result += fmt.Sprintf("  Updated At: %s\n", volume.UpdatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create volume tool
	createVolumeTool := mcp.NewTool("bizflycloud_create_volume",
		mcp.WithDescription("Create a new Bizfly Cloud volume"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the volume"),
		),
		mcp.WithNumber("size",
			mcp.Required(),
			mcp.Description("Size of the volume in GB"),
		),
		mcp.WithString("volume_type",
			mcp.Required(),
			mcp.Description("Type of the volume"),
		),
	)
	s.AddTool(createVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		size, ok := request.Params.Arguments["size"].(float64)
		if !ok {
			return nil, errors.New("size must be a number")
		}
		volumeType, ok := request.Params.Arguments["volume_type"].(string)
		if !ok {
			return nil, errors.New("volume_type must be a string")
		}

		volume, err := client.CloudServer.Volumes().Create(ctx, &gobizfly.VolumeCreateRequest{
			Name:       name,
			Size:       int(size),
			VolumeType: volumeType,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create volume: %v", err)), nil
		}

		result := fmt.Sprintf("Volume created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", volume.Name)
		result += fmt.Sprintf("  ID: %s\n", volume.ID)
		result += fmt.Sprintf("  Size: %d GB\n", volume.Size)
		result += fmt.Sprintf("  Type: %s\n", volume.VolumeType)
		return mcp.NewToolResultText(result), nil
	})

	// Resize volume tool
	resizeVolumeTool := mcp.NewTool("bizflycloud_resize_volume",
		mcp.WithDescription("Resize a Bizfly Cloud volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to resize"),
		),
		mcp.WithNumber("new_size",
			mcp.Required(),
			mcp.Description("New size of the volume in GB"),
		),
	)
	s.AddTool(resizeVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		newSize, ok := request.Params.Arguments["new_size"].(float64)
		if !ok {
			return nil, errors.New("new_size must be a number")
		}

		_, err := client.CloudServer.Volumes().ExtendVolume(ctx, volumeID, int(newSize))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resize volume: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Volume %s resized to %d GB successfully", volumeID, int(newSize))), nil
	})

	// Delete volume tool
	deleteVolumeTool := mcp.NewTool("bizflycloud_delete_volume",
		mcp.WithDescription("Delete a Bizfly Cloud volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to delete"),
		),
	)
	s.AddTool(deleteVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		err := client.CloudServer.Volumes().Delete(ctx, volumeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete volume: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Volume %s deleted successfully", volumeID)), nil
	})

	// List snapshots tool
	listSnapshotsTool := mcp.NewTool("bizflycloud_list_snapshots",
		mcp.WithDescription("List all Bizfly Cloud volume snapshots"),
	)
	s.AddTool(listSnapshotsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := &gobizfly.ListSnasphotsOptions{}
		snapshots, err := client.CloudServer.Snapshots().List(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list snapshots: %v", err)), nil
		}

		result := "Available snapshots:\n\n"
		for _, snapshot := range snapshots {
			result += fmt.Sprintf("Snapshot: %s\n", snapshot.Name)
			result += fmt.Sprintf("  ID: %s\n", snapshot.ID)
			result += fmt.Sprintf("  Status: %s\n", snapshot.Status)
			result += fmt.Sprintf("  Volume ID: %s\n", snapshot.VolumeID)
			result += fmt.Sprintf("  Size: %d GB\n", snapshot.Size)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create snapshot tool
	createSnapshotTool := mcp.NewTool("bizflycloud_create_snapshot",
		mcp.WithDescription("Create a snapshot of a Bizfly Cloud volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to snapshot"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the snapshot"),
		),
	)
	s.AddTool(createSnapshotTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}

		snapshot, err := client.CloudServer.Snapshots().Create(ctx, &gobizfly.SnapshotCreateRequest{
			VolumeID: volumeID,
			Name:     name,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create snapshot: %v", err)), nil
		}

		result := fmt.Sprintf("Snapshot created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", snapshot.Name)
		result += fmt.Sprintf("  ID: %s\n", snapshot.ID)
		result += fmt.Sprintf("  Status: %s\n", snapshot.Status)
		result += fmt.Sprintf("  Volume ID: %s\n", snapshot.VolumeID)
		result += fmt.Sprintf("  Size: %d GB\n", snapshot.Size)
		return mcp.NewToolResultText(result), nil
	})

	// Delete snapshot tool
	deleteSnapshotTool := mcp.NewTool("bizflycloud_delete_snapshot",
		mcp.WithDescription("Delete a Bizfly Cloud volume snapshot"),
		mcp.WithString("snapshot_id",
			mcp.Required(),
			mcp.Description("ID of the snapshot to delete"),
		),
	)
	s.AddTool(deleteSnapshotTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		snapshotID, ok := request.Params.Arguments["snapshot_id"].(string)
		if !ok {
			return nil, errors.New("snapshot_id must be a string")
		}
		err := client.CloudServer.Snapshots().Delete(ctx, snapshotID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete snapshot: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Snapshot %s deleted successfully", snapshotID)), nil
	})

	// Get volume tool
	getVolumeTool := mcp.NewTool("bizflycloud_get_volume",
		mcp.WithDescription("Get details of a Bizfly Cloud volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to get details for"),
		),
	)
	s.AddTool(getVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		volume, err := client.CloudServer.Volumes().Get(ctx, volumeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get volume: %v", err)), nil
		}

		result := fmt.Sprintf("Volume Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", volume.Name)
		result += fmt.Sprintf("ID: %s\n", volume.ID)
		result += fmt.Sprintf("Status: %s\n", volume.Status)
		result += fmt.Sprintf("Size: %d GB\n", volume.Size)
		result += fmt.Sprintf("Type: %s\n", volume.VolumeType)
		result += fmt.Sprintf("Zone: %s\n", volume.AvailabilityZone)
		result += fmt.Sprintf("Category: %s\n", volume.Category)
		if len(volume.Attachments) > 0 {
			result += fmt.Sprintf("Attachments:\n")
			for _, attachment := range volume.Attachments {
				result += fmt.Sprintf("  - Server ID: %s, Device: %s\n", attachment.ServerID, attachment.Device)
			}
		}
		result += fmt.Sprintf("Created At: %s\n", volume.CreatedAt)
		result += fmt.Sprintf("Updated At: %s\n", volume.UpdatedAt)
		return mcp.NewToolResultText(result), nil
	})

	// Attach volume tool
	attachVolumeTool := mcp.NewTool("bizflycloud_attach_volume",
		mcp.WithDescription("Attach a Bizfly Cloud volume to a server"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to attach"),
		),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to attach the volume to"),
		),
	)
	s.AddTool(attachVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.Volumes().Attach(ctx, volumeID, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to attach volume: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Volume %s attached to server %s successfully", volumeID, serverID)), nil
	})

	// Detach volume tool
	detachVolumeTool := mcp.NewTool("bizflycloud_detach_volume",
		mcp.WithDescription("Detach a Bizfly Cloud volume from a server"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to detach"),
		),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("ID of the server to detach the volume from"),
		),
	)
	s.AddTool(detachVolumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok {
			return nil, errors.New("volume_id must be a string")
		}
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			return nil, errors.New("server_id must be a string")
		}
		_, err := client.CloudServer.Volumes().Detach(ctx, volumeID, serverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to detach volume: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Volume %s detached from server %s successfully", volumeID, serverID)), nil
	})
} 