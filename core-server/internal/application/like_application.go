package application

import (
	"context"
	"core-server/internal/infras/clog"
	"time"

	"core-server/internal/config"
	"core-server/internal/domain"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type LikeService struct {
	cfg         *config.Config
	log         *clog.Log
	repo        domain.LikeRepoDomain
	articleRepo domain.ArticleRepoDomain
	commentRepo domain.CommentRepoDomain
	userRepo    domain.UserRepoDomain
}

func NewLikeService(
	log *clog.Log,
	repo domain.LikeRepoDomain,
	articleRepo domain.ArticleRepoDomain,
	commentRepo domain.CommentRepoDomain,
	userRepo domain.UserRepoDomain,
	cfg *config.Config,
) (*LikeService, error) {
	return &LikeService{
		cfg:         cfg,
		repo:        repo,
		articleRepo: articleRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		log:         log,
	}, nil
}

// =====================================================================================================================
// 点赞状态

func (s *LikeService) HasThumbUp(ctx context.Context, userID uint64, objectType string, objectID uint64) (bool, error) {
	interaction, err := s.repo.QueryWithCondition(
		ctx,
		userID,
		objectType,
		objectID,
		entity.LikeStatusTypeThumbUp.String(),
	)
	if err != nil {
		return false, err
	}
	return interaction != nil, nil
}

// =====================================================================================================================
// 点赞操作

func (s *LikeService) ThumbUp(ctx context.Context, userID uint64, objectType string, objectID uint64) error {
	// 首先判断点赞是否已经存在
	exists, err := s.HasThumbUp(ctx, userID, objectType, objectID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyLiked
	}

	like := &entity.InteractionLike{
		UserID:     userID,
		ObjectType: enum.ParseObjectType(objectType),
		ObjectID:   objectID,
		Status:     entity.LikeStatusTypeThumbUp,
		Version:    time.Now().UnixMicro(),
	}

	// 开始操作点赞
	err = s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		affected, err := s.repo.Upsert(ctx, like)
		if err != nil {
			return err
		}
		// upsert 0表示无变化或未生效，1表示成功，2表示更新了原有的点赞状态（这里1，2都是1）
		if affected == 0 {
			// 无影响，则计数不需要+1
			return nil
		}

		// 更新实体获赞数
		if err := s.updateObjectLikeCount(ctx, objectType, objectID, 1); err != nil {
			return err
		}

		// 更新用户点赞数
		if err := s.userRepo.UpdateLikeCount(ctx, userID, 1); err != nil {
			return err
		}
		return s.updateArticleAuthorReceiveLikeCount(ctx, objectType, objectID, 1)
	})
	if err != nil {
		s.log.Error("Error in ThumbUp transaction")
		return err
	}

	return nil
}

func (s *LikeService) CancelThumbUp(ctx context.Context, userID uint64, objectType string, objectID uint64) error {
	exists, err := s.HasThumbUp(ctx, userID, objectType, objectID)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	like := &entity.InteractionLike{
		UserID:     userID,
		ObjectType: enum.ParseObjectType(objectType),
		ObjectID:   objectID,
		Status:     entity.LikeStatusTypeNothing,
		Version:    time.Now().UnixMicro(),
	}

	err = s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		affected, err := s.repo.UpdateWithCondition(ctx, entity.LikeStatusTypeThumbUp.String(), like)
		if err != nil {
			return err
		}
		if affected == 0 {
			// 无影响，则计数不需要-1
			return nil
		}
		if err := s.updateObjectLikeCount(ctx, objectType, objectID, -1); err != nil {
			return err
		}
		if err := s.userRepo.UpdateLikeCount(ctx, userID, -1); err != nil {
			return err
		}
		return s.updateArticleAuthorReceiveLikeCount(ctx, objectType, objectID, -1)
	})
	if err != nil {
		s.log.Error("Error in CancelThumbUp transaction")
		return err
	}

	return nil
}

func (s *LikeService) updateObjectLikeCount(ctx context.Context, objectType string, objectID uint64, delta int) error {
	switch enum.ParseObjectType(objectType) {
	case enum.ObjectTypeArticle:
		return s.articleRepo.UpdateLikeCount(ctx, objectID, delta)
	case enum.ObjectTypeComment:
		return s.commentRepo.UpdateLikeCount(ctx, objectID, delta)
	default:
		return ErrUnsupportedLikeObject
	}
}

func (s *LikeService) updateArticleAuthorReceiveLikeCount(ctx context.Context, objectType string, objectID uint64, delta int) error {
	if enum.ParseObjectType(objectType) != enum.ObjectTypeArticle {
		return nil
	}
	article, err := s.articleRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}
	if article == nil || article.Article == nil || article.Article.AuthorID == 0 {
		return nil
	}
	return s.userRepo.UpdateReceiveLikeCount(ctx, article.Article.AuthorID, int64(delta))
}

// =====================================================================================================================
// 点赞列表方面

func (s *LikeService) UserLikeList(ctx context.Context, userID uint64, objectType string, page, pageSize int) ([]*aggregate.ArticleAggregate, error) {
	offset := (page - 1) * pageSize

	likes, err := s.repo.PageQueryLikeObjects(ctx, userID, objectType, offset, pageSize)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(likes))
	for _, like := range likes {
		ids = append(ids, like.ObjectID)
	}

	articles, err := s.articleRepo.ListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
