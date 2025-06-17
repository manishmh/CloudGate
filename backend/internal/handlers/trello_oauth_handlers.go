package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
)

// In-memory store for OAuth 1.0a request token secrets (for demo purposes)
var (
	requestTokenSecrets = make(map[string]string)
	requestTokenMutex   = sync.RWMutex{}
)

// TrelloOAuthConfig holds Trello OAuth 1.0a configuration
type TrelloOAuthConfig struct {
	APIKey          string // Client ID
	APISecret       string // Client Secret
	CallbackURL     string
	RequestTokenURL string
	AuthorizeURL    string
	AccessTokenURL  string
}

// TrelloRequestTokenResponse represents Trello's request token response
type TrelloRequestTokenResponse struct {
	OAuthToken             string `json:"oauth_token"`
	OAuthTokenSecret       string `json:"oauth_token_secret"`
	OAuthCallbackConfirmed string `json:"oauth_callback_confirmed"`
}

// TrelloAccessTokenResponse represents Trello's access token response
type TrelloAccessTokenResponse struct {
	OAuthToken       string `json:"oauth_token"`
	OAuthTokenSecret string `json:"oauth_token_secret"`
}

// TrelloUserInfo represents Trello user information
type TrelloUserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	URL      string `json:"url"`
}

// getTrelloOAuthConfig returns Trello OAuth 1.0a configuration
func getTrelloOAuthConfig() *TrelloOAuthConfig {
	return &TrelloOAuthConfig{
		APIKey:          getEnv("TRELLO_CLIENT_ID", ""),
		APISecret:       getEnv("TRELLO_CLIENT_SECRET", ""),
		CallbackURL:     getEnv("TRELLO_REDIRECT_URI", getEnv("NEXT_PUBLIC_API_URL", "http://localhost:8081")+"/oauth/trello/callback"),
		RequestTokenURL: "https://trello.com/1/OAuthGetRequestToken",
		AuthorizeURL:    "https://trello.com/1/OAuthAuthorizeToken",
		AccessTokenURL:  "https://trello.com/1/OAuthGetAccessToken",
	}
}

// generateNonce generates a random nonce for OAuth 1.0a
func generateNonce() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

// generateTimestamp generates current timestamp for OAuth 1.0a
func generateTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// generateSignature generates OAuth 1.0a signature
func generateSignature(method, baseURL string, params map[string]string, consumerSecret, tokenSecret string) string {
	// Sort parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	var paramPairs []string
	for _, k := range keys {
		paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
	}
	paramString := strings.Join(paramPairs, "&")

	// Build signature base string
	signatureBaseString := fmt.Sprintf("%s&%s&%s",
		strings.ToUpper(method),
		url.QueryEscape(baseURL),
		url.QueryEscape(paramString))

	// Build signing key
	signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(consumerSecret), url.QueryEscape(tokenSecret))

	// Generate HMAC-SHA1 signature
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(signatureBaseString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}

// buildOAuthHeader builds OAuth 1.0a authorization header
func buildOAuthHeader(params map[string]string) string {
	var headerParts []string
	for k, v := range params {
		if strings.HasPrefix(k, "oauth_") {
			headerParts = append(headerParts, fmt.Sprintf(`%s="%s"`, k, url.QueryEscape(v)))
		}
	}
	return "OAuth " + strings.Join(headerParts, ", ")
}

// TrelloOAuthInitHandler initiates Trello OAuth 1.0a flow
func TrelloOAuthInitHandler(c *gin.Context) {
	config := getTrelloOAuthConfig()

	if config.APIKey == "" || config.APISecret == "" {
		log.Printf("Trello OAuth not configured - missing APIKey or APISecret")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Trello OAuth not configured",
			"message": "OAuth credentials not set up for this provider",
		})
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Step 1: Get request token
	requestToken, requestTokenSecret, err := getTrelloRequestToken(config)
	if err != nil {
		log.Printf("Error getting Trello request token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initiate Trello OAuth",
		})
		return
	}

	// Store request token secret for callback (in production, use Redis/database)
	requestTokenMutex.Lock()
	requestTokenSecrets[requestToken] = requestTokenSecret
	requestTokenMutex.Unlock()

	// Build authorization URL
	authURL := fmt.Sprintf("%s?oauth_token=%s&scope=read,write&expiration=30days&name=CloudGate",
		config.AuthorizeURL,
		url.QueryEscape(requestToken))

	c.JSON(http.StatusOK, gin.H{
		"auth_url":    authURL,
		"provider":    "trello",
		"oauth_token": requestToken,
	})
}

