package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
	"cloudgate-backend/pkg/types"
)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c *gin.Context) {
	response := types.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "cloudgate-backend",
	}
	c.JSON(http.StatusOK, response)
}

// PrivacyPolicyHandler serves the privacy policy
func PrivacyPolicyHandler(c *gin.Context) {
	privacyPolicy := `
<h1>CloudGate Privacy Policy</h1>

<p><strong>Effective Date:</strong> January 2025<br>
<strong>Last Updated:</strong> January 2025</p>

<h2>1. Introduction</h2>

<p>CloudGate ("we," "our," or "us") is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our Single Sign-On (SSO) service.</p>

<h2>2. Information We Collect</h2>

<h3>2.1 Personal Information</h3>
<ul>
<li><strong>Account Information:</strong> Name, email address, username</li>
<li><strong>Profile Information:</strong> Profile picture, preferences, settings</li>
<li><strong>Authentication Data:</strong> Login credentials, session tokens</li>
<li><strong>Contact Information:</strong> Email for verification and communication</li>
</ul>

<h3>2.2 Technical Information</h3>
<ul>
<li><strong>Log Data:</strong> IP addresses, browser type, device information</li>
<li><strong>Usage Data:</strong> Login times, accessed applications, session duration</li>
<li><strong>Security Data:</strong> Failed login attempts, security events</li>
</ul>

<h3>2.3 Third-Party Integration Data</h3>
<ul>
<li><strong>OAuth Tokens:</strong> Access tokens for connected SaaS applications</li>
<li><strong>Application Data:</strong> Connection status, usage patterns</li>
</ul>

<h2>3. How We Use Your Information</h2>

<p>We use your information to:</p>
<ul>
<li>Provide and maintain our SSO service</li>
<li>Authenticate and authorize access to applications</li>
<li>Improve security and prevent fraud</li>
<li>Send important notifications and updates</li>
<li>Provide customer support</li>
<li>Comply with legal obligations</li>
</ul>

<h2>4. Information Sharing and Disclosure</h2>

<p>We do not sell, trade, or rent your personal information. We may share information:</p>
<ul>
<li><strong>With Your Consent:</strong> When you explicitly authorize sharing</li>
<li><strong>Service Providers:</strong> Third-party vendors who assist our operations</li>
<li><strong>Legal Requirements:</strong> When required by law or to protect rights</li>
<li><strong>Business Transfers:</strong> In case of merger, acquisition, or sale</li>
</ul>

<h2>5. Data Security</h2>

<p>We implement appropriate security measures including:</p>
<ul>
<li>Encryption of data in transit and at rest</li>
<li>Regular security audits and monitoring</li>
<li>Access controls and authentication</li>
<li>Secure data centers and infrastructure</li>
</ul>

<h2>6. Data Retention</h2>

<p>We retain your information:</p>
<ul>
<li><strong>Account Data:</strong> Until account deletion</li>
<li><strong>Log Data:</strong> For 90 days for security purposes</li>
<li><strong>Audit Logs:</strong> For 7 years for compliance</li>
<li><strong>Session Data:</strong> Until session expiration</li>
</ul>

<h2>7. Your Rights</h2>

<p>You have the right to:</p>
<ul>
<li>Access your personal information</li>
<li>Correct inaccurate information</li>
<li>Delete your account and data</li>
<li>Export your data</li>
<li>Opt-out of non-essential communications</li>
</ul>

<h2>8. Cookies and Tracking</h2>

<p>We use cookies and similar technologies for:</p>
<ul>
<li>Authentication and session management</li>
<li>Security and fraud prevention</li>
<li>Analytics and performance monitoring</li>
<li>User preferences and settings</li>
</ul>

<h2>9. Third-Party Services</h2>

<p>Our service integrates with third-party applications. Each has their own privacy policies:</p>
<ul>
<li>Google Workspace</li>
<li>Microsoft 365</li>
<li>Slack</li>
<li>Salesforce</li>
<li>Other connected applications</li>
</ul>

<h2>10. International Data Transfers</h2>

<p>Your information may be transferred to and processed in countries other than your own. We ensure appropriate safeguards are in place.</p>

<h2>11. Children's Privacy</h2>

<p>Our service is not intended for children under 13. We do not knowingly collect information from children under 13.</p>

<h2>12. Changes to This Policy</h2>

<p>We may update this Privacy Policy. We will notify you of significant changes via email or service notifications.</p>

<h2>13. Contact Information</h2>

<p><strong>Data Controller:</strong> Manish Kumar Saw<br>
<strong>Email:</strong> manishmh982@gmail.com<br>
<strong>Address:</strong> [Your Business Address]</p>

<p>For privacy-related questions or requests, please contact us at the above email address.</p>

<h2>14. Legal Basis for Processing (GDPR)</h2>

<p>For EU users, our legal basis for processing includes:</p>
<ul>
<li><strong>Consent:</strong> For optional features and communications</li>
<li><strong>Contract:</strong> To provide our SSO service</li>
<li><strong>Legitimate Interest:</strong> For security and service improvement</li>
<li><strong>Legal Obligation:</strong> For compliance requirements</li>
</ul>

<h2>15. Data Protection Officer</h2>

<p>For GDPR-related inquiries, contact our Data Protection Officer at: manishmh982@gmail.com</p>

<hr>

<p><em>This Privacy Policy is part of our Terms of Service and governs your use of CloudGate.</em></p>
`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, privacyPolicy)
}

