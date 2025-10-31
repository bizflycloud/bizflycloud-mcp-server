package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAutoScalingTools registers all AutoScaling-related tools with the MCP server
func RegisterAutoScalingTools(s *server.MCPServer, client *gobizfly.Client) {
	// List auto scaling groups tool
	listGroupsTool := mcp.NewTool("bizflycloud_list_autoscaling_groups",
		mcp.WithDescription("List all Bizfly Cloud AutoScaling groups"),
		mcp.WithBoolean("all",
			mcp.Description("List all groups including deleted ones (default: false)"),
		),
	)
	s.AddTool(listGroupsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		all, _ := request.Params.Arguments["all"].(bool)

		groups, err := client.AutoScaling.AutoScalingGroups().List(ctx, all)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list auto scaling groups: %v", err)), nil
		}

		result := "Available AutoScaling groups:\n\n"
		for _, group := range groups {
			result += fmt.Sprintf("Group: %s\n", group.Name)
			result += fmt.Sprintf("  ID: %s\n", group.ID)
			result += fmt.Sprintf("  Status: %s\n", group.Status)
			result += fmt.Sprintf("  Min Size: %d\n", group.MinSize)
			result += fmt.Sprintf("  Max Size: %d\n", group.MaxSize)
			result += fmt.Sprintf("  Desired Capacity: %d\n", group.DesiredCapacity)
			result += fmt.Sprintf("  Current Nodes: %d\n", len(group.NodeIDs))
			result += fmt.Sprintf("  Profile ID: %s\n", group.ProfileID)
			result += fmt.Sprintf("  Created At: %s\n", group.Created)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get auto scaling group tool
	getGroupTool := mcp.NewTool("bizflycloud_get_autoscaling_group",
		mcp.WithDescription("Get details of a Bizfly Cloud AutoScaling group"),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("ID of the auto scaling group"),
		),
	)
	s.AddTool(getGroupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		groupID, ok := request.Params.Arguments["group_id"].(string)
		if !ok {
			return nil, errors.New("group_id must be a string")
		}
		group, err := client.AutoScaling.AutoScalingGroups().Get(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get auto scaling group: %v", err)), nil
		}

		result := fmt.Sprintf("AutoScaling Group Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", group.Name)
		result += fmt.Sprintf("ID: %s\n", group.ID)
		result += fmt.Sprintf("Status: %s\n", group.Status)
		result += fmt.Sprintf("Min Size: %d\n", group.MinSize)
		result += fmt.Sprintf("Max Size: %d\n", group.MaxSize)
		result += fmt.Sprintf("Desired Capacity: %d\n", group.DesiredCapacity)
		result += fmt.Sprintf("Current Nodes: %d\n", len(group.NodeIDs))
		result += fmt.Sprintf("Profile ID: %s\n", group.ProfileID)
		result += fmt.Sprintf("Profile Name: %s\n", group.ProfileName)
		result += fmt.Sprintf("Created At: %s\n", group.Created)
		result += fmt.Sprintf("Updated At: %s\n", group.Updated)
		return mcp.NewToolResultText(result), nil
	})

	// Create auto scaling group tool
	createGroupTool := mcp.NewTool("bizflycloud_create_autoscaling_group",
		mcp.WithDescription("Create a new Bizfly Cloud AutoScaling group"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the auto scaling group"),
		),
		mcp.WithString("profile_id",
			mcp.Required(),
			mcp.Description("ID of the launch configuration profile"),
		),
		mcp.WithNumber("min_size",
			mcp.Required(),
			mcp.Description("Minimum number of nodes"),
		),
		mcp.WithNumber("max_size",
			mcp.Required(),
			mcp.Description("Maximum number of nodes"),
		),
		mcp.WithNumber("desired_capacity",
			mcp.Required(),
			mcp.Description("Desired number of nodes"),
		),
	)
	s.AddTool(createGroupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		profileID, ok := request.Params.Arguments["profile_id"].(string)
		if !ok {
			return nil, errors.New("profile_id must be a string")
		}
		minSize, ok := request.Params.Arguments["min_size"].(float64)
		if !ok {
			return nil, errors.New("min_size must be a number")
		}
		maxSize, ok := request.Params.Arguments["max_size"].(float64)
		if !ok {
			return nil, errors.New("max_size must be a number")
		}
		desiredCapacity, ok := request.Params.Arguments["desired_capacity"].(float64)
		if !ok {
			return nil, errors.New("desired_capacity must be a number")
		}

		group, err := client.AutoScaling.AutoScalingGroups().Create(ctx, &gobizfly.AutoScalingGroupCreateRequest{
			Name:           name,
			ProfileID:      profileID,
			MinSize:        int(minSize),
			MaxSize:        int(maxSize),
			DesiredCapacity: int(desiredCapacity),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create auto scaling group: %v", err)), nil
		}

		result := fmt.Sprintf("AutoScaling group created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", group.Name)
		result += fmt.Sprintf("  ID: %s\n", group.ID)
		result += fmt.Sprintf("  Status: %s\n", group.Status)
		result += fmt.Sprintf("  Min Size: %d\n", group.MinSize)
		result += fmt.Sprintf("  Max Size: %d\n", group.MaxSize)
		result += fmt.Sprintf("  Desired Capacity: %d\n", group.DesiredCapacity)
		return mcp.NewToolResultText(result), nil
	})

	// Delete auto scaling group tool
	deleteGroupTool := mcp.NewTool("bizflycloud_delete_autoscaling_group",
		mcp.WithDescription("Delete a Bizfly Cloud AutoScaling group"),
		mcp.WithString("group_id",
			mcp.Required(),
			mcp.Description("ID of the auto scaling group to delete"),
		),
	)
	s.AddTool(deleteGroupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		groupID, ok := request.Params.Arguments["group_id"].(string)
		if !ok {
			return nil, errors.New("group_id must be a string")
		}
		err := client.AutoScaling.AutoScalingGroups().Delete(ctx, groupID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete auto scaling group: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("AutoScaling group %s deleted successfully", groupID)), nil
	})
}

