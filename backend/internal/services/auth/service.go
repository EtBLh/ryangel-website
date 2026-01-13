package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/ryangel/ryangel-backend/internal/auth"
	"github.com/ryangel/ryangel-backend/internal/config"
	"github.com/ryangel/ryangel-backend/internal/models"
	"github.com/ryangel/ryangel-backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrOTPNotFound        = errors.New("OTP not found")
	ErrPhoneNotFound      = errors.New("phone number not registered")
)

type AdminLoginResult struct {
	Token     string
	ExpiresAt time.Time
	Admin     *models.Admin
}

type ClientLoginResult struct {
	Token     string
	ExpiresAt time.Time
	Client    *models.Client
	CartID    string
}

type Service struct {
	admins  *repository.AdminRepository
	clients *repository.ClientRepository
	carts   *repository.CartRepository
	cfg     *config.Config
	twilio  *twilio.RestClient
}

func NewService(admins *repository.AdminRepository, clients *repository.ClientRepository, carts *repository.CartRepository, cfg *config.Config) *Service {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &Service{
		admins:  admins,
		clients: clients,
		carts:   carts,
		cfg:     cfg,
		twilio:  client,
	}
}

func (s *Service) AdminLogin(ctx context.Context, identifier, password string) (*AdminLoginResult, error) {
	admin, err := s.admins.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !admin.IsActive {
		return nil, ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	raw, hash, err := auth.GenerateOpaqueToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(s.cfg.TokenTTL())
	if err := s.admins.UpdateToken(ctx, admin.ID, hash, expiresAt); err != nil {
		return nil, err
	}

	admin.TokenHash = &hash
	admin.TokenExpiry = &expiresAt
	admin.LastLogin = &now

	return &AdminLoginResult{Token: raw, ExpiresAt: expiresAt, Admin: admin}, nil
}

func (s *Service) AdminLogout(ctx context.Context, adminID int64) error {
	return s.admins.ClearToken(ctx, adminID)
}

func (s *Service) ValidateAdminToken(ctx context.Context, token string) (*models.Admin, error) {
	hash := auth.HashToken(token)
	admin, err := s.admins.GetByTokenHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return admin, nil
}

func (s *Service) ClientLogin(ctx context.Context, phone string) error {
	// Check if client exists
	client, err := s.clients.GetByPhone(ctx, phone)
	if err != nil && err != repository.ErrNotFound {
		return err
	}

	// If client doesn't exist, this is fine - we'll create one during registration
	// But for login, we require the client to exist
	if err == repository.ErrNotFound {
		return ErrPhoneNotFound
	}

	if !client.IsActive {
		return ErrInactiveAccount
	}

	// Generate 6-digit OTP
	otpCode := s.generateOTP()

	// Set OTP expiry (5 minutes from now)
	expiry := time.Now().Add(5 * time.Minute)

	// Store OTP in database
	if err := s.clients.SetOTP(ctx, phone, otpCode, expiry); err != nil {
		return err
	}

	// Send SMS via Twilio
	return s.sendSMS(phone, fmt.Sprintf("你的RyAngel驗證碼為: %s", otpCode))
}

func (s *Service) ClientLoginPassword(ctx context.Context, identifier string, password string, cartID string) (*ClientLoginResult, error) {
	// Try to find by phone first
	client, err := s.clients.GetByPhone(ctx, identifier)
	if err == repository.ErrNotFound {
		// Then by username
		client, err = s.clients.GetByUsername(ctx, identifier)
	}
	
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !client.IsActive {
		return nil, ErrInactiveAccount
	}

	if client.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*client.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if cartID != "" {
		if err := s.carts.AssignCartToClient(ctx, cartID, client.ID); err != nil {
			fmt.Printf("Failed to assign cart %s to client %d: %v\n", cartID, client.ID, err)
			cartID = "" // Reset logic to fetch/create
		}
	}

	return s.generateClientToken(ctx, client, cartID)
}

func (s *Service) ClientLogout(ctx context.Context, clientID int64) error {
	return s.clients.ClearToken(ctx, clientID)
}

func (s *Service) ValidateClientToken(ctx context.Context, token string) (*models.Client, error) {
	hash := auth.HashToken(token)
	client, err := s.clients.GetByTokenHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return client, nil
}

func (s *Service) ClientRegister(ctx context.Context, phone string, email *string, username string, password string) (*models.Client, error) {
	// Check if phone already exists
	existingClient, err := s.clients.GetByPhone(ctx, phone)
	if err == nil {
		// Phone exists. Check if active.
		if existingClient.IsActive {
			return nil, errors.New("phone number already registered")
		}
		// Inactive: assume retry registration. Update logic below.
	} else if err != repository.ErrNotFound {
		return nil, err
	} else {
		existingClient = nil
	}

	// Check if username already exists (and not owned by this inactive client)
	u, err := s.clients.GetByUsername(ctx, username)
	if err == nil {
		if existingClient == nil || u.ID != existingClient.ID {
			return nil, errors.New("username already taken")
		}
	} else if err != repository.ErrNotFound {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hashed)
	uName := &username

	var client *models.Client

	if existingClient != nil {
		// Update existing legacy/inactive client with new credentials
		client, err = s.clients.UpdateClientRegistration(ctx, existingClient.ID, uName, &hashStr)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new inactive client
		// Note: passing false for isActive
		client, err = s.clients.CreateClient(ctx, &phone, email, uName, &hashStr, nil, false)
		if err != nil {
			return nil, err
		}
	}

	// Generate & Send OTP
	otpCode := s.generateOTP()
	expiry := time.Now().Add(5 * time.Minute)

	if err := s.clients.SetOTP(ctx, *client.Phone, otpCode, expiry); err != nil {
		// If we failed to set OTP, we should probably rollback or error out.
		// For now return error.
		return nil, err
	}

	// Send SMS
	if err := s.sendSMS(*client.Phone, fmt.Sprintf("你的RyAngel驗證碼為: %s", otpCode)); err != nil {
		fmt.Printf("Failed to send SMS to %s: %v\n", *client.Phone, err)
		// Return success anyway so user can try "Resend OTP" or similar?
		// Or return error?
		// Better to return success but client won't receive it. They can click "Resend".
	}

	return client, nil
}

func (s *Service) generateClientToken(ctx context.Context, client *models.Client, cartID string) (*ClientLoginResult, error) {
	if !client.IsActive {
		return nil, ErrInactiveAccount
	}

	raw, hash, err := auth.GenerateOpaqueToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.cfg.TokenTTL())
	if err := s.clients.UpdateToken(ctx, client.ID, hash, expiresAt); err != nil {
		return nil, err
	}

	client.TokenHash = &hash
	client.TokenExpiry = &expiresAt

    // Handle Cart ID Logic
    finalCartID := cartID
    if finalCartID == "" {
        // Look up latest cart for client
        cart, err := s.carts.GetCartByClientID(ctx, client.ID)
        if err == nil {
            finalCartID = cart.CartID
        } else if err == repository.ErrNotFound {
            // Create new cart
            newCart, err := s.carts.CreateCart(ctx, &client.ID)
            if err != nil {
                // Log and proceed? Or fail? 
                // Failing login because cart creation failed might be too harsh, but cart is essential.
                // Let's log and return empty cartID? No, better return error or handle it.
                // Assuming CreateCart won't fail often.
                fmt.Printf("Failed to create cart for client %d: %v\n", client.ID, err)
            } else {
                finalCartID = newCart.CartID
            }
        } else {
             fmt.Printf("Failed to get cart for client %d: %v\n", client.ID, err)
        }
    }

	return &ClientLoginResult{Token: raw, ExpiresAt: expiresAt, Client: client, CartID: finalCartID}, nil
}

