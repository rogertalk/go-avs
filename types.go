//
package avs

import (
	"strings"
	"time"
)

// Alert type.
type AlertType string

// Possible values for AlertType, used by the SetAlert directive.
const (
	// Alarms are scheduled for specific times (e.g., to wake the user up).
	AlertTypeAlarm = AlertType("ALARM")
	// Timers count down a certain amount of time (e.g., "timer for 5 minutes").
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
	// Clears all items in the queue, including the current one.
	ClearBehaviorClearAll = ClearBehavior("CLEAR_ALL")
	// Clears all queued items in the queue (not including the current one).
	ClearBehaviorClearEnqueued = ClearBehavior("CLEAR_ENQUEUED")
)

// How an audio item should be inserted into the play queue.
type PlayBehavior string

// Possible values for PlayBehavior.
const (
	// Play after current queue of audio.
	PlayBehaviorEnqueue = PlayBehavior("ENQUEUE")
	// Play immediately, throwing away all queued audio.
	PlayBehaviorReplaceAll = PlayBehavior("REPLACE_ALL")
	// Play after the currently playing audio, replacing all queued audio.
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

// An audio stream which can either be attached with the response or a remote URL.
type Stream struct {
	ExpiryTime           string         `json:"expiryTime"`
	OffsetInMilliseconds float64        `json:"offsetInMilliseconds"`
	ProgressReport       ProgressReport `json:"progressReport"`
	Token                string         `json:"token"`
	URL                  string         `json:"url"`
}

// If the audio is attached with the response, this returns the content id of that audio.
func (s *Stream) ContentId() string {
	if !strings.HasPrefix(s.URL, "cid:") {
		return ""
	}
	return s.URL[4:]
}
