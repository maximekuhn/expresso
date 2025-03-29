package sqlite

import (
	"context"
	"testing"
	"time"

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

func TestAuthStore_SaveDuplicate(t *testing.T) {
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

func TestAuthStore_GetByEmailOk(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewAuthStore(db)
	entry := authEntry()
	err := store.Save(context.TODO(), entry)
	assert.NoError(t, err)

	e, found, err := store.GetByEmail(context.TODO(), entry.Email)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, *e, entry)
}

func TestAuthStore_GetByEmailNotFound(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewAuthStore(db)

	_, found, err := store.GetByEmail(context.TODO(), "bill@gmail.com")
	assert.NoError(t, err)
	assert.False(t, found, "should not find an auth entry for email: bill@gmail.com")
}

func TestAuthStore_GetBySessionIdOk(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()

	store := NewAuthStore(db)
	entry := authEntry()
	err := store.Save(context.TODO(), entry)
	assert.NoError(t, err)

	entryWithSession := addSessionToEntry(entry)
	err = store.Update(context.TODO(), entry, entryWithSession)
	assert.NoError(t, err)

	e, found, err := store.GetBySessionID(context.TODO(), entryWithSession.SessionID)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, *e, entryWithSession)
}

func TestAuthStore_GetBySessionIdNotFound(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()

	store := NewAuthStore(db)

	sessionId := uuid.NewString()
	_, found, err := store.GetBySessionID(context.TODO(), sessionId)
	assert.NoError(t, err)
	assert.Falsef(t, found, "should not find an auth entry for sessionId: %s", sessionId)
}

func TestAuthStore_UpdateOk(t *testing.T) {
	scenarios := []struct {
		title               string
		entry               auth.Entry
		newEmail            string     // empty if no update
		newHashedPassword   []byte     // nil if no update
		newSessionID        string     // empty if no update
		newSessionExpiresAt *time.Time // nil if no update
	}{
		{
			title:               "update email only",
			entry:               authEntry(),
			newEmail:            "jeff2@gmail.com",
			newHashedPassword:   nil,
			newSessionID:        "",
			newSessionExpiresAt: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.title, func(t *testing.T) {
			db := createTmpDbWithAllMigrationsApplied()
			store := NewAuthStore(db)
			assert.NoError(t, store.Save(context.TODO(), scenario.entry))

			old := scenario.entry

			e := old.Email
			if scenario.newEmail != "" {
				e = scenario.newEmail
			}
			hp := old.HashedPassword
			if scenario.newHashedPassword != nil {
				hp = scenario.newHashedPassword
			}
			sId := old.SessionID
			if scenario.newSessionID != "" {
				sId = scenario.newSessionID
			}
			sExpiresAt := old.SessionExpiresAt
			if scenario.newSessionExpiresAt != nil {
				sExpiresAt = scenario.newSessionExpiresAt
			}

			new, err := auth.NewEntry(e, hp, old.UserID, sId, sExpiresAt)
			if err != nil {
				panic(err)
			}

			err = store.Update(context.TODO(), old, *new)
			assert.NoError(t, err)

			inDb, found, err := store.GetByEmail(context.TODO(), new.Email)
			assert.NoError(t, err)
			assert.True(t, found)
			assert.Equal(t, *new, *inDb)
		})
	}
}

func TestAuthStore_UpdateErr(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	store := NewAuthStore(db)

	entry := authEntry()
	assert.NoError(t, store.Save(context.TODO(), entry))

	// try to change user ID
	new, err := auth.NewEntry(entry.Email, entry.HashedPassword, uuid.New(), entry.SessionID, entry.SessionExpiresAt)
	if err != nil {
		panic(err)
	}
	err = store.Update(context.TODO(), entry, *new)
	assert.Error(t, err)
	assert.Equal(t, "changing UserID is not supported/forbidden as it is the primary key", err.Error())
}

func authEntry() auth.Entry {
	e, err := auth.NewEntry(
		"jeff@gmail.com",
		[]byte{1, 2, 3, 4, 5, 6},
		uuid.New(),
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	return *e
}

func addSessionToEntry(e auth.Entry) auth.Entry {
	expiresAt := time.Now().UTC().Add(time.Hour * 24)
	e.SessionID = uuid.NewString()
	e.SessionExpiresAt = &expiresAt
	return e
}
