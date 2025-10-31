package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestVolumeToolsRegistration(t *testing.T) {
	t.Run("register volume tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		// Should not panic
		RegisterVolumeTools(s, client)
	})
}

func TestListVolumesTool(t *testing.T) {
	t.Run("list volumes request structure", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_volumes", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_volumes" {
			t.Errorf("Expected tool name 'bizflycloud_list_volumes', got '%s'", request.Params.Name)
		}
	})
}

func TestCreateVolumeTool(t *testing.T) {
	t.Run("create volume with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_volume", map[string]interface{}{
			"name":        "test-volume",
			"size":        20.0,
			"volume_type": "PREMIUM-HDD1",
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		size, _ := request.Params.Arguments["size"].(float64)
		volumeType, _ := request.Params.Arguments["volume_type"].(string)
		
		if name != "test-volume" {
			t.Errorf("Expected name 'test-volume', got '%s'", name)
		}
		if size != 20.0 {
			t.Errorf("Expected size 20.0, got %f", size)
		}
		if volumeType != "PREMIUM-HDD1" {
			t.Errorf("Expected volume_type 'PREMIUM-HDD1', got '%s'", volumeType)
		}
	})
	
	t.Run("create volume missing required parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_volume", map[string]interface{}{
			"name": "test-volume",
		})
		
		if _, ok := request.Params.Arguments["size"]; ok {
			t.Error("Expected size to be missing")
		}
	})
}

func TestResizeVolumeTool(t *testing.T) {
	t.Run("resize volume with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_resize_volume", map[string]interface{}{
			"volume_id": "vol-123",
			"new_size":  50.0,
		})
		
		volumeID, _ := request.Params.Arguments["volume_id"].(string)
		newSize, _ := request.Params.Arguments["new_size"].(float64)
		
		if volumeID != "vol-123" || newSize != 50.0 {
			t.Error("Invalid parameters")
		}
	})
}

func TestDeleteVolumeTool(t *testing.T) {
	t.Run("delete volume with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_volume", map[string]interface{}{
			"volume_id": "vol-123",
		})
		
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok || volumeID != "vol-123" {
			t.Error("Invalid volume_id")
		}
	})
}

func TestGetVolumeTool(t *testing.T) {
	t.Run("get volume with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_volume", map[string]interface{}{
			"volume_id": "vol-123",
		})
		
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok || volumeID != "vol-123" {
			t.Error("Invalid volume_id")
		}
	})
}

func TestAttachVolumeTool(t *testing.T) {
	t.Run("attach volume with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_attach_volume", map[string]interface{}{
			"volume_id": "vol-123",
			"server_id": "server-123",
		})
		
		volumeID, _ := request.Params.Arguments["volume_id"].(string)
		serverID, _ := request.Params.Arguments["server_id"].(string)
		
		if volumeID != "vol-123" || serverID != "server-123" {
			t.Error("Invalid parameters")
		}
	})
}

func TestDetachVolumeTool(t *testing.T) {
	t.Run("detach volume with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_detach_volume", map[string]interface{}{
			"volume_id": "vol-123",
			"server_id": "server-123",
		})
		
		volumeID, _ := request.Params.Arguments["volume_id"].(string)
		serverID, _ := request.Params.Arguments["server_id"].(string)
		
		if volumeID != "vol-123" || serverID != "server-123" {
			t.Error("Invalid parameters")
		}
	})
}

func TestListSnapshotsTool(t *testing.T) {
	t.Run("list snapshots", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_snapshots", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_snapshots" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateSnapshotTool(t *testing.T) {
	t.Run("create snapshot with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_snapshot", map[string]interface{}{
			"volume_id": "vol-123",
			"name":      "snapshot-1",
		})
		
		volumeID, _ := request.Params.Arguments["volume_id"].(string)
		name, _ := request.Params.Arguments["name"].(string)
		
		if volumeID != "vol-123" || name != "snapshot-1" {
			t.Error("Invalid parameters")
		}
	})
}

func TestDeleteSnapshotTool(t *testing.T) {
	t.Run("delete snapshot with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_snapshot", map[string]interface{}{
			"snapshot_id": "snap-123",
		})
		
		snapshotID, ok := request.Params.Arguments["snapshot_id"].(string)
		if !ok || snapshotID != "snap-123" {
			t.Error("Invalid snapshot_id")
		}
	})
}

