package manager

import (
	"log"
	"time"

	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/store"
	"github.com/mmaelzer/opencam/transform"
	"github.com/mmaelzer/opencam/writer"
)

var p *pipeline.Pipeline

func Start() {
	if p != nil {
		Stop()
	}

	var cameras []pipeline.Camera
	if err := store.Get("camera", &cameras); err != nil {
		log.Fatal("Unable to get cameras from database", err)
	}

	if len(cameras) == 0 {
		return
	}

	camps := make([]*pipeline.Camera, len(cameras))
	for i := range cameras {
		camps[i] = &cameras[i]
	}

	p = pipeline.New().
		AddCameras(camps).
		Pipe(transform.Motion()).
		Pipe(transform.Batch(time.Second*4, time.Second*30)).
		AddWriter(writer.Writer)

	p.Start()
}

func Cameras() []*pipeline.Camera {
	if p == nil {
		return []*pipeline.Camera{}
	}
	return p.Cameras()
}

func Stop() {
	p.Stop()
	p = nil
}
