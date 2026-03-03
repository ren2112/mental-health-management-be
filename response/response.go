package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func CommonResp(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, response{code, msg, data})
}
