package main

import (
	"LMSSystem/entitiesDTO"
	"LMSSystem/keycloak"
	"context"
	"log"
	"net/http"

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
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/register", func(c *gin.Context) {
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

	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		user, errTwo := kc.GetUser(ctx, id)
		if errTwo != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.DELETE("/user/:id", func(c *gin.Context) {
		id := c.Param("id")

		errTwo := kc.DeleteUser(ctx, id)
		if errTwo != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTwo.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.PATCH("/user", func(c *gin.Context) {
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
