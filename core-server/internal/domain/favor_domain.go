package domain

import (
	"context"
	"core-server/internal/model/entity"
)

type FavorRepo interface {
	Upsert(ctx context.Context, favor *entity.InteractionFavor) (int, error)
	UpdateWithCondition(ctx context.Context, favor *entity.InteractionFavor) (int, error)
	Exists(ctx context.Context, userID, articleID uint64) (bool, error)
	CountByUser(ctx context.Context, userID, folderID uint64) (int64, error)
	PageQueryArticles(ctx context.Context, userID, folderID uint64, offset, limit int) ([]*entity.InteractionFavor, error)
	EnsureDefaultFolder(ctx context.Context, userID uint64) (*entity.FavoriteFolder, error)
	ListFolders(ctx context.Context, userID uint64) ([]*entity.FavoriteFolder, error)
	CreateFolder(ctx context.Context, userID uint64, name string) (*entity.FavoriteFolder, error)
	GetFolder(ctx context.Context, userID, folderID uint64) (*entity.FavoriteFolder, error)
	DeleteFolder(ctx context.Context, userID, folderID uint64) error
	ListActiveByFolder(ctx context.Context, userID, folderID uint64) ([]*entity.InteractionFavor, error)
	CountByFolder(ctx context.Context, userID, folderID uint64) (int64, error)
}

type FavorRepoDomain interface {
	ITransaction
	FavorRepo
}
