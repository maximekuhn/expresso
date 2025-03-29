package webapp

import (
	"database/sql"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/common"
	"github.com/maximekuhn/expresso/internal/database/sqlite"
	"github.com/maximekuhn/expresso/internal/group"
	"github.com/maximekuhn/expresso/internal/transaction"
	usecaseGroup "github.com/maximekuhn/expresso/internal/usecase/group"
	usecaseUser "github.com/maximekuhn/expresso/internal/usecase/user"
	"github.com/maximekuhn/expresso/internal/user"
)

type application struct {
	registerUsecaseHandler    *usecaseUser.RegisterUseCaseHandler
	loginUsecaseHandler       *usecaseUser.LoginUseCaseHandler
	createGroupUsecaseHandler *usecaseGroup.CreateUseCaseRequestHandler

	authService *auth.Service
	userService *user.Service

	sessionProvider transaction.SessionProvider
	cookieProvider  auth.CookieProvider
}

func newApplication(db *sql.DB, isProd bool) *application {
	idProvider := common.IdProvider{}
	datetimeProvider := common.DatetimeProvider{}
	sessionProvider := sqlite.NewSqliteSessionProvider(db)
	authStore := sqlite.NewAuthStore(db)
	authService := auth.NewService(authStore, datetimeProvider)
	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore, idProvider, datetimeProvider)
	groupStore := sqlite.NewGroupStore(db)
	groupService := group.NewService(groupStore, idProvider, datetimeProvider)

	registerUseCaseHandler := usecaseUser.NewRegisterUseCaseHandler(sessionProvider, authService, userService)
	loginUsecaseHandler := usecaseUser.NewLoginUseCaseHandler(sessionProvider, authService, datetimeProvider)
	createGroupUsecaseHandler := usecaseGroup.NewCreateUseCaseRequestHandler(sessionProvider, groupService)

	cookieProvider := auth.NewLocalhostCookieProvider()
	if isProd {
		panic("TODO: handle prod deployment")
	}

	return &application{
		registerUsecaseHandler:    registerUseCaseHandler,
		loginUsecaseHandler:       loginUsecaseHandler,
		createGroupUsecaseHandler: createGroupUsecaseHandler,
		authService:               authService,
		userService:               userService,
		sessionProvider:           sessionProvider,
		cookieProvider:            cookieProvider,
	}
}
