package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"imageshow/db"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	if name == "" || password == "" {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	user, err := db.GetUser(name)
	if err != nil || user.Password != password {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "登录失败"})
	} else {
		// 保存session
		session := sessions.Default(c)
		session.Set("loginUser", user)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "登录成功"})
	}
}
