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
	handleAssets()

	l := logger.Setup()

	requestIdMw := middleware.NewRequestIdMiddleware()
	loggerMw := middleware.NewLoggerMiddleware(l)
	chain := middleware.Chain(requestIdMw, loggerMw)

	registerHandler := handlers.NewRegisterHandler(l, s.app.registerUsecaseHandler)
	http.Handle("/register", chain.Middleware(registerHandler))

	loginHandler := handlers.NewLoginHandler(l, s.app.loginUsecaseHandler)
	http.Handle("/login", chain.Middleware(loginHandler))

}

func handleAssets() {
	fs := http.FileServer(http.Dir("internal/webapp/ui/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

}

func Run(db *sql.DB) error {
	server := &server{app: newApplication(db)}
	server.setup()
	fmt.Println("server is running on port 5092")
	return http.ListenAndServe(":5092", nil)
}
