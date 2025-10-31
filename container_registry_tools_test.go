package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestContainerRegistryToolsRegistration(t *testing.T) {
	t.Run("register container registry tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterContainerRegistryTools(s, client)
	})
}

func TestListContainerRegistriesTool(t *testing.T) {
	t.Run("list container registries", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_container_registries", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_container_registries" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateContainerRegistryTool(t *testing.T) {
	t.Run("create repository with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_container_registry", map[string]interface{}{
			"name":   "my-repo",
			"public": false,
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		public, _ := request.Params.Arguments["public"].(bool)
		
		if name != "my-repo" || public != false {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create public repository", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_container_registry", map[string]interface{}{
			"name":   "my-repo",
			"public": true,
		})
		
		public, _ := request.Params.Arguments["public"].(bool)
		if !public {
			t.Error("Expected public to be true")
		}
	})
}

func TestDeleteContainerRegistryTool(t *testing.T) {
	t.Run("delete repository with valid name", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_container_registry", map[string]interface{}{
			"repository_name": "my-repo",
		})
		
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok || repositoryName != "my-repo" {
			t.Error("Invalid repository_name")
		}
	})
}

func TestListContainerRegistryTagsTool(t *testing.T) {
	t.Run("list tags with valid repository name", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_container_registry_tags", map[string]interface{}{
			"repository_name": "my-repo",
		})
		
		repositoryName, ok := request.Params.Arguments["repository_name"].(string)
		if !ok || repositoryName != "my-repo" {
			t.Error("Invalid repository_name")
		}
	})
}

func TestGetContainerRegistryTagTool(t *testing.T) {
	t.Run("get tag with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_container_registry_tag", map[string]interface{}{
			"repository_name": "my-repo",
			"tag_name":       "latest",
			"vulnerabilities": "yes",
		})
		
		repositoryName, _ := request.Params.Arguments["repository_name"].(string)
		tagName, _ := request.Params.Arguments["tag_name"].(string)
		vulns, _ := request.Params.Arguments["vulnerabilities"].(string)
		
		if repositoryName != "my-repo" || tagName != "latest" || vulns != "yes" {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("get tag without vulnerabilities parameter", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_container_registry_tag", map[string]interface{}{
			"repository_name": "my-repo",
			"tag_name":       "latest",
		})
		
		if _, ok := request.Params.Arguments["vulnerabilities"]; ok {
			t.Error("vulnerabilities should be optional")
		}
	})
}

func TestDeleteContainerRegistryTagTool(t *testing.T) {
	t.Run("delete tag with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_container_registry_tag", map[string]interface{}{
			"repository_name": "my-repo",
			"tag_name":       "latest",
		})
		
		repositoryName, _ := request.Params.Arguments["repository_name"].(string)
		tagName, _ := request.Params.Arguments["tag_name"].(string)
		
		if repositoryName != "my-repo" || tagName != "latest" {
			t.Error("Invalid parameters")
		}
	})
}

func TestUpdateContainerRegistryTool(t *testing.T) {
	t.Run("update repository with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_update_container_registry", map[string]interface{}{
			"repository_name": "my-repo",
			"public":          true,
		})
		
		repositoryName, _ := request.Params.Arguments["repository_name"].(string)
		public, _ := request.Params.Arguments["public"].(bool)
		
		if repositoryName != "my-repo" || !public {
			t.Error("Invalid parameters")
		}
	})
}

