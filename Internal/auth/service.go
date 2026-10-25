package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//holds the rule for the application

type Service struct {
	repo         *Repository
	secret       []byte
	ttl          time.Duration
	loginLimiter *Limiter // counts login attempts per phone number

}

func NewSrvice(repo *Repository, secret []byte, ttl time.Duration) *Service {
	return &Service{repo: repo, secret: secret, ttl: ttl}
}
func (s *Service) SetLoginLimiter(l *Limiter) { s.loginLimiter = l }

func (s *Service) SignupClient(ctx context.Context, req ClientRegister) (*AuthResponse, error) {
	phone, err := Normalizephone(req.Phone, req.Country)
	if err != nil {
		return nil, err
	}
	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	acc := Account{
		ID:           uuid.NewString(),
		Role:         RoleClient,
		Email:        Normalizemail(req.Email),
		Phone:        phone,
		PasswordHash: hash,
	}
	profile := ClientProfile{
		FullName: req.FullName,
		State:    req.State,
		Address:  req.Address,
	}
	if err := s.repo.CreateClient(ctx, acc, profile); err != nil {
		return nil, err
	}
	token, err := NewToken(s.secret, acc.ID, acc.Role, s.ttl)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, Role: string(acc.Role)}, nil
}

func (s *Service) Login(ctx context.Context, req Login) (*AuthResponse, error) {
	phone, err := Normalizephone(req.Phone, req.Country)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	//per-account limit. it runs before the database and bcrypt, so a
	//flood of guesses cannot use up your cpu
	if s.loginLimiter != nil {
		if ok, _ := s.loginLimiter.Allow(phone); !ok {
			return nil, ErrTooManyAttempts
		}
	}

	acc, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !CheckPassword(acc.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	//correct password: the real owner gets a fresh bag of tickets
	if s.loginLimiter != nil {
		s.loginLimiter.Reset(phone)
	}

	token, err := NewToken(s.secret, acc.ID, acc.Role, s.ttl)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, Role: string(acc.Role)}, nil
}

func (s *Service) Me(ctx context.Context, id string, role Role) (any, error) {
	switch role {
	case RoleClient:
		return s.repo.FindClientByID(ctx, id)
	case RoleShop:
		return s.repo.FindShopByID(ctx, id)
	case RoleRider:
		return s.repo.FindRiderByID(ctx, id)
	default:
		return nil, fmt.Errorf("unknown role %q", role)

	}

}
