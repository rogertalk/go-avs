package avs

import (
	"io"
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
