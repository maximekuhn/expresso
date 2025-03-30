package group

import (
	"context"
	"fmt"
	"time"

	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/transaction"
	"github.com/maximekuhn/expresso/internal/user"
)

type ListUseCaseRequest struct {
	User *user.User
}

type ListUseCaseResponse struct {
	Groups []ListUseCaseResponseGroup
}

type ListUseCaseResponseGroup struct {
	Owner     *user.User
	Name      string
	CreatedAt time.Time
}

type ListUseCaseRequestHandler struct {
	seessionProvider transaction.SessionProvider
	groupService     *group.Service
	userService      *user.Service
}

func NewListUseCaseRequestHandler(
	seessionProvider transaction.SessionProvider,
	groupService *group.Service,
	userService *user.Service,
) *ListUseCaseRequestHandler {
	return &ListUseCaseRequestHandler{
		seessionProvider: seessionProvider,
		groupService:     groupService,
		userService:      userService,
	}
}

func (h *ListUseCaseRequestHandler) Handle(ctx context.Context, r *ListUseCaseRequest) (*ListUseCaseResponse, error) {
	var res *ListUseCaseResponse

	session := h.seessionProvider.Provide(ctx)
	err := session.Transaction(ctx, func(ctx context.Context) error {
		groups, err := h.groupService.ListGroupOfUser(ctx, r.User.ID)
		if err != nil {
			return err
		}

		// TODO: fetch all user at once (provide a list of IDs)
		resGroups := make([]ListUseCaseResponseGroup, 0)
		for _, g := range groups {
			owner, found, err := h.userService.Get(ctx, g.Owner)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("user not found[id=%s]", g.Owner)
			}
			resGroups = append(resGroups, ListUseCaseResponseGroup{
				Owner:     owner,
				Name:      g.Name,
				CreatedAt: g.CreatedAt,
			})
		}

		res = &ListUseCaseResponse{
			Groups: resGroups,
		}

		return nil
	})
	return res, err
}
