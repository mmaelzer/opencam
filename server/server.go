package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/server/controller"
	"github.com/mmaelzer/opencam/settings"
)

func config(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/config.html")
}

func logHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		t := &controller.TransparentResponseWriter{Writer: w}
		h.ServeHTTP(t, r)
		log.Printf("%s %d %s %s %d", r.Method, t.Status, r.URL.String(), time.Since(start), t.Size)
	}
}

func Serve(cameras []*pipeline.Camera) {
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("client/")))
	vpath := settings.GetString("output")
	videos := http.StripPrefix("/video/", http.FileServer(http.Dir(vpath)))

	http.HandleFunc("/", logHandler(controller.LivePage()))
	http.HandleFunc("/events", logHandler(controller.EventPage()))
	http.HandleFunc("/live", logHandler(controller.LivePage()))
	http.HandleFunc("/config", logHandler(controller.ConfigPage()))
	http.HandleFunc("/static/", logHandler(static))
	http.HandleFunc("/video/", logHandler(videos))
	http.HandleFunc("/api/event/", logHandler(controller.Event()))
	http.HandleFunc("/api/events", logHandler(controller.Events()))
	http.HandleFunc("/api/camera", logHandler(controller.Cameras()))
	http.HandleFunc("/stream/", controller.Stream(cameras))
	http.HandleFunc("/blended/", logHandler(controller.Blended(cameras)))
	http.HandleFunc("/frame/", logHandler(controller.Frame(cameras)))

	port := settings.GetInt("http.port")
	log.Printf("http listening on %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
