package repo

import (
	"context"
	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/userfollow"
	"core-server/internal/model/entity"
	stdsql "database/sql"
	"fmt"
)

type FollowRepo struct{ *EntClient }

func NewFollowRepo(client *EntClient) *FollowRepo { return &FollowRepo{client} }
func (r *FollowRepo) Upsert(ctx context.Context, follow *entity.UserFollow) (int, error) {
	const query = `INSERT INTO user_follow (follower_id, following_id, status, version) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = IF(status = 'following' AND version >= VALUES(version), status, VALUES(status)), version = IF(status = 'following' AND version >= VALUES(version), version, VALUES(version)), updated_at = IF(status = 'following' AND version >= VALUES(version), updated_at, NOW(3))`
	driver, ok := r.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return 0, fmt.Errorf("ent driver does not support ExecContext")
	}
	result, err := driver.ExecContext(ctx, query, follow.FollowerID, follow.FollowingID, follow.Status, follow.Version)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return 0, err
	}
	return 1, nil
}
func (r *FollowRepo) UpdateWithCondition(ctx context.Context, follow *entity.UserFollow) (int, error) {
	return r.DB(ctx).UserFollow.Update().Where(userfollow.FollowerIDEQ(follow.FollowerID), userfollow.FollowingIDEQ(follow.FollowingID), userfollow.StatusEQ(userfollow.StatusFollowing), userfollow.VersionLTE(follow.Version)).SetStatus(userfollow.StatusNothing).SetVersion(follow.Version).Save(ctx)
}
func (r *FollowRepo) Exists(ctx context.Context, followerID, followingID uint64) (bool, error) {
	return r.DB(ctx).UserFollow.Query().Where(userfollow.FollowerIDEQ(followerID), userfollow.FollowingIDEQ(followingID), userfollow.StatusEQ(userfollow.StatusFollowing)).Exist(ctx)
}
func toEntityUserFollow(node *ent.UserFollow) *entity.UserFollow {
	if node == nil {
		return nil
	}
	return &entity.UserFollow{ID: node.ID, CreatedAt: node.CreatedAt, UpdatedAt: node.UpdatedAt, FollowerID: node.FollowerID, FollowingID: node.FollowingID, Status: string(node.Status), Version: node.Version}
}
func toEntityUserFollows(nodes []*ent.UserFollow) []*entity.UserFollow {
	result := make([]*entity.UserFollow, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, toEntityUserFollow(node))
	}
	return result
}
func (r *FollowRepo) ListFollowing(ctx context.Context, userID uint64, offset, limit int) ([]*entity.UserFollow, int64, error) {
	query := r.DB(ctx).UserFollow.Query().Where(userfollow.FollowerIDEQ(userID), userfollow.StatusEQ(userfollow.StatusFollowing))
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	nodes, err := query.Order(ent.Desc(userfollow.FieldUpdatedAt)).Offset(offset).Limit(limit).All(ctx)
	return toEntityUserFollows(nodes), int64(total), err
}
func (r *FollowRepo) ListFollowers(ctx context.Context, userID uint64, offset, limit int) ([]*entity.UserFollow, int64, error) {
	query := r.DB(ctx).UserFollow.Query().Where(userfollow.FollowingIDEQ(userID), userfollow.StatusEQ(userfollow.StatusFollowing))
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	nodes, err := query.Order(ent.Desc(userfollow.FieldUpdatedAt)).Offset(offset).Limit(limit).All(ctx)
	return toEntityUserFollows(nodes), int64(total), err
}
