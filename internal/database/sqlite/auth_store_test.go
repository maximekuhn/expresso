package sqlite

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthStore_Save(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewAuthStore(db)
	err := store.Save(context.TODO(), authEntry())
	assert.NoError(t, err)
}

func TestAuthStore_Duplicate(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewAuthStore(db)

	entry := authEntry()

	err := store.Save(context.TODO(), entry)
	assert.NoError(t, err)

	err = store.Save(context.TODO(), entry)
	assert.ErrorIs(t, err, auth.EntryAlreadyExistsError{
		Email:  entry.Email,
		UserID: entry.UserID,
	})
}

func authEntry() auth.Entry {
	e, err := auth.NewEntry(
		"jeff@gmail.com",
		[]byte{1, 2, 3, 4, 5, 6},
		uuid.New(),
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return *e
}
