package html

import (
	"io"
)

// MailHTMLBody represents HTML text for use in an HTML email.
type MailHTMLBody string

// TemplateFromFile will return a MailHTMLBody based upon the template at the given filename using
// the provided data.
func TemplateFromFile(filename string, data interface{}) MailHTMLBody {
	return MailHTMLBody("")
}

// TemplateFromReader will return a MailHTMLBody based upon the template using the provided data.
func TemplateFromReader(reader io.Reader, data interface{}) MailHTMLBody {
	return MailHTMLBody("")
}
