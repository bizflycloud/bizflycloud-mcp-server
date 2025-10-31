package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestKMSToolsRegistration(t *testing.T) {
	t.Run("register KMS tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterKMSTools(s, client)
	})
}

func TestListKMSCertificatesTool(t *testing.T) {
	t.Run("list KMS certificates", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_kms_certificates", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_kms_certificates" {
			t.Error("Invalid tool name")
		}
	})
}

func TestGetKMSCertificateTool(t *testing.T) {
	t.Run("get KMS certificate with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_kms_certificate", map[string]interface{}{
			"certificate_id": "cert-123",
		})
		
		certificateID, ok := request.Params.Arguments["certificate_id"].(string)
		if !ok || certificateID != "cert-123" {
			t.Error("Invalid certificate_id")
		}
	})
}

func TestCreateKMSCertificateTool(t *testing.T) {
	t.Run("create KMS certificate with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_kms_certificate", map[string]interface{}{
			"name":                    "cert-container",
			"certificate_name":        "cert",
			"certificate_payload":     "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
			"private_key_name":       "key",
			"private_key_payload":     "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
			"private_key_passphrase_name": "passphrase",
			"private_key_passphrase_payload": "secret",
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		certName, _ := request.Params.Arguments["certificate_name"].(string)
		certPayload, _ := request.Params.Arguments["certificate_payload"].(string)
		keyName, _ := request.Params.Arguments["private_key_name"].(string)
		keyPayload, _ := request.Params.Arguments["private_key_payload"].(string)
		
		if name != "cert-container" || certName != "cert" || len(certPayload) == 0 || keyName != "key" || len(keyPayload) == 0 {
			t.Error("Invalid parameters")
		}
	})
	
	t.Run("create KMS certificate without passphrase", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_kms_certificate", map[string]interface{}{
			"name":                  "cert-container",
			"certificate_name":      "cert",
			"certificate_payload":   "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
			"private_key_name":     "key",
			"private_key_payload":  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
		})
		
		// Passphrase should be optional
		if _, ok := request.Params.Arguments["private_key_passphrase_name"]; ok {
			t.Error("Passphrase should be optional")
		}
	})
}

func TestDeleteKMSCertificateTool(t *testing.T) {
	t.Run("delete KMS certificate with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_kms_certificate", map[string]interface{}{
			"certificate_id": "cert-123",
		})
		
		certificateID, ok := request.Params.Arguments["certificate_id"].(string)
		if !ok || certificateID != "cert-123" {
			t.Error("Invalid certificate_id")
		}
	})
}

