package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/mmaelzer/opencam/settings"
	"github.com/mmaelzer/opencam/store"
)

func Event() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/event/"):]
		ev, err := store.GetEvent(id)
		if err != nil {
			sendText(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if ev.ID == 0 {
			sendText(w, http.StatusNotFound, "Not Found")
			return
		}
		base := settings.GetString("output")
		dir := path.Join(base, ev.Filepath)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			sendText(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		frames := make([]string, 0)
		for i := range files {
			if files[i].IsDir() {
				continue
			}
			frames = append(frames, path.Join("/video", ev.Filepath, files[i].Name()))
		}
		ev.Frames = frames
		sendJSON(w, ev)
	}
}

func Events() http.HandlerFunc {
	sqldURL, err := url.Parse(settings.GetString("sqldURL"))
	if err != nil {
		log.Fatal(err)
	}

	sqldURL.Path = path.Join(sqldURL.Path, "/event")
	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = sqldURL
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t := &TransparentResponseWriter{Writer: w}
		proxy.ServeHTTP(t, r)
	}
}
