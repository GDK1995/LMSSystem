package repositories

import (
	"MainService/entities"

	"gorm.io/gorm"
)

type LessonRepository interface {
	AddLesson(lesson entities.Lesson) (uint, error)
	GetLessons() ([]entities.Lesson, error)
	GetLessonsByChapterID(chapterID uint) ([]entities.Lesson, error)
	GetLessonByID(lessonID uint) (entities.Lesson, error)
	DeleteLesson(lessonID uint) error
	UpdateLesson(updLesson entities.Lesson) error
}

type lessonRepository struct {
	gormDB *gorm.DB
}

func NewLessonRepository(gormDB *gorm.DB) LessonRepository {
	return &lessonRepository{gormDB: gormDB}
}

func (lr *lessonRepository) AddLesson(lesson entities.Lesson) (uint, error) {
	if err := lr.gormDB.Create(&lesson).Error; err != nil {
		return 0, err
	}

	return lesson.ID, nil
}

func (lr *lessonRepository) GetLessons() ([]entities.Lesson, error) {
	var lessons []entities.Lesson
	if err := lr.gormDB.Preload("Chapter").Find(&lessons).Error; err != nil {
		return nil, err
	}

	return lessons, nil
}

func (lr *lessonRepository) GetLessonsByChapterID(chapterID uint) ([]entities.Lesson, error) {
	var lessons []entities.Lesson
	if err := lr.gormDB.Preload("Chapter").Where("chapter_id = ?", chapterID).Find(&lessons).Error; err != nil {
		return nil, err
	}

	return lessons, nil
}

func (lr *lessonRepository) GetLessonByID(lessonID uint) (entities.Lesson, error) {
	var lesson entities.Lesson

	if err := lr.gormDB.Preload("Chapter").First(&lesson, lessonID).Error; err != nil {
		return entities.Lesson{}, err
	}

	return lesson, nil
}

func (lr *lessonRepository) DeleteLesson(lessonID uint) error {
	result := lr.gormDB.Delete(&entities.Lesson{}, lessonID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (lr *lessonRepository) UpdateLesson(updLesson entities.Lesson) error {
	result := lr.gormDB.Model(&entities.Lesson{}).Where("id == ?", updLesson.ID).Updates(updLesson)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
