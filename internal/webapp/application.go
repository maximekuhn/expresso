package webapp

import (
	"database/sql"

	"github.com/maximekuhn/expresso/internal/auth"
	"github.com/maximekuhn/expresso/internal/common"
	"github.com/maximekuhn/expresso/internal/database/sqlite"
	usecaseUser "github.com/maximekuhn/expresso/internal/usecase/user"
	"github.com/maximekuhn/expresso/internal/user"
)

type application struct {
	registerUsecaseHandler *usecaseUser.RegisterUseCaseHandler
	loginUsecaseHandler    *usecaseUser.LoginUseCaseHandler
}

func newApplication(db *sql.DB) *application {
	idProvider := common.IdProvider{}
	datetimeProvider := common.DatetimeProvider{}
	sessionProvider := sqlite.NewSqliteSessionProvider(db)
	authStore := sqlite.NewAuthStore(db)
	authService := auth.NewService(authStore, datetimeProvider)
	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore, idProvider, datetimeProvider)

	registerUseCaseHandler := usecaseUser.NewRegisterUseCaseHandler(sessionProvider, authService, userService)
	loginUsecaseHandler := usecaseUser.NewLoginUseCaseHandler(sessionProvider, authService, datetimeProvider)

	return &application{
		registerUsecaseHandler: registerUseCaseHandler,
		loginUsecaseHandler:    loginUsecaseHandler,
	}
}
