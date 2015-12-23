package client

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"
	"strconv"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/camotion"
	"github.com/mmaelzer/opencam/settings"
)

func getCamera(cameras []*cam.Camera, w http.ResponseWriter, r *http.Request) *cam.Camera {
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

func stream(cameras []*cam.Camera) http.HandlerFunc {
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

func frame(cameras []*cam.Camera) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(frame.Bytes)))
		w.Write(frame.Bytes)
	}
}

func config(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/config.html")
}

func blended(cameras []*cam.Camera) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
		w.Write(b.Bytes())
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

func Serve(cameras []*cam.Camera) {
	static := http.FileServer(http.Dir("client/"))
	http.Handle("/", static)

	vpath := settings.GetString("output")
	videos := http.FileServer(http.Dir(vpath))
	http.Handle("/video/", http.StripPrefix("/video/", videos))

	http.HandleFunc("/stream/", stream(cameras))
	http.HandleFunc("/blended/", blended(cameras))
	http.HandleFunc("/frame/", frame(cameras))
	http.HandleFunc("/config", config)

	port := settings.GetInt("http.port")
	log.Printf("http listening on %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
