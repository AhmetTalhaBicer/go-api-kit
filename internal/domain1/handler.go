package domain1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/username/go-api-kit/internal/response"
	"github.com/username/go-api-kit/internal/validator"
)

// CreateRequest holds the fields accepted when creating a Domain1 record.
type CreateRequest struct {
	Name string `json:"name"`
}

// UpdateRequest holds the fields accepted when updating a Domain1 record.
type UpdateRequest struct {
	Name string `json:"name"`
}

// Handler holds the HTTP handlers for Domain1.
type Handler struct {
	log     *slog.Logger
	service *Service
}

// NewHandler creates a new Domain1 Handler.
// svc must not be nil; pass a real Service backed by a Repository implementation.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	if svc == nil {
		panic("domain1: NewHandler requires a non-nil Service")
	}
	return &Handler{log: log, service: svc}
}

// RegisterRoutes registers all Domain1 routes on the provided ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/domain1", h.List)
	mux.HandleFunc("POST /api/v1/domain1", h.Create)
	mux.HandleFunc("GET /api/v1/domain1/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/v1/domain1/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/domain1/{id}", h.Delete)
}

// List returns all Domain1 records.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll(r.Context())
	if err != nil {
		h.log.Error("domain1: failed to list", "error", err)
		response.InternalError(w)
		return
	}
	response.Success(w, items)
}

// GetByID returns a single Domain1 record by ID, or 404 if not found.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "id must be an integer")
		return
	}
	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w)
			return
		}
		h.log.Error("domain1: failed to get by id", "id", id, "error", err)
		response.InternalError(w)
		return
	}
	response.Success(w, item)
}

// Create persists a new Domain1 record.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	v := validator.New().Required("name", req.Name).MinLength("name", req.Name, 1)
	if !v.Valid() {
		response.BadRequest(w, v.Errors().Error())
		return
	}

	item, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		h.log.Error("domain1: failed to create", "error", err)
		response.InternalError(w)
		return
	}
	response.Created(w, item)
}

// Update modifies an existing Domain1 record.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "id must be an integer")
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	v := validator.New().Required("name", req.Name).MinLength("name", req.Name, 1)
	if !v.Valid() {
		response.BadRequest(w, v.Errors().Error())
		return
	}

	item, err := h.service.Update(r.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w)
			return
		}
		h.log.Error("domain1: failed to update", "id", id, "error", err)
		response.InternalError(w)
		return
	}
	response.Success(w, item)
}

// Delete removes a Domain1 record by ID.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "id must be an integer")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w)
			return
		}
		h.log.Error("domain1: failed to delete", "id", id, "error", err)
		response.InternalError(w)
		return
	}
	response.Success(w, map[string]any{"id": id, "deleted": true})
}
