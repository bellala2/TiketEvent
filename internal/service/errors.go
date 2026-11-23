package service

import "github.com/bellala2/TiketEvent/internal/model"

// Kode error yang dipakai di seluruh aplikasi (sesuai tabel di brief).
const (
	CodeValidation         = "VALIDATION_ERROR"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeVersionConflict    = "VERSION_CONFLICT"
	CodeInvalidState       = "INVALID_STATE"
	CodeInsufficientQuota  = "INSUFFICIENT_QUOTA"
	CodeTicketLimit        = "TICKET_LIMIT_EXCEEDED"
	CodeRequestInProgress  = "REQUEST_IN_PROGRESS"
	CodeIdempotencyReused  = "IDEMPOTENCY_KEY_REUSED"
	CodeInternal           = "INTERNAL_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// AppError adalah error buatan kita. Layer service hanya tahu "kode"
// dan "pesan", tidak tahu soal HTTP. Pemetaan ke status HTTP dilakukan
// di layer handler.
type AppError struct {
	Code    string
	Message string
	Details []model.FieldError
}

func (e *AppError) Error() string {
	return e.Code + ": " + e.Message
}

// NewError membuat AppError biasa (details kosong).
func NewError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewValidationError membuat error validasi beserta daftar field yang salah.
func NewValidationError(message string, details ...model.FieldError) *AppError {
	return &AppError{Code: CodeValidation, Message: message, Details: details}
}
