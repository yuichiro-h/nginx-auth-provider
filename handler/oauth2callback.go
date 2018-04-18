package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuichiro-h/nginx-auth-provider/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Oauth2Callback struct {
	oauth *oauth2.Config
	log   *zap.Logger
}

func NewOauth2Callback(oauth *oauth2.Config, logger *zap.Logger) *Oauth2Callback {
	return &Oauth2Callback{
		oauth: oauth,
		log:   logger,
	}
}

// Google側で認証成功後にコールバックされる前提
// HeaderのHost情報がGoogle側から呼び出されているため、認証対象のアプリになっていないので、
// 認証対象のアプリに一度リダイレクトして、再度nginx-auth-providerを呼び出させる
func (h *Oauth2Callback) Handle(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		h.log.Warn("oauth2 callback code is empty")
		c.Status(403)
		return
	}

	token, err := h.oauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		h.log.Warn("failed to exchange token", zap.Error(err))
		c.Status(403)
		return
	}

	idTokens := strings.Split(token.Extra("id_token").(string), ".")
	if len(idTokens) < 2 {
		c.Status(403)
		return
	}

	data, err := base64Decode(idTokens[1])
	if err != nil {
		h.log.Warn("failed to decode token", zap.Error(err))
		c.Status(403)
		return
	}

	var info map[string]interface{}
	if err := json.Unmarshal(data, &info); err != nil {
		h.log.Warn("failed to decode json", zap.Error(err))
		c.Status(403)
		return
	}

	email := info["email"].(string)
	if email == "" {
		h.log.Warn("email is empty")
		c.Status(403)
		return
	}

	state := c.Query("state")
	stateDataJSON, found := stateCache.Get(state)
	if !found {
		h.log.Warn("not found state data.")
		c.Status(403)
		return
	}
	var stateData StateData
	if err := json.Unmarshal([]byte(stateDataJSON), &stateData); err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}

	accept := false
	domain := strings.Split(email, "@")[1]
	if len(stateData.AcceptDomains) > 0 {
		for _, acceptDomain := range stateData.AcceptDomains {
			if acceptDomain == domain {
				accept = true
				break
			}
		}
	} else if config.Get().GoogleDomain == domain {
		accept = true
	}

	if !accept {
		h.log.Warn("unexpected email domain", zap.String("request_domain", domain))
		c.String(403, "Invalid domain.")
		return
	}

	stateData.Email = &email
	newStateDataJSON, err := json.Marshal(&stateData)
	if err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}
	stateCache.Set(state, string(newStateDataJSON))

	callbackURL, err := url.Parse(stateData.Callback)
	if err != nil {
		h.log.Warn(err.Error())
		c.Status(403)
		return
	}
	query := callbackURL.Query()
	query.Add("state", state)
	callbackURL.RawQuery = query.Encode()

	c.Redirect(302, callbackURL.String())
}

// base64Decode decodes the Base64url encoded string
// steel from code.google.com/p/goauth2/oauth/jwt
func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
