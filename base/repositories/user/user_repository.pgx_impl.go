package user_repository

import (
	"context"
	"errors"
	"strings"

	"gaudiot.com/fonli/core/database"

	"github.com/jackc/pgx/v5"
)

type pgxUserRepository struct {
	db *database.DB
}

func NewPgxUserRepository(db *database.DB) UserRepository {
	return &pgxUserRepository{db: db}
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrUsernameExists    = errors.New("username already exists")
)

// GetUserByID fetches a user by its ID.
func (r *pgxUserRepository) GetUserByID(id string) (*User, error) {
	ctx := context.Background()
	query := `
		SELECT id, email, password, username, canonical_username, lifestyle, lifestyle_topics
		FROM users WHERE id = $1
	`
	row := r.db.Pool.QueryRow(ctx, query, id)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.CanonicalUsername,
		&u.Lifestyle,
		&u.LifestyleTopics,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail fetches a user by email (case insensitive).
func (r *pgxUserRepository) GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	query := `
		SELECT id, email, password, username, canonical_username, lifestyle, lifestyle_topics
		FROM users WHERE lower(email) = lower($1)
	`
	row := r.db.Pool.QueryRow(ctx, query, email)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.CanonicalUsername,
		&u.Lifestyle,
		&u.LifestyleTopics,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByUsername fetches a user by username (case insensitive).
func (r *pgxUserRepository) GetUserByUsername(username string) (*User, error) {
	ctx := context.Background()
	query := `
		SELECT id, email, password, username, canonical_username, lifestyle, lifestyle_topics
		FROM users WHERE lower(username) = lower($1)
	`
	row := r.db.Pool.QueryRow(ctx, query, username)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.CanonicalUsername,
		&u.Lifestyle,
		&u.LifestyleTopics,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmailOrUsername fetches a user by matching email or username (case insensitive).
func (r *pgxUserRepository) GetUserByEmailOrUsername(text string) (*User, error) {
	ctx := context.Background()
	val := strings.ToLower(text)
	query := `
		SELECT id, email, password, username, canonical_username, lifestyle, lifestyle_topics
		FROM users
		WHERE lower(email) = $1 OR lower(username) = $1
	`
	row := r.db.Pool.QueryRow(ctx, query, val)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.CanonicalUsername,
		&u.Lifestyle,
		&u.LifestyleTopics,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser inserts a new user and returns it.
func (r *pgxUserRepository) CreateUser(email, password, username string) (*User, error) {
	ctx := context.Background()
	canonicalUsername := strings.ToLower(username)
	query := `
		INSERT INTO users (email, password, username, canonical_username)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password, username, canonical_username, lifestyle, lifestyle_topics
	`
	row := r.db.Pool.QueryRow(ctx, query, email, password, username, canonicalUsername)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Username,
		&u.CanonicalUsername,
		&u.Lifestyle,
		&u.LifestyleTopics,
	)
	if err != nil {
		// Optionally, you can check for duplicate key errors here if needed, based on the database.
		return nil, err
	}
	return &u, nil
}

// DeleteUser deletes a user by ID.
func (r *pgxUserRepository) DeleteUser(userID string) error {
	ctx := context.Background()
	query := `DELETE FROM users WHERE id = $1`
	ct, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UpdateUserLifestyle updates lifestyle and lifestyle_topics for a user by ID.
func (r *pgxUserRepository) UpdateUserLifestyle(userID, lifestyle, lifestyleTopics string) error {
	ctx := context.Background()
	query := `
		UPDATE users SET lifestyle = $1, lifestyle_topics = $2, updated_at = NOW() WHERE id = $3
	`
	ct, err := r.db.Pool.Exec(ctx, query, lifestyle, lifestyleTopics, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
