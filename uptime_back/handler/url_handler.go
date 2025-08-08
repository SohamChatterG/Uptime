package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SohamChatterG/uptime/middleware"
	"github.com/SohamChatterG/uptime/service"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{urlService: svc}
}

// Helper to write JSON error responses
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Helper to write successful JSON responses
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper to get user ID from context
func getUserIDFromContext(r *http.Request) (primitive.ObjectID, bool) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return primitive.NilObjectID, false
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return primitive.NilObjectID, false
	}
	return userID, true
}

func (h *URLHandler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeError(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newURL, err := h.urlService.CreateURL(r.Context(), reqBody.Name, reqBody.URL, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, newURL, http.StatusCreated)
}

func (h *URLHandler) GetURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeError(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	urls, err := h.urlService.GetURLsForUser(r.Context(), userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, urls, http.StatusOK)
}

func (h *URLHandler) DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeError(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	urlID := mux.Vars(r)["id"]
	err := h.urlService.DeleteURL(r.Context(), urlID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]string{"message": "URL deleted successfully"}, http.StatusOK)
}

func (h *URLHandler) GetURLHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		writeError(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	urlID := mux.Vars(r)["id"]
	history, err := h.urlService.GetURLHistory(r.Context(), urlID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, history, http.StatusOK)
}
