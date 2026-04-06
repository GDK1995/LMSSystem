package main

import (
	"LMSSystem/entitiesDTO"
	"LMSSystem/keycloak"
	"LMSSystem/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	ctx := context.Background()

	kc, err := keycloak.NewKeycloakService(ctx)
	if err != nil {
		log.Fatalf("Keycloak service initialization failed: %v", err)
	}

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	authGroup := router.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware(os.Getenv("KEYCLOAK_URL"), os.Getenv("KEYCLOAK_REALM")))

	router.POST("/create-admin", func(c *gin.Context) {
		var signUpReq entitiesDTO.SignupReq

		if errThree := c.BindJSON(&signUpReq); errThree != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errThree.Error()})
			return
		}

		signUpReq.Role = "ROLE_ADMIN"

		userID, errTwo := kc.CreateUser(ctx, signUpReq)

		if errTwo != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
			return
		}

		fmt.Println("userID 2", userID)

		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	router.POST("/login", func(c *gin.Context) {
		var loginReq entitiesDTO.LoginReq

		if errTwo := c.BindJSON(&loginReq); errTwo != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errTwo.Error()})
			return
		}

		token, errThree := kc.LoginUser(ctx, loginReq)
		if errThree != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errThree.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	router.POST("/refresh", func(c *gin.Context) {
		var refreshToken string

		if errTwo := c.BindJSON(&refreshToken); errTwo != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errTwo.Error()})
			return
		}

		token, errThree := kc.RefreshToken(ctx, refreshToken)
		if errThree != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errThree.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	authGroup.POST("/register", func(c *gin.Context) {
		role, ok := c.Get("roles")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		roleSlice, ok := role.([]string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		hasAdmin := false
		for _, r := range roleSlice {
			if r == "ROLE_ADMIN" {
				hasAdmin = true
				break
			}
		}

		if hasAdmin {
			var signUpReq entitiesDTO.SignupReq

			if errThree := c.BindJSON(&signUpReq); errThree != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": errThree.Error()})
				return
			}

			userID, errTwo := kc.CreateUser(ctx, signUpReq)

			if errTwo != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Don't have permission to create user"})
			return
		}

	})

	authGroup.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		user, errTwo := kc.GetUser(ctx, id)
		if errTwo != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	authGroup.DELETE("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		roles, _ := c.Get("roles")
		fmt.Println("Запрос выполняет пользователь с ролями:", roles)

		errTwo := kc.DeleteUser(ctx, id)
		if errTwo != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	authGroup.PATCH("/user", func(c *gin.Context) {
		var user entitiesDTO.User
		if errTwo := c.BindJSON(&user); errTwo != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errTwo.Error()})
			return
		}

		errThree := kc.UpdateUser(ctx, user)
		if errThree != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errThree.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.Run(":8082")
}
