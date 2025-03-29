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

func TestGroupStore_SaveSameId(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)
	group1 := group1Empty()
	err := store.Save(context.TODO(), group1)
	assert.NoError(t, err)

	// change group name, otherwise AnotherGroupWithSameNameAlreadyExistsError will come up
	group1.Name = "another name"
	err = store.Save(context.TODO(), group1)
	assert.ErrorIs(t, err, group.GroupAlreadyExistsError{ID: group1.ID})
}

func TestGroupStore_SaveNameNotUnique(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()

	store := NewGroupStore(db)

	err := store.Save(context.TODO(), group1Empty())
	assert.NoError(t, err)

	g, err := group.New(
		uuid.New(),
		"group 1", // already taken
		jeff().ID,
		make([]uuid.UUID, 0),
		[]byte{1, 2, 3, 4},
		time.Now(),
	)
	assert.NoError(t, err)
	err = store.Save(context.TODO(), *g)
	assert.ErrorIs(t, err, group.AnotherGroupWithSameNameAlreadyExistsError{Name: "group 1"})
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
