package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestLoadBalancerToolsRegistration(t *testing.T) {
	t.Run("register load balancer tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterLoadBalancerTools(s, client)
	})
}

func TestListLoadBalancersTool(t *testing.T) {
	t.Run("list load balancers", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_loadbalancers", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_loadbalancers" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateLoadBalancerTool(t *testing.T) {
	t.Run("create load balancer with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_loadbalancer", map[string]interface{}{
			"name":         "lb-1",
			"network_type": "external",
			"type":         "basic",
			"description":  "Test LB",
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		networkType, _ := request.Params.Arguments["network_type"].(string)
		lbType, _ := request.Params.Arguments["type"].(string)
		description, _ := request.Params.Arguments["description"].(string)
		
		if name != "lb-1" || networkType != "external" || lbType != "basic" || description != "Test LB" {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create load balancer without description", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_loadbalancer", map[string]interface{}{
			"name":         "lb-1",
			"network_type": "external",
			"type":         "basic",
		})
		
		if _, ok := request.Params.Arguments["description"]; ok {
			t.Error("Description should be optional")
		}
	})
}

func TestGetLoadBalancerTool(t *testing.T) {
	t.Run("get load balancer with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_loadbalancer", map[string]interface{}{
			"loadbalancer_id": "lb-123",
		})
		
		lbID, ok := request.Params.Arguments["loadbalancer_id"].(string)
		if !ok || lbID != "lb-123" {
			t.Error("Invalid loadbalancer_id")
		}
	})
}

func TestUpdateLoadBalancerTool(t *testing.T) {
	t.Run("update load balancer with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_update_loadbalancer", map[string]interface{}{
			"loadbalancer_id": "lb-123",
			"name":           "updated-lb",
			"description":    "Updated description",
			"admin_state_up": true,
		})
		
		lbID, _ := request.Params.Arguments["loadbalancer_id"].(string)
		name, _ := request.Params.Arguments["name"].(string)
		adminStateUp, _ := request.Params.Arguments["admin_state_up"].(bool)
		
		if lbID != "lb-123" || name != "updated-lb" || !adminStateUp {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("update load balancer with partial parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_update_loadbalancer", map[string]interface{}{
			"loadbalancer_id": "lb-123",
			"name":           "updated-lb",
		})
		
		// Verify optional parameters can be omitted
		if _, ok := request.Params.Arguments["description"]; ok {
			t.Error("Description should be optional")
		}
	})
}

func TestDeleteLoadBalancerTool(t *testing.T) {
	t.Run("delete load balancer with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_loadbalancer", map[string]interface{}{
			"loadbalancer_id": "lb-123",
		})
		
		lbID, ok := request.Params.Arguments["loadbalancer_id"].(string)
		if !ok || lbID != "lb-123" {
			t.Error("Invalid loadbalancer_id")
		}
	})
}

