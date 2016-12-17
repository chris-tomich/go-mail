package message

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

	"github.com/chris-tomich/go-mail/message/attachments"
	"github.com/chris-tomich/go-mail/message/body"
	"github.com/chris-tomich/go-mail/message/headers"
)

// Message represents a mail message to be sent.
// At the very minimum a message must contain a from address, a to address. Everything else is optional.
type Message struct {
	buffer      *bytes.Buffer
	writer      *multipart.Writer
	headers     textproto.MIMEHeader
	textBody    string
	htmlBody    string
	images      map[string]*attachments.EmbeddedBinaryObject
	attachments map[string]*attachments.EmbeddedBinaryObject
}

// NewEmpty creates a new go-mail Message object that is completely empty (no headers). You must use this to create
// an empty Message object as there are internal members that need to be instantiated.
func NewEmpty() *Message {
	m := &Message{}
	m.buffer = &bytes.Buffer{}
	m.writer = multipart.NewWriter(m.buffer)
	m.headers = make(textproto.MIMEHeader)
	return m
}

// NewSimple creates a simple go-mail Message with a bare minimum of headers set. The following lists the headers set and their values.
// "From" - with the address provided.
// "To" - with the address provided. To add more, use the AddMailHeader method and headers.To() helper function.
// "MIME-Version" - "1.0"
//
// Addresses can be in the following formats -
// Simple Email - chris.tomich@email.com
// Email with Name - Chris Tomich <chris.tomich@email.com>
func NewSimple(from string, to string, subject string, bodyText string, attachmentFilePaths ...string) *Message {
	m := NewEmpty()
	m.AddMailHeader(headers.From(from))
	m.AddMailHeader(headers.To(to))
	m.AddMailHeader(headers.Subject(subject))
	m.AddMailHeader(headers.MIMEVersion(1.0))
	m.SetTextBody(body.FromString(bodyText))

	for _, filePath := range attachmentFilePaths {
		m.AddAttachment(attachments.FileAttachmentFromFile(filePath))
	}

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
func (m *Message) SetTextBody(body *body.MailBody) {
	m.textBody = body.Body
}

// SetHTMLBody will set/overwrite the body for the HTML portion of the email and override all currently stored images.
// This body is additional to the text-only body and does not impact anything set with SetTextBody().
func (m *Message) SetHTMLBody(body *body.MailBody, images ...*attachments.EmbeddedBinaryObject) {
	m.htmlBody = body.Body

	m.images = make(map[string]*attachments.EmbeddedBinaryObject)

	for _, inlineImage := range images {
		m.images[inlineImage.FileName] = inlineImage
	}
}

// AddAttachment will add an attachment to the email.
func (m *Message) AddAttachment(attachment *attachments.EmbeddedBinaryObject) {
	m.attachments[attachment.FileName] = attachment
}

// RemoveAttachment will remove the attachment from the email.
func (m *Message) RemoveAttachment(filename string) {
	delete(m.attachments, filename)
}

func serialiseHeaders(w io.Writer, headers textproto.MIMEHeader) error {
	for k, vs := range headers {
		if len(vs) == 1 {
			_, err := fmt.Fprintf(w, "%v: %v\r\n", k, vs[0])

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Message) GenerateMessage() *bytes.Buffer {
	//err := serialiseHeaders(m.buffer, m.headers)

	return nil
}
