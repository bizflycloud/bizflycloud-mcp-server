# Bizfly Cloud MCP Server

A comprehensive Model Context Protocol (MCP) server implementation that connects to Bizfly Cloud to manage cloud resources. Built using the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) SDK.

## Features

- 🔧 **Complete Cloud Management**: Manage servers, volumes, load balancers, databases, Kubernetes clusters, and more
- 🔒 **Secure Authentication**: Uses environment variables for credentials
- 🐳 **Docker Support**: Ready-to-use Docker image for easy deployment
- 📦 **10 Services Supported**: Server, Volume, Load Balancer, Kubernetes, Database, DNS, CDN, KMS, Container Registry, AutoScaling, and Alert services
- ✅ **Fully Tested**: Comprehensive test suite with 196+ test cases
- 🚀 **MCP Protocol**: Compatible with Cursor and Claude Desktop
- 🌐 **Multi-Region Support**: Configurable region and API endpoints

## Prerequisites

### Local Development
-   Go 1.23 or later
-   Bizfly Cloud account credentials
-   Cursor or Claude Desktop installed

### Docker Deployment
-   Docker 20.10 or later
-   Docker Compose (optional)

## Setup

### Option 1: Docker (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/bizflycloud-mcp-server.git
   cd bizflycloud-mcp-server
   ```

2. **Build the Docker image:**
   ```bash
   docker build -t bizfly-mcp-server:latest .
   ```

3. **Run the container:**
   ```bash
   docker run -it --rm \
     -e BIZFLY_USERNAME=your_username \
     -e BIZFLY_PASSWORD=your_password \
     -e BIZFLY_REGION=HaNoi \
     bizfly-mcp-server:latest
   ```

4. **Using Docker Compose:**

   Create a `.env` file:
   ```bash
   BIZFLY_USERNAME=your_username
   BIZFLY_PASSWORD=your_password
   BIZFLY_REGION=HaNoi
   BIZFLY_API_URL=https://manage.bizflycloud.vn
   ```

   Run with docker-compose:
   ```bash
   docker-compose up
   ```

### Option 2: Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/bizflycloud-mcp-server.git
   cd bizflycloud-mcp-server
   ```

2. **Set up environment variables:**
   ```bash
   export BIZFLY_USERNAME=your_username
   export BIZFLY_PASSWORD=your_password
   export BIZFLY_REGION=HaNoi  # Optional, defaults to HaNoi
   export BIZFLY_API_URL=https://manage.bizflycloud.vn  # Optional, defaults to https://manage.bizflycloud.vn
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Build the server:**
   ```bash
   go build -o bizfly-mcp-server
   ```

## Running the Server

### For Cursor/Claude Desktop Integration

1. Build the server:

    ```bash
    go build -o bizfly-mcp-server
    ```

2. Configure your MCP client (Cursor or Claude Desktop) by adding the following to the configuration:

For Cursor:

```json
{
    "mcpServers": {
        "bizfly": {
            "command": "/absolute/path/to/bizfly-mcp-server",
            "env": {
                "BIZFLY_USERNAME": "your_username",
                "BIZFLY_PASSWORD": "your_password",
                "BIZFLY_REGION": "HaNoi"
            }
        }
    }
}
```

For Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
    "mcpServers": {
        "bizfly": {
            "command": "/absolute/path/to/bizfly-mcp-server",
            "env": {
                "BIZFLY_USERNAME": "your_username",
                "BIZFLY_PASSWORD": "your_password",
                "BIZFLY_REGION": "HaNoi"
            }
        }
    }
}
```

## Available Tools

The server provides comprehensive MCP tools for managing all Bizfly Cloud services. All tool names are prefixed with `bizflycloud_` for consistency.

### 🖥️ Server Management (`bizflycloud_*`)

-   `bizflycloud_list_servers` - List all Bizfly Cloud servers
-   `bizflycloud_get_server` - Get detailed information about a server
-   `bizflycloud_start_server` - Start a stopped server
-   `bizflycloud_stop_server` - Stop a running server
-   `bizflycloud_reboot_server` - Soft reboot a server
-   `bizflycloud_hard_reboot_server` - Hard reboot a server
-   `bizflycloud_delete_server` - Delete a server
-   `bizflycloud_resize_server` - Resize a server to a different flavor
-   `bizflycloud_list_flavors` - List available server flavors

