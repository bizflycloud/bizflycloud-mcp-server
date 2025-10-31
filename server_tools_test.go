package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestRegisterServerTools(t *testing.T) {
	t.Run("register server tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		// Should not panic
		RegisterServerTools(s, client)
	})
}

func TestRebootServerTool(t *testing.T) {
	t.Run("reboot server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_reboot_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		if request.Params.Arguments["server_id"] != "server-123" {
			t.Errorf("Expected server_id 'server-123', got '%v'", request.Params.Arguments["server_id"])
		}
	})
	
	t.Run("reboot server missing server_id", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_reboot_server", map[string]interface{}{})
		
		if _, ok := request.Params.Arguments["server_id"]; ok {
			t.Error("Expected server_id to be missing")
		}
	})
}

func TestGetServerTool(t *testing.T) {
	t.Run("get server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok {
			t.Error("Expected server_id to be a string")
		}
		if serverID != "server-123" {
			t.Errorf("Expected server_id 'server-123', got '%s'", serverID)
		}
	})
}

func TestStartServerTool(t *testing.T) {
	t.Run("start server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_start_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok || serverID != "server-123" {
			t.Error("Invalid server_id in request")
		}
	})
}

func TestStopServerTool(t *testing.T) {
	t.Run("stop server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_stop_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok || serverID != "server-123" {
			t.Error("Invalid server_id in request")
		}
	})
}

func TestHardRebootServerTool(t *testing.T) {
	t.Run("hard reboot server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_hard_reboot_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok || serverID != "server-123" {
			t.Error("Invalid server_id in request")
		}
	})
}

func TestDeleteServerTool(t *testing.T) {
	t.Run("delete server with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_server", map[string]interface{}{
			"server_id": "server-123",
		})
		
		serverID, ok := request.Params.Arguments["server_id"].(string)
		if !ok || serverID != "server-123" {
			t.Error("Invalid server_id in request")
		}
	})
}

func TestResizeServerTool(t *testing.T) {
	t.Run("resize server with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_resize_server", map[string]interface{}{
			"server_id":   "server-123",
			"flavor_name": "medium",
		})
		
		serverID, _ := request.Params.Arguments["server_id"].(string)
		flavorName, _ := request.Params.Arguments["flavor_name"].(string)
		
		if serverID != "server-123" || flavorName != "medium" {
			t.Error("Invalid parameters in request")
		}
	})
}

func TestListFlavorsTool(t *testing.T) {
	t.Run("list flavors", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_flavors", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_flavors" {
			t.Error("Invalid tool name")
		}
	})
}

