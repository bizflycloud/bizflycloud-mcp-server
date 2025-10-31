package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDNSTools registers all DNS-related tools with the MCP server
func RegisterDNSTools(s *server.MCPServer, client *gobizfly.Client) {
	// List DNS zones tool
	listZonesTool := mcp.NewTool("bizflycloud_list_dns_zones",
		mcp.WithDescription("List all Bizfly Cloud DNS zones"),
	)
	s.AddTool(listZonesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zones, err := client.DNS.ListZones(ctx, &gobizfly.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list DNS zones: %v", err)), nil
		}

		result := "Available DNS zones:\n\n"
		for _, zone := range zones.Zones {
			result += fmt.Sprintf("Zone: %s\n", zone.Name)
			result += fmt.Sprintf("  ID: %s\n", zone.ID)
			result += fmt.Sprintf("  Active: %v\n", zone.Active)
			result += fmt.Sprintf("  TTL: %d\n", zone.TTL)
			if len(zone.NameServer) > 0 {
				result += fmt.Sprintf("  Name Servers: %v\n", zone.NameServer)
			}
			result += fmt.Sprintf("  Created At: %s\n", zone.CreatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create DNS zone tool
	createZoneTool := mcp.NewTool("bizflycloud_create_dns_zone",
		mcp.WithDescription("Create a new Bizfly Cloud DNS zone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the DNS zone (e.g., example.com)"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the DNS zone"),
		),
	)
	s.AddTool(createZoneTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		description, _ := request.Params.Arguments["description"].(string)

		zone, err := client.DNS.CreateZone(ctx, &gobizfly.CreateZonePayload{
			Name:        name,
			Description: description,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create DNS zone: %v", err)), nil
		}

		result := fmt.Sprintf("DNS zone created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", zone.Name)
		result += fmt.Sprintf("  ID: %s\n", zone.ID)
		result += fmt.Sprintf("  Active: %v\n", zone.Active)
		result += fmt.Sprintf("  TTL: %d\n", zone.TTL)
		return mcp.NewToolResultText(result), nil
	})

	// Get DNS zone tool
	getZoneTool := mcp.NewTool("bizflycloud_get_dns_zone",
		mcp.WithDescription("Get details of a Bizfly Cloud DNS zone"),
		mcp.WithString("zone_id",
			mcp.Required(),
			mcp.Description("ID of the DNS zone"),
		),
	)
	s.AddTool(getZoneTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, ok := request.Params.Arguments["zone_id"].(string)
		if !ok {
			return nil, errors.New("zone_id must be a string")
		}
		zone, err := client.DNS.GetZone(ctx, zoneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get DNS zone: %v", err)), nil
		}

		result := fmt.Sprintf("DNS Zone Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", zone.Name)
		result += fmt.Sprintf("ID: %s\n", zone.ID)
		result += fmt.Sprintf("Active: %v\n", zone.Active)
		result += fmt.Sprintf("TTL: %d\n", zone.TTL)
		if len(zone.NameServer) > 0 {
			result += fmt.Sprintf("Name Servers: %v\n", zone.NameServer)
		}
		result += fmt.Sprintf("Records Count: %d\n", len(zone.RecordsSet))
		result += fmt.Sprintf("Created At: %s\n", zone.CreatedAt)
		result += fmt.Sprintf("Updated At: %s\n", zone.UpdatedAt)
		return mcp.NewToolResultText(result), nil
	})

	// Delete DNS zone tool
	deleteZoneTool := mcp.NewTool("bizflycloud_delete_dns_zone",
		mcp.WithDescription("Delete a Bizfly Cloud DNS zone"),
		mcp.WithString("zone_id",
			mcp.Required(),
			mcp.Description("ID of the DNS zone to delete"),
		),
	)
	s.AddTool(deleteZoneTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, ok := request.Params.Arguments["zone_id"].(string)
		if !ok {
			return nil, errors.New("zone_id must be a string")
		}
		err := client.DNS.DeleteZone(ctx, zoneID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete DNS zone: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("DNS zone %s deleted successfully", zoneID)), nil
	})

	// Create DNS record tool
	createRecordTool := mcp.NewTool("bizflycloud_create_dns_record",
		mcp.WithDescription("Create a DNS record in a Bizfly Cloud DNS zone"),
		mcp.WithString("zone_id",
			mcp.Required(),
			mcp.Description("ID of the DNS zone"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the DNS record"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of DNS record (A, AAAA, CNAME, MX, TXT, SRV, etc.)"),
		),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("DNS record data (comma-separated for multiple values)"),
		),
		mcp.WithNumber("ttl",
			mcp.Description("TTL for the DNS record"),
		),
	)
	s.AddTool(createRecordTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, ok := request.Params.Arguments["zone_id"].(string)
		if !ok {
			return nil, errors.New("zone_id must be a string")
		}
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		recordType, ok := request.Params.Arguments["type"].(string)
		if !ok {
			return nil, errors.New("type must be a string")
		}
		dataStr, ok := request.Params.Arguments["data"].(string)
		if !ok {
			return nil, errors.New("data must be a string")
		}

		ttl := 3600 // Default TTL
		if ttlVal, ok := request.Params.Arguments["ttl"].(float64); ok {
			ttl = int(ttlVal)
		}

		var payload interface{}
		payload = &gobizfly.CreateNormalRecordPayload{
			BaseCreateRecordPayload: gobizfly.BaseCreateRecordPayload{
				Name: name,
				Type: recordType,
				TTL:  ttl,
			},
			Data: []string{dataStr},
		}

		record, err := client.DNS.CreateRecord(ctx, zoneID, payload)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create DNS record: %v", err)), nil
		}

		result := fmt.Sprintf("DNS record created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", record.Name)
		result += fmt.Sprintf("  ID: %s\n", record.ID)
		result += fmt.Sprintf("  Type: %s\n", record.Type)
		result += fmt.Sprintf("  TTL: %d\n", record.TTL)
		return mcp.NewToolResultText(result), nil
	})

	// Get DNS record tool
	getRecordTool := mcp.NewTool("bizflycloud_get_dns_record",
		mcp.WithDescription("Get details of a Bizfly Cloud DNS record"),
		mcp.WithString("record_id",
			mcp.Required(),
			mcp.Description("ID of the DNS record"),
		),
	)
	s.AddTool(getRecordTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		recordID, ok := request.Params.Arguments["record_id"].(string)
		if !ok {
			return nil, errors.New("record_id must be a string")
		}
		record, err := client.DNS.GetRecord(ctx, recordID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get DNS record: %v", err)), nil
		}

		result := fmt.Sprintf("DNS Record Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", record.Name)
		result += fmt.Sprintf("ID: %s\n", record.ID)
		result += fmt.Sprintf("Type: %s\n", record.Type)
		result += fmt.Sprintf("TTL: %d\n", record.TTL)
		return mcp.NewToolResultText(result), nil
	})

	// Delete DNS record tool
	deleteRecordTool := mcp.NewTool("bizflycloud_delete_dns_record",
		mcp.WithDescription("Delete a Bizfly Cloud DNS record"),
		mcp.WithString("record_id",
			mcp.Required(),
			mcp.Description("ID of the DNS record to delete"),
		),
	)
	s.AddTool(deleteRecordTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		recordID, ok := request.Params.Arguments["record_id"].(string)
		if !ok {
			return nil, errors.New("record_id must be a string")
		}
		err := client.DNS.DeleteRecord(ctx, recordID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete DNS record: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("DNS record %s deleted successfully", recordID)), nil
	})
}

