package middleware

import (
	"context"

	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

// resolveRole 从 core-server 读取用户当前角色。
// 不信任 token 里的 role 声明，避免账号角色变更后旧 token 继续沿用旧权限。
func (m *AuthMiddleware) resolveRole(ctx context.Context, userID uint64) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, m.rpcClient.GetRequestTimeout())
	defer cancel()
	profile, err := m.rpcClient.GetUserClient().GetProfile(callCtx, &userpb.GetProfileRequest{UserId: userID})
	if err != nil {
		return "", err
	}
	return profile.GetRole(), nil
}

// Role 校验当前用户是否为管理员。
// 角色已在 RequireAuth 中按 core-server 的当前值解析并写入上下文。
func (m *AuthMiddleware) Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetRole(c) != "admin" {
			// 必须 Abort，否则后续 controller 还会再写一次 200 响应，
			// 浏览器会收到两个拼接的 JSON，前端解析直接失败。
			c.Abort()
			reponse.Forbidden(c)
			return
		}

		c.Next()
	}
}

