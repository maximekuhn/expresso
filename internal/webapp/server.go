package webapp

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/expresso/internal/webapp/handlers"
	"github.com/maximekuhn/expresso/internal/webapp/middleware"
)

type server struct {
	app *application
}

func (s *server) setup(l *slog.Logger) {
	handleAssets()

	requestIdMw := middleware.NewRequestIdMiddleware()
	loggerMw := middleware.NewLoggerMiddleware(l)
	loggedInMw := middleware.NewLoggedInMiddleware(
		s.app.authService,
		s.app.userService,
		s.app.sessionProvider,
	)
	chain := middleware.Chain(requestIdMw, loggerMw)
	loggedInChain := middleware.Chain(requestIdMw, loggerMw, loggedInMw)

	registerHandler := handlers.NewRegisterHandler(l, s.app.registerUsecaseHandler)
	http.Handle("/register", chain.Middleware(registerHandler))

	loginHandler := handlers.NewLoginHandler(l, s.app.loginUsecaseHandler)
	http.Handle("/login", chain.Middleware(loginHandler))

	indexHandler := handlers.NewIndexHandler(l)
	http.Handle("/", loggedInChain.Middleware(indexHandler))
}

func handleAssets() {
	fs := http.FileServer(http.Dir("internal/webapp/ui/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

}

func Run(db *sql.DB, l *slog.Logger) error {
	server := &server{app: newApplication(db)}
	server.setup(l)
	fmt.Println("server is running on port 5092")
	return http.ListenAndServe(":5092", nil)
}
