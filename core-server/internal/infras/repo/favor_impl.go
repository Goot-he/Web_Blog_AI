package repo

import (
	"context"
	stdsql "database/sql"
	"fmt"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/favoritefolder"
	"core-server/internal/infras/repo/ent/interactionfavor"
	"core-server/internal/model/entity"
)

type FavorRepo struct{ *EntClient }

func NewFavorRepo(client *EntClient) *FavorRepo { return &FavorRepo{EntClient: client} }

func (r *FavorRepo) Upsert(ctx context.Context, favor *entity.InteractionFavor) (int, error) {
	const query = `INSERT INTO interaction_favor (user_id, article_id, folder_id, status, version)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		status = IF(status = 'favored' AND version >= VALUES(version), status, VALUES(status)),
		version = IF(status = 'favored' AND version >= VALUES(version), version, VALUES(version)),
		updated_at = IF(status = 'favored' AND version >= VALUES(version), updated_at, NOW(3))`
	driver, ok := r.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return 0, fmt.Errorf("ent driver does not support ExecContext")
	}
	result, err := driver.ExecContext(ctx, query, favor.UserID, favor.ArticleID, favor.FolderID, favor.Status, favor.Version)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return 0, err
	}
	return 1, nil
}

func (r *FavorRepo) UpdateWithCondition(ctx context.Context, favor *entity.InteractionFavor) (int, error) {
	return r.DB(ctx).InteractionFavor.Update().Where(
		interactionfavor.UserIDEQ(favor.UserID),
		interactionfavor.ArticleIDEQ(favor.ArticleID),
		interactionfavor.StatusEQ(interactionfavor.StatusFavored),
		interactionfavor.VersionLTE(favor.Version),
	).SetStatus(interactionfavor.StatusNothing).SetVersion(favor.Version).Save(ctx)
}

func (r *FavorRepo) Exists(ctx context.Context, userID, articleID uint64) (bool, error) {
	return r.DB(ctx).InteractionFavor.Query().Where(
		interactionfavor.UserIDEQ(userID), interactionfavor.ArticleIDEQ(articleID), interactionfavor.StatusEQ(interactionfavor.StatusFavored),
	).Exist(ctx)
}

func (r *FavorRepo) CountByUser(ctx context.Context, userID, folderID uint64) (int64, error) {
	query := r.DB(ctx).InteractionFavor.Query().Where(interactionfavor.UserIDEQ(userID), interactionfavor.StatusEQ(interactionfavor.StatusFavored))
	if folderID > 0 {
		query.Where(interactionfavor.FolderIDEQ(folderID))
	}
	count, err := query.Count(ctx)
	return int64(count), err
}

func (r *FavorRepo) PageQueryArticles(ctx context.Context, userID, folderID uint64, offset, limit int) ([]*entity.InteractionFavor, error) {
	query := r.DB(ctx).InteractionFavor.Query().Where(interactionfavor.UserIDEQ(userID), interactionfavor.StatusEQ(interactionfavor.StatusFavored))
	if folderID > 0 {
		query.Where(interactionfavor.FolderIDEQ(folderID))
	}
	nodes, err := query.Order(ent.Desc(interactionfavor.FieldVersion), ent.Desc(interactionfavor.FieldUpdatedAt)).Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*entity.InteractionFavor, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, &entity.InteractionFavor{ID: node.ID, CreatedAt: node.CreatedAt, UpdatedAt: node.UpdatedAt, UserID: node.UserID, ArticleID: node.ArticleID, FolderID: node.FolderID, Status: node.Status.String(), Version: node.Version})
	}
	return items, nil
}

func toFavoriteFolder(node *ent.FavoriteFolder) *entity.FavoriteFolder {
	if node == nil {
		return nil
	}
	return &entity.FavoriteFolder{ID: node.ID, UserID: node.UserID, Name: node.Name, IsDefault: node.IsDefault}
}
func (r *FavorRepo) EnsureDefaultFolder(ctx context.Context, userID uint64) (*entity.FavoriteFolder, error) {
	node, err := r.DB(ctx).FavoriteFolder.Query().Where(favoritefolder.UserIDEQ(userID), favoritefolder.IsDefaultEQ(true)).Only(ctx)
	if err == nil {
		return toFavoriteFolder(node), nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}
	node, err = r.DB(ctx).FavoriteFolder.Create().SetUserID(userID).SetName("默认收藏夹").SetIsDefault(true).Save(ctx)
	if err != nil {
		return nil, err
	}
	return toFavoriteFolder(node), nil
}
func (r *FavorRepo) ListFolders(ctx context.Context, userID uint64) ([]*entity.FavoriteFolder, error) {
	nodes, err := r.DB(ctx).FavoriteFolder.Query().Where(favoritefolder.UserIDEQ(userID)).Order(ent.Desc(favoritefolder.FieldIsDefault), ent.Asc(favoritefolder.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*entity.FavoriteFolder, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, toFavoriteFolder(node))
	}
	return items, nil
}
func (r *FavorRepo) CreateFolder(ctx context.Context, userID uint64, name string) (*entity.FavoriteFolder, error) {
	node, err := r.DB(ctx).FavoriteFolder.Create().SetUserID(userID).SetName(name).Save(ctx)
	if err != nil {
		return nil, err
	}
	return toFavoriteFolder(node), nil
}
func (r *FavorRepo) GetFolder(ctx context.Context, userID, folderID uint64) (*entity.FavoriteFolder, error) {
	node, err := r.DB(ctx).FavoriteFolder.Query().Where(favoritefolder.IDEQ(folderID), favoritefolder.UserIDEQ(userID)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toFavoriteFolder(node), nil
}
func (r *FavorRepo) DeleteFolder(ctx context.Context, userID, folderID uint64) error {
	_, err := r.DB(ctx).FavoriteFolder.Delete().Where(favoritefolder.IDEQ(folderID), favoritefolder.UserIDEQ(userID), favoritefolder.IsDefaultEQ(false)).Exec(ctx)
	return err
}
func (r *FavorRepo) ListActiveByFolder(ctx context.Context, userID, folderID uint64) ([]*entity.InteractionFavor, error) {
	return r.PageQueryArticles(ctx, userID, folderID, 0, 100000)
}
func (r *FavorRepo) CountByFolder(ctx context.Context, userID, folderID uint64) (int64, error) {
	return r.CountByUser(ctx, userID, folderID)
}
