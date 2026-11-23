package repository

import "errors"

// Error dari layer repository. Layer service yang menerjemahkannya
// menjadi AppError (kode + pesan untuk client).
var (
	ErrNotFound       = errors.New("data tidak ditemukan")
	ErrDuplicateEmail = errors.New("email sudah terdaftar")
)
