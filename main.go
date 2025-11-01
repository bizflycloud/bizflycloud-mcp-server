package main

import (
	"context"
	"log"
	"os"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/server"
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
		log.Fatalf("Failed to create BizflyCloud client: %v", err)
	}

	// Initialize token
	ctx := context.Background()
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

	// Create MCP server
	s := server.NewMCPServer(
		"BizflyCloud MCP",
		"1.0.0",
	)

	// Register tools
	RegisterServerTools(s, client)
	RegisterVolumeTools(s, client)
	RegisterKubernetesTools(s, client)
	RegisterDatabaseTools(s, client)
	RegisterLoadBalancerTools(s, client)
	RegisterDNSTools(s, client)
	RegisterCDNTools(s, client)
	RegisterKMSTools(s, client)
	RegisterContainerRegistryTools(s, client)
	RegisterAutoScalingTools(s, client)
	RegisterAlertTools(s, client)
	RegisterResourceSummaryTools(s, client)

	// Start the stdio server for Cursor/Claude Desktop integration
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
} 