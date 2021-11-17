package arc

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection rest by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path, zap.Any("error", err))
					c.Error(err.(error))
					c.Abort()
					return
				}

				status := http.StatusInternalServerError

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				log.Println()
				zap.L().Error("[ Recovery From Panic ]")
				for _, txt := range strings.Split(string(httpRequest), "\r\n") {
					zap.L().Error(txt)
				}
				if os.Getenv("GIN_MODE") != "release" {
					zap.L().Error("[ Debug of Stack ]")
					for _, txt := range strings.Split(string(debug.Stack()), "\r\n") {
						zap.L().Error(txt)
					}
				}
				e, ok := err.(Error)
				if !ok {
					c.AbortWithStatusJSON(status, gin.H{"error": err})
				} else {
					c.AbortWithStatusJSON(int(e.Type), e.JSON())
				}
			}
		}()
	}

}

func CrosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

	}
}
