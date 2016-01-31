package controller

import (
	"bytes"
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
	"github.com/mmaelzer/opencam/pipeline"
)

func getCamera(cameras []*pipeline.Camera, w http.ResponseWriter, r *http.Request) *pipeline.Camera {
	cIDStr := path.Base(r.URL.Path)
	cID, err := strconv.Atoi(cIDStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	id := uint64(cID)
	for i := range cameras {
		if cameras[i].ID == id {
			return cameras[i]
		}
	}

	http.NotFound(w, r)
	return nil
}

func Stream(cameras []*pipeline.Camera) http.HandlerFunc {
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

func Frame(cameras []*pipeline.Camera) http.HandlerFunc {
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

func Blended(cameras []*pipeline.Camera) http.HandlerFunc {
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

		blended := getBlendedImage(&frame1, &frame2, camera)

		var b bytes.Buffer
		err = jpeg.Encode(&b, blended, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendImage(w, b.Bytes())
	}
}

func getBlendedImage(frame1, frame2 *cam.Frame, camera *pipeline.Camera) image.Image {
	jpg1, err := jpeg.Decode(bytes.NewReader(frame1.Bytes))
	if err != nil {
		return nil
	}
	jpg2, err := jpeg.Decode(bytes.NewReader(frame2.Bytes))
	if err != nil {
		return nil
	}
	return camotion.Blended(jpg1, jpg2, camera.Threshold)
}
