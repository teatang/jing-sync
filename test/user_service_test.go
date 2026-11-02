package test

import (
	"jing-sync/boot/database"
	"jing-sync/internal/models"
	"jing-sync/internal/services/db_services"
	"testing"
)

func TestUserServiceCreateAndUpdate(t *testing.T) {
	db := database.GetDB()
	us := db_services.NewUserService(db)

	// 1. 创建测试用户
	testUsername := "test_user_for_login"
	testPassword := "plain_password_123"

	// 清理旧数据（如果存在）
	db.Where("username = ?", testUsername).Delete(&models.User{})

	user := &models.User{
		Username: testUsername,
		Password: testPassword,
	}

	err := us.Create(user)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	if user.ID == 0 {
		t.Fatalf("用户创建后ID应大于0")
	}

	// 验证用明文密码登录可以成功获取用户
	loginUser, err := us.GetUserByUsernamePassword(testUsername, testPassword)
	if err != nil || loginUser == nil {
		t.Fatalf("创建后使用明文密码登录失败: %v", err)
	}

	// 2. 编辑用户信息但未提供新密码（Password 为空）
	updatedUsername := "test_user_renamed"
	updatePayload := &models.User{
		BaseModel: models.BaseModel{ID: user.ID},
		Username:  updatedUsername,
		Password:  "", // 留空，不修改密码
	}

	err = us.Update(updatePayload)
	if err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}

	// 验证用户名已更新，且原密码仍可正常登录
	loginUserAfterRename, err := us.GetUserByUsernamePassword(updatedUsername, testPassword)
	if err != nil || loginUserAfterRename == nil {
		t.Fatalf("修改用户名后使用原密码登录失败: %v", err)
	}

	// 3. 编辑用户并提供新密码
	newPassword := "brand_new_password_456"
	updatePasswordPayload := &models.User{
		BaseModel: models.BaseModel{ID: user.ID},
		Username:  updatedUsername,
		Password:  newPassword,
	}

	err = us.Update(updatePasswordPayload)
	if err != nil {
		t.Fatalf("更新密码失败: %v", err)
	}

	// 验证旧密码不可登录
	oldLoginUser, err := us.GetUserByUsernamePassword(updatedUsername, testPassword)
	if err == nil && oldLoginUser != nil && oldLoginUser.ID != 0 {
		t.Fatalf("更新密码后使用旧密码不应该能登录成功")
	}

	// 验证新密码可以正常登录
	newLoginUser, err := us.GetUserByUsernamePassword(updatedUsername, newPassword)
	if err != nil || newLoginUser == nil {
		t.Fatalf("更新密码后使用新密码登录失败: %v", err)
	}

	// 清理测试数据
	db.Where("id = ?", user.ID).Delete(&models.User{})
}
