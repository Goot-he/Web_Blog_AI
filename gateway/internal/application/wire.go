package application

import (
	"gateway/internal/application/ai-module"
	"gateway/internal/application/blog-module"
	"gateway/internal/application/interaction-module"
	"gateway/internal/application/user-module"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthService,
	ai_module.NewAiChatService,
	ai_module.NewVectorService,
	blog_module.NewArticleService,
	blog_module.NewCategoryService,
	blog_module.NewCommentService,
	interaction_module.NewLikeService,
	interaction_module.NewFavorService,
	interaction_module.NewFollowService,
	user_module.NewUserService,
	user_module.NewAuthService,
)
