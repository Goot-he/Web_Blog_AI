package application

import (
	"context"
	"time"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
)

type FavorService struct {
	log         *clog.Log
	repo        domain.FavorRepoDomain
	articleRepo domain.ArticleRepoDomain
	userRepo    domain.UserRepoDomain
}

func NewFavorService(log *clog.Log, repo domain.FavorRepoDomain, articleRepo domain.ArticleRepoDomain, userRepo domain.UserRepoDomain) *FavorService {
	return &FavorService{log: log, repo: repo, articleRepo: articleRepo, userRepo: userRepo}
}

func (s *FavorService) HasFavor(ctx context.Context, userID, articleID uint64) (bool, error) {
	return s.repo.Exists(ctx, userID, articleID)
}

func (s *FavorService) Favor(ctx context.Context, userID, articleID, folderID uint64) error {
	exists, err := s.HasFavor(ctx, userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		if folderID == 0 {
			folder, err := s.repo.EnsureDefaultFolder(tx, userID)
			if err != nil {
				return err
			}
			folderID = folder.ID
		} else if folder, err := s.repo.GetFolder(tx, userID, folderID); err != nil || folder == nil {
			if err != nil {
				return err
			}
			return ErrFavoriteFolderNotFound
		}
		affected, err := s.repo.Upsert(tx, &entity.InteractionFavor{UserID: userID, ArticleID: articleID, FolderID: folderID, Status: entity.FavorStatusFavored, Version: time.Now().UnixMicro()})
		if err != nil || affected == 0 {
			return err
		}
		if err := s.articleRepo.UpdateFavorCount(tx, articleID, 1); err != nil {
			return err
		}
		if err := s.userRepo.UpdateFavorCount(tx, userID, 1); err != nil {
			return err
		}
		article, err := s.articleRepo.GetByID(tx, articleID)
		if err != nil || article == nil || article.Article == nil {
			return err
		}
		return s.userRepo.UpdateReceiveFavorCount(tx, article.Article.AuthorID, 1)
	})
}

func (s *FavorService) CancelFavor(ctx context.Context, userID, articleID uint64) error {
	exists, err := s.HasFavor(ctx, userID, articleID)
	if err != nil || !exists {
		return err
	}
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		affected, err := s.repo.UpdateWithCondition(tx, &entity.InteractionFavor{UserID: userID, ArticleID: articleID, Status: entity.FavorStatusNothing, Version: time.Now().UnixMicro()})
		if err != nil || affected == 0 {
			return err
		}
		if err := s.articleRepo.UpdateFavorCount(tx, articleID, -1); err != nil {
			return err
		}
		if err := s.userRepo.UpdateFavorCount(tx, userID, -1); err != nil {
			return err
		}
		article, err := s.articleRepo.GetByID(tx, articleID)
		if err != nil || article == nil || article.Article == nil {
			return err
		}
		return s.userRepo.UpdateReceiveFavorCount(tx, article.Article.AuthorID, -1)
	})
}

func (s *FavorService) UserFavorList(ctx context.Context, userID, folderID uint64, page, pageSize int) ([]*aggregate.ArticleAggregate, int64, error) {
	page, pageSize = Page(page), Size(pageSize)
	total, err := s.repo.CountByUser(ctx, userID, folderID)
	if err != nil {
		return nil, 0, err
	}
	favors, err := s.repo.PageQueryArticles(ctx, userID, folderID, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ids := make([]uint64, 0, len(favors))
	for _, favor := range favors {
		ids = append(ids, favor.ArticleID)
	}
	articles, err := s.articleRepo.ListByIDs(ctx, ids)
	return articles, total, err
}

func (s *FavorService) ListFolders(ctx context.Context, userID uint64) ([]*entity.FavoriteFolder, error) {
	if _, err := s.repo.EnsureDefaultFolder(ctx, userID); err != nil {
		return nil, err
	}
	folders, err := s.repo.ListFolders(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, folder := range folders {
		folder.ArticleCount, err = s.repo.CountByFolder(ctx, userID, folder.ID)
		if err != nil {
			return nil, err
		}
	}
	return folders, nil
}

func (s *FavorService) CreateFolder(ctx context.Context, userID uint64, name string) (*entity.FavoriteFolder, error) {
	if _, err := s.repo.EnsureDefaultFolder(ctx, userID); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, ErrFavoriteFolderName
	}
	return s.repo.CreateFolder(ctx, userID, name)
}

func (s *FavorService) DeleteFolder(ctx context.Context, userID, folderID uint64) error {
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		folder, err := s.repo.GetFolder(tx, userID, folderID)
		if err != nil {
			return err
		}
		if folder == nil {
			return ErrFavoriteFolderNotFound
		}
		if folder.IsDefault {
			return ErrDefaultFavoriteFolder
		}
		favors, err := s.repo.ListActiveByFolder(tx, userID, folderID)
		if err != nil {
			return err
		}
		for _, favor := range favors {
			affected, err := s.repo.UpdateWithCondition(tx, &entity.InteractionFavor{UserID: userID, ArticleID: favor.ArticleID, Status: entity.FavorStatusNothing, Version: time.Now().UnixMicro()})
			if err != nil {
				return err
			}
			if affected == 0 {
				continue
			}
			if err := s.articleRepo.UpdateFavorCount(tx, favor.ArticleID, -1); err != nil {
				return err
			}
			if err := s.userRepo.UpdateFavorCount(tx, userID, -1); err != nil {
				return err
			}
			article, err := s.articleRepo.GetByID(tx, favor.ArticleID)
			if err != nil || article == nil || article.Article == nil {
				return err
			}
			if err := s.userRepo.UpdateReceiveFavorCount(tx, article.Article.AuthorID, -1); err != nil {
				return err
			}
		}
		return s.repo.DeleteFolder(tx, userID, folderID)
	})
}
