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

func TestAddCourse(t *testing.T) {
	repo := new(mocks.CourseRepository)
	service := NewCourseService(repo)

	course := entities.Course{
		Name:        "Golang",
		Description: "В этом курсе узнаете все про язык програмирования",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("AddCourse", course).Return(uint(1), nil).Once()
		id, err := service.AddCourseS(course)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), id)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("AddCourse", mock.Anything).Return(uint(0), gorm.ErrInvalidData).Once()
		id, err := service.AddCourseS(course)
		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		assert.Equal(t, uint(0), id)
	})
}

func TestGetCurses(t *testing.T) {
	repo := new(mocks.CourseRepository)
	service := NewCourseService(repo)

	t.Run("Success", func(t *testing.T) {
		courses := []entities.Course{
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

		repo.On("GetCourses").Return(courses, nil).Once()
		results, err := service.GetCoursesS()

		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Python", results[2].Name)
		assert.Equal(t, "В этом курсе узнаете все про язык програмирования Java", results[1].Description)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetCourses").Return(nil, gorm.ErrInvalidDB).Once()
		results, err := service.GetCoursesS()

		assert.Nil(t, results)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestGetCoursesByID(t *testing.T) {
	repo := new(mocks.CourseRepository)
	service := NewCourseService(repo)

	t.Run("Success", func(t *testing.T) {
		course := entities.Course{
			ID:          2,
			Name:        "Java",
			Description: "В этом курсе узнаете все про язык програмирования Java",
		}

		repo.On("GetCourseByID", uint(2)).Return(course, nil).Once()
		result, err := service.GetCourseByIDS(uint(2))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(2), result.ID)
		assert.Equal(t, "Java", result.Name)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("GetCourseByID", uint(504)).Return(entities.Course{}, gorm.ErrRecordNotFound).Once()
		result, err := service.GetCourseByIDS(uint(504))

		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrCourseNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repo.On("GetCourseByID", mock.Anything).Return(entities.Course{}, gorm.ErrInvalidDB).Once()
		result, err := service.GetCourseByIDS(uint(1))
		assert.Nil(t, result)
		assert.Error(t, err, errorsEntities.ErrInternalServer)
	})
}

func TestDeleteCourse(t *testing.T) {
	repo := new(mocks.CourseRepository)
	service := NewCourseService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.On("DeleteCourse", uint(3)).Return(nil).Once()
		err := service.DeleteCourseS(uint(3))

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("DeleteCourse", uint(10)).Return(gorm.ErrRecordNotFound).Once()
		err := service.DeleteCourseS(uint(10))

		assert.ErrorIs(t, err, errorsEntities.ErrCourseNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo.On("DeleteCourse", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.DeleteCourseS(uint(3))

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}

func TestUpdateCourse(t *testing.T) {
	updCourse := entitiesDTO.CourseDTO{
		ID:   1,
		Name: "Язык программирования Golang",
	}

	existCourse := entities.Course{
		ID:          1,
		Name:        "Golang",
		Description: "В этом курсе узнаете все про язык програмирования Golang",
	}

	t.Run("Success", func(t *testing.T) {
		repo := new(mocks.CourseRepository)
		service := NewCourseService(repo)

		repo.On("GetCourseByID", updCourse.ID).Return(existCourse, nil).Once()

		repo.On("UpdateCurse", mock.MatchedBy(func(c entities.Course) bool {
			return c.ID == updCourse.ID &&
				c.Name == updCourse.Name &&
				c.Description == existCourse.Description
		})).Return(nil).Once()
		err := service.UpdateCurseS(updCourse)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		repo := new(mocks.CourseRepository)
		service := NewCourseService(repo)

		repo.On("GetCourseByID", updCourse.ID).Return(entities.Course{}, gorm.ErrRecordNotFound).Once()
		err := service.UpdateCurseS(updCourse)

		assert.ErrorIs(t, err, errorsEntities.ErrCourseNotFound)
		repo.AssertNotCalled(t, "UpdateCurse", mock.Anything)
	})

	t.Run("Internal Repository Error", func(t *testing.T) {
		repo := new(mocks.CourseRepository)
		service := NewCourseService(repo)

		repo.On("GetCourseByID", updCourse.ID).Return(existCourse, nil).Once()
		repo.On("UpdateCurse", mock.Anything).Return(gorm.ErrInvalidDB).Once()
		err := service.UpdateCurseS(updCourse)

		assert.ErrorIs(t, err, errorsEntities.ErrInternalServer)
		repo.AssertExpectations(t)
	})
}
