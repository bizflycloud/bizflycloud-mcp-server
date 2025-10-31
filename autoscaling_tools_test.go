package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestAutoScalingToolsRegistration(t *testing.T) {
	t.Run("register autoscaling tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterAutoScalingTools(s, client)
	})
}

func TestListAutoScalingGroupsTool(t *testing.T) {
	t.Run("list autoscaling groups", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_autoscaling_groups", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_autoscaling_groups" {
			t.Error("Invalid tool name")
		}
	})
	
	t.Run("list autoscaling groups with all flag", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_autoscaling_groups", map[string]interface{}{
			"all": true,
		})
		
		all, ok := request.Params.Arguments["all"].(bool)
		if !ok || !all {
			t.Error("Expected all to be true")
		}
	})
}

func TestGetAutoScalingGroupTool(t *testing.T) {
	t.Run("get autoscaling group with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_autoscaling_group", map[string]interface{}{
			"group_id": "group-123",
		})
		
		groupID, ok := request.Params.Arguments["group_id"].(string)
		if !ok || groupID != "group-123" {
			t.Error("Invalid group_id")
		}
	})
}

func TestCreateAutoScalingGroupTool(t *testing.T) {
	t.Run("create autoscaling group with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_autoscaling_group", map[string]interface{}{
			"name":            "group-1",
			"profile_id":      "profile-123",
			"min_size":        1.0,
			"max_size":        10.0,
			"desired_capacity": 3.0,
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		profileID, _ := request.Params.Arguments["profile_id"].(string)
		minSize, _ := request.Params.Arguments["min_size"].(float64)
		maxSize, _ := request.Params.Arguments["max_size"].(float64)
		desiredCapacity, _ := request.Params.Arguments["desired_capacity"].(float64)
		
		if name != "group-1" || profileID != "profile-123" || minSize != 1.0 || maxSize != 10.0 || desiredCapacity != 3.0 {
			t.Error("Invalid parameters")
		}
	})
}

func TestDeleteAutoScalingGroupTool(t *testing.T) {
	t.Run("delete autoscaling group with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_autoscaling_group", map[string]interface{}{
			"group_id": "group-123",
		})
		
		groupID, ok := request.Params.Arguments["group_id"].(string)
		if !ok || groupID != "group-123" {
			t.Error("Invalid group_id")
		}
	})
}

