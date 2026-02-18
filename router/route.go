package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(userCtrl *controller.UserController, commCtrl *controller.CommunityController) *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//加载前端
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// v1 := r.Group("/api/v1")
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// // 注册
	// v1.POST("/signup", controller.SignUpHandler)
	// // 登录
	// v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", userCtrl.SignUpHandler)
	v1.POST("/login", userCtrl.LoginHandler)

	// community routes
	community := v1.Group("/communities")
	{
		community.GET("", commCtrl.GetCommunityListHandler)
		community.GET("/:id", commCtrl.GetCommunityDetailHandler)
	}

	return r
}
