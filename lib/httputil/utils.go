package httputil

import (
	"log"

	"github.com/gin-gonic/gin"
)

//MiddlewareHandler handle middleware error
func MiddlewareHandler(c *gin.Context) {
	c.Next()
	defer func(c *gin.Context) {
		log.Printf("something: %d - %d", len(c.Errors), c.Writer.Status())
		if len(c.Errors) > 0 {
			c.JSON(
				c.Writer.Status(),
				c.Errors,
			)
		}
	}(c)
}

//ResponseFailure sets response code and error to the given one in parameter.
func ResponseFailure(c *gin.Context, code int, err error) {
	c.JSON(
		code,
		gin.H{"error": err.Error()},
	)
}
