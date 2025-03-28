package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/group"
	"github.com/stretchr/testify/assert"
)

func TestGroupStore_Save(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)
	err := store.Save(context.TODO(), group1Empty())
	assert.NoError(t, err)
}

func TestGroupStore_SaveDuplicate(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)
	group1 := group1Empty()
	err := store.Save(context.TODO(), group1)
	assert.NoError(t, err)

	err = store.Save(context.TODO(), group1)
	assert.ErrorIs(t, err, group.GroupAlreadyExistsError{ID: group1.ID})
}

func group1Empty() group.Group {
	g, err := group.New(
		uuid.New(),
		"group 1",
		jeff().ID,
		make([]uuid.UUID, 0),
		[]byte{7, 2, 12, 23},
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	return *g
}
