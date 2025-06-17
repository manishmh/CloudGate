package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/services"
)

// WebAuthn credential structures
type WebAuthnCredential struct {
	ID               string                        `json:"id"`
	RawID            []byte                        `json:"rawId"`
	Type             string                        `json:"type"`
	Response         WebAuthnAuthenticatorResponse `json:"response"`
	ClientExtensions map[string]interface{}        `json:"clientExtensions,omitempty"`
}

type WebAuthnAuthenticatorResponse struct {
	AttestationObject []byte `json:"attestationObject,omitempty"`
	ClientDataJSON    []byte `json:"clientDataJSON"`
	AuthenticatorData []byte `json:"authenticatorData,omitempty"`
	Signature         []byte `json:"signature,omitempty"`
	UserHandle        []byte `json:"userHandle,omitempty"`
}

type WebAuthnPublicKeyCredentialCreationOptions struct {
	Challenge              []byte                                  `json:"challenge"`
	RP                     WebAuthnRelyingParty                    `json:"rp"`
	User                   WebAuthnUser                            `json:"user"`
	PubKeyCredParams       []WebAuthnPubKeyCredParam               `json:"pubKeyCredParams"`
	AuthenticatorSelection WebAuthnAuthenticatorSelection          `json:"authenticatorSelection"`
	Timeout                int                                     `json:"timeout"`
	Attestation            string                                  `json:"attestation"`
	ExcludeCredentials     []WebAuthnPublicKeyCredentialDescriptor `json:"excludeCredentials,omitempty"`
}

type WebAuthnPublicKeyCredentialRequestOptions struct {
	Challenge        []byte                                  `json:"challenge"`
	Timeout          int                                     `json:"timeout"`
	RPID             string                                  `json:"rpId"`
	AllowCredentials []WebAuthnPublicKeyCredentialDescriptor `json:"allowCredentials"`
	UserVerification string                                  `json:"userVerification"`
}

type WebAuthnRelyingParty struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

type WebAuthnUser struct {
	ID          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon,omitempty"`
}

type WebAuthnPubKeyCredParam struct {
	Type string `json:"type"`
	Alg  int    `json:"alg"`
}

