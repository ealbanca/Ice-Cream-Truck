package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type User struct {
	ID        int
	Name      string
	Lastname  string
	Username  string
	Email     string
	Phone     string
	Password  string // Plain password only for registration, not stored
	CreatedAt time.Time
}

// Register a new user with hashed password
func (u *User) Register() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (name, lastname, username, email, phone, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	err = config.DB.QueryRow(query, u.Name, u.Lastname, u.Username, u.Email, u.Phone, string(hash), time.Now()).Scan(&u.ID, &u.CreatedAt)
	return err
}

// Authenticate user by username/email and password
func AuthenticateUser(db *sql.DB, usernameOrEmail, password string) (*User, error) {
	var u User
	var hash string
	query := `SELECT id, name, lastname, username, email, phone, password_hash, created_at FROM users WHERE username=$1 OR email=$1`
	err := db.QueryRow(query, usernameOrEmail).Scan(&u.ID, &u.Name, &u.Lastname, &u.Username, &u.Email, &u.Phone, &hash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, err
	}
	return &u, nil
}

// RegisterUser creates a new user and saves it to the database
func RegisterUser(user *User) error {
	return user.Register()
}
