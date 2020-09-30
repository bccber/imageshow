package index

import (
	"github.com/gin-gonic/gin"
	"imageshow/db"
	"imageshow/models"
	"imageshow/utils"
	"net/http"
	"strconv"
)

func IndexHandler(c *gin.Context) {
	strAction := c.Query("a")
	strPageIndex := c.Query("page")
	strLastMaxId := c.Query("id1")
	strLastMinId := c.Query("id2")

	pagePrev := 1
	pageNext := 2

	lasMaxId, err := strconv.ParseInt(strLastMaxId, 10, 64)
	lasMinId, err := strconv.ParseInt(strLastMinId, 10, 64)
	pageIndex, err := strconv.Atoi(strPageIndex)

	if pageIndex <= 1 {
		strAction = ""
		lasMaxId = 0
		lasMinId = 0
	} else {
		pagePrev = pageIndex - 1
		pageNext = pageIndex + 1
	}

	list, err := db.GetImages(strAction, lasMaxId, lasMinId)
	if err != nil || list == nil || len(list) <= 0 {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"list":     nil,
			"maxId":    0,
			"minId":    0,
			"pagePrev": pagePrev,
			"pageNext": pageNext,
		})
		return
	}

	maxId := list[0].Id
	minId := list[len(list)-1].Id

	c.HTML(http.StatusOK, "index.html", gin.H{
		"list":     list,
		"maxId":    maxId,
		"minId":    minId,
		"pagePrev": pagePrev,
		"pageNext": pageNext,
	})
}

// UpdateLikeHandler 更新点赞次数
func UpdateLikeHandler(c *gin.Context) {
	// 判断用户是否已登录
	loginUser, ok := utils.GetSession(c, "loginUser").(models.User)
	if !ok || loginUser.Id <= 0 {
		c.JSON(200, gin.H{"code": 200, "message": "请先登录"})
		return
	}

	strImgId := c.Query("id")
	imgId, err := strconv.ParseInt(strImgId, 10, 64)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	err = db.UpdateLike(loginUser.Id, imgId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "点赞失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "点赞成功"})
	}
}
