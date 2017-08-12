package common

import (
	"github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	RequestIdKey = "X-Request-Id"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := uuid.NewV4().String()
		c.Set(RequestIdKey, rid)
		c.Header(RequestIdKey, rid)
	}
}
