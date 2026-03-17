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

func TestAddChapterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.ChapterService)
	handler := NewChapterHandler(service)

	t.Run("Success", func(t *testing.T) {
		inpData := entities.Chapter{
			Name:        "Структура",
			Description: "Полная информация о структуре и о его типах",
			Order:       3,
			CourseID:    1,
		}
		marshal, _ := json.Marshal(inpData)

		service.On("AddChapterS", mock.Anything).Return(uint(1), nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/chapter", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddChapterH(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		fmt.Println(w.Body.String())
		assert.Equal(t, float64(1), response["chapter id"])

		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		brokenJson := []byte(`{"name": "Broken JSON",,,, }`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/chapter", bytes.NewBuffer(brokenJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddChapterH(c)

		assert.NotEmpty(t, c.Errors)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		inpData := entities.Chapter{
			Name:        "Структура",
			Description: "Полная информация о структуре и о его типах",
			Order:       3,
			CourseID:    1,
		}
		marshal, _ := json.Marshal(inpData)

		service.On("AddChapterS", mock.Anything).Return(uint(0), errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/chapter", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddChapterH(c)

		assert.Equal(t, errorsEntities.ErrInternalServer, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})
}

func TestGetChaptersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.ChapterService)
	handler := NewChapterHandler(service)

	t.Run("Success", func(t *testing.T) {
		chapters := []entitiesDTO.ChapterDTO{
			{
				ID:          1,
				Name:        "Основы языка Go",
				Description: "Полная информация о переменных, константах и типах данных и тд",
				Order:       1,
				CourseID:    1,
			},
			{
				ID:          2,
				Name:        "Функции",
				Description: "Полная информация о функциях, рекурсии и замыкании",
				Order:       2,
				CourseID:    1,
			},
			{
				ID:          3,
				Name:        "Указатели",
				Description: "Полная информация об указателях",
				Order:       1,
				CourseID:    2,
			},
		}

		service.On("GetChaptersS").Return(chapters, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter", nil)

		handler.GetChaptersH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entitiesDTO.ChapterDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 3)
		assert.Equal(t, "Основы языка Go", response[0].Name)

		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service.On("GetChaptersS").Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter", nil)

		handler.GetChaptersH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
	})
}

func TestGetChapterByCourseIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		chapters := []entitiesDTO.ChapterDTO{
			{
				ID:          1,
				Name:        "Основы языка Go",
				Description: "Полная информация о переменных, константах и типах данных и тд",
				Order:       1,
				CourseID:    1,
			},
			{
				ID:          2,
				Name:        "Функции",
				Description: "Полная информация о функциях, рекурсии и замыкании",
				Order:       2,
				CourseID:    1,
			},
		}
		service.On("GetChaptersByCourseIDS", uint(1)).Return(chapters, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "courseId", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/course/1", nil)

		handler.GetChaptersByCourseIDH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entitiesDTO.ChapterDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "Функции", response[1].Name)

		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "courseId", Value: "a"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/course/a", nil)

		handler.GetChaptersByCourseIDH(c)

		assert.NotEmpty(t, c.Errors)
		assert.Equal(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "GetChaptersByCourseIDS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("GetChaptersByCourseIDS", uint(1)).Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "courseId", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/course/1", nil)

		handler.GetChaptersByCourseIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
	})
}

func TestGetChapterByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Succss", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		chapter := entitiesDTO.ChapterDTO{
			ID:          2,
			Name:        "Функции",
			Description: "Полная информация о функциях, рекурсии и замыкании",
			Order:       2,
			CourseID:    1,
		}

		service.On("GetChapterByIDS", uint(2)).Return(&chapter, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/2", nil)

		handler.GetChapterByIDH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entitiesDTO.ChapterDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Функции", response.Name)

		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "a"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/a", nil)

		handler.GetChapterByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "GetChapterByIDS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("GetChapterByIDS", uint(2)).Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/2", nil)

		handler.GetChapterByIDH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})

	t.Run("Chapter Not Found", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("GetChapterByIDS", uint(66)).Return(nil, errorsEntities.ErrChapterNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "66"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/chapter/66", nil)

		handler.GetChapterByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrChapterNotFound)
		service.AssertExpectations(t)
	})
}

func TestDeleteChapterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("DeleteChapterS", uint(5)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "5"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/chapter/5", nil)

		handler.DeleteChapterH(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		service.AssertExpectations(t)
	})

	t.Run("Chapter Not Found", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("DeleteChapterS", uint(55)).Return(errorsEntities.ErrChapterNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "55"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/chapter/55", nil)

		handler.DeleteChapterH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrChapterNotFound)
		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "wrong"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/chapter/wrong", nil)

		handler.DeleteChapterH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "DeleteChapterS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		service.On("DeleteChapterS", uint(2)).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/chapter/2", nil)

		handler.DeleteChapterH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}

func TestUpdateChapterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		updDTO := entitiesDTO.ChapterDTO{ID: 1, Name: "Основы Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateChapterS", updDTO).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/chapter", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateChapterH(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		body := []byte(`{"id": 1, "name": 2`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/chapter", bytes.NewBuffer(body))

		handler.UpdateChapterH(c)

		assert.NotEmpty(t, c.Errors)
		assert.Equal(t, errorsEntities.ErrBadRequest, c.Errors.Last().Err)
	})

	t.Run("Chapter Not Found", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		updDTO := entitiesDTO.ChapterDTO{ID: 99, Name: "Основы Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateChapterS", updDTO).Return(errorsEntities.ErrChapterNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/chapter", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateChapterH(c)

		assert.Equal(t, errorsEntities.ErrChapterNotFound, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.ChapterService)
		handler := NewChapterHandler(service)

		updDTO := entitiesDTO.ChapterDTO{ID: 1, Name: "Основы Golang"}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateChapterS", updDTO).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("PATCH", "/api/v1/chapter", bytes.NewBuffer(body))

		handler.UpdateChapterH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}