### 💾 Volume Management (`bizflycloud_*`)

-   `bizflycloud_list_volumes` - List all volumes
-   `bizflycloud_get_volume` - Get detailed information about a volume
-   `bizflycloud_create_volume` - Create a new volume
-   `bizflycloud_delete_volume` - Delete a volume
-   `bizflycloud_resize_volume` - Resize a volume
-   `bizflycloud_attach_volume` - Attach a volume to a server
-   `bizflycloud_detach_volume` - Detach a volume from a server
-   `bizflycloud_list_snapshots` - List all volume snapshots
-   `bizflycloud_create_snapshot` - Create a volume snapshot
-   `bizflycloud_delete_snapshot` - Delete a volume snapshot

### ⚖️ Load Balancer Management (`bizflycloud_*`)

-   `bizflycloud_list_loadbalancers` - List all load balancers
-   `bizflycloud_get_loadbalancer` - Get detailed information about a load balancer
-   `bizflycloud_create_loadbalancer` - Create a new load balancer
-   `bizflycloud_update_loadbalancer` - Update load balancer properties
-   `bizflycloud_delete_loadbalancer` - Delete a load balancer

### ☸️ Kubernetes Management (`bizflycloud_*`)

-   `bizflycloud_list_kubernetes_clusters` - List all Kubernetes clusters
-   `bizflycloud_get_kubernetes_cluster` - Get detailed information about a cluster
-   `bizflycloud_create_kubernetes_cluster` - Create a new Kubernetes cluster
-   `bizflycloud_delete_kubernetes_cluster` - Delete a Kubernetes cluster
-   `bizflycloud_list_kubernetes_nodes` - List nodes in a cluster/pool
-   `bizflycloud_update_kubernetes_pool` - Update worker pool configuration
-   `bizflycloud_resize_kubernetes_pool` - Resize a worker pool
-   `bizflycloud_delete_kubernetes_pool` - Delete a worker pool

### 🗄️ Database Management (`bizflycloud_*`)

-   `bizflycloud_list_databases` - List all database instances
-   `bizflycloud_list_datastores` - List available database engines and versions
-   `bizflycloud_get_database` - Get detailed information about a database
-   `bizflycloud_create_database` - Create a new database instance
-   `bizflycloud_delete_database` - Delete a database instance
-   `bizflycloud_list_database_backups` - List backups for a database instance
-   `bizflycloud_create_database_backup` - Create a backup for a database instance

### 🌐 DNS Management (`bizflycloud_*`)

-   `bizflycloud_list_dns_zones` - List all DNS zones
-   `bizflycloud_get_dns_zone` - Get detailed information about a DNS zone
-   `bizflycloud_create_dns_zone` - Create a new DNS zone
-   `bizflycloud_delete_dns_zone` - Delete a DNS zone
-   `bizflycloud_create_dns_record` - Create a DNS record (A, AAAA, CNAME, MX, SRV, TXT, NS, PTR)
-   `bizflycloud_get_dns_record` - Get detailed information about a DNS record
-   `bizflycloud_delete_dns_record` - Delete a DNS record

### 🚀 CDN Management (`bizflycloud_*`)

-   `bizflycloud_list_cdn_domains` - List all CDN domains
-   `bizflycloud_get_cdn_domain` - Get detailed information about a CDN domain
-   `bizflycloud_create_cdn_domain` - Create a new CDN domain
-   `bizflycloud_update_cdn_domain` - Update CDN domain configuration
-   `bizflycloud_delete_cdn_domain` - Delete a CDN domain
-   `bizflycloud_delete_cdn_cache` - Delete CDN cache (all or specific files)

### 🔐 KMS (Key Management Service) (`bizflycloud_*`)

-   `bizflycloud_list_kms_certificates` - List all KMS certificates
-   `bizflycloud_get_kms_certificate` - Get detailed information about a certificate
-   `bizflycloud_create_kms_certificate` - Create a new KMS certificate container
-   `bizflycloud_delete_kms_certificate` - Delete a KMS certificate

