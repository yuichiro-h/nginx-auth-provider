package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type StateData struct {
	BackTo   string
	Callback string
	Email    *string
}

type Initiate struct {
	oauth *oauth2.Config
	log   *zap.Logger
}

func NewInitiate(oauth *oauth2.Config, logger *zap.Logger) *Initiate {
	return &Initiate{
		oauth: oauth,
		log:   logger,
	}
}

// authリソースへのGETリクエストで401が返った場合に、Nginxから内部リクエストで呼び出される前提
// Googleの認証ページにリダイレクトさせる
func (h *Initiate) Handle(c *gin.Context) {
	backTo := c.Request.Header.Get(headerNameInitiateBackTo)
	callback := c.Request.Header.Get(headerNameInitiateCallback)
	if backTo == "" || callback == "" {
		h.log.Warn("backTo and callback URL is empty.")
		c.Status(403)
		return
	}

	stateData := StateData{
		BackTo:   backTo,
		Callback: callback,
	}
	stateDataJSON, err := json.Marshal(&stateData)
	if err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}

	stateKey := uuid.NewV4().String()
	stateCache.Set(stateKey, string(stateDataJSON))
	url := h.oauth.AuthCodeURL(stateKey)
	c.Redirect(302, url)
}
