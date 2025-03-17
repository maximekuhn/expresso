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
}

func newApplication(db *sql.DB) *application {
	sessionProvider := sqlite.NewSqliteSessionProvider(db)
	authStore := sqlite.NewAuthStore(db)
	authService := auth.NewService(authStore)
	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore, common.IdProvider{}, common.DatetimeProvider{})

	registerUseCaseHandler := usecaseUser.NewRegisterUseCaseHandler(sessionProvider, authService, userService)

	return &application{
		registerUsecaseHandler: registerUseCaseHandler,
	}
}
