package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/hex/gatekeeper-go/internal/audit"
	"github.com/hex/gatekeeper-go/internal/auth"
	"github.com/hex/gatekeeper-go/internal/service/access"
)

type authorizeRequest struct {
	Action     string            `json:"action"`
	Resource   string            `json:"resource"`
	Attributes map[string]string `json:"attributes"`
}

func registerRoutes(mux *http.ServeMux, service *access.Service, auditStore *audit.Store) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	mux.HandleFunc("/v1/authorize", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		subject, err := auth.SubjectFromRequest(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		var req authorizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		if req.Action == "" || req.Resource == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "action and resource are required"})
			return
		}
		if req.Attributes == nil {
			req.Attributes = map[string]string{}
		}

		response := service.Authorize(access.Request{
			Subject:    subject,
			Action:     req.Action,
			Resource:   req.Resource,
			Attributes: req.Attributes,
		})
		writeJSON(w, http.StatusOK, response)
	})

	mux.HandleFunc("/v1/audit/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		writeJSON(w, http.StatusOK, auditStore.List())
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
