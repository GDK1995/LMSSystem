package entities

import "time"

type Course struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chapter struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	CourseID    uint      `json:"course_id" gorm:"not null;index"`
	Course      Course    `json:"course" gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Lesson struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Content     string    `json:"content" gorm:"not null"`
	Order       int       `json:"order" gorm:"not null"`
	ChapterID   uint      `json:"chapter_id" gorm:"not null;index"`
	Chapter     Chapter   `json:"chapter" gorm:"foreignKey:ChapterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Attachment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	URL       string    `json:"url" gorm:"not null"`
	LessonID  uint      `json:"lesson_id" gorm:"not null;index"`
	Lesson    Lesson    `json:"lesson" gorm:"foreignKey:LessonID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt time.Time `json:"created_at"`
}

type UserLesson struct {
	UserID   string `json:"user_id" gorm:"primaryKey"`
	LessonID uint   `json:"lesson_id" gorm:"primaryKey"`
}
