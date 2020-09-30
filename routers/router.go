package routers

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"imageshow/controllers/article"
	"imageshow/controllers/checkimg"
	"imageshow/controllers/index"
	"imageshow/controllers/login"
	"imageshow/controllers/register"
	"imageshow/models"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	gob.Register(models.User{})
	store := memstore.NewStore([]byte("ImageShow.com"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("static/html/*.html")
	router.Static("/static", "static")

	router.GET("/", index.IndexHandler)
	router.GET("/api/like", index.UpdateLikeHandler)

	router.GET("/article", article.IndexHandler)
	router.POST("/api/comment", article.PostCommentHandler)

	router.GET("/login", login.IndexHandler)
	router.POST("/api/login", login.LoginHandler)

	router.GET("/register", register.IndexHandler)
	router.POST("/api/regUser", register.RegUserHandler)

	router.GET("/checkImg", checkimg.IndexHandler)
	router.GET("/api/checkImg", checkimg.CheckImgHandler)

	return router
}
