package domain

import (
	"context"
	"core-server/internal/model/entity"
)

type FollowRepo interface {
	Upsert(context.Context, *entity.UserFollow) (int, error)
	UpdateWithCondition(context.Context, *entity.UserFollow) (int, error)
	Exists(context.Context, uint64, uint64) (bool, error)
	ListFollowing(context.Context, uint64, int, int) ([]*entity.UserFollow, int64, error)
	ListFollowers(context.Context, uint64, int, int) ([]*entity.UserFollow, int64, error)
}
type FollowRepoDomain interface {
	ITransaction
	FollowRepo
}