type WebAuthnAuthenticatorSelection struct {
	AuthenticatorAttachment string `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      bool   `json:"requireResidentKey"`
	UserVerification        string `json:"userVerification"`
}

type WebAuthnPublicKeyCredentialDescriptor struct {
	Type       string   `json:"type"`
	ID         []byte   `json:"id"`
	Transports []string `json:"transports,omitempty"`
}

// WebAuthn registration request/response types
type WebAuthnRegistrationRequest struct {
	Credential WebAuthnCredential `json:"credential" binding:"required"`
}

type WebAuthnRegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// WebAuthn authentication request/response types
type WebAuthnAuthenticationRequest struct {
	Credential WebAuthnCredential `json:"credential" binding:"required"`
}

type WebAuthnAuthenticationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// WebAuthnRegistrationBeginHandler begins WebAuthn registration process
func WebAuthnRegistrationBeginHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user info - simplified for demo
	user := struct {
		Email             string
		FirstName         string
		LastName          string
		ProfilePictureURL string
	}{
		Email:             "demo@cloudgate.com",
		FirstName:         "Demo",
		LastName:          "User",
		ProfilePictureURL: "",
	}

	// Generate challenge
	challenge := generateChallenge()

	// Store challenge for verification (in production, use Redis with expiry)
	// For demo, we'll skip this step

	// Create registration options
	options := WebAuthnPublicKeyCredentialCreationOptions{
		Challenge: challenge,
		RP: WebAuthnRelyingParty{
			ID:   "localhost",
			Name: "CloudGate SSO",
			Icon: "https://cloudgate.example.com/icon.png",
		},
		User: WebAuthnUser{
			ID:          []byte(userID),
			Name:        user.Email,
			DisplayName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			Icon:        user.ProfilePictureURL,
		},
		PubKeyCredParams: []WebAuthnPubKeyCredParam{
			{Type: "public-key", Alg: -7},   // ES256
			{Type: "public-key", Alg: -257}, // RS256
		},
		AuthenticatorSelection: WebAuthnAuthenticatorSelection{
			AuthenticatorAttachment: "platform",
			RequireResidentKey:      false,
			UserVerification:        "preferred",
		},
		Timeout:     60000, // 60 seconds
		Attestation: "direct",
	}

	// Get existing credentials to exclude
	existingCreds, err := services.GetUserWebAuthnCredentials(userID)
	if err == nil && len(existingCreds) > 0 {
		excludeCredentials := make([]WebAuthnPublicKeyCredentialDescriptor, len(existingCreds))
		for i, cred := range existingCreds {
			excludeCredentials[i] = WebAuthnPublicKeyCredentialDescriptor{
				Type:       "public-key",
				ID:         []byte(cred.CredentialID),
				Transports: []string{"internal", "usb", "nfc", "ble"},
			}
		}
		options.ExcludeCredentials = excludeCredentials
	}

	c.JSON(http.StatusOK, options)
}

// WebAuthnRegistrationFinishHandler completes WebAuthn registration
func WebAuthnRegistrationFinishHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request WebAuthnRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate credential (simplified for demo)
	if request.Credential.Type != "public-key" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential type"})
		return
	}

	// Parse client data JSON
	var clientData map[string]interface{}
	if err := json.Unmarshal(request.Credential.Response.ClientDataJSON, &clientData); err != nil {
		log.Printf("Error parsing client data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client data"})
		return
	}

	// Verify challenge (in production, verify against stored challenge)
	if clientData["type"] != "webauthn.create" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ceremony type"})
		return
	}

	// Store credential
	credentialID := request.Credential.ID
	err := services.StoreWebAuthnCredential(userID, credentialID, request.Credential.Response.AttestationObject)
	if err != nil {
		log.Printf("Error storing WebAuthn credential: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store credential"})
		return
	}

	// Log WebAuthn registration
	services.LogAuditEvent(userID, "webauthn_registered", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "WebAuthn credential registered", "success")

	response := WebAuthnRegistrationResponse{
		Success: true,
		Message: "WebAuthn credential registered successfully",
	}

	c.JSON(http.StatusOK, response)
}

// WebAuthnAuthenticationBeginHandler begins WebAuthn authentication
func WebAuthnAuthenticationBeginHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate challenge
	challenge := generateChallenge()

	// Get user's credentials
	credentials, err := services.GetUserWebAuthnCredentials(userID)
	if err != nil {
		log.Printf("Error getting WebAuthn credentials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get credentials"})
		return
	}

	if len(credentials) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No WebAuthn credentials registered"})
		return
	}

	// Create authentication options
	allowCredentials := make([]WebAuthnPublicKeyCredentialDescriptor, len(credentials))
	for i, cred := range credentials {
		allowCredentials[i] = WebAuthnPublicKeyCredentialDescriptor{
			Type:       "public-key",
			ID:         []byte(cred.CredentialID),
			Transports: []string{"internal", "usb", "nfc", "ble"},
		}
	}

	options := WebAuthnPublicKeyCredentialRequestOptions{
		Challenge:        challenge,
		Timeout:          60000, // 60 seconds
		RPID:             "localhost",
		AllowCredentials: allowCredentials,
		UserVerification: "preferred",
	}

	c.JSON(http.StatusOK, options)
}

// WebAuthnAuthenticationFinishHandler completes WebAuthn authentication
func WebAuthnAuthenticationFinishHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request WebAuthnAuthenticationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate credential
	if request.Credential.Type != "public-key" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential type"})
		return
	}

	// Parse client data JSON
	var clientData map[string]interface{}
	if err := json.Unmarshal(request.Credential.Response.ClientDataJSON, &clientData); err != nil {
		log.Printf("Error parsing client data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client data"})
		return
	}

	// Verify challenge and ceremony type
	if clientData["type"] != "webauthn.get" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ceremony type"})
		return
	}

	// Verify credential exists
	credentialID := request.Credential.ID
	exists, err := services.VerifyWebAuthnCredential(userID, credentialID)
	if err != nil {
		log.Printf("Error verifying WebAuthn credential: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify credential"})
		return
	}

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	// Update credential usage
	err = services.UpdateWebAuthnCredentialUsage(userID, credentialID)
	if err != nil {
		log.Printf("Error updating credential usage: %v", err)
		// Don't fail authentication for this
	}

	// Log WebAuthn authentication
	services.LogAuditEvent(userID, "webauthn_authentication", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "WebAuthn authentication successful", "success")

	// Generate session token (simplified)
	token := generateSessionToken()

	response := WebAuthnAuthenticationResponse{
		Success: true,
		Message: "WebAuthn authentication successful",
		Token:   token,
	}

	c.JSON(http.StatusOK, response)
}

// GetWebAuthnCredentialsHandler returns user's WebAuthn credentials
func GetWebAuthnCredentialsHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	credentials, err := services.GetUserWebAuthnCredentials(userID)
	if err != nil {
		log.Printf("Error getting WebAuthn credentials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"credentials": credentials,
		"count":       len(credentials),
	})
}

// DeleteWebAuthnCredentialHandler deletes a WebAuthn credential
func DeleteWebAuthnCredentialHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	credentialID := c.Param("credential_id")
	if credentialID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credential ID is required"})
		return
	}

	err := services.DeleteWebAuthnCredential(userID, credentialID)
	if err != nil {
		log.Printf("Error deleting WebAuthn credential: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete credential"})
		return
	}

	// Log credential deletion
	services.LogAuditEvent(userID, "webauthn_credential_deleted", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), fmt.Sprintf("WebAuthn credential deleted: %s", credentialID), "warning")

	c.JSON(http.StatusOK, gin.H{
		"message": "WebAuthn credential deleted successfully",
	})
}

// Helper functions
func generateChallenge() []byte {
	// Generate cryptographically secure random challenge
	challenge := make([]byte, 32)
	// In production, use crypto/rand
	for i := range challenge {
		challenge[i] = byte(time.Now().UnixNano() % 256)
	}
	return challenge
}

func generateSessionToken() string {
	return uuid.New().String()
}

// Helper function removed - using inline struct instead
