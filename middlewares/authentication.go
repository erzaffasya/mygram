package middlewares

import (
	"net/http"

	"github.com/erzaffasya/mygram/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helpers.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			return
		} else {
			c.Set("userData", verifyToken)
			c.Next()
		}
	}
}
