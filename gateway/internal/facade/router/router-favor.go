package router

import (
	"gateway/internal/facade/controller"
	"github.com/gin-gonic/gin"
)

func NewFavorRouter(v *gin.RouterGroup, ct *controller.FavorController, authMiddleware gin.HandlerFunc) {
	favor := v.Group("/favor")
	favor.Use(authMiddleware)
	{
		favor.POST("", ct.Favor)
		favor.POST("/cancel", ct.CancelFavor)
		favor.POST("/list", ct.UserFavorList)
		favor.GET("/folders", ct.ListFolders)
		favor.POST("/folders", ct.CreateFolder)
		favor.POST("/folders/delete", ct.DeleteFolder)
	}
}
