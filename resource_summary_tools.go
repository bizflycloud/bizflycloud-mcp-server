package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterResourceSummaryTools registers a tool to list all resources in table format
func RegisterResourceSummaryTools(s *server.MCPServer, client *gobizfly.Client) {
	// List all resources tool
	listAllResourcesTool := mcp.NewTool("bizflycloud_list_all_resources",
		mcp.WithDescription("List all Bizfly Cloud resources in a formatted table"),
	)
	s.AddTool(listAllResourcesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[DEBUG] List All Resources tool called")
		
		var result strings.Builder
		result.WriteString("# Bizfly Cloud Resources Summary\n\n")
		
		// Declare variables for summary
		var servers []*gobizfly.Server
		var volumes []*gobizfly.Volume
		var clusters []*gobizfly.Cluster
		var databases []*gobizfly.CloudDatabaseInstance
		var repos []*gobizfly.Repository
		var cdnDomains *gobizfly.DomainsResp
		var snapshots []*gobizfly.Snapshot
		
		// Servers
		result.WriteString("## 1. Servers\n\n")
		var err error
		servers, err = client.CloudServer.List(ctx, &gobizfly.ServerListOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			result.WriteString(fmt.Sprintf("**Total: %d servers**\n\n", len(servers)))
			if len(servers) > 0 {
				result.WriteString("| Name | Status | Flavor | Zone | WAN IP | LAN IP | Created At |\n")
				result.WriteString("|------|--------|--------|------|--------|--------|------------|\n")
				for _, srv := range servers {
					wanIP := "-"
					if len(srv.IPAddresses.WanV4Addresses) > 0 {
						wanIP = string(srv.IPAddresses.WanV4Addresses[0].Address)
					}
					lanIP := "-"
					if len(srv.IPAddresses.LanAddresses) > 0 {
						lanIP = string(srv.IPAddresses.LanAddresses[0].Address)
					}
					createdAt := formatDate(srv.CreatedAt)
					result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
						srv.Name, srv.Status, srv.FlavorName, srv.AvailabilityZone, wanIP, lanIP, createdAt))
				}
			} else {
				result.WriteString("(No servers found)\n")
			}
		}
		result.WriteString("\n")
		
		// Volumes
		result.WriteString("## 2. Volumes\n\n")
		volumes, err = client.CloudServer.Volumes().List(ctx, &gobizfly.VolumeListOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			inUseCount := 0
			availableCount := 0
			for _, vol := range volumes {
				if vol.Status == "in-use" {
					inUseCount++
				} else {
					availableCount++
				}
			}
			result.WriteString(fmt.Sprintf("**Total: %d volumes** (%d in-use, %d available)\n\n", len(volumes), inUseCount, availableCount))
			if len(volumes) > 0 {
				result.WriteString("| Name | Status | Size | Type | Zone | Created At |\n")
				result.WriteString("|------|--------|------|------|------|------------|\n")
				for _, vol := range volumes {
					createdAt := formatDate(vol.CreatedAt)
					result.WriteString(fmt.Sprintf("| %s | %s | %d GB | %s | %s | %s |\n",
						vol.Name, vol.Status, vol.Size, vol.VolumeType, vol.AvailabilityZone, createdAt))
				}
			} else {
				result.WriteString("(No volumes found)\n")
			}
		}
		result.WriteString("\n")
		
		// Kubernetes Clusters
		result.WriteString("## 3. Kubernetes Clusters\n\n")
		clusters, err = client.KubernetesEngine.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			result.WriteString(fmt.Sprintf("**Total: %d clusters**\n\n", len(clusters)))
			if len(clusters) > 0 {
				result.WriteString("| Name | Status | Version | Node Pools | Created At |\n")
				result.WriteString("|------|--------|---------|------------|------------|\n")
				for _, cluster := range clusters {
					createdAt := formatDate(cluster.CreatedAt)
					result.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s |\n",
						cluster.Name, cluster.ClusterStatus, getK8sVersion(cluster), cluster.WorkerPoolsCount, createdAt))
				}
			} else {
				result.WriteString("(No clusters found)\n")
			}
		}
		result.WriteString("\n")
		
		// Databases
		result.WriteString("## 4. Databases\n\n")
		databases, err = client.CloudDatabase.Instances().List(ctx, &gobizfly.CloudDatabaseListOption{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			result.WriteString(fmt.Sprintf("**Total: %d databases**\n\n", len(databases)))
			if len(databases) > 0 {
				result.WriteString("| Name | Type | Status | Created At |\n")
				result.WriteString("|------|------|--------|------------|\n")
				for _, db := range databases {
					createdAt := formatDate(db.CreatedAt)
					result.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						db.Name, db.Datastore.Type, db.Status, createdAt))
				}
			} else {
				result.WriteString("(No databases found)\n")
			}
		}
		result.WriteString("\n")
		
		// Container Registries
		result.WriteString("## 5. Container Registries\n\n")
		repos, err = client.ContainerRegistry.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			publicCount := 0
			privateCount := 0
			for _, repo := range repos {
				if repo.Public {
					publicCount++
				} else {
					privateCount++
				}
			}
			result.WriteString(fmt.Sprintf("**Total: %d repositories** (%d public, %d private)\n\n", len(repos), publicCount, privateCount))
			if len(repos) > 0 {
				result.WriteString("| Repository | Public | Pulls | Last Push |\n")
				result.WriteString("|------------|--------|-------|-----------|\n")
				for _, repo := range repos {
					public := "❌"
					if repo.Public {
						public = "✅"
					}
					lastPush := "-"
					if repo.LastPush != "" {
						lastPush = formatDate(repo.LastPush)
					}
					result.WriteString(fmt.Sprintf("| %s | %s | %d | %s |\n",
						repo.Name, public, repo.Pulls, lastPush))
				}
			} else {
				result.WriteString("(No repositories found)\n")
			}
		}
		result.WriteString("\n")
		
		// CDN Domains
		result.WriteString("## 6. CDN Domains\n\n")
		cdnDomains, err = client.CDN.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			if cdnDomains != nil && len(cdnDomains.Domains) > 0 {
				result.WriteString(fmt.Sprintf("**Total: %d domains**\n\n", len(cdnDomains.Domains)))
				result.WriteString("| Domain | CDN Domain |\n")
				result.WriteString("|--------|------------|\n")
				for _, domain := range cdnDomains.Domains {
					result.WriteString(fmt.Sprintf("| %s | %s |\n", domain.Domain, domain.DomainCDN))
				}
			} else {
				if cdnDomains == nil {
					result.WriteString("(CDN service not available or no domains found)\n")
				} else {
					result.WriteString("(No CDN domains found)\n")
				}
			}
		}
		result.WriteString("\n")
		
		// KMS Certificates
		result.WriteString("## 7. KMS Certificates\n\n")
		if client.KMS != nil && client.KMS.Certificates() != nil {
			certificates, err := client.KMS.Certificates().List(ctx)
			if err != nil {
				result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
			} else {
				result.WriteString(fmt.Sprintf("**Total: %d certificates**\n\n", len(certificates)))
				if len(certificates) > 0 {
					result.WriteString("| Certificate Name | Container ID |\n")
					result.WriteString("|------------------|--------------|\n")
					for _, cert := range certificates {
						result.WriteString(fmt.Sprintf("| %s | %s |\n", cert.Name, cert.ContainerID))
					}
				} else {
					result.WriteString("(No certificates found)\n")
				}
			}
		} else {
			result.WriteString("(KMS service not available)\n")
		}
		result.WriteString("\n")
		
		// Auto Scaling Groups
		result.WriteString("## 8. Auto Scaling Groups\n\n")
		if client.AutoScaling != nil && client.AutoScaling.AutoScalingGroups() != nil {
			groups, err := client.AutoScaling.AutoScalingGroups().List(ctx, false)
			if err != nil {
				result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
			} else {
				result.WriteString(fmt.Sprintf("**Total: %d groups**\n\n", len(groups)))
				if len(groups) > 0 {
					result.WriteString("| Name | Status | Min/Max Size | Desired | Current Nodes |\n")
					result.WriteString("|------|--------|--------------|---------|---------------|\n")
					for _, group := range groups {
						result.WriteString(fmt.Sprintf("| %s | %s | %d/%d | %d | %d |\n",
							group.Name, group.Status, group.MinSize, group.MaxSize, group.DesiredCapacity, len(group.NodeIDs)))
					}
				} else {
					result.WriteString("(No groups found)\n")
				}
			}
		} else {
			result.WriteString("(AutoScaling service not available)\n")
		}
		result.WriteString("\n")
		
		// Snapshots
		result.WriteString("## 9. Snapshots\n\n")
		snapshots, err = client.CloudServer.Snapshots().List(ctx, &gobizfly.ListSnasphotsOptions{})
		if err != nil {
			result.WriteString(fmt.Sprintf("❌ Error: %v\n\n", err))
		} else {
			result.WriteString(fmt.Sprintf("**Total: %d snapshots**\n\n", len(snapshots)))
			if len(snapshots) > 0 {
				result.WriteString("| Name | Status | Size | Volume ID |\n")
				result.WriteString("|------|--------|------|-----------|\n")
				for _, snap := range snapshots {
					result.WriteString(fmt.Sprintf("| %s | %s | %d GB | %s |\n",
						snap.Name, snap.Status, snap.Size, snap.VolumeID))
				}
			} else {
				result.WriteString("(No snapshots found)\n")
			}
		}
		result.WriteString("\n")
		
		// Summary
		result.WriteString("---\n\n")
		result.WriteString("## Summary\n\n")
		result.WriteString(fmt.Sprintf("- **Servers**: %d\n", len(servers)))
		result.WriteString(fmt.Sprintf("- **Volumes**: %d\n", len(volumes)))
		result.WriteString(fmt.Sprintf("- **Kubernetes Clusters**: %d\n", len(clusters)))
		result.WriteString(fmt.Sprintf("- **Databases**: %d\n", len(databases)))
		result.WriteString(fmt.Sprintf("- **Container Registries**: %d\n", len(repos)))
		cdnDomainCount := 0
		if cdnDomains != nil {
			cdnDomainCount = len(cdnDomains.Domains)
		}
		result.WriteString(fmt.Sprintf("- **CDN Domains**: %d\n", cdnDomainCount))
		if client.KMS != nil && client.KMS.Certificates() != nil {
			certs, _ := client.KMS.Certificates().List(ctx)
			result.WriteString(fmt.Sprintf("- **KMS Certificates**: %d\n", len(certs)))
		}
		if client.AutoScaling != nil && client.AutoScaling.AutoScalingGroups() != nil {
			grps, _ := client.AutoScaling.AutoScalingGroups().List(ctx, false)
			result.WriteString(fmt.Sprintf("- **Auto Scaling Groups**: %d\n", len(grps)))
		}
		result.WriteString(fmt.Sprintf("- **Snapshots**: %d\n", len(snapshots)))
		
		return mcp.NewToolResultText(result.String()), nil
	})
}

// Helper functions
func formatDate(dateStr string) string {
	if dateStr == "" {
		return "-"
	}
	// Try to format date - extract just the date part
	if len(dateStr) >= 10 {
		return dateStr[:10]
	}
	return dateStr
}

func getK8sVersion(cluster *gobizfly.Cluster) string {
	if cluster.Version.K8SVersion != "" {
		return cluster.Version.K8SVersion
	}
	if cluster.Version.Name != "" {
		return cluster.Version.Name
	}
	return "-"
}

