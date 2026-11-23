package facade

import (
	"gateway/internal/client/rpc"
	"gateway/internal/config"
	"gateway/internal/facade/controller"
	"gateway/internal/facade/middleware"
	"gateway/internal/facade/router"
	"gateway/internal/infras/cache"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config,
	health *controller.HealthController,
	aiChat *controller.AiChatController,
	vectorCtrl *controller.VectorController,
	userCtrl *controller.UserController,
	categoryCtrl *controller.CategoryController,
	articleCtrl *controller.ArticleController,
	commentCtrl *controller.CommentController,
	likeCtrl *controller.LikeController,
	favorCtrl *controller.FavorController,
	followCtrl *controller.FollowController,
	authCtrl *controller.AuthController,
	BlackList *cache.Blacklist,
	rpcClient *rpc.Client,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 认证中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth, BlackList, rpcClient)
	auth := authMiddleware.RequireAuth()
	optionalAuth := authMiddleware.OptionalAuth()
	role := authMiddleware.Role()

	// Prometheus metrics middleware
	//r.Use(infraProm.GinMiddleware())
	// r.Use(middleware.Recovery(), middleware.CORS())

	// 注册静态路由（含 static/upload 上传目录）
	r.Static("/static", "./static")

	// metrics endpoint (Prometheus)
	//r.GET("/metrics", gin.WrapH(infraProm.Handler()))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/ping", health.Ping)

		// 注册路由

		// RAG模块
		router.NewAIRouter(v1, aiChat, auth)
		router.NewVectorRouter(v1, vectorCtrl, auth, role)

		// 博客模块
		router.NewAuthRouter(v1, authCtrl, auth)
		router.NewUserRouter(v1, userCtrl, auth, role)
		router.NewCategoryRouter(v1, categoryCtrl, auth, role)
		router.NewArticleRouter(v1, articleCtrl, auth, optionalAuth)
		router.NewCommentRouter(v1, commentCtrl, auth, optionalAuth)
		router.NewLikeRouter(v1, likeCtrl, auth)
		router.NewFavorRouter(v1, favorCtrl, auth)
		router.NewFollowRouter(v1, followCtrl, auth)

	}

	return r
}
