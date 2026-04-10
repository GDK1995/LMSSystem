package repositories

import (
	"MainService/entities"
	"MainService/errorsEntities"

	"errors"
	"gorm.io/gorm"
)

type UserLessonRepository interface {
	AddUserLesson(item entities.UserLesson) error
	CheckUserLesson(item entities.UserLesson) (bool, error)
}

type userLessonRepository struct {
	gormDB *gorm.DB
}

func NewUserLessonRepo(gormDB *gorm.DB) UserLessonRepository {
	return &userLessonRepository{gormDB: gormDB}
}

func (ul *userLessonRepository) AddUserLesson(item entities.UserLesson) error {
	err := ul.gormDB.Create(&item).Error
	if err != nil {
		return errorsEntities.ErrInternalServer
	}

	return nil
}

func (ul *userLessonRepository) CheckUserLesson(item entities.UserLesson) (bool, error) {
	var userLesson entities.UserLesson
	err := ul.gormDB.Where("user_id = ? AND lesson_id = ?", item.UserID, item.LessonID).First(&userLesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errorsEntities.ErrInternalServer
	}

	return true, nil
}