// getTrelloRequestToken gets request token from Trello (Step 1 of OAuth 1.0a)
func getTrelloRequestToken(config *TrelloOAuthConfig) (string, string, error) {
	// OAuth 1.0a parameters
	params := map[string]string{
		"oauth_callback":         config.CallbackURL,
		"oauth_consumer_key":     config.APIKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        generateTimestamp(),
		"oauth_version":          "1.0",
	}

	// Generate signature
	signature := generateSignature("POST", config.RequestTokenURL, params, config.APISecret, "")
	params["oauth_signature"] = signature

	// Build authorization header
	authHeader := buildOAuthHeader(params)

	// Make request
	req, err := http.NewRequest("POST", config.RequestTokenURL, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("request token failed: %s", string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Parse URL-encoded response
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", err
	}

	requestToken := values.Get("oauth_token")
	requestTokenSecret := values.Get("oauth_token_secret")

	if requestToken == "" || requestTokenSecret == "" {
		return "", "", fmt.Errorf("invalid response from Trello: missing tokens")
	}

	return requestToken, requestTokenSecret, nil
}

// TrelloOAuthCallbackHandler handles Trello OAuth 1.0a callback
func TrelloOAuthCallbackHandler(c *gin.Context) {
	config := getTrelloOAuthConfig()

	oauthToken := c.Query("oauth_token")
	oauthVerifier := c.Query("oauth_verifier")

	if oauthToken == "" || oauthVerifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing OAuth parameters",
		})
		return
	}

	// Step 3: Exchange for access token
	// Retrieve the stored request token secret
	requestTokenMutex.RLock()
	requestTokenSecret, exists := requestTokenSecrets[oauthToken]
	requestTokenMutex.RUnlock()

	if !exists {
		log.Printf("Request token secret not found for token: %s", oauthToken)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OAuth state - request token not found",
		})
		return
	}

	// Clean up the stored secret
	requestTokenMutex.Lock()
	delete(requestTokenSecrets, oauthToken)
	requestTokenMutex.Unlock()

	accessToken, accessTokenSecret, err := getTrelloAccessToken(config, oauthToken, oauthVerifier, requestTokenSecret)
	if err != nil {
		log.Printf("Error getting Trello access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange OAuth tokens",
		})
		return
	}

	// Get user information from Trello
	userInfo, err := getTrelloUserInfo(config, accessToken, accessTokenSecret)
	if err != nil {
		log.Printf("Error getting Trello user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID // In production, get from JWT
	err = storeTrelloTokens(userID, accessToken, accessTokenSecret, userInfo)
	if err != nil {
		log.Printf("Error storing Trello tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=trello&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Username))
	c.Redirect(http.StatusFound, redirectURL)
}

// getTrelloAccessToken exchanges request token for access token (Step 3 of OAuth 1.0a)
func getTrelloAccessToken(config *TrelloOAuthConfig, oauthToken, oauthVerifier, requestTokenSecret string) (string, string, error) {
	// OAuth 1.0a parameters
	params := map[string]string{
		"oauth_consumer_key":     config.APIKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        generateTimestamp(),
		"oauth_token":            oauthToken,
		"oauth_verifier":         oauthVerifier,
		"oauth_version":          "1.0",
	}

	// Generate signature (using request token secret)
	signature := generateSignature("POST", config.AccessTokenURL, params, config.APISecret, requestTokenSecret)
	params["oauth_signature"] = signature

	// Build authorization header
	authHeader := buildOAuthHeader(params)

	// Make request
	req, err := http.NewRequest("POST", config.AccessTokenURL, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("access token failed: %s", string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Parse URL-encoded response
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", err
	}

	accessToken := values.Get("oauth_token")
	accessTokenSecret := values.Get("oauth_token_secret")

	if accessToken == "" || accessTokenSecret == "" {
		return "", "", fmt.Errorf("invalid response from Trello: missing access tokens")
	}

	return accessToken, accessTokenSecret, nil
}

// getTrelloUserInfo retrieves user information from Trello API
func getTrelloUserInfo(config *TrelloOAuthConfig, accessToken, accessTokenSecret string) (*TrelloUserInfo, error) {
	userInfoURL := "https://api.trello.com/1/members/me"

	// OAuth 1.0a parameters for API call
	params := map[string]string{
		"oauth_consumer_key":     config.APIKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        generateTimestamp(),
		"oauth_token":            accessToken,
		"oauth_version":          "1.0",
	}

	// Generate signature
	signature := generateSignature("GET", userInfoURL, params, config.APISecret, accessTokenSecret)
	params["oauth_signature"] = signature

	// Build authorization header
	authHeader := buildOAuthHeader(params)

	// Make request
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)

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

	var userInfo TrelloUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// storeTrelloTokens stores Trello OAuth 1.0a tokens in database
func storeTrelloTokens(userID string, accessToken, accessTokenSecret string, userInfo *TrelloUserInfo) error {
	// Create app connection record
	connection := map[string]interface{}{
		"status":              constants.StatusConnected,
		"access_token":        accessToken,
		"access_token_secret": accessTokenSecret, // OAuth 1.0a specific
		"token_type":          "OAuth1.0a",
		"scope":               "read,write",
		"user_id":             userInfo.ID,
		"username":            userInfo.Username,
		"user_name":           userInfo.FullName,
		"user_email":          userInfo.Email,
		"connected_at":        time.Now().UTC().Format(time.RFC3339),
	}

	// Update user app connection
	err := services.UpdateUserAppConnection(userID, "trello", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	// Log the connection event
	log.Printf("Trello OAuth successful for user %s (username: %s)", userID, userInfo.Username)

	return nil
}
