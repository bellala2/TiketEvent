package model

// Meta dipakai untuk response daftar data (paginasi).
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// FieldError menjelaskan satu field yang gagal validasi.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorBody adalah isi objek "error" pada envelope error.
type ErrorBody struct {
	Code    string       `json:"code"`
	Details []FieldError `json:"details"`
}

// SuccessResponse adalah envelope untuk response sukses (2xx).
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"` // hanya muncul untuk daftar data
}

// ErrorResponse adalah envelope untuk response gagal (4xx dan 5xx).
type ErrorResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Error   ErrorBody `json:"error"`
}
