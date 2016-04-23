package avs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// A streamable audio item.
type AudioItem struct {
	AudioItemId string `json:"audioItemId"`
	Stream      Stream `json:"stream"`
}

// A value for how an audio item should be inserted into the play queue.
type PlayBehavior string

// TODO: Complete these constants.
const (
	PlayBehaviorReplaceAll = PlayBehavior("REPLACE_ALL")
)

// A value for what state the device player is in.
type PlayerActivity string

// TODO: Complete these constants.
const (
	PlayerActivityPlaying  = PlayerActivity("PLAYING")
	PlayerActivityFinished = PlayerActivity("FINISHED")
)

type ProgressReport struct {
	ProgressReportIntervalInMilliseconds float64 `json:"progressReportIntervalInMilliseconds"`
}

type Stream struct {
	ExpiryTime           string         `json:"expiryTime"`
	OffsetInMilliseconds float64        `json:"offsetInMilliseconds"`
	ProgressReport       ProgressReport `json:"progressReport"`
	Token                string         `json:"token"`
	URL                  string         `json:"url"`
}

// An interface that represents both raw Parcel objects and more specifically
// typed ones. Usually, values of this interface are used with a type switch:
//	switch d := typedParcel.(type) {
//	case *Speak:
//		fmt.Printf("We got a spoken response in format %s.\n", d.Payload.Format)
//	}
//
type TypedParcel interface {
	GetParcel() *Parcel
	Typed() TypedParcel
}

// A general structure for contexts, events and directives.
type Parcel struct {
	Header  map[string]string `json:"header"`
	Payload json.RawMessage   `json:"payload,omitempty"`
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
	case "AudioPlayer.Play":
		return fill(new(Play), p)
	case "AudioPlayer.PlaybackState":
		return fill(new(PlaybackState), p)
	case "SpeechRecognizer.ExpectSpeech":
		return fill(new(ExpectSpeech), p)
	case "SpeechRecognizer.Recognize":
		return fill(new(Recognize), p)
	case "SpeechSynthesizer.Speak":
		return fill(new(Speak), p)
	default:
		return p
	}
}

// The Play directive.
type Play struct {
	*Parcel
	Payload struct {
		AudioItem    AudioItem    `json:"audioItem"`
		PlayBehavior PlayBehavior `json:"playBehavior"`
	} `json:"payload"`
}

func (p *Play) DialogRequestId() string {
	return p.Header["dialogRequestId"]
}

func (p *Play) MessageId() string {
	return p.Header["messageId"]
}

// The Recognize event.
type Recognize struct {
	*Parcel
	Payload struct {
		Profile string `json:"profile"`
		Format  string `json:"format"`
	} `json:"payload"`
}

func NewRecognize(messageId, dialogRequestId string) *Recognize {
	p := new(Recognize)
	p.Parcel = &Parcel{
		Header: map[string]string{
			"namespace":       "SpeechRecognizer",
			"name":            "Recognize",
			"messageId":       messageId,
			"dialogRequestId": dialogRequestId,
		},
		Payload: nil,
	}
	p.Payload.Format = "AUDIO_L16_RATE_16000_CHANNELS_1"
	p.Payload.Profile = "CLOSE_TALK"
	return p
}

// The Speak directive.
type Speak struct {
	*Parcel
	Payload struct {
		Format string
		URL    string
	} `json:"payload"`
}

func (p *Speak) ContentId() string {
	if !strings.HasPrefix(p.Payload.URL, "cid:") {
		return ""
	}
	return p.Payload.URL[4:]
}

// The ExpectSpeech directive.
type ExpectSpeech struct {
	*Parcel
	Payload struct {
		TimeoutInMilliseconds float64 `json:"timeoutInMilliseconds"`
	} `json:"payload"`
}

// The PlaybackState context.
type PlaybackState struct {
	*Parcel
	Payload struct {
		Token                string         `json:"token"`
		OffsetInMilliseconds float64        `json:"offsetInMilliseconds"`
		PlayerActivity       PlayerActivity `json:"playerActivity"`
	}
}

func NewPlaybackState(token string, offset time.Duration, activity PlayerActivity) *PlaybackState {
	p := new(PlaybackState)
	p.Parcel = &Parcel{
		Header: map[string]string{
			"namespace": "AudioPlayer",
			"name":      "PlaybackState",
		},
		Payload: nil,
	}
	p.Payload.OffsetInMilliseconds = offset.Seconds() * 1000
	p.Payload.PlayerActivity = activity
	p.Payload.Token = token
	return p
}

// Convenience method to set up an empty typed parcel object from a raw Parcel.
func fill(dst TypedParcel, src *Parcel) TypedParcel {
	v := reflect.ValueOf(dst).Elem()
	v.FieldByName("Parcel").Set(reflect.ValueOf(src))
	payload := v.FieldByName("Payload")
	if payload.Kind() != reflect.Struct {
		return dst
	}
	json.Unmarshal(src.Payload, payload.Addr().Interface())
	return dst
}
