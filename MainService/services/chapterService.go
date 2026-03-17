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

type ChapterService interface {
	AddChapterS(chapter entities.Chapter) (uint, error)
	GetChaptersS() ([]entitiesDTO.ChapterDTO, error)
	GetChaptersByCourseIDS(courseID uint) ([]entitiesDTO.ChapterDTO, error)
	GetChapterByIDS(chapterID uint) (*entitiesDTO.ChapterDTO, error)
	DeleteChapterS(chapterID uint) error
	UpdateChapterS(updChapter entitiesDTO.ChapterDTO) error
}

type chapterService struct {
	chapterRepository repositories.ChapterRepository
}

func NewChapterService(chapterRepository repositories.ChapterRepository) ChapterService {
	return &chapterService{chapterRepository: chapterRepository}
}

func (chs *chapterService) AddChapterS(chapter entities.Chapter) (uint, error) {
	logrus.Info("Creating new chapter")
	logrus.Debugf("Chapter details: %+v", chapter)
	chapterID, err := chs.chapterRepository.AddChapter(chapter)
	if err != nil {
		logrus.Error("Failed to add chapter: ", err)
		return 0, errorsEntities.ErrInternalServer
	}

	logrus.Info("Chapter added successfully")
	return chapterID, nil
}

func (chs *chapterService) GetChaptersS() ([]entitiesDTO.ChapterDTO, error) {
	logrus.Info("Getting chapters")
	chapters, err := chs.chapterRepository.GetChapters()
	if err != nil {
		logrus.Error("Failed to get chapters from repository: ", err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found %d chapters: %+v", len(chapters), chapters)
	chaptersDTO := mappers.ChaptersToDTO(chapters)

	logrus.Info("Chapters successfully converted to DTO")
	return chaptersDTO, nil
}

func (chs *chapterService) GetChaptersByCourseIDS(courseID uint) ([]entitiesDTO.ChapterDTO, error) {
	logrus.Infof("Getting chapters by course id %d", courseID)
	chapters, err := chs.chapterRepository.GetChaptersByCourseID(courseID)
	if err != nil {
		logrus.Errorf("Failed to get chapters by course id %d: %v", courseID, err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found %d chapters: %+v", len(chapters), chapters)
	chaptersDTO := mappers.ChaptersToDTO(chapters)

	logrus.Info("Chapters successfully converted to DTO")
	return chaptersDTO, nil
}

func (chs *chapterService) GetChapterByIDS(chapterID uint) (*entitiesDTO.ChapterDTO, error) {
	logrus.Infof("Getting chapter by id %d", chapterID)
	chapter, err := chs.chapterRepository.GetChapterByID(chapterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Chapter with id %d not found", chapterID)
			return nil, errorsEntities.ErrChapterNotFound
		}

		logrus.Errorf("Failed to get chapter by id %d: %v", chapterID, err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found chapter: %+v", chapter)
	chapterDTO := mappers.ChapterToDTO(chapter)
	logrus.Info("Chapter successfully converted to DTO")

	return &chapterDTO, nil
}

func (chs *chapterService) DeleteChapterS(chapterID uint) error {
	logrus.Infof("Deleting chapter by id %d", chapterID)
	err := chs.chapterRepository.DeleteChapter(chapterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Chapter with id %d not found", chapterID)
			return errorsEntities.ErrChapterNotFound
		}

		logrus.Error("Failed to delete chapter from repository: ", err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Info("Chapter successfully deleted")
	return nil
}

func (chs *chapterService) UpdateChapterS(updChapter entitiesDTO.ChapterDTO) error {
	logrus.Infof("Updating chapter with id %d", updChapter.ID)
	chapter, err := chs.chapterRepository.GetChapterByID(updChapter.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Chapter with id %d not found", updChapter.ID)
			return errorsEntities.ErrChapterNotFound
		}

		logrus.Errorf("Failed to get chapter by id %d: %v", updChapter.ID, err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Current chapter data: %+v", chapter)

	if updChapter.Name != "" && updChapter.Name != chapter.Name {
		logrus.Debugf("Updating Name: %s to %s", chapter.Name, updChapter.Name)
		chapter.Name = updChapter.Name
	}

	if updChapter.Description != "" && updChapter.Description != chapter.Description {
		logrus.Debugf("Updating Description: %s to %s", chapter.Description, updChapter.Description)
		chapter.Description = updChapter.Description
	}

	if updChapter.Order != 0 && updChapter.Order != chapter.Order {
		logrus.Debugf("Updating Order: %d to %d", chapter.Order, updChapter.Order)
		chapter.Order = updChapter.Order
	}

	if updChapter.CourseID != 0 && updChapter.CourseID != chapter.CourseID {
		logrus.Debugf("Updating CourseID: %d to %d", chapter.CourseID, updChapter.CourseID)
		chapter.CourseID = updChapter.CourseID
	}

	errTwo := chs.chapterRepository.UpdateChapter(chapter)
	if errTwo != nil {
		if errors.Is(errTwo, gorm.ErrRecordNotFound) {
			logrus.Warnf("Chapter with id %d not found during update", updChapter.ID)
			return errorsEntities.ErrChapterNotFound
		}

		logrus.Errorf("Failed to update chapter id %d: %v", updChapter.ID, errTwo)
		return errorsEntities.ErrInternalServer
	}

	logrus.Infof("Chapter with id %d successfully updated", updChapter.ID)
	return nil
}
