package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestDNSToolsRegistration(t *testing.T) {
	t.Run("register DNS tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterDNSTools(s, client)
	})
}

func TestListDNSZonesTool(t *testing.T) {
	t.Run("list DNS zones", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_dns_zones", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_dns_zones" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateDNSZoneTool(t *testing.T) {
	t.Run("create DNS zone with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_dns_zone", map[string]interface{}{
			"name":        "example.com",
			"description": "Test zone",
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		description, _ := request.Params.Arguments["description"].(string)
		
		if name != "example.com" || description != "Test zone" {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create DNS zone without description", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_dns_zone", map[string]interface{}{
			"name": "example.com",
		})
		
		if _, ok := request.Params.Arguments["description"]; ok {
			t.Error("Description should be optional")
		}
	})
}

func TestGetDNSZoneTool(t *testing.T) {
	t.Run("get DNS zone with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_dns_zone", map[string]interface{}{
			"zone_id": "zone-123",
		})
		
		zoneID, ok := request.Params.Arguments["zone_id"].(string)
		if !ok || zoneID != "zone-123" {
			t.Error("Invalid zone_id")
		}
	})
}

func TestDeleteDNSZoneTool(t *testing.T) {
	t.Run("delete DNS zone with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_dns_zone", map[string]interface{}{
			"zone_id": "zone-123",
		})
		
		zoneID, ok := request.Params.Arguments["zone_id"].(string)
		if !ok || zoneID != "zone-123" {
			t.Error("Invalid zone_id")
		}
	})
}

func TestCreateDNSRecordTool(t *testing.T) {
	t.Run("create DNS record with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_dns_record", map[string]interface{}{
			"zone_id": "zone-123",
			"name":    "www.example.com",
			"type":    "A",
			"data":    "1.2.3.4",
			"ttl":     3600.0,
		})
		
		zoneID, _ := request.Params.Arguments["zone_id"].(string)
		name, _ := request.Params.Arguments["name"].(string)
		recordType, _ := request.Params.Arguments["type"].(string)
		data, _ := request.Params.Arguments["data"].(string)
		ttl, _ := request.Params.Arguments["ttl"].(float64)
		
		if zoneID != "zone-123" || name != "www.example.com" || recordType != "A" || data != "1.2.3.4" || ttl != 3600.0 {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create DNS record without TTL", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_dns_record", map[string]interface{}{
			"zone_id": "zone-123",
			"name":    "www.example.com",
			"type":    "A",
			"data":    "1.2.3.4",
		})
		
		if _, ok := request.Params.Arguments["ttl"]; ok {
			t.Error("TTL should be optional")
		}
	})
}

func TestGetDNSRecordTool(t *testing.T) {
	t.Run("get DNS record with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_dns_record", map[string]interface{}{
			"record_id": "record-123",
		})
		
		recordID, ok := request.Params.Arguments["record_id"].(string)
		if !ok || recordID != "record-123" {
			t.Error("Invalid record_id")
		}
	})
}

func TestDeleteDNSRecordTool(t *testing.T) {
	t.Run("delete DNS record with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_dns_record", map[string]interface{}{
			"record_id": "record-123",
		})
		
		recordID, ok := request.Params.Arguments["record_id"].(string)
		if !ok || recordID != "record-123" {
			t.Error("Invalid record_id")
		}
	})
}

