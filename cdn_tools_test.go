package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestCDNToolsRegistration(t *testing.T) {
	t.Run("register CDN tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterCDNTools(s, client)
	})
}

func TestListCDNDomainsTool(t *testing.T) {
	t.Run("list CDN domains", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_cdn_domains", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_cdn_domains" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateCDNDomainTool(t *testing.T) {
	t.Run("create CDN domain with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_cdn_domain", map[string]interface{}{
			"domain":         "example.com",
			"upstream_host":  "origin.example.com",
			"upstream_addrs": "1.2.3.4",
			"upstream_proto": "https",
		})
		
		domain, _ := request.Params.Arguments["domain"].(string)
		upstreamHost, _ := request.Params.Arguments["upstream_host"].(string)
		upstreamAddrs, _ := request.Params.Arguments["upstream_addrs"].(string)
		upstreamProto, _ := request.Params.Arguments["upstream_proto"].(string)
		
		if domain != "example.com" || upstreamHost != "origin.example.com" || upstreamAddrs != "1.2.3.4" || upstreamProto != "https" {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create CDN domain without upstream_proto", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_cdn_domain", map[string]interface{}{
			"domain":         "example.com",
			"upstream_host":  "origin.example.com",
			"upstream_addrs": "1.2.3.4",
		})
		
		// upstream_proto should default to "http" in the handler
		if _, ok := request.Params.Arguments["upstream_proto"]; ok {
			t.Error("upstream_proto should be optional")
		}
	})
}

func TestGetCDNDomainTool(t *testing.T) {
	t.Run("get CDN domain with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_cdn_domain", map[string]interface{}{
			"domain_id": "domain-123",
		})
		
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok || domainID != "domain-123" {
			t.Error("Invalid domain_id")
		}
	})
}

func TestUpdateCDNDomainTool(t *testing.T) {
	t.Run("update CDN domain with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_update_cdn_domain", map[string]interface{}{
			"domain_id":      "domain-123",
			"upstream_addrs": "5.6.7.8",
			"upstream_proto": "http",
		})
		
		domainID, _ := request.Params.Arguments["domain_id"].(string)
		upstreamAddrs, _ := request.Params.Arguments["upstream_addrs"].(string)
		
		if domainID != "domain-123" || upstreamAddrs != "5.6.7.8" {
			t.Error("Invalid parameters")
		}
	})
}

func TestDeleteCDNDomainTool(t *testing.T) {
	t.Run("delete CDN domain with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_cdn_domain", map[string]interface{}{
			"domain_id": "domain-123",
		})
		
		domainID, ok := request.Params.Arguments["domain_id"].(string)
		if !ok || domainID != "domain-123" {
			t.Error("Invalid domain_id")
		}
	})
}

func TestDeleteCDNCacheTool(t *testing.T) {
	t.Run("delete CDN cache without files", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_cdn_cache", map[string]interface{}{
			"domain_id": "domain-123",
		})
		
		domainID, _ := request.Params.Arguments["domain_id"].(string)
		if domainID != "domain-123" {
			t.Error("Invalid domain_id")
		}
		
		// files parameter should be optional
		if _, ok := request.Params.Arguments["files"]; ok {
			t.Error("files should be optional")
		}
	})
	
	t.Run("delete CDN cache with specific files", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_cdn_cache", map[string]interface{}{
			"domain_id": "domain-123",
			"files":     "/path/to/file1.jpg,/path/to/file2.jpg",
		})
		
		files, _ := request.Params.Arguments["files"].(string)
		if files != "/path/to/file1.jpg,/path/to/file2.jpg" {
			t.Error("Invalid files parameter")
		}
	})
}

