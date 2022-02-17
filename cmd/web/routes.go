package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/ezratameno/lets_go/middleware"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {

	metricsMiddleware := middleware.NewMetricsMiddleware()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(metricsMiddleware.Metrics, app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our *dynamic application routes*. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddelware := alice.New(app.session.Enable)

	// Swap the route declarations to use the application struct's methods as the
	// handler functions.
	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux := pat.New()
	mux.Get("/", dynamicMiddelware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddelware.ThenFunc(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamicMiddelware.ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddelware.ThenFunc(http.HandlerFunc(app.showSnippet)))
	mux.Get("/metrics", promhttp.Handler())
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
