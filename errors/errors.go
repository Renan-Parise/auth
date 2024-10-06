package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
	Method  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Error while validating: field '%s': %s. Event occurred in method: %s.", e.Field, e.Message, e.Method)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message, Method: getCallerMethodName()}
}

type QueryError struct {
	Reason string
	Method string
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("Query error: %s. Event occurred in method: %s.", e.Reason, e.Method)
}

func NewQueryError(reason string) *QueryError {
	return &QueryError{Reason: reason, Method: getCallerMethodName()}
}

type DatabaseError struct {
	Reason string
	Method string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Database error: %s. Event occurred in method: %s.", e.Reason, e.Method)
}

func NewDatabaseError(reason string) *DatabaseError {
	return &DatabaseError{Reason: reason, Method: getCallerMethodName()}
}

type ServiceError struct {
	Reason string
	Method string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("Service error: %s. Event occurred in method: %s.", e.Reason, e.Method)
}

func NewServiceError(reason string) *ServiceError {
	return &ServiceError{Reason: reason, Method: getCallerMethodName()}
}

func getCallerMethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, "/")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		if idx := strings.Index(lastPart, "."); idx != -1 {
			return lastPart[idx:]
		}
	}
	return fullName
}
