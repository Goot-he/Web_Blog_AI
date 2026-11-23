package application

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserService,
	NewLikeService,
	NewFavorService,
	NewFollowService,
	NewCategoryService,
	NewArticleService,
	NewCommentService,
)
