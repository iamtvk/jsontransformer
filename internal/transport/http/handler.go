package http

import (
	"encoding/json"
	"log"
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

func (h *Handler) TransformHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("http transform req recived from", r.RemoteAddr)
	var req models.TransformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := h.transformerService.Transform(r.Context(), &req)

	if err != nil {
		log.Printf("error transforming %v", err.Error())
		http.Error(w, "Error Transforming", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
	log.Println("tranform req fulfilled for", r.RemoteAddr)
}

func (h *Handler) CreateScriptHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("http create script req recived from", r.RemoteAddr)
	var req models.CreateScriptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.transformerService.CreateScript(r.Context(), models.TransformationScript{
		Identifier:  req.ScriptIdentifier,
		Description: req.Description,
		Script:      req.Script,
		CreatedBy:   req.CreatedBy,
		Name:        req.Name,
	})
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		resp := models.CreateScriptResponse{Error: err.Error(), Success: false}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := models.CreateScriptResponse{Error: "", Success: true}
	json.NewEncoder(w).Encode(resp)
	log.Println("create-script req fulfilled for", r.RemoteAddr)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/transform":
		h.TransformHandler(w, r)
	case "/create-script":
		h.CreateScriptHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
