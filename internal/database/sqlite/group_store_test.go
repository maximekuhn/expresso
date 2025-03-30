package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/user"
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

func TestGroupStore_GetAllWhereUserIsOwner(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	jeff := jeff()
	bill := bill()

	group1Empty := createGroup(jeff, "group 1")
	group2Empty := createGroup(jeff, "group 2")
	group3Empty := createGroup(bill, "group 3")

	assert.NoError(t, store.Save(context.TODO(), group1Empty))
	assert.NoError(t, store.Save(context.TODO(), group2Empty))
	assert.NoError(t, store.Save(context.TODO(), group3Empty))

	groups, err := store.GetAllWhereUserIsOwner(context.TODO(), jeff.ID)
	assert.NoError(t, err)

	expectedGroups := []group.Group{
		group1Empty,
		group2Empty,
		// group3 is owned by bill, not jeff
	}
	assert.ElementsMatch(t, expectedGroups, groups)
}

func TestGroupStore_GetAllWhereUserIsOwnerEmptyList(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	groups, err := store.GetAllWhereUserIsOwner(context.TODO(), jeff().ID)
	assert.NoError(t, err)
	assert.Empty(t, groups)
}

func TestGroupStore_GetAllWhereUserIsMember(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	jeff := jeff()
	bill := bill()

	group1 := createGroupWithOneMember(jeff, bill, "group 1")
	group2 := createGroup(jeff, "group 2")
	group3 := createGroupWithOneMember(bill, jeff, "group 3")

	assert.NoError(t, store.Save(context.TODO(), group1))
	assert.NoError(t, store.Save(context.TODO(), group2))
	assert.NoError(t, store.Save(context.TODO(), group3))

	groups, err := store.GetAllWhereUserIsMember(context.TODO(), jeff.ID)
	assert.NoError(t, err)

	expectedGroups := []group.Group{
		group3,
		// jeff is only member of group 3
		// he is owner of group 1 and group 2
	}
	assert.ElementsMatch(t, expectedGroups, groups)
}

func TestGroupStore_GetAllWhereUserIsMemberEmptyList(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	groups, err := store.GetAllWhereUserIsMember(context.TODO(), jeff().ID)
	assert.NoError(t, err)
	assert.Empty(t, groups)
}

func TestGroupStore_GetByGroupName(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	g := createGroupWithOneMember(jeff(), bill(), "Jeff & Bill group")
	assert.NoError(t, store.Save(context.TODO(), g))

	gr, found, err := store.GetByGroupName(context.TODO(), g.Name)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, g, *gr)
}

func TestGroupStore_GetByGroupNameNotFound(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	_, found, err := store.GetByGroupName(context.TODO(), "non existing group name")
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestGroupStore_AddMember(t *testing.T) {
	db := createTmpDbWithAllMigrationsApplied()
	defer db.Close()
	store := NewGroupStore(db)

	g := createGroupWithOneMember(jeff(), bill(), "Jeff & Bill group")
	assert.NoError(t, store.Save(context.TODO(), g))

	bob := bob()
	assert.NoError(t, store.AddMember(context.TODO(), g.ID, bob.ID))

	gr, found, err := store.GetByGroupName(context.TODO(), g.Name)
	assert.NoError(t, err)
	assert.True(t, found)
	g.Members = append(g.Members, bob.ID)
	assert.ElementsMatch(t, g.Members, gr.Members)
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

func createGroup(owner user.User, name string) group.Group {
	g, err := group.New(
		uuid.New(),
		name,
		owner.ID,
		make([]uuid.UUID, 0),
		[]byte{7, 2, 12, 23},
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	return *g
}

func createGroupWithOneMember(owner, member user.User, name string) group.Group {
	members := make([]uuid.UUID, 0)
	members = append(members, member.ID)
	g, err := group.New(
		uuid.New(),
		name,
		owner.ID,
		members,
		[]byte{7, 2, 12, 23},
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	return *g
}
