package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bizflycloud/gobizfly"
)

const (
	defaultRegion = "HaNoi"
	defaultAPIURL = "https://manage.bizflycloud.vn"
)

func main() {
	// Load environment variables
	username := os.Getenv("BIZFLY_USERNAME")
	password := os.Getenv("BIZFLY_PASSWORD")
	region := os.Getenv("BIZFLY_REGION")
	apiURL := os.Getenv("BIZFLY_API_URL")

	if username == "" || password == "" {
		log.Fatal("BIZFLY_USERNAME and BIZFLY_PASSWORD environment variables are required")
	}

	// Set defaults if not provided
	if region == "" {
		region = defaultRegion
	}
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	// Initialize Bizfly client
	client, err := gobizfly.NewClient(
		gobizfly.WithAPIURL(apiURL),
		gobizfly.WithRegionName(region),
	)
	if err != nil {
		log.Fatalf("Failed to create Bizfly Cloud client: %v", err)
	}

	// Authenticate with username and password
	ctx := context.Background()
	fmt.Println("=== Authenticating with username and password ===")
	token, err := client.Token.Init(
		ctx,
		&gobizfly.TokenCreateRequest{
			AuthMethod: "password",
			Username:   username,
			Password:   password,
		})
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	// Set the token
	client.SetKeystoneToken(token)
	fmt.Println("✅ Authentication successful\n")

	fmt.Println("=== Listing Kubernetes Clusters ===")
	fmt.Println()

	// List clusters
	clusters, err := client.KubernetesEngine.List(ctx, &gobizfly.ListOptions{})
	if err != nil {
		// Check if it's a 404 error
		errStr := err.Error()
		if containsAny(errStr, "404", "<svg", "Resource not found") {
			fmt.Println("❌ Error: Kubernetes Engine service may not be enabled or API endpoint not available")
			fmt.Printf("   Details: %v\n", err)
			fmt.Println()
			fmt.Println("Please check:")
			fmt.Println("  - Kubernetes Engine service is enabled on your account")
			fmt.Println("  - Your credentials have permission to access Kubernetes Engine")
			fmt.Println("  - BIZFLY_USERNAME and BIZFLY_PASSWORD are set correctly")
			fmt.Printf("  - Region: %s\n", region)
			fmt.Printf("  - API URL: %s\n", apiURL)
			os.Exit(1)
		}
		log.Fatalf("Failed to list clusters: %v", err)
	}

	// Print summary
	fmt.Printf("Total clusters found: %d\n\n", len(clusters))

	if len(clusters) == 0 {
		fmt.Println("No clusters found.")
		fmt.Println()
		fmt.Println("Possible reasons:")
		fmt.Println("  - No clusters exist in this account/project")
		fmt.Println("  - Clusters are in a different project/region")
		fmt.Println("  - Insufficient permissions")
		os.Exit(0)
	}

	// Print raw JSON for debugging
	fmt.Println("=== Raw JSON Response ===")
	clustersJSON, err := json.MarshalIndent(clusters, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal clusters to JSON: %v", err)
	} else {
		fmt.Println(string(clustersJSON))
		fmt.Println()
	}

	// Print formatted output
	fmt.Println("=== Formatted Cluster List ===")
	for i, c := range clusters {
		fmt.Printf("\n--- Cluster %d ---\n", i+1)
		fmt.Printf("Name:              %s\n", c.Name)
		fmt.Printf("UID:               %s\n", c.UID)
		fmt.Printf("Status:            %s\n", c.ClusterStatus)
		fmt.Printf("Provision Status:  %s\n", c.ProvisionStatus)
		
		if c.Version.K8SVersion != "" {
			fmt.Printf("Version:           %s\n", c.Version.K8SVersion)
		} else if c.Version.Name != "" {
			fmt.Printf("Version:           %s\n", c.Version.Name)
		}
		
		fmt.Printf("Worker Pools:      %d\n", c.WorkerPoolsCount)
		fmt.Printf("Created At:        %s\n", c.CreatedAt)
		
		if c.VPCNetworkID != "" {
			fmt.Printf("VPC Network ID:   %s\n", c.VPCNetworkID)
		}
		
		if c.SubnetID != "" {
			fmt.Printf("Subnet ID:        %s\n", c.SubnetID)
		}

		// Try to get full details
		fmt.Println("\nFetching full details...")
		cluster, err := client.KubernetesEngine.Get(ctx, c.UID)
		if err != nil {
			fmt.Printf("  ⚠️  Warning: Could not fetch full details: %v\n", err)
		} else {
			fmt.Printf("  ✅ Successfully fetched full details\n")
			if len(cluster.WorkerPools) > 0 {
				fmt.Println("\n  Worker Pools:")
				for j, pool := range cluster.WorkerPools {
					fmt.Printf("    Pool %d:\n", j+1)
					fmt.Printf("      Name:         %s\n", pool.Name)
					fmt.Printf("      UID:          %s\n", pool.UID)
					fmt.Printf("      Flavor:       %s\n", pool.Flavor)
					fmt.Printf("      Profile Type: %s\n", pool.ProfileType)
					fmt.Printf("      Volume Type:  %s\n", pool.VolumeType)
					fmt.Printf("      Volume Size:  %d GB\n", pool.VolumeSize)
					fmt.Printf("      Desired Size: %d nodes\n", pool.DesiredSize)
					fmt.Printf("      Auto Scaling: %v\n", pool.EnableAutoScaling)
					if pool.EnableAutoScaling {
						fmt.Printf("      Min Size:     %d\n", pool.MinSize)
						fmt.Printf("      Max Size:     %d\n", pool.MaxSize)
					}
				}
			}
		}
	}
	fmt.Println()
}

func containsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

