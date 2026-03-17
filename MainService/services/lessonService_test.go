package services

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddLesson(t *testing.T) {
	repo := new(mocks.LessonRepository)
	service := NewLessonService(repo)

	lesson := entities.Lesson{
		Name:        "Переменные",
		Description: "Урок полностью описывающий все типы переменных",
		Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
		Order:       1,
		ChapterID:   2,
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("AddLesson", lesson).Return(uint(1), nil).Once()
		id, err := service.AddLessonS(lesson)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("AddLesson", mock.Anything).Return(uint(0), gorm.ErrInvalidData).Once()
		id, err := service.AddLessonS(lesson)
		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		assert.Equal(t, uint(0), id)
	})
}

func TestGetLessons(t *testing.T) {
	repo := new(mocks.LessonRepository)
	service := NewLessonService(repo)

	t.Run("Success", func(t *testing.T) {
		lessons := []entities.Lesson{
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

		repo.On("GetLessons").Return(lessons, nil).Once()
		results, err := service.GetLessonsS()

		assert.NoError(t, err)
		assert.Len(t, results, 4)
		assert.Equal(t, "Каналы", results[2].Name)
		assert.Equal(t, 1, results[0].Order)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetLessons").Return(nil, gorm.ErrInvalidDB).Once()
		results, err := service.GetLessonsS()

		assert.Nil(t, results)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestGetLessonsByChapterID(t *testing.T) {
	repo := new(mocks.LessonRepository)
	service := NewLessonService(repo)

	t.Run("Success", func(t *testing.T) {
		chaptersLesson := []entities.Lesson{
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

		repo.On("GetLessonsByChapterID", uint(2)).Return(chaptersLesson, nil).Once()
		results, err := service.GetLessonsByChapterIDS(uint(2))

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, uint(2), results[0].ChapterID)
		assert.Equal(t, "Переменные bool", results[1].Name)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetLessonsByChapterID", mock.Anything).Return(nil, gorm.ErrInvalidDB).Once()
		results, err := service.GetLessonsByChapterIDS(uint(2))

		assert.Nil(t, results)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestGetLessonsByID(t *testing.T) {
	repo := new(mocks.LessonRepository)
	service := NewLessonService(repo)

	t.Run("Success", func(t *testing.T) {
		lesson := entities.Lesson{
			ID:          4,
			Name:        "Переменные bool",
			Description: "Про тип переменных правда ложь",
			Content:     "По значению отвечает правда или ложь",
			Order:       2,
			ChapterID:   2,
		}

		repo.On("GetLessonByID", uint(4)).Return(lesson, nil).Once()
		result, err := service.GetLessonByIDS(uint(4))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(4), result.ID)
		assert.Equal(t, "Про тип переменных правда ложь", result.Description)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("GetLessonByID", uint(404)).Return(entities.Lesson{}, gorm.ErrRecordNotFound).Once()
		result, err := service.GetLessonByIDS(uint(404))

		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrLessonNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetLessonByID", mock.Anything).Return(entities.Lesson{}, gorm.ErrInvalidDB).Once()
		result, err := service.GetLessonByIDS(uint(3))
		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestDeleteLesson(t *testing.T) {
	repo := new(mocks.LessonRepository)
	service := NewLessonService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.On("DeleteLesson", uint(3)).Return(nil).Once()
		err := service.DeleteLessonS(uint(3))

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("DeleteLesson", uint(99)).Return(gorm.ErrRecordNotFound).Once()
		err := service.DeleteLessonS(uint(99))

		assert.ErrorIs(t, err, errorsEntities.ErrLessonNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo.On("DeleteLesson", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.DeleteLessonS(uint(3))

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}

func TestUpdateLesson(t *testing.T) {
	updLesson := entitiesDTO.LessonDTO{
		ID:    1,
		Name:  "Переменные Golang",
		Order: 0,
	}

	existLesson := entities.Lesson{
		ID:          1,
		Name:        "Переменные",
		Description: "Урок полностью описывающий все типы переменных",
		Content:     "В Golang есть такие типы переменных как: int, float, string, bool и тд",
		Order:       1,
		ChapterID:   2,
	}

	t.Run("Success", func(t *testing.T) {
		repo := new(mocks.LessonRepository)
		service := NewLessonService(repo)

		repo.On("GetLessonByID", updLesson.ID).Return(existLesson, nil).Once()

		repo.On("UpdateLesson", mock.MatchedBy(func(c entities.Lesson) bool {
			return c.ID == updLesson.ID &&
				c.Name == updLesson.Name &&
				c.Description == existLesson.Description &&
				c.Content == existLesson.Content &&
				c.ChapterID == existLesson.ChapterID &&
				c.Order == 1
		})).Return(nil).Once()
		err := service.UpdateLessonS(updLesson)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo := new(mocks.LessonRepository)
		service := NewLessonService(repo)

		repo.On("GetLessonByID", updLesson.ID).Return(entities.Lesson{}, gorm.ErrRecordNotFound).Once()
		err := service.UpdateLessonS(updLesson)

		assert.ErrorIs(t, err, errorsEntities.ErrLessonNotFound)
		repo.AssertNotCalled(t, "UpdateLesson", mock.Anything)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo := new(mocks.LessonRepository)
		service := NewLessonService(repo)

		repo.On("GetLessonByID", updLesson.ID).Return(existLesson, nil).Once()
		repo.On("UpdateLesson", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.UpdateLessonS(updLesson)

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}
