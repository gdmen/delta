package common

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
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
