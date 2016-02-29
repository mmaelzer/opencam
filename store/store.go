package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/settings"
)

const (
	DATE_FORMAT = "2006-01-02"
	TIME_FORMAT = "15_04_05"
)

var sqldURL string
var cameraCache map[int]model.Camera

func Initialize(url string) {
	sqldURL = strings.TrimSuffix(url, "/")
	CacheCameras()
}

func CacheCameras() {
	var cameras []model.Camera
	err := Get("camera", &cameras)
	if err != nil {
		log.Printf("Unable to cache cameras: %s", err.Error())
		return
	}
	cameraCache = make(map[int]model.Camera, len(cameras))
	for i := range cameras {
		cameraCache[cameras[i].ID] = cameras[i]
	}
}

func FolderForBlock(block pipeline.Block) string {
	if len(block.Frames) == 0 {
		return ""
	}

	ff := block.Frames[0]
	lf := block.Frames[len(block.Frames)-1]
	duration := lf.Timestamp.Unix() - ff.Timestamp.Unix()
	folder := fmt.Sprintf(
		"%d-%s-%d",
		block.Camera.ID,
		ff.Timestamp.Format(TIME_FORMAT),
		duration,
	)

	dateDir := ff.Timestamp.Format(DATE_FORMAT)
	relativePath := path.Join(dateDir, folder)
	return relativePath
}

func FilenameForFrame(frame cam.Frame, ext string) string {
	return fmt.Sprintf("%d.%s", frame.Timestamp.UnixNano(), ext)
}

func SaveImageFile(data []byte, folder, filename string) error {
	base := settings.GetString("output")
	fullPath := path.Join(base, folder)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}
	file := path.Join(fullPath, filename)
	log.Printf("saving %s", file)
	return ioutil.WriteFile(file, data, 0755)
}

func Get(name string, i interface{}) error {
	return get(fmt.Sprintf("%s/%s", sqldURL, name), i)
}

func GetCameras() []model.Camera {
	cameras := make([]model.Camera, 0)
	for _, camera := range cameraCache {
		cameras = append(cameras, camera)
	}
	return cameras
}

func GetEvent(id string) (*model.Event, error) {
	var events []model.Event
	if err := get(fmt.Sprintf("%s/event/%s", sqldURL, id), &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return &model.Event{}, nil
	}
	return &events[0], nil
}

func GetID(name string, id string, i interface{}) error {
	return get(fmt.Sprintf("%s/%s/%s", sqldURL, name, id), i)
}

func Save(name string, data interface{}) ([]byte, error) {
	return post(fmt.Sprintf("%s/%s", sqldURL, name), data)
}

func Update(name string, data interface{}, id string) ([]byte, error) {
	return put(fmt.Sprintf("%s/%s/%s", sqldURL, name, id), data)
}

func Delete(name string, params map[string]string) ([]byte, error) {
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	url := fmt.Sprintf("%s/%s?%s", sqldURL, name, v.Encode())
	return send(url, nil, "DELETE")
}

func put(url string, i interface{}) ([]byte, error) {
	return send(url, i, "PUT")
}

func post(url string, i interface{}) ([]byte, error) {
	return send(url, i, "POST")
}

func send(url string, i interface{}, method string) ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Non-2xx response when attempting to perform db query: %d %s", res.StatusCode, string(body))
	}

	return body, nil
}

func get(url string, i interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.Unmarshal(body, &i)
}

func RawWrite(query string) ([]byte, error) {
	return rawRequest("write", query)
}

func RawRead(query string) ([]byte, error) {
	return rawRequest("read", query)
}

func rawRequest(t string, query string) ([]byte, error) {
	query = strings.Replace(query, "\n", " ", -1)
	query = strings.Replace(query, "\t", " ", -1)
	return post(sqldURL, map[string]string{t: query})
}
