package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	issuer := os.Getenv("KEYCLOAK_URL")
	provider := waitForKeycloak(issuer)

	verifier := provider.Verifier(&oidc.Config{
		ClientID:          os.Getenv("KEYCLOAK_CLIENT_ID"),
		SkipClientIDCheck: true,
	})

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, errTwo := verifier.Verify(c.Request.Context(), tokenStr)
		if errTwo != nil {
			fmt.Printf("Verification failed: %v", errTwo)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var claims map[string]interface{}
		if errThree := idToken.Claims(&claims); errThree != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Next()
	}
}

func waitForKeycloak(issuer string) *oidc.Provider {
	var provider *oidc.Provider
	var err error
	maxAttempts := 36
	for i := 1; i <= maxAttempts; i++ {
		provider, err = oidc.NewProvider(context.Background(), issuer)
		if err == nil {
			fmt.Println("Connected to Keycloak!")
			return provider
		}
		fmt.Printf("Waiting for Keycloak (%d/%d): %v\n", i, maxAttempts, err)
		time.Sleep(5 * time.Second)
	}
	panic("Failed to connect to Keycloak after retries: " + err.Error())
}
