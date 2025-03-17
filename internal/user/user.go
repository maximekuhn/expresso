package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

func New(id uuid.UUID, name string, createdAt time.Time) (*User, error) {
	u := &User{
		ID:        id,
		Name:      strings.TrimSpace(name),
		CreatedAt: createdAt.UTC(),
	}
	return u, u.validate()
}

func (u *User) validate() error {
	return ValidateName(u.Name)
}
