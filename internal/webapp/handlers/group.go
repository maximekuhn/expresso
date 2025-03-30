package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/logger"
	usecaseGroup "github.com/maximekuhn/expresso/internal/usecase/group"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
	"github.com/maximekuhn/expresso/internal/webapp/ui/components/lists"
)

type GroupHandler struct {
	logger        *slog.Logger
	createUseCase *usecaseGroup.CreateUseCaseRequestHandler
	listUseCase   *usecaseGroup.ListUseCaseRequestHandler
}

func NewGroupHandler(
	l *slog.Logger,
	createUseCase *usecaseGroup.CreateUseCaseRequestHandler,
	listUseCase *usecaseGroup.ListUseCaseRequestHandler,
) *GroupHandler {
	return &GroupHandler{
		logger:        l.With(slog.String(logger.LoggerNameField, "GroupHandler")),
		createUseCase: createUseCase,
		listUseCase:   listUseCase,
	}
}

func (h *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.createGroup(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.getList(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *GroupHandler) createGroup(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)
	if err := r.ParseForm(); err != nil {
		l.Error("failed to parse form", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loggedUser := extractUserOrReturnInternalError(l, w, r)
	if loggedUser == nil {
		return
	}

	groupname := r.PostForm.Get("name")
	password := r.PostForm.Get("password")

	if err := h.createUseCase.Handle(r.Context(), &usecaseGroup.CreateUseCaseRequest{
		Owner:     loggedUser,
		GroupName: groupname,
		Password:  password,
	}); err != nil {
		h.handleCreateError(l, w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	l.Info("group created")
}

func (h *GroupHandler) getList(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeWithRequestId(r.Context(), middleware.RequestIdKey{}, h.logger)

	loggedUser := extractUserOrReturnInternalError(l, w, r)
	if loggedUser == nil {
		return
	}

	res, err := h.listUseCase.Handle(r.Context(), &usecaseGroup.ListUseCaseRequest{
		User: loggedUser,
	})
	if err != nil {
		l.Error("failed to retrieve list", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info("found groups", slog.Int("count", len(res.Groups)))

	if err := lists.GroupsList(lists.GroupsFromUseCaseResponse(res)).Render(r.Context(), w); err != nil {
		l.Error("failed to render lists.GroupsList")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (_ *GroupHandler) handleCreateError(l *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var groupAlreadyExistsError group.GroupAlreadyExistsError
	if errors.As(err, &groupAlreadyExistsError) {
		l.Info(
			"ID already taken by another group",
			slog.String("id", string(groupAlreadyExistsError.ID.String())),
		)

		// this error should never happen, as group ID are UUID
		// we will return internal error for now
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var anotherGroupSameNameExistsError group.AnotherGroupWithSameNameAlreadyExistsError
	if errors.As(err, &anotherGroupSameNameExistsError) {
		l.Info(
			"another group with same name already exists",
			slog.String("name", anotherGroupSameNameExistsError.Name),
		)
		returnConflictAndBoxError(
			"Name not available. Try again with a different one",
			l, w, r,
		)
		w.WriteHeader(http.StatusConflict)
		return
	}

	l.Error("internal error", slog.String("err", err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
