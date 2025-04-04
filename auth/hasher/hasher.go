package hasher

import "golang.org/x/crypto/bcrypt"

type Hash interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashsedPassword, plainPassword string) bool
}

type hash struct {
}

func NewHasher() Hash {
	return &hash{}
}

func (h *hash) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *hash) ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
