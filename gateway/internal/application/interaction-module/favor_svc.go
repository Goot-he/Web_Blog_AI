package interaction_module

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/favorpb"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"

	"go.uber.org/zap"
)

type FavorService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewFavorService(rpcClient *rpc.Client, log *clog.Log) *FavorService {
	return &FavorService{rpc: rpcClient, log: log}
}

func (s *FavorService) Favor(ctx context.Context, req dto.FavorRequest) (*dto.FavorBoolResponse, error) {
	resp, err := s.rpc.FavorClient.Favor(ctx, &favorpb.FavorRequest{UserID: req.UserID, ObjectID: req.ObjectID, FolderID: req.FolderID})
	if err != nil {
		s.log.Error("FavorService/Favor error", zap.Error(err))
		return nil, err
	}
	return dto.ToFavorBoolResponse(resp.GetSuccess()), nil
}

func (s *FavorService) CancelFavor(ctx context.Context, req dto.FavorRequest) (*dto.FavorBoolResponse, error) {
	resp, err := s.rpc.FavorClient.CancelFavor(ctx, &favorpb.CancelFavorRequest{UserID: req.UserID, ObjectID: req.ObjectID})
	if err != nil {
		s.log.Error("FavorService/CancelFavor error", zap.Error(err))
		return nil, err
	}
	return dto.ToFavorBoolResponse(resp.GetSuccess()), nil
}

func (s *FavorService) UserFavorList(ctx context.Context, req dto.UserFavorListRequest) (*dto.UserFavorListResponse, error) {
	resp, err := s.rpc.FavorClient.PageQueryUserFavorList(ctx, &favorpb.PageQueryUserFavorListRequest{UserID: req.UserID, Page: uint32(req.Page), PageSize: uint32(req.PageSize), FolderID: req.FolderID})
	if err != nil {
		s.log.Error("FavorService/UserFavorList error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserFavorListResponse(resp), nil
}

func (s *FavorService) ListFolders(ctx context.Context, userID uint64) (*dto.ListFavoriteFoldersResponse, error) {
	resp, err := s.rpc.FavorClient.ListFavoriteFolders(ctx, &favorpb.ListFavoriteFoldersRequest{UserID: userID})
	if err != nil {
		return nil, err
	}
	return &dto.ListFavoriteFoldersResponse{Folders: dto.ToFavoriteFolders(resp.GetFolders())}, nil
}

func (s *FavorService) CreateFolder(ctx context.Context, req dto.CreateFavoriteFolderRequest) (*dto.FavoriteFolder, error) {
	resp, err := s.rpc.FavorClient.CreateFavoriteFolder(ctx, &favorpb.CreateFavoriteFolderRequest{UserID: req.UserID, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return dto.ToFavoriteFolder(resp.GetFolder()), nil
}

func (s *FavorService) DeleteFolder(ctx context.Context, req dto.DeleteFavoriteFolderRequest) (*dto.FavorBoolResponse, error) {
	resp, err := s.rpc.FavorClient.DeleteFavoriteFolder(ctx, &favorpb.DeleteFavoriteFolderRequest{UserID: req.UserID, FolderID: req.FolderID})
	if err != nil {
		return nil, err
	}
	return dto.ToFavorBoolResponse(resp.GetSuccess()), nil
}
