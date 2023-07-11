package main

import (
	"net/http"

	"snippetbox.lets-go/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// this method initializes our servemux with our routes for the web application
func (app *application) routes() http.Handler {

	// initialize our router
	router := httprouter.New()

	// handler for 404s
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to http.FileServer()
	// func to create the file server handler
	fileServer := http.FileServer(http.FS(ui.Files))

	// our static files are contained in the "static" folder of the ui.Files
	// embedded file system. So, for example, our CSS stylesheet is located at "static/css/main.css". This means
	// we no longer need to strip the prefix from the request URL. any requests that start with /static/ can
	// just be passed directly to the file server and the correspondign static file will be served (as long as it exists)
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// unprotected app routes use the "dynamic" middleware chain
	dynamic := alice.New(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	// signup
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))

	// login
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// protected (authenticated-only) app routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)

	// snippet create
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// logout
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// create a middleware chain containing the standard middleware which will be used for
	// every request that our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
