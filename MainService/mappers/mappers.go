package mappers

import (
	"MainService/entities"
	"MainService/entitiesDTO"
)

func CoursesToDTO(courses []entities.Course) []entitiesDTO.CourseDTO {
	if len(courses) == 0 {
		return []entitiesDTO.CourseDTO{}
	}

	var coursesDTO []entitiesDTO.CourseDTO
	for i := 0; i < len(courses); i++ {
		var courseDTO entitiesDTO.CourseDTO
		courseDTO.ID = courses[i].ID
		courseDTO.Name = courses[i].Name
		courseDTO.Description = courses[i].Description
		coursesDTO = append(coursesDTO, courseDTO)
	}

	return coursesDTO
}

func CourseToDTO(course entities.Course) entitiesDTO.CourseDTO {
	var courseDTO entitiesDTO.CourseDTO
	courseDTO.ID = course.ID
	courseDTO.Name = course.Name
	courseDTO.Description = course.Description

	return courseDTO
}

func ChaptersToDTO(chapters []entities.Chapter) []entitiesDTO.ChapterDTO {
	if len(chapters) == 0 {
		return []entitiesDTO.ChapterDTO{}
	}

	var chaptersDTO []entitiesDTO.ChapterDTO
	for i := 0; i < len(chapters); i++ {
		var chapterDTO entitiesDTO.ChapterDTO
		chapterDTO.ID = chapters[i].ID
		chapterDTO.Name = chapters[i].Name
		chapterDTO.Description = chapters[i].Description
		chapterDTO.Order = chapters[i].Order
		chapterDTO.CourseID = chapters[i].CourseID
		chapterDTO.Course = chapters[i].Course
		chaptersDTO = append(chaptersDTO, chapterDTO)
	}

	return chaptersDTO
}

func ChapterToDTO(chapter entities.Chapter) entitiesDTO.ChapterDTO {
	var chapterDTO entitiesDTO.ChapterDTO
	chapterDTO.ID = chapter.ID
	chapterDTO.Name = chapter.Name
	chapterDTO.Description = chapter.Description
	chapterDTO.Order = chapter.Order
	chapterDTO.CourseID = chapter.CourseID
	chapterDTO.Course = chapter.Course

	return chapterDTO
}

func LessonsToDTO(lessons []entities.Lesson) []entitiesDTO.LessonDTO {
	if len(lessons) == 0 {
		return []entitiesDTO.LessonDTO{}
	}

	var lessonsDTO []entitiesDTO.LessonDTO
	for i := 0; i < len(lessons); i++ {
		var lessonDTO entitiesDTO.LessonDTO
		lessonDTO.ID = lessons[i].ID
		lessonDTO.Name = lessons[i].Name
		lessonDTO.Description = lessons[i].Description
		lessonDTO.Content = lessons[i].Content
		lessonDTO.Order = lessons[i].Order
		lessonDTO.ChapterID = lessons[i].ChapterID
		lessonDTO.Chapter = lessons[i].Chapter
		lessonsDTO = append(lessonsDTO, lessonDTO)
	}

	return lessonsDTO
}

func LessonToDTO(lesson entities.Lesson) entitiesDTO.LessonDTO {
	var lessonDTO entitiesDTO.LessonDTO
	lessonDTO.ID = lesson.ID
	lessonDTO.Name = lesson.Name
	lessonDTO.Description = lesson.Description
	lessonDTO.Content = lesson.Content
	lessonDTO.Order = lesson.Order
	lessonDTO.ChapterID = lesson.ChapterID
	lessonDTO.Chapter = lesson.Chapter

	return lessonDTO
}
