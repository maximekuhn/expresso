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

func jeff() user.User {
	u, err := user.New(uuid.New(), "jeff", time.Now())
	if err != nil {
		panic(err)
	}
	return *u
}
