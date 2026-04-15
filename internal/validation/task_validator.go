package validation

import (
	"time"

	"task-manager/internal/domain"
)

func ValidateDueDate(dueDate time.Time) error {
	if dueDate.IsZero() {
		return domain.ErrInvalidDueDate
	}
	if dueDate.Before(time.Now()) {
		return domain.ErrInvalidDueDate
	}
	return nil
}
