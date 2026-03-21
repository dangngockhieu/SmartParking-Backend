package errors

// Format Response Error
type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewBadRequest(message string) *AppError {
	return NewAppError(message, 400)
}

func NewUnauthorized(message string) *AppError {
	return NewAppError(message, 401)
}

func NewForbidden(message string) *AppError {
	return NewAppError(message, 403)
}

func NewNotFound(message string) *AppError {
	return NewAppError(message, 404)
}

func NewConflict(message string) *AppError {
	return NewAppError(message, 409)
}

func NewInternal(message string) *AppError {
	return NewAppError(message, 500)
}
