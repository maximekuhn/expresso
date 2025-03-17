package webapp

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/maximekuhn/expresso/internal/logger"
	"github.com/maximekuhn/expresso/internal/webapp/handlers"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
)

type server struct {
	app *application
}

func (s *server) setup() {
	l := logger.Setup()

	requestIdMw := middleware.NewRequestIdMiddleware()
	loggerMw := middleware.NewLoggerMiddleware(l)
	chain := middleware.Chain(requestIdMw, loggerMw)

	registerHandler := handlers.NewRegisterHandler(l, s.app.registerUsecaseHandler)
	http.Handle("/register", chain.Middleware(registerHandler))
}

func Run(db *sql.DB) error {
	server := &server{app: newApplication(db)}
	server.setup()
	fmt.Println("server is running on port 5092")
	return http.ListenAndServe(":5092", nil)
}
