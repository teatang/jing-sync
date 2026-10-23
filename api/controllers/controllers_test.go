package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jing-sync/internal/models"
	"jing-sync/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupControllerTestDB(t *testing.T) *gorm.DB {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.User{}, &models.Engine{}, &models.Job{}, &models.SyncLog{})
	require.NoError(t, err)
	return db
}

func TestEngineTestConnectionEndpoint(t *testing.T) {
	// 模拟 OpenList HTTP 服务
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/me" {
			w.Write([]byte(`{"code": 200, "message": "success"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	db := setupControllerTestDB(t)
	ec := NewEngineController(db)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", &utils.Claims{UserId: 1, Username: "admin"})
		c.Next()
	})
	r.POST("/api/engine/test", ec.TestEngine)

	// 测试请求
	reqBody, _ := json.Marshal(map[string]string{
		"url":   mockServer.URL,
		"token": "valid-token",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/engine/test", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, float64(0), resp["code"])
}

func TestJobToggleEndpoint(t *testing.T) {
	db := setupControllerTestDB(t)
	jc := NewJobController(db)

	job := &models.Job{
		Remark:   "Toggle Job",
		SrcPath:  "/src",
		DstPath:  "/dst",
		EngineId: 1,
		UserId:   1,
		Enabled:  1,
	}
	require.NoError(t, db.Create(job).Error)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", &utils.Claims{UserId: 1, Username: "admin"})
		c.Next()
	})
	r.PUT("/api/job/:id/toggle", jc.ToggleJob)

	req, _ := http.NewRequest(http.MethodPut, "/api/job/1/toggle", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedJob models.Job
	require.NoError(t, db.First(&updatedJob, 1).Error)
	assert.Equal(t, uint8(0), updatedJob.Enabled) // 状态被切换为 0
}

func TestSyncLogListAllEndpoint(t *testing.T) {
	db := setupControllerTestDB(t)
	sc := NewSyncLogController(db)

	log1 := &models.SyncLog{JobId: 1, Status: 1, FilesAdded: 5}
	log2 := &models.SyncLog{JobId: 2, Status: 0, ErrorMsg: "failed"}
	require.NoError(t, db.Create(log1).Error)
	require.NoError(t, db.Create(log2).Error)

	r := gin.New()
	r.GET("/api/sync-logs", sc.GetAllSyncLogs)

	req, _ := http.NewRequest(http.MethodGet, "/api/sync-logs?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, float64(0), resp["code"])
}

func TestGetSyncLogsByJobIdEndpoint(t *testing.T) {
	db := setupControllerTestDB(t)
	sc := NewSyncLogController(db)

	log1 := &models.SyncLog{JobId: 10, Status: 1, FilesAdded: 3}
	require.NoError(t, db.Create(log1).Error)

	r := gin.New()
	r.GET("/api/job/:id/sync-logs", sc.GetSyncLogsByJobId)

	req, _ := http.NewRequest(http.MethodGet, "/api/job/10/sync-logs?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, float64(0), resp["code"])
}
