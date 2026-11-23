package controller

import (
	"gateway/internal/application/interaction-module"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type FavorController struct {
	svc *interaction_module.FavorService
}

func NewFavorController(svc *interaction_module.FavorService) *FavorController {
	return &FavorController{svc: svc}
}

func (ct *FavorController) Favor(c *gin.Context) {
	var req dto.FavorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.Favor(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *FavorController) CancelFavor(c *gin.Context) {
	var req dto.FavorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.CancelFavor(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *FavorController) UserFavorList(c *gin.Context) {
	var req dto.UserFavorListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.UserFavorList(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *FavorController) ListFolders(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	result, err := ct.svc.ListFolders(c.Request.Context(), userID)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *FavorController) CreateFolder(c *gin.Context) {
	var req dto.CreateFavoriteFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.CreateFolder(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *FavorController) DeleteFolder(c *gin.Context) {
	var req dto.DeleteFavoriteFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.DeleteFolder(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}
