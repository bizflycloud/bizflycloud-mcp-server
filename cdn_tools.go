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

// RegisterCDNTools registers all CDN-related tools with the MCP server
func RegisterCDNTools(s *server.MCPServer, client *gobizfly.Client) {
	// List CDN domains tool
	listDomainsTool := mcp.NewTool("bizflycloud_list_cdn_domains",
		mcp.WithDescription("List all Bizfly Cloud CDN domains"),
	)
	s.AddTool(listDomainsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domains, err := client.CDN.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			// Check if error is 404 or service not available
			errStr := strings.ToLower(err.Error())
			if strings.Contains(errStr, "404") || 
			   strings.Contains(errStr, "not found") || 
			   strings.Contains(errStr, "resource not found") ||
			   strings.Contains(errStr, "<svg") ||
			   strings.Contains(errStr, "<html") {
				return mcp.NewToolResultText("Available CDN domains:\n\n(No CDN domains found or CDN service is not enabled)"), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list CDN domains: %v", err)), nil
		}

		result := "Available CDN domains:\n\n"
		if domains == nil || len(domains.Domains) == 0 {
			result += "(No CDN domains found)\n"
		} else {
			for _, domain := range domains.Domains {
				result += fmt.Sprintf("Domain: %s\n", domain.Domain)
				result += fmt.Sprintf("  ID: %s\n", domain.DomainID)
				result += fmt.Sprintf("  Slug: %s\n", domain.Slug)
				result += fmt.Sprintf("  CDN Domain: %s\n", domain.DomainCDN)
				result += "\n"
			}
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create CDN domain tool
	createDomainTool := mcp.NewTool("bizflycloud_create_cdn_domain",
		mcp.WithDescription("Create a new Bizfly Cloud CDN domain"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("Domain name for CDN (e.g., example.com)"),
		),
		mcp.WithString("upstream_host",
			mcp.Required(),
			mcp.Description("Upstream host for the origin"),
		),
		mcp.WithString("upstream_addrs",
			mcp.Required(),
			mcp.Description("Upstream addresses (comma-separated IPs or domains)"),
		),
		mcp.WithString("upstream_proto",
			mcp.Description("Upstream protocol (http or https, default: http)"),
		),
	)
	s.AddTool(createDomainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domain, ok := request.Params.Arguments["domain"].(string)
		if !ok {
			return nil, errors.New("domain must be a string")
		}
		upstreamHost, ok := request.Params.Arguments["upstream_host"].(string)
		if !ok {
			return nil, errors.New("upstream_host must be a string")
		}
		upstreamAddrs, ok := request.Params.Arguments["upstream_addrs"].(string)
		if !ok {
			return nil, errors.New("upstream_addrs must be a string")
		}
		upstreamProto, _ := request.Params.Arguments["upstream_proto"].(string)
		if upstreamProto == "" {
			upstreamProto = "http"
		}

		resp, err := client.CDN.Create(ctx, &gobizfly.CreateDomainPayload{
			Domain: domain,
			Origin: &gobizfly.Origin{
				Name:          domain,
				UpstreamHost:  upstreamHost,
				UpstreamAddrs: upstreamAddrs,
				UpstreamProto: upstreamProto,
			},
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create CDN domain: %v", err)), nil
		}

		result := fmt.Sprintf("CDN domain created successfully:\n")
		result += fmt.Sprintf("  Domain: %s\n", resp.Domain.Domain)
		result += fmt.Sprintf("  ID: %s\n", resp.Domain.DomainID)
		result += fmt.Sprintf("  CDN Domain: %s\n", resp.Domain.DomainCDN)
		result += fmt.Sprintf("  Message: %s\n", resp.Message)
		return mcp.NewToolResultText(result), nil
	})

	// Get CDN domain tool
	getDomainTool := mcp.NewTool("bizflycloud_get_cdn_domain",
		mcp.WithDescription("Get details of a Bizfly Cloud CDN domain"),
		mcp.WithString("domain_id",
			mcp.Required(),
			mcp.Description("ID of the CDN domain"),
		),
	)
	s.AddTool(getDomainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok {
			return nil, errors.New("domain_id must be a string")
		}
		domain, err := client.CDN.Get(ctx, domainID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get CDN domain: %v", err)), nil
		}

		result := fmt.Sprintf("CDN Domain Details:\n\n")
		result += fmt.Sprintf("Domain: %s\n", domain.Domain)
		result += fmt.Sprintf("ID: %s\n", domain.DomainID)
		result += fmt.Sprintf("Slug: %s\n", domain.Slug)
		result += fmt.Sprintf("CDN Domain: %s\n", domain.DomainCDN)
		return mcp.NewToolResultText(result), nil
	})

	// Update CDN domain tool
	updateDomainTool := mcp.NewTool("bizflycloud_update_cdn_domain",
		mcp.WithDescription("Update a Bizfly Cloud CDN domain"),
		mcp.WithString("domain_id",
			mcp.Required(),
			mcp.Description("ID of the CDN domain to update"),
		),
		mcp.WithString("upstream_addrs",
			mcp.Description("New upstream addresses (comma-separated)"),
		),
		mcp.WithString("upstream_proto",
			mcp.Description("Upstream protocol (http or https)"),
		),
	)
	s.AddTool(updateDomainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok {
			return nil, errors.New("domain_id must be a string")
		}

		payload := &gobizfly.UpdateDomainPayload{}
		if upstreamAddrs, ok := request.Params.Arguments["upstream_addrs"].(string); ok && upstreamAddrs != "" {
			if upstreamProto, ok := request.Params.Arguments["upstream_proto"].(string); ok && upstreamProto != "" {
				payload.Origin = &gobizfly.Origin{
					UpstreamAddrs: upstreamAddrs,
					UpstreamProto: upstreamProto,
				}
			}
		}

		resp, err := client.CDN.Update(ctx, domainID, payload)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update CDN domain: %v", err)), nil
		}

		result := fmt.Sprintf("CDN domain updated successfully:\n")
		result += fmt.Sprintf("  Domain: %s\n", resp.Domain.Domain)
		result += fmt.Sprintf("  Message: %s\n", resp.Message)
		return mcp.NewToolResultText(result), nil
	})

	// Delete CDN domain tool
	deleteDomainTool := mcp.NewTool("bizflycloud_delete_cdn_domain",
		mcp.WithDescription("Delete a Bizfly Cloud CDN domain"),
		mcp.WithString("domain_id",
			mcp.Required(),
			mcp.Description("ID of the CDN domain to delete"),
		),
	)
	s.AddTool(deleteDomainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok {
			return nil, errors.New("domain_id must be a string")
		}
		err := client.CDN.Delete(ctx, domainID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete CDN domain: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("CDN domain %s deleted successfully", domainID)), nil
	})

	// Delete CDN cache tool
	deleteCacheTool := mcp.NewTool("bizflycloud_delete_cdn_cache",
		mcp.WithDescription("Delete cache for a Bizfly Cloud CDN domain"),
		mcp.WithString("domain_id",
			mcp.Required(),
			mcp.Description("ID of the CDN domain"),
		),
		mcp.WithString("files",
			mcp.Description("Comma-separated list of file paths to purge (leave empty to purge all)"),
		),
	)
	s.AddTool(deleteCacheTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok {
			return nil, errors.New("domain_id must be a string")
		}

		filesStr, _ := request.Params.Arguments["files"].(string)
		var files *gobizfly.Files
		if filesStr != "" {
			fileList := strings.Split(filesStr, ",")
			for i, f := range fileList {
				fileList[i] = strings.TrimSpace(f)
			}
			files = &gobizfly.Files{
				Files: fileList,
			}
		} else {
			files = &gobizfly.Files{
				Files: []string{},
			}
		}

		err := client.CDN.DeleteCache(ctx, domainID, files)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete CDN cache: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("CDN cache for domain %s deleted successfully", domainID)), nil
	})
}

