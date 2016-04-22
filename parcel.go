package avs

import (
	"fmt"
	"strings"
	"time"
)

type TypedParcel interface {
	GetParcel() *Parcel
	Typed() TypedParcel
}

// A general structure for contexts, events and directives.
type Parcel struct {
	Header  map[string]string      `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

func (p *Parcel) GetParcel() *Parcel {
	return p
}

func (p *Parcel) String() string {
	return fmt.Sprintf("%s.%s", p.Header["namespace"], p.Header["name"])
}

// Returns a more specific type for this context, event or directive.
func (p *Parcel) Typed() TypedParcel {
	switch p.String() {
	case "AudioPlayer.PlaybackState":
		return &PlaybackState{p}
	case "SpeechRecognizer.ExpectSpeech":
		return &ExpectSpeech{p}
	case "SpeechRecognizer.Recognize":
		return &Recognize{p}
	case "SpeechSynthesizer.Speak":
		return &Speak{p}
	default:
		return p
	}
}

// The Recognize event.
type Recognize struct {
	*Parcel
}

func NewRecognize(messageId, dialogRequestId string) *Recognize {
	parcel := &Parcel{
		Header: map[string]string{
			"namespace":       "SpeechRecognizer",
			"name":            "Recognize",
			"messageId":       messageId,
			"dialogRequestId": dialogRequestId,
		},
		Payload: map[string]interface{}{
			"profile": "CLOSE_TALK",
			"format":  "AUDIO_L16_RATE_16000_CHANNELS_1",
		},
	}
	return &Recognize{parcel}
}

// The Speak directive.
type Speak struct {
	*Parcel
}

func (p *Speak) ContentId() string {
	url := p.Payload["url"].(string)
	if !strings.HasPrefix(url, "cid:") {
		return ""
	}
	return url[4:]
}

func (p *Speak) Format() string {
	return p.Payload["format"].(string)
}

func (p *Speak) URL() string {
	return p.Payload["url"].(string)
}

// The ExpectSpeech directive.
type ExpectSpeech struct {
	*Parcel
}

func (p *ExpectSpeech) Timeout() time.Duration {
	return time.Duration(p.Payload["timeoutInMilliseconds"].(float64)) * time.Millisecond
}

// The PlaybackState context.
type PlaybackState struct {
	*Parcel
}

type PlaybackStateActivity string

const (
	PlaybackStateActivityPlaying  = PlaybackStateActivity("PLAYING")
	PlaybackStateActivityFinished = PlaybackStateActivity("FINISHED")
)

func NewPlaybackState(token string, offset time.Duration, activity PlaybackStateActivity) *PlaybackState {
	parcel := &Parcel{
		Header: map[string]string{
			"namespace": "AudioPlayer",
			"name":      "PlaybackState",
		},
		Payload: map[string]interface{}{
			"token":                token,
			"offsetInMilliseconds": float64(offset.Seconds() * 1000),
			"playerActivity":       string(activity),
		},
	}
	return &PlaybackState{parcel}
}

func (p *PlaybackState) Offset() time.Duration {
	return time.Duration(p.Payload["offsetInMilliseconds"].(float64)) * time.Millisecond
}
