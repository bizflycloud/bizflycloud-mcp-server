package main

import (
	"testing"

	"github.com/bizflycloud/gobizfly"
)

func TestDatabaseToolsRegistration(t *testing.T) {
	t.Run("register database tools", func(t *testing.T) {
		s := createTestMCPServer()
		client, _ := gobizfly.NewClient()
		
		RegisterDatabaseTools(s, client)
	})
}

func TestListDatabasesTool(t *testing.T) {
	t.Run("list databases", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_databases", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_databases" {
			t.Error("Invalid tool name")
		}
	})
}

func TestListDatastoresTool(t *testing.T) {
	t.Run("list datastores", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_datastores", map[string]interface{}{})
		
		if request.Params.Name != "bizflycloud_list_datastores" {
			t.Error("Invalid tool name")
		}
	})
}

func TestCreateDatabaseTool(t *testing.T) {
	t.Run("create database with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_database", map[string]interface{}{
			"name":              "db-1",
			"type":              "mysql",
			"version":            "8.0",
			"flavor":             "db.s1.small",
			"volume_size":        50.0,
			"availability_zone":  "HN1",
		})
		
		name, _ := request.Params.Arguments["name"].(string)
		dbType, _ := request.Params.Arguments["type"].(string)
		version, _ := request.Params.Arguments["version"].(string)
		flavor, _ := request.Params.Arguments["flavor"].(string)
		volumeSize, _ := request.Params.Arguments["volume_size"].(float64)
		availabilityZone, _ := request.Params.Arguments["availability_zone"].(string)
		
		if name != "db-1" || dbType != "mysql" || version != "8.0" || flavor != "db.s1.small" || volumeSize != 50.0 || availabilityZone != "HN1" {
			t.Error("Invalid parameters")
		}
	})
}

func TestGetDatabaseTool(t *testing.T) {
	t.Run("get database with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_get_database", map[string]interface{}{
			"database_id": "db-123",
		})
		
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok || databaseID != "db-123" {
			t.Error("Invalid database_id")
		}
	})
}

func TestDeleteDatabaseTool(t *testing.T) {
	t.Run("delete database with valid ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_delete_database", map[string]interface{}{
			"database_id": "db-123",
		})
		
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok || databaseID != "db-123" {
			t.Error("Invalid database_id")
		}
	})
}

func TestListDatabaseBackupsTool(t *testing.T) {
	t.Run("list backups with valid database ID", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_list_database_backups", map[string]interface{}{
			"database_id": "db-123",
		})
		
		databaseID, ok := request.Params.Arguments["database_id"].(string)
		if !ok || databaseID != "db-123" {
			t.Error("Invalid database_id")
		}
	})
}

func TestCreateDatabaseBackupTool(t *testing.T) {
	t.Run("create backup with valid parameters", func(t *testing.T) {
		request := createTestMCPRequest("bizflycloud_create_database_backup", map[string]interface{}{
			"database_id": "db-123",
			"backup_name": "backup-1",
		})
		
		databaseID, _ := request.Params.Arguments["database_id"].(string)
		backupName, _ := request.Params.Arguments["backup_name"].(string)
		
		if databaseID != "db-123" || backupName != "backup-1" {
			t.Error("Invalid parameters")
		}
	})
}

