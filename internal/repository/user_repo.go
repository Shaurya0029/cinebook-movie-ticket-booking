package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"movieticketbooking/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateEmail = errors.New("email already registered")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, firstName, lastName, email, passwordHash string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, first_name, last_name, email, password_hash, created_at`,
		firstName, lastName, email, passwordHash).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, password_hash, created_at FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, password_hash, created_at FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
