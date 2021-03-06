package main

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yuichiro-h/nginx-auth-provider/config"
	"github.com/yuichiro-h/nginx-auth-provider/handler"
	"github.com/yuichiro-h/nginx-auth-provider/middleware"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	if config.Get().Debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}

	if !config.Get().Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	cookieStore := sessions.NewCookieStore([]byte(config.Get().CookieSecret))
	r.Use(sessions.Sessions("session", cookieStore))

	logMiddleware := middleware.NewLog(middleware.LogConfig{
		IgnorePaths: []string{"/"},
		Logger:      logger,
	})
	r.Use(logMiddleware.Handle)

	oauthConfig := &oauth2.Config{
		ClientID:     config.Get().GoogleClientID,
		ClientSecret: config.Get().GoogleClientSecret,
		RedirectURL:  fmt.Sprintf("%s/oauth2callback", config.Get().GoogleCallbackURL),
		Scopes:       []string{"email"},
		Endpoint:     google.Endpoint,
	}

	authHandler := handler.NewAuth(logger).Handle
	initiateHandler := handler.NewInitiate(oauthConfig, logger).Handle
	oauth2CallbackHandler := handler.NewOauth2Callback(oauthConfig, logger).Handle
	callbackHandler := handler.NewCallback(logger).Handle

	r.GET("/", func(c *gin.Context) { c.Status(200) })
	r.GET("/auth", authHandler)
	r.GET("/initiate", initiateHandler)
	r.GET("/oauth2callback", oauth2CallbackHandler)
	r.GET("/callback", callbackHandler)

	r.Run(fmt.Sprintf(":%d", config.Get().Port))
}
