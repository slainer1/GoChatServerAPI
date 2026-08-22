package main

import (
	"log"
	"net/http"
)

func (app *application) statusInternalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s err: %v", r.Method, r.URL.Path, err)
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "The server encountered an problem")
}

func (app *application) badRequestResponsee(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s err: %v", r.Method, r.URL.Path, err)
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s err: %v", r.Method, r.URL.Path, err)
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict Error: %s path: %s err: %v", r.Method, r.URL.Path, err)
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusConflict, "resource not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden response", "method", r.Method, "path", r.URL.Path, "error")
	writeJSONError(w, http.StatusForbidden, "Forbidden")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("Rate Limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("retry after", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after "+retryAfter)
}
