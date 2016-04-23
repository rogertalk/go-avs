package avs

import (
	"strings"
	"time"
)

// A streamable audio item.
type AudioItem struct {
	AudioItemId string `json:"audioItemId"`
	Stream      Stream `json:"stream"`
}

// How the play queue should be cleared.
type ClearBehavior string

// Possible values for ClearBehavior.
// TODO: Complete these constants.
const (
	ClearBehaviorClearEnqueued = ClearBehavior("CLEAR_ENQUEUED")
)

// How an audio item should be inserted into the play queue.
type PlayBehavior string

// Possible values for PlayBehavior.
// TODO: Complete these constants.
const (
	PlayBehaviorEnqueue    = PlayBehavior("ENQUEUE")
	PlayBehaviorReplaceAll = PlayBehavior("REPLACE_ALL")
)

// A value for what state the device player is in.
type PlayerActivity string

// Possible values for PlayerActivity.
// TODO: Complete these constants.
const (
	PlayerActivityPlaying  = PlayerActivity("PLAYING")
	PlayerActivityFinished = PlayerActivity("FINISHED")
)

type ProgressReport struct {
	ProgressReportIntervalInMilliseconds float64 `json:"progressReportIntervalInMilliseconds"`
}

func (p *ProgressReport) Interval() time.Duration {
	return time.Duration(p.ProgressReportIntervalInMilliseconds) * time.Millisecond
}

type Stream struct {
	ExpiryTime           string         `json:"expiryTime"`
	OffsetInMilliseconds float64        `json:"offsetInMilliseconds"`
	ProgressReport       ProgressReport `json:"progressReport"`
	Token                string         `json:"token"`
	URL                  string         `json:"url"`
}

func (s *Stream) ContentId() string {
	if !strings.HasPrefix(s.URL, "cid:") {
		return ""
	}
	return s.URL[4:]
}
