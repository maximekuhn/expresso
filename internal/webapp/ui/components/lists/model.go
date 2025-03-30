package lists

import (
	"time"

	usecaseGroup "github.com/maximekuhn/expresso/internal/usecase/group"
	"github.com/maximekuhn/expresso/internal/user"
)

type Group struct {
	Owner     *user.User
	Name      string
	CreatedAt time.Time
}

func GroupsFromUseCaseResponse(res *usecaseGroup.ListUseCaseResponse) []Group {
	groups := make([]Group, 0)
	for _, g := range res.Groups {
		groups = append(groups, Group{
			Owner:     g.Owner,
			Name:      g.Name,
			CreatedAt: g.CreatedAt,
		})
	}
	return groups
}
