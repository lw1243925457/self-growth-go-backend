package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorResponse(c *gin.Context, code int16, message string) {
	c.JSON(int(code), gin.H{
		"code":    code,
		"message": message,
	})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

func SuccessResponseWithoutData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": "{}",
	})
}

func SuccessResponseWithPage(c *gin.Context, data interface{}, pageSize, pageIndex, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
		"pageSize": pageSize,
		"pageIndex": pageIndex,
		"total": total,
	})
}

func GetLoginUserName(c *gin.Context) string {
	return c.GetHeader("userName")
}
