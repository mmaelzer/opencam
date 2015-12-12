package transform

import (
	"bytes"
	"image/jpeg"
	"log"
	"sync"
	"time"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/camotion"
	"github.com/mmaelzer/opencam/pipeline"
)

type Lockit struct {
	mutex    sync.Mutex
	Unlocked bool
}

func (l *Lockit) Lock() {
	l.set(false)
}

func (l *Lockit) Unlock() {
	l.set(true)
}

func (l *Lockit) set(lock bool) {
	l.mutex.Lock()
	l.Unlocked = lock
	l.mutex.Unlock()
}

func Motion() pipeline.TransformFunc {
	return func(in chan pipeline.Block) chan pipeline.Block {
		out := make(chan pipeline.Block)
		go readFrames(in, out)
		return out
	}
}

func readFrames(in chan pipeline.Block, out chan pipeline.Block) {
	var fellBehind uint64
	lock := &Lockit{}
	lock.Unlock()

	fcache := make([]cam.Frame, 0)
	lastTest := time.Now()

	for block := range in {
		for _, frame := range block.Frames {
			fcache = append(fcache, frame)
			if len(fcache) > 1 && time.Since(lastTest) > time.Second*2 {
				if lock.Unlocked {
					lock.Lock()
					go processFrames(fcache, block.Camera, out, lock)
				} else {
					fellBehind++
					log.Printf("Motion processing fell behind. Occurence %d", fellBehind)
				}

				lastTest = time.Now()
				fcache = make([]cam.Frame, 0)
			}
		}
	}
}

func processFrames(frames []cam.Frame, camera *cam.Camera, out chan pipeline.Block, lock sync.Locker) {
	firstFrame := frames[0]
	lastFrame := frames[len(frames)-1]
	motion := testFrames(&firstFrame, &lastFrame)

	if motion {
		log.Printf(
			"Event found %s - %s",
			firstFrame.Timestamp.Format(time.RFC3339),
			lastFrame.Timestamp.Format(time.RFC3339),
		)

		out <- pipeline.Block{
			Camera: camera,
			Frames: frames,
		}
	}

	lock.Unlock()
}

func testFrames(frame1, frame2 *cam.Frame) bool {
	jpg1, err := jpeg.Decode(bytes.NewReader(frame1.Bytes))
	if err != nil {
		return false
	}
	jpg2, err := jpeg.Decode(bytes.NewReader(frame2.Bytes))
	if err != nil {
		return false
	}
	return camotion.MotionWithStep(jpg1, jpg2, 6, 2500, 10)
}
