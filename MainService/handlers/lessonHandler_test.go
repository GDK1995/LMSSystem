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

func TestAddLessonHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.LessonService)
	handler := NewLessonHandler(service)

	t.Run("Success", func(t *testing.T) {
		inpData := entities.Lesson{
			Name:        "Переменные",
			Description: "Урок полностью описывающий все типы переменных",
			Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
			Order:       1,
			ChapterID:   2,
		}
		marshal, _ := json.Marshal(inpData)

		service.On("AddLessonS", mock.Anything).Return(uint(1), nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/lesson", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddLessonH(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		fmt.Println(w.Body.String())
		assert.Equal(t, float64(1), response["lesson_id"])

		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		brokenJson := []byte(`{"name": "Broken JSON",,,, }`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/lesson", bytes.NewBuffer(brokenJson))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddLessonH(c)

		assert.NotEmpty(t, c.Errors)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		inpData := entities.Lesson{
			Name:        "Переменные",
			Description: "Урок полностью описывающий все типы переменных",
			Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
			Order:       1,
			ChapterID:   2,
		}
		marshal, _ := json.Marshal(inpData)

		service.On("AddLessonS", mock.Anything).Return(uint(0), errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/lesson", bytes.NewBuffer(marshal))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddLessonH(c)

		assert.Equal(t, errorsEntities.ErrInternalServer, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})
}

func TestGetLessonsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := new(mocks.LessonService)
	handler := NewLessonHandler(service)

	t.Run("Success", func(t *testing.T) {
		lessons := []entitiesDTO.LessonDTO{
			{
				ID:          1,
				Name:        "Переменные",
				Description: "Урок полностью описывающий все типы переменных",
				Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
				Order:       1,
				ChapterID:   2,
			},
			{
				ID:          2,
				Name:        "Фнкции",
				Description: "Фнкции и его особенности",
				Content:     "Функция помогает не повторять блок кода несколько раз",
				Order:       1,
				ChapterID:   1,
			},
			{
				ID:          3,
				Name:        "Каналы",
				Description: "Каналы и его виды",
				Content:     "Каналы делятся на буфферизированные и небуферизированные",
				Order:       1,
				ChapterID:   3,
			},
			{
				ID:          4,
				Name:        "Переменные bool",
				Description: "Про тип переменных правда ложь",
				Content:     "По значению отвечает правда или ложь",
				Order:       2,
				ChapterID:   2,
			},
		}

		service.On("GetLessonsS").Return(lessons, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson", nil)

		handler.GetLessonsH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entitiesDTO.LessonDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 4)
		assert.Equal(t, "Переменные", response[0].Name)

		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service.On("GetLessonsS").Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson", nil)

		handler.GetLessonsH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
	})
}

func TestGetLessonByChapterIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		lessons := []entitiesDTO.LessonDTO{
			{
				ID:          1,
				Name:        "Переменные",
				Description: "Урок полностью описывающий все типы переменных",
				Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
				Order:       1,
				ChapterID:   2,
			},
			{
				ID:          4,
				Name:        "Переменные bool",
				Description: "Про тип переменных правда ложь",
				Content:     "По значению отвечает правда или ложь",
				Order:       2,
				ChapterID:   2,
			},
		}
		service.On("GetLessonsByChapterIDS", uint(2)).Return(lessons, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "chapterId", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/chapter/2", nil)

		handler.GetLessonsByChapterIDH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entitiesDTO.LessonDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, uint(4), response[1].ID)

		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "chapterId", Value: "a"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/chapter/a", nil)

		handler.GetLessonsByChapterIDH(c)

		assert.NotEmpty(t, c.Errors)
		assert.Equal(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "GetLessonsByChapterIDS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("GetLessonsByChapterIDS", uint(1)).Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "chapterId", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/chapter/1", nil)

		handler.GetLessonsByChapterIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
	})
}

func TestGetLessonByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Succss", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		lesson := entitiesDTO.LessonDTO{
			ID:          4,
			Name:        "Переменные bool",
			Description: "Про тип переменных правда ложь",
			Content:     "По значению отвечает правда или ложь",
			Order:       2,
			ChapterID:   2,
		}

		service.On("GetLessonByIDS", uint(4)).Return(&lesson, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "4"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/4", nil)

		handler.GetLessonByIDH(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entitiesDTO.LessonDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Переменные bool", response.Name)

		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "a"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/a", nil)

		handler.GetLessonByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "GetLessonByIDS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("GetLessonByIDS", uint(2)).Return(nil, errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/2", nil)

		handler.GetLessonByIDH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})

	t.Run("Lesson Not Found", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("GetLessonByIDS", uint(66)).Return(nil, errorsEntities.ErrLessonNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "66"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/lesson/66", nil)

		handler.GetLessonByIDH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrLessonNotFound)
		service.AssertExpectations(t)
	})
}

func TestDeleteLessonHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("DeleteLessonS", uint(2)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/lesson/2", nil)

		handler.DeleteLessonH(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		service.AssertExpectations(t)
	})

	t.Run("Lesson Not Found", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("DeleteLessonS", uint(55)).Return(errorsEntities.ErrLessonNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "55"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/lesson/55", nil)

		handler.DeleteLessonH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrLessonNotFound)
		service.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "wrong"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/lesson/wrong", nil)

		handler.DeleteLessonH(c)

		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrBadRequest)
		service.AssertNotCalled(t, "DeleteLessonS", mock.Anything)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		service.On("DeleteLessonS", uint(2)).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/lesson/2", nil)

		handler.DeleteLessonH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}

func TestUpdateLessonHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		updDTO := entitiesDTO.LessonDTO{
			ID:    1,
			Name:  "Переменные Golang",
			Order: 0,
		}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateLessonS", updDTO).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/lesson", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLessonH(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
		service.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		body := []byte(`{"id": 1, "name": 2`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/lesson", bytes.NewBuffer(body))

		handler.UpdateLessonH(c)

		assert.NotEmpty(t, c.Errors)
		assert.Equal(t, errorsEntities.ErrBadRequest, c.Errors.Last().Err)
	})

	t.Run("Lesson Not Found", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		updDTO := entitiesDTO.LessonDTO{
			ID:    99,
			Name:  "Переменные Golang",
			Order: 0,
		}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateLessonS", updDTO).Return(errorsEntities.ErrLessonNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PATCH", "/api/v1/lesson", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLessonH(c)

		assert.Equal(t, errorsEntities.ErrLessonNotFound, c.Errors.Last().Err)
		service.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		service := new(mocks.LessonService)
		handler := NewLessonHandler(service)

		updDTO := entitiesDTO.LessonDTO{
			ID:   1,
			Name: "Переменные Golang",
		}
		body, _ := json.Marshal(updDTO)

		service.On("UpdateLessonS", updDTO).Return(errorsEntities.ErrInternalServer).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("PATCH", "/api/v1/lesson", bytes.NewBuffer(body))

		handler.UpdateLessonH(c)

		assert.NotEmpty(t, c.Errors)
		assert.ErrorIs(t, c.Errors.Last().Err, errorsEntities.ErrInternalServer)
		service.AssertExpectations(t)
	})
}
