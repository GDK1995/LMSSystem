package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakClaims struct {
	jwt.RegisteredClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	PreferredUsername string `json:"preferred_username"`
}

func AuthMiddleware(keycloakURL, realm string) gin.HandlerFunc {
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakURL, realm)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			fmt.Println("JWKS refresh error:", err)
		},
	})

	if err != nil {
		panic("failed to get JWKS: " + err.Error())
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Authorization:", c.GetHeader("Authorization"))

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer {token}"})
			return
		}

		claims := &KeycloakClaims{}
		fmt.Println("claims", claims)

		token, errTwo := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)

		if errTwo != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		fmt.Println("role", claims.RealmAccess.Roles)
		c.Set("user_id", claims.Subject)
		c.Set("username", claims.PreferredUsername)
		c.Set("roles", claims.RealmAccess.Roles)

		c.Next()
	}
}
