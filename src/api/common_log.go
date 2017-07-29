package api

import (
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/gin-gonic/gin.v1"
	"runtime"
	"strings"
)

func getRequestId(c *gin.Context) string {
	rid, exists := c.Get(RequestIdKey)
	if !exists {
		glog.Errorf("Couldn't find RequestIdKey")
		return "unknown"
	}
	return rid.(string)
}

func getFuncName(depth int) string {
	function, _, _, _ := runtime.Caller(depth + 1)
	split := strings.Split(runtime.FuncForPC(function).Name(), ".")
	return split[len(split)-1]
}

func GetLogPrefix(c *gin.Context) string {
	rid := getRequestId(c)
	fcn := getFuncName(1)
	return fmt.Sprintf("[rid=%s | fcn=%s]", rid, fcn)
}
