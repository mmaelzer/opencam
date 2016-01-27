package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/store"
)

func Cameras() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			camera, err := saveCamera(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				sendJSON(w, camera)
			}
		}
	}
}

func saveCamera(r *http.Request) (*model.Camera, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var camera model.Camera
	err = json.Unmarshal(body, &camera)
	if err != nil {
		return nil, err
	}
	var res []byte
	if camera.ID > 0 {
		_, err = store.Update("camera", &camera)
		if err != nil {
			return nil, err
		}
		return &camera, nil
	} else {
		res, err = store.Save("camera", &camera)
	}

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &camera)
	return &camera, err
}
