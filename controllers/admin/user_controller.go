package admin

import (
	"net/http"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

var users []models.User

func (u *UserController) Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		// 处理 GET 请求，返回登录页面
		c.HTML(http.StatusOK, "admin/login.html", pongo2.Context{
			"msg":            "",
			"serverName":     global.GVA_CONFIG.Server.ServerName,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	} else if c.Request.Method == "POST" {
		// 处理 POST 请求，获取表单提交的用户名和密码
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 查询数据库，根据用户名获取用户信息
		user, login := models.NewUser().Login(username, password)
		if login == false {
			// 用户不存在，或者其他数据库错误处理
			c.HTML(http.StatusOK, "admin/login.html", pongo2.Context{
				"msg":            "用户名或密码错误",
				csrf.TemplateTag: csrf.TemplateField(c.Request),
			})
			return
		}

		// 如果用户名和密码验证成功，设置会话信息表示用户已登录
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Set("is_logged_in", true)
		session.Save()

		// 重定向到首页或其他受保护的页面
		c.Redirect(http.StatusFound, "/admin/")
	}
}

func (u *UserController) Logout(c *gin.Context) {
	// 获取 session
	session := sessions.Default(c)

	// 删除用户信息
	session.Delete("user_id")
	session.Delete("is_logged_in")

	// 保存 session
	if err := session.Save(); err != nil {
		global.GVA_LOG.Warn("Failed to save session", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false})
		return
	}

	// 重定向到登录页面或其他页面
	c.Redirect(http.StatusSeeOther, "login")
}

func (u *UserController) UserSetting(c *gin.Context) {
	session := sessions.Default(c)
	isLogged := session.Get("is_logged_in")
	userId := session.Get("user_id")
	if isLogged != true {
		global.GVA_LOG.Warn("User is not logged in")
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "msg": "用户未登录"})
		return
	}

	if userId == nil {
		global.GVA_LOG.Warn("User ID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户ID不存在"})
		return
	}
	i, ok := userId.(uint)
	if !ok {
		global.GVA_LOG.Warn("Failed to convert user ID to integer")
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户ID不合法"})
		return
	}
	user, err := models.NewUser().FindByUserId(int(i))
	if err != nil {
		global.GVA_LOG.Warn("Failed to finc user by id", userId)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户不存在"})
		return
	}
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/user/user-setting.html", pongo2.Context{
			"user":      user,
			"csrfToken": csrf.Token(c.Request),
		})
		return
	} else if c.Request.Method == "POST" {
		// 处理 POST 请求，获取表单提交的用户名和密码
		username := c.PostForm("username")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		remark := c.PostForm("remark")

		if username == "" || phone == "" {
			global.GVA_LOG.Warn("username and phone is required")
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户名和手机号是必须的"})
			return
		}

		updates := map[string]interface{}{
			"username":    username,
			"email":       email,
			"phone":       phone,
			"description": remark,
		}
		if err := user.Update(updates); err != nil {
			global.GVA_LOG.Warn("Failed to save user setting", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户信息更新失败"})
			return
		}
		// 重定向到登录页面或其他页面
		return
	}
}

func (u *UserController) UserPassword(c *gin.Context) {
	session := sessions.Default(c)
	isLogged := session.Get("is_logged_in")
	userId := session.Get("user_id")
	if isLogged != true {
		global.GVA_LOG.Warn("User is not logged in")
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "msg": "用户未登录"})
		return
	}

	if userId == nil {
		global.GVA_LOG.Warn("User ID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户ID不存在"})
		return
	}
	i, ok := userId.(uint)
	if !ok {
		global.GVA_LOG.Warn("Failed to convert user ID to integer")
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户ID不合法"})
		return
	}
	user, err := models.NewUser().FindByUserId(int(i))
	if err != nil {
		global.GVA_LOG.Warn("Failed to finc user by id", userId)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "用户不存在"})
		return
	}
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/user/user-password.html", pongo2.Context{
			"user":      user,
			"csrfToken": csrf.Token(c.Request),
		})
		return
	} else if c.Request.Method == "POST" {
		// 处理 POST 请求，获取表单提交的用户名和密码
		oldPassword := c.PostForm("old_password")
		newPassword := c.PostForm("new_password")
		againPassword := c.PostForm("again_password")

		if oldPassword == "" || newPassword == "" || againPassword == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "密码是必须的"})
			return
		}
		if newPassword != againPassword {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "密码不一致"})
			return
		}
		if newPassword == oldPassword {
			global.GVA_LOG.Warn("The password no change")
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "密码没有改变"})
			return
		}
		if user.CheckPassword(oldPassword) == false {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "原始密码错误"})
			return
		}

		updates := map[string]interface{}{
			"password": newPassword,
		}
		if err := user.Update(updates); err != nil {
			global.GVA_LOG.Warn("Failed to update user password", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "更新密码失败"})
			return
		}
		// 重定向到登录页面或其他页面
		return
	}
}
