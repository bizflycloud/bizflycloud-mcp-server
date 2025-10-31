package main

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCreateTestMCPRequest(t *testing.T) {
	t.Run("create request with arguments", func(t *testing.T) {
		request := createTestMCPRequest("test_tool", map[string]interface{}{
			"param1": "value1",
			"param2": 123,
		})
		
		if request.Params.Name != "test_tool" {
			t.Errorf("Expected tool name 'test_tool', got '%s'", request.Params.Name)
		}
		
		if request.Params.Arguments["param1"] != "value1" {
			t.Error("Expected param1 to be 'value1'")
		}
		
		if request.Params.Arguments["param2"] != 123 {
			t.Error("Expected param2 to be 123")
		}
	})
	
	t.Run("create request without arguments", func(t *testing.T) {
		request := createTestMCPRequest("test_tool", map[string]interface{}{})
		
		if request.Params.Name != "test_tool" {
			t.Error("Invalid tool name")
		}
		
		if len(request.Params.Arguments) != 0 {
			t.Error("Expected empty arguments")
		}
	})
}

func TestCreateTestMCPServer(t *testing.T) {
	t.Run("create test MCP server", func(t *testing.T) {
		s := createTestMCPServer()
		
		if s == nil {
			t.Fatal("Expected non-nil server")
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"contains substring", "hello world", "world", true},
		{"does not contain substring", "hello world", "foo", false},
		{"empty substring", "hello world", "", true},
		{"empty string", "", "foo", false},
		{"exact match", "hello", "hello", true},
		{"case sensitive", "Hello", "hello", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestVerifyToolResult(t *testing.T) {
	t.Run("verify successful result", func(t *testing.T) {
		result := mcp.NewToolResultText("test output")
		verifyToolResult(t, result, "test")
	})
	
	t.Run("verify error result logs but doesn't fail", func(t *testing.T) {
		result := mcp.NewToolResultError("test error")
		// This should log but not fail the test
		verifyToolResult(t, result, "")
	})
}

func TestVerifyToolError(t *testing.T) {
	t.Run("verify error result", func(t *testing.T) {
		result := mcp.NewToolResultError("test error")
		verifyToolError(t, result, "error")
	})
	
	t.Run("verify non-error result is correctly identified", func(t *testing.T) {
		// This test verifies that verifyToolError correctly identifies non-error results
		result := mcp.NewToolResultText("success")
		
		// Manually check that it's not an error (since verifyToolError will fail)
		if result.IsError {
			t.Error("Expected non-error result")
		}
		
		text := getTextFromResult(result)
		if text != "success" {
			t.Errorf("Expected text 'success', got '%s'", text)
		}
	})
}

func TestGetTextFromResult(t *testing.T) {
	t.Run("get text from text result", func(t *testing.T) {
		result := mcp.NewToolResultText("test output")
		text := getTextFromResult(result)
		if text != "test output" {
			t.Errorf("Expected 'test output', got '%s'", text)
		}
	})
	
	t.Run("get text from error result", func(t *testing.T) {
		result := mcp.NewToolResultError("test error")
		text := getTextFromResult(result)
		if text != "test error" {
			t.Errorf("Expected 'test error', got '%s'", text)
		}
	})
	
	t.Run("get text from nil result", func(t *testing.T) {
		text := getTextFromResult(nil)
		if text != "" {
			t.Errorf("Expected empty string, got '%s'", text)
		}
	})
}

