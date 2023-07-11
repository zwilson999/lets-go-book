package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// middleware function to set security headers
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// method to handle panics and recover with proper error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// deferred func will always run in the event of a panic as Go unwinds the call stack
		defer func() {

			// use builtin recover func to cehck if there has been a panic or not
			// if there has, set appropriate headers and send server error
			if err := recover(); err != nil {

				// set connection: close header on the response
				w.Header().Set("Connection", "close")

				// call our application's serverError() helper to return a 500 status
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// method to conditionally show user login page (if not authenticated)
// and refrain from caching items in the browser cache
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers
		// in the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// otherwise set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or other intermediary)
		w.Header().Add("Cache-Control", "no-store")

		// and call next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// create a middleware func which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve the authenticatedUserID value from the session
		// by using the GetInt() method. this will return zero value for an int (0)
		// if no "authenticatedUserID" value is in the session, in which case we
		// call the next handler in the chain as normal and return
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// otherwise we check to see if a user with that ID exists in our database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// if a matching user is found, we know that the request is coming from an authenticated user whom
		// exists in our database. we create a new copy of the request (with an isAuthenticatedContextKey)
		// value of true in the request context) and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		// call the next handler in the middleware chain
		next.ServeHTTP(w, r)
	})
}
