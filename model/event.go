package model

type Event struct {
	ID         int      `json:"id"`
	CameraID   uint64   `json:"camera_id"`
	Type       string   `json:"type"`
	StartTime  string   `json:"start_time"`
	EndTime    string   `json:"end_time"`
	Duration   int64    `json:"duration"`
	Filepath   string   `json:"filepath"`
	FirstFrame string   `json:"first_frame"`
	LastFrame  string   `json:"last_frame"`
	FrameCount int      `json:"frame_count"`
	Frames     []string `json:"frames"`
}
