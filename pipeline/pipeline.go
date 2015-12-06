package pipeline

import (
	"fmt"

	"github.com/mmaelzer/cam"
)

type Block struct {
	Camera *cam.Camera
	Frames []cam.Frame
}

type TransformFunc func(chan Block) chan Block

type Pipeline struct {
	transforms []TransformFunc
	cameras    []cam.Camera
	writer     func(Block)
}

func New() *Pipeline {
	transforms := make([]TransformFunc, 0)
	cameras := make([]cam.Camera, 0)
	writer := func(frames Block) {}
	return &Pipeline{
		transforms: transforms,
		cameras:    cameras,
		writer:     writer,
	}
}

func (p *Pipeline) AddCamera(camera *cam.Camera) *Pipeline {
	p.cameras = append(p.cameras, *camera)
	return p
}

func (p *Pipeline) AddWriter(w func(Block)) *Pipeline {
	p.writer = w
	return p
}

func (p *Pipeline) Start() []error {
	errors := make([]error, 0)
	for i := 0; i < len(p.cameras); i++ {
		c := p.cameras[i]
		frames, err := initCamera(&c)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		for j := 0; j < len(p.transforms); j++ {
			frames = p.transforms[j](frames)
		}

		go p.writeFrames(frames)
	}

	if len(errors) > 0 {
		return errors
	} else {
		return nil
	}
}

func (p *Pipeline) AddTransform(t TransformFunc) *Pipeline {
	p.transforms = append(p.transforms, t)
	return p
}

func (p *Pipeline) writeFrames(frames chan Block) {
	fmt.Println("write frames")
	for f := range frames {
		p.writer(f)
	}
}

func initCamera(camera *cam.Camera) (chan Block, error) {
	frames, err := camera.Subscribe()
	if err != nil {
		return nil, err
	}
	c := make(chan Block)
	go func() {
		for frame := range frames {
			c <- Block{
				Camera: camera,
				Frames: []cam.Frame{frame},
			}
		}
	}()
	return c, nil
}
