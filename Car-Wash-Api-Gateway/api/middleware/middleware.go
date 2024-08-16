package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/exam-5/Car-Wash-Api-Gateway/api/token"
	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		t := ctx.GetHeader("Authorization")

		url := ctx.Request.URL.Path

		if strings.Contains(url, "swagger") {

			ctx.Next()

			return

		} else if _, err := token.ExtractClaim(t); err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{

				"error": err.Error(),
			})

			return

		}

		ctx.Next()

	}

}
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// claims, exists := c.Get("claims")
		// if !exists {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// 	return
		// }
		claims, _ := token.ExtractClaim(c.Request.Header.Get("Authorization"))
		sub := claims["role"].(string)
		fmt.Println("role:", sub)
		obj := c.Request.URL.Path
		act := c.Request.Method

		allowed, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during authorization"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}
