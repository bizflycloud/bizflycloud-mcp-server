package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAlertTools registers all Alert (CloudWatcher)-related tools with the MCP server
func RegisterAlertTools(s *server.MCPServer, client *gobizfly.Client) {
	// List alarms tool
	listAlarmsTool := mcp.NewTool("bizflycloud_list_alarms",
		mcp.WithDescription("List all Bizfly Cloud alarms"),
	)
	s.AddTool(listAlarmsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		alarms, err := client.CloudWatcher.Alarms().List(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list alarms: %v", err)), nil
		}

		result := "Available alarms:\n\n"
		for _, alarm := range alarms {
			result += fmt.Sprintf("Alarm: %s\n", alarm.Name)
			result += fmt.Sprintf("  ID: %s\n", alarm.ID)
			result += fmt.Sprintf("  Resource Type: %s\n", alarm.ResourceType)
			result += fmt.Sprintf("  Enable: %v\n", alarm.Enable)
			result += fmt.Sprintf("  Alert Interval: %d\n", alarm.AlertInterval)
			result += fmt.Sprintf("  Created At: %s\n", alarm.Created)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get alarm tool
	getAlarmTool := mcp.NewTool("bizflycloud_get_alarm",
		mcp.WithDescription("Get details of a Bizfly Cloud alarm"),
		mcp.WithString("alarm_id",
			mcp.Required(),
			mcp.Description("ID of the alarm"),
		),
	)
	s.AddTool(getAlarmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		alarmID, ok := request.Params.Arguments["alarm_id"].(string)
		if !ok {
			return nil, errors.New("alarm_id must be a string")
		}
		alarm, err := client.CloudWatcher.Alarms().Get(ctx, alarmID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get alarm: %v", err)), nil
		}

		result := fmt.Sprintf("Alarm Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", alarm.Name)
		result += fmt.Sprintf("ID: %s\n", alarm.ID)
		result += fmt.Sprintf("Resource Type: %s\n", alarm.ResourceType)
		result += fmt.Sprintf("Enable: %v\n", alarm.Enable)
		result += fmt.Sprintf("Alert Interval: %d\n", alarm.AlertInterval)
		result += fmt.Sprintf("Created At: %s\n", alarm.Created)
		if len(alarm.Receivers) > 0 {
			result += fmt.Sprintf("Receivers:\n")
			for _, receiver := range alarm.Receivers {
				result += fmt.Sprintf("  - %s (ID: %s)\n", receiver.Name, receiver.ReceiverID)
			}
		}
		return mcp.NewToolResultText(result), nil
	})

	// List receivers tool
	listReceiversTool := mcp.NewTool("bizflycloud_list_receivers",
		mcp.WithDescription("List all Bizfly Cloud alert receivers"),
	)
	s.AddTool(listReceiversTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		receivers, err := client.CloudWatcher.Receivers().List(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list receivers: %v", err)), nil
		}

		result := "Available receivers:\n\n"
		for _, receiver := range receivers {
			result += fmt.Sprintf("Receiver: %s\n", receiver.Name)
			result += fmt.Sprintf("  ID: %s\n", receiver.ReceiverID)
			if receiver.EmailAddress != "" {
				result += fmt.Sprintf("  Type: Email\n")
				result += fmt.Sprintf("  Email: %s\n", receiver.EmailAddress)
			} else if receiver.WebhookURL != "" {
				result += fmt.Sprintf("  Type: Webhook\n")
			} else if receiver.SMSNumber != "" {
				result += fmt.Sprintf("  Type: SMS\n")
			} else if receiver.TelegramChatID != "" {
				result += fmt.Sprintf("  Type: Telegram\n")
			}
			result += fmt.Sprintf("  Created At: %s\n", receiver.Created)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get receiver tool
	getReceiverTool := mcp.NewTool("bizflycloud_get_receiver",
		mcp.WithDescription("Get details of a Bizfly Cloud alert receiver"),
		mcp.WithString("receiver_id",
			mcp.Required(),
			mcp.Description("ID of the receiver"),
		),
	)
	s.AddTool(getReceiverTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		receiverID, ok := request.Params.Arguments["receiver_id"].(string)
		if !ok {
			return nil, errors.New("receiver_id must be a string")
		}
		receiver, err := client.CloudWatcher.Receivers().Get(ctx, receiverID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get receiver: %v", err)), nil
		}

		result := fmt.Sprintf("Receiver Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", receiver.Name)
		result += fmt.Sprintf("ID: %s\n", receiver.ReceiverID)
		if receiver.EmailAddress != "" {
			result += fmt.Sprintf("Type: Email\n")
			result += fmt.Sprintf("Email: %s\n", receiver.EmailAddress)
			result += fmt.Sprintf("Verified: %v\n", receiver.VerifiedEmailDddress)
		} else if receiver.WebhookURL != "" {
			result += fmt.Sprintf("Type: Webhook\n")
			result += fmt.Sprintf("Webhook URL: %s\n", receiver.WebhookURL)
			result += fmt.Sprintf("Verified: %v\n", receiver.VerifiedWebhookURL)
		} else if receiver.SMSNumber != "" {
			result += fmt.Sprintf("Type: SMS\n")
			result += fmt.Sprintf("SMS Number: %s\n", receiver.SMSNumber)
			result += fmt.Sprintf("Verified: %v\n", receiver.VerifiedSMSNumber)
		} else if receiver.TelegramChatID != "" {
			result += fmt.Sprintf("Type: Telegram\n")
			result += fmt.Sprintf("Telegram Chat ID: %s\n", receiver.TelegramChatID)
			result += fmt.Sprintf("Verified: %v\n", receiver.VerifiedTelegramChatID)
		}
		result += fmt.Sprintf("Created At: %s\n", receiver.Created)
		return mcp.NewToolResultText(result), nil
	})
}

