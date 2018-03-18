package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Log struct {
	logger *zap.Logger
}

func NewLog(logger *zap.Logger) *Log {
	return &Log{
		logger: logger,
	}
}

func (m *Log) Handle(c *gin.Context) {
	m.logger.Info("", zap.Any("Header", c.Request.Header))
}
