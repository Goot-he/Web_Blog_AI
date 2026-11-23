package domain

import (
	"context"
	"core-server/internal/model/entity"
)

type LikeRepo interface {
	Upsert(ctx context.Context, like *entity.InteractionLike) (int, error)
	UpdateWithCondition(ctx context.Context, condition string, like *entity.InteractionLike) (int, error)
	QueryWithCondition(ctx context.Context, userID uint64, objectType string, objectID uint64, status string) (*entity.InteractionLike, error)
	CountUserLiked(ctx context.Context, userID uint64, objectType string) (int64, error)
	PageQueryLikeObjects(ctx context.Context, userID uint64, objectType string, offset, limit int) ([]*entity.InteractionLike, error)
}

type LikeRepoDomain interface {
	ITransaction
	LikeRepo
}
