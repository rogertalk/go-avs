/*
Package avs makes requests to Amazon's AVS API using HTTP/2 and also supports
creating a downchannel which receives directives from AVS based on the user's
speech and/or actions in the companion app.

You can make requests to AVS with the Do method:

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
			ioutil.WriteFile("./speak.wav", response.Content[d.ContentId()], 0666)
		default:
			fmt.Println("No code to handle directive:", d)
		}
	}

To create a downchannel, which is a long-lived request for AVS to deliver
directives, use the CreateDownchannel method of the Client type:

	directives, _ := avs.DefaultClient.CreateDownchannel(ACCESS_TOKEN)
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
)

// The different endpoints that are supported by the AVS API.
const (
	// EndpointURL is the base endpoint URL for the API, including version.
	// You can find the latest versioning information on the AVS API overview page:
	// https://developer.amazon.com/public/solutions/alexa/alexa-voice-service/content/avs-api-overview
	EndpointURL   = "https://avs-alexa-na.amazon.com/v20160207"
	DirectivesURL = EndpointURL + "/directives"
	EventsURL     = EndpointURL + "/events"
	PingURL       = EndpointURL + "/ping"
)

// A Request represents an event and optional context to send to AVS.
type Request struct {
	// Access token for the user that this request should be made for.
	AccessToken string         `json:"-"`
	Audio       io.Reader      `json:"-"`
	Context     []TypedMessage `json:"context"`
	Event       TypedMessage   `json:"event"`
}

// NewRequest returns a new Request given an access token.
//
// The Request is suitable for use with Client.Do.
func NewRequest(accessToken string) *Request {
	return &Request{
		AccessToken: accessToken,
		Context:     []TypedMessage{},
	}
}

// AddContext adds a context Message to the Request.
func (r *Request) AddContext(m TypedMessage) {
	r.Context = append(r.Context, m)
}

// Response represents a response from the AVS API.
type Response struct {
	// The Amazon request id (for debugging purposes).
	RequestId string
	// All the directives in the response.
	Directives []TypedMessage
	// Attachments (usually audio). Key is the Content-ID header value.
	Content map[string][]byte
}

// Multipart object returned by AVS.
type responsePart struct {
	Directive *Message
}

// Client enables making requests and creating downchannels to AVS.
type Client struct {
}

// DefaultClient is the default Client.
var DefaultClient = &Client{}

// CreateDownchannel establishes a persistent connection with AVS and returns a
// read-only channel through which AVS will deliver directives.
func (c *Client) CreateDownchannel(accessToken string) (<-chan TypedMessage, error) {
	req, err := http.NewRequest("GET", DirectivesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	directives := make(chan TypedMessage)
	go func() {
		defer close(directives)
		defer resp.Body.Close()
		mr, err := newMultipartReaderFromResponse(resp)
		if err != nil {
			return
		}
		// TODO: Consider reporting errors.
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}
			data, err := ioutil.ReadAll(p)
			if err != nil {
				break
			}
			var response responsePart
			err = json.Unmarshal(data, &response)
			if err != nil {
				break
			}
			directives <- response.Directive.Typed()
		}
	}()
	return directives, nil
}

// Do posts a request to the AVS service.
func (c *Client) Do(request *Request) (*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writeJSON(writer, "metadata", request)
	if err != nil {
		return nil, err
	}
	if request.Audio != nil {
		p, err := writer.CreateFormFile("audio", "audio.wav")
		if err != nil {
			return nil, err
		}
		// TODO: Write the audio data directly to the HTTP/2 socket for faster delivery.
		_, err = io.Copy(p, request.Audio)
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	}
	// Send the request to AVS.
	req, err := http.NewRequest("POST", EventsURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.AccessToken))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		// Keep going.
	case 204:
		// No content.
		return nil, nil
	default:
		// Attempt to parse the response as a System.Exception message.
		data, _ := ioutil.ReadAll(resp.Body)
		var exception Exception
		json.Unmarshal(data, &exception)
		if exception.Payload.Code != "" {
			return nil, &exception
		}
		// Fallback error.
		return nil, fmt.Errorf("request failed with %s", resp.Status)
	}
	// Parse the multipart response.
	mr, err := newMultipartReaderFromResponse(resp)
	if err != nil {
		return nil, err
	}
	response := &Response{
		RequestId:  resp.Header.Get("x-amzn-requestid"),
		Directives: []TypedMessage{},
		Content:    map[string][]byte{},
	}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		mediatype, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(p)
		if err != nil {
			return nil, err
		}
		if contentId := p.Header.Get("Content-ID"); contentId != "" {
			// This part is a referencable piece of content.
			// XXX: Content-ID generally always has angle brackets, but there may be corner cases?
			response.Content[contentId[1:len(contentId)-1]] = data
		} else if mediatype == "application/json" {
			// This is a directive.
			var resp responsePart
			err = json.Unmarshal(data, &resp)
			if err != nil {
				return nil, err
			}
			if resp.Directive == nil {
				return nil, fmt.Errorf("missing directive %s", string(data))
			}
			response.Directives = append(response.Directives, resp.Directive)
		} else {
			return nil, fmt.Errorf("unhandled part %s", p)
		}
	}
	return response, nil
}

// Ping will ping AVS on behalf of a user to indicate that the connection is
// still alive.
func (c *Client) Ping(accessToken string) error {
	// TODO: Once Go supports sending PING frames, that would be a better alternative.
	req, err := http.NewRequest("GET", PingURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
