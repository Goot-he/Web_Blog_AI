package interaction_module

import (
	"context"
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/followpb"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/model/dto"
)

type FollowService struct{ rpc *rpc.Client }

func NewFollowService(r *rpc.Client) *FollowService {
	return &FollowService{
		rpc: r,
	}
}

func (s *FollowService) Follow(ctx context.Context, req dto.FollowRequest) error {
	_, e := s.rpc.FollowClient.Follow(ctx, &followpb.FollowRequest{UserId: req.UserID, TargetUserId: req.TargetUserID})
	return e
}

func (s *FollowService) Unfollow(ctx context.Context, req dto.FollowRequest) error {
	_, e := s.rpc.FollowClient.Unfollow(ctx, &followpb.FollowRequest{UserId: req.UserID, TargetUserId: req.TargetUserID})
	return e
}

func (s *FollowService) Status(ctx context.Context, req dto.FollowRequest) (bool, error) {
	r, e := s.rpc.FollowClient.HasFollow(ctx, &followpb.HasFollowRequest{UserId: req.UserID, TargetUserId: req.TargetUserID})
	if e != nil {
		return false, e
	}
	return r.GetIsFollowing(), nil
}

func (s *FollowService) ListFollowing(ctx context.Context, req dto.FollowListRequest) ([]*userpb.UserInfo, uint64, error) {
	r, e := s.rpc.FollowClient.ListFollowing(ctx, &followpb.ListFollowRequest{UserId: req.UserID, Page: int32(req.Page), PageSize: int32(req.PageSize)})
	if e != nil {
		return nil, 0, e
	}
	return r.GetUsers(), r.GetTotal(), nil
}

func (s *FollowService) ListFollowers(ctx context.Context, req dto.FollowListRequest) ([]*userpb.UserInfo, uint64, error) {
	r, e := s.rpc.FollowClient.ListFollowers(ctx, &followpb.ListFollowRequest{UserId: req.UserID, Page: int32(req.Page), PageSize: int32(req.PageSize)})
	if e != nil {
		return nil, 0, e
	}
	return r.GetUsers(), r.GetTotal(), nil
}
