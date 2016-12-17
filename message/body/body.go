package body

import (
	"io"
	"net/textproto"

	"github.com/chris-tomich/go-mail/message/headers"
	"github.com/chris-tomich/go-mail/message/headers/mime"
)

// MailBody represents HTML text for use in an HTML email.
type MailBody struct {
	Body    string
	Headers textproto.MIMEHeader
}

// FromFile will return a MailBody based upon the text at the given filename.
func FromFile(filename string) *MailBody {
	return &MailBody{}
}

// FromReader will return a MailBody based upon the provided data.
func FromReader(reader io.Reader) *MailBody {
	return &MailBody{}
}

// FromString will return a MailBody based upon the text provided.
func FromString(body string) *MailBody {
	b := &MailBody{
		Body: body,
	}

	b.Headers = make(textproto.MIMEHeader)
	b.Headers.Add(headers.ContentType(mime.TextPlain))

	return b
}

// TemplateFromFile will return a MailBody based upon the template at the given filename using
// the provided data.
func TemplateFromFile(filename string, data interface{}) *MailBody {
	return &MailBody{}
}

// TemplateFromReader will return a MailBody based upon the template using the provided data.
func TemplateFromReader(reader io.Reader, data interface{}) *MailBody {
	return &MailBody{}
}
