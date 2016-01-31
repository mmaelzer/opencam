package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mmaelzer/opencam/model"
	"github.com/mmaelzer/opencam/store"

	"github.com/mmaelzer/opencam/pipeline/manager"
)

func Cameras() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			camera, err := saveCamera(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				sendJSON(w, camera)
				go func() {
					store.CacheCameras()
					manager.Restart()
				}()
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
		_, err = store.Update("camera", &camera, fmt.Sprintf("%d", camera.ID))
		if err != nil {
			return nil, err
		}
		return &camera, nil
	} else {
		res, err = store.Save("camera", struct {
			MinChange int    `json:"min_change"`
			Name      string `json:"name"`
			Password  string `json:"password"`
			Threshold int    `json:"threshold"`
			URL       string `json:"url"`
			Username  string `json:"username"`
		}{
			camera.MinChange,
			camera.Name,
			camera.Password,
			camera.Threshold,
			camera.URL,
			camera.Username,
		})
	}

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &camera)
	return &camera, err
}
