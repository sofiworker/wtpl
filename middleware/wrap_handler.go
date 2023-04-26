package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WrapHandler func(c *gin.Context) (interface{}, error)

func WrapJsonHandler(handler WrapHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
