package http

import (
	"FOLLOWER/internal/models"
	"FOLLOWER/internal/service"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc service.FollowService
}

func NewHandler(svc service.FollowService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	// Protected routes (Assuming middleware sets 'subject' in context)
	r.Post("/v1/social/follow", h.FollowUser)
	r.Post("/v1/social/unfollow", h.UnfollowUser)
	r.Get("/v1/social/recommendations", h.GetRecommendations)

	// Public or Protected routes
	r.Get("/v1/social/users/{user_id}/stats", h.GetUserStats)
	r.Get("/v1/social/users/{user_id}/followers", h.GetFollowers)
	r.Get("/v1/social/users/{user_id}/following", h.GetFollowing)
}

func ctxSubject(ctx context.Context) (string, bool) {
	v := ctx.Value("subject")
	sub, ok := v.(string)
	return sub, ok
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, ok := ctxSubject(r.Context())
	if !ok || followerID == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	var req models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.svc.FollowUser(r.Context(), followerID, req.TargetID); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "followed"})
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, ok := ctxSubject(r.Context())
	if !ok || followerID == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	var req models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.svc.UnfollowUser(r.Context(), followerID, req.TargetID); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "unfollowed"})
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, ok := ctxSubject(r.Context())
	if !ok || userID == "" {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	recs, err := h.svc.GetRecommendations(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, recs)
}

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	stats, err := h.svc.GetUserStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	list, err := h.svc.GetFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	list, err := h.svc.GetFollowing(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
