package controller

import (
	"net/http"

	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/store"
	"github.com/mmaelzer/opencam/templater"
)

func unableToRender(w http.ResponseWriter) {
	sendText(w, http.StatusInternalServerError, "Unable to render page")
}

func ConfigPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cameras := store.GetCameras()
		html, err := templater.Build("config", struct {
			Cameras []model.Camera
			Page    string
		}{
			cameras,
			"config",
		})
		if err != nil {
			unableToRender(w)
		} else {
			sendHTML(w, html)
		}
	}
}

func EventPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cameras := store.GetCameras()
		html, err := templater.Build("events", struct {
			Cameras []model.Camera
			Page    string
		}{
			cameras,
			"events",
		})
		if err != nil {
			unableToRender(w)
		} else {
			sendHTML(w, html)
		}
	}
}

func LivePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cameras []model.Camera
		err := store.Get("camera", &cameras)
		html, err := templater.Build("live", struct {
			Cameras []model.Camera
			Page    string
		}{
			cameras,
			"live",
		})
		if err != nil {
			unableToRender(w)
		} else {
			sendHTML(w, html)
		}
	}
}
