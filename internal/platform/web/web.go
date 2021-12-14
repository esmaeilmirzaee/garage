package web

import (
	"context"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

// In order to have access to the status code in the middleware
// we should attach the status code into the context

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how represent values or stored/retrieved.
const KeyValues ctxKey = 1

type Values struct {
	Start      time.Time
	StatusCode int
}

// ************************************************************

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

		// Attaching status code into the context
		v := Values{
			Start: time.Now(),
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyValues, &v)
		r = r.WithContext(ctx)

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
