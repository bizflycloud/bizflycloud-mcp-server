package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterKMSTools registers all KMS-related tools with the MCP server
func RegisterKMSTools(s *server.MCPServer, client *gobizfly.Client) {
	// List KMS certificates tool
	listCertificatesTool := mcp.NewTool("bizflycloud_list_kms_certificates",
		mcp.WithDescription("List all Bizfly Cloud KMS certificates"),
	)
	s.AddTool(listCertificatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certificates, err := client.KMS.Certificates().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list KMS certificates: %v", err)), nil
		}

		result := "Available KMS certificates:\n\n"
		for _, cert := range certificates {
			result += fmt.Sprintf("Certificate: %s\n", cert.Name)
			result += fmt.Sprintf("  Container ID: %s\n", cert.ContainerID)
			result += "\n"
		}
		return mcp.NewToolResultText(result), nil
	})

	// Get KMS certificate tool
	getCertificateTool := mcp.NewTool("bizflycloud_get_kms_certificate",
		mcp.WithDescription("Get details of a Bizfly Cloud KMS certificate"),
		mcp.WithString("certificate_id",
			mcp.Required(),
			mcp.Description("Container ID of the KMS certificate"),
		),
	)
	s.AddTool(getCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certificateID, ok := request.Params.Arguments["certificate_id"].(string)
		if !ok {
			return nil, errors.New("certificate_id must be a string")
		}
		cert, err := client.KMS.Certificates().Get(ctx, certificateID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get KMS certificate: %v", err)), nil
		}

		result := fmt.Sprintf("KMS Certificate Details:\n\n")
		result += fmt.Sprintf("Name: %s\n", cert.Name)
		result += fmt.Sprintf("Container ID: %s\n", cert.ContainerID)
		result += fmt.Sprintf("Certificate: %s\n", cert.Certificate)
		return mcp.NewToolResultText(result), nil
	})

	// Create KMS certificate tool
	createCertificateTool := mcp.NewTool("bizflycloud_create_kms_certificate",
		mcp.WithDescription("Create a new Bizfly Cloud KMS certificate"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the certificate container"),
		),
		mcp.WithString("certificate_name",
			mcp.Required(),
			mcp.Description("Name for the certificate"),
		),
		mcp.WithString("certificate_payload",
			mcp.Required(),
			mcp.Description("Certificate content (PEM format)"),
		),
		mcp.WithString("private_key_name",
			mcp.Required(),
			mcp.Description("Name for the private key"),
		),
		mcp.WithString("private_key_payload",
			mcp.Required(),
			mcp.Description("Private key content (PEM format)"),
		),
		mcp.WithString("private_key_passphrase_name",
			mcp.Description("Name for the private key passphrase"),
		),
		mcp.WithString("private_key_passphrase_payload",
			mcp.Description("Private key passphrase"),
		),
	)
	s.AddTool(createCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}
		certName, ok := request.Params.Arguments["certificate_name"].(string)
		if !ok {
			return nil, errors.New("certificate_name must be a string")
		}
		certPayload, ok := request.Params.Arguments["certificate_payload"].(string)
		if !ok {
			return nil, errors.New("certificate_payload must be a string")
		}
		keyName, ok := request.Params.Arguments["private_key_name"].(string)
		if !ok {
			return nil, errors.New("private_key_name must be a string")
		}
		keyPayload, ok := request.Params.Arguments["private_key_payload"].(string)
		if !ok {
			return nil, errors.New("private_key_payload must be a string")
		}

		req := &gobizfly.KMSCertificateContainerCreateRequest{
			CertContainer: gobizfly.KMSCertContainer{
				Name: name,
				Certificate: gobizfly.KMSCertificateCreateReqest{
					Name:    certName,
					Payload: certPayload,
				},
				PrivateKey: gobizfly.KMSPrivateKeyCreateReqest{
					Name:    keyName,
					Payload: keyPayload,
				},
			},
		}

		if passphraseName, ok := request.Params.Arguments["private_key_passphrase_name"].(string); ok && passphraseName != "" {
			if passphrasePayload, ok := request.Params.Arguments["private_key_passphrase_payload"].(string); ok && passphrasePayload != "" {
				req.CertContainer.PrivateKeyPassphrase = gobizfly.KMSPrivateKeyPassphraseCreateReqest{
					Name:    passphraseName,
					Payload: passphrasePayload,
				}
			}
		}

		resp, err := client.KMS.Certificates().Create(ctx, req)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create KMS certificate: %v", err)), nil
		}

		result := fmt.Sprintf("KMS certificate created successfully:\n")
		result += fmt.Sprintf("  Certificate Href: %s\n", resp.CertificateHref)
		return mcp.NewToolResultText(result), nil
	})

	// Delete KMS certificate tool
	deleteCertificateTool := mcp.NewTool("bizflycloud_delete_kms_certificate",
		mcp.WithDescription("Delete a Bizfly Cloud KMS certificate"),
		mcp.WithString("certificate_id",
			mcp.Required(),
			mcp.Description("Container ID of the KMS certificate to delete"),
		),
	)
	s.AddTool(deleteCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certificateID, ok := request.Params.Arguments["certificate_id"].(string)
		if !ok {
			return nil, errors.New("certificate_id must be a string")
		}
		err := client.KMS.Certificates().Delete(ctx, certificateID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete KMS certificate: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("KMS certificate %s deleted successfully", certificateID)), nil
	})
}

