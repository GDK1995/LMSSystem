package handlers

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LessonHandler interface {
	AddLessonH(c *gin.Context)
	GetLessonsH(c *gin.Context)
	GetLessonsByChapterIDH(c *gin.Context)
	GetLessonByIDH(c *gin.Context)
	DeleteLessonH(c *gin.Context)
	UpdateLessonH(c *gin.Context)
}

type lessonHandler struct {
	lessonService services.LessonService
}

func NewLessonHandler(lessonService services.LessonService) LessonHandler {
	return &lessonHandler{lessonService: lessonService}
}

// AddLessonH godoc
// @Summary Add new lesson
// @Description Add a new lesson to the database
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body entities.Lesson true "Lesson info"
// @Success 200 {object} map[string]uint "lesson_id"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/lesson [post]
func (lh *lessonHandler) AddLessonH(c *gin.Context) {
	var lesson entities.Lesson
	if err := c.BindJSON(&lesson); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	lessonID, errTwo := lh.lessonService.AddLessonS(lesson)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"lesson_id": lessonID,
	})
}

// GetLessonsH godoc
// @Summary Get all lessons
// @Description Retrieve all lessons
// @Tags lessons
// @Produce json
// @Success 200 {array} entities.Lesson
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/lesson [get]
func (lh *lessonHandler) GetLessonsH(c *gin.Context) {
	lessons, err := lh.lessonService.GetLessonsS()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, lessons)
}

// GetLessonsByChapterIDH godoc
// @Summary Get lessons by chapter ID
// @Description Retrieve all lessons for a specific chapter
// @Tags lessons
// @Produce json
// @Param chapterID path int true "Chapter ID"
// @Success 200 {array} entities.Lesson
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/lesson/chapter/{chapterID} [get]
func (lh *lessonHandler) GetLessonsByChapterIDH(c *gin.Context) {
	strID := c.Param("chapterId")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	lessons, errTwo := lh.lessonService.GetLessonsByChapterIDS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, lessons)
}

// GetLessonByIDH godoc
// @Summary Get lesson by ID
// @Description Retrieve a single lesson by its ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} entities.Lesson
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Router /api/v1/lesson/{id} [get]
func (lh *lessonHandler) GetLessonByIDH(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	lesson, errTwo := lh.lessonService.GetLessonByIDS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, lesson)
}

// DeleteLessonH godoc
// @Summary Delete lesson
// @Description Delete a lesson by its ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/lesson/{id} [delete]
func (lh *lessonHandler) DeleteLessonH(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := lh.lessonService.DeleteLessonS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "success",
	})
}

// UpdateLessonH godoc
// @Summary Update lesson
// @Description Update existing lesson info
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body entitiesDTO.LessonDTO true "Lesson info"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/lesson [patch]
func (lh *lessonHandler) UpdateLessonH(c *gin.Context) {
	var lessonDTO entitiesDTO.LessonDTO
	if err := c.BindJSON(&lessonDTO); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := lh.lessonService.UpdateLessonS(lessonDTO)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
