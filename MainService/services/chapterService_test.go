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

func TestAddChapter(t *testing.T) {
	repo := new(mocks.ChapterRepository)
	service := NewChapterService(repo)

	chapter := entities.Chapter{
		Name:        "Context и как он работает",
		Description: "В этом разделе полностью описывается как работает context",
		Order:       4,
		CourseID:    2,
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("AddChapter", chapter).Return(uint(1), nil).Once()
		id, err := service.AddChapterS(chapter)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("AddChapter", mock.Anything).Return(uint(0), gorm.ErrInvalidData).Once()
		id, err := service.AddChapterS(chapter)
		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		assert.Equal(t, uint(0), id)
	})
}

func TestGetChapters(t *testing.T) {
	repo := new(mocks.ChapterRepository)
	service := NewChapterService(repo)

	t.Run("Success", func(t *testing.T) {
		chapters := []entities.Chapter{
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

		repo.On("GetChapters").Return(chapters, nil).Once()
		results, err := service.GetChaptersS()

		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Основы языка Go", results[0].Name)
		assert.Equal(t, "Полная информация о функциях, рекурсии и замыкании", results[1].Description)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetChapters").Return(nil, gorm.ErrInvalidDB).Once()
		results, err := service.GetChaptersS()

		assert.Nil(t, results)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestGetChaptersByCourseID(t *testing.T) {
	repo := new(mocks.ChapterRepository)
	service := NewChapterService(repo)

	t.Run("Success", func(t *testing.T) {
		coursesChapter := []entities.Chapter{
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

		repo.On("GetChaptersByCourseID", uint(1)).Return(coursesChapter, nil).Once()
		results, err := service.GetChaptersByCourseIDS(uint(1))

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, uint(1), results[0].CourseID)
		assert.Equal(t, "Функции", results[1].Name)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetChaptersByCourseID", mock.Anything).Return(nil, gorm.ErrInvalidDB).Once()
		results, err := service.GetChaptersByCourseIDS(uint(1))

		assert.Nil(t, results)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestGetChaptersByID(t *testing.T) {
	repo := new(mocks.ChapterRepository)
	service := NewChapterService(repo)

	t.Run("Success", func(t *testing.T) {
		chapter := entities.Chapter{
			ID:          3,
			Name:        "Указатели",
			Description: "Полная информация об указателях",
			Order:       1,
			CourseID:    2,
		}

		repo.On("GetChapterByID", uint(3)).Return(chapter, nil).Once()
		result, err := service.GetChapterByIDS(uint(3))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(3), result.ID)
		assert.Equal(t, "Полная информация об указателях", result.Description)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("GetChapterByID", uint(10)).Return(entities.Chapter{}, gorm.ErrRecordNotFound).Once()
		result, err := service.GetChapterByIDS(uint(10))

		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrChapterNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetChapterByID", mock.Anything).Return(entities.Chapter{}, gorm.ErrInvalidDB).Once()
		result, err := service.GetChapterByIDS(uint(3))
		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestDeleteChapter(t *testing.T) {
	repo := new(mocks.ChapterRepository)
	service := NewChapterService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.On("DeleteChapter", uint(3)).Return(nil).Once()
		err := service.DeleteChapterS(uint(3))

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("DeleteChapter", uint(10)).Return(gorm.ErrRecordNotFound).Once()
		err := service.DeleteChapterS(uint(10))

		assert.ErrorIs(t, err, errorsEntities.ErrChapterNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo.On("DeleteChapter", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.DeleteChapterS(uint(3))

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}

func TestUpdateChapter(t *testing.T) {
	updChapter := entitiesDTO.ChapterDTO{
		ID:    1,
		Name:  "Основы языка Golang",
		Order: 0,
	}

	existChapter := entities.Chapter{
		ID:          1,
		Name:        "Основы языка Go",
		Description: "Полная информация о переменных, константах и типах данных и тд",
		Order:       1,
		CourseID:    1,
	}

	t.Run("Success", func(t *testing.T) {
		repo := new(mocks.ChapterRepository)
		service := NewChapterService(repo)

		repo.On("GetChapterByID", updChapter.ID).Return(existChapter, nil).Once()

		repo.On("UpdateChapter", mock.MatchedBy(func(c entities.Chapter) bool {
			return c.ID == updChapter.ID &&
				c.Name == updChapter.Name &&
				c.Description == existChapter.Description &&
				c.CourseID == existChapter.CourseID &&
				c.Order == 1
		})).Return(nil).Once()
		err := service.UpdateChapterS(updChapter)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo := new(mocks.ChapterRepository)
		service := NewChapterService(repo)

		repo.On("GetChapterByID", updChapter.ID).Return(entities.Chapter{}, gorm.ErrRecordNotFound).Once()
		err := service.UpdateChapterS(updChapter)

		assert.ErrorIs(t, err, errorsEntities.ErrChapterNotFound)
		repo.AssertNotCalled(t, "UpdateChapter", mock.Anything)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo := new(mocks.ChapterRepository)
		service := NewChapterService(repo)

		repo.On("GetChapterByID", updChapter.ID).Return(existChapter, nil).Once()
		repo.On("UpdateChapter", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.UpdateChapterS(updChapter)

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}
