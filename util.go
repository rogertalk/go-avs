package avs

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/rogertalk/go-avs/multipart2"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// UUID holds a 16 byte unique identifier.
type UUID []byte

// String returns a string representation of the UUID.
func (uuid UUID) String() string {
	if len(uuid) != 16 {
		return ""
	}
	b := []byte(uuid)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// NewUUID returns a randomly generated UUID.
func NewUUID() (UUID, error) {
	b := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		return nil, err
	}
	b[8] = b[8]&^0xc0 | 0x80
	b[6] = b[6]&^0xf0 | 0x40
	return UUID(b), nil
}

func RandomUUIDString() string {
	uuid, err := NewUUID()
	if err != nil {
		panic("RandomUUIDString: unable to create UUID")
	}
	return uuid.String()
}

func newMultipartReaderFromResponse(resp *http.Response) (*multipart2.Reader, error) {
	// Work around bug in Amazon's downchannel server.
	contentType := strings.Replace(resp.Header.Get("Content-Type"), "type=application/json", `type="application/json"`, 1)
	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediatype, "multipart/") {
		return nil, fmt.Errorf("unexpected content type %s", mediatype)
	}
	return multipart2.NewReader(resp.Body, params["boundary"]), nil
}

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
