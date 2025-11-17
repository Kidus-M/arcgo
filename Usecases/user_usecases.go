package Usecases

import (
	"context"
	"errors"

	"task_manager1/Domain"
	"task_manager1/Infrastructure/security"
	"task_manager1/Repositories"
)

// UserUsecase holds dependencies for user business rules.
type UserUsecase struct {
	repo Repositories.UserRepository
	pw   *security.PasswordService
}

func NewUserUsecase(r Repositories.UserRepository, pw *security.PasswordService) *UserUsecase {
	return &UserUsecase{repo: r, pw: pw}
}

// Register creates a new user. First created user becomes admin.
func (u *UserUsecase) Register(ctx context.Context, username, password string) (Domain.User, error) {
	if username == "" || password == "" {
		return Domain.User{}, errors.New("username and password required")
	}

	// hash password
	hash, err := u.pw.HashPassword(password)
	if err != nil {
		return Domain.User{}, err
	}

	user := Domain.User{
		Username:     username,
		PasswordHash: hash,
		Role:         "user",
	}

	// If first user → assign admin role
	cnt, err := u.repo.Count(ctx)
	if err != nil {
		return Domain.User{}, err
	}
	if cnt == 0 {
		user.Role = "admin"
	}

	return u.repo.Create(ctx, user)
}

// Authenticate checks username + password and returns user without hash.
func (u *UserUsecase) Authenticate(ctx context.Context, username, password string) (Domain.User, error) {
	found, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return Domain.User{}, err
	}

	if found.Username == "" {
		return Domain.User{}, errors.New("invalid credentials")
	}

	// Compare hashed password
	if !u.pw.ComparePassword(found.PasswordHash, password) {
		return Domain.User{}, errors.New("invalid credentials")
	}

	// Never return password hash
	found.PasswordHash = ""
	return found, nil
}

// Promote makes a user an admin.
func (u *UserUsecase) Promote(ctx context.Context, username string) (Domain.User, error) {
	updated, err := u.repo.PromoteToAdmin(ctx, username)
	if err != nil {
		return Domain.User{}, err
	}
	return updated, nil
}
