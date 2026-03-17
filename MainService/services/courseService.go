package services

import (
	"MainService/entities"
	"MainService/entitiesDTO"
	"MainService/errorsEntities"
	"MainService/mappers"
	"MainService/repositories"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseService interface {
	AddCourseS(course entities.Course) (uint, error)
	GetCoursesS() ([]entitiesDTO.CourseDTO, error)
	GetCourseByIDS(courseID uint) (*entitiesDTO.CourseDTO, error)
	DeleteCourseS(courseID uint) error
	UpdateCurseS(updCourse entitiesDTO.CourseDTO) error
}

type courseService struct {
	courseRepository repositories.CourseRepository
}

func NewCourseService(courseRepository repositories.CourseRepository) CourseService {
	return &courseService{courseRepository: courseRepository}
}

func (cs *courseService) AddCourseS(course entities.Course) (uint, error) {
	logrus.Info("Creating new course")
	logrus.Debugf("Course details: %+v", course)
	courseID, err := cs.courseRepository.AddCourse(course)
	if err != nil {
		logrus.Error("Failed to add course: ", err)
		return 0, errorsEntities.ErrInternalServer
	}

	logrus.Info("Course added successfully")
	return courseID, nil
}

func (cs *courseService) GetCoursesS() ([]entitiesDTO.CourseDTO, error) {
	logrus.Info("Getting courses")
	courses, err := cs.courseRepository.GetCourses()
	if err != nil {
		logrus.Error("Failed to get courses from repository: ", err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found %d courses: %+v", len(courses), courses)

	coursesDTO := mappers.CoursesToDTO(courses)

	logrus.Info("Courses successfully converted to DTO")
	return coursesDTO, nil
}

func (cs *courseService) GetCourseByIDS(courseID uint) (*entitiesDTO.CourseDTO, error) {
	logrus.Infof("Getting course by id %d", courseID)
	course, err := cs.courseRepository.GetCourseByID(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Course with id %d not found", courseID)
			return nil, errorsEntities.ErrCourseNotFound
		}

		logrus.Errorf("Failed to get course by id %d: %v", courseID, err)
		return nil, errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Found course: %+v", course)

	courseDTO := mappers.CourseToDTO(course)
	logrus.Info("Course successfully converted to DTO")

	return &courseDTO, nil
}

func (cs *courseService) DeleteCourseS(courseID uint) error {
	logrus.Infof("Deleting course by id %d", courseID)
	err := cs.courseRepository.DeleteCourse(courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Course with id %d not found", courseID)
			return errorsEntities.ErrCourseNotFound
		}

		logrus.Error("Failed to delete course from repository: ", err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Info("Course successfully deleted")
	return nil
}

func (cs *courseService) UpdateCurseS(updCourse entitiesDTO.CourseDTO) error {
	logrus.Infof("Updating course with id %d", updCourse.ID)
	course, err := cs.courseRepository.GetCourseByID(updCourse.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("Course with id %d not found", updCourse.ID)
			return errorsEntities.ErrCourseNotFound
		}

		logrus.Errorf("Failed to get course by id %d: %v", updCourse.ID, err)
		return errorsEntities.ErrInternalServer
	}

	logrus.Debugf("Current course data: %+v", course)

	if updCourse.Name != "" && course.Name != updCourse.Name {
		logrus.Debugf("Updating Name: %s to %s", course.Name, updCourse.Name)
		course.Name = updCourse.Name
	}

	if updCourse.Description != "" && course.Description != updCourse.Description {
		logrus.Debugf("Updating Description: %s to %s", course.Description, updCourse.Description)
		course.Description = updCourse.Description
	}

	errTwo := cs.courseRepository.UpdateCurse(course)

	if errTwo != nil {
		if errors.Is(errTwo, gorm.ErrRecordNotFound) {
			logrus.Warnf("Course with id %d not found during update", updCourse.ID)
			return errorsEntities.ErrCourseNotFound
		}

		logrus.Errorf("Failed to update course id %d: %v", updCourse.ID, errTwo)
		return errorsEntities.ErrInternalServer
	}

	logrus.Infof("Course with id %d successfully updated", updCourse.ID)
	return nil
}
