package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterKubernetesTools registers all Kubernetes-related tools with the MCP server
func RegisterKubernetesTools(s *server.MCPServer, client *gobizfly.Client) {
	// List clusters tool
	listClustersTool := mcp.NewTool("bizflycloud_list_kubernetes_clusters",
		mcp.WithDescription("List all Bizfly Cloud Kubernetes clusters"),
	)
	s.AddTool(listClustersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusters, err := client.KubernetesEngine.List(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list clusters: %v", err)), nil
		}

		result := "Available Kubernetes clusters:\n\n"
		for _, c := range clusters {
			// Get full cluster details
			cluster, err := client.KubernetesEngine.Get(ctx, c.UID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster details: %v", err)), nil
			}

			result += fmt.Sprintf("Cluster: %s\n", cluster.Name)
			result += fmt.Sprintf("  ID: %s\n", cluster.UID)
			result += fmt.Sprintf("  Status: %s\n", cluster.ClusterStatus)
			result += fmt.Sprintf("  Version: %s\n", cluster.Version)
			result += fmt.Sprintf("  Node Pools Count: %d\n", cluster.WorkerPoolsCount)
			result += fmt.Sprintf("  Created At: %s\n", cluster.CreatedAt)
			result += "\nWorker Pools:\n"
			for _, pool := range cluster.WorkerPools {
				result += fmt.Sprintf("  - Name: %s\n", pool.Name)
				result += fmt.Sprintf("    ID: %s\n", pool.UID)
				result += fmt.Sprintf("    Flavor: %s\n", pool.Flavor)
				result += fmt.Sprintf("    Profile Type: %s\n", pool.ProfileType)
				result += fmt.Sprintf("    Volume Type: %s\n", pool.VolumeType)
				result += fmt.Sprintf("    Volume Size: %d GB\n", pool.VolumeSize)
				result += fmt.Sprintf("    Desired Size: %d nodes\n", pool.DesiredSize)
				result += "\n"
			}
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Create cluster tool
	createClusterTool := mcp.NewTool("bizflycloud_create_kubernetes_cluster",
		mcp.WithDescription("Create a new Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the cluster"),
		),
		mcp.WithString("version",
			mcp.Required(),
			mcp.Description("Kubernetes version"),
		),
		mcp.WithString("worker_flavor",
			mcp.Required(),
			mcp.Description("Flavor for worker nodes"),
		),
		mcp.WithNumber("worker_count",
			mcp.Required(),
			mcp.Description("Number of worker nodes"),
		),
	)
	s.AddTool(createClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		version, ok := request.Params.Arguments["version"].(string)
		if !ok {
			return nil, errors.New("version must be a string")
		}
		workerFlavor, ok := request.Params.Arguments["worker_flavor"].(string)
		if !ok {
			return nil, errors.New("worker_flavor must be a string")
		}
		workerCount, ok := request.Params.Arguments["worker_count"].(float64)
		if !ok {
			return nil, errors.New("worker_count must be a number")
		}

		// Get the flavor ID from the name
		flavors, err := client.CloudServer.Flavors().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get flavors: %v", err)), nil
		}

		var flavorID string
		for _, flavor := range flavors {
			if flavor.Name == workerFlavor {
				flavorID = flavor.ID
				break
			}
		}
		if flavorID == "" {
			return mcp.NewToolResultError(fmt.Sprintf("Flavor '%s' not found", workerFlavor)), nil
		}

		cluster, err := client.KubernetesEngine.Create(ctx, &gobizfly.ClusterCreateRequest{
			Name:    name,
			Version: version,
			WorkerPools: []gobizfly.WorkerPool{
				{
					Name:              "default-pool",
					Flavor:            flavorID,
					ProfileType:       "premium",
					VolumeType:        "PREMIUM-HDD1",
					VolumeSize:        50,
					DesiredSize:       int(workerCount),
					EnableAutoScaling: false,
					MinSize:           int(workerCount),
					MaxSize:           int(workerCount),
					NetworkPlan:       "free_plan",
					BillingPlan:       "on_demand",
					AvailabilityZone:  "HN1",
				},
			},
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create cluster: %v", err)), nil
		}

		result := fmt.Sprintf("Cluster created successfully:\n")
		result += fmt.Sprintf("  Name: %s\n", cluster.Name)
		result += fmt.Sprintf("  ID: %s\n", cluster.UID)
		result += fmt.Sprintf("  Status: %s\n", cluster.ClusterStatus)
		result += fmt.Sprintf("  Version: %s\n", cluster.Version)
		result += fmt.Sprintf("  Node Pools Count: %d\n", cluster.WorkerPoolsCount)
		return mcp.NewToolResultText(result), nil
	})

	// Delete cluster tool
	deleteClusterTool := mcp.NewTool("bizflycloud_delete_kubernetes_cluster",
		mcp.WithDescription("Delete a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster to delete"),
		),
	)
	s.AddTool(deleteClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		err := client.KubernetesEngine.Delete(ctx, clusterID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete cluster: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Cluster %s deleted successfully", clusterID)), nil
	})

	// List cluster nodes tool
	listClusterNodesTool := mcp.NewTool("bizflycloud_list_kubernetes_nodes",
		mcp.WithDescription("List nodes in a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster"),
		),
		mcp.WithString("pool_id",
			mcp.Required(),
			mcp.Description("ID of the node pool"),
		),
	)
	s.AddTool(listClusterNodesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		poolID, ok := request.Params.Arguments["pool_id"].(string)
		if !ok {
			return nil, errors.New("pool_id must be a string")
		}

		// Get cluster details to find the pool
		cluster, err := client.KubernetesEngine.Get(ctx, clusterID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster: %v", err)), nil
		}

		// Find the pool
		var pool *gobizfly.ExtendedWorkerPool
		for _, p := range cluster.WorkerPools {
			if p.Name == poolID {
				pool = &p
				break
			}
		}
		if pool == nil {
			return mcp.NewToolResultError(fmt.Sprintf("Pool %s not found in cluster %s", poolID, clusterID)), nil
		}

		result := fmt.Sprintf("Worker Pool Details:\n")
		result += fmt.Sprintf("  Name: %s\n", pool.Name)
		result += fmt.Sprintf("  Flavor: %s\n", pool.Flavor)
		result += fmt.Sprintf("  Profile Type: %s\n", pool.ProfileType)
		result += fmt.Sprintf("  Volume Type: %s\n", pool.VolumeType)
		result += fmt.Sprintf("  Volume Size: %d GB\n", pool.VolumeSize)
		result += fmt.Sprintf("  Availability Zone: %s\n", pool.AvailabilityZone)
		result += fmt.Sprintf("  Desired Size: %d\n", pool.DesiredSize)
		result += fmt.Sprintf("  Auto Scaling: %v\n", pool.EnableAutoScaling)
		if pool.EnableAutoScaling {
			result += fmt.Sprintf("  Min Size: %d\n", pool.MinSize)
			result += fmt.Sprintf("  Max Size: %d\n", pool.MaxSize)
		}
		if len(pool.Tags) > 0 {
			result += fmt.Sprintf("  Tags: %v\n", pool.Tags)
		}
		if len(pool.Labels) > 0 {
			result += fmt.Sprintf("  Labels: %v\n", pool.Labels)
		}
		result += fmt.Sprintf("  Network Plan: %s\n", pool.NetworkPlan)
		result += fmt.Sprintf("  Billing Plan: %s\n", pool.BillingPlan)
		result += "\n"

		// Get nodes in the pool
		nodes, err := client.KubernetesEngine.GetClusterWorkerPool(ctx, clusterID, poolID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list nodes: %v", err)), nil
		}

		result += fmt.Sprintf("Nodes in pool %s:\n\n", pool.Name)
		for _, node := range nodes.Nodes {
			result += fmt.Sprintf("Node: %s\n", node.Name)
			result += fmt.Sprintf("  Status: %s\n", node.Status)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get cluster tool
	getClusterTool := mcp.NewTool("bizflycloud_get_kubernetes_cluster",
		mcp.WithDescription("Get details of a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster to get details for"),
		),
	)
	s.AddTool(getClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		cluster, err := client.KubernetesEngine.Get(ctx, clusterID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get cluster: %v", err)), nil
		}

		result := fmt.Sprintf("Cluster Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", cluster.Name)
		result += fmt.Sprintf("ID: %s\n", cluster.UID)
		result += fmt.Sprintf("Status: %s\n", cluster.ClusterStatus)
		result += fmt.Sprintf("Provision Status: %s\n", cluster.ProvisionStatus)
		result += fmt.Sprintf("Version: %s\n", cluster.Version.K8SVersion)
		result += fmt.Sprintf("Worker Pools Count: %d\n", cluster.WorkerPoolsCount)
		result += fmt.Sprintf("Auto Upgrade: %v\n", cluster.AutoUpgrade)
		result += fmt.Sprintf("Created At: %s\n", cluster.CreatedAt)
		result += "\nWorker Pools:\n"
		for _, pool := range cluster.WorkerPools {
			result += fmt.Sprintf("  - Name: %s\n", pool.Name)
			result += fmt.Sprintf("    ID: %s\n", pool.UID)
			result += fmt.Sprintf("    Flavor: %s\n", pool.Flavor)
			result += fmt.Sprintf("    Profile Type: %s\n", pool.ProfileType)
			result += fmt.Sprintf("    Volume Type: %s\n", pool.VolumeType)
			result += fmt.Sprintf("    Volume Size: %d GB\n", pool.VolumeSize)
			result += fmt.Sprintf("    Desired Size: %d nodes\n", pool.DesiredSize)
			result += fmt.Sprintf("    Auto Scaling: %v\n", pool.EnableAutoScaling)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Update pool tool
	updatePoolTool := mcp.NewTool("bizflycloud_update_kubernetes_pool",
		mcp.WithDescription("Update a worker pool in a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster"),
		),
		mcp.WithString("pool_id",
			mcp.Required(),
			mcp.Description("ID of the pool to update"),
		),
		mcp.WithNumber("desired_size",
			mcp.Description("Desired number of nodes in the pool"),
		),
		mcp.WithBoolean("enable_autoscaling",
			mcp.Description("Enable auto scaling for the pool"),
		),
		mcp.WithNumber("min_size",
			mcp.Description("Minimum number of nodes for auto scaling"),
		),
		mcp.WithNumber("max_size",
			mcp.Description("Maximum number of nodes for auto scaling"),
		),
	)
	s.AddTool(updatePoolTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		poolID, ok := request.Params.Arguments["pool_id"].(string)
		if !ok {
			return nil, errors.New("pool_id must be a string")
		}

		req := &gobizfly.UpdateWorkerPoolRequest{}
		if desiredSize, ok := request.Params.Arguments["desired_size"].(float64); ok {
			req.DesiredSize = int(desiredSize)
		}
		if enableAutoScaling, ok := request.Params.Arguments["enable_autoscaling"].(bool); ok {
			req.EnableAutoScaling = enableAutoScaling
		}
		if minSize, ok := request.Params.Arguments["min_size"].(float64); ok {
			req.MinSize = int(minSize)
		}
		if maxSize, ok := request.Params.Arguments["max_size"].(float64); ok {
			req.MaxSize = int(maxSize)
		}

		err := client.KubernetesEngine.UpdateClusterWorkerPool(ctx, clusterID, poolID, req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update pool: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Pool %s in cluster %s updated successfully", poolID, clusterID)), nil
	})

	// Resize pool tool (uses update with desired_size)
	resizePoolTool := mcp.NewTool("bizflycloud_resize_kubernetes_pool",
		mcp.WithDescription("Resize a worker pool in a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster"),
		),
		mcp.WithString("pool_id",
			mcp.Required(),
			mcp.Description("ID of the pool to resize"),
		),
		mcp.WithNumber("desired_size",
			mcp.Required(),
			mcp.Description("New desired number of nodes in the pool"),
		),
	)
	s.AddTool(resizePoolTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		poolID, ok := request.Params.Arguments["pool_id"].(string)
		if !ok {
			return nil, errors.New("pool_id must be a string")
		}
		desiredSize, ok := request.Params.Arguments["desired_size"].(float64)
		if !ok {
			return nil, errors.New("desired_size must be a number")
		}

		req := &gobizfly.UpdateWorkerPoolRequest{
			DesiredSize: int(desiredSize),
		}
		err := client.KubernetesEngine.UpdateClusterWorkerPool(ctx, clusterID, poolID, req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to resize pool: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Pool %s in cluster %s resized to %d nodes successfully", poolID, clusterID, int(desiredSize))), nil
	})

	// Delete pool tool
	deletePoolTool := mcp.NewTool("bizflycloud_delete_kubernetes_pool",
		mcp.WithDescription("Delete a worker pool from a Bizfly Cloud Kubernetes cluster"),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("ID of the cluster"),
		),
		mcp.WithString("pool_id",
			mcp.Required(),
			mcp.Description("ID of the pool to delete"),
		),
	)
	s.AddTool(deletePoolTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clusterID, ok := request.Params.Arguments["cluster_id"].(string)
		if !ok {
			return nil, errors.New("cluster_id must be a string")
		}
		poolID, ok := request.Params.Arguments["pool_id"].(string)
		if !ok {
			return nil, errors.New("pool_id must be a string")
		}
		err := client.KubernetesEngine.DeleteClusterWorkerPool(ctx, clusterID, poolID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete pool: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Pool %s deleted from cluster %s successfully", poolID, clusterID)), nil
	})
}