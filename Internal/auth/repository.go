package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type ClientProfile struct {
	FullName string
	State    string
	Address  string
}

// CreateClient saves an account and its client profile as one unit
func (r *Repository) CreatClient(ctx context.Context, acc Account, p ClientProfile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO account (id, role, email, phone, password_hash)
	VALUES ($1, $2, $3, $4, $5)`,
		acc.ID, string(acc.Role), acc.Email, acc.Phone, acc.PasswordHash,
	)
	if err != nil {
		return mapInsertError(err)
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO client_profile (account_id, full_name, state, address)
VALUES ($1, $2, $3, $4)`, acc.ID, p.FullName, p.State, p.Address,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindByPhine returns the account with this (already normalized) phone number
func (r *Repository) FindByPhone(ctx context.Context, phone string) (*Account, error) {
	var acc Account
	var role string

	err := r.db.QueryRow(ctx,
		`SELECT id, role, email, phone, password_hash, created_at
FROM account
WHERE phone = $1`,
		phone,
	).Scan(&acc.ID, &role, &acc.Email, &acc.Phone, &acc.PasswordHash, &acc.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	acc.Role = Role(role)
	return &acc, nil
}

func mapInsertError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "accounts_email_key", "account_email_lower_idx":
			return ErrEmailTaken
		case "accounts_phone_key":
			return ErrPhoneTaken
		}
	}
	return err
}
