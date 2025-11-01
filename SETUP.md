# Setup Guide for Bizfly Cloud MCP Server

This guide will help you configure the Bizfly Cloud MCP Server for Cursor or Claude Desktop.

## Prerequisites

-   Docker installed (for Docker setup) OR Go 1.23+ (for local setup)
-   Cursor IDE or Claude Desktop installed
-   Bizfly Cloud account credentials

## Step 1: Choose Your Setup Method

### Option A: Using Docker (Recommended)

1. **Build the Docker image:**

    ```bash
    docker build -t bizfly-mcp-server:latest .
    ```

2. **Verify the image was created:**
    ```bash
    docker images bizfly-mcp-server
    ```

### Option B: Using Local Binary

1. **Build the binary:**

    ```bash
    go build -o bizfly-mcp-server
    ```

2. **Note the absolute path:**
    ```bash
    pwd  # Save this path
    # Example: /Users/username/WORK/Github/bizflycloud-mcp-server
    ```

## Step 2: Configure for Cursor

### Method 1: Docker Configuration

1. **Create or edit the Cursor MCP config file:**

    ```bash
    mkdir -p ~/.cursor
    nano ~/.cursor/mcp.json  # or use your preferred editor
    ```

2. **Add the following configuration:**
    ```json
    {
        "mcpServers": {
            "bizflycloud": {
                "command": "docker",
                "args": [
                    "run",
                    "-i",
                    "--rm",
                    "-e",
                    "BIZFLY_USERNAME",
                    "-e",
                    "BIZFLY_PASSWORD",
                    "-e",
                    "BIZFLY_REGION",
                    "-e",
                    "BIZFLY_API_URL",
                    "bizfly-mcp-server:latest"
                ],
                "env": {
                    "BIZFLY_USERNAME": "your_username",
                    "BIZFLY_PASSWORD": "your_password",
                    "BIZFLY_REGION": "HaNoi",
                    "BIZFLY_API_URL": "https://manage.bizflycloud.vn"
                }
            }
        }
    }
    ```

### Method 2: Local Binary Configuration

1. **Create or edit the Cursor MCP config file:**

    ```bash
    mkdir -p ~/.cursor
    nano ~/.cursor/mcp.json
    ```

2. **Add the following configuration (replace the path with your actual path):**

    ```json
    {
        "mcpServers": {
            "bizflycloud": {
                "command": "/Users/your-username/WORK/Github/bizflycloud-mcp-server/bizfly-mcp-server",
                "env": {
                    "BIZFLY_USERNAME": "your_username",
                    "BIZFLY_PASSWORD": "your_password",
                    "BIZFLY_REGION": "HaNoi",
                    "BIZFLY_API_URL": "https://manage.bizflycloud.vn"
                }
            }
        }
    }
    ```

3. **Make sure the binary is executable:**
    ```bash
    chmod +x bizfly-mcp-server
    ```

## Step 3: Configure for Claude Desktop

### Method 1: Docker Configuration

1. **Create or edit the Claude Desktop config file (macOS):**

    ```bash
    mkdir -p ~/Library/Application\ Support/Claude
    nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
    ```

    **Windows:**

    ```powershell
    # Path: %APPDATA%\Claude\claude_desktop_config.json
    ```

    **Linux:**

    ```bash
    mkdir -p ~/.config/Claude
    nano ~/.config/Claude/claude_desktop_config.json
    ```

2. **Add the following configuration:**
    ```json
    {
        "mcpServers": {
            "bizflycloud": {
                "command": "docker",
                "args": [
                    "run",
                    "-i",
                    "--rm",
                    "-e",
                    "BIZFLY_USERNAME",
                    "-e",
                    "BIZFLY_PASSWORD",
                    "-e",
                    "BIZFLY_REGION",
                    "-e",
                    "BIZFLY_API_URL",
                    "bizfly-mcp-server:latest"
                ],
                "env": {
                    "BIZFLY_USERNAME": "your_username",
                    "BIZFLY_PASSWORD": "your_password",
                    "BIZFLY_REGION": "HaNoi",
                    "BIZFLY_API_URL": "https://manage.bizflycloud.vn"
                }
            }
        }
    }
    ```

### Method 2: Local Binary Configuration

1. **Create or edit the Claude Desktop config file (see paths above)**

