package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterContainerRegistryTools registers all Container Registry-related tools with the MCP server
func RegisterContainerRegistryTools(s *server.MCPServer, client *gobizfly.Client) {
	// List repositories tool
	listRepositoriesTool := mcp.NewTool("bizflycloud_list_container_registries",
		mcp.WithDescription("List all Bizfly Cloud Container Registry repositories"),
	)
	s.AddTool(listRepositoriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositories, err := client.ContainerRegistry.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list repositories: %v", err)), nil
		}

		result := "Available repositories:\n\n"
		for _, repo := range repositories {
			result += fmt.Sprintf("Repository: %s\n", repo.Name)
			result += fmt.Sprintf("  Public: %v\n", repo.Public)
			result += fmt.Sprintf("  Pulls: %d\n", repo.Pulls)
			result += fmt.Sprintf("  Last Push: %s\n", repo.LastPush)
			result += fmt.Sprintf("  Created At: %s\n", repo.CreatedAt)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create repository tool
	createRepositoryTool := mcp.NewTool("bizflycloud_create_container_registry",
		mcp.WithDescription("Create a new Bizfly Cloud Container Registry repository"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the repository"),
		),
		mcp.WithBoolean("public",
			mcp.Description("Whether the repository is public (default: false)"),
		),
	)
	s.AddTool(createRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		public, _ := request.Params.Arguments["public"].(bool)

		err := client.ContainerRegistry.Create(ctx, &gobizfly.CreateRepositoryPayload{
			Name:   name,
			Public: public,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create repository: %v", err)), nil
		}

		result := fmt.Sprintf("Repository created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", name)
		result += fmt.Sprintf("  Public: %v\n", public)
		return mcp.NewToolResultText(result), nil
	})

	// Delete repository tool
	deleteRepositoryTool := mcp.NewTool("bizflycloud_delete_container_registry",
		mcp.WithDescription("Delete a Bizfly Cloud Container Registry repository"),
		mcp.WithString("repository_name",
			mcp.Required(),
			mcp.Description("Name of the repository to delete"),
		),
	)
	s.AddTool(deleteRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok {
			return nil, errors.New("repository_name must be a string")
		}
		err := client.ContainerRegistry.Delete(ctx, repositoryName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete repository: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Repository %s deleted successfully", repositoryName)), nil
	})

	// Get repository tags tool
	getTagsTool := mcp.NewTool("bizflycloud_list_container_registry_tags",
		mcp.WithDescription("List tags for a Bizfly Cloud Container Registry repository"),
		mcp.WithString("repository_name",
			mcp.Required(),
			mcp.Description("Name of the repository"),
		),
	)
	s.AddTool(getTagsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok {
			return nil, errors.New("repository_name must be a string")
		}
		tags, err := client.ContainerRegistry.GetTags(ctx, repositoryName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get tags: %v", err)), nil
		}

		result := fmt.Sprintf("Repository: %s\n\n", tags.Repository.Name)
		result += fmt.Sprintf("Tags:\n\n")
		for _, tag := range tags.Tags {
			result += fmt.Sprintf("Tag: %s\n", tag.Name)
			result += fmt.Sprintf("  Author: %s\n", tag.Author)
			result += fmt.Sprintf("  Created At: %s\n", tag.CreatedAt)
			result += fmt.Sprintf("  Last Updated: %s\n", tag.LastUpdated)
			result += fmt.Sprintf("  Vulnerabilities: %d\n", tag.Vulnerabilities)
			result += fmt.Sprintf("  Fixes: %d\n", tag.Fixes)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get tag details tool
	getTagTool := mcp.NewTool("bizflycloud_get_container_registry_tag",
		mcp.WithDescription("Get details of a Container Registry tag"),
		mcp.WithString("repository_name",
			mcp.Required(),
			mcp.Description("Name of the repository"),
		),
		mcp.WithString("tag_name",
			mcp.Required(),
			mcp.Description("Name of the tag"),
		),
		mcp.WithString("vulnerabilities",
			mcp.Description("Include vulnerabilities (yes/no, default: no)"),
		),
	)
	s.AddTool(getTagTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok {
			return nil, errors.New("repository_name must be a string")
		}
		tagName, ok := request.Params.Arguments["tag_name"].(string)
		if !ok {
			return nil, errors.New("tag_name must be a string")
		}
		vulns, _ := request.Params.Arguments["vulnerabilities"].(string)
		if vulns != "yes" {
			vulns = "no"
		}

		image, err := client.ContainerRegistry.GetTag(ctx, repositoryName, tagName, vulns)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get tag: %v", err)), nil
		}

		result := fmt.Sprintf("Tag Details:\n\n")
		result += fmt.Sprintf("Repository: %s\n", image.Repository.Name)
		result += fmt.Sprintf("Tag: %s\n", image.Tag.Name)
		result += fmt.Sprintf("Author: %s\n", image.Tag.Author)
		result += fmt.Sprintf("Created At: %s\n", image.Tag.CreatedAt)
		result += fmt.Sprintf("Last Updated: %s\n", image.Tag.LastUpdated)
		result += fmt.Sprintf("Scan Status: %s\n", image.Tag.ScanStatus)
		result += fmt.Sprintf("Vulnerabilities: %d\n", image.Tag.Vulnerabilities)
		result += fmt.Sprintf("Fixes: %d\n", image.Tag.Fixes)
		if len(image.Vulnerabilities) > 0 {
			result += fmt.Sprintf("\nVulnerabilities:\n")
			for _, vuln := range image.Vulnerabilities {
				result += fmt.Sprintf("  - %s (%s): %s\n", vuln.Name, vuln.Severity, vuln.Description)
			}
		}
		return mcp.NewToolResultText(result), nil
	})

	// Delete tag tool
	deleteTagTool := mcp.NewTool("bizflycloud_delete_container_registry_tag",
		mcp.WithDescription("Delete a tag from a Bizfly Cloud Container Registry repository"),
		mcp.WithString("repository_name",
			mcp.Required(),
			mcp.Description("Name of the repository"),
		),
		mcp.WithString("tag_name",
			mcp.Required(),
			mcp.Description("Name of the tag to delete"),
		),
	)
	s.AddTool(deleteTagTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok {
			return nil, errors.New("repository_name must be a string")
		}
		tagName, ok := request.Params.Arguments["tag_name"].(string)
		if !ok {
			return nil, errors.New("tag_name must be a string")
		}
		err := client.ContainerRegistry.DeleteTag(ctx, tagName, repositoryName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete tag: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Tag %s deleted from repository %s successfully", tagName, repositoryName)), nil
	})

	// Update repository tool
	updateRepositoryTool := mcp.NewTool("bizflycloud_update_container_registry",
		mcp.WithDescription("Update a Bizfly Cloud Container Registry repository"),
		mcp.WithString("repository_name",
			mcp.Required(),
			mcp.Description("Name of the repository to update"),
		),
		mcp.WithBoolean("public",
			mcp.Required(),
			mcp.Description("Whether the repository should be public"),
		),
	)
	s.AddTool(updateRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok {
			return nil, errors.New("repository_name must be a string")
		}
		public, ok := request.Params.Arguments["public"].(bool)
		if !ok {
			return nil, errors.New("public must be a boolean")
		}

		err := client.ContainerRegistry.EditRepo(ctx, repositoryName, &gobizfly.EditRepositoryPayload{
			Public: public,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update repository: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Repository %s updated successfully", repositoryName)), nil
	})
}

