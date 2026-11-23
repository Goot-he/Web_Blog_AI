package rpc

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAuthRPC,
	NewUserRPC,
	NewLikeRPC,
	NewFavorRPC,
	NewFollowRPC,
	NewCategoryRPC,
	NewArticleRPC,
	NewCommentRPC,
	NewServer,
)
