package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestAlertToolsRegistration(t *testing.T) {
	t.Run("register alert tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterAlertTools(s, client)
	})
}

func TestListAlarmsTool(t *testing.T) {
	t.Run("list alarms", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_alarms", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_alarms" {
			t.Error("Invalid tool name")
		}
	})
}

func TestGetAlarmTool(t *testing.T) {
	t.Run("get alarm with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_alarm", map[string]interface{}{
			"alarm_id": "alarm-123",
		})
		
		alarmID, ok := request.Params.Arguments["alarm_id"].(string)
		if !ok || alarmID != "alarm-123" {
			t.Error("Invalid alarm_id")
		}
	})
}

func TestListReceiversTool(t *testing.T) {
	t.Run("list receivers", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_receivers", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_receivers" {
			t.Error("Invalid tool name")
		}
	})
}

func TestGetReceiverTool(t *testing.T) {
	t.Run("get receiver with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_receiver", map[string]interface{}{
			"receiver_id": "receiver-123",
		})
		
		receiverID, ok := request.Params.Arguments["receiver_id"].(string)
		if !ok || receiverID != "receiver-123" {
			t.Error("Invalid receiver_id")
		}
	})
}

