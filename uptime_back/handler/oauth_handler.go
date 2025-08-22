package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SohamChatterG/uptime/config"
	"github.com/SohamChatterG/uptime/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	googleOAuthConfig *oauth2.Config
	githubOAuthConfig *oauth2.Config
	userService       *service.UserService
}

func NewOAuthHandler(cfg *config.Config, svc *service.UserService) *OAuthHandler {
	return &OAuthHandler{
		googleOAuthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		githubOAuthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
		userService: svc,
	}
}

func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.googleOAuthConfig.AuthCodeURL("random-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	jwtToken, err := h.userService.FindOrCreateUser(context.Background(), userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/auth/callback?token=%s", jwtToken), http.StatusSeeOther)
}

func (h *OAuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	url := h.githubOAuthConfig.AuthCodeURL("random-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := h.githubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Login string `json:"login"`
	}
	json.Unmarshal(body, &userInfo)

	if userInfo.Email == "" {
		userInfo.Email = fmt.Sprintf("%s@github.com", userInfo.Login) // Create a placeholder email
	}
	if userInfo.Name == "" {
		userInfo.Name = userInfo.Login
	}

	jwtToken, err := h.userService.FindOrCreateUser(context.Background(), userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/auth/callback?token=%s", jwtToken), http.StatusSeeOther)
}
