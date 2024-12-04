package repository

import (
	"database/sql"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/constants"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/utils"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *database.DB
}

func NewAuthRepo(db *database.DB) AuthRepository {
	return &AuthRepo{db}
}

func (r *AuthRepo) Register(input *model.RegisterInput) (*model.AuthResponse, error) {
	// Check if email exists
	exists, err := r.emailExists(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, constants.ErrEmailExists
	}
	input.RoleID = 2
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert user
	var user model.AuthOutput
	err = tx.QueryRow(`
		INSERT INTO users (role_id, email, password, full_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, email, full_name, role_id, status, created_at, updated_at
	`, input.RoleID, input.Email, string(hashedPassword), input.FullName).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.RoleID,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateNewTokens(uint(user.ID), user.Email, uint(user.RoleID), user.RoleID == 1)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

func (r *AuthRepo) Login(input *model.LoginInput) (*model.AuthResponse, error) {
	user, err := r.GetUserByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	// Check if user is active
	if !user.Status {
		return nil, constants.ErrUserDeactivated
	}

	// Get stored password
	var hashedPassword string
	err = r.db.QueryRow("SELECT password FROM users WHERE email = $1", input.Email).Scan(&hashedPassword)
	if err != nil {
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		return nil, constants.ErrInvalidPassword
	}

	// Generate token
	token, err := utils.GenerateNewTokens(uint(user.ID), user.Email, uint(user.RoleID), user.RoleID == 1)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (r *AuthRepo) GetUserByEmail(email string) (*model.AuthOutput, error) {
	var user model.AuthOutput
	err := r.db.QueryRow(`
		SELECT id, email, full_name, role_id, status, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.RoleID,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Helper function to check if email exists
func (r *AuthRepo) emailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}
