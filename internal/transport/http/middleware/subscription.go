package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	polargo "github.com/polarsource/polar-go"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
)

func SubscriptionMiddleware(polar *config.PolarConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := polargo.New(
			polargo.WithServer("sandbox"),
			polargo.WithSecurity(polar.AccessToken),
		)
		res, err := s.Customers.GetStateExternal(c, c.MustGet("user_id").(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res.CustomerState == nil || len(res.CustomerState.ActiveSubscriptions) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No active subscription found. Please upgrade to a pro plan"})
			return
		}
		c.Next()
	}
}
