package pipeline

import "github.com/mmaelzer/cam"

type Block struct {
	Camera *Camera
	Frames []cam.Frame
}

type StreamFunc func(chan Block) chan Block
type WriterFunc func(Block)

type Camera struct {
	ID uint64 `json:"id"`
	cam.Camera
}

type Pipeline struct {
	streams []StreamFunc
	cameras []*Camera
	writers []WriterFunc
}

func New() *Pipeline {
	p := &Pipeline{}
	return p.init()
}

func (p *Pipeline) Cameras() []*Camera {
	return p.cameras
}

func (p *Pipeline) init() *Pipeline {
	p.streams = make([]StreamFunc, 0)
	p.cameras = make([]*Camera, 0)
	p.writers = make([]WriterFunc, 0)
	return p
}

func (p *Pipeline) AddCamera(camera *Camera) *Pipeline {
	p.cameras = append(p.cameras, camera)
	return p
}

func (p *Pipeline) AddCameras(cameras []*Camera) *Pipeline {
	p.cameras = append(p.cameras, cameras...)
	return p
}

func (p *Pipeline) AddWriter(w WriterFunc) *Pipeline {
	p.writers = append(p.writers, w)
	return p
}

func (p *Pipeline) Start() []error {
	errors := make([]error, 0)
	for i := range p.cameras {
		c := p.cameras[i]
		frames, err := initCamera(c)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		for j := 0; j < len(p.streams); j++ {
			frames = p.streams[j](frames)
		}

		go p.writeFrames(frames)
	}

	if len(errors) > 0 {
		return errors
	} else {
		return nil
	}
}

func (p *Pipeline) Stop() {
	for i := range p.cameras {
		p.cameras[i].Stop()
	}
	p.init()
}

func (p *Pipeline) Pipe(t StreamFunc) *Pipeline {
	p.streams = append(p.streams, t)
	return p
}

func (p *Pipeline) writeFrames(frames chan Block) {
	for f := range frames {
		for _, w := range p.writers {
			w(f)
		}
	}
}

func initCamera(camera *Camera) (chan Block, error) {
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
