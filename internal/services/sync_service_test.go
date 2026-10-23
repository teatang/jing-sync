package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jing-sync/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncServiceExecuteSyncOverlapProtection(t *testing.T) {
	db := setupTestDB(t)

	engine := &models.Engine{
		Remark: "Engine",
		Url:    "http://127.0.0.1:9999",
		Token:  "token",
		UserId: 1,
	}
	require.NoError(t, db.Create(engine).Error)

	// 重叠路径：源路径是目标路径的上级
	job := &models.Job{
		Remark:   "重叠作业",
		SrcPath:  "/data",
		DstPath:  "/data/backup",
		EngineId: engine.ID,
		UserId:   1,
		SyncMode: models.SyncModeAddOnly,
		Enabled:  1,
	}
	require.NoError(t, db.Create(job).Error)

	syncService := NewSyncService(db)
	syncLog, err := syncService.ExecuteSync(job.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "overlap")
	require.NotNil(t, syncLog)
	assert.Equal(t, models.SyncStatusFailed, syncLog.Status)
}

func TestSyncServiceAddOnlyDiffAndSync(t *testing.T) {
	var copiedFiles []string

	// 模拟 OpenList HTTP 服务
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/api/fs/list" {
			var req map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&req)
			path, _ := req["path"].(string)

			if path == "/source" {
				// 源目录有 3 个文件
				w.Write([]byte(`{
					"code": 200,
					"message": "success",
					"data": {
						"total": 3,
						"content": [
							{"name": "exist.jpg", "size": 1000, "is_dir": false},
							{"name": "new.jpg", "size": 2000, "is_dir": false},
							{"name": "ignored.tmp", "size": 500, "is_dir": false}
						]
					}
				}`))
				return
			}

			if path == "/target" {
				// 目标目录已有 exist.jpg
				w.Write([]byte(`{
					"code": 200,
					"message": "success",
					"data": {
						"total": 1,
						"content": [
							{"name": "exist.jpg", "size": 1000, "is_dir": false}
						]
					}
				}`))
				return
			}
		}

		if r.URL.Path == "/api/fs/copy" {
			var req map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&req)
			names, _ := req["names"].([]interface{})
			for _, n := range names {
				copiedFiles = append(copiedFiles, n.(string))
			}
			w.Write([]byte(`{"code": 200, "message": "success", "data": {"task_id": "copy_ok"}}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	db := setupTestDB(t)
	engine := &models.Engine{
		Remark: "TestEngine",
		Url:    server.URL,
		Token:  "token",
		UserId: 1,
	}
	require.NoError(t, db.Create(engine).Error)

	minSize := int64(800)
	job := &models.Job{
		Remark:   "增量同步作业",
		SrcPath:  "/source",
		DstPath:  "/target",
		EngineId: engine.ID,
		UserId:   1,
		SyncMode: models.SyncModeAddOnly,
		MinSize:  &minSize,
		Excludes: "*.tmp",
		Enabled:  1,
	}
	require.NoError(t, db.Create(job).Error)

	syncService := NewSyncService(db)
	syncLog, err := syncService.ExecuteSync(job.ID)

	require.NoError(t, err)
	require.NotNil(t, syncLog)
	assert.Equal(t, models.SyncStatusSuccess, syncLog.Status)

	// 验证：
	// 1. exist.jpg 在目标端已存在 -> 跳过
	// 2. ignored.tmp 被 *.tmp 规则排除 -> 跳过
	// 3. new.jpg 符合大于 800 字节且不存在于目标端 -> 复制
	assert.Equal(t, 1, syncLog.FilesAdded)
	assert.Equal(t, []string{"new.jpg"}, copiedFiles)
}
