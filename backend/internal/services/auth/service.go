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
}

type Service struct {
	admins  *repository.AdminRepository
	clients *repository.ClientRepository
	cfg     *config.Config
	twilio  *twilio.RestClient
}

func NewService(admins *repository.AdminRepository, clients *repository.ClientRepository, cfg *config.Config) *Service {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &Service{
		admins:  admins,
		clients: clients,
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

func (s *Service) ClientRegister(ctx context.Context, phone string, email *string, username *string) (*models.Client, error) {
	// Check if phone already exists
	_, err := s.clients.GetByPhone(ctx, phone)
	if err == nil {
		return nil, errors.New("phone number already registered")
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	return s.clients.CreateClient(ctx, phone, email, username)
}

func (s *Service) UpdateClient(ctx context.Context, clientID int64, email *string, username *string) (*models.Client, error) {
	return s.clients.UpdateClient(ctx, clientID, email, username)
}

func (s *Service) VerifyOTP(ctx context.Context, phone, otpCode string) (*ClientLoginResult, error) {
	client, err := s.clients.VerifyOTP(ctx, phone, otpCode)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidOTP
		}
		return nil, err
	}

	if !client.IsActive {
		return nil, ErrInactiveAccount
	}

	// Generate token
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

	return &ClientLoginResult{Token: raw, ExpiresAt: expiresAt, Client: client}, nil
}

func (s *Service) ResendOTP(ctx context.Context, phone string) error {
	// Check if client exists and is active
	client, err := s.clients.GetByPhone(ctx, phone)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrPhoneNotFound
		}
		return err
	}

	if !client.IsActive {
		return ErrInactiveAccount
	}

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
