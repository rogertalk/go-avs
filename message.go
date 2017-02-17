package avs

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// TypedMessage is an interface that represents both raw Message objects and
// more specifically typed ones. Usually, values of this interface are used
// with a type switch:
//	switch d := typedMessage.(type) {
//	case *Speak:
//		fmt.Printf("We got a spoken response in format %s.\n", d.Payload.Format)
//	}
//
type TypedMessage interface {
	GetMessage() *Message
	Typed() TypedMessage
}

// Message is a general structure for contexts, events and directives.
type Message struct {
	Header  map[string]string `json:"header"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}

// GetMessage returns a pointer to the underlying Message object.
func (m *Message) GetMessage() *Message {
	return m
}

// String returns the namespace and name as a single string.
func (m *Message) String() string {
	return fmt.Sprintf("%s.%s", m.Header["namespace"], m.Header["name"])
}

// Typed returns a more specific type for this message.
//
// This only parses directives as they're the only type of message sent by AVS.
func (m *Message) Typed() TypedMessage {
	switch m.String() {
	case "Alerts.DeleteAlert":
		return fill(new(DeleteAlert), m)
	case "Alerts.SetAlert":
		return fill(new(SetAlert), m)
	case "AudioPlayer.ClearQueue":
		return fill(new(ClearQueue), m)
	case "AudioPlayer.Play":
		return fill(new(Play), m)
	case "AudioPlayer.Stop":
		return fill(new(Stop), m)
	case "Speaker.AdjustVolume":
		return fill(new(AdjustVolume), m)
	case "Speaker.SetMute":
		return fill(new(SetMute), m)
	case "Speaker.SetVolume":
		return fill(new(SetVolume), m)
	case "SpeechRecognizer.ExpectSpeech":
		return fill(new(ExpectSpeech), m)
	case "SpeechRecognizer.StopCapture":
		return fill(new(StopCapture), m)
	case "SpeechSynthesizer.Speak":
		return fill(new(Speak), m)
	case "System.Exception":
		// Exception is not a directive, but may also be sent by AVS.
		return fill(new(Exception), m)
	case "System.SetEndpoint":
		return fill(new(SetEndpoint), m)
	case "System.ResetUserInactivity":
		return fill(new(ResetUserInactivity), m)
	default:
		return m
	}
}

// The Exception message.
type Exception struct {
	*Message
	Payload struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"payload"`
}

// Error returns the Exception formatted as a human readable string.
func (m *Exception) Error() string {
	return fmt.Sprintf("%s: %s", m.Payload.Code, m.Payload.Description)
}

// Convenience function to set up an empty typed message object from a raw Message.
func fill(dst TypedMessage, src *Message) TypedMessage {
	v := reflect.ValueOf(dst).Elem()
	v.FieldByName("Message").Set(reflect.ValueOf(src))
	payload := v.FieldByName("Payload")
	if payload.Kind() != reflect.Struct {
		return dst
	}
	json.Unmarshal(src.Payload, payload.Addr().Interface())
	return dst
}
