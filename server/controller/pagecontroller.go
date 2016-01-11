package controller

import (
	"net/http"

	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/store"
	"github.com/mmaelzer/opencam/templater"
)

func EventPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cameras := store.GetCameras()
		html, err := templater.Build("events", struct {
			Cameras []model.Camera
		}{
			cameras,
		})
		if err != nil {
			unableToRender(w)
		} else {
			sendHTML(w, html)
		}
	}
}

func unableToRender(w http.ResponseWriter) {
	sendText(w, http.StatusInternalServerError, "Unable to render page")
}

func LivePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cameras []model.Camera
		err := store.Get("camera", &cameras)
		html, err := templater.Build("live", cameras)
		if err != nil {
			unableToRender(w)
		} else {
			sendHTML(w, html)
		}
	}
}
