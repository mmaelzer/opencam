package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
