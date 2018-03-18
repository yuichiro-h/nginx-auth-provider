package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Log struct {
	config LogConfig
}

type LogConfig struct {
	Logger      *zap.Logger
	IgnorePaths []string
}

func NewLog(config LogConfig) *Log {
	return &Log{
		config: config,
	}
}

func (m *Log) Handle(c *gin.Context) {
	for _, p := range m.config.IgnorePaths {
		if c.Request.URL.Path == p {
			return
		}
	}

	m.config.Logger.Info("request",
		zap.String("path", c.Request.URL.Path),
		zap.Any("Header", c.Request.Header))
}
