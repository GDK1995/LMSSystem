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

type ChapterHandler interface {
	AddChapterH(c *gin.Context)
	GetChaptersH(c *gin.Context)
	GetChaptersByCourseIDH(c *gin.Context)
	GetChapterByIDH(c *gin.Context)
	DeleteChapterH(c *gin.Context)
	UpdateChapterH(c *gin.Context)
}

type chapterHandler struct {
	chapterService services.ChapterService
}

func NewChapterHandler(chapterService services.ChapterService) ChapterHandler {
	return &chapterHandler{chapterService: chapterService}
}

// AddChapterH godoc
// @Summary Add new chapter
// @Description Add a new chapter to the database
// @Tags chapters
// @Accept json
// @Produce json
// @Param chapter body entities.Chapter true "Chapter info"
// @Success 200 {object} map[string]uint "chapter_id"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/chapter [post]
func (chh *chapterHandler) AddChapterH(c *gin.Context) {
	var chapter entities.Chapter
	if err := c.BindJSON(&chapter); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	chapterID, errTwo := chh.chapterService.AddChapterS(chapter)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"chapter id": chapterID,
	})
}

// GetChaptersH godoc
// @Summary Get all chapters
// @Description Retrieve all chapters
// @Tags chapters
// @Produce json
// @Success 200 {array} entities.Chapter
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/chapter [get]
func (chh *chapterHandler) GetChaptersH(c *gin.Context) {
	chapters, err := chh.chapterService.GetChaptersS()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chapters)
}

// GetChaptersByCourseIDH godoc
// @Summary Get chapters by course ID
// @Description Retrieve all chapters for a specific course
// @Tags chapters
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {array} entities.Chapter
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/chapter/course/{courseId} [get]
func (chh *chapterHandler) GetChaptersByCourseIDH(c *gin.Context) {
	strID := c.Param("courseId")

	id64, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	chapters, errTwo := chh.chapterService.GetChaptersByCourseIDS(uint(id64))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, chapters)
}

// GetChapterByIDH godoc
// @Summary Get chapter by ID
// @Description Retrieve a single chapter by its ID
// @Tags chapters
// @Produce json
// @Param id path int true "Chapter ID"
// @Success 200 {object} entities.Chapter
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Router /api/v1/chapter/{id} [get]
func (chh *chapterHandler) GetChapterByIDH(c *gin.Context) {
	strID := c.Param("id")

	id64, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	chapter, errTwo := chh.chapterService.GetChapterByIDS(uint(id64))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// DeleteChapterH godoc
// @Summary Delete chapter
// @Description Delete a chapter by its ID
// @Tags chapters
// @Produce json
// @Param id path int true "Chapter ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/chapter/{id} [delete]
func (chh *chapterHandler) DeleteChapterH(c *gin.Context) {
	strID := c.Param("id")

	id64, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := chh.chapterService.DeleteChapterS(uint(id64))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "no content",
	})
}

// UpdateChapterH godoc
// @Summary Update chapter
// @Description Update existing chapter info
// @Tags chapters
// @Accept json
// @Produce json
// @Param chapter body entitiesDTO.ChapterDTO true "Chapter info"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/chapter [patch]
func (chh *chapterHandler) UpdateChapterH(c *gin.Context) {
	var chapterDTO entitiesDTO.ChapterDTO
	if err := c.BindJSON(&chapterDTO); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := chh.chapterService.UpdateChapterS(chapterDTO)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