2. **Add the following configuration:**
    ```json
    {
        "mcpServers": {
            "bizflycloud": {
                "command": "/Users/your-username/WORK/Github/bizflycloud-mcp-server/bizfly-mcp-server",
                "env": {
                    "BIZFLY_USERNAME": "your_username",
                    "BIZFLY_PASSWORD": "your_password",
                    "BIZFLY_REGION": "HaNoi",
                    "BIZFLY_API_URL": "https://manage.bizflycloud.vn"
                }
            }
        }
    }
    ```

## Step 4: Replace Credentials

**Important:** Replace the placeholder values with your actual credentials:

-   `your_username` → Your Bizfly Cloud username
-   `your_password` → Your Bizfly Cloud password
-   `HaNoi` → Your preferred region (optional, defaults to HaNoi)
-   `/Users/your-username/...` → Your actual absolute path (if using local binary)

## Step 5: Restart Application

After saving the configuration file:

1. **For Cursor:** Restart Cursor completely
2. **For Claude Desktop:** Restart Claude Desktop completely

## Step 6: Verify Setup

1. **Check MCP Status:**

    - In Cursor: Look for MCP indicators in the UI
    - In Claude Desktop: Check the MCP server status

2. **Test with a query:**
    - "List all my Bizfly Cloud servers"
    - "Show me available server flavors"

## Troubleshooting

### Docker Issues

1. **Verify Docker is running:**

    ```bash
    docker ps
    ```

2. **Check if image exists:**

    ```bash
    docker images bizfly-mcp-server
    ```

3. **Test Docker command manually:**
    ```bash
    docker run -i --rm \
      -e BIZFLY_USERNAME=your_username \
      -e BIZFLY_PASSWORD=your_password \
      bizfly-mcp-server:latest
    ```

### Local Binary Issues

1. **Check if binary exists and is executable:**

    ```bash
    ls -lh bizfly-mcp-server
    chmod +x bizfly-mcp-server  # if not executable
    ```

2. **Test binary manually:**

    ```bash
    export BIZFLY_USERNAME=your_username
    export BIZFLY_PASSWORD=your_password
    ./bizfly-mcp-server
    ```

3. **Check the path in config:**
    - Must be absolute path
    - Must not contain `~` (use full path)
    - Must have execute permissions

### Configuration Issues

1. **Validate JSON:**

    ```bash
    # Use a JSON validator
    cat ~/.cursor/mcp.json | python -m json.tool
    ```

2. **Check file permissions:**

    ```bash
    ls -l ~/.cursor/mcp.json
    ```

3. **Check logs:**
    - Cursor: Check Cursor logs for MCP errors
    - Claude Desktop: Check application logs

## Security Notes

⚠️ **Important Security Considerations:**

1. **Never commit credentials** to version control
2. **Use environment variables** or secure credential storage
3. **Limit file permissions:**

    ```bash
    chmod 600 ~/.cursor/mcp.json  # Read/write for owner only
    ```

4. **Consider using Docker secrets** or key management services for production

## Quick Setup Script (macOS/Linux)

Save this as `setup-mcp.sh` and run it:

```bash
#!/bin/bash

# Set your credentials here
BIZFLY_USERNAME="your_username"
BIZFLY_PASSWORD="your_password"
BIZFLY_REGION="HaNoi"
CURRENT_DIR=$(pwd)

# Create Cursor config
mkdir -p ~/.cursor
cat > ~/.cursor/mcp.json << EOF
{
  "mcpServers": {
    "bizflycloud": {
      "command": "$CURRENT_DIR/bizfly-mcp-server",
      "env": {
        "BIZFLY_USERNAME": "$BIZFLY_USERNAME",
        "BIZFLY_PASSWORD": "$BIZFLY_PASSWORD",
        "BIZFLY_REGION": "$BIZFLY_REGION"
      }
    }
  }
}
EOF

chmod 600 ~/.cursor/mcp.json
echo "✅ Cursor MCP configuration created at ~/.cursor/mcp.json"
echo "⚠️  Please replace placeholder credentials with your actual ones!"
echo "🔄 Restart Cursor to apply changes"
```

Make it executable and run:

```bash
chmod +x setup-mcp.sh
./setup-mcp.sh
```

## Next Steps

Once configured, you can:

1. Ask Cursor/Claude to list your Bizfly Cloud resources
2. Create, update, or delete cloud resources
3. Manage servers, volumes, databases, and more
4. Use natural language to interact with your cloud infrastructure

See the [README.md](README.md) for a complete list of available tools and examples.
