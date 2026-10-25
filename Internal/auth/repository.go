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

// type ClientProfile struct {
// 	FullName string
// 	State    string
// 	Address  string
// }

// CreateClient saves an account and its client profile as one unit
func (r *Repository) withTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// return tx.commit(ctx)

func insertAccount(ctx context.Context, tx pgx.Tx, acc Account) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO accounts (id, role, email, phone, password_hash)
	VALUES ($1, $2, $3, $4, $5)`,
		acc.ID, string(acc.Role), acc.Email, acc.Phone, acc.PasswordHash,
	)
	// err := insertAccount(ctx, tx, acc); err

	return mapInsertError(err)
}

func (r *Repository) CreateClient(ctx context.Context, acc Account, p ClientProfile) error {
	return r.withTx(ctx, func(tx pgx.Tx) error {
		if err := insertAccount(ctx, tx, acc); err != nil {
			return err
		}

		_, err := tx.Exec(ctx,
			`INSERT INTO client_profile (account_id, full_name, state, address)
VALUES ($1, $2, $3, $4)`, acc.ID, p.FullName, p.State, p.Address,
		)
		// if err != nil {
		// 	return err
		// }
		return err
	})

}

func (r *Repository) CreateShop(ctx context.Context, acc Account, p ShopProfile) error {

	return r.withTx(ctx, func(tx pgx.Tx) error {
		if err := insertAccount(ctx, tx, acc); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO shop_profiles(account_id, business_name, state, address, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			acc.ID, p.BusinessName, p.State, p.Address, p.Latitude, p.Longitude,
		)
		return err
	})
}

func (r *Repository) CreateRider(ctx context.Context, acc Account, p RiderProfile) error {
	return r.withTx(ctx, func(tx pgx.Tx) error {
		if err := insertAccount(ctx, tx, acc); err != nil {
			return err
		}
		_, err := tx.Exec(ctx,
			`INSERT INTO rider_profiles (account_id, full_name, registered_company, plate_number)
			VALUES ($1,$2,$3,$4)`,
			acc.ID, p.FullName, p.RegisteredCompany, p.PlateNumber,
		)
		return err
	})
}

// FindByPhine returns the account with this (already normalized) phone number
func (r *Repository) FindByPhone(ctx context.Context, phone string) (*Account, error) {
	var acc Account
	var role string

	err := r.db.QueryRow(ctx,
		`SELECT id, role, email, phone, password_hash, created_at
FROM accounts
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
func (r *Repository) FindClientByID(ctx context.Context, id string) (*ClientResponse, error) {
	var c ClientResponse

	err := r.db.QueryRow(ctx,
		`SELECT a.id::text, c.full_name, a.email, a.phone, c.state, c.address
	FROM accounts a
	JOIN client_profiles c ON c.account_id = a.id
	WHERE a.id = $1`,
		id,
	).Scan(&c.ID, &c.FullName, &c.Email, &c.Phone, &c.State, &c.Address)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindShopByID(ctx context.Context, id string) (*ShopResponse, error) {
	var s ShopResponse

	err := r.db.QueryRow(ctx,
		`SELECT a.id::text, s.business_name, a.email, a.phone, s.state, s.address, s.latitude, s.longitude, s.is_verified
	FROM accounts a
	JOIN shop_profiles s ON s.account_id = a.id
	WHERE a.id = $1`,
		id).Scan(&s.ID, &s.BusinessName, &s.Email, &s.Phone, &s.State,
		&s.Address, &s.Latitude, &s.Longitude, &s.IsVerified,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) FindRiderByID(ctx context.Context, id string) (*RiderResponse, error) {
	var rd RiderResponse
	err := r.db.QueryRow(ctx,
		`SELECT a.id::text, p.full_name, a.email, a.phone, p.registered_company,
	p.plate_number, p.is_verified
	FROM accounts a
	JOIN rider_profiles p ON p.account_id = a.id
	WHERE a.id = &1`,
		id).Scan(&rd.ID, &rd.FullName, &rd.Email, &rd.Phone, &rd.RegisteredCompany, &rd.PlateNumber, &rd.IsVerified)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rd, nil
}

// turns a duplicate database error into our own error
func mapInsertError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "accounts_email_key", "accounts_email_lower_idx":
			return ErrEmailTaken
		case "accounts_phone_key":
			return ErrPhoneTaken
		}
	}
	return err
}
