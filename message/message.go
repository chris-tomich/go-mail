package message

import (
	"bytes"
	"mime/multipart"
	"net/textproto"

	"github.com/chris-tomich/go-mail/message/attachments"
	"github.com/chris-tomich/go-mail/message/html"
)

// Message represents a mail message to be sent.
// At the very minimum a message must contain a from address, a to address. Everything else is optional.
type Message struct {
	buffer      *bytes.Buffer
	writer      *multipart.Writer
	headers     textproto.MIMEHeader
	textBody    string
	htmlBody    string
	images      map[string]bytes.Buffer
	attachments map[string]bytes.Buffer
}

// New creats a new go-mail Message object.
func New() *Message {
	m := &Message{}
	m.buffer = &bytes.Buffer{}
	m.writer = multipart.NewWriter(m.buffer)
	m.headers = make(textproto.MIMEHeader)
	return m
}

// AddMailHeader adds a header to the email given the key and value.
// Look in the headers package to see a number of helpers for potential headers.
func (m *Message) AddMailHeader(key, value string) {
	if v := m.headers.Get(key); v != "" {
		appendedVal := v + "; " + value
		m.headers.Set(key, appendedVal)
	} else {
		m.headers.Add(key, value)
	}
}

// SetTextBody will set/overwrite the body for the text-only portion of the email.
// This body is additional to the HTML body and does not impact anything set with SetHTMLBody().
func (m *Message) SetTextBody(body string) {

}

// SetHTMLBody will set/overwrite the body for the HTML portion of the email.
// This body is additional to the text-only body and does not impact anything set with SetTextBody().
func (m *Message) SetHTMLBody(body html.MailHTMLBody, images ...attachments.InlineImage) {

}

// AddAttachment will add an attachment to the email.
func (m *Message) AddAttachment(attachment attachments.FileAttachment) {

}

// RemoveAttachment will remove the attachment from the email.
func (m *Message) RemoveAttachment(filename string) {

}
