package repositories

import (
	"MainService/entities"

	"gorm.io/gorm"
)

type FileRepository interface {
	SaveFile(attach entities.Attachment) (uint, error)
	GetFileByLessonID(lessonID uint) ([]entities.Attachment, error)
}

type fileRepository struct {
	gormDB *gorm.DB
}

func NewFileRepository(gormDB *gorm.DB) FileRepository {
	return &fileRepository{gormDB: gormDB}
}

func (fr *fileRepository) SaveFile(attach entities.Attachment) (uint, error) {
	err := fr.gormDB.Create(&attach).Error
	if err != nil {
		return 0, err
	}

	return attach.ID, nil
}

func (fr *fileRepository) GetFileByLessonID(lessonID uint) ([]entities.Attachment, error) {
	var attach []entities.Attachment
	err := fr.gormDB.Where("lesson_id = ?", lessonID).Find(&attach).Error
	if err != nil {
		return []entities.Attachment{}, err
	}

	return attach, nil
}
