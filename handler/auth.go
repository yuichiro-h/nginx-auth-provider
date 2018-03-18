package handler

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Auth struct {
	log *zap.Logger
}

func NewAuth(logger *zap.Logger) *Auth {
	return &Auth{
		log: logger,
	}
}

// Hostヘッダーに認証するアプリのHostが設定されている前提
func (h *Auth) Handle(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get(sessionKeyUser) == nil {
		c.Status(401)
		return
	}
	c.Status(200)
}
