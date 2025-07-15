package controllers

import (
	"jing-sync/models"
	"jing-sync/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser 获取单个用户
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers 获取用户列表
func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetAll()
	if err != nil {
		uc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc.Success(c, users)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.Update(id, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.userService.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
