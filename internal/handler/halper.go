package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bellala2/TiketEvent/internal/model"
	"github.com/bellala2/TiketEvent/internal/service"
)

// decodeJSON membaca body request ke struct tujuan.
// Body dibatasi 1 MB supaya tidak bisa dipakai untuk membanjiri memori server.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return service.NewValidationError("Body request tidak valid",
			model.FieldError{Field: "body", Message: "harus berupa JSON yang valid"})
	}
	return nil
}
