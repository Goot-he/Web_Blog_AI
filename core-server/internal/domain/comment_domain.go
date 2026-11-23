package domain

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
)

type CommentRepo interface {
	// 评论
	CreateComment(ctx context.Context, comment *entity.Comment) (uint64, error)
	DeleteComments(ctx context.Context, ids []uint64) (int64, error)

	GetByID(ctx context.Context, id uint64) (*aggregate.CommentAggregate, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]*aggregate.CommentAggregate, error)
	ListReplyIDsByRoot(ctx context.Context, rootID uint64) ([]uint64, error)

	ListCommentByArticle(ctx context.Context, articleID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error)
	ListCommentByReply(ctx context.Context, rootID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error)

	// 评论的stat

	UpdateChildCount(ctx context.Context, rootID uint64, delta int) error
	UpdateLikeCount(ctx context.Context, id uint64, delta int) error
}

type CommentRepoDomain interface {
	ITransaction
	CommentRepo
}
