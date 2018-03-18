package handler

import (
	"encoding/json"

	"github.com/yuichiro-h/nginx-auth-provider/config"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Callback struct {
	log *zap.Logger
}

func NewCallback(logger *zap.Logger) *Callback {
	return &Callback{
		log: logger,
	}
}

// 認証完了後、認証情報を認証対象のアプリのCookieに設定するため呼び出される前提
func (h *Callback) Handle(c *gin.Context) {
	state := c.Query("state")
	stateDataJSON, found := stateCache.Get(state)
	if !found {
		h.log.Warn("not found state data")
		c.Status(403)
		return
	}
	var stateData StateData
	if err := json.Unmarshal([]byte(stateDataJSON), &stateData); err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}

	session := sessions.Default(c)
	session.Set(sessionKeyUser, *stateData.Email)

	secure := c.GetHeader("X-Forwarded-Proto") == "https"
	cookiePath := c.GetHeader(headerNameCallbackCookiePath)
	if cookiePath == "" {
		cookiePath = "/"
	}
	maxAge := config.Get().CookieMaxAge
	if maxAge == 0 {
		maxAge = 60 * 60 * 24 * 7 // 1 week
	}

	session.Options(sessions.Options{
		Path:   cookiePath,
		Secure: secure,
		MaxAge: maxAge,
	})

	if err := session.Save(); err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}

	c.Redirect(302, stateData.BackTo)
}
