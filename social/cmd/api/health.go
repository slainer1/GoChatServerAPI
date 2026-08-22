package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		health check
//	@Description	health check endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		//err
		app.statusInternalServerErrorHandler(w, r, err)
	}
}
