package article

import (
	"github.com/gin-gonic/gin"
	"imageshow/db"
	"imageshow/idworker"
	"imageshow/models"
	"imageshow/utils"
	"net/http"
	"strconv"
	"time"
)

func IndexHandler(c *gin.Context) {
	strId := c.Query("id")
	strPageIndex := c.Query("page")
	id, err := strconv.ParseInt(strId, 10, 64)
	pageIndex, err := strconv.Atoi(strPageIndex)
	if pageIndex <= 0 {
		pageIndex = 1
	}

	pagePrev := 1
	pageNext := 2
	if pageIndex > 1 {
		pagePrev = pageIndex - 1
		pageNext = pageIndex + 1
	}

	img, err := db.GetImage(id)
	if err != nil || img == nil {
		c.HTML(http.StatusOK, "article.html", gin.H{
			"img":      nil,
			"pagePrev": pagePrev,
			"pageNext": pageNext,
		})
		return
	}

	commentList, err := db.GetComments(id, pageIndex)
	c.HTML(http.StatusOK, "article.html", gin.H{
		"img":      img,
		"comments": commentList,
		"pagePrev": pagePrev,
		"pageNext": pageNext,
	})
}

// PostCommentHandler 增加图片评论
func PostCommentHandler(c *gin.Context) {
	// 判断用户是否已登录
	loginUser, ok := utils.GetSession(c, "loginUser").(models.User)
	if !ok || loginUser.Id <= 0 {
		c.JSON(200, gin.H{"code": 200, "message": "请先登录"})
		return
	}

	strImgId := c.PostForm("id")
	strContent := c.PostForm("content")
	imgId, err := strconv.ParseInt(strImgId, 10, 64)
	if err != nil || strContent == "" {
		c.JSON(200, gin.H{"code": 0, "message": "参数不正确"})
		return
	}

	// 使用imgid作为分库基因，能通过id和imgid定位分库
	id := idworker.GetIdByDNA(imgId)
	comment := models.Comment{
		Id:           id,
		UId:          loginUser.Id,
		ImgId:        imgId,
		UserName:     loginUser.Name,
		Content:      strContent,
		Created_Time: time.Now().Unix(),
	}
	err = db.AppendComment(comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "评论失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "评论成功"})
	}
}
