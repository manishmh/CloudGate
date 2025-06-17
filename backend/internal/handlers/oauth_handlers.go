package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
)

// OAuthState stores OAuth state information
type OAuthState struct {
	State    string `json:"state"`
	Provider string `json:"provider"`
	UserID   string `json:"user_id"`
	Created  int64  `json:"created"`
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        string
}

// GoogleTokenResponse represents Google's token response
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// getGoogleOAuthConfig returns Google OAuth configuration from environment
func getGoogleOAuthConfig() *GoogleOAuthConfig {
	return &GoogleOAuthConfig{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		RedirectURI:  getEnv("GOOGLE_REDIRECT_URI", "http://localhost:8081/oauth/google/callback"),
		Scope:        "openid email profile https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/drive.readonly https://www.googleapis.com/auth/calendar.readonly",
	}
}

// generateOAuthState generates a secure random state parameter
func generateOAuthState() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating OAuth state: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

// GoogleOAuthInitHandler initiates Google OAuth flow
func GoogleOAuthInitHandler(c *gin.Context) {
	config := getGoogleOAuthConfig()

	if config.ClientID == "" || config.ClientSecret == "" {
		log.Printf("Google OAuth not configured - missing ClientID or ClientSecret")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Google OAuth not configured",
			"message": "OAuth credentials not set up for this provider",
		})
		return
	}

	// Get user ID from context (in production, extract from JWT)
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate state parameter
	state := generateOAuthState()
	if state == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in session/cache (for demo, we'll skip this)
	// In production, store state with expiry in Redis or database

	// Build Google OAuth URL
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s&access_type=offline&prompt=consent",
		url.QueryEscape(config.ClientID),
		url.QueryEscape(config.RedirectURI),
		url.QueryEscape(config.Scope),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "google",
	})
}

// GoogleOAuthCallbackHandler handles Google OAuth callback
func GoogleOAuthCallbackHandler(c *gin.Context) {
	config := getGoogleOAuthConfig()

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Google OAuth error: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "OAuth authorization failed",
			"details": errorParam,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing authorization code or state",
		})
		return
	}

	// In production, validate state parameter here
	// For demo, we'll skip state validation

	// Exchange authorization code for access token
	tokenResp, err := exchangeGoogleCode(config, code)
	if err != nil {
		log.Printf("Error exchanging Google code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Google
	userInfo, err := getGoogleUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Google user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID // In production, get from JWT
	err = storeGoogleTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Google tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=google&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// exchangeGoogleCode exchanges authorization code for access token
func exchangeGoogleCode(config *GoogleOAuthConfig, code string) (*GoogleTokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", config.RedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// getGoogleUserInfo retrieves user information from Google
func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// storeGoogleTokens stores Google OAuth tokens in database
func storeGoogleTokens(userID string, tokenResp *GoogleTokenResponse, userInfo *GoogleUserInfo) error {
	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Create app connection record
	connection := map[string]interface{}{
		"status":        constants.StatusConnected,
		"access_token":  tokenResp.AccessToken,
		"refresh_token": tokenResp.RefreshToken,
		"token_type":    tokenResp.TokenType,
		"scope":         tokenResp.Scope,
		"expires_at":    expiresAt.UTC().Format(time.RFC3339),
		"user_email":    userInfo.Email,
		"user_name":     userInfo.Name,
		"connected_at":  time.Now().UTC().Format(time.RFC3339),
	}

	// Update user app connection
	err := services.UpdateUserAppConnection(userID, "google-workspace", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	// Log the connection event
	log.Printf("Google OAuth successful for user %s (email: %s)", userID, userInfo.Email)

	return nil
}

// MicrosoftOAuthInitHandler initiates Microsoft OAuth flow
func MicrosoftOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("MICROSOFT_CLIENT_ID", "")
	redirectURI := getEnv("MICROSOFT_REDIRECT_URI", "http://localhost:8081/oauth/microsoft/callback")

	if clientID == "" {
		log.Printf("Microsoft OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Microsoft OAuth not configured",
			"message": "OAuth credentials not set up for this provider",
		})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	state := generateOAuthState()
	scope := "openid email profile User.Read Mail.Read Calendars.Read Files.Read"

	authURL := fmt.Sprintf(
		"https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "microsoft",
	})
}

// SlackOAuthInitHandler initiates Slack OAuth flow
func SlackOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("SLACK_CLIENT_ID", "")
	redirectURI := getEnv("SLACK_REDIRECT_URI", "http://localhost:8081/oauth/slack/callback")

	if clientID == "" {
		log.Printf("Slack OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Slack OAuth not configured",
			"message": "OAuth credentials not set up for this provider",
		})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	state := generateOAuthState()
	scope := "channels:read,chat:write,users:read,users:read.email"

	authURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(scope),
		url.QueryEscape(redirectURI),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "slack",
	})
}

// GitHubOAuthInitHandler initiates GitHub OAuth flow
func GitHubOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("GITHUB_CLIENT_ID", "")
	redirectURI := getEnv("GITHUB_REDIRECT_URI", "http://localhost:8081/oauth/github/callback")

	if clientID == "" {
		log.Printf("GitHub OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "GitHub OAuth not configured",
			"message": "OAuth credentials not set up for this provider",
		})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	state := generateOAuthState()
	scope := "user:email,repo,read:org"

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "github",
	})
}

