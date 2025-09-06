package users

import (
	"database/sql"

	"github.com/yourusername/saas-billing/internal/auth"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Password is never returned in JSON
}

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(email, password string) error {
	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	// Insert the user
	_, err = s.db.Exec(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
	`, email, hashedPassword)

	return err
}

func (s *UserService) Login(email, password string) (string, error) {
	var user User
	var hashedPassword string

	// Get the user
	err := s.db.QueryRow(`
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &hashedPassword)

	if err == sql.ErrNoRows {
		return "", err
	}

	if err != nil {
		return "", err
	}

	// Check password
	if !auth.CheckPasswordHash(password, hashedPassword) {
		return "", err
	}

	// Generate JWT token
	return auth.GenerateToken(user.ID)
}
