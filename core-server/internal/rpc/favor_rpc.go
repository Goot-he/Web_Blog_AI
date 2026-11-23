package rpc

import (
	"context"

	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
	"core-server/internal/rpc/articlepb"
	"core-server/internal/rpc/favorpb"
)

type FavorRPC struct {
	favorpb.UnimplementedFavorServiceServer
	service *application.FavorService
}

func NewFavorRPC(service *application.FavorService) *FavorRPC { return &FavorRPC{service: service} }

func (r *FavorRPC) Favor(ctx context.Context, req *favorpb.FavorRequest) (*favorpb.FavorResponse, error) {
	if err := r.service.Favor(ctx, req.GetUserID(), req.GetObjectID(), req.GetFolderID()); err != nil {
		return nil, err
	}
	return &favorpb.FavorResponse{Success: true}, nil
}

func (r *FavorRPC) CancelFavor(ctx context.Context, req *favorpb.CancelFavorRequest) (*favorpb.FavorResponse, error) {
	if err := r.service.CancelFavor(ctx, req.GetUserID(), req.GetObjectID()); err != nil {
		return nil, err
	}
	return &favorpb.FavorResponse{Success: true}, nil
}

func (r *FavorRPC) HasFavor(ctx context.Context, req *favorpb.HasFavorRequest) (*favorpb.HasFavorResponse, error) {
	favored, err := r.service.HasFavor(ctx, req.GetUserID(), req.GetObjectID())
	if err != nil {
		return nil, err
	}
	return &favorpb.HasFavorResponse{IsFavored: favored}, nil
}

func (r *FavorRPC) PageQueryUserFavorList(ctx context.Context, req *favorpb.PageQueryUserFavorListRequest) (*favorpb.PageQueryUserFavorListResponse, error) {
	articles, total, err := r.service.UserFavorList(ctx, req.GetUserID(), req.GetFolderID(), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	items := make([]*articlepb.Article, 0, len(articles))
	for _, article := range articles {
		if article != nil && article.Article != nil {
			items = append(items, toArticlePB(dto.ToArticleResponse(article)))
		}
	}
	return &favorpb.PageQueryUserFavorListResponse{Articles: items, Total: total}, nil
}

func toFavoriteFolderPB(folder *entity.FavoriteFolder) *favorpb.FavoriteFolder {
	if folder == nil {
		return nil
	}
	return &favorpb.FavoriteFolder{Id: folder.ID, Name: folder.Name, IsDefault: folder.IsDefault, ArticleCount: folder.ArticleCount}
}

func (r *FavorRPC) ListFavoriteFolders(ctx context.Context, req *favorpb.ListFavoriteFoldersRequest) (*favorpb.ListFavoriteFoldersResponse, error) {
	folders, err := r.service.ListFolders(ctx, req.GetUserID())
	if err != nil {
		return nil, err
	}
	items := make([]*favorpb.FavoriteFolder, 0, len(folders))
	for _, folder := range folders {
		items = append(items, toFavoriteFolderPB(folder))
	}
	return &favorpb.ListFavoriteFoldersResponse{Folders: items}, nil
}

func (r *FavorRPC) CreateFavoriteFolder(ctx context.Context, req *favorpb.CreateFavoriteFolderRequest) (*favorpb.CreateFavoriteFolderResponse, error) {
	folder, err := r.service.CreateFolder(ctx, req.GetUserID(), req.GetName())
	if err != nil {
		return nil, err
	}
	return &favorpb.CreateFavoriteFolderResponse{Folder: toFavoriteFolderPB(folder)}, nil
}

func (r *FavorRPC) DeleteFavoriteFolder(ctx context.Context, req *favorpb.DeleteFavoriteFolderRequest) (*favorpb.FavorResponse, error) {
	if err := r.service.DeleteFolder(ctx, req.GetUserID(), req.GetFolderID()); err != nil {
		return nil, err
	}
	return &favorpb.FavorResponse{Success: true}, nil
}
