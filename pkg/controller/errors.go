package controller

type ErrorType int

const (
	ErrUnsupportedMediaType ErrorType = iota + 1
	ErrDecodingJSON
	ErrParsingForm
	ErrMissingName
	ErrMissingRecord
	ErrMissingDateOfBirth
	ErrInvalidDateOfBirthFormat
	ErrDatabaseError
)

type TypedError struct {
	kind ErrorType
	message string
}

func NewTypedError(kind ErrorType, message string) *TypedError{
	return &TypedError{kind: kind, message: message}
}

func (e *TypedError) Error() string {
	return e.message
}
