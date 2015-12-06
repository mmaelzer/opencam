package client

import (
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/opencam/settings"
)

func stream(camera *cam.Camera) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		boundary := "gocamera"
		w.Header().Set(
			"Content-Type",
			"multipart/x-mixed-replace; boundary="+boundary,
		)
		writer := multipart.NewWriter(w)
		writer.SetBoundary(boundary)

		closed := false

		cn := w.(http.CloseNotifier).CloseNotify()

		go func() {
			<-cn
			closed = true
		}()

		frames, err := camera.Subscribe()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for frame := range frames {
			if closed {
				writer.Close()
				return
			}

			mh := make(textproto.MIMEHeader)
			mh.Set("Content-Type", "image/jpeg")
			mh.Set("Content-Length", strconv.Itoa(len(frame.Bytes)))
			pw, err := writer.CreatePart(mh)

			if err != nil {
				writer.Close()
				return
			}
			pw.Write(frame.Bytes)
		}
	}
}

func Serve(camera *cam.Camera) {
	static := http.FileServer(http.Dir("client/"))
	http.Handle("/", static)

	vpath := settings.GetString("output")
	videos := http.FileServer(http.Dir(vpath))
	http.Handle("/video/", http.StripPrefix("/video/", videos))

	http.HandleFunc("/stream", stream(camera))
	panic(http.ListenAndServe(":8080", nil))
}
