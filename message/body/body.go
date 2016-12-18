package body

import (
	"io"
	"net/textproto"

	"github.com/chris-tomich/go-mail/message/headers"
	"github.com/chris-tomich/go-mail/message/headers/mime"
)

// MailBody represents HTML text for use in an HTML email.
type MailBody struct {
	body    string
	Headers textproto.MIMEHeader
}

// GenerateBody will perform any parsing and return the message body.
func (b *MailBody) GenerateBody() string {
	return b.body
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
		body: body,
	}

	b.Headers = make(textproto.MIMEHeader)
	b.Headers.Add(headers.ContentType(mime.TextPlain))

	return b
}

// TemplateFromFile will return a MailBody based upon the template at the given filename.
// A pointer to an interface{} is accepted so that this object can be updated with different data.
// This allows for generating multiple emails with different data from the same template.
func TemplateFromFile(filename string, data *interface{}) *MailBody {
	return &MailBody{}
}

// TemplateFromReader will return a MailBody based upon the template using the provided data.
func TemplateFromReader(reader io.Reader, data interface{}) *MailBody {
	return &MailBody{}
}

// TemplateFromString
func TemplateFromString(body string) *MailBody {
	return &MailBody{
		body: body,
	}
}
