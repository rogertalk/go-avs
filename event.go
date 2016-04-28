package avs

import (
	"time"
)

// NewEvent creates a Message suited for being used as an event value.
func NewEvent(namespace, name, messageId, dialogRequestId string) *Message {
	m := &Message{
		Header: map[string]string{
			"namespace": namespace,
			"name":      name,
			"messageId": messageId,
		},
		Payload: nil,
	}
	if dialogRequestId != "" {
		m.Header["dialogRequestId"] = dialogRequestId
	}
	return m
}

/********** Alerts **********/

// The AlertEnteredBackground event.
type AlertEnteredBackground struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewAlertEnteredBackground(messageId, token string) *AlertEnteredBackground {
	m := new(AlertEnteredBackground)
	m.Message = NewEvent("Alerts", "AlertEnteredBackground", messageId, "")
	m.Payload.Token = token
	return m
}

// The AlertEnteredForeground event.
type AlertEnteredForeground struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewAlertEnteredForeground(messageId, token string) *AlertEnteredForeground {
	m := new(AlertEnteredForeground)
	m.Message = NewEvent("Alerts", "AlertEnteredForeground", messageId, "")
	m.Payload.Token = token
	return m
}

// The AlertStarted event.
type AlertStarted struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewAlertStarted(messageId, token string) *AlertStarted {
	m := new(AlertStarted)
	m.Message = NewEvent("Alerts", "AlertStarted", messageId, "")
	m.Payload.Token = token
	return m
}

// The AlertStopped event.
type AlertStopped struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewAlertStopped(messageId, token string) *AlertStopped {
	m := new(AlertStopped)
	m.Message = NewEvent("Alerts", "AlertStopped", messageId, "")
	m.Payload.Token = token
	return m
}

// The DeleteAlertFailed event.
type DeleteAlertFailed struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewDeleteAlertFailed(messageId, token string) *DeleteAlertFailed {
	m := new(DeleteAlertFailed)
	m.Message = NewEvent("Alerts", "DeleteAlertFailed", messageId, "")
	m.Payload.Token = token
	return m
}

// The DeleteAlertSucceeded event.
type DeleteAlertSucceeded struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewDeleteAlertSucceeded(messageId, token string) *DeleteAlertSucceeded {
	m := new(DeleteAlertSucceeded)
	m.Message = NewEvent("Alerts", "DeleteAlertSucceeded", messageId, "")
	m.Payload.Token = token
	return m
}

// The SetAlertFailed event.
type SetAlertFailed struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewSetAlertFailed(messageId, token string) *SetAlertFailed {
	m := new(SetAlertFailed)
	m.Message = NewEvent("Alerts", "SetAlertFailed", messageId, "")
	m.Payload.Token = token
	return m
}

// The SetAlertSucceeded event.
type SetAlertSucceeded struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewSetAlertSucceeded(messageId, token string) *SetAlertSucceeded {
	m := new(SetAlertSucceeded)
	m.Message = NewEvent("Alerts", "SetAlertSucceeded", messageId, "")
	m.Payload.Token = token
	return m
}

/********** Speaker **********/

// The MuteChanged event.
type MuteChanged struct {
	*Message
	Payload struct {
		Volume int  `json:"volume"`
		Muted  bool `json:"muted"`
	} `json:"payload"`
}

func NewMuteChanged(messageId string, volume int, muted bool) *MuteChanged {
	m := new(MuteChanged)
	m.Message = NewEvent("Speaker", "MuteChanged", messageId, "")
	m.Payload.Volume = volume
	m.Payload.Muted = muted
	return m
}

// The VolumeChanged event.
type VolumeChanged struct {
	*Message
	Payload struct {
		Volume int  `json:"volume"`
		Muted  bool `json:"muted"`
	} `json:"payload"`
}

func NewVolumeChanged(messageId string, volume int, muted bool) *VolumeChanged {
	m := new(VolumeChanged)
	m.Message = NewEvent("Speaker", "VolumeChanged", messageId, "")
	m.Payload.Volume = volume
	m.Payload.Muted = muted
	return m
}

/********** SpeechRecognizer **********/

// The ExpectSpeechTimedOut event.
type ExpectSpeechTimedOut struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewExpectSpeechTimedOut(messageId string) *ExpectSpeechTimedOut {
	m := new(ExpectSpeechTimedOut)
	m.Message = NewEvent("SpeechRecognizer", "ExpectSpeechTimedOut", messageId, "")
	return m
}

// The Recognize event.
type Recognize struct {
	*Message
	Payload struct {
		Profile string `json:"profile"`
		Format  string `json:"format"`
	} `json:"payload"`
}

func NewRecognize(messageId, dialogRequestId string) *Recognize {
	m := new(Recognize)
	m.Message = NewEvent("SpeechRecognizer", "Recognize", messageId, dialogRequestId)
	m.Payload.Format = "AUDIO_L16_RATE_16000_CHANNELS_1"
	m.Payload.Profile = "CLOSE_TALK"
	return m
}

/********** SpeechSynthesizer **********/

// The SpeechFinished event.
type SpeechFinished struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewSpeechFinished(messageId, token string) *SpeechFinished {
	m := new(SpeechFinished)
	m.Message = NewEvent("SpeechSynthesizer", "SpeechFinished", messageId, "")
	m.Payload.Token = token
	return m
}

// The SpeechStarted event.
type SpeechStarted struct {
	*Message
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

func NewSpeechStarted(messageId, token string) *SpeechStarted {
	m := new(SpeechStarted)
	m.Message = NewEvent("SpeechSynthesizer", "SpeechStarted", messageId, "")
	m.Payload.Token = token
	return m
}

/********** System **********/

// The ExceptionEncountered event.
type ExceptionEncountered struct {
	*Message
	Payload struct {
		UnparsedDirective string `json:"unparsedDirective"`
		Error             struct {
			Type    ErrorType `json:"type"`
			Message string    `json:"message"`
		} `json:"error"`
	} `json:"payload"`
}

func NewExceptionEncountered(messageId, directive string, errorType ErrorType, errorMessage string) *ExceptionEncountered {
	m := new(ExceptionEncountered)
	m.Message = NewEvent("System", "ExceptionEncountered", messageId, "")
	m.Payload.UnparsedDirective = directive
	m.Payload.Error.Type = errorType
	m.Payload.Error.Message = errorMessage
	return m
}

// The SynchronizeState event.
type SynchronizeState struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewSynchronizeState(messageId string) *SynchronizeState {
	m := new(SynchronizeState)
	m.Message = NewEvent("System", "SynchronizeState", messageId, "")
	return m
}

// The UserInactivityReport event.
type UserInactivityReport struct {
	*Message
	Payload struct {
		InactiveTimeInSeconds int `json:"inactiveTimeInSeconds"`
	} `json:"payload"`
}

func NewUserInactivityReport(messageId string, inactiveTime time.Duration) *UserInactivityReport {
	m := new(UserInactivityReport)
	m.Message = NewEvent("System", "UserInactivityReport", messageId, "")
	m.Payload.InactiveTimeInSeconds = int(inactiveTime.Seconds())
	return m
}
