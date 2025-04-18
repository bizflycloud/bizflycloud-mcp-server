package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterLoadBalancerTools registers all load balancer-related tools with the MCP server
func RegisterLoadBalancerTools(s *server.MCPServer, client *gobizfly.Client) {
	// List load balancers tool
	listLoadBalancersTool := mcp.NewTool("bizflycloud_list_loadbalancers",
		mcp.WithDescription("List all Bizfly Cloud load balancers"),
	)
	s.AddTool(listLoadBalancersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		loadbalancers, err := client.CloudLoadBalancer.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list load balancers: %v", err)), nil
		}

		result := "Available load balancers:\n\n"
		for _, lb := range loadbalancers {
			result += fmt.Sprintf("Load Balancer: %s\n", lb.Name)
			result += fmt.Sprintf("  ID: %s\n", lb.ID)
			result += fmt.Sprintf("  Provider Status: %s\n", lb.ProvisioningStatus)
			result += fmt.Sprintf("  Operating Status: %s\n", lb.OperatingStatus)
			result += fmt.Sprintf("  Type: %s\n", lb.Type)
			result += fmt.Sprintf("  Network Type: %s\n", lb.NetworkType)
			result += fmt.Sprintf("  Created At: %s\n", lb.CreatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create load balancer tool
	createLoadBalancerTool := mcp.NewTool("bizflycloud_create_loadbalancer",
		mcp.WithDescription("Create a new Bizfly Cloud load balancer"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the load balancer"),
		),
		mcp.WithString("network_type",
			mcp.Required(),
			mcp.Description("Network type (external, internal)"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of load balancer"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the load balancer"),
		),
	)
	s.AddTool(createLoadBalancerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		networkType, ok := request.Params.Arguments["network_type"].(string)
		if !ok {
			return nil, errors.New("network_type must be a string")
		}
		lbType, ok := request.Params.Arguments["type"].(string)
		if !ok {
			return nil, errors.New("type must be a string")
		}
		description, _ := request.Params.Arguments["description"].(string)

		loadbalancer, err := client.CloudLoadBalancer.Create(ctx, &gobizfly.LoadBalancerCreateRequest{
			Name:        name,
			NetworkType: networkType,
			Type:        lbType,
			Description: description,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create load balancer: %v", err)), nil
		}

		result := fmt.Sprintf("Load balancer created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", loadbalancer.Name)
		result += fmt.Sprintf("  ID: %s\n", loadbalancer.ID)
		result += fmt.Sprintf("  Provider Status: %s\n", loadbalancer.ProvisioningStatus)
		result += fmt.Sprintf("  Operating Status: %s\n", loadbalancer.OperatingStatus)
		result += fmt.Sprintf("  Type: %s\n", loadbalancer.Type)
		result += fmt.Sprintf("  Network Type: %s\n", loadbalancer.NetworkType)
		return mcp.NewToolResultText(result), nil
	})

	// Delete load balancer tool
	deleteLoadBalancerTool := mcp.NewTool("bizflycloud_delete_loadbalancer",
		mcp.WithDescription("Delete a Bizfly Cloud load balancer"),
		mcp.WithString("loadbalancer_id",
			mcp.Required(),
			mcp.Description("ID of the load balancer to delete"),
		),
	)
	s.AddTool(deleteLoadBalancerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		loadbalancerID, ok := request.Params.Arguments["loadbalancer_id"].(string)
		if !ok {
			return nil, errors.New("loadbalancer_id must be a string")
		}
		err := client.CloudLoadBalancer.Delete(ctx, &gobizfly.LoadBalancerDeleteRequest{
			ID: loadbalancerID,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete load balancer: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Load balancer %s deleted successfully", loadbalancerID)), nil
	})
} 