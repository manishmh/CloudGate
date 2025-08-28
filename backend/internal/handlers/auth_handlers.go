package handlers

import (
	"net/http"
	"os"
	"time"

	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RegisterHandler registers a new local user
func RegisterHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := userService.CreateUserWithPassword(req.Email, req.Username, req.FirstName, req.LastName, string(hash))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user_id": user.ID})
	}
}

// LoginHandler authenticates a user and returns tokens
func LoginHandler(userService *services.UserService, sessionService *services.SessionService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Create a session (used as refresh token)
		session, err := sessionService.CreateSession(user.ID, c.ClientIP(), c.GetHeader("User-Agent"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		// Issue access token
		accessToken, expiresIn, err := generateAccessToken(cfg, user.ID.String(), user.Email, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue token"})
			return
		}

		// Configure cookie attributes for deployment (e.g., Render)
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		if cookieSecure {
			c.SetSameSite(http.SameSiteNoneMode)
		} else {
			c.SetSameSite(http.SameSiteLaxMode)
		}

		// Set tokens as HTTP-only cookies for browser auth
		c.SetCookie("access_token", accessToken, expiresIn, "/", cookieDomain, cookieSecure, true)
		c.SetCookie("refresh_token", session.SessionToken, cfg.RefreshTokenTTLHour*3600, "/", cookieDomain, cookieSecure, true)

		c.JSON(http.StatusOK, tokenResponse{
			AccessToken:  accessToken,
			RefreshToken: session.SessionToken,
			ExpiresIn:    expiresIn,
			TokenType:    "Bearer",
		})
	}
}

// RefreshHandler exchanges a refresh token (session token) for a new access token
func RefreshHandler(sessionService *services.SessionService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session, err := sessionService.GetSessionByToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		// Rotate/refresh session expiry
		if _, err := sessionService.RefreshSession(req.RefreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
			return
		}

		accessToken, expiresIn, err := generateAccessToken(cfg, session.User.ID.String(), session.User.Email, session.User.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue token"})
			return
		}
		// Update access token cookie
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		if cookieSecure {
			c.SetSameSite(http.SameSiteNoneMode)
		} else {
			c.SetSameSite(http.SameSiteLaxMode)
		}
		c.SetCookie("access_token", accessToken, expiresIn, "/", cookieDomain, cookieSecure, true)

		c.JSON(http.StatusOK, tokenResponse{
			AccessToken:  accessToken,
			RefreshToken: req.RefreshToken,
			ExpiresIn:    expiresIn,
			TokenType:    "Bearer",
		})
	}
}

// LogoutHandler invalidates a refresh token
func LogoutHandler(sessionService *services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := sessionService.InvalidateSession(req.RefreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}
		// Clear cookies
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		if cookieSecure {
			c.SetSameSite(http.SameSiteNoneMode)
		} else {
			c.SetSameSite(http.SameSiteLaxMode)
		}
		c.SetCookie("access_token", "", -1, "/", cookieDomain, cookieSecure, true)
		c.SetCookie("refresh_token", "", -1, "/", cookieDomain, cookieSecure, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}

func generateAccessToken(cfg *config.Config, sub, email, username string) (string, int, error) {
	ttl := time.Duration(cfg.AccessTokenTTLMin) * time.Minute
	expiresAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"sub":      sub,
		"email":    email,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"typ":      "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", 0, err
	}
	return signed, int(ttl.Seconds()), nil
}
