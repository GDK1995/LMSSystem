package handlers

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/middleware"
	"MainService/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileHandler interface {
	GetFileByLessonID(c *gin.Context)
	Upload(c *gin.Context)
	Download(c *gin.Context)
}

type fileHandler struct {
	fileService services.FileService
}

func NewFleHandler(fileService services.FileService) FileHandler {
	return &fileHandler{fileService: fileService}
}

func (f *fileHandler) GetFileByLessonID(c *gin.Context) {
	strID := c.Param("lessonId")

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	files, errTwo := f.fileService.GetFileByLessonIDS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(200, files)
}

func (f *fileHandler) Upload(c *gin.Context) {
	hasAccess, err := middleware.CheckAccess(c)
	if err != nil {
		c.Error(errorsEntities.ErrForbidden)
		return
	}

	if !hasAccess {
		c.Error(errorsEntities.ErrForbidden)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}
	defer file.Close()

	name := c.PostForm("name")
	lessonIDStr := c.PostForm("lesson_id")

	lessonID, err := strconv.Atoi(lessonIDStr)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	attached := entities.Attachment{
		Name:     name,
		LessonID: uint(lessonID),
	}

	fileID, err := f.fileService.SaveFileS(file, header, attached)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file id": fileID,
	})
}

func (f *fileHandler) Download(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.Error(errorsEntities.ErrUnauthorized)
		return
	}

	var req entitiesDTO.DownloadReq
	if err := c.BindJSON(&req); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	lessonID, err := strconv.ParseUint(c.Param("lessonId"), 10, 32)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	downloadReq := entitiesDTO.DownloadReqDTO{
		UserID:   userID,
		Name:     req.Name,
		LessonID: uint(lessonID),
	}

	errTwo := f.fileService.DownloadS(downloadReq, c.Writer)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}
}
