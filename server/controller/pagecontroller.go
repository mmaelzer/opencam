package controller

import (
	"net/http"

	"github.com/mmaelzer/opencam/templater"
)

func Main() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html, err := templater.Build("events", nil)
		if err != nil {
			sendText(w, http.StatusInternalServerError, "Unable to render page")
		}
		sendHTML(w, html)
	}
}