func (s *Service) VerifyOTP(ctx context.Context, phone, otpCode string, cartID string) (*ClientLoginResult, error) {
	client, err := s.clients.VerifyOTP(ctx, phone, otpCode)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidOTP
		}
		return nil, err
	}

    // Activate the client if they were verifying registration
    if !client.IsActive {
        if err := s.clients.ActivateClient(ctx, client.ID); err != nil {
            return nil, err
        }
        client.IsActive = true
    }

	if cartID != "" {
		if err := s.carts.AssignCartToClient(ctx, cartID, client.ID); err != nil {
			fmt.Printf("Failed to assign cart %s to client %d: %v\n", cartID, client.ID, err)
            cartID = ""
		}
	}

	return s.generateClientToken(ctx, client, cartID)
}

func (s *Service) LoginWithGoogle(ctx context.Context, googleID, email, name string) (*ClientLoginResult, error) {
	client, err := s.clients.GetByGoogleID(ctx, googleID)
	if err == nil {
		return s.generateClientToken(ctx, client, "")
	}

	if email != "" {
		client, err = s.clients.GetByEmail(ctx, email)
		if err == nil {
			if err := s.clients.UpdateGoogleID(ctx, client.ID, googleID); err != nil {
				return nil, err
			}
			return s.generateClientToken(ctx, client, "")
		}
	}

	// Create new (Active by default for Google Login)
    var pEmail, pName, pGoogleID *string
    if email != "" { pEmail = &email }
    if name != "" { pName = &name }
    pGoogleID = &googleID

	client, err = s.clients.CreateClient(ctx, nil, pEmail, pName, nil, pGoogleID, true)
	if err != nil {
		return nil, err
	}

	return s.generateClientToken(ctx, client, "")
}

func (s *Service) UpdateClient(ctx context.Context, clientID int64, email *string, username *string, dateOfBirth *time.Time) (*models.Client, error) {
    return s.clients.UpdateClient(ctx, clientID, email, username, dateOfBirth)
}

func (s *Service) ResendOTP(ctx context.Context, phone string) error {
	// Check if client exists
	_, err := s.clients.GetByPhone(ctx, phone)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrPhoneNotFound
		}
		return err
	}

    // We allow resending OTP for inactive accounts (assuming they are unverified)
    // If we had a "Banned" status separate from "Inactive/Unverified", we should check that.
	// if !client.IsActive {
	// 	return ErrInactiveAccount
	// }

	// Generate new OTP
	otpCode := s.generateOTP()
	expiry := time.Now().Add(5 * time.Minute)

	// Store OTP
	if err := s.clients.SetOTP(ctx, phone, otpCode, expiry); err != nil {
		return err
	}

	// Send SMS
	return s.sendSMS(phone, fmt.Sprintf("你的RyAngel驗證碼為: %s", otpCode))
}

func (s *Service) generateOTP() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%10]
	}
	return string(b)
}

func (s *Service) sendSMS(to, message string) error {
	if s.cfg.SkipSMSSending {
		// In development mode, skip SMS sending and just log the message
		fmt.Printf("[DEV MODE] Would send SMS to %s: %s\n", to, message)
		return nil
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.cfg.TwilioPhoneNumber)
	params.SetBody(message)

	_, err := s.twilio.Api.CreateMessage(params)
	return err
}
