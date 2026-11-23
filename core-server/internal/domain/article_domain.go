package domain

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
)

type ArticleRepo interface {
	// 文章写操作
	Create(ctx context.Context, article *entity.Article) error
	Edit(ctx context.Context, article *entity.Article) error
	DeleteByID(ctx context.Context, id, authorID uint64, role string) error

	// 文章读操作
	GetByID(ctx context.Context, id uint64) (*aggregate.ArticleAggregate, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]*aggregate.ArticleAggregate, error)
	List(ctx context.Context, offset, limit int) ([]*aggregate.ArticleAggregate, error)
	ListByAuthor(ctx context.Context, authorID uint64, offset, limit int, isPublished *bool) ([]*aggregate.ArticleAggregate, uint64, error)
	ListByCategory(ctx context.Context, categoryID uint64, offset, limit int) ([]*aggregate.ArticleAggregate, error)

	Search(ctx context.Context, name string, offset, limit int) ([]*aggregate.ArticleAggregate, error)

	// 草稿箱方面操作
	PublishDraft(ctx context.Context, id, authorID uint64) error
	DeleteDraftByID(ctx context.Context, id, authorID uint64) error

	// article stat 内容
	UpdateViewCount(ctx context.Context, id uint64, delta int) error
	UpdateLikeCount(ctx context.Context, id uint64, delta int) error
	UpdateFavorCount(ctx context.Context, id uint64, delta int) error
	UpdateCommentCount(ctx context.Context, id uint64, delta int) error
}

type ArticleRepoDomain interface {
	ITransaction
	ArticleRepo
}
