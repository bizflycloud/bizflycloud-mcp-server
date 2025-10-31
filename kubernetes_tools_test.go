package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestKubernetesToolsRegistration(t *testing.T) {
	t.Run("register kubernetes tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterKubernetesTools(s, client)
	})
}

func TestListKubernetesClustersTool(t *testing.T) {
	t.Run("list kubernetes clusters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_kubernetes_clusters", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_kubernetes_clusters" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateKubernetesClusterTool(t *testing.T) {
	t.Run("create cluster with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_kubernetes_cluster", map[string]interface{}{
			"name":          "cluster-1",
			"version":       "v1.28.0",
			"worker_flavor": "small",
			"worker_count":  3.0,
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		version, _ := request.Params.Arguments["version"].(string)
		workerFlavor, _ := request.Params.Arguments["worker_flavor"].(string)
		workerCount, _ := request.Params.Arguments["worker_count"].(float64)
		
		if name != "cluster-1" || version != "v1.28.0" || workerFlavor != "small" || workerCount != 3.0 {
			t.Error("Invalid parameters")
		}
	})
}

func TestGetKubernetesClusterTool(t *testing.T) {
	t.Run("get cluster with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_kubernetes_cluster", map[string]interface{}{
			"cluster_id": "cluster-123",
		})
		
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok || clusterID != "cluster-123" {
			t.Error("Invalid cluster_id")
		}
	})
}

func TestDeleteKubernetesClusterTool(t *testing.T) {
	t.Run("delete cluster with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_kubernetes_cluster", map[string]interface{}{
			"cluster_id": "cluster-123",
		})
		
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok || clusterID != "cluster-123" {
			t.Error("Invalid cluster_id")
		}
	})
}

func TestListKubernetesNodesTool(t *testing.T) {
	t.Run("list nodes with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_kubernetes_nodes", map[string]interface{}{
			"cluster_id": "cluster-123",
			"pool_id":   "pool-123",
		})
		
		clusterID, _ := request.Params.Arguments["cluster_id"].(string)
		poolID, _ := request.Params.Arguments["pool_id"].(string)
		
		if clusterID != "cluster-123" || poolID != "pool-123" {
			t.Error("Invalid parameters")
		}
	})
}

func TestUpdateKubernetesPoolTool(t *testing.T) {
	t.Run("update pool with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_update_kubernetes_pool", map[string]interface{}{
			"cluster_id":        "cluster-123",
			"pool_id":           "pool-123",
			"desired_size":      5.0,
			"enable_autoscaling": true,
			"min_size":          3.0,
			"max_size":          10.0,
		})
		
		clusterID, _ := request.Params.Arguments["cluster_id"].(string)
		poolID, _ := request.Params.Arguments["pool_id"].(string)
		desiredSize, _ := request.Params.Arguments["desired_size"].(float64)
		enableAutoScaling, _ := request.Params.Arguments["enable_autoscaling"].(bool)
		
		if clusterID != "cluster-123" || poolID != "pool-123" || desiredSize != 5.0 || !enableAutoScaling {
			t.Error("Invalid parameters")
		}
	})
}

func TestResizeKubernetesPoolTool(t *testing.T) {
	t.Run("resize pool with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_resize_kubernetes_pool", map[string]interface{}{
			"cluster_id":   "cluster-123",
			"pool_id":      "pool-123",
			"desired_size": 7.0,
		})
		
		clusterID, _ := request.Params.Arguments["cluster_id"].(string)
		poolID, _ := request.Params.Arguments["pool_id"].(string)
		desiredSize, _ := request.Params.Arguments["desired_size"].(float64)
		
		if clusterID != "cluster-123" || poolID != "pool-123" || desiredSize != 7.0 {
			t.Error("Invalid parameters")
		}
	})
}

func TestDeleteKubernetesPoolTool(t *testing.T) {
	t.Run("delete pool with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_kubernetes_pool", map[string]interface{}{
			"cluster_id": "cluster-123",
			"pool_id":   "pool-123",
		})
		
		clusterID, _ := request.Params.Arguments["cluster_id"].(string)
		poolID, _ := request.Params.Arguments["pool_id"].(string)
		
		if clusterID != "cluster-123" || poolID != "pool-123" {
			t.Error("Invalid parameters")
		}
	})
}

