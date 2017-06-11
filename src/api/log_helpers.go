package api

import (
	"github.com/golang/glog"
	"gopkg.in/gin-gonic/gin.v1"
	"runtime"
	"strings"
)

func GetRequestId(c *gin.Context) string {
	rid, exists := c.Get(RequestIdKey)
	if !exists {
		glog.Errorf("Couldn't find RequestIdKey")
		return "unknown"
	}
	return rid.(string)
}

func GetFuncName() string {
	function, _, _, _ := runtime.Caller(1)
	split := strings.Split(runtime.FuncForPC(function).Name(), ".")
	return split[len(split)-1]
}
