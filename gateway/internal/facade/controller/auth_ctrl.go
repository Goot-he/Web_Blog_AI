package controller

import (
	"net/http"
	"strings"

	"gateway/internal/application/user-module"
	"gateway/internal/config"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	cfg *config.Config
	svc *user_module.AuthService
}

func NewAuthController(svc *user_module.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		cfg: cfg,
		svc: svc,
	}
}

func (a *AuthController) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	if err := a.svc.SendEmailCode(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	if err := a.svc.Register(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := a.svc.Login(c.Request.Context(), req)
	if err != nil {
		if isUserBlockedError(err) {
			reponse.Fail(c, http.StatusForbidden, "该账号已被管理员拉黑，无法登录。")
			return
		}
		reponse.InternalServerError(c, err.Error())
		return
	}
	// 在cookie设置refresh_token
	utils.SetRefreshCookie(c.Writer, result.RefreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) EmailLogin(c *gin.Context) {
	var req dto.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)

		return
	}
	result, err := a.svc.EmailLogin(c.Request.Context(), req)
	if err != nil {
		if isUserBlockedError(err) {
			reponse.Fail(c, http.StatusForbidden, "该账号已被管理员拉黑，无法登录。")
			return
		}
		reponse.InternalServerError(c, err.Error())
		return
	}
	utils.SetRefreshCookie(c.Writer, result.RefreshToken, a.cfg.Auth)
	reponse.Success(c, result)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)

		return
	}
	if err := a.svc.ForgotPassword(c.Request.Context(), req); err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (a *AuthController) Logout(c *gin.Context) {
	refreshToken, _ := utils.RefreshTokenFromCookie(c.Request)
	utils.ClearRefreshCookie(c.Writer, a.cfg.Auth)
	if err := a.svc.Logout(c.Request.Context(), c.GetHeader("Authorization"), refreshToken); err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

// isUserBlockedError 判断 core-server 返回的是否为“用户被拉黑”错误。
// core-server 以普通 error 返回 USER_BLOCKED，gRPC 会包装成
// "rpc error: code = Unknown desc = USER_BLOCKED"，所以这里用包含匹配。
func isUserBlockedError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "USER_BLOCKED")
}
