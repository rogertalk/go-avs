package avs

import (
	"strings"
	"time"
)

/********** Alerts **********/

// The DeleteAlert directive.
type DeleteAlert struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

// The SetAlert directive.
type SetAlert struct {
	*Message
	Payload Alert `json:"payload"`
}

/********** AudioPlayer **********/

// The ClearQueue directive.
type ClearQueue struct {
	*Message
	Payload struct {
		ClearBehavior ClearBehavior `json:"clearBehavior"`
	} `json:"payload"`
}

// The Play directive.
type Play struct {
	*Message
	Payload struct {
		AudioItem    AudioItem    `json:"audioItem"`
		PlayBehavior PlayBehavior `json:"playBehavior"`
	} `json:"payload"`
}

// The Stop directive.
type Stop struct {
	*Message
	Payload struct{} `json:"payload"`
}

/********** Speaker **********/

// The AdjustVolume directive.
type AdjustVolume struct {
	*Message
	Payload struct {
		Volume int `json:"volume"`
	} `json:"payload"`
}

// The SetMute directive.
type SetMute struct {
	*Message
	Payload struct {
		Mute bool `json:"mute"`
	} `json:"payload"`
}

// The SetVolume directive.
type SetVolume struct {
	*Message
	Payload struct {
		Volume int `json:"volume"`
	} `json:"payload"`
}

/********** SpeechRecognizer **********/

// The ExpectSpeech directive.
type ExpectSpeech struct {
	*Message
	Payload struct {
		TimeoutInMilliseconds int `json:"timeoutInMilliseconds"`
	} `json:"payload"`
}

func (m *ExpectSpeech) Timeout() time.Duration {
	return time.Duration(m.Payload.TimeoutInMilliseconds) * time.Millisecond
}

// The StopCapture directive.
type StopCapture struct {
	*Message
	Payload struct{} `json:"payload"`
}

/********** SpeechSynthesizer **********/

// The Speak directive.
type Speak struct {
	*Message
	Payload struct {
		Format string `json:"format"`
		URL    string `json:"url"`
		Token  string `json:"token"`
	} `json:"payload"`
}

func (m *Speak) ContentId() string {
	if !strings.HasPrefix(m.Payload.URL, "cid:") {
		return ""
	}
	return m.Payload.URL[4:]
}

/********** System **********/

// The SetEndpoint directive.
type SetEndpoint struct {
	*Message
	Payload struct {
		Endpoint string `json:"endpoint"`
	} `json:"payload"`
}

// The ResetUserInactivity directive.
type ResetUserInactivity struct {
	*Message
	Payload struct{} `json:"payload"`
}
