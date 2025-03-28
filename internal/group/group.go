package group

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID             uuid.UUID
	Name           string
	Owner          uuid.UUID
	Members        []uuid.UUID
	HashedPassword []byte
	CreatedAt      time.Time
}

func New(
	id uuid.UUID,
	name string,
	owner uuid.UUID,
	members []uuid.UUID,
	hashedPassword []byte,
	createdAt time.Time,
) (*Group, error) {
	g := &Group{
		ID:             id,
		Name:           strings.TrimSpace(name),
		Owner:          owner,
		Members:        members,
		HashedPassword: hashedPassword,
		CreatedAt:      createdAt.UTC(),
	}
	return g, g.validate()
}

func (g *Group) validate() error {
	return ValidateName(g.Name)
}
