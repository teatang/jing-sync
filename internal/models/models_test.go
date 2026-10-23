package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestJobModelMigrationAndFields(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&Job{}, &SyncLog{})
	require.NoError(t, err)

	minSize := int64(1024)
	maxSize := int64(1048576)
	job := &Job{
		Remark:   "测试同步作业",
		SrcPath:  "/source/photos",
		DstPath:  "/backup/photos",
		EngineId: 1,
		UserId:   1,
		Speed:    1,
		SyncMode: SyncModeAddOnly,
		Interval: 30,
		IsCron:   1,
		CronExpr: "0 2 * * *",
		MinSize:  &minSize,
		MaxSize:  &maxSize,
		Excludes: "*.tmp\n.git/*",
		Enabled:  1,
	}

	err = db.Create(job).Error
	require.NoError(t, err)
	assert.NotZero(t, job.ID)

	var queryJob Job
	err = db.First(&queryJob, job.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "测试同步作业", queryJob.Remark)
	assert.Equal(t, SyncModeAddOnly, queryJob.SyncMode)
	assert.Equal(t, uint8(1), queryJob.Enabled)
	assert.Equal(t, "0 2 * * *", queryJob.CronExpr)
	require.NotNil(t, queryJob.MinSize)
	assert.Equal(t, int64(1024), *queryJob.MinSize)
}

func TestSyncLogModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&SyncLog{})
	require.NoError(t, err)

	log := &SyncLog{
		JobId:      1,
		Status:     SyncStatusSuccess,
		ErrorMsg:   "",
		FilesAdded: 15,
		Details:    `{"files_found": 15, "files_added": 15}`,
	}
	err = db.Create(log).Error
	require.NoError(t, err)

	var queryLog SyncLog
	err = db.First(&queryLog, log.ID).Error
	require.NoError(t, err)
	assert.Equal(t, SyncStatusSuccess, queryLog.Status)
	assert.Equal(t, 15, queryLog.FilesAdded)
}
