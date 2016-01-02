package writer

import (
	"log"
	"path"
	"time"

	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/store"
)

func Writer(block pipeline.Block) {
	if len(block.Frames) == 0 {
		return
	}

	ff := block.Frames[0]
	lf := block.Frames[len(block.Frames)-1]
	duration := lf.Timestamp.Unix() - ff.Timestamp.Unix()
	folder := store.FolderForBlock(block)

	firstImage := path.Join(folder, store.FilenameForFrame(ff, "jpg"))
	lastImage := path.Join(folder, store.FilenameForFrame(lf, "jpg"))

	event := map[string]interface{}{
		"camera_id":   block.Camera.ID,
		"type":        "motion",
		"start_time":  ff.Timestamp.Format(time.RFC3339),
		"end_time":    lf.Timestamp.Format(time.RFC3339),
		"duration":    duration,
		"filepath":    folder,
		"first_frame": firstImage,
		"last_frame":  lastImage,
		"frame_count": len(block.Frames),
	}

	log.Printf("Saving event: %v", event)
	go store.Save("event", event)

	for _, frame := range block.Frames {
		file := store.FilenameForFrame(frame, "jpg")
		store.SaveImageFile(frame.Bytes, folder, file)
	}
}
