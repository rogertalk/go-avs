package avs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
)

// Multipart object returned by AVS.
type responsePart struct {
	Directive *Message
}

// Client enables making requests and creating downchannels to AVS.
type Client struct {
	EndpointURL string
}

// CreateDownchannel establishes a persistent connection with AVS and returns a
// read-only channel through which AVS will deliver directives.
func (c *Client) CreateDownchannel(accessToken string) (<-chan *Message, error) {
	req, err := http.NewRequest("GET", c.EndpointURL + DirectivesPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if more, err := checkStatusCode(resp); !more {
		resp.Body.Close()
		return nil, err
	}
	directives := make(chan *Message)
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
			directives <- response.Directive
		}
	}()
	return directives, nil
}

// Do posts a request to the AVS service's /events endpoint.
func (c *Client) Do(request *Request) (*Response, error) {
	body, bodyIn := io.Pipe()
	writer := multipart.NewWriter(bodyIn)
	go func() {
		// Write to pipe must be parallel to allow HTTP request to read
		err := writeJSON(writer, "metadata", request)
		if err != nil {
			bodyIn.CloseWithError(err)
			return
		}
		if request.Audio != nil {
			p, err := writer.CreateFormFile("audio", "audio.wav")
			if err != nil {
				bodyIn.CloseWithError(err)
				return
			}
			// Run io.Copy in goroutine so audio can be streamed
			_, err = io.Copy(p, request.Audio)
			if err != nil {
				bodyIn.CloseWithError(err)
				return
			}
		}
		err = writer.Close()
		if err != nil {
			bodyIn.CloseWithError(err)
			return
		}
		bodyIn.Close()
	}()
	// Send the request to AVS.
	req, err := http.NewRequest("POST", c.EndpointURL + EventsPath, body)
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
	more, err := checkStatusCode(resp)
	if err != nil {
		return nil, err
	}
	response := &Response{
		RequestId:  resp.Header.Get("x-amzn-requestid"),
		Directives: []*Message{},
		Content:    map[string][]byte{},
	}
	if !more {
		// AVS returned an empty response, so there's nothing to parse.
		return response, nil
	}
	// Parse the multipart response.
	mr, err := newMultipartReaderFromResponse(resp)
	if err != nil {
		return nil, err
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
	req, err := http.NewRequest("GET", c.EndpointURL + PingPath, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = checkStatusCode(resp)
	return err
}

// Checks the status code of the response and returns whether the caller should
// expect there to be more content, as well as any error.
//
// This function should only be called before the body has been read.
func checkStatusCode(resp *http.Response) (more bool, err error) {
	switch resp.StatusCode {
	case 200:
		// Keep going.
		return true, nil
	case 204:
		// No content.
		return false, nil
	default:
		// Attempt to parse the response as a System.Exception message.
		data, _ := ioutil.ReadAll(resp.Body)
		var exception Exception
		json.Unmarshal(data, &exception)
		if exception.Payload.Code != "" {
			return false, &exception
		}
		// Fallback error.
		return false, fmt.Errorf("request failed with %s", resp.Status)
	}
}
