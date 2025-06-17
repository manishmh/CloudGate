package handlers

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
)

// SAML Request structure
type SAMLRequest struct {
	XMLName                     xml.Name   `xml:"urn:oasis:names:tc:SAML:2.0:protocol AuthnRequest"`
	ID                          string     `xml:"ID,attr"`
	Version                     string     `xml:"Version,attr"`
	IssueInstant                string     `xml:"IssueInstant,attr"`
	Destination                 string     `xml:"Destination,attr"`
	ProtocolBinding             string     `xml:"ProtocolBinding,attr"`
	AssertionConsumerServiceURL string     `xml:"AssertionConsumerServiceURL,attr"`
	Issuer                      SAMLIssuer `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
}

type SAMLIssuer struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Value   string   `xml:",chardata"`
}

// SAML Response structure
type SAMLResponse struct {
	XMLName      xml.Name      `xml:"urn:oasis:names:tc:SAML:2.0:protocol Response"`
	ID           string        `xml:"ID,attr"`
	Version      string        `xml:"Version,attr"`
	IssueInstant string        `xml:"IssueInstant,attr"`
	Destination  string        `xml:"Destination,attr"`
	InResponseTo string        `xml:"InResponseTo,attr"`
	Issuer       SAMLIssuer    `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Status       SAMLStatus    `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`
	Assertion    SAMLAssertion `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
}

type SAMLStatus struct {
	XMLName    xml.Name       `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`
	StatusCode SAMLStatusCode `xml:"urn:oasis:names:tc:SAML:2.0:protocol StatusCode"`
}

type SAMLStatusCode struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol StatusCode"`
	Value   string   `xml:"Value,attr"`
}

type SAMLAssertion struct {
	XMLName            xml.Name               `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
	ID                 string                 `xml:"ID,attr"`
	Version            string                 `xml:"Version,attr"`
	IssueInstant       string                 `xml:"IssueInstant,attr"`
	Issuer             SAMLIssuer             `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Subject            SAMLSubject            `xml:"urn:oasis:names:tc:SAML:2.0:assertion Subject"`
	Conditions         SAMLConditions         `xml:"urn:oasis:names:tc:SAML:2.0:assertion Conditions"`
	AttributeStatement SAMLAttributeStatement `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeStatement"`
	AuthnStatement     SAMLAuthnStatement     `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnStatement"`
}

type SAMLSubject struct {
	XMLName             xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:assertion Subject"`
	NameID              SAMLNameID              `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	SubjectConfirmation SAMLSubjectConfirmation `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmation"`
}

type SAMLNameID struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	Format  string   `xml:"Format,attr"`
	Value   string   `xml:",chardata"`
}

type SAMLSubjectConfirmation struct {
	XMLName                 xml.Name                    `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmation"`
	Method                  string                      `xml:"Method,attr"`
	SubjectConfirmationData SAMLSubjectConfirmationData `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmationData"`
}

type SAMLSubjectConfirmationData struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmationData"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr"`
	Recipient    string   `xml:"Recipient,attr"`
	InResponseTo string   `xml:"InResponseTo,attr"`
}

type SAMLConditions struct {
	XMLName             xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:assertion Conditions"`
	NotBefore           string                  `xml:"NotBefore,attr"`
	NotOnOrAfter        string                  `xml:"NotOnOrAfter,attr"`
	AudienceRestriction SAMLAudienceRestriction `xml:"urn:oasis:names:tc:SAML:2.0:assertion AudienceRestriction"`
}

type SAMLAudienceRestriction struct {
	XMLName  xml.Name     `xml:"urn:oasis:names:tc:SAML:2.0:assertion AudienceRestriction"`
	Audience SAMLAudience `xml:"urn:oasis:names:tc:SAML:2.0:assertion Audience"`
}

type SAMLAudience struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Audience"`
	Value   string   `xml:",chardata"`
}

type SAMLAttributeStatement struct {
	XMLName    xml.Name        `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeStatement"`
	Attributes []SAMLAttribute `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute"`
}

type SAMLAttribute struct {
	XMLName        xml.Name           `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute"`
	Name           string             `xml:"Name,attr"`
	AttributeValue SAMLAttributeValue `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue"`
}

type SAMLAttributeValue struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue"`
	Type    string   `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	Value   string   `xml:",chardata"`
}

type SAMLAuthnStatement struct {
	XMLName      xml.Name         `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnStatement"`
	AuthnInstant string           `xml:"AuthnInstant,attr"`
	SessionIndex string           `xml:"SessionIndex,attr"`
	AuthnContext SAMLAuthnContext `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContext"`
}

type SAMLAuthnContext struct {
	XMLName              xml.Name                 `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContext"`
	AuthnContextClassRef SAMLAuthnContextClassRef `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContextClassRef"`
}

type SAMLAuthnContextClassRef struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContextClassRef"`
	Value   string   `xml:",chardata"`
}

