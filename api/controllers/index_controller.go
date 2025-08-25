package controllers

import (
	"jing-sync/internal/models"
	"jing-sync/internal/services/db_services"
	"jing-sync/internal/utils"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IndexController struct {
	BaseController
	userService *db_services.UserService
}

func NewIndexController(db *gorm.DB) *IndexController {
	return &IndexController{
		userService: db_services.NewUserService(db),
	}
}

// 用户登陆
func (uc *IndexController) Login(c *gin.Context) {
	var postUser models.User
	if err := c.ShouldBindJSON(&postUser); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if postUser.Username == "" || postUser.Password == "" {
		uc.Error(c, http.StatusBadRequest, "用户名或者密码不能为空")
		return
	}

	user, err := uc.userService.GetUserByUsernamePassword(postUser.Username, postUser.Password)
	if err != nil {
		uc.Error(c, http.StatusInternalServerError, utils.GetI18nMsg("error_msg.username_password_incorrect", c))
		return
	}

	// 创建Token声明
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &utils.Claims{
		Username: user.Username,
		UserId:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "access",
		},
	}

	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, s_err := utils.GetSecretKey()
	if s_err != nil {
		uc.Error(c, http.StatusInternalServerError, s_err.Error())
		return
	}
	tokenString, t_err := token.SignedString([]byte(jwtSecret))
	if t_err != nil {
		uc.Error(c, http.StatusInternalServerError, t_err.Error())
		return
	}

	uc.Success(c, map[string]interface{}{"token": tokenString, "expires_at": expirationTime})
}