// TermsHandler serves the terms and conditions
func TermsHandler(c *gin.Context) {
	terms := `
<h1>CloudGate Terms of Service</h1>

<p><strong>Effective Date:</strong> January 2025<br>
<strong>Last Updated:</strong> January 2025</p>

<h2>1. Acceptance of Terms</h2>

<p>By accessing or using CloudGate ("Service"), you agree to be bound by these Terms of Service ("Terms"). If you do not agree to these Terms, do not use the Service.</p>

<h2>2. Description of Service</h2>

<p>CloudGate is a Single Sign-On (SSO) platform that provides:</p>
<ul>
<li>Centralized authentication for multiple applications</li>
<li>User identity management</li>
<li>Security and access controls</li>
<li>Integration with third-party SaaS applications</li>
</ul>

<h2>3. User Accounts</h2>

<h3>3.1 Account Creation</h3>
<ul>
<li>You must provide accurate and complete information</li>
<li>You are responsible for maintaining account security</li>
<li>You must be at least 13 years old to use the Service</li>
<li>One person may not maintain multiple accounts</li>
</ul>

<h3>3.2 Account Security</h3>
<ul>
<li>Keep your login credentials confidential</li>
<li>Notify us immediately of any unauthorized access</li>
<li>You are responsible for all activities under your account</li>
<li>Use strong passwords and enable available security features</li>
</ul>

<h2>4. Acceptable Use</h2>

<h3>4.1 Permitted Uses</h3>
<ul>
<li>Access authorized applications through SSO</li>
<li>Manage your profile and preferences</li>
<li>Use security features as intended</li>
<li>Integrate with supported third-party applications</li>
</ul>

<h3>4.2 Prohibited Uses</h3>
<ul>
<li>Violate any laws or regulations</li>
<li>Infringe on intellectual property rights</li>
<li>Attempt to gain unauthorized access</li>
<li>Distribute malware or harmful code</li>
<li>Interfere with service operation</li>
<li>Use the service for illegal activities</li>
</ul>

<h2>5. Privacy and Data Protection</h2>

<p>Your privacy is important to us. Our Privacy Policy explains how we collect, use, and protect your information. By using the Service, you consent to our data practices as described in the Privacy Policy.</p>

<h2>6. Third-Party Integrations</h2>

<h3>6.1 Connected Applications</h3>
<ul>
<li>You may connect third-party applications through OAuth</li>
<li>Each application has its own terms and privacy policy</li>
<li>We are not responsible for third-party application behavior</li>
<li>You grant us permission to facilitate these connections</li>
</ul>

<h3>6.2 Data Sharing</h3>
<ul>
<li>We share minimal necessary data with connected applications</li>
<li>You control which applications to connect</li>
<li>You can revoke access at any time</li>
<li>Review each application's data access requirements</li>
</ul>

<h2>7. Service Availability</h2>

<h3>7.1 Uptime</h3>
<ul>
<li>We strive for high availability but cannot guarantee 100% uptime</li>
<li>Scheduled maintenance will be announced in advance</li>
<li>We are not liable for service interruptions</li>
</ul>

<h3>7.2 Support</h3>
<ul>
<li>Support is provided via email: manishmh982@gmail.com</li>
<li>We aim to respond within 24-48 hours</li>
<li>Premium support may be available for enterprise customers</li>
</ul>

<h2>8. Intellectual Property</h2>

<h3>8.1 Our Rights</h3>
<ul>
<li>CloudGate and related trademarks are our property</li>
<li>The Service and its technology are protected by intellectual property laws</li>
<li>You may not copy, modify, or reverse engineer our Service</li>
</ul>

<h3>8.2 Your Rights</h3>
<ul>
<li>You retain ownership of your data and content</li>
<li>You grant us license to use your data to provide the Service</li>
<li>You can export your data at any time</li>
</ul>

<h2>9. Fees and Payment</h2>

<h3>9.1 Free Tier</h3>
<ul>
<li>Basic SSO functionality is provided free of charge</li>
<li>Usage limits may apply to free accounts</li>
<li>We reserve the right to modify free tier limitations</li>
</ul>

<h3>9.2 Paid Plans</h3>
<ul>
<li>Premium features may require payment</li>
<li>Fees are charged in advance</li>
<li>Refunds are provided according to our refund policy</li>
</ul>

<h2>10. Termination</h2>

<h3>10.1 By You</h3>
<ul>
<li>You may terminate your account at any time</li>
<li>Data deletion will occur according to our retention policy</li>
<li>Some data may be retained for legal compliance</li>
</ul>

<h3>10.2 By Us</h3>
<ul>
<li>We may terminate accounts for Terms violations</li>
<li>We may suspend service for security reasons</li>
<li>We will provide notice when reasonably possible</li>
</ul>

<h2>11. Disclaimers</h2>

<h3>11.1 Service Warranty</h3>
<ul>
<li>The Service is provided "as is" without warranties</li>
<li>We disclaim all warranties, express or implied</li>
<li>We do not guarantee error-free operation</li>
</ul>

<h3>11.2 Third-Party Services</h3>
<ul>
<li>We are not responsible for third-party application failures</li>
<li>Connected services have their own terms and limitations</li>
<li>Integration issues may occur beyond our control</li>
</ul>

<h2>12. Limitation of Liability</h2>

<p>To the maximum extent permitted by law:</p>
<ul>
<li>Our liability is limited to the amount you paid for the Service</li>
<li>We are not liable for indirect, incidental, or consequential damages</li>
<li>Some jurisdictions do not allow liability limitations</li>
</ul>

<h2>13. Indemnification</h2>

<p>You agree to indemnify and hold us harmless from claims arising from:</p>
<ul>
<li>Your use of the Service</li>
<li>Your violation of these Terms</li>
<li>Your violation of third-party rights</li>
<li>Your negligent or wrongful conduct</li>
</ul>

<h2>14. Governing Law</h2>

<p>These Terms are governed by the laws of [Your Jurisdiction]. Disputes will be resolved in the courts of [Your Jurisdiction].</p>

<h2>15. Changes to Terms</h2>

<h3>15.1 Modifications</h3>
<ul>
<li>We may modify these Terms at any time</li>
<li>Significant changes will be communicated via email</li>
<li>Continued use constitutes acceptance of new Terms</li>
</ul>

<h3>15.2 Notice Period</h3>
<ul>
<li>Changes take effect 30 days after notification</li>
<li>You may terminate your account if you disagree with changes</li>
</ul>

<h2>16. Contact Information</h2>

<p><strong>Service Provider:</strong> Manish Kumar Saw<br>
<strong>Email:</strong> manishmh982@gmail.com<br>
<strong>Address:</strong> [Your Business Address]</p>

<p>For questions about these Terms, contact us at the above email address.</p>

<h2>17. Severability</h2>

<p>If any provision of these Terms is found unenforceable, the remaining provisions will continue in full force and effect.</p>

<h2>18. Entire Agreement</h2>

<p>These Terms, together with our Privacy Policy, constitute the entire agreement between you and CloudGate regarding the Service.</p>

<hr>

<p><em>By using CloudGate, you acknowledge that you have read, understood, and agree to be bound by these Terms of Service.</em></p>
`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, terms)
}

