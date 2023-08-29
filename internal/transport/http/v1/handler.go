package handler

import (
	"item-service/internal/app"
	"time"
)

type HttpHandler struct {
	app     app.Application
	timeOut time.Duration
}

func New(app app.Application) HttpHandler {
	return HttpHandler{
		app:     app,
		timeOut: 15 * time.Second,
	}
}
