package config

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var e = gin.Default()

func R() *gin.Engine {
	return e
}

// 成功
func Ok(c *gin.Context, cm CM, data interface{}) {
	baseHandler(c, http.StatusOK, cm, data)
}
func Ok2(c *gin.Context, data interface{}) {
	baseHandler(c, http.StatusOK, &SUCCESS, data)
}

// 失败
func Error(c *gin.Context, cm CM, data interface{}) {
	baseHandler(c, http.StatusInternalServerError, cm, data)
}
func Error2(c *gin.Context, data interface{}) {
	baseHandler(c, http.StatusInternalServerError, &SERVER_ERROR, data)
}

// 自定义
func Custom(c *gin.Context, status int, cm CM, data interface{}) {
	baseHandler(c, status, cm, data)
}

// 基础的私有模型
func baseHandler(c *gin.Context, status int, cm CM, data interface{}) {
	c.JSON(status, gin.H{
		"code":    cm.State(),
		"message": cm.Message(),
		"body":    data,
	})
}
