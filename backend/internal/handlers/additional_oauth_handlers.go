package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
)

// Salesforce OAuth handlers
func SalesforceOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("SALESFORCE_CLIENT_ID", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/salesforce/callback"

	if clientID == "" {
		log.Printf("Salesforce OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Salesforce OAuth not configured",
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
	scope := "openid email profile api"

	authURL := fmt.Sprintf(
		"https://login.salesforce.com/services/oauth2/authorize?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "salesforce",
	})
}

func SalesforceOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("SALESFORCE_CLIENT_ID", "")
	clientSecret := getEnv("SALESFORCE_CLIENT_SECRET", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/salesforce/callback"

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Salesforce OAuth error: %s", errorParam)
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
	tokenResp, err := exchangeSalesforceCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Salesforce code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Salesforce
	userInfo, err := getSalesforceUserInfo(tokenResp.AccessToken, tokenResp.InstanceURL)
	if err != nil {
		log.Printf("Error getting Salesforce user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID
	err = storeSalesforceTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Salesforce tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=salesforce&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// Jira OAuth handlers
func JiraOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("JIRA_CLIENT_ID", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/jira/callback"

	if clientID == "" {
		log.Printf("Jira OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Jira OAuth not configured",
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
	scope := "read:jira-user read:jira-work write:jira-work"

	authURL := fmt.Sprintf(
		"https://auth.atlassian.com/authorize?audience=api.atlassian.com&client_id=%s&scope=%s&redirect_uri=%s&state=%s&response_type=code&prompt=consent",
		url.QueryEscape(clientID),
		url.QueryEscape(scope),
		url.QueryEscape(redirectURI),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "jira",
	})
}

func JiraOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("JIRA_CLIENT_ID", "")
	clientSecret := getEnv("JIRA_CLIENT_SECRET", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/jira/callback"

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Jira OAuth error: %s", errorParam)
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
	tokenResp, err := exchangeJiraCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Jira code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Jira
	userInfo, err := getJiraUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Jira user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID
	err = storeJiraTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Jira tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=jira&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.EmailAddress))
	c.Redirect(http.StatusFound, redirectURL)
}

// Notion OAuth handlers
func NotionOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("NOTION_CLIENT_ID", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/notion/callback"

	if clientID == "" {
		log.Printf("Notion OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Notion OAuth not configured",
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

	authURL := fmt.Sprintf(
		"https://api.notion.com/v1/oauth/authorize?client_id=%s&response_type=code&owner=user&redirect_uri=%s&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "notion",
	})
}

func NotionOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("NOTION_CLIENT_ID", "")
	clientSecret := getEnv("NOTION_CLIENT_SECRET", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/notion/callback"

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Notion OAuth error: %s", errorParam)
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
	tokenResp, err := exchangeNotionCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Notion code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Notion
	userInfo, err := getNotionUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Notion user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID
	err = storeNotionTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Notion tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=notion&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Person.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// Dropbox OAuth handlers
func DropboxOAuthInitHandler(c *gin.Context) {
	clientID := getEnv("DROPBOX_CLIENT_ID", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/dropbox/callback"

	if clientID == "" {
		log.Printf("Dropbox OAuth not configured - missing ClientID")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Dropbox OAuth not configured",
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

	authURL := fmt.Sprintf(
		"https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		state,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"provider": "dropbox",
	})
}

func DropboxOAuthCallbackHandler(c *gin.Context) {
	clientID := getEnv("DROPBOX_CLIENT_ID", "")
	clientSecret := getEnv("DROPBOX_CLIENT_SECRET", "")
	redirectURI := getEnv("BACKEND_URL", "http://localhost:8081") + "/oauth/dropbox/callback"

	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		log.Printf("Dropbox OAuth error: %s", errorParam)
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
	tokenResp, err := exchangeDropboxCode(clientID, clientSecret, redirectURI, code)
	if err != nil {
		log.Printf("Error exchanging Dropbox code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange authorization code",
		})
		return
	}

	// Get user information from Dropbox
	userInfo, err := getDropboxUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting Dropbox user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// Store tokens in database
	userID := constants.DemoUserID
	err = storeDropboxTokens(userID, tokenResp, userInfo)
	if err != nil {
		log.Printf("Error storing Dropbox tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store tokens",
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=dropbox&email=%s&code=success", frontendURL, url.QueryEscape(userInfo.Email))
	c.Redirect(http.StatusFound, redirectURL)
}

// Type definitions for additional OAuth providers

// Salesforce types
type SalesforceTokenResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

type SalesforceUserInfo struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

// Jira types
type JiraTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type JiraUserInfo struct {
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

// Notion types
type NotionTokenResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	BotID       string      `json:"bot_id"`
	WorkspaceID string      `json:"workspace_id"`
	Owner       NotionOwner `json:"owner"`
}

type NotionOwner struct {
	Type   string       `json:"type"`
	Person NotionPerson `json:"person"`
}

type NotionPerson struct {
	Email string `json:"email"`
}

type NotionUserInfo struct {
	Object string       `json:"object"`
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Person NotionPerson `json:"person"`
}

// Dropbox types
type DropboxTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	UID         string `json:"uid"`
	AccountID   string `json:"account_id"`
}

type DropboxUserInfo struct {
	AccountID string             `json:"account_id"`
	Name      DropboxUserName    `json:"name"`
	Email     string             `json:"email"`
	Country   string             `json:"country"`
	Locale    string             `json:"locale"`
	Profile   DropboxUserProfile `json:"profile_photo_url"`
}

type DropboxUserName struct {
	GivenName    string `json:"given_name"`
	Surname      string `json:"surname"`
	FamiliarName string `json:"familiar_name"`
	DisplayName  string `json:"display_name"`
}

type DropboxUserProfile struct {
	URL string `json:"url"`
}

// Token exchange functions
func exchangeSalesforceCode(clientID, clientSecret, redirectURI, code string) (*SalesforceTokenResponse, error) {
	tokenURL := "https://login.salesforce.com/services/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

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

	var tokenResp SalesforceTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func exchangeJiraCode(clientID, clientSecret, redirectURI, code string) (*JiraTokenResponse, error) {
	tokenURL := "https://auth.atlassian.com/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
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

	var tokenResp JiraTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func exchangeNotionCode(clientID, clientSecret, redirectURI, code string) (*NotionTokenResponse, error) {
	tokenURL := "https://api.notion.com/v1/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodeBasicAuth(clientID, clientSecret))

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

	var tokenResp NotionTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func exchangeDropboxCode(clientID, clientSecret, redirectURI, code string) (*DropboxTokenResponse, error) {
	tokenURL := "https://api.dropboxapi.com/oauth2/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
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

	var tokenResp DropboxTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// User info retrieval functions
func getSalesforceUserInfo(accessToken, instanceURL string) (*SalesforceUserInfo, error) {
	userInfoURL := instanceURL + "/services/oauth2/userinfo"

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

	var userInfo SalesforceUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func getJiraUserInfo(accessToken string) (*JiraUserInfo, error) {
	userInfoURL := "https://api.atlassian.com/me"

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

	var userInfo JiraUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func getNotionUserInfo(accessToken string) (*NotionUserInfo, error) {
	userInfoURL := "https://api.notion.com/v1/users/me"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Notion-Version", "2022-06-28")

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

	var userInfo NotionUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func getDropboxUserInfo(accessToken string) (*DropboxUserInfo, error) {
	userInfoURL := "https://api.dropboxapi.com/2/users/get_current_account"

	req, err := http.NewRequest("POST", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

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

	var userInfo DropboxUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// Token storage functions
func storeSalesforceTokens(userID string, tokenResp *SalesforceTokenResponse, userInfo *SalesforceUserInfo) error {
	connection := map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"instance_url": tokenResp.InstanceURL,
		"user_email":   userInfo.Email,
		"user_name":    userInfo.DisplayName,
		"username":     userInfo.Username,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "salesforce", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Salesforce OAuth successful for user %s (email: %s)", userID, userInfo.Email)
	return nil
}

func storeJiraTokens(userID string, tokenResp *JiraTokenResponse, userInfo *JiraUserInfo) error {
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	connection := map[string]interface{}{
		"status":        constants.StatusConnected,
		"access_token":  tokenResp.AccessToken,
		"refresh_token": tokenResp.RefreshToken,
		"token_type":    tokenResp.TokenType,
		"scope":         tokenResp.Scope,
		"expires_at":    expiresAt.UTC().Format(time.RFC3339),
		"user_email":    userInfo.EmailAddress,
		"user_name":     userInfo.DisplayName,
		"account_id":    userInfo.AccountID,
		"connected_at":  time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "jira", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Jira OAuth successful for user %s (email: %s)", userID, userInfo.EmailAddress)
	return nil
}

func storeNotionTokens(userID string, tokenResp *NotionTokenResponse, userInfo *NotionUserInfo) error {
	connection := map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"bot_id":       tokenResp.BotID,
		"workspace_id": tokenResp.WorkspaceID,
		"user_email":   userInfo.Person.Email,
		"user_name":    userInfo.Name,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "notion", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Notion OAuth successful for user %s (email: %s)", userID, userInfo.Person.Email)
	return nil
}

func storeDropboxTokens(userID string, tokenResp *DropboxTokenResponse, userInfo *DropboxUserInfo) error {
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	connection := map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"scope":        tokenResp.Scope,
		"expires_at":   expiresAt.UTC().Format(time.RFC3339),
		"user_email":   userInfo.Email,
		"user_name":    userInfo.Name.DisplayName,
		"account_id":   userInfo.AccountID,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	}

	err := services.UpdateUserAppConnection(userID, "dropbox", connection)
	if err != nil {
		return fmt.Errorf("failed to update app connection: %w", err)
	}

	log.Printf("Dropbox OAuth successful for user %s (email: %s)", userID, userInfo.Email)
	return nil
}

// Helper function for basic auth encoding
func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
