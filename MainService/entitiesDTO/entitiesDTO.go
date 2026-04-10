package entitiesDTO

import "MainService/entities"

type CourseDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChapterDTO struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Order       int             `json:"order"`
	CourseID    uint            `json:"course_id"`
	Course      entities.Course `json:"course"`
}

type LessonDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Content     string           `json:"content"`
	Order       int              `json:"order"`
	ChapterID   uint             `json:"chapter_id"`
	Chapter     entities.Chapter `json:"chapter"`
}

type CheckRequest struct {
	LessonID uint `json:"lesson_id"`
}

type DownloadReq struct {
	Name string `json:"name"`
}

type DownloadReqDTO struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	LessonID uint   `json:"lesson_id"`
}
