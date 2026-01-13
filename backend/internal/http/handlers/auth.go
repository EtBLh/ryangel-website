package handlers

import (
	"encoding/json"
    "fmt"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/oauth2"
    "golang.org/x/oauth2/google"

	"github.com/ryangel/ryangel-backend/internal/config"
	httpmw "github.com/ryangel/ryangel-backend/internal/http/middleware"
	"github.com/ryangel/ryangel-backend/internal/models"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
)

// AuthHandler exposes admin and client authentication endpoints.
type AuthHandler struct {
	Service *authsvc.Service
	Config  *config.Config
}

// RegisterAdminRoutes wires admin auth endpoints beneath /admin.
func (h AuthHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	admin.POST("/login", h.adminLogin)

	protected := admin.Group("")
	protected.Use(httpmw.AdminAuth(h.Service))
	protected.GET("/me", h.adminMe)
	protected.POST("/logout", h.adminLogout)
}

// RegisterClientRoutes wires client auth endpoints beneath /clients.
func (h AuthHandler) RegisterClientRoutes(rg *gin.RouterGroup) {
	clients := rg.Group("/clients")
	clients.POST("/register", h.clientRegister)
	clients.POST("/login", h.clientLogin)
	clients.POST("/verify-otp", h.verifyOTP)
	clients.POST("/resend-otp", h.resendOTP)

	authenticated := clients.Group("")
	authenticated.Use(httpmw.ClientAuth(h.Service))
	authenticated.GET("/me", h.clientMe)
	authenticated.PATCH("/me", h.clientUpdate)
	authenticated.POST("/logout", h.clientLogout)
}

type loginRequest struct {
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
	CartID   string `json:"cart_id"`
}

