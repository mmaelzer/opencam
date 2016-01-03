package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/camotion"
	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/settings"
	"github.com/mmaelzer/opencam/store"
)

type TransparentResponseWriter struct {
	Status int
	Writer http.ResponseWriter
	Size   int
}

func (t *TransparentResponseWriter) Write(b []byte) (int, error) {
	size, err := t.Writer.Write(b)
	t.Size = size
	return size, err
}

func (t *TransparentResponseWriter) Header() http.Header {
	return t.Writer.Header()
}

func (t *TransparentResponseWriter) WriteHeader(code int) {
	t.Status = code
	t.Writer.WriteHeader(code)
}

func sendImage(w http.ResponseWriter, image []byte) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(image)))
	w.WriteHeader(http.StatusOK)
	w.Write(image)
}

func sendJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

func sendText(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func getCamera(cameras []*pipeline.Camera, w http.ResponseWriter, r *http.Request) *pipeline.Camera {
	cIDStr := path.Base(r.URL.Path)
	cID, err := strconv.Atoi(cIDStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	if cID > len(cameras) {
		http.NotFound(w, r)
		return nil
	}

	return cameras[cID]
}

func stream(cameras []*pipeline.Camera) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		camera := getCamera(cameras, w, r)
		if camera == nil {
			return
		}

		log.Printf("[%s] Opening HTTP stream to client", camera.Name)

		boundary := "gocamera"
		w.Header().Set(
			"Content-Type",
			"multipart/x-mixed-replace; boundary="+boundary,
		)
		writer := multipart.NewWriter(w)
		writer.SetBoundary(boundary)

		cn := w.(http.CloseNotifier).CloseNotify()
		block := make(chan bool)

		frames, err := camera.Subscribe()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			<-cn
			log.Printf("[%s] Closing HTTP stream to client", camera.Name)
			block <- true
			camera.Unsubscribe(frames)
			writer.Close()
		}()

		go func() {
			for frame := range frames {
				mh := make(textproto.MIMEHeader)
				mh.Set("Content-Type", "image/jpeg")
				mh.Set("Content-Length", strconv.Itoa(len(frame.Bytes)))
				pw, err := writer.CreatePart(mh)

				if err != nil {
					writer.Close()
					return
				}
				log.Printf("[%s:%d] sending frame %d", frame.CameraName, frame.Number, len(frame.Bytes))
				pw.Write(frame.Bytes)
			}
		}()
		<-block
	}
}

func frame(cameras []*pipeline.Camera) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		camera := getCamera(cameras, w, r)
		if camera == nil {
			return
		}

		frames, err := camera.Subscribe()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer camera.Unsubscribe(frames)
		frame := <-frames
		sendImage(w, frame.Bytes)
	}
}

func config(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/config.html")
}

func blended(cameras []*pipeline.Camera) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		camera := getCamera(cameras, w, r)
		if camera == nil {
			return
		}

		frames, err := camera.Subscribe()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		frame1 := <-frames
		frame2 := <-frames
		camera.Unsubscribe(frames)

		blended := getBlendedImage(&frame1, &frame2)

		var b bytes.Buffer
		err = jpeg.Encode(&b, blended, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendImage(w, b.Bytes())
	}
}

func getBlendedImage(frame1, frame2 *cam.Frame) image.Image {
	jpg1, err := jpeg.Decode(bytes.NewReader(frame1.Bytes))
	if err != nil {
		return nil
	}
	jpg2, err := jpeg.Decode(bytes.NewReader(frame2.Bytes))
	if err != nil {
		return nil
	}
	return camotion.Blended(jpg1, jpg2, 2500)
}

func event(w http.ResponseWriter, r *http.Request) {
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

func l(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		t := &TransparentResponseWriter{Writer: w}
		fn(t, r)
		log.Printf("%s %d %s %s %d", r.Method, t.Status, r.URL.String(), time.Since(start), t.Size)
	}
}

func events() http.HandlerFunc {
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

func Serve(cameras []*pipeline.Camera) {
	static := http.FileServer(http.Dir("client/"))
	http.Handle("/", static)

	vpath := settings.GetString("output")
	videos := http.FileServer(http.Dir(vpath))
	http.Handle("/video/", http.StripPrefix("/video/", videos))
	http.HandleFunc("/api/event/", l(event))
	http.HandleFunc("/api/events", l(events()))
	http.HandleFunc("/stream/", stream(cameras))
	http.HandleFunc("/blended/", l(blended(cameras)))
	http.HandleFunc("/frame/", l(frame(cameras)))
	http.HandleFunc("/config", l(config))

	port := settings.GetInt("http.port")
	log.Printf("http listening on %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
