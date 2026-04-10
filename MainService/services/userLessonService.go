package services

import (
	"MainService/entities"
	"MainService/repositories"
)

type UserLessonService interface {
	AddUserLessonS(item entities.UserLesson) error
	CheckUserLessonS(item entities.UserLesson) (bool, error)
}

type userLessonService struct {
	userLessonRepo repositories.UserLessonRepository
}

func NewUserLessonService(userLessonRepo repositories.UserLessonRepository) UserLessonService {
	return &userLessonService{userLessonRepo: userLessonRepo}
}

func (uls *userLessonService) AddUserLessonS(item entities.UserLesson) error {
	err := uls.userLessonRepo.AddUserLesson(item)
	if err != nil {
		return err
	}

	return nil
}

func (uls *userLessonService) CheckUserLessonS(item entities.UserLesson) (bool, error) {
	isExit, err := uls.userLessonRepo.CheckUserLesson(item)
	if err != nil {
		return false, err
	}

	return isExit, nil
}