// APIInfoHandler provides information about the API endpoints
func APIInfoHandler(c *gin.Context) {
	response := types.APIInfoResponse{
		Service:     "CloudGate SSO Backend",
		Version:     "1.0.0",
		Description: "Enterprise SSO Portal Backend API",
		Endpoints: []string{
			"GET /health - Health check",
			"GET /privacy-policy - Privacy Policy",
			"GET /terms - Terms of Service",
			"POST /token/introspect - Token introspection",
			"GET /user/info - User information",
			"GET /api/info - API information",
			"GET /apps - List SaaS applications",
			"POST /apps/connect - Connect to a SaaS application",
			"POST /apps/launch - Launch a SaaS application",
			"POST /apps/callback - OAuth callback handler",
		},
	}
	c.JSON(http.StatusOK, response)
}

// TokenIntrospectionHandler handles JWT token introspection requests
func TokenIntrospectionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.TokenIntrospectionRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Prepare introspection request to Keycloak
		introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
			cfg.KeycloakURL, cfg.KeycloakRealm)

		data := url.Values{}
		data.Set("token", request.Token)
		data.Set("client_id", cfg.KeycloakClientID)

		req, err := http.NewRequest("POST", introspectionURL, strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("Error creating introspection request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making introspection request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to introspect token"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading introspection response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var introspectionResp types.TokenIntrospectionResponse
		if err := json.Unmarshal(body, &introspectionResp); err != nil {
			log.Printf("Error parsing introspection response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, introspectionResp)
	}
}

// UserInfoHandler handles user information requests
func UserInfoHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := tokenParts[1]

		// Get user info from Keycloak
		userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
			cfg.KeycloakURL, cfg.KeycloakRealm)

		req, err := http.NewRequest("GET", userInfoURL, nil)
		if err != nil {
			log.Printf("Error creating userinfo request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making userinfo request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading userinfo response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var userInfo map[string]interface{}
		if err := json.Unmarshal(body, &userInfo); err != nil {
			log.Printf("Error parsing userinfo response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}

// GetAppsHandler returns all SaaS applications with user connection status
func GetAppsHandler(c *gin.Context) {
	// Get user ID from token (simplified for demo)
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	apps := services.GetAppsWithUserStatus(userID)
	c.JSON(http.StatusOK, gin.H{
		"apps":  apps,
		"count": len(apps),
	})
}

// ConnectAppHandler initiates OAuth connection to a SaaS application
func ConnectAppHandler(c *gin.Context) {
	var request types.AppConnectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get app configuration
	app, exists := services.GetSaaSApp(request.AppID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Create or update user connection
	services.CreateUserAppConnection(userID, request.AppID)

	// Generate OAuth URL
	state := services.GenerateState()
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s",
		app.Config["auth_url"],
		app.Config["client_id"],
		url.QueryEscape("http://localhost:8081/apps/callback"),
		url.QueryEscape(app.Config["scope"]),
		state,
	)

	// Store state for validation (in production, use Redis or database)
	// For demo, we'll skip state validation

	response := types.AppConnectionResponse{
		AuthURL: authURL,
		State:   state,
	}

	c.JSON(http.StatusOK, response)
}

// LaunchAppHandler handles application launch requests
func LaunchAppHandler(c *gin.Context) {
	var request types.AppLaunchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get app configuration
	app, exists := services.GetSaaSApp(request.AppID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check if user is connected to the app
	connection, exists := services.GetUserAppConnection(userID, request.AppID)
	if !exists || connection.Status != constants.StatusConnected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not connected to this application"})
		return
	}

	// Get launch URL from constants
	launchURL, exists := constants.LaunchURLs[request.AppID]
	if !exists {
		launchURL = app.LaunchURL
	}

	// If still no launch URL, use a default
	if launchURL == "" {
		launchURL = "https://example.com"
	}

	// Update last access time
	services.UpdateUserAppConnection(userID, request.AppID, map[string]interface{}{
		"last_access_at": time.Now().UTC().Format(time.RFC3339),
	})

	response := types.AppLaunchResponse{
		LaunchURL: launchURL,
		Method:    "redirect",
		ExpiresIn: 3600,
	}

	c.JSON(http.StatusOK, response)
}

// OAuthCallbackHandler handles OAuth callbacks from SaaS applications
func OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	appID := c.Query("app_id")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})
		return
	}

	// For demo purposes, we'll simulate successful OAuth completion
	// In production, you would exchange the code for tokens

	// Simulate finding the user (in production, you'd validate the state)
	userID := constants.DemoUserID // This would come from the state parameter

	// Update connection status
	err := services.UpdateUserAppConnection(userID, appID, map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": constants.DemoAccessToken,
		"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connection"})
		return
	}

	// Redirect back to frontend
	c.Redirect(http.StatusFound, "http://localhost:3000/dashboard?connected="+appID)
}

// Helper function to extract user ID from context
// In production, this would parse the JWT token
func getUserIDFromContext(c *gin.Context) string {
	// For demo purposes, return a fixed user ID
	// In production, you would extract this from the JWT token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Simplified: just return demo user ID if auth header exists
	return constants.DemoUserID
}

// DatabaseHealthCheckHandler checks database connectivity
func DatabaseHealthCheckHandler(c *gin.Context) {
	if err := services.DatabaseHealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}

// AdminStatsHandler returns system statistics (placeholder)
func AdminStatsHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	sessionService := services.NewSessionService(services.GetDB())
	stats, err := sessionService.GetSessionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// AdminUsersHandler returns user list (placeholder)
func AdminUsersHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin users endpoint - not implemented yet",
	})
}

// AdminSessionsHandler returns session list (placeholder)
func AdminSessionsHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin sessions endpoint - not implemented yet",
	})
}
