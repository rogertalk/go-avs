package avs

import (
	"time"
)

// newContext creates a Message suited for being used as a context value.
func newContext(namespace, name string) *Message {
	return &Message{
		Header: map[string]string{
			"namespace": namespace,
			"name":      name,
		},
		Payload: nil,
	}
}

/********** Alerts **********/

// The AlertsState context.
type AlertsState struct {
	*Message
	Payload struct {
		AllAlerts    []Alert `json:"allAlerts"`
		ActiveAlerts []Alert `json:"activeAlerts"`
	} `json:"payload"`
}

func NewAlertsState(allAlerts, activeAlerts []Alert) *AlertsState {
	m := new(AlertsState)
	m.Message = newContext("Alerts", "AlertsState")
	m.Payload.AllAlerts = allAlerts
	m.Payload.ActiveAlerts = activeAlerts
	return m
}

/********** AudioPlayer **********/

// The PlaybackState context.
type PlaybackState struct {
	*Message
	Payload playbackState `json:"payload"`
}

func NewPlaybackState(token string, offset time.Duration, activity PlayerActivity) *PlaybackState {
	m := new(PlaybackState)
	m.Message = newContext("AudioPlayer", "PlaybackState")
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	m.Payload.PlayerActivity = activity
	m.Payload.Token = token
	return m
}

/********** Speaker **********/

// The VolumeState context.
type VolumeState struct {
	*Message
	Payload struct {
		Volume int  `json:"volume"`
		Muted  bool `json:"muted"`
	} `json:"payload"`
}

func NewVolumeState(volume int, muted bool) *VolumeState {
	m := new(VolumeState)
	m.Message = newContext("Speaker", "VolumeState")
	m.Payload.Volume = volume
	m.Payload.Muted = muted
	return m
}

/********** SpeechSynthesizer **********/

// The SpeechState context.
type SpeechState struct {
	*Message
	Payload struct {
		Token                string         `json:"token"`
		OffsetInMilliseconds int            `json:"offsetInMilliseconds"`
		PlayerActivity       PlayerActivity `json:"playerActivity"`
	} `json:"payload"`
}

func NewSpeechState(token string, offset time.Duration, playerActivity PlayerActivity) *SpeechState {
	m := new(SpeechState)
	m.Message = newContext("SpeechSynthesizer", "SpeechState")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	m.Payload.PlayerActivity = playerActivity
	return m
}
