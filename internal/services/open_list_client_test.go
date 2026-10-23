package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jing-sync/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Engine{}, &models.Job{}, &models.SyncLog{})
	require.NoError(t, err)
	return db
}

func TestOpenListClientOperations(t *testing.T) {
	// 启动 Mock HTTP 服务器模拟 OpenList API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/me":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "success", "data": {"username": "admin"}}`))

		case "/api/fs/list":
			var req map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&req)
			targetPath, _ := req["path"].(string)

			if targetPath == "/not_found" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code": 404, "message": "object not found"}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": 200,
				"message": "success",
				"data": {
					"total": 3,
					"content": [
						{"name": "subfolder", "size": 0, "is_dir": true, "modified": "2026-01-01T00:00:00Z"},
						{"name": "photo.jpg", "size": 2048, "is_dir": false, "modified": "2026-01-02T00:00:00Z"},
						{"name": "doc.pdf", "size": 4096, "is_dir": false, "modified": "2026-01-03T00:00:00Z"}
					]
				}
			}`))

		case "/api/fs/copy":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "success", "data": {"task_id": "task_copy_123"}}`))

		case "/api/fs/move":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "success", "data": {"task_id": "task_move_456"}}`))

		case "/api/fs/remove":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "success"}`))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": 404, "message": "not found"}`))
		}
	}))
	defer server.Close()

	db := setupTestDB(t)
	engine := &models.Engine{
		Remark: "Test Engine",
		Url:    server.URL,
		Token:  "test-token",
		UserId: 1,
	}
	require.NoError(t, db.Create(engine).Error)

	client := &OpenListClient{
		Engine: engine,
		db:     db,
	}

	// 1. 测试连通性
	err := client.TestConnection()
	require.NoError(t, err)

	// 2. 测试获取目录与文件列表（支持区分文件和目录及获取文件大小）
	files, err := client.ListFilesDetailed("/photos", 0)
	require.NoError(t, err)
	assert.Len(t, files, 3)

	var dirCount, fileCount int
	for _, f := range files {
		if f.IsDir {
			dirCount++
			assert.Equal(t, "subfolder", f.Name)
		} else {
			fileCount++
			assert.True(t, f.Size > 0)
		}
	}
	assert.Equal(t, 1, dirCount)
	assert.Equal(t, 2, fileCount)

	// 3. 测试复制文件
	copyTaskId, err := client.CopyFile("/photos", "/backup/photos", "photo.jpg")
	require.NoError(t, err)
	assert.Equal(t, "task_copy_123", copyTaskId)

	// 4. 测试移动文件
	moveTaskId, err := client.MoveFile("/photos", "/backup/photos", "doc.pdf")
	require.NoError(t, err)
	assert.Equal(t, "task_move_456", moveTaskId)

	// 5. 测试删除文件
	err = client.RemoveFile("/backup/photos", "old.pdf")
	require.NoError(t, err)
}
