package application

import (
	"context"
	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/entity"
	"errors"
	"time"
)

var ErrCannotFollowSelf = errors.New("cannot follow yourself")

type FollowService struct {
	repo     domain.FollowRepoDomain
	userRepo domain.UserRepoDomain
	log      *clog.Log
}

func NewFollowService(repo domain.FollowRepoDomain, userRepo domain.UserRepoDomain, log *clog.Log) *FollowService {
	return &FollowService{repo: repo, userRepo: userRepo, log: log}
}
func (s *FollowService) HasFollow(ctx context.Context, follower, following uint64) (bool, error) {
	return s.repo.Exists(ctx, follower, following)
}
func (s *FollowService) Follow(ctx context.Context, follower, following uint64) error {
	if follower == following {
		return ErrCannotFollowSelf
	}
	exists, err := s.HasFollow(ctx, follower, following)
	if err != nil || exists {
		return err
	}
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		affected, e := s.repo.Upsert(tx, &entity.UserFollow{FollowerID: follower, FollowingID: following, Status: entity.FollowStatusFollowing, Version: time.Now().UnixMicro()})
		if e != nil || affected == 0 {
			return e
		}
		if e = s.userRepo.UpdateFollowingCount(tx, follower, 1); e != nil {
			return e
		}
		return s.userRepo.UpdateFollowersCount(tx, following, 1)
	})
}
func (s *FollowService) Unfollow(ctx context.Context, follower, following uint64) error {
	exists, err := s.HasFollow(ctx, follower, following)
	if err != nil || !exists {
		return err
	}
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		affected, e := s.repo.UpdateWithCondition(tx, &entity.UserFollow{FollowerID: follower, FollowingID: following, Version: time.Now().UnixMicro()})
		if e != nil || affected == 0 {
			return e
		}
		if e = s.userRepo.UpdateFollowingCount(tx, follower, -1); e != nil {
			return e
		}
		return s.userRepo.UpdateFollowersCount(tx, following, -1)
	})
}
func (s *FollowService) ListFollowing(ctx context.Context, userID uint64, page, pageSize int) ([]*entity.User, int64, error) {
	return s.list(ctx, userID, page, pageSize, true)
}
func (s *FollowService) ListFollowers(ctx context.Context, userID uint64, page, pageSize int) ([]*entity.User, int64, error) {
	return s.list(ctx, userID, page, pageSize, false)
}
func (s *FollowService) list(ctx context.Context, userID uint64, page, pageSize int, following bool) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 10 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	var rels []*entity.UserFollow
	var total int64
	var err error
	if following {
		rels, total, err = s.repo.ListFollowing(ctx, userID, offset, pageSize)
	} else {
		rels, total, err = s.repo.ListFollowers(ctx, userID, offset, pageSize)
	}
	if err != nil {
		return nil, 0, err
	}
	users := make([]*entity.User, 0, len(rels))
	for _, rel := range rels {
		id := rel.FollowerID
		if following {
			id = rel.FollowingID
		}
		u, e := s.userRepo.GetByID(ctx, id)
		if e != nil {
			continue
		}
		users = append(users, u)
	}
	return users, total, nil
}
