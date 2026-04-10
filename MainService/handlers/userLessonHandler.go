package handlers

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLessonHandler interface {
	AddUserLessonH(c *gin.Context)
	CheckUserLessonH(c *gin.Context)
}

type userLessonHandler struct {
	userLessonService services.UserLessonService
}

func NewUserLessonHandler(userLessonService services.UserLessonService) UserLessonHandler {
	return &userLessonHandler{userLessonService: userLessonService}
}

func (ulh *userLessonHandler) AddUserLessonH(c *gin.Context) {
	var userLesson entities.UserLesson

	if err := c.BindJSON(&userLesson); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := ulh.userLessonService.AddUserLessonS(userLesson)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

func (ulh *userLessonHandler) CheckUserLessonH(c *gin.Context) {
	var req entitiesDTO.CheckRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.Error(errorsEntities.ErrUnauthorized)
		return
	}

	item := entities.UserLesson{
		LessonID: req.LessonID,
		UserID:   userID,
	}

	isExit, err := ulh.userLessonService.CheckUserLessonS(item)
	if err != nil {
		c.Error(errorsEntities.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isExit": isExit,
	})
}
