package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//holds the rule for the application

type Service struct {
	repo   *Repository
	secret []byte
	ttl    time.Duration
}

func NewSrvice(repo *Repository, secret []byte, ttl time.Duration) *Service {
	return &Service{repo: repo, secret: secret, ttl: ttl}
}

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
	if err := s.repo.CreatClient(ctx, acc, profile); err != nil {
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
	token, err := NewToken(s.secret, acc.ID, acc.Role, s.ttl)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, Role: string(acc.Role)}, nil
}
