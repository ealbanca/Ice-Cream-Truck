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
	UserType  string
}

// Register a new user with hashed password
// Check if a username already exists
func UsernameExists(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := config.DB.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// Check if an email already exists
func EmailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := config.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// Check if a phone number already exists
func PhoneExists(phone string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)"
	err := config.DB.QueryRow(query, phone).Scan(&exists)
	return exists, err
}
func (u *User) Register() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (name, lastname, username, email, phone, password_hash, created_at, user_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, user_type`
	err = config.DB.QueryRow(query, u.Name, u.Lastname, u.Username, u.Email, u.Phone, string(hash), time.Now(), u.UserType).Scan(&u.ID, &u.CreatedAt, &u.UserType)
	return err
}

// Authenticate user by username/email and password
func AuthenticateUser(db *sql.DB, usernameOrEmail, password string) (*User, error) {
	var u User
	var hash string
	query := `SELECT id, name, lastname, username, email, phone, password_hash, created_at, user_type FROM users WHERE username=$1 OR email=$1`
	err := db.QueryRow(query, usernameOrEmail).Scan(&u.ID, &u.Name, &u.Lastname, &u.Username, &u.Email, &u.Phone, &hash, &u.CreatedAt, &u.UserType)
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

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (*User, error) {
	var u User
	query := `SELECT id, name, lastname, username, email, phone, created_at, user_type FROM users WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Lastname, &u.Username, &u.Email, &u.Phone, &u.CreatedAt, &u.UserType)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUser updates the user's account information in the database
func UpdateUser(user *User) error {
	query := `UPDATE users SET name=$1, lastname=$2, username=$3, email=$4, phone=$5 WHERE id=$6`
	_, err := config.DB.Exec(query, user.Name, user.Lastname, user.Username, user.Email, user.Phone, user.ID)
	return err
}

// UpdateUserPassword updates the user's password hash in the database
func UpdateUserPassword(userID int, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `UPDATE users SET password_hash=$1 WHERE id=$2`
	_, err = config.DB.Exec(query, string(hash), userID)
	return err
}
