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

type CourseHandler interface {
	AddCourseH(c *gin.Context)
	GetCourseH(c *gin.Context)
	GetCourseByIDH(c *gin.Context)
	DeleteCourseH(c *gin.Context)
	UpdateCourseH(c *gin.Context)
}

type courseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) CourseHandler {
	return &courseHandler{courseService: courseService}
}

// AddCourseH godoc
// @Summary Add new course
// @Description Add a new course to the database
// @Tags courses
// @Accept json
// @Produce json
// @Param course body entities.Course true "Course info"
// @Success 200 {object} map[string]uint "course_id"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/course [post]
func (ch *courseHandler) AddCourseH(c *gin.Context) {
	var course entities.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	courseID, errTwo := ch.courseService.AddCourseS(course)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"course_id": courseID,
	})
}

// GetCourseH godoc
// @Summary Get all courses
// @Description Retrieve all courses from the database
// @Tags courses
// @Produce json
// @Success 200 {array} entities.Course
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Failure 404 {object} middleware.AppError "Course not found"
// @Router /api/v1/course [get]
func (ch *courseHandler) GetCourseH(c *gin.Context) {
	courses, err := ch.courseService.GetCoursesS()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, courses)
}

// GetCourseByIDH godoc
// @Summary Get course by ID
// @Description Retrieve a single course by its ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} entities.Course
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Course not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/course/{id} [get]
func (ch *courseHandler) GetCourseByIDH(c *gin.Context) {
	strID := c.Param("id")

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	course, errTwo := ch.courseService.GetCourseByIDS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(200, course)
}

// DeleteCourseH godoc
// @Summary Delete course
// @Description Delete a course by its ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 404 {object} middleware.AppError "Course not found"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/course/{id} [delete]
func (ch *courseHandler) DeleteCourseH(c *gin.Context) {
	strID := c.Param("id")

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := ch.courseService.DeleteCourseS(uint(id))
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "no content",
	})
}

// UpdateCourseH godoc
// @Summary Update course
// @Description Update existing course info
// @Tags courses
// @Accept json
// @Produce json
// @Param course body entitiesDTO.CourseDTO true "Course info"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} middleware.AppError "Bad request"
// @Failure 500 {object} middleware.AppError "Internal server error"
// @Router /api/v1/course [patch]
func (ch *courseHandler) UpdateCourseH(c *gin.Context) {
	var courseDTO entitiesDTO.CourseDTO
	err := c.ShouldBindJSON(&courseDTO)
	if err != nil {
		c.Error(errorsEntities.ErrBadRequest)
		return
	}

	errTwo := ch.courseService.UpdateCurseS(courseDTO)
	if errTwo != nil {
		c.Error(errTwo)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
