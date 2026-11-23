package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/bellala2/TiketEvent/internal/response"
	"github.com/bellala2/TiketEvent/internal/service"
)

type HealthHandler struct {
	db            *sql.DB
	queueCapacity int
}

func NewHealthHandler(db *sql.DB, queueCapacity int) *HealthHandler {
	return &HealthHandler{db: db, queueCapacity: queueCapacity}
}

// Check menjawab GET /health.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.WriteError(w, service.NewError(service.CodeServiceUnavailable, "Database tidak dapat dihubungi"))
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Service sehat", map[string]any{
		"status":   "up",
		"database": "up",
		// TODO (Hari 4): "waiting" diganti jumlah job yang benar-benar antre.
		"import_queue": map[string]int{"waiting": 0, "capacity": h.queueCapacity},
	}, nil)
}