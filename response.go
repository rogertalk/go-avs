package avs

// Response represents a response from the AVS API.
type Response struct {
	// The Amazon request id (for debugging purposes).
	RequestId string
	// All the directives in the response.
	Directives []*Message
	// Attachments (usually audio). Key is the Content-ID header value.
	Content map[string][]byte
}
