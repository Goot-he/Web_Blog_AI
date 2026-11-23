package rpc

import (
	"context"
	"core-server/internal/application"
	"core-server/internal/model/entity"
	"core-server/internal/rpc/followpb"
	"core-server/internal/rpc/userpb"
)

type FollowRPC struct {
	followpb.UnimplementedFollowServiceServer
	service *application.FollowService
}

func NewFollowRPC(service *application.FollowService) *FollowRPC { return &FollowRPC{service: service} }
func (r *FollowRPC) Follow(ctx context.Context, req *followpb.FollowRequest) (*followpb.FollowResponse, error) {
	if err := r.service.Follow(ctx, req.GetUserId(), req.GetTargetUserId()); err != nil {
		return nil, err
	}
	return &followpb.FollowResponse{Success: true}, nil
}
func (r *FollowRPC) Unfollow(ctx context.Context, req *followpb.FollowRequest) (*followpb.FollowResponse, error) {
	if err := r.service.Unfollow(ctx, req.GetUserId(), req.GetTargetUserId()); err != nil {
		return nil, err
	}
	return &followpb.FollowResponse{Success: true}, nil
}
func (r *FollowRPC) HasFollow(ctx context.Context, req *followpb.HasFollowRequest) (*followpb.HasFollowResponse, error) {
	ok, err := r.service.HasFollow(ctx, req.GetUserId(), req.GetTargetUserId())
	return &followpb.HasFollowResponse{IsFollowing: ok}, err
}
func followUserPB(u *entity.User) *userpb.UserInfo {
	if u == nil {
		return nil
	}
	return &userpb.UserInfo{Id: u.ID, Username: u.Name, Email: u.Email, Phone: u.Phone, Avatar: u.Avatar, Sex: u.Sex.String(), Birthday: formatBirthday(u.Birthday), Signature: u.Signature, Role: u.Role.String(), Status: u.Status.String()}
}
func (r *FollowRPC) ListFollowing(ctx context.Context, req *followpb.ListFollowRequest) (*followpb.ListFollowResponse, error) {
	users, total, err := r.service.ListFollowing(ctx, req.GetUserId(), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	items := make([]*userpb.UserInfo, 0, len(users))
	for _, u := range users {
		items = append(items, followUserPB(u))
	}
	return &followpb.ListFollowResponse{Users: items, Total: uint64(total)}, nil
}
func (r *FollowRPC) ListFollowers(ctx context.Context, req *followpb.ListFollowRequest) (*followpb.ListFollowResponse, error) {
	users, total, err := r.service.ListFollowers(ctx, req.GetUserId(), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	items := make([]*userpb.UserInfo, 0, len(users))
	for _, u := range users {
		items = append(items, followUserPB(u))
	}
	return &followpb.ListFollowResponse{Users: items, Total: uint64(total)}, nil
}
