package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type UserAPIController struct{}

func NewUserAPIController() *UserAPIController {
	return &UserAPIController{}
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserAPIController) Login(c *gin.Context) {
	// 1. gin中获取json中的username和password
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.LoginError)
		return
	}

	// 2. 查询数据库，根据用户名获取用户信息
	user, login := models.NewUser().Login(form.Username, form.Password)
	if !login {
		c.JSON(http.StatusUnauthorized, common.LoginError)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &common.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(global.GVA_CONFIG.Server.SecretKey)
	tokenString, err := token.SignedString(key)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("URL: %s, generate token error: %v", c.Request.URL.Path, err))
		c.JSON(http.StatusOK, common.TokenGenerateError)
		return
	}

	c.JSON(
		http.StatusOK,
		common.SuccessMsg{Code: common.SuccessCode, Data: map[string]string{"token": tokenString}})
}

func (u *UserAPIController) Logout(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		common.SuccessMsg{Code: common.SuccessCode, Data: "success"})
}

func (u *UserAPIController) UserInfo(c *gin.Context) {
	// 从上下文中获取存储的Claims
	claims, exists := c.Get("claims")
	if !exists {
		// 如果Claims不存在，返回错误信息
		c.JSON(http.StatusOK, common.TokenInfoGetError)
		return
	}

	// 在这里处理基于Claims的逻辑
	// 例如，你可以使用 claims 中的信息进行授权、身份验证等操作
	username := claims.(*common.Claims).Username
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: map[string]string{"name": username, "avatar": "https://bookandmusic.cn/logo.png"}})
}
