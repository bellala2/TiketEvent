// Package response berisi penulis envelope JSON yang dipakai handler DAN middleware.
// Dipisah dari package handler supaya tidak terjadi import melingkar
// (handler memakai middleware, middleware memakai response).
package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bellala2/TiketEvent/internal/model"
	"github.com/bellala2/TiketEvent/internal/service"
)

// statusByCode memetakan kode error ke status HTTP. Hanya ada di satu tempat ini.
var statusByCode = map[string]int{
	service.CodeValidation:         http.StatusBadRequest,
	service.CodeUnauthorized:       http.StatusUnauthorized,
	service.CodeForbidden:          http.StatusForbidden,
	service.CodeNotFound:           http.StatusNotFound,
	service.CodeConflict:           http.StatusConflict,
	service.CodeVersionConflict:    http.StatusConflict,
	service.CodeInvalidState:       http.StatusConflict,
	service.CodeInsufficientQuota:  http.StatusConflict,
	service.CodeTicketLimit:        http.StatusConflict,
	service.CodeRequestInProgress:  http.StatusConflict,
	service.CodeIdempotencyReused:  http.StatusUnprocessableEntity,
	service.CodeInternal:           http.StatusInternalServerError,
	service.CodeServiceUnavailable: http.StatusServiceUnavailable,
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("gagal menulis response", "error", err)
	}
}

// WriteSuccess menulis envelope sukses. meta boleh nil (untuk data tunggal).
func WriteSuccess(w http.ResponseWriter, status int, message string, data any, meta *model.Meta) {
	writeJSON(w, status, model.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// WriteError mengubah error apa pun menjadi envelope error.
// Error yang bukan AppError dianggap tak terduga: detail teknisnya hanya
// ditulis ke log, client cukup menerima pesan umum.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *service.AppError
	if !errors.As(err, &appErr) {
		slog.Error("error tak terduga", "error", err)
		appErr = service.NewError(service.CodeInternal, "Terjadi kesalahan pada server")
	}

	status, ok := statusByCode[appErr.Code]
	if !ok {
		status = http.StatusInternalServerError
	}

	details := appErr.Details
	if details == nil {
		details = []model.FieldError{} // harus [] bukan null
	}

	writeJSON(w, status, model.ErrorResponse{
		Success: false,
		Message: appErr.Message,
		Error:   model.ErrorBody{Code: appErr.Code, Details: details},
	})
}
