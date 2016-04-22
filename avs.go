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
	"net/textproto"
	"strings"
)

const (
	EndpointURL   = "https://avs-alexa-na.amazon.com/v20160207"
	DirectivesURL = EndpointURL + "/directives"
	EventsURL     = EndpointURL + "/events"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// Encodes a JSON value and writes it to a field with the provided multipart writer.
func writeJSON(writer *multipart.Writer, fieldname string, value interface{}) error {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname)))
	h.Set("Content-Type", "application/json; charset=UTF-8")
	p, err := writer.CreatePart(h)
	if err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = p.Write(data)
	return err
}

// An AVS client.
type Client struct {
}

var DefaultClient = &Client{}

// Request to AVS.
type Request struct {
	// Access token for the user that this request should be made for.
	AccessToken string        `json:"-"`
	Audio       io.Reader     `json:"-"`
	Context     []TypedParcel `json:"context"`
	Event       TypedParcel   `json:"event"`
}

func NewRequest(accessToken string) *Request {
	return &Request{
		AccessToken: accessToken,
		Context:     []TypedParcel{},
	}
}

func (r *Request) AddContext(p TypedParcel) {
	r.Context = append(r.Context, p.GetParcel())
}

// Representation of an AVS response.
type Response struct {
	// The Amazon request id (for debugging purposes).
	RequestId string
	// All the directives in the response.
	Directives []TypedParcel
	// Attachments (usually audio). Key is the Content-ID header value.
	Content map[string][]byte
}

// Multipart object returned by AVS.
type responsePart struct {
	Directive *Parcel
}

// Posts a request to the AVS service.
func (c *Client) Do(request *Request) (*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writeJSON(writer, "metadata", request)
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
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with %s", resp.Status)
	}
	// Parse the multipart response.
	mediatype, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediatype, "multipart/") {
		err = fmt.Errorf("unexpected content type %s", mediatype)
		return nil, err
	}
	mr := multipart.NewReader(resp.Body, params["boundary"])
	response := &Response{
		RequestId:  resp.Header.Get("x-amzn-requestid"),
		Directives: []TypedParcel{},
		Content:    make(map[string][]byte),
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