package router

import (
	"gateway/internal/facade/controller"
	"github.com/gin-gonic/gin"
)

func NewFollowRouter(v *gin.RouterGroup, c *controller.FollowController, auth gin.HandlerFunc) {
	g := v.Group("/follow")
	g.Use(auth)
	g.POST("", c.Follow)
	g.POST("/cancel", c.Unfollow)
	g.GET("/status/:id", c.Status)
	g.GET("/following/:id", c.Following)
	g.GET("/followers/:id", c.Followers)
}
