package avs

import (
	"time"
)

// newEvent creates a Message suited for being used as an event value.
func newEvent(namespace, name, messageId, dialogRequestId string) *Message {
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
	m.Message = newEvent("Alerts", "AlertEnteredBackground", messageId, "")
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
	m.Message = newEvent("Alerts", "AlertEnteredForeground", messageId, "")
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
	m.Message = newEvent("Alerts", "AlertStarted", messageId, "")
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
	m.Message = newEvent("Alerts", "AlertStopped", messageId, "")
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
	m.Message = newEvent("Alerts", "DeleteAlertFailed", messageId, "")
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
	m.Message = newEvent("Alerts", "DeleteAlertSucceeded", messageId, "")
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
	m.Message = newEvent("Alerts", "SetAlertFailed", messageId, "")
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
	m.Message = newEvent("Alerts", "SetAlertSucceeded", messageId, "")
	m.Payload.Token = token
	return m
}

/********** AudioPlayer **********/
// Also used by the PlaybackState context.
type playbackState struct {
	Token                string         `json:"token"`
	OffsetInMilliseconds int            `json:"offsetInMilliseconds"`
	PlayerActivity       PlayerActivity `json:"playerActivity"`
}

// The PlaybackFailed event.
type PlaybackFailed struct {
	*Message
	Payload struct {
		Token                string        `json:"token"`
		CurrentPlaybackState playbackState `json:"currentPlaybackState"`
		Error struct {
			Type    MediaErrorType `json:"type"`
			Message string         `json:"message"`
		} `json:"error"`
	} `json:"payload"`
}

func NewPlaybackFailed(messageId, token string, errorType MediaErrorType, errorMessage string) *PlaybackFailed {
	m := new(PlaybackFailed)
	m.Message = newEvent("AudioPlayer", "PlaybackFailed", messageId, "")
	m.Payload.Token = token
	m.Payload.Error.Type = errorType
	m.Payload.Error.Message = errorMessage
	return m
}

