package dto

import "gateway/internal/client/rpc/core-rpc/favorpb"

type FavorRequest struct {
	UserID   uint64 `json:"-"`
	ObjectID uint64 `json:"objectId" binding:"required"`
	FolderID uint64 `json:"folderId"`
}

type UserFavorListRequest struct {
	UserID   uint64 `json:"-"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"pageSize"`
	FolderID uint64 `json:"folderId"`
}
type FavoriteFolder struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	IsDefault    bool   `json:"isDefault"`
	ArticleCount int64  `json:"articleCount"`
}
type CreateFavoriteFolderRequest struct {
	UserID uint64 `json:"-"`
	Name   string `json:"name" binding:"required,max=50"`
}
type DeleteFavoriteFolderRequest struct {
	UserID   uint64 `json:"-"`
	FolderID uint64 `json:"folderId" binding:"required"`
}
type ListFavoriteFoldersResponse struct {
	Folders []*FavoriteFolder `json:"folders"`
}

type FavorBoolResponse struct {
	Success bool `json:"success"`
}

func ToFavoriteFolder(item *favorpb.FavoriteFolder) *FavoriteFolder {
	if item == nil {
		return nil
	}
	return &FavoriteFolder{ID: item.GetId(), Name: item.GetName(), IsDefault: item.GetIsDefault(), ArticleCount: item.GetArticleCount()}
}
func ToFavoriteFolders(items []*favorpb.FavoriteFolder) []*FavoriteFolder {
	result := make([]*FavoriteFolder, 0, len(items))
	for _, item := range items {
		if folder := ToFavoriteFolder(item); folder != nil {
			result = append(result, folder)
		}
	}
	return result
}

type UserFavorListResponse struct {
	Articles []*Article `json:"articles"`
	Total    int64      `json:"total"`
}

func ToFavorBoolResponse(success bool) *FavorBoolResponse {
	return &FavorBoolResponse{Success: success}
}
func ToUserFavorListResponse(resp *favorpb.PageQueryUserFavorListResponse) *UserFavorListResponse {
	if resp == nil {
		return &UserFavorListResponse{}
	}
	return &UserFavorListResponse{Articles: ToArticles(resp.GetArticles()), Total: resp.GetTotal()}
}
