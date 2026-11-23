package controller

import (
	"gateway/internal/application/interaction-module"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	svc *interaction_module.FollowService
}

func NewFollowController(s *interaction_module.FollowService) *FollowController {
	return &FollowController{svc: s}
}
func (c *FollowController) Follow(ctx *gin.Context) {
	uid, ok := middleware.GetUserID(ctx)
	if !ok {
		reponse.Unauthorized(ctx)
		return
	}
	var req dto.FollowRequest
	if ctx.ShouldBindJSON(&req) != nil {
		reponse.StatusBadRequest(ctx)
		return
	}
	req.UserID = uid
	if err := c.svc.Follow(ctx, req); err != nil {
		reponse.InternalServerError(ctx, err.Error())
		return
	}
	reponse.Success(ctx, gin.H{"success": true})
}

func (c *FollowController) Unfollow(ctx *gin.Context) {
	uid, ok := middleware.GetUserID(ctx)
	if !ok {
		reponse.Unauthorized(ctx)
		return
	}
	var req dto.FollowRequest
	if ctx.ShouldBindJSON(&req) != nil {
		reponse.StatusBadRequest(ctx)
		return
	}
	req.UserID = uid
	if err := c.svc.Unfollow(ctx, req); err != nil {
		reponse.InternalServerError(ctx, err.Error())
		return
	}
	reponse.Success(ctx, gin.H{"success": true})
}

func (c *FollowController) Status(ctx *gin.Context) {
	uid, ok := middleware.GetUserID(ctx)
	if !ok {
		reponse.Unauthorized(ctx)
		return
	}
	target, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	v, e := c.svc.Status(ctx, dto.FollowRequest{UserID: uid, TargetUserID: target})
	if e != nil {
		reponse.InternalServerError(ctx, e.Error())
		return
	}
	reponse.Success(ctx, gin.H{"isFollowing": v})
}

func (c *FollowController) list(ctx *gin.Context, followers bool) {
	currentUserID, ok := middleware.GetUserID(ctx)
	if !ok {
		reponse.Unauthorized(ctx)
		return
	}
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if id == 0 {
		reponse.StatusBadRequest(ctx)
		return
	}
	if id != currentUserID {
		reponse.Forbidden(ctx)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	req := dto.FollowListRequest{UserID: id, Page: page, PageSize: size}
	var u interface{}
	var total uint64
	var e error
	if followers {
		u, total, e = c.svc.ListFollowers(ctx, req)
	} else {
		u, total, e = c.svc.ListFollowing(ctx, req)
	}
	if e != nil {
		reponse.InternalServerError(ctx, e.Error())
		return
	}
	reponse.Success(ctx, gin.H{"users": u, "total": total})
}

func (c *FollowController) Following(ctx *gin.Context) { c.list(ctx, false) }

func (c *FollowController) Followers(ctx *gin.Context) { c.list(ctx, true) }
