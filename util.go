package avs

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
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
