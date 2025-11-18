package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"stakeholders/internal/model"
	"stakeholders/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc service.ProfileService
}

func NewHandler(svc service.ProfileService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Get("/v1/profiles/me", h.GetMyProfile)
	r.Get("/v1/profiles/{user_id}", h.GetProfile)
	r.Put("/v1/profiles/me", h.UpdateMyProfile)
	r.Post("/v1/admin/profiles/{user_id}/block", h.BlockUser)
}

func ctxSubject(ctx context.Context) (string, bool) {
	v := ctx.Value("subject")
	sub, ok := v.(string)
	return sub, ok
}

func ctxRoles(ctx context.Context) ([]string, bool) {
	v := ctx.Value("roles")
	rs, ok := v.([]string)
	return rs, ok
}
func hasRole(ctx context.Context, want string) bool {
	rs, ok := ctxRoles(ctx)
	if !ok {
		return false
	}
	for _, r := range rs {
		if r == want {
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	sub, ok := ctxSubject(r.Context())
	if !ok || sub == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	p, err := h.svc.GetProfile(r.Context(), sub)
	if err != nil {
		if errors.Is(err, service.ErrProfileExists) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	sub, _ := ctxSubject(r.Context())
	if sub != userID && !hasRole(r.Context(), string(model.RoleAdmin)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	p, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	sub, ok := ctxSubject(r.Context())
	if !ok || sub == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	var in struct {
		FirstName      *string `json:"first_name"`
		LastName       *string `json:"last_name"`
		ProfilePicture *string `json:"profile_picture"`
		Biography      *string `json:"biography"`
		Motto          *string `json:"motto"`
		Role           *string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Load existing
	existing, err := h.svc.GetProfile(r.Context(), sub)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Apply updates
	if in.FirstName != nil {
		existing.FirstName = *in.FirstName
	}
	if in.LastName != nil {
		existing.LastName = *in.LastName
	}
	if in.ProfilePicture != nil {
		existing.ProfilePicture = *in.ProfilePicture
	}
	if in.Biography != nil {
		existing.Biography = *in.Biography
	}
	if in.Motto != nil {
		existing.Motto = *in.Motto
	}
	// Role change allowed only for admins
	if in.Role != nil && hasRole(r.Context(), string(model.RoleAdmin)) {
		existing.Role = model.Role(*in.Role)
	}
	updated, err := h.svc.UpdateProfile(r.Context(), existing)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) BlockUser(w http.ResponseWriter, r *http.Request) {
	if !hasRole(r.Context(), string(model.RoleAdmin)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	userID := chi.URLParam(r, "user_id")
	adminID, _ := ctxSubject(r.Context())
	if adminID == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	if err := h.svc.BlockUser(r.Context(), userID, adminID); err != nil {
		if errors.Is(err, service.ErrCannotBlockAdmin) {
			http.Error(w, "cannot block admin", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