### 📦 Container Registry Management (`bizflycloud_*`)

-   `bizflycloud_list_container_registries` - List all container registries/repositories
-   `bizflycloud_create_container_registry` - Create a new repository
-   `bizflycloud_update_container_registry` - Update repository settings (e.g., public/private)
-   `bizflycloud_delete_container_registry` - Delete a repository
-   `bizflycloud_list_container_registry_tags` - List tags in a repository
-   `bizflycloud_get_container_registry_tag` - Get detailed information about a tag
-   `bizflycloud_delete_container_registry_tag` - Delete a tag

### 📈 AutoScaling Management (`bizflycloud_*`)

-   `bizflycloud_list_autoscaling_groups` - List all auto scaling groups
-   `bizflycloud_get_autoscaling_group` - Get detailed information about an auto scaling group
-   `bizflycloud_create_autoscaling_group` - Create a new auto scaling group
-   `bizflycloud_delete_autoscaling_group` - Delete an auto scaling group

### 🚨 Alert/CloudWatcher Management (`bizflycloud_*`)

-   `bizflycloud_list_alarms` - List all alarms
-   `bizflycloud_get_alarm` - Get detailed information about an alarm
-   `bizflycloud_list_receivers` - List all notification receivers
-   `bizflycloud_get_receiver` - Get detailed information about a receiver

## Docker Configuration

### Using Docker Image with Cursor/Claude Desktop

For Cursor (`~/.cursor/mcp.json`):

```json
{
    "mcpServers": {
        "bizfly": {
            "command": "docker",
            "args": [
                "run", "-i", "--rm",
                "-e", "BIZFLY_USERNAME",
                "-e", "BIZFLY_PASSWORD",
                "-e", "BIZFLY_REGION",
                "bizfly-mcp-server:latest"
            ],
            "env": {
                "BIZFLY_USERNAME": "your_username",
                "BIZFLY_PASSWORD": "your_password",
                "BIZFLY_REGION": "HaNoi"
            }
        }
    }
}
```

For Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
    "mcpServers": {
        "bizfly": {
            "command": "docker",
            "args": [
                "run", "-i", "--rm",
                "-e", "BIZFLY_USERNAME",
                "-e", "BIZFLY_PASSWORD",
                "-e", "BIZFLY_REGION",
                "bizfly-mcp-server:latest"
            ],
            "env": {
                "BIZFLY_USERNAME": "your_username",
                "BIZFLY_PASSWORD": "your_password",
                "BIZFLY_REGION": "HaNoi"
            }
        }
    }
}
```

## Example Usage

You can interact with the server through natural language queries in Cursor or Claude Desktop:

### Server Management
-   "Show me all my Bizfly Cloud servers"
-   "Start server server-123"
-   "Reboot the server named production-web"
-   "List available server flavors"

### Volume Management
-   "List all volumes in my Bizfly Cloud account"
-   "Create a 100GB volume named data-storage"
-   "Attach volume vol-123 to server server-456"
-   "Show me all snapshots"

### Load Balancer
-   "List all load balancers"
-   "Create a new load balancer for my web servers"

### Kubernetes
-   "List all Kubernetes clusters"
-   "Show me the nodes in cluster cluster-123"
-   "Resize the worker pool to 5 nodes"

### Database
-   "Show me all databases"
-   "Create a MySQL 8.0 database"
-   "List backups for database db-123"

### DNS
-   "List all DNS zones"
-   "Create an A record for www.example.com pointing to 1.2.3.4"
-   "Show me all records in zone example.com"

### CDN
-   "List all CDN domains"
-   "Create a CDN domain for example.com"
-   "Clear CDN cache for domain cdn-123"

### Container Registry
-   "List all container repositories"
-   "Create a public repository named my-app"
-   "Show me all tags in repository my-app"

### AutoScaling & Alerts
-   "List all auto scaling groups"
-   "Show me all alarms"
-   "List notification receivers"

## MCP Implementation Details

This server uses the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) SDK to implement the Model Context Protocol:

1. **Standard I/O Transport**: Uses stdin/stdout for communication with MCP clients
2. **Tool Definitions**: Clear tool descriptions and parameters
3. **Error Handling**: Proper error reporting in MCP format
4. **Text Formatting**: Human-readable output for resource listings

## Environment Variables

The server uses environment variables for configuration:

### Required Variables
-   `BIZFLY_USERNAME`: Your Bizfly Cloud username
-   `BIZFLY_PASSWORD`: Your Bizfly Cloud password

### Optional Variables
-   `BIZFLY_REGION`: Region name (defaults to "HaNoi")
  - Available regions: `HaNoi`, `HoChiMinh`, etc.
-   `BIZFLY_API_URL`: API endpoint URL (defaults to "https://manage.bizflycloud.vn")

### Security Best Practices

⚠️ **Important**: Keep your credentials secure and never commit them to version control.

1. **For Docker**: Use `.env` files or Docker secrets
   ```bash
   # .env file (add to .gitignore)
   BIZFLY_USERNAME=your_username
   BIZFLY_PASSWORD=your_password
   ```

2. **For Local Development**: Use environment variables or a `.env` file loader
   ```bash
   export BIZFLY_USERNAME=your_username
   export BIZFLY_PASSWORD=your_password
   ```

3. **For Production**: Consider using:
   - Docker secrets
   - Kubernetes secrets
   - HashiCorp Vault
   - AWS Secrets Manager / Azure Key Vault / GCP Secret Manager

## Testing

The project includes comprehensive test coverage with 196+ test cases:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v ./... -run TestServerToolsRegistration

# Run tests with coverage
go test -cover ./...
```

