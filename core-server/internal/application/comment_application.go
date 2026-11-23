package application

import (
	"context"
	"core-server/internal/config"
	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"

	"go.uber.org/zap"
)

type CommentService struct {
	cfg         *config.Config
	log         *clog.Log
	repo        domain.CommentRepoDomain
	articleRepo domain.ArticleRepoDomain
	userRepo    domain.UserRepoDomain
}

func NewCommentService(
	log *clog.Log,
	repo domain.CommentRepoDomain,
	articleRepo domain.ArticleRepoDomain,
	userRepo domain.UserRepoDomain,
	cfg *config.Config,
) (*CommentService, error) {
	return &CommentService{
		cfg:         cfg,
		log:         log,
		repo:        repo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}, nil
}

func (s *CommentService) CreateComment(ctx context.Context, req *dto.CreateCommentRequest) (bool, error) {
	// 1.校验这篇文章是否存在
	article, err := s.articleRepo.GetByID(ctx, req.ArticleID)
	if err != nil {
		return false, err
	}

	// 2.存在则创建评论，更新计数
	err = s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repo.CreateComment(ctx, &entity.Comment{
			ArticleID: req.ArticleID,
			ReplyToID: article.Article.AuthorID,
			UserID:    req.UserID,
			Content:   req.Content,
			RootID:    0,
		})
		if err != nil {
			return err
		}

		if err := s.articleRepo.UpdateCommentCount(ctx, req.ArticleID, 1); err != nil {
			return err
		}
		return s.userRepo.UpdateReceiveCommentCount(ctx, article.Article.AuthorID, 1)
	})
	if err != nil {
		s.log.Error("CreateComment", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) error {
	err := s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		commentAggregate, err := s.repo.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if commentAggregate == nil || commentAggregate.Comment == nil || commentAggregate.Comment.UserID != req.UserID {
			return ErrCommentNotFound
		}
		comment := commentAggregate.Comment
		if !comment.IsTopLevel() {
			return ErrCommentNotFound
		}

		article, err := s.articleRepo.GetByID(ctx, comment.ArticleID)
		if err != nil {
			return err
		}

		replyIDs, err := s.repo.ListReplyIDsByRoot(ctx, comment.ID)
		if err != nil {
			return err
		}
		ids := append([]uint64{comment.ID}, replyIDs...)

		deletedCount, err := s.repo.DeleteComments(ctx, ids)
		if err != nil {
			return err
		}

		if err := s.articleRepo.UpdateCommentCount(ctx, comment.ArticleID, -int(deletedCount)); err != nil {
			return err
		}
		return s.userRepo.UpdateReceiveCommentCount(ctx, article.Article.AuthorID, -deletedCount)
	})
	if err != nil {
		s.log.Error("DeleteComment", zap.Error(err))
		return err
	}
	return nil
}

func (s *CommentService) CreateReply(ctx context.Context, req *dto.CreateReplyRequest) (bool, error) {
	// 1.首先校验父评论是否存在
	rootComment, err := s.repo.GetByID(ctx, req.RootID)
	if err != nil || rootComment == nil || rootComment.Comment == nil || !rootComment.Comment.IsTopLevel() {
		return false, err
	}

	AuthorID := rootComment.Comment.ReplyToID
	ArticleID := rootComment.Comment.ArticleID

	// 2.校验回复评论存在
	replyComment, err := s.repo.GetByID(ctx, req.ReplyToID)
	if err != nil || replyComment == nil {
		return false, err
	}

	// 4.没问题，则创建回复，修改计数
	err = s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		_, err = s.repo.CreateComment(ctx, &entity.Comment{
			ArticleID: rootComment.Comment.ArticleID,
			UserID:    req.UserID,
			RootID:    rootComment.Comment.ID,
			ReplyToID: req.ReplyToID,
			Content:   req.Content,
		})
		if err != nil {
			return err
		}
		if err = s.repo.UpdateChildCount(ctx, rootComment.Comment.ID, 1); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCommentCount(ctx, ArticleID, 1); err != nil {
			return err
		}
		return s.userRepo.UpdateReceiveCommentCount(ctx, AuthorID, 1)
	})
	if err != nil {
		s.log.Error("CreateReply", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (s *CommentService) DeleteReply(ctx context.Context, req *dto.DeleteReplyRequest) error {
	err := s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		replyAggregate, err := s.repo.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if replyAggregate == nil || replyAggregate.Comment == nil || replyAggregate.Comment.UserID != req.UserID {
			return ErrCommentNotFound
		}
		reply := replyAggregate.Comment
		if reply.IsTopLevel() {
			return ErrCommentNotFound
		}

		article, err := s.articleRepo.GetByID(ctx, reply.ArticleID)
		if err != nil {
			return err
		}

		if _, err := s.repo.DeleteComments(ctx, []uint64{reply.ID}); err != nil {
			return err
		}
		if err := s.repo.UpdateChildCount(ctx, reply.RootID, -1); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCommentCount(ctx, reply.ArticleID, -1); err != nil {
			return err
		}
		return s.userRepo.UpdateReceiveCommentCount(ctx, article.Article.AuthorID, -1)
	})
	if err != nil {
		s.log.Error("DeleteReply", zap.Error(err))
		return err
	}
	return nil
}

// =====================================================================================================================
// 列表函数

func (s *CommentService) GetArticleComments(ctx context.Context, req *dto.GetArticleCommentsRequest) (*dto.GetArticleCommentsResponse, error) {
	offset := (req.Page - 1) * req.Size

	// 1. 获取评论列表
	comments, err := s.repo.ListCommentByArticle(ctx, req.ArticleID, offset, req.Size)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	// 2.直接返回具体的聚合根，后续在Rpc层拆分
	return &dto.GetArticleCommentsResponse{
		Comments: comments,
	}, nil
}

func (s *CommentService) GetCommentReplies(ctx context.Context, req *dto.GetCommentRepliesRequest) (*dto.GetCommentRepliesResponse, error) {
	offset := (req.Page - 1) * req.Size

	replies, err := s.repo.ListCommentByReply(ctx, req.RootID, offset, req.Size)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	replyToIDs := make([]uint64, 0, len(replies))
	for _, reply := range replies {
		if reply != nil && reply.Comment != nil && reply.Comment.ReplyToID != 0 {
			replyToIDs = append(replyToIDs, reply.Comment.ReplyToID)
		}
	}
	replyToComments, err := s.repo.ListByIDs(ctx, replyToIDs)
	if err != nil {
		s.log.Error("load reply targets error", zap.Error(err))
		return nil, err
	}
	replyToCommentMap := make(map[uint64]*aggregate.CommentAggregate, len(replyToComments))
	for _, replyToComment := range replyToComments {
		if replyToComment != nil && replyToComment.Comment != nil {
			replyToCommentMap[replyToComment.Comment.ID] = replyToComment
		}
	}
	items := make([]*dto.CommentInfo, 0, len(replies))
	for _, reply := range replies {
		if reply == nil || reply.Comment == nil {
			continue
		}
		item := &dto.CommentInfo{Comment: reply}
		if replyToComment := replyToCommentMap[reply.Comment.ReplyToID]; replyToComment != nil && replyToComment.Author != nil {
			item.ReplyToUserName = replyToComment.Author.Name
		}
		items = append(items, item)
	}

	return &dto.GetCommentRepliesResponse{
		Replies: items,
	}, nil
}
