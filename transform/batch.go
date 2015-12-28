package transform

import (
	"time"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/opencam/pipeline"
)

func Batch(wait, max time.Duration) pipeline.StreamFunc {
	return func(in chan pipeline.Block) chan pipeline.Block {
		out := make(chan pipeline.Block)
		go batchBlocks(wait, max, in, out)
		return out
	}
}

func batchBlocks(wait, max time.Duration, in, out chan pipeline.Block) {
	frames := make(map[*cam.Camera][]cam.Frame)
	ticks := make(map[*cam.Camera]time.Time)
	starts := make(map[*cam.Camera]time.Time)

	go func() {
		for {
			for camera, f := range frames {
				if len(f) == 0 {
					continue
				}
				tick := ticks[camera]
				start := starts[camera]
				if time.Since(tick) > wait || time.Since(start) > max {
					out <- pipeline.Block{camera, f}
					frames[camera] = make([]cam.Frame, 0)
					ticks[camera] = time.Now()
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	for block := range in {
		if len(frames[block.Camera]) == 0 {
			starts[block.Camera] = time.Now()
		}
		frames[block.Camera] = append(frames[block.Camera], block.Frames...)
		ticks[block.Camera] = time.Now()
	}
}
