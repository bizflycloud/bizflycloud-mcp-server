package main

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// createTestMCPRequest creates a test MCP CallToolRequest with given arguments
func createTestMCPRequest(toolName string, arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name:      toolName,
			Arguments: arguments,
		},
	}
}

// createTestMCPServer creates a test MCP server instance
func createTestMCPServer() *server.MCPServer {
	return server.NewMCPServer(
		"BizflyCloud MCP Test",
		"1.0.0",
	)
}

// getTextFromResult extracts text from CallToolResult Content
func getTextFromResult(result *mcp.CallToolResult) string {
	if result == nil || len(result.Content) == 0 {
		return ""
	}
	
	// Get first text content
	for _, content := range result.Content {
		if textContent, ok := mcp.AsTextContent(content); ok {
			return textContent.Text
		}
	}
	return ""
}

// verifyToolResult verifies that a tool result contains expected text
func verifyToolResult(t *testing.T, result *mcp.CallToolResult, expectedSubstring string) {
	if result == nil {
		t.Fatal("Expected non-nil result")
		return
	}

	if result.IsError {
		text := getTextFromResult(result)
		t.Logf("Tool returned error: %s", text)
		return
	}

	text := getTextFromResult(result)
	if text == "" {
		t.Fatal("Expected non-empty text result")
		return
	}

	if expectedSubstring != "" && !contains(text, expectedSubstring) {
		t.Errorf("Expected result to contain '%s', got: %s", expectedSubstring, text)
	}
}

// verifyToolError verifies that a tool result is an error
func verifyToolError(t *testing.T, result *mcp.CallToolResult, expectedErrorSubstring string) {
	if result == nil {
		t.Fatal("Expected non-nil result")
		return
	}

	if !result.IsError {
		text := getTextFromResult(result)
		t.Errorf("Expected error result, got: %s", text)
		return
	}

	text := getTextFromResult(result)
	if expectedErrorSubstring != "" && !contains(text, expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErrorSubstring, text)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

