package errors

import "errors"

//
// -------- Common Errors --------
//

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)

//
// -------- Task Domain Errors --------
//

var (
	ErrInvalidTitle   = errors.New("title cannot be empty")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidDueDate = errors.New("invalid due date")
)

//
// -------- Resource Errors --------
//

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

//
// -------- Pagination Errors --------
//

var (
	ErrInvalidLimit  = errors.New("invalid limit")
	ErrInvalidOffset = errors.New("invalid offset")
)
