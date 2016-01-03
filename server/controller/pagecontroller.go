package controller

import (
	"net/http"

	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/store"
	"github.com/mmaelzer/opencam/templater"
)

func EventPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html, err := templater.Build("events", nil)
		if err != nil {
			sendText(w, http.StatusInternalServerError, "Unable to render page")
		}
		sendHTML(w, html)
	}
}

func LivePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cameras []model.Camera
		err := store.Get("camera", &cameras)
		html, err := templater.Build("live", cameras)
		if err != nil {
			sendText(w, http.StatusInternalServerError, "Unable to render page")
		}
		sendHTML(w, html)
	}
}
