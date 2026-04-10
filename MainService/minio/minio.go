package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio interface {
	UploadMinio(file multipart.File, header *multipart.FileHeader) (string, string, error)
	DownloadMinio(objectName string, w http.ResponseWriter) error
}

type minios struct {
	minioClient *minio.Client
	bucketName  string
	ctx         context.Context
}

func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string) Minio {
	minioClient, bucketName, ctx := InitMinio(endpoint, accessKeyID, secretAccessKey)
	return &minios{minioClient: minioClient, bucketName: bucketName, ctx: ctx}
}

func InitMinio(endpoint string, accessKeyID string, secretAccessKey string) (*minio.Client, string, context.Context) {
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err)
	}

	bucketName := os.Getenv("MINIO_BUCKET")

	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
	}

	if !exists {
		log.Println("Bucket не существует, создаём...")
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
		}
	}

	return minioClient, bucketName, ctx
}

func (m *minios) UploadMinio(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	objectName := uuid.New().String() + "_" + header.Filename

	_, err := m.minioClient.PutObject(
		m.ctx,
		m.bucketName,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("http://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), m.bucketName, objectName)
	return url, objectName, nil
}

func (m *minios) DownloadMinio(objectName string, w http.ResponseWriter) error {
	object, err := m.minioClient.GetObject(
		m.ctx,
		m.bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		fmt.Println("err1", err)
		return err
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		fmt.Println("err2", err)
		return err
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+objectName)
	w.Header().Set("Content-Type", stat.ContentType)

	_, err = io.Copy(w, object)
	fmt.Println("err3", err)
	return err
}
