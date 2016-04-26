package avs

import (
	"strings"
	"time"
)

// Alert type.
type AlertType string

// Possible values for AlertType.
const (
	AlertTypeAlarm = AlertType("ALARM")
	AlertTypeTimer = AlertType("TIMER")
)

// A streamable audio item.
type AudioItem struct {
	AudioItemId string `json:"audioItemId"`
	Stream      Stream `json:"stream"`
}

// How the play queue should be cleared.
type ClearBehavior string

// Possible values for ClearBehavior.
const (
	ClearBehaviorClearAll      = ClearBehavior("CLEAR_ALL")
	ClearBehaviorClearEnqueued = ClearBehavior("CLEAR_ENQUEUED")
)

// How an audio item should be inserted into the play queue.
type PlayBehavior string

// Possible values for PlayBehavior.
const (
	PlayBehaviorEnqueue         = PlayBehavior("ENQUEUE")
	PlayBehaviorReplaceAll      = PlayBehavior("REPLACE_ALL")
	PlayBehaviorReplaceEnqueued = PlayBehavior("REPLACE_ENQUEUED")
)

// A value for what state the device player is in.
type PlayerActivity string

// Possible values for PlayerActivity.
const (
	PlayerActivityBufferUnderrun = PlayerActivity("BUFFER_UNDERRUN")
	PlayerActivityIdle           = PlayerActivity("IDLE")
	PlayerActivityPaused         = PlayerActivity("PAUSED")
	PlayerActivityPlaying        = PlayerActivity("PLAYING")
	PlayerActivityFinished       = PlayerActivity("FINISHED")
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