See [TEST_COVERAGE.md](TEST_COVERAGE.md) for detailed test coverage information.

## Docker Commands

### Build the Image
```bash
docker build -t bizfly-mcp-server:latest .
```

### Run Container
```bash
docker run -it --rm \
  -e BIZFLY_USERNAME=your_username \
  -e BIZFLY_PASSWORD=your_password \
  -e BIZFLY_REGION=HaNoi \
  bizfly-mcp-server:latest
```

### Build for Specific Platform
```bash
# For ARM64 (Apple Silicon)
docker build --platform linux/arm64 -t bizfly-mcp-server:arm64 .

# For AMD64
docker build --platform linux/amd64 -t bizfly-mcp-server:amd64 .
```

### Using Docker Compose
```bash
# Start the container
docker-compose up

# Start in detached mode
docker-compose up -d

# Stop the container
docker-compose down

# View logs
docker-compose logs -f
```

## MCP Features

1. **Standard I/O Transport**: Uses stdin/stdout for seamless integration with Cursor/Claude Desktop
2. **Standardized Response Format**: All responses follow the MCP format with context, type, data, and root fields
3. **Resource Organization**: Resources are organized under root paths
4. **Type Safety**: Strong typing for all resources
5. **Error Handling**: Standardized error responses in MCP format
6. **Comprehensive Tool Coverage**: 80+ tools covering 10 Bizfly Cloud services
7. **Security**: Non-root user in Docker container, secure credential handling

## Development

### Project Structure
```
.
├── main.go                    # Entry point
├── server_tools.go           # Server management tools
├── volume_tools.go           # Volume management tools
├── loadbalancer_tools.go     # Load balancer tools
├── kubernetes_tools.go       # Kubernetes management tools
├── database_tools.go          # Database management tools
├── dns_tools.go              # DNS service tools
├── cdn_tools.go              # CDN service tools
├── kms_tools.go              # KMS service tools
├── container_registry_tools.go # Container registry tools
├── autoscaling_tools.go      # AutoScaling tools
├── alert_tools.go            # Alert/CloudWatcher tools
├── *_test.go                 # Test files
├── test_helpers.go           # Test utilities
├── Dockerfile                # Docker image definition
├── docker-compose.yml        # Docker Compose configuration
└── README.md                 # This file
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new features
5. Ensure all tests pass: `go test ./...`
6. Submit a pull request

## License

[Add your license here]

## Support

For issues, questions, or contributions, please open an issue on GitHub.
