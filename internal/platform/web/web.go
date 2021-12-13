package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

// App is the entrypoint for all web applications
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

// NewApp knows how to construct for an App.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}
}

// Handle connects a method and URL pattern to a particular application handler.
func (a *App) Handle(method, pattern string, h Handler) {

	// Everytime a request comes from the routes would be
	// wrapped via middleware.
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			a.log.Printf("Unhandled Errors: %v", err)
		}
	}

	a.mux.MethodFunc(method, pattern, fn)
}

// ServeHttp handles http service
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
