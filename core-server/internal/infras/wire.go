package infras

import (
	"core-server/internal/domain"
	"core-server/internal/infras/cache"
	"core-server/internal/infras/clog"
	"core-server/internal/infras/repo"

	"github.com/google/wire"
)

var RepoProviderSet = wire.NewSet(
	// v2 ent 版本
	repo.NewEntClient,
	repo.NewLikeRepo,
	repo.NewFavorRepo,
	repo.NewFollowRepo,
	repo.NewUserRepo,
	repo.NewCategoryRepo,
	repo.NewArticleRepo,
	repo.NewCommentRepo,
	// todo 新操作

	wire.Bind(new(domain.LikeRepoDomain), new(*repo.LikeRepo)),
	wire.Bind(new(domain.FavorRepoDomain), new(*repo.FavorRepo)),
	wire.Bind(new(domain.FollowRepoDomain), new(*repo.FollowRepo)),
	wire.Bind(new(domain.UserRepoDomain), new(*repo.UserRepo)),
	wire.Bind(new(domain.CategoryRepoDomain), new(*repo.CategoryRepo)),
	wire.Bind(new(domain.ArticleRepoDomain), new(*repo.ArticleRepo)),
	wire.Bind(new(domain.CommentRepoDomain), new(*repo.CommentRepo)),
)

var CacheProviderSet = wire.NewSet(
	cache.NewClient,
)

var LogProviderSet = wire.NewSet(
	clog.NewLog,
)

var JobProviderSet = wire.NewSet(
	RepoProviderSet,
	CacheProviderSet,
	LogProviderSet,
)