type adminLoginRequest struct {
	Username string `json:"username" binding:"required_without=Email"`
	Email    string `json:"email" binding:"required_without=Username"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Phone    string  `json:"phone" binding:"required"`
	Email    *string `json:"email"`
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
}

type updateClientRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	DateOfBirth *string `json:"date_of_birth"` // YYYY-MM-DDT00:00:00Z
}

type verifyOTPRequest struct {
	Phone  string `json:"phone" binding:"required"`
	OTP    string `json:"otp" binding:"required,len=6"`
	CartID string `json:"cart_id"`
}

func (r adminLoginRequest) identifier() string {
	if trimmed := strings.TrimSpace(r.Username); trimmed != "" {
		return trimmed
	}
	return strings.TrimSpace(r.Email)
}

func (h AuthHandler) adminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	identifier := req.identifier()
	if identifier == "" {
		writeValidationError(c, errors.New("username or email is required"))
		return
	}

	result, err := h.Service.AdminLogin(c.Request.Context(), identifier, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrInvalidCredentials):
			writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Invalid username/email or password.", nil)
		case errors.Is(err, authsvc.ErrInactiveAccount):
			writeError(c, http.StatusForbidden, "AUTH_INACTIVE_ACCOUNT", "Account is inactive.", nil)
		default:
			writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to login right now.", nil)
		}
		return
	}

	response := gin.H{
		"token":      result.Token,
		"token_type": "Bearer",
		"expires_in": int(h.Config.TokenTTL().Seconds()),
		"admin":      toAdminPayload(result.Admin),
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) adminMe(c *gin.Context) {
	admin, ok := httpmw.AdminFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Missing authentication context.", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"admin": toAdminPayload(admin)})
}

func (h AuthHandler) adminLogout(c *gin.Context) {
	admin, ok := httpmw.AdminFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Missing authentication context.", nil)
		return
	}

	if err := h.Service.AdminLogout(c.Request.Context(), admin.ID); err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to logout right now.", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h AuthHandler) clientRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	client, err := h.Service.ClientRegister(c.Request.Context(), req.Phone, req.Email, req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") || strings.Contains(err.Error(), "username already taken") {
			writeError(c, http.StatusConflict, "CLIENT_EXISTS", err.Error(), nil)
			return
		}
		writeError(c, http.StatusInternalServerError, "REGISTRATION_ERROR", "Unable to register right now.", nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "OTP sent to " + *client.Phone,
		"otp_expires_in": 300,
	})
}

func (h AuthHandler) clientLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	if req.Password != "" {
		identifier := req.Phone
		if identifier == "" {
			identifier = req.Username
		}
		if identifier == "" {
			writeValidationError(c, errors.New("phone or username required for password login"))
			return
		}

		result, err := h.Service.ClientLoginPassword(c.Request.Context(), identifier, req.Password, req.CartID)
		if err != nil {
			writeError(c, http.StatusUnauthorized, "LOGIN_FAILED", "Invalid credentials", nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":      result.Token,
			"expires_in": int(h.Config.TokenTTL().Seconds()),
			"client":     result.Client,
            "cart_id":    result.CartID,
		})
		return
	}

	err := h.Service.ClientLogin(c.Request.Context(), req.Phone)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrPhoneNotFound):
			writeError(c, http.StatusNotFound, "PHONE_NOT_FOUND", "Phone number not registered.", nil)
		case errors.Is(err, authsvc.ErrInactiveAccount):
			writeError(c, http.StatusForbidden, "AUTH_INACTIVE_ACCOUNT", "Account is inactive.", nil)
		default:
			writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to send OTP right now.", nil)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "OTP sent to " + req.Phone,
		"otp_expires_in": 300,
	})
}

func (h AuthHandler) verifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	result, err := h.Service.VerifyOTP(c.Request.Context(), req.Phone, req.OTP, req.CartID)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrInvalidOTP):
			writeError(c, http.StatusUnauthorized, "INVALID_OTP", "Invalid or expired OTP.", nil)
		case errors.Is(err, authsvc.ErrInactiveAccount):
			writeError(c, http.StatusForbidden, "AUTH_INACTIVE_ACCOUNT", "Account is inactive.", nil)
		default:
			writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to verify OTP right now.", nil)
		}
		return
	}

	response := gin.H{
		"token":      result.Token,
		"token_type": "Bearer",
		"expires_in": int(h.Config.TokenTTL().Seconds()),
		"client":     toClientPayload(result.Client),
        "cart_id":    result.CartID,
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) resendOTP(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	err := h.Service.ResendOTP(c.Request.Context(), req.Phone)
	if err != nil {
		switch {
		case errors.Is(err, authsvc.ErrPhoneNotFound):
			writeError(c, http.StatusNotFound, "PHONE_NOT_FOUND", "Phone number not registered.", nil)
		case errors.Is(err, authsvc.ErrInactiveAccount):
			writeError(c, http.StatusForbidden, "AUTH_INACTIVE_ACCOUNT", "Account is inactive.", nil)
		default:
			writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to resend OTP right now.", nil)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "OTP resent to " + req.Phone,
		"otp_expires_in": 300,
	})
}

func (h AuthHandler) clientMe(c *gin.Context) {
	client, ok := httpmw.ClientFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Missing authentication context.", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": toClientPayload(client)})
}

func (h AuthHandler) clientUpdate(c *gin.Context) {
	client, ok := httpmw.ClientFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Missing authentication context.", nil)
		return
	}

	var req updateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	var parsedDOB *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		// Handle formats: "2006-01-02T15:04:05.000Z" (ISO) or "2006-01-02"
		t, err := time.Parse(time.RFC3339, *req.DateOfBirth)
		if err != nil {
			// Try simple date
			t, err = time.Parse("2006-01-02", *req.DateOfBirth)
			if err != nil {
				writeError(c, http.StatusBadRequest, "INVALID_DATE", "Invalid date format. Use ISO8601 or YYYY-MM-DD", nil)
				return
			}
		}
		parsedDOB = &t
	}

	updatedClient, err := h.Service.UpdateClient(c.Request.Context(), client.ID, req.Email, req.Username, parsedDOB)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "username") || strings.Contains(pgErr.Detail, "username") {
				writeError(c, http.StatusConflict, "USERNAME_EXISTS", "This username is already taken.", nil)
				return
			}
			if strings.Contains(pgErr.ConstraintName, "email") || strings.Contains(pgErr.Detail, "email") {
				writeError(c, http.StatusConflict, "EMAIL_EXISTS", "This email is already taken.", nil)
				return
			}
			writeError(c, http.StatusConflict, "DUPLICATE_ENTRY", "This value is already in use.", nil)
			return
		}

		writeError(c, http.StatusInternalServerError, "UPDATE_ERROR", "Unable to update client right now.", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": toClientPayload(updatedClient)})
}

func (h AuthHandler) clientLogout(c *gin.Context) {
	client, ok := httpmw.ClientFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Missing authentication context.", nil)
		return
	}

	if err := h.Service.ClientLogout(c.Request.Context(), client.ID); err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_INTERNAL_ERROR", "Unable to logout right now.", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func toAdminPayload(admin *models.Admin) gin.H {
	return gin.H{
		"admin_id":   admin.ID,
		"username":   admin.Username,
		"email":      admin.Email,
		"is_active":  admin.IsActive,
		"last_login": admin.LastLogin,
	}
}

func toClientPayload(client *models.Client) gin.H {
	return gin.H{
		"client_id":     client.ID,
		"email":         client.Email,
		"username":      client.Username,
		"phone":         client.Phone,
		"is_active":     client.IsActive,
		"date_of_birth": client.DateOfBirth,
	}
}

func writeValidationError(c *gin.Context, err error) {
	writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request payload.", gin.H{"error": err.Error()})
}

func writeError(c *gin.Context, status int, code, message string, details gin.H) {
	payload := gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}

	if len(details) > 0 {
		payload["error"].(gin.H)["details"] = details
	}

	c.JSON(status, payload)
}

func (h AuthHandler) getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.Config.GoogleClientID,
		ClientSecret: h.Config.GoogleClientSecret,
		RedirectURL:  h.Config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (h AuthHandler) RegisterGoogleRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	google := auth.Group("/google")
	google.GET("/login", h.googleLogin)
	google.GET("/callback", h.googleCallback)
}

func (h AuthHandler) googleLogin(c *gin.Context) {
	conf := h.getOAuthConfig()
	url := conf.AuthCodeURL("state-csrf-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h AuthHandler) googleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "Code not found", nil)
		return
	}

	conf := h.getOAuthConfig()
	token, err := conf.Exchange(c.Request.Context(), code)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_ERROR", "Failed to exchange token", nil)
		return
	}

	client := conf.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_ERROR", "Failed to get user info", nil)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_ERROR", "Failed to parse user info", nil)
		return
	}

	res, err := h.Service.LoginWithGoogle(c.Request.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "AUTH_ERROR", err.Error(), nil)
		return
	}

	// Redirect to frontend with token
	frontendURL := "http://localhost:5173"
	redirectURL := fmt.Sprintf("%s/google-callback?token=%s&cart_id=%s", frontendURL, res.Token, res.CartID)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
