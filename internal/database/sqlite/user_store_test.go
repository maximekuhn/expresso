package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserStore_Save(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewUserStore(db)
	err := store.Save(context.TODO(), jeff())
	assert.NoError(t, err)
}

func TestUserStore_SaveDuplicate(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewUserStore(db)

	jeff := jeff()

	err := store.Save(context.TODO(), jeff)
	assert.NoError(t, err)

	err = store.Save(context.TODO(), jeff)
	assert.ErrorIs(t, err, user.UserAlreadyExistsError{ID: jeff.ID})
}

func TestUserStore_GetByIdOk(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()

	store := NewUserStore(db)

	jeff := jeff()

	err := store.Save(context.TODO(), jeff)
	assert.NoError(t, err)

	u, found, err := store.GetById(context.TODO(), jeff.ID)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, *u, jeff)
}

func TestUserStore_NotFound(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()

	store := NewUserStore(db)

	jeff := jeff()

	_, found, err := store.GetById(context.TODO(), jeff.ID)
	assert.NoError(t, err)
	assert.False(t, found)
}

func jeff() user.User {
	u, err := user.New(uuid.New(), "jeff", time.Now())
	if err != nil {
		panic(err)
	}
	return *u
}

func bill() user.User {
	u, err := user.New(uuid.New(), "bill", time.Now())
	if err != nil {
		panic(err)
	}
	return *u
}
