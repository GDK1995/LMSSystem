package services

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/mappers"
	"MainService/repositories"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LessonService interface {
	AddLessonS(lesson entities.Lesson) (uint, error)
	GetLessonsS() ([]entitiesDTO.LessonDTO, error)
	GetLessonsByChapterIDS(chapterID uint) ([]entitiesDTO.LessonDTO, error)
	GetLessonByIDS(lessonID uint) (*entitiesDTO.LessonDTO, error)
	DeleteLessonS(lessonID uint) error
	UpdateLessonS(updLesson entitiesDTO.LessonDTO) error
}

type lessonService struct {
	lessonRepository repositories.LessonRepository
}

func NewLessonService(lessonRepository repositories.LessonRepository) LessonService {
	return &lessonService{lessonRepository: lessonRepository}
}

func (ls *lessonService) AddLessonS(lesson entities.Lesson) (uint, error) {
	logrus.Info("Creating new lesson")
	logrus.Debugf("Lesson details: %+v", lesson)
	lessonID, err := ls.lessonRepository.AddLesson(lesson)
	if err != nil {
		logrus.Error("Failed to add lesson: ", err)
		return 0, errorsEntities.ErrInternalServer
	}

	logrus.Info("Lesson added successfully")
	return lessonID, nil
}

func (ls *lessonService) GetLessonsS() ([]entitiesDTO.LessonDTO, error) {
	logrus.Info("Getting lessons")
	lessons, err := ls.lessonRepository.GetLessons()
	if err != nil {
		logrus.Error("Failed to get lessons from repository: ", err)
		return nil, err
	}

	logrus.Debugf("Found %d lessons: %+v", len(lessons), lessons)
	lessonDTO := mappers.LessonsToDTO(lessons)

	logrus.Info("Lessons successfully converted to DTO")
	return lessonDTO, nil
}

func (ls *lessonService) GetLessonsByChapterIDS(chapterID uint) ([]entitiesDTO.LessonDTO, error) {
	logrus.Infof("Getting lessons by chapter id %d", chapterID)
	lessons, err := ls.lessonRepository.GetLessonsByChapterID(chapterID)
	if err != nil {
		logrus.Errorf("Failed to get lessons by chapter id %d: %v", chapterID, err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found %d lessons: %+v", len(lessons), lessons)
	lessonsDTO := mappers.LessonsToDTO(lessons)

	logrus.Info("Lessons successfully converted to DTO")
	return lessonsDTO, nil
}

func (ls *lessonService) GetLessonByIDS(lessonID uint) (*entitiesDTO.LessonDTO, error) {
	logrus.Infof("Getting lesson by id %d", lessonID)
	lesson, err := ls.lessonRepository.GetLessonByID(lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Lesson with id %d not found", lessonID)
			return nil, errorsEntities.ErrLessonNotFound
		}

		logrus.Errorf("Failed to get lesson by id %d: %v", lessonID, err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found lesson: %+v", lesson)
	lessonDTO := mappers.LessonToDTO(lesson)
	logrus.Info("Lesson successfully converted to DTO")

	return &lessonDTO, nil
}

func (ls *lessonService) DeleteLessonS(lessonID uint) error {
	logrus.Infof("Deleting lesson by id %d", lessonID)
	err := ls.lessonRepository.DeleteLesson(lessonID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Lesson with id %d not found", lessonID)
			return errorsEntities.ErrLessonNotFound
		}

		logrus.Error("Failed to delete lesson from repository: ", err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Info("Lesson successfully deleted")
	return nil
}

func (ls *lessonService) UpdateLessonS(updLesson entitiesDTO.LessonDTO) error {
	logrus.Infof("Updating lesson with id %d", updLesson.ID)
	lesson, err := ls.lessonRepository.GetLessonByID(updLesson.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Lesson with id %d not found", updLesson.ID)
			return errorsEntities.ErrLessonNotFound
		}

		logrus.Errorf("Failed to get lesson by id %d: %v", updLesson.ID, err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Current lesson data: %+v", lesson)

	if updLesson.Name != "" && updLesson.Name != lesson.Name {
		logrus.Debugf("Updating Name: %s to %s", lesson.Name, updLesson.Name)
		lesson.Name = updLesson.Name
	}

	if updLesson.Description != "" && updLesson.Description != lesson.Description {
		logrus.Debugf("Updating Description: %s to %s", lesson.Description, updLesson.Description)
		lesson.Description = updLesson.Description
	}

	if updLesson.Content != "" && updLesson.Content != lesson.Content {
		logrus.Debugf("Updating Content: %s to %s", lesson.Content, updLesson.Content)
		lesson.Content = updLesson.Content
	}

	if updLesson.Order != 0 && updLesson.Order != lesson.Order {
		logrus.Debugf("Updating Order: %d to %d", lesson.Order, updLesson.Order)
		lesson.Order = updLesson.Order
	}

	if updLesson.ChapterID != 0 && updLesson.ChapterID != lesson.ChapterID {
		logrus.Debugf("Updating ChapterID: %d to %d", lesson.ChapterID, updLesson.ChapterID)
		lesson.ChapterID = updLesson.ChapterID
	}

	errTwo := ls.lessonRepository.UpdateLesson(lesson)
	if errTwo != nil {
		if errors.Is(errTwo, gorm.ErrRecordNotFound) {
			logrus.Warnf("Lesson with id %d not found during update", updLesson.ID)
			return errorsEntities.ErrLessonNotFound
		}

		logrus.Errorf("Failed to update lesson id %d: %v", updLesson.ID, errTwo)
		return errorsEntities.ErrInternalServer
	}

	logrus.Infof("Leson with id %d successfully updated", updLesson.ID)
	return nil
}
