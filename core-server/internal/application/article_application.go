package application

import (
	"context"
	"core-server/internal/config"
	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type ArticleService struct {
	log      *clog.Log
	cfg      *config.Config
	ArtRepo  domain.ArticleRepoDomain
	userRepo domain.UserRepoDomain
}

func NewArticleService(
	log *clog.Log,
	cfg *config.Config,
	ArtRepo domain.ArticleRepoDomain,
	userRepo domain.UserRepoDomain,
) (*ArticleService, error) {
	return &ArticleService{
		cfg:      cfg,
		log:      log,
		ArtRepo:  ArtRepo,
		userRepo: userRepo,
	}, nil
}

// =====================================================================================================================
// 文章增删改查

func (s *ArticleService) CreateArticle(ctx context.Context, req *dto.CreateArticleRequest) error {
	a := &entity.Article{
		AuthorID:    req.AuthorID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublished,
	}

	// 用来记录这篇文章的状态
	if a.IsPublished {
		a.PublishedAt = sql.NullTime{
			Valid: true, // 表示当前这个字段有数据
			Time:  time.Now(),
		}
	}

	if err := s.ArtRepo.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.ArtRepo.Create(ctx, a); err != nil {
			return err
		}
		if a.IsPublished {
			return s.userRepo.UpdateArticleCount(ctx, a.AuthorID, 1)
		}
		return nil
	}); err != nil {
		s.log.Error("CreateArticle error", zap.Error(err))
		return err
	}

	return nil
}

func (s *ArticleService) EditorArticle(ctx context.Context, req *dto.EditorArticleRequest) error {
	a := &entity.Article{
		ID:          req.ID,
		AuthorID:    req.AuthorID,
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublished,
		CoverImage:  req.CoverImage,
	}
	if err := s.ArtRepo.Edit(ctx, a); err != nil {
		s.log.Error("EditorArticle error", zap.Error(err))
		return err
	}
	return nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, req *dto.DelArticleRequest) error {
	article, err := s.ArtRepo.GetByID(ctx, req.ID)
	if err != nil {
		s.log.Error("DeleteArticle query error", zap.Error(err))
		return err
	}

	// 删除文章和作者统计必须处于同一个事务中；权限校验仍由 ArticleRepo 负责。
	if err := s.ArtRepo.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.ArtRepo.DeleteByID(ctx, req.ID, req.UserID, req.Role); err != nil {
			return err
		}
		return s.userRepo.UpdateArticleCount(ctx, article.Article.AuthorID, -1)
	}); err != nil {
		s.log.Error("DeleteArticle error", zap.Error(err))
		return err
	}
	return nil
}

func (s *ArticleService) GetArticle(ctx context.Context, id uint64) (*dto.GetArticleResponse, error) {
	res, err := s.ArtRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("GetArticle info", zap.Error(err))
		return nil, err
	}

	// 浏览量+1，有问题也没事
	err = s.ArtRepo.UpdateViewCount(ctx, id, 1)
	if err != nil {
		// 有问题记录，但是不中断
		s.log.Warn("GetArticle update view count error", zap.Error(err))
	}

	return dto.ToGetArticleResponse(res), nil
}

// =====================================================================================================================
// 草稿箱 (草稿箱对了一个isPublished的字段, 这个字段为false就是草稿箱, true就是发布的文章)

func (s *ArticleService) PublishDraft(ctx context.Context, req *dto.PublishDraftRequest) error {
	// 发布草稿箱内的文章
	if err := s.ArtRepo.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.ArtRepo.PublishDraft(ctx, req.ID, req.AuthorID); err != nil {
			return err
		}
		return s.userRepo.UpdateArticleCount(ctx, req.AuthorID, 1)
	}); err != nil {
		s.log.Error("PublishDraft error", zap.Error(err))
		return err
	}
	return nil
}

func (s *ArticleService) DeleteDraft(ctx context.Context, req *dto.DeleteDraftRequest) error {
	if err := s.ArtRepo.DeleteDraftByID(ctx, req.ID, req.AuthorID); err != nil {
		s.log.Error("DeleteDraft error", zap.Error(err))
		return err
	}
	return nil
}

// =====================================================================================================================
// 列表函数

// ListArticles 公开列表
func (s *ArticleService) ListArticles(ctx context.Context, req *dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, err := s.ArtRepo.List(ctx, offset, req.PageSize)
	if err != nil {
		s.log.Error("ListArticles error", zap.Error(err))
		return nil, err
	}

	return dto.ToListArticlesResponse(articles), nil
}

// ListMyArticles 个人列表
func (s *ArticleService) ListMyArticles(ctx context.Context, req *dto.ListMyArticlesRequest) (*dto.ListMyArticlesResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.ArtRepo.ListByAuthor(ctx, req.AuthorID, offset, req.PageSize, req.IsPublished)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return dto.ToListMyArticlesResponse(articles, total), nil
}

// =====================================================================================================================
// 查询函数

func (s *ArticleService) SearchArticlesByCategory(ctx context.Context, req *dto.ListArticlesByCategoryRequest) (*dto.ListArticlesByCategoryResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, err := s.ArtRepo.ListByCategory(ctx, req.CategoryID, offset, req.PageSize)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	return dto.ToListArticlesByCategoryResponse(articles), nil
}

func (s *ArticleService) SearchArticles(ctx context.Context, req *dto.SearchArticlesRequest) (*dto.SearchArticlesResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, err := s.ArtRepo.Search(ctx, req.Query, offset, req.PageSize)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	return dto.ToSearchArticlesResponse(articles), nil
}
