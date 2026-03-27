package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IvanTime-Kai/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	svc *service.LinkService
}

func NewLinkHandler(svc *service.LinkService) *LinkHandler {
	return &LinkHandler{
		svc: svc,
	}
}

func (h *LinkHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	link, err := h.svc.Shorten(r.Context(), req.URL)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, link)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.svc.Resolve(r.Context(), code)

	if err != nil {
		writeError(w, http.StatusBadRequest, "link nou found")
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.List(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch links")
		return
	}

	writeJSON(w, http.StatusOK, links)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if err := h.svc.Delete(r.Context(), code); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete link")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
