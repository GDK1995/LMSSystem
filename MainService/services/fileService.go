package services

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/minio"
	"MainService/repositories"
	"fmt"
	"mime/multipart"
	"net/http"
)

type FileService interface {
	SaveFileS(file multipart.File, header *multipart.FileHeader, attach entities.Attachment) (uint, error)
	DownloadS(downloadReq entitiesDTO.DownloadReqDTO, w http.ResponseWriter) error
	GetFileByLessonIDS(lessonID uint) ([]entities.Attachment, error)
}

type fileService struct {
	fileRepository    repositories.FileRepository
	minioClient       minio.Minio
	userLessonService UserLessonService
}

func NewFileService(fileRepository repositories.FileRepository, minioClient minio.Minio, userLessonService UserLessonService) FileService {
	return &fileService{fileRepository: fileRepository, minioClient: minioClient, userLessonService: userLessonService}
}

func (fs *fileService) GetFileByLessonIDS(lessonID uint) ([]entities.Attachment, error) {
	fileList, err := fs.fileRepository.GetFileByLessonID(lessonID)
	if err != nil {
		return nil, errorsEntities.ErrInternalServer
	}

	return fileList, nil
}

func (fs *fileService) SaveFileS(file multipart.File, header *multipart.FileHeader, attach entities.Attachment) (uint, error) {
	url, name, err := fs.minioClient.UploadMinio(file, header)
	if err != nil {
		return 0, errorsEntities.ErrInternalServer
	}

	attach.URL = url
	attach.Name = name

	fileID, err := fs.fileRepository.SaveFile(attach)
	if err != nil {
		return 0, errorsEntities.ErrInternalServer
	}

	return fileID, nil
}

func (fs *fileService) DownloadS(downloadReq entitiesDTO.DownloadReqDTO, w http.ResponseWriter) error {
	checkReq := entities.UserLesson{
		UserID:   downloadReq.UserID,
		LessonID: downloadReq.LessonID,
	}
	isExist, errThree := fs.userLessonService.CheckUserLessonS(checkReq)
	if errThree != nil {
		return errorsEntities.ErrInternalServer
	}

	if isExist {
		err := fs.minioClient.DownloadMinio(downloadReq.Name, w)
		if err != nil {
			fmt.Println("co", err)
			return errorsEntities.ErrInternalServer
		}
	} else {
		return errorsEntities.ErrForbidden
	}

	return nil
}
