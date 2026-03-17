package handlers

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddCourseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.CourseService)
	handler := NewCourseHandler(service)

	t.Run("Success", func(t *testing.T) {
		inpData := entities.Course{
			Name:        "Golang",
			Description: "В этом курсе узнаете все про язык програмирования",
		}
		marshal, _ := json.Marshal(inpData)

		service.On("AddCourseS", mock.Anything).Return(uint(1), nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/course", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCourseH(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		fmt.Println(w.Body.String())
		assert.Equal(t, float64(1), response["course_id"])

		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		brokenJson := []byte(`{"name": "Broken JSON",,,, }`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/course", bytes.NewBuffer(brokenJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCourseH(c)

		assert.NotEmpty(t, c.Errors)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		inpData := entities.Course{
			Name:        "Golang",
			Description: "В этом курсе узнаете все про язык програмирования",
		}

		marshal, _ := json.Marshal(inpData)

		service.On("AddCourseS", mock.Anything).Return(uint(0), errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/course", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCourseH(c)

		assert.Equal(t, errorsEntities.ErrInternalServer, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})
}

func TestGetCoursesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.CourseService)
	handler := NewCourseHandler(service)

	t.Run("Success", func(t *testing.T) {
		courses := []entitiesDTO.CourseDTO{
			{
				ID:          1,
				Name:        "Golang",
				Description: "В этом курсе узнаете все про язык програмирования Golang",
			},
			{
				ID:          2,
				Name:        "Java",
				Description: "В этом курсе узнаете все про язык програмирования Java",
			},
			{
				ID:          3,
				Name:        "Python",
				Description: "В этом курсе узнаете все про язык програмирования Python",
			},
		}

		service.On("GetCoursesS").Return(courses, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/course", nil)

		handler.GetCourseH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entitiesDTO.CourseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 3)
		assert.Equal(t, "Python", response[2].Name)

		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service.On("GetCoursesS").Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/course", nil)

		handler.GetCourseH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
	})
}

func TestGetCourseByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		course := entitiesDTO.CourseDTO{
			ID:          2,
			Name:        "Java",
			Description: "В этом курсе узнаете все про язык програмирования Java",
		}

		service.On("GetCourseByIDS", uint(2)).Return(&course, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/course/2", nil)

		handler.GetCourseByIDH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entitiesDTO.CourseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Java", response.Name)

		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "a"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/course/a", nil)

		handler.GetCourseByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "GetCourseByIDS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		service.On("GetCourseByIDS", uint(2)).Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/course/2", nil)

		handler.GetCourseByIDH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})

	t.Run("Course Not Found", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		service.On("GetCourseByIDS", uint(66)).Return(nil, errorsEntities.ErrCourseNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "66"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/course/66", nil)

		handler.GetCourseByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrCourseNotFound)
		service.AssertExpectations(t)
	})
}

func TestDeleteCourseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		service.On("DeleteCourseS", uint(2)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/course/2", nil)

		handler.DeleteCourseH(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		service.AssertExpectations(t)
	})

	t.Run("Course Not Found", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		service.On("DeleteCourseS", uint(55)).Return(errorsEntities.ErrCourseNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "55"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/course/55", nil)

		handler.DeleteCourseH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrCourseNotFound)
		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "wrong"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/course/wrong", nil)

		handler.DeleteCourseH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "DeleteCourseS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		service.On("DeleteCourseS", uint(2)).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/course/2", nil)

		handler.DeleteCourseH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}

func TestUpdateCourseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		updDTO := entitiesDTO.CourseDTO{ID: 1, Name: "Язык программирования Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateCurseS", updDTO).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/course", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCourseH(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		body := []byte(`{"id": 1, "name": 2`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/course", bytes.NewBuffer(body))

		handler.UpdateCourseH(c)

		assert.NotEmpty(t, c.Errors)
		assert.Equal(t, errorsEntities.ErrBadRequest, c.Errors.Last().Err)
	})

	t.Run("Course Not Found", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		updDTO := entitiesDTO.CourseDTO{ID: 99, Name: "Язык программирования Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateCurseS", updDTO).Return(errorsEntities.ErrCourseNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/course", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCourseH(c)

		assert.Equal(t, errorsEntities.ErrCourseNotFound, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.CourseService)
		handler := NewCourseHandler(service)

		updDTO := entitiesDTO.CourseDTO{ID: 1, Name: "Язык программирования Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateCurseS", updDTO).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("PATCH", "/api/v1/course", bytes.NewBuffer(body))

		handler.UpdateCourseH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}
