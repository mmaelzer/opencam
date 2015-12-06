package transform

import (
	"bytes"
	"image/jpeg"
	"log"
	"time"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/camotion"
	"github.com/mmaelzer/opencam/pipeline"
)

func Motion() pipeline.TransformFunc {
	return func(blocks chan pipeline.Block) chan pipeline.Block {
		return readFrames(blocks)
	}
}

func readFrames(blocks chan pipeline.Block) chan pipeline.Block {
	fcache := make([]cam.Frame, 0)

	lastTest := time.Now()
	lastMotion := time.Now()
	startMotion := time.Now()
	inMotion := false

	out := make(chan pipeline.Block)

	var lastFrame cam.Frame
	go func() {
		for block := range blocks {
			for _, frame := range block.Frames {
				if len(fcache) > 0 && time.Since(lastTest) > time.Second*1 {
					motion := testFrames(&lastFrame, &frame)
					if motion {
						lastMotion = time.Now()
						if !inMotion {
							log.Printf("Event begin %s", startMotion.Format("15_04_05"))
							startMotion = time.Now()
							inMotion = true
						}
					}

					// log.Printf("time since last motion: %s", time.Since(lastMotion))
					// log.Printf("cache size: %d\n", len(fcache))
					if time.Since(lastMotion) < time.Second*4 {
						out <- pipeline.Block{
							Camera: block.Camera,
							Frames: fcache,
						}
					} else {
						if inMotion {
							log.Printf("Event end. Duration: %s", time.Since(startMotion).String())
						}

						inMotion = false
					}

					fcache = make([]cam.Frame, 0)
					lastTest = time.Now()
				}
				fcache = append(fcache, frame)
				lastFrame = frame
			}
		}
	}()
	return out
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
	return camotion.MotionWithStep(jpg1, jpg2, 4, 2500, 10)
}
