package repo

import (
	"context"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/comment"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
)

type CommentRepo struct{ *EntClient }

func NewCommentRepo(client *EntClient) *CommentRepo {
	return &CommentRepo{EntClient: client}
}

func (r *CommentRepo) CreateComment(ctx context.Context, value *entity.Comment) (uint64, error) {
	node, err := r.DB(ctx).Comment.Create().
		SetArticleID(value.ArticleID).
		SetUserID(value.UserID).
		SetRootID(value.RootID).
		SetReplyToID(value.ReplyToID).
		SetContent(value.Content).
		SetLikeCount(value.LikeCount).
		SetChildCount(value.ChildCount).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return node.ID, nil
}

func (r *CommentRepo) DeleteComments(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	affected, err := r.DB(ctx).Comment.Update().
		Where(comment.IDIn(ids...), comment.DeletedAtIsNil()).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		return 0, ErrCommentNotFound
	}
	return int64(affected), nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id uint64) (*aggregate.CommentAggregate, error) {
	node, err := r.DB(ctx).Comment.Query().
		Where(comment.IDEQ(id), comment.DeletedAtIsNil()).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return aggregate.NewCommentAggregate(toEntityComment(node), toEntityUser(node.Edges.User)), nil
}

func (r *CommentRepo) ListByIDs(ctx context.Context, ids []uint64) ([]*aggregate.CommentAggregate, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	nodes, err := r.DB(ctx).Comment.Query().
		Where(comment.IDIn(ids...), comment.DeletedAtIsNil()).
		WithUser().
		WithArticle().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityCommentAggregates(nodes), nil
}

func (r *CommentRepo) ListReplyIDsByRoot(ctx context.Context, rootID uint64) ([]uint64, error) {
	return r.DB(ctx).Comment.Query().
		Where(comment.RootIDEQ(rootID), comment.DeletedAtIsNil()).
		IDs(ctx)
}

// =====================================================================================================================
// comment stat部分

func (r *CommentRepo) UpdateChildCount(ctx context.Context, rootID uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Comment.Update().
		Where(comment.IDEQ(rootID), comment.RootIDEQ(0), comment.DeletedAtIsNil())
	if delta < 0 {
		query.Where(comment.ChildCountGTE(uint32(-delta)))
	}
	affected, err := query.AddChildCount(int32(delta)).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrCommentNotFound
	}
	return nil
}

func (r *CommentRepo) UpdateLikeCount(ctx context.Context, id uint64, delta int) error {
	if delta == 0 {
		return nil
	}
	query := r.DB(ctx).Comment.Update().
		Where(comment.IDEQ(id), comment.DeletedAtIsNil())
	if delta < 0 {
		query.Where(comment.LikeCountGTE(uint32(-delta)))
	}
	affected, err := query.AddLikeCount(int32(delta)).Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		return ErrCommentNotFound
	}
	return nil
}

// =====================================================================================================================
// 列表函数

func (r *CommentRepo) ListCommentByArticle(ctx context.Context, articleID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error) {
	nodes, err := r.DB(ctx).Comment.Query().
		Where(comment.ArticleIDEQ(articleID), comment.RootIDEQ(0), comment.DeletedAtIsNil()).
		WithUser().
		WithArticle().
		Order(ent.Desc(comment.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityCommentAggregates(nodes), nil
}

func (r *CommentRepo) ListCommentByReply(ctx context.Context, rootID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error) {
	nodes, err := r.DB(ctx).Comment.Query().
		Where(comment.RootIDEQ(rootID), comment.DeletedAtIsNil()).
		WithUser().
		WithArticle().
		Order(ent.Asc(comment.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityCommentAggregates(nodes), nil
}
