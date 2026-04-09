package minio

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio interface {
	UploadMinio()
	DownloadMinio()
}

type minios struct {
	minioClient *minio.Client
	bucketName  string
	ctx         context.Context
}

func NewMinioClient() Minio {
	minioClient, bucketName, ctx := InitMinio()
	return &minios{minioClient: minioClient, bucketName: bucketName, ctx: ctx}
}

func InitMinio() (*minio.Client, string, context.Context) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
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

func (m *minios) UploadMinio() {
	filePath := "example.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString("Hello from MinIO + Go!")
	if err != nil {
		fmt.Println(err)
	}

	uploadInfo, err := m.minioClient.FPutObject(m.ctx, m.bucketName, "example.txt", filePath, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Файл загружен: %s, размер: %d\n", uploadInfo.Key, uploadInfo.Size)

}

func (m *minios) DownloadMinio() {
	err := m.minioClient.FGetObject(m.ctx, m.bucketName, "example.txt", "downloaded_example.txt", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Файл скачан как downloaded_example.txt")
}
