package http

import (
	"encoding/json"
	"net/http"

	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/iamtvk/jsontransformer/internal/service"
)

type Handler struct {
	transformerService *service.TransformerService
	// logger             Logger
}

func NewHandler(transformerService *service.TransformerService) *Handler {
	return &Handler{
		transformerService: transformerService,
	}
}

func (h *Handler) Transform(w http.ResponseWriter, r *http.Request) {
	var req models.TransformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := h.transformerService.Transform(r.Context(), &req)
	if err != nil {
		http.Error(w, "Error Transforming", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/transform":
		h.Transform(w, r)
	default:
		http.NotFound(w, r)
	}
}
