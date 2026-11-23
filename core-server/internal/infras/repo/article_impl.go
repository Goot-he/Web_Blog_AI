package repo

import (
	"context"
	"core-server/internal/infras/repo/ent/user"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/enum"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/article"
	"core-server/internal/infras/repo/ent/predicate"
	"core-server/internal/model/entity"
)

type ArticleRepo struct {
	*EntClient
}

func NewArticleRepo(client *EntClient) *ArticleRepo {
	return &ArticleRepo{
		EntClient: client,
	}
}

// =====================================================================================================================
// 文章写操作

func (r *ArticleRepo) Create(ctx context.Context, value *entity.Article) error {
	// 校验作者存在，但不在文章仓储中维护用户聚合统计。
	exists, err := r.DB(ctx).User.Query().
		Where(user.IDEQ(value.AuthorID)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	create := r.DB(ctx).Article.Create().
		SetTitle(value.Title).
		SetSummary(value.Summary).
		SetContent(value.Content).
		SetCoverImage(value.CoverImage).
		SetAuthorID(value.AuthorID).
		SetCategoryID(value.CategoryID).
		SetIsTop(value.IsTop).
		SetIsPublished(value.IsPublished)
	if value.PublishedAt.Valid {
		create.SetPublishedAt(value.PublishedAt.Time)
	}
	_, err = create.Save(ctx)
	return err
}

func (r *ArticleRepo) Edit(ctx context.Context, value *entity.Article) error {
	update := r.DB(ctx).Article.Update().
		SetTitle(value.Title).
		SetSummary(value.Summary).
		SetContent(value.Content).
		SetCoverImage(value.CoverImage).
		SetCategoryID(value.CategoryID).
		SetIsPublished(value.IsPublished).
		SetIsTop(value.IsTop).
		Where(article.IDEQ(value.ID), article.AuthorIDEQ(value.AuthorID), article.DeletedAtIsNil())
	if value.IsPublished {
		update.SetPublishedAt(time.Now())
	} else {
		update.ClearPublishedAt()
	}

	affected, err := update.Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrArticleNotFound
	}
	return err
}

func (r *ArticleRepo) DeleteByID(ctx context.Context, id, authorID uint64, role string) error {
	now := time.Now()
	predicates := []predicate.Article{
		article.IDEQ(id),
		article.IsPublishedEQ(true),
		article.DeletedAtIsNil(),
	}
	if role != enum.UserRoleAdmin.String() {
		predicates = append(predicates, article.AuthorIDEQ(authorID))
	}

	affected, err := r.DB(ctx).Article.Update().
		Where(predicates...).
		SetDeletedAt(now).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

// =====================================================================================================================
// 文章读操作

func (r *ArticleRepo) List(ctx context.Context, page, pageSize int) ([]*aggregate.ArticleAggregate, error) {
	nodes, err := r.DB(ctx).Article.Query().
		Where(article.DeletedAtIsNil(), article.IsPublishedEQ(true)).
		WithCategory().
		WithUser().
		Order(ent.Desc(article.FieldID)).
		Offset(page).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityArticleAggregates(nodes), nil
}

func (r *ArticleRepo) GetByID(ctx context.Context, id uint64) (*aggregate.ArticleAggregate, error) {
	node, err := r.DB(ctx).Article.Query().
		Where(article.IDEQ(id), article.DeletedAtIsNil(), article.IsPublishedEQ(true)).
		WithCategory().
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return toEntityArticleAggregate(node), nil
}

func (r *ArticleRepo) ListByIDs(ctx context.Context, ids []uint64) ([]*aggregate.ArticleAggregate, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	nodes, err := r.DB(ctx).Article.Query().
		Where(article.IDIn(ids...), article.DeletedAtIsNil(), article.IsPublishedEQ(true)).
		WithCategory().
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityArticleAggregates(nodes), nil
}

func (r *ArticleRepo) ListByAuthor(ctx context.Context, authorID uint64, offset, limit int, isPublished *bool) ([]*aggregate.ArticleAggregate, uint64, error) {
	predicates := []predicate.Article{
		article.AuthorIDEQ(authorID),
		article.DeletedAtIsNil(),
	}
	if isPublished != nil {
		predicates = append(predicates, article.IsPublishedEQ(*isPublished))
	}

	// 查询博客实体
	nodes, err := r.DB(ctx).Article.Query().
		Where(predicates...).
		WithCategory().
		WithUser().
		Order(ent.Desc(article.FieldID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 查询计数
	total, err := r.DB(ctx).Article.Query().
		Where(predicates...).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toEntityArticleAggregates(nodes), uint64(total), nil
}

func (r *ArticleRepo) ListByCategory(ctx context.Context, categoryID uint64, offset, limit int) ([]*aggregate.ArticleAggregate, error) {
	nodes, err := r.DB(ctx).Article.Query().
		Where(article.CategoryIDEQ(categoryID), article.DeletedAtIsNil(), article.IsPublishedEQ(true)).
		WithCategory().
		WithUser().
		Order(ent.Desc(article.FieldID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityArticleAggregates(nodes), nil
}

// 查询方面

func (r *ArticleRepo) Search(ctx context.Context, name string, offset, limit int) ([]*aggregate.ArticleAggregate, error) {
	nodes, err := r.DB(ctx).Article.Query().
		Where(
			article.DeletedAtIsNil(),
			article.IsPublishedEQ(true),
			article.Or(
				article.TitleContains(name),
				article.SummaryContains(name),
				article.ContentContains(name),
			),
		).
		WithCategory().
		WithUser().
		Order(ent.Desc(article.FieldID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityArticleAggregates(nodes), nil
}

// =====================================================================================================================
// 草稿箱操作

func (r *ArticleRepo) PublishDraft(ctx context.Context, id, authorID uint64) error {
	affected, err := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.AuthorIDEQ(authorID), article.IsPublishedEQ(false), article.DeletedAtIsNil()).
		SetIsPublished(true).
		SetPublishedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (r *ArticleRepo) DeleteDraftByID(ctx context.Context, id, authorID uint64) error {
	now := time.Now()
	affected, err := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.AuthorIDEQ(authorID), article.IsPublishedEQ(false), article.DeletedAtIsNil()).
		SetDeletedAt(now).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrArticleNotFound
	}
	return nil
}

// =====================================================================================================================
// 文章计数方面

func (r *ArticleRepo) UpdateViewCount(ctx context.Context, id uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.DeletedAtIsNil(), article.IsPublishedEQ(true))
	if delta < 0 {
		query.Where(article.ViewCountGTE(-delta))
	}
	affected, err := query.AddViewCount(delta).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (r *ArticleRepo) UpdateLikeCount(ctx context.Context, id uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.DeletedAtIsNil(), article.IsPublishedEQ(true))
	if delta < 0 {
		query.Where(article.LikeCountGTE(-delta))
	}
	affected, err := query.AddLikeCount(delta).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (r *ArticleRepo) UpdateFavorCount(ctx context.Context, id uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.DeletedAtIsNil(), article.IsPublishedEQ(true))
	if delta < 0 {
		query.Where(article.FavorCountGTE(-delta))
	}
	affected, err := query.AddFavorCount(delta).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (r *ArticleRepo) UpdateCommentCount(ctx context.Context, id uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Article.Update().
		Where(article.IDEQ(id), article.DeletedAtIsNil(), article.IsPublishedEQ(true))
	if delta < 0 {
		query.Where(article.CommentCountGTE(-delta))
	}
	affected, err := query.AddCommentCount(delta).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrArticleNotFound
	}
	return nil
}
