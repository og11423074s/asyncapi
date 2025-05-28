package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	Id                   uuid.UUID `db:"id"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hashed_password"`
	CreatedAt            time.Time `db:"created_at"`
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (u *User) ComparePassword(password string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(u.HashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("error failed to decode hashed password: %w", err)
	}

	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	const dml = `INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *;`
	var user User

	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error failed to hash password: %w", err)
	}

	// base64 encode the hashed password
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)

	err = s.db.GetContext(ctx, &user, dml, email, hashedPasswordBase64)
	if err != nil {
		return nil, fmt.Errorf("error failed to insert user: %w", err)
	}

	return &user, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	const query = `SELECT * FROM users WHERE email = $1;`
	var user User
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("error failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserStore) ById(ctx context.Context, userId uuid.UUID) (*User, error) {
	const query = `SELECT * FROM users WHERE id = $1;`
	var user User
	err := s.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error failed to get user by id %s: %w", userId, err)
	}
	return &user, nil
}
