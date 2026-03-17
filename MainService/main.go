package main

import (
	"MainService/handlers"
	"MainService/middleware"
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

var gormDB *gorm.DB

func InitDB() {
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

	gormDB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatal("Error of GORM:", err)
	}
}
func CloseDB() {
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

	InitDB()
	defer CloseDB()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("Logrus is configured")

	courseRepo := repositories.NewCourseRepository(gormDB)
	courseServ := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseServ)

	chapterRepo := repositories.NewChapterRepository(gormDB)
	chapterServ := services.NewChapterService(chapterRepo)
	chapterHandler := handlers.NewChapterHandler(chapterServ)

	lessonRepo := repositories.NewLessonRepository(gormDB)
	lessonServ := services.NewLessonService(lessonRepo)
	lessonHandler := handlers.NewLessonHandler(lessonServ)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	{
		api.POST("/course", courseHandler.AddCourseH)
		api.GET("/course", courseHandler.GetCourseH)
		api.GET("/course/:id", courseHandler.GetCourseByIDH)
		api.DELETE("/course/:id", courseHandler.DeleteCourseH)
		api.PATCH("course", courseHandler.UpdateCourseH)

		api.POST("/chapter", chapterHandler.AddChapterH)
		api.GET("/chapter", chapterHandler.GetChaptersH)
		api.GET("/chapter/course/:courseId", chapterHandler.GetChaptersByCourseIDH)
		api.GET("/chapter/:id", chapterHandler.GetChapterByIDH)
		api.DELETE("/chapter/:id", chapterHandler.DeleteChapterH)
		api.PATCH("/chapter", chapterHandler.UpdateChapterH)

		api.POST("/lesson", lessonHandler.AddLessonH)
		api.GET("/lesson", lessonHandler.GetLessonsH)
		api.GET("/lesson/chapter/:chapterId", lessonHandler.GetLessonsByChapterIDH)
		api.GET("/lesson/:id", lessonHandler.GetLessonByIDH)
		api.DELETE("/lesson/:id", lessonHandler.DeleteLessonH)
		api.PATCH("/lesson", lessonHandler.UpdateLessonH)
	}

	router.Run(":8083")
}
