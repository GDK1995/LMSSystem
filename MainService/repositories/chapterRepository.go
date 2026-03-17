package repositories

import (
	"MainService/entities"

	"gorm.io/gorm"
)

type ChapterRepository interface {
	AddChapter(chapter entities.Chapter) (uint, error)
	GetChapters() ([]entities.Chapter, error)
	GetChaptersByCourseID(courseID uint) ([]entities.Chapter, error)
	GetChapterByID(chapterID uint) (entities.Chapter, error)
	DeleteChapter(chapterID uint) error
	UpdateChapter(updChapter entities.Chapter) error
}

type chapterRepository struct {
	gormDB *gorm.DB
}

func NewChapterRepository(gormDB *gorm.DB) ChapterRepository {
	return &chapterRepository{gormDB: gormDB}
}

func (chr *chapterRepository) AddChapter(chapter entities.Chapter) (uint, error) {
	err := chr.gormDB.Create(&chapter).Error
	if err != nil {
		return 0, err
	}

	return chapter.ID, nil
}

func (chr *chapterRepository) GetChapters() ([]entities.Chapter, error) {
	var chapters []entities.Chapter
	err := chr.gormDB.Preload("Course").Find(&chapters).Error
	if err != nil {
		return nil, err
	}

	return chapters, nil
}

func (chr *chapterRepository) GetChaptersByCourseID(courseID uint) ([]entities.Chapter, error) {
	var chapters []entities.Chapter
	err := chr.gormDB.Preload("Course").Where("course_id = ?", courseID).Find(&chapters).Error
	if err != nil {
		return nil, err
	}

	return chapters, nil
}

func (chr *chapterRepository) GetChapterByID(chapterID uint) (entities.Chapter, error) {
	var chapter entities.Chapter
	err := chr.gormDB.Preload("Course").First(&chapter, chapterID).Error
	if err != nil {
		return entities.Chapter{}, err
	}

	return chapter, nil
}

func (chr *chapterRepository) DeleteChapter(chapterID uint) error {
	result := chr.gormDB.Delete(entities.Chapter{}, chapterID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (chr *chapterRepository) UpdateChapter(updChapter entities.Chapter) error {
	result := chr.gormDB.Model(&entities.Chapter{}).Where("id = ?", updChapter.ID).Updates(updChapter)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
