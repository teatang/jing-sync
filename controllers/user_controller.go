package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jing-sync/models"
	"jing-sync/services"
	"net/http"
	"strconv"
)

type UserController struct {
	BaseController
	userService *services.UserService
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := uc.userService.Create(&user); err != nil {
		uc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc.Success(c, user)
}

// GetUser 获取单个用户
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "User not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, user)
}

// GetPageUsers 分页获取用户列表
func (uc *UserController) GetPageUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi("2")
	users, err := uc.userService.GetPageList(page, size)
	if err != nil {
		uc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc.Success(c, users)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := uc.userService.Update(&user); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "User not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	id := strconv.Itoa(int(user.ID))
	if err := uc.userService.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "User not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, user)
}