// SAMLInitHandler initiates SAML SSO for legacy applications
func SAMLInitHandler(c *gin.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get app configuration
	app, exists := services.GetSaaSApp(appID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if app.Protocol != constants.ProtocolSAML {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application does not support SAML"})
		return
	}

	// Generate SAML request
	requestID := generateSAMLID()
	issueInstant := time.Now().UTC().Format(time.RFC3339)

	samlRequest := SAMLRequest{
		ID:                          requestID,
		Version:                     "2.0",
		IssueInstant:                issueInstant,
		Destination:                 app.Config["sso_url"],
		ProtocolBinding:             "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
		AssertionConsumerServiceURL: fmt.Sprintf("http://localhost:8081/saml/%s/acs", appID),
		Issuer: SAMLIssuer{
			Value: "CloudGate-SSO",
		},
	}

	// Convert to XML
	xmlData, err := xml.MarshalIndent(samlRequest, "", "  ")
	if err != nil {
		log.Printf("Error marshaling SAML request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate SAML request"})
		return
	}

	// Store request for validation (in production, use Redis)
	// For demo, we'll skip this step

	// Create HTML form for auto-submission
	samlRequestB64 := encodeBase64(xmlData)
	relayState := generateSAMLID()

	htmlForm := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>CloudGate SAML SSO</title>
</head>
<body onload="document.forms[0].submit()">
    <form method="post" action="%s">
        <input type="hidden" name="SAMLRequest" value="%s" />
        <input type="hidden" name="RelayState" value="%s" />
        <noscript>
            <p>Your browser does not support JavaScript. Please click the button below to continue.</p>
            <input type="submit" value="Continue" />
        </noscript>
    </form>
    <p>Redirecting to %s...</p>
</body>
</html>`, app.Config["sso_url"], samlRequestB64, relayState, app.Name)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlForm)
}

// SAMLACSHandler handles SAML Assertion Consumer Service responses
func SAMLACSHandler(c *gin.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	// Get SAML response from form data
	samlResponse := c.PostForm("SAMLResponse")
	_ = c.PostForm("RelayState") // RelayState for future use

	if samlResponse == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SAML response is required"})
		return
	}

	// Decode base64 SAML response
	xmlData, err := decodeBase64(samlResponse)
	if err != nil {
		log.Printf("Error decoding SAML response: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SAML response"})
		return
	}

	// Parse SAML response
	var response SAMLResponse
	if err := xml.Unmarshal(xmlData, &response); err != nil {
		log.Printf("Error parsing SAML response: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SAML response format"})
		return
	}

	// Validate SAML response
	if response.Status.StatusCode.Value != "urn:oasis:names:tc:SAML:2.0:status:Success" {
		log.Printf("SAML authentication failed: %s", response.Status.StatusCode.Value)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "SAML authentication failed"})
		return
	}

	// Extract user information from assertion
	userEmail := response.Assertion.Subject.NameID.Value
	userID := constants.DemoUserID // In production, map from SAML attributes

	// Create or update app connection
	services.CreateUserAppConnection(userID, appID)
	err = services.UpdateUserAppConnection(userID, appID, map[string]interface{}{
		"status":       constants.StatusConnected,
		"user_email":   userEmail,
		"connected_at": time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		log.Printf("Error updating app connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connection"})
		return
	}

	// Log SAML authentication
	services.LogAuditEvent(userID, "saml_authentication", "app", appID, c.ClientIP(), c.GetHeader("User-Agent"), fmt.Sprintf("SAML authentication successful for %s", appID), "success")

	// Redirect to frontend with success
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := fmt.Sprintf("%s/oauth/callback?provider=%s&email=%s&code=success", frontendURL, appID, url.QueryEscape(userEmail))
	c.Redirect(http.StatusFound, redirectURL)
}

// SAMLMetadataHandler provides SAML metadata for CloudGate as IdP
func SAMLMetadataHandler(c *gin.Context) {
	// Get base URL from environment or use default
	baseURL := getEnv("BACKEND_URL", "http://localhost:8081")

	metadata := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     entityID="CloudGate-SSO">
    <md:IDPSSODescriptor WantAuthnRequestsSigned="false"
                         protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:KeyDescriptor use="signing">
            <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
                <ds:X509Data>
                    <ds:X509Certificate><!-- Certificate would go here --></ds:X509Certificate>
                </ds:X509Data>
            </ds:KeyInfo>
        </md:KeyDescriptor>
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
                               Location="%s/saml/sso"/>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                               Location="%s/saml/sso"/>
    </md:IDPSSODescriptor>
</md:EntityDescriptor>`, baseURL, baseURL)

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, metadata)
}

// Helper functions
func generateSAMLID() string {
	return "_" + uuid.New().String()
}

func encodeBase64(data []byte) string {
	return fmt.Sprintf("%s", data) // Simplified for demo - should use proper base64 encoding
}

func decodeBase64(data string) ([]byte, error) {
	return []byte(data), nil // Simplified for demo - should use proper base64 decoding
}
