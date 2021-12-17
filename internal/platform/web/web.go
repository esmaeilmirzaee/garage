package web

import (
	"context"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/trace"

	// Distributed tracing
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
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

// Handler is our customized handler
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
	och *ochttp.Handler // Distributed tracing
}

// NewApp constructs an App to handle a set of routes. Any Middleware provided
// will be ran for every request.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	app := App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}

	// Create an Opencensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if a client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}

// Handle connects a method and URL pattern to a particular application handler.
// To handle authorization mechanism a new argument which is a middleware should
// be part of the following handler
func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {

	// First wrap specific middleware around this handler
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	// Everytime a request comes from the routes would be
	// wrapped via middleware.
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Attaching status code into the context
		v := Values{
			Start: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		// Trace the application
		ctx, span := trace.StartSpan(ctx, "internal.platform.web")
		// End for the span could be called immediately or in the line 76
		// (after the if) but the best practice is
		defer span.End()

		if err := h(ctx, w, r); err != nil {
			a.log.Printf("Unhandled Errors: %v", err)
		}
	}

	a.mux.MethodFunc(method, pattern, fn)
}

// ServeHttp handles http service
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}
