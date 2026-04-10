package main

import (
	"MainService/handlers"
	"MainService/middleware"
	"MainService/minio"
	"MainService/repositories"
	"MainService/services"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "MainService/docs"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	logrus "github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDB() *gorm.DB {
	gormUser := os.Getenv("GORM_USER")
	gormPassword := os.Getenv("GORM_PASSWORD")
	gormName := os.Getenv("GORM_NAME")
	gormHost := os.Getenv("GORM_HOST")
	gormPort := os.Getenv("GORM_PORT")

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", gormUser, gormPassword, gormHost, gormPort, gormName)

	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connection)
		if err != nil {
			log.Printf("Connection error for migrations:", err)
		} else {
			break
		}
		fmt.Println("Waiting for DB...")
		time.Sleep(3 * time.Second)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if errTwo := goose.SetDialect("postgres"); errTwo != nil {
		log.Fatal(errTwo)
	}
	log.Println("Launching migrations...")

	if errThree := goose.Up(db, "migrations"); errThree != nil {
		log.Fatal("Migration execution error:", errThree)
	}
	log.Println("Migrations successfully applied")

	gormDB, errFour := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if errFour != nil {
		log.Fatal("Error of GORM:", errFour)
	}

	return gormDB
}
func CloseDB(gormDB *gorm.DB) {
	s, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	errTwo := s.Close()
	if errTwo != nil {
		log.Fatal(errTwo)
	}
}

// @title LMSSysytem API
// @version 1.0
// @description This is a LMSSystem server.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8083

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	gormDB := InitDB()
	defer CloseDB(gormDB)

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("Logrus is configured")

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

	minios := minio.NewMinioClient(endpoint, accessKeyID, secretAccessKey)

	courseRepo := repositories.NewCourseRepository(gormDB)
	courseServ := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseServ)

	chapterRepo := repositories.NewChapterRepository(gormDB)
	chapterServ := services.NewChapterService(chapterRepo)
	chapterHandler := handlers.NewChapterHandler(chapterServ)

	lessonRepo := repositories.NewLessonRepository(gormDB)
	lessonServ := services.NewLessonService(lessonRepo)
	lessonHandler := handlers.NewLessonHandler(lessonServ)

	userLessonRepo := repositories.NewUserLessonRepo(gormDB)
	userLessonServ := services.NewUserLessonService(userLessonRepo)
	userLessonHandl := handlers.NewUserLessonHandler(userLessonServ)

	fileRepo := repositories.NewFileRepository(gormDB)
	fileServ := services.NewFileService(fileRepo, minios, userLessonServ)
	fileHandler := handlers.NewFleHandler(fileServ)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	{
		api.GET("/course", courseHandler.GetCourseH)
		api.GET("/course/:id", courseHandler.GetCourseByIDH)

		api.GET("/chapter", chapterHandler.GetChaptersH)
		api.GET("/chapter/course/:courseId", chapterHandler.GetChaptersByCourseIDH)
		api.GET("/chapter/:id", chapterHandler.GetChapterByIDH)

		api.GET("/lesson", lessonHandler.GetLessonsH)
		api.GET("/lesson/chapter/:chapterId", lessonHandler.GetLessonsByChapterIDH)
		api.GET("/lesson/:id", lessonHandler.GetLessonByIDH)

		api.POST("/upload", fileHandler.Upload)
		api.POST("/download/:lessonId", fileHandler.Download)
		api.GET("/get-file/:lessonId", fileHandler.GetFileByLessonID)

		api.POST("/add-user-lesson", userLessonHandl.AddUserLessonH)
		api.GET("/check-access/:lessonId", userLessonHandl.CheckUserLessonH)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AdminMiddleware())

	{
		admin.POST("/course", courseHandler.AddCourseH)
		admin.DELETE("/course/:id", courseHandler.DeleteCourseH)
		admin.PATCH("/course", courseHandler.UpdateCourseH)

		admin.POST("/chapter", chapterHandler.AddChapterH)
		admin.DELETE("/chapter/:id", chapterHandler.DeleteChapterH)
		admin.PATCH("/chapter", chapterHandler.UpdateChapterH)

		admin.POST("/lesson", lessonHandler.AddLessonH)
		admin.DELETE("/lesson/:id", lessonHandler.DeleteLessonH)
		admin.PATCH("/lesson", lessonHandler.UpdateLessonH)
	}

	router.Run(":" + os.Getenv("PORT"))
}
