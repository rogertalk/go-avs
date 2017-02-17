/*
Package avs makes requests to Amazon's AVS API using HTTP/2 and also supports
creating a downchannel which receives directives from AVS based on the user's
speech and/or actions in the companion app.

To send a Recognize event to AVS, you can use the PostRecognize function:

	audio, _ := os.Open("./request.wav")
	response, err := avs.PostRecognize(ACCESS_TOKEN, "abc123",
	                                   "abc123dialog", audio)

For simple events you can use the PostEvent function:

	event := NewVolumeChanged("abc123", 100, false)
	response, err := avs.PostEvent(ACCESS_TOKEN, event)

You can also make requests to AVS with the Client.Do method:

	request := avs.NewRequest(ACCESS_TOKEN)
	request.Event = avs.NewRecognize("abc123", "abc123dialog")
	request.Audio, _ = os.Open("./request.wav")
	response, err := avs.DefaultClient.Do(request)

A Response will contain a list of directives from AVS. The list contains untyped
Message instances which hold the raw response data and headers, but it can be
typed by calling the Typed method of Message:

	for _, directive := range response.Directives {
		switch d := directive.Typed().(type) {
		case *avs.Speak:
			cid := d.ContentId()
			ioutil.WriteFile("./speak.mp3", response.Content[cid], 0666)
		default:
			fmt.Println("No code to handle directive:", d)
		}
	}

To create a downchannel, a long-lived request for AVS to deliver directives,
use the CreateDownchannel method of the Client type:

	directives, _ := avs.CreateDownchannel(ACCESS_TOKEN)
	for directive := range directives {
		switch d := directive.Typed().(type) {
		case *avs.DeleteAlert:
			fmt.Println("Delete alert:", d.Payload.Token)
		case *avs.SetAlert:
			fmt.Println("Set an alert for:", d.Payload.ScheduledTime)
		default:
			fmt.Println("No code to handle directive:", d)
		}
	}
*/
package avs

import (
	"io"
)

// The different endpoints that are supported by the AVS API.
const (
	// You can find the latest versioning information on the AVS API overview page:
	// https://developer.amazon.com/public/solutions/alexa/alexa-voice-service/content/avs-api-overview
	Version        = "/v20160207"
	DirectivesPath = Version + "/directives"
	EventsPath     = Version + "/events"
	PingPath       = "/ping"
)

// DefaultClient is the default Client.
var DefaultClient = &Client{
	// EndpointURL is the base endpoint URL for the AVS API.
	EndpointURL: "https://avs-alexa-na.amazon.com",
}

// CreateDownchannel establishes a persistent connection with AVS and returns a
// read-only channel through which AVS will deliver directives.
//
// CreateDownchannel is a wrapper around DefaultClient.CreateDownchannel.
func CreateDownchannel(accessToken string) (<-chan *Message, error) {
	return DefaultClient.CreateDownchannel(accessToken)
}

// PostEvent will post an event to AVS.
//
// PostEvent is a wrapper around DefaultClient.Do.
func PostEvent(accessToken string, event TypedMessage) (*Response, error) {
	request := NewRequest(accessToken)
	request.Event = event
	return DefaultClient.Do(request)
}

// PostRecognize will post a Recognize event to AVS.
//
// PostRecognize is a wrapper around DefaultClient.Do.
func PostRecognize(accessToken, messageId, dialogRequestId string, audio io.Reader) (*Response, error) {
	request := NewRequest(accessToken)
	request.Event = NewRecognize(messageId, dialogRequestId)
	request.Audio = audio
	return DefaultClient.Do(request)
}

// PostSynchronizeState will post a SynchronizeState event with the provided
// context to AVS.
//
// PostSynchronizeState is a wrapper around DefaultClient.Do.
func PostSynchronizeState(accessToken, messageId string, context []TypedMessage) (*Response, error) {
	request := NewRequest(accessToken)
	request.Event = NewSynchronizeState(messageId)
	request.Context = context
	return DefaultClient.Do(request)
}
