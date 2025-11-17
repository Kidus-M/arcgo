package usecases

import (
	"context"
	"errors"

	"task_manager/domain"
	"task_manager/infrastructure"
	"task_manager/repositories"
)

// UserUsecase defines user-related business rules
type UserUsecase struct {
	repo repositories.UserRepository
	pw   *infrastructure.PasswordService
}

func NewUserUsecase(r repositories.UserRepository, pw *infrastructure.PasswordService) *UserUsecase {
	return &UserUsecase{repo: r, pw: pw}
}

// Register creates a user. First user becomes admin.
func (u *UserUsecase) Register(ctx context.Context, username, password string) (domain.User, error) {
	if username == "" || password == "" {
		return domain.User{}, errors.New("username and password required")
	}
	// hash password
	hash, err := u.pw.Hash(password)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		Username:     username,
		PasswordHash: hash,
		Role:         "user",
	}
	// if first user, promote to admin
	cnt, err := u.repo.Count(ctx)
	if err != nil {
		return domain.User{}, err
	}
	if cnt == 0 {
		user.Role = "admin"
	}

	return u.repo.Create(ctx, user)
}

// Authenticate checks credentials and returns user (without hash)
func (u *UserUsecase) Authenticate(ctx context.Context, username, password string) (domain.User, error) {
	found, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	if found.Username == "" {
		return domain.User{}, errors.New("invalid credentials")
	}
	if err := u.pw.Compare(found.PasswordHash, password); err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}
	found.PasswordHash = ""
	return found, nil
}

func (u *UserUsecase) Promote(ctx context.Context, username string) (domain.User, error) {
	updated, err := u.repo.PromoteToAdmin(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	if updated.Username == "" {
		return domain.User{}, nil
	}
	return updated, nil
}
