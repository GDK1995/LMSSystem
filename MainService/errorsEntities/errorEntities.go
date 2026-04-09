package errorsEntities

import (
	"MainService/middleware"
	"net/http"
)

var ErrCourseNotFound = &middleware.AppError{
	Code:    http.StatusNotFound,
	Message: "course not found",
}

var ErrChapterNotFound = &middleware.AppError{
	Code:    http.StatusNotFound,
	Message: "chapter not found",
}

var ErrLessonNotFound = &middleware.AppError{
	Code:    http.StatusNotFound,
	Message: "lesson not found",
}

var ErrInternalServer = &middleware.AppError{
	Code:    http.StatusInternalServerError,
	Message: "internal server error",
}

var ErrBadRequest = &middleware.AppError{
	Code:    http.StatusBadRequest,
	Message: "bad request",
}

var ErrForbidden = &middleware.AppError{
	Code:    http.StatusForbidden,
	Message: "forbidden action",
}
