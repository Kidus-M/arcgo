package infrastructure

import "golang.org/x/crypto/bcrypt"

type PasswordService struct{}

func NewPasswordService() *PasswordService { return &PasswordService{} }

func (p *PasswordService) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (p *PasswordService) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
