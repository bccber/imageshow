package register

import (
	"github.com/gin-gonic/gin"
	"imageshow/db"
	"imageshow/idworker"
	"imageshow/models"
	"imageshow/utils"
	"net/http"
	"strings"
	"time"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// RegUserHandler 注册用户
func RegUserHandler(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	remark := c.PostForm("remark")

	if name == "" || password == "" {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	// 创建用户id，后期操作能通过id和name定位分库。
	nameCrc32 := int64(utils.CRC32(strings.ToLower(name)))
	uid := idworker.GetIdByDNA(nameCrc32)
	if uid <= 0 {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	userInfo := models.User{
		Id:           uid,
		Name:         name,
		Password:     password,
		Remark:       remark,
		Created_Time: time.Now().Unix(),
	}
	err := db.AppendUser(userInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "注册失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "注册成功"})
	}
}
