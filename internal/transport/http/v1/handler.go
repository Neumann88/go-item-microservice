package handler

import (
	"item-service/internal/app"
)

type HttpHandler struct {
	app app.Application
}

func New(app app.Application) HttpHandler {
	return HttpHandler{
		app: app,
	}
}