// getEnv helper function
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// MicrosoftOAuthCallbackHandler handles Microsoft OAuth callback
func MicrosoftOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("MICROSOFT_CLIENT_ID", "")
	clientSecret := getEnv("MICROSOFT_CLIENT_SECRET", "")
	redirectURI := getEnv("MICROSOFT_REDIRECT_URI", "http://localhost:8081/oauth/microsoft/callback")

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Microsoft OAuth error: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "OAuth authorization failed",
			"details": errorParam,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing authorization code or state",
		})
		return
	}

	// Exchange authorization code for access token
	tokenResp, err := exchangeMicrosoftCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Microsoft code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Microsoft Graph
	userInfo, err := getMicrosoftUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Microsoft user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID // In production, get from JWT
	err = storeMicrosoftTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Microsoft tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=microsoft&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// SlackOAuthCallbackHandler handles Slack OAuth callback
func SlackOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("SLACK_CLIENT_ID", "")
	clientSecret := getEnv("SLACK_CLIENT_SECRET", "")
	redirectURI := getEnv("SLACK_REDIRECT_URI", "http://localhost:8081/oauth/slack/callback")

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Slack OAuth error: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "OAuth authorization failed",
			"details": errorParam,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing authorization code or state",
		})
		return
	}

	// Exchange authorization code for access token
	tokenResp, err := exchangeSlackCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Slack code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Slack
	userInfo, err := getSlackUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Slack user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID // In production, get from JWT
	err = storeSlackTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Slack tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=slack&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.User.Profile.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// GitHubOAuthCallbackHandler handles GitHub OAuth callback
func GitHubOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("GITHUB_CLIENT_ID", "")
	clientSecret := getEnv("GITHUB_CLIENT_SECRET", "")
	redirectURI := getEnv("GITHUB_REDIRECT_URI", "http://localhost:8081/oauth/github/callback")

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("GitHub OAuth error: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "OAuth authorization failed",
			"details": errorParam,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing authorization code or state",
		})
		return
	}

	// Exchange authorization code for access token
	tokenResp, err := exchangeGitHubCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging GitHub code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from GitHub
	userInfo, err := getGitHubUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting GitHub user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID // In production, get from JWT
	err = storeGitHubTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing GitHub tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	email := userInfo.Email
	if email == "" {
		email = userInfo.Login // Use username if email not available
	}
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=github&email=%s&code=success", frontendURL, url.QueryEscape(email))
	c.Redirect(http.StatusFound, redirectURL)
}

// Microsoft Token Response and User Info types
type MicrosoftTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type MicrosoftUserInfo struct {
	ID                string `json:"id"`
	Email             string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	JobTitle          string `json:"jobTitle"`
}

// Slack Token Response and User Info types
type SlackTokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Scope       string    `json:"scope"`
	Team        SlackTeam `json:"team"`
}

type SlackTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackUserInfo struct {
	OK   bool      `json:"ok"`
	User SlackUser `json:"user"`
}

type SlackUser struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	RealName string           `json:"real_name"`
	Profile  SlackUserProfile `json:"profile"`
}

type SlackUserProfile struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image192  string `json:"image_192"`
}

// GitHub Token Response and User Info types
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
	Location  string `json:"location"`
}

// Token exchange functions
func exchangeMicrosoftCode(clientID, clientSecret, redirectURI, code string) (*MicrosoftTokenResponse, error) {
	tokenURL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp MicrosoftTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func exchangeSlackCode(clientID, clientSecret, redirectURI, code string) (*SlackTokenResponse, error) {
	tokenURL := "https://slack.com/api/oauth.v2.access"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp SlackTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func exchangeGitHubCode(clientID, clientSecret, redirectURI, code string) (*GitHubTokenResponse, error) {
	tokenURL := "https://github.com/login/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// User info retrieval functions
func getMicrosoftUserInfo(accessToken string) (*MicrosoftUserInfo, error) {
	userInfoURL := "https://graph.microsoft.com/v1.0/me"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo MicrosoftUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// Use userPrincipalName as email if mail is empty
	if userInfo.Email == "" {
		userInfo.Email = userInfo.UserPrincipalName
	}

	return &userInfo, nil
}

func getSlackUserInfo(accessToken string) (*SlackUserInfo, error) {
	userInfoURL := "https://slack.com/api/users.identity"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo SlackUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func getGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	userInfoURL := "https://api.github.com/user"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// Token storage functions
func storeMicrosoftTokens(userID string, tokenResp *MicrosoftTokenResponse, userInfo *MicrosoftUserInfo) error {
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	connection := map[string]interface{}{
		"status":        constants.StatusConnected,
		"access_token":  tokenResp.AccessToken,
		"refresh_token": tokenResp.RefreshToken,
		"token_type":    tokenResp.TokenType,
		"scope":         tokenResp.Scope,
		"expires_at":    expiresAt.UTC().Format(time.RFC3339),
		"user_email":    userInfo.Email,
		"user_name":     userInfo.DisplayName,
		"connected_at":  time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "microsoft-365", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Microsoft OAuth successful for user %s (email: %s)", userID, userInfo.Email)
	return nil
}

func storeSlackTokens(userID string, tokenResp *SlackTokenResponse, userInfo *SlackUserInfo) error {
	connection := map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"scope":        tokenResp.Scope,
		"user_email":   userInfo.User.Profile.Email,
		"user_name":    userInfo.User.RealName,
		"team_name":    tokenResp.Team.Name,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "slack", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Slack OAuth successful for user %s (email: %s)", userID, userInfo.User.Profile.Email)
	return nil
}

func storeGitHubTokens(userID string, tokenResp *GitHubTokenResponse, userInfo *GitHubUserInfo) error {
	connection := map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"scope":        tokenResp.Scope,
		"user_email":   userInfo.Email,
		"user_name":    userInfo.Name,
		"username":     userInfo.Login,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "github", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("GitHub OAuth successful for user %s (email: %s)", userID, userInfo.Email)
	return nil
}
