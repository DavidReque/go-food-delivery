package errors

import (
	"errors"
	"fmt"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(
		fmt.Sprintf("domain_events event already exists in event registry"),
	)
	InvalidEventTypeError = errors.New("invalid event type")
)
