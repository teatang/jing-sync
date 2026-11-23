package db_services

import (
	"os"
	"strconv"
	"testing"

	"jing-sync/boot/logger"
	"jing-sync/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	logger.LoggerInit()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.User{}, &models.Engine{}, &models.Job{}, &models.SyncLog{})
	require.NoError(t, err)
	return db
}

func TestEngineServiceCRUD(t *testing.T) {
	db := setupTestDB(t)
	svc := NewEngineService(db)

	// 1. Create
	engine := &models.Engine{
		Remark: "测试同步引擎",
		Url:    "http://127.0.0.1:5244",
		Token:  "test-token",
		UserId: 1,
	}
	err := svc.Create(engine)
	require.NoError(t, err)
	assert.NotZero(t, engine.ID)
	assert.Equal(t, uint8(1), engine.Status)

	idStr := strconv.FormatUint(uint64(engine.ID), 10)

	// 2. GetByID
	fetched, err := svc.GetByID(idStr)
	require.NoError(t, err)
	assert.Equal(t, "测试同步引擎", fetched.Remark)
	assert.Equal(t, "http://127.0.0.1:5244", fetched.Url)

	// 3. Update
	fetched.Remark = "已更新的同步引擎"
	err = svc.Update(fetched)
	require.NoError(t, err)

	updated, err := svc.GetByID(idStr)
	require.NoError(t, err)
	assert.Equal(t, "已更新的同步引擎", updated.Remark)

	// 4. GetPageList
	pageList, err := svc.GetPageList(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), pageList.Pagination.Total)
	assert.Len(t, pageList.List, 1)

	// 5. Delete (Soft delete: status -> 0)
	err = svc.Delete(idStr)
	require.NoError(t, err)

	// 验证软删除后通过 GetPageList 无法查出
	pageListAfterDelete, err := svc.GetPageList(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), pageListAfterDelete.Pagination.Total)
	assert.Empty(t, pageListAfterDelete.List)
}

func TestJobServiceCRUD(t *testing.T) {
	db := setupTestDB(t)
	svc := NewJobService(db)

	job := &models.Job{
		Remark:   "照片同步定时作业",
		SrcPath:  "/source/photos",
		DstPath:  "/backup/photos",
		EngineId: 1,
		UserId:   1,
		SyncMode: models.SyncModeAddOnly,
		Enabled:  1,
	}
	err := svc.Create(job)
	require.NoError(t, err)
	assert.NotZero(t, job.ID)

	idStr := strconv.FormatUint(uint64(job.ID), 10)

	fetched, err := svc.GetByID(idStr)
	require.NoError(t, err)
	assert.Equal(t, "/source/photos", fetched.SrcPath)

	fetched.Enabled = 0
	err = svc.Update(fetched)
	require.NoError(t, err)

	updated, err := svc.GetByID(idStr)
	require.NoError(t, err)
	assert.Equal(t, uint8(0), updated.Enabled)

	err = svc.Delete(idStr)
	require.NoError(t, err)

	pageList, err := svc.GetPageList(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), pageList.Pagination.Total)
}

func TestUserServicePasswordHandling(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("data")
	})

	db := setupTestDB(t)
	svc := NewUserService(db)

	// 1. 创建用户（明文密码应自动被加密存储）
	user := &models.User{
		Username: "sync_tester",
		Password: "InitialPassword123!",
	}
	err := svc.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotEqual(t, "InitialPassword123!", user.Password)

	// 2. 使用明文密码验证登录
	loginUser, err := svc.GetUserByUsernamePassword("sync_tester", "InitialPassword123!")
	require.NoError(t, err)
	require.NotNil(t, loginUser)
	assert.Equal(t, user.ID, loginUser.ID)

	// 3. 错误密码无法登录
	wrongUser, err := svc.GetUserByUsernamePassword("sync_tester", "WrongPassword")
	assert.Error(t, err)
	assert.Equal(t, uint(0), wrongUser.ID)

	// 4. 更新用户信息但不修改密码 (Password 为空)
	updatePayload := &models.User{
		BaseModel: models.BaseModel{ID: user.ID},
		Username:  "sync_tester_renamed",
		Password:  "",
	}
	err = svc.Update(updatePayload)
	require.NoError(t, err)

	// 原密码仍可登录新用户名
	afterRename, err := svc.GetUserByUsernamePassword("sync_tester_renamed", "InitialPassword123!")
	require.NoError(t, err)
	require.NotNil(t, afterRename)
	assert.Equal(t, user.ID, afterRename.ID)

	// 5. 更新新密码
	updateNewPass := &models.User{
		BaseModel: models.BaseModel{ID: user.ID},
		Username:  "sync_tester_renamed",
		Password:  "BrandNewPassword999!",
	}
	err = svc.Update(updateNewPass)
	require.NoError(t, err)

	// 旧密码失效
	oldLogin, err := svc.GetUserByUsernamePassword("sync_tester_renamed", "InitialPassword123!")
	assert.Error(t, err)
	assert.Equal(t, uint(0), oldLogin.ID)

	// 新密码成功登录
	newLogin, err := svc.GetUserByUsernamePassword("sync_tester_renamed", "BrandNewPassword999!")
	require.NoError(t, err)
	require.NotNil(t, newLogin)
}
