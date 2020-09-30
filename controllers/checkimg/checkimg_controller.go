package checkimg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imageshow/db"
	"imageshow/idworker"
	"imageshow/models"
	"imageshow/utils"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func IndexHandler(c *gin.Context) {
	list, err := db.GetSpiderImages()
	if err != nil {
		c.HTML(http.StatusOK, "checkimg.html", gin.H{
			"list": nil,
		})
		return
	}

	c.HTML(http.StatusOK, "checkimg.html", gin.H{
		"list": list,
	})
}

func CheckImgHandler(c *gin.Context) {
	// 判断用户是否已登录
	loginUser, ok := utils.GetSession(c, "loginUser").(models.User)
	if !ok || loginUser.Id <= 0 {
		c.JSON(200, gin.H{"code": 200, "message": "请先登录"})
		return
	}

	strId := c.Query("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	image, err := db.GetSpiderImage(id)
	if err != nil || image == nil {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	ext := path.Ext(image.Url)
	newUrl := fmt.Sprintf("/images/%c/%c/%s%s", image.MD5[0], image.MD5[1], image.MD5, ext)

	// 创建图片id
	newMD5 := utils.MD5(newUrl)
	md5Crc32 := int64(utils.CRC32(strings.ToLower(newMD5)))
	newId := idworker.GetIdByDNA(md5Crc32)
	newImg := models.Image{
		Id:            newId,
		Title:         image.Title,
		Url:           newUrl,
		MD5:           newMD5,
		Comment_Count: 0,
		Like_Count:    0,
		Created_Time:  time.Now().Unix(),
	}

	err = db.AppendImage(newImg)
	if err != nil && err.Error() != "记录已存在" {
		c.JSON(200, gin.H{"code": 0, "message": "审核失败"})
		return
	}

	err = db.CheckImg(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审核失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "审核通过"})
	}
}