// The PlaybackFinished event.
type PlaybackFinished struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackFinished(messageId, token string, offset time.Duration) *PlaybackFinished {
	m := new(PlaybackFinished)
	m.Message = newEvent("AudioPlayer", "PlaybackFinished", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackNearlyFinished event.
type PlaybackNearlyFinished struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackNearlyFinished(messageId, token string, offset time.Duration) *PlaybackNearlyFinished {
	m := new(PlaybackNearlyFinished)
	m.Message = newEvent("AudioPlayer", "PlaybackNearlyFinished", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackPaused event.
type PlaybackPaused struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackPaused(messageId, token string, offset time.Duration) *PlaybackPaused {
	m := new(PlaybackPaused)
	m.Message = newEvent("AudioPlayer", "PlaybackPaused", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackQueueCleared event.
type PlaybackQueueCleared struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewPlaybackQueueCleared(messageId string) *PlaybackQueueCleared {
	m := new(PlaybackQueueCleared)
	m.Message = newEvent("AudioPlayer", "PlaybackQueueCleared", messageId, "")
	return m
}

// The PlaybackResumed event.
type PlaybackResumed struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackResumed(messageId, token string, offset time.Duration) *PlaybackResumed {
	m := new(PlaybackResumed)
	m.Message = newEvent("AudioPlayer", "PlaybackResumed", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackStarted event.
type PlaybackStarted struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackStarted(messageId, token string, offset time.Duration) *PlaybackStarted {
	m := new(PlaybackStarted)
	m.Message = newEvent("AudioPlayer", "PlaybackStarted", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackStopped event.
type PlaybackStopped struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackStopped(messageId, token string, offset time.Duration) *PlaybackStopped {
	m := new(PlaybackStopped)
	m.Message = newEvent("AudioPlayer", "PlaybackStopped", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackStutterStarted event.
type PlaybackStutterStarted struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackStutterStarted(messageId, token string, offset time.Duration) *PlaybackStutterStarted {
	m := new(PlaybackStutterStarted)
	m.Message = newEvent("AudioPlayer", "PlaybackStutterStarted", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The PlaybackStutterFinished event.
type PlaybackStutterFinished struct {
	*Message
	Payload struct {
		Token                         string `json:"token"`
		OffsetInMilliseconds          int    `json:"offsetInMilliseconds"`
		StutterDurationInMilliseconds int    `json:"stutterDurationInMilliseconds"`
	} `json:"payload"`
}

func NewPlaybackStutterFinished(messageId, token string, offset, stutterDuration time.Duration) *PlaybackStutterFinished {
	m := new(PlaybackStutterFinished)
	m.Message = newEvent("AudioPlayer", "PlaybackStutterFinished", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	m.Payload.StutterDurationInMilliseconds = int(stutterDuration.Seconds() * 1000)
	return m
}

// The ProgressReportDelayElapsed event.
type ProgressReportDelayElapsed struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewProgressReportDelayElapsed(messageId, token string, offset time.Duration) *ProgressReportDelayElapsed {
	m := new(ProgressReportDelayElapsed)
	m.Message = newEvent("AudioPlayer", "ProgressReportDelayElapsed", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The ProgressReportIntervalElapsed event.
type ProgressReportIntervalElapsed struct {
	*Message
	Payload struct {
		Token                string `json:"token"`
		OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
	} `json:"payload"`
}

func NewProgressReportIntervalElapsed(messageId, token string, offset time.Duration) *ProgressReportIntervalElapsed {
	m := new(ProgressReportIntervalElapsed)
	m.Message = newEvent("AudioPlayer", "ProgressReportIntervalElapsed", messageId, "")
	m.Payload.Token = token
	m.Payload.OffsetInMilliseconds = int(offset.Seconds() * 1000)
	return m
}

// The StreamMetadataExtracted event.
type StreamMetadataExtracted struct {
	*Message
	Payload struct {
		Token    string                 `json:"token"`
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"payload"`
}

func NewStreamMetadataExtracted(messageId, token string, metadata map[string]interface{}) *StreamMetadataExtracted {
	m := new(StreamMetadataExtracted)
	m.Message = newEvent("AudioPlayer", "StreamMetadataExtracted", messageId, "")
	m.Payload.Token = token
	m.Payload.Metadata = metadata
	return m
}

/********** PlaybackController **********/

// The NextCommandIssued event.
type NextCommandIssued struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewNextCommandIssued(messageId string) *NextCommandIssued {
	m := new(NextCommandIssued)
	m.Message = newEvent("PlaybackController", "NextCommandIssued", messageId, "")
	return m
}

// The PauseCommandIssued event.
type PauseCommandIssued struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewPauseCommandIssued(messageId string) *PauseCommandIssued {
	m := new(PauseCommandIssued)
	m.Message = newEvent("PlaybackController", "PauseCommandIssued", messageId, "")
	return m
}

// The PlayCommandIssued event.
type PlayCommandIssued struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewPlayCommandIssued(messageId string) *PlayCommandIssued {
	m := new(PlayCommandIssued)
	m.Message = newEvent("PlaybackController", "PlayCommandIssued", messageId, "")
	return m
}

// The PreviousCommandIssued event.
type PreviousCommandIssued struct {
	*Message
	Payload struct{} `json:"payload"`
}

func NewPreviousCommandIssued(messageId string) *PreviousCommandIssued {
	m := new(PreviousCommandIssued)
	m.Message = newEvent("PlaybackController", "PreviousCommandIssued", messageId, "")
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
	m.Message = newEvent("Speaker", "MuteChanged", messageId, "")
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
	m.Message = newEvent("Speaker", "VolumeChanged", messageId, "")
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
	m.Message = newEvent("SpeechRecognizer", "ExpectSpeechTimedOut", messageId, "")
	return m
}


// RecognizeProfile identifies the ASR profile associated with your product.
type RecognizeProfile string

// Possible values for RecognizeProfile.
// Supports three distinct profiles optimized for speech at varying distances.
const (
	RecognizeProfileCloseTalk = RecognizeProfile("CLOSE_TALK")
	RecognizeProfileNearField = RecognizeProfile("NEAR_FIELD")
	RecognizeProfileFarField  = RecognizeProfile("FAR_FIELD")
)

// The Recognize event.
type Recognize struct {
	*Message
	Payload struct {
		Profile RecognizeProfile `json:"profile"`
		Format  string           `json:"format"`
	} `json:"payload"`
}

func NewRecognize(messageId, dialogRequestId string) *Recognize {
	return NewRecognizeWithProfile(messageId, dialogRequestId, RecognizeProfileCloseTalk)
}

func NewRecognizeWithProfile(messageId, dialogRequestId string, profile RecognizeProfile) *Recognize {
	m := new(Recognize)
	m.Message = newEvent("SpeechRecognizer", "Recognize", messageId, dialogRequestId)
	m.Payload.Format = "AUDIO_L16_RATE_16000_CHANNELS_1"
	m.Payload.Profile = profile
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
	m.Message = newEvent("SpeechSynthesizer", "SpeechFinished", messageId, "")
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
	m.Message = newEvent("SpeechSynthesizer", "SpeechStarted", messageId, "")
	m.Payload.Token = token
	return m
}

/********** Settings **********/

// The SettingsUpdated event.
type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingsUpdated struct {
	*Message
	Payload struct {
		Settings []Setting `json:"settings"`
	} `json:"payload"`
}

type SettingLocale string
// Possible values for SettingLocale.
const (
	SettingLocaleUS = SettingLocale("en-US")
	SettingLocaleGB = SettingLocale("en-GB")
	SettingLocaleDE = SettingLocale("de-DE")
)

func NewLocaleSettingsUpdated(messageId string, locale SettingLocale) *SettingsUpdated {
	m := new(SettingsUpdated)
	m.Message = newEvent("Settings", "SettingsUpdated", messageId, "")
	m.Payload.Settings = append(m.Payload.Settings, Setting{
		Key: "locale",
		Value: string(locale),
	})
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
	m.Message = newEvent("System", "ExceptionEncountered", messageId, "")
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
	m.Message = newEvent("System", "SynchronizeState", messageId, "")
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
	m.Message = newEvent("System", "UserInactivityReport", messageId, "")
	m.Payload.InactiveTimeInSeconds = int(inactiveTime.Seconds())
	return m
}
