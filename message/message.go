package message

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

	"strconv"
	"strings"
	"time"

	"github.com/chris-tomich/go-mail/message/attachments"
	"github.com/chris-tomich/go-mail/message/body"
	"github.com/chris-tomich/go-mail/message/headers"
	"github.com/chris-tomich/go-mail/message/headers/encoding"
	"github.com/chris-tomich/go-mail/message/headers/mime"
	"github.com/chris-tomich/go-mail/message/headers/params"
	"github.com/pkg/errors"
)

// Message represents a mail message to be sent.
// At the very minimum a message must contain a from address, a to address. Everything else is optional.
type Message struct {
	headers     textproto.MIMEHeader
	textBody    *body.MailBody
	htmlBody    *body.MailBody
	images      map[string]*attachments.EmbeddedBinaryObject
	attachments map[string]*attachments.EmbeddedBinaryObject
}

// NewEmpty creates a new go-mail Message object that is completely empty (no headers). You must use this to create
// an empty Message object as there are internal members that need to be instantiated.
func NewEmpty() *Message {
	m := &Message{}
	m.headers = make(textproto.MIMEHeader)
	m.images = make(map[string]*attachments.EmbeddedBinaryObject)
	m.attachments = make(map[string]*attachments.EmbeddedBinaryObject)
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
func NewSimple(from string, to string, subject string, bodyText string, attachmentFilePaths ...string) (*Message, error) {
	m := NewEmpty()
	m.AddMailHeader(headers.From(from))
	m.AddMailHeader(headers.To(to))
	m.AddMailHeader(headers.Subject(subject))
	m.AddMailHeader(headers.MIMEVersion(1.0))
	m.SetTextBody(body.FromString(bodyText))

	for _, filePath := range attachmentFilePaths {
		err := m.AddAttachment(attachments.FileAttachmentFromFile(filePath))

		if err != nil {
			return nil, err
		}
	}

	return m, nil
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
	m.textBody = body
}

// SetHTMLBody will set/overwrite the body for the HTML portion of the email and override all currently stored images.
// This body is additional to the text-only body and does not impact anything set with SetTextBody().
func (m *Message) SetHTMLBody(body *body.MailBody, images ...*attachments.EmbeddedBinaryObject) {
	m.htmlBody = body

	for _, inlineImage := range images {
		m.images[inlineImage.Filename] = inlineImage
	}
}

// AddAttachment will add an attachment to the email.
func (m *Message) AddAttachment(attachment *attachments.EmbeddedBinaryObject, err error) error {
	if err != nil {
		return err
	}

	m.attachments[attachment.Filename] = attachment

	return nil
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

func writeTextBody(alternativeBodyW *multipart.Writer, body string) error {
	textBodyHeaders := make(textproto.MIMEHeader)
	textBodyHeaders.Add(headers.ContentType(mime.TextPlain))

	p, err := alternativeBodyW.CreatePart(textBodyHeaders)

	if err != nil {
		return err
	}

	p.Write([]byte(body))

	return nil
}

func writeHtmlBody(alternativeBodyW *multipart.Writer, body string, images map[string]*attachments.EmbeddedBinaryObject) (map[string]string, error) {
	contentIds := make(map[string]string)

	htmlBodyHeaders := make(textproto.MIMEHeader)
	htmlBodyHeaders.Add(headers.ContentType(mime.TextHTML))
	htmlBodyP, err := alternativeBodyW.CreatePart(htmlBodyHeaders)

	if err != nil {
		return nil, err
	}

	counter := 10000

	if len(images) > 0 {
		for _, image := range images {
			contentIds[image.Filename] = image.Filename + "@" + strconv.Itoa(counter)
			body = strings.Replace(body, image.Filename, "cid:"+contentIds[image.Filename], -1)
		}
	}

	htmlBodyP.Write([]byte(body))

	return contentIds, nil
}

func writeInlineImages(relatedBodyW *multipart.Writer, images map[string]*attachments.EmbeddedBinaryObject, contentIds map[string]string, currentTime time.Time) error {
	for _, image := range images {
		embeddedImageHeaders := make(textproto.MIMEHeader)
		embeddedImageHeaders.Add(headers.ContentType(image.MIMEType, params.StringValue("name", image.Filename)))
		embeddedImageHeaders.Add(headers.ContentDescription(image.Filename))
		embeddedImageHeaders.Add(headers.ContentDisposition(params.DispInline, params.StringValue("filename", image.Filename), params.IntValue("size", image.Data.Len()), params.DateValue("creation-date", currentTime), params.DateValue("modification-date", currentTime)))
		embeddedImageHeaders.Add(headers.ContentId(contentIds[image.Filename]))
		embeddedImageHeaders.Add(headers.ContentTransferEncoding(encoding.Base64))
		embeddedImageP, err := relatedBodyW.CreatePart(embeddedImageHeaders)

		if err != nil {
			return err
		}

		embeddedImageP.Write(image.Data.Bytes())
	}

	return nil
}

func writeAttachments(w *multipart.Writer, attachments map[string]*attachments.EmbeddedBinaryObject, currentTime time.Time) error {
	for _, attachment := range attachments {
		attachmentHeaders := make(textproto.MIMEHeader)
		attachmentHeaders.Add(headers.ContentType(attachment.MIMEType, params.StringValue("name", attachment.Filename)))
		attachmentHeaders.Add(headers.ContentDescription(attachment.Filename))
		attachmentHeaders.Add(headers.ContentDisposition(params.DispAttachment, params.StringValue("filename", attachment.Filename), params.IntValue("size", attachment.Data.Len()), params.DateValue("creation-date", currentTime), params.DateValue("modification-date", currentTime)))
		attachmentHeaders.Add(headers.ContentTransferEncoding(encoding.Base64))
		attachmentP, err := w.CreatePart(attachmentHeaders)

		if err != nil {
			return err
		}

		attachmentP.Write(attachment.Data.Bytes())
	}

	return nil
}

type emailBody interface {
	WriteContent(io.Writer, *multipart.Writer) error
	Close(io.Writer) error
}

type emailPart interface {
	WriteFull(io.Writer, *multipart.Writer) error
	Close(io.Writer) error
}

type container struct {
	MIMEType mime.Type
	Children []emailPart
}

// region textContainer is a container for plain text emails.

type textContainer struct {
	TextBody *body.MailBody
	container
}

func (c *textContainer) WriteContent(p io.Writer, _ *multipart.Writer) error {
	p.Write([]byte(c.TextBody.GenerateBody()))

	return nil
}

func (c *textContainer) WriteFull(_ io.Writer, pw *multipart.Writer) error {
	h := make(textproto.MIMEHeader)
	h.Add(headers.ContentType(c.MIMEType))

	p, err := pw.CreatePart(h)

	if err != nil {
		return err
	}

	return c.WriteContent(p, nil)
}

func (c *textContainer) Close(_ io.Writer) error {
	return nil
}

// endregion

// region htmlContainer is a container for the HTML portion of an HTML email. The images are stored in a separate map.

type htmlContainer struct {
	HTMLBody   *body.MailBody
	ContentIds map[string]string
	Images     map[string]*attachments.EmbeddedBinaryObject
	container
}

func (c *htmlContainer) WriteContent(p io.Writer, _ *multipart.Writer) error {
	counter := 10000

	htmlBody := c.HTMLBody.GenerateBody()

	if len(c.Images) > 0 {
		for _, image := range c.Images {
			c.ContentIds[image.Filename] = image.Filename + "@" + strconv.Itoa(counter)
			htmlBody = strings.Replace(htmlBody, image.Filename, "cid:"+c.ContentIds[image.Filename], -1)
		}
	}

	p.Write([]byte(htmlBody))

	return nil
}

func (c *htmlContainer) WriteFull(_ io.Writer, pw *multipart.Writer) error {
	h := make(textproto.MIMEHeader)
	h.Add(headers.ContentType(c.MIMEType))
	p, err := pw.CreatePart(h)

	if err != nil {
		return err
	}

	return c.WriteContent(p, nil)
}

func (c *htmlContainer) Close(_ io.Writer) error {
	return nil
}

// endregion

// region alternativeContainer is a container for emails that both have plain text and HTML components.

type alternativeContainer struct {
	container
}

func (c *alternativeContainer) WriteContent(buf io.Writer, w *multipart.Writer) error {
	for _, child := range c.Children {
		child.WriteFull(buf, w)
	}

	w.Close()

	return nil
}

func (c *alternativeContainer) WriteFull(buf io.Writer, pw *multipart.Writer) error {
	w := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	h.Add(headers.ContentType(mime.MultipartAlternative, params.StringValue("boundary", w.Boundary())))

	pw.CreatePart(h)

	return c.WriteContent(buf, w)
}

func (c *alternativeContainer) Close(_ io.Writer) error {
	return nil
}

// endregion

// region relatedContainer is a container for emails that have a HTML part with inline images.

type relatedContainer struct {
	ContentIds map[string]string
	Images     map[string]*attachments.EmbeddedBinaryObject
	Time       time.Time
	container
}

func (c *relatedContainer) WriteContent(buf io.Writer, w *multipart.Writer) error {
	for _, child := range c.Children {
		child.WriteFull(buf, w)
	}

	err := writeInlineImages(w, c.Images, c.ContentIds, c.Time)

	if err != nil {
		return err
	}

	w.Close()

	return nil
}

func (c *relatedContainer) WriteFull(buf io.Writer, pw *multipart.Writer) error {
	w := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	h.Add(headers.ContentType(mime.MultipartRelated, params.StringValue("boundary", w.Boundary())))

	pw.CreatePart(h)

	return c.WriteContent(buf, w)
}

func (c *relatedContainer) Close(_ io.Writer) error {
	return nil
}

// endregion

// region mixedContainer is a container for emails that have attachments.

type mixedContainer struct {
	Attachments map[string]*attachments.EmbeddedBinaryObject
	Time        time.Time
	container
}

func (c *mixedContainer) WriteContent(buf io.Writer, w *multipart.Writer) error {
	for _, child := range c.Children {
		child.WriteFull(buf, w)
	}

	err := writeAttachments(w, c.Attachments, c.Time)

	if err != nil {
		return err
	}

	w.Close()

	return nil
}

func (c *mixedContainer) WriteFull(buf io.Writer, w *multipart.Writer) error {
	h := make(textproto.MIMEHeader)
	h.Add(headers.ContentType(mime.MultipartMixed, params.StringValue("boundary", w.Boundary())))

	err := serialiseHeaders(buf, h)

	if err != nil {
		return err
	}

	return c.WriteContent(buf, w)
}

func (c *mixedContainer) Close(_ io.Writer) error {
	return nil
}

// endregion

func newContainer() container {
	return container{
		Children: make([]emailPart, 0),
	}
}

// GenerateMessage will create a buffer containing the email message in it's current state.
func (m *Message) GenerateMessage() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

	var main emailBody
	var mainContentType mime.Type
	var text *textContainer
	var html *htmlContainer
	var alternative *alternativeContainer
	var related *relatedContainer
	var mixed *mixedContainer
	contentIds := make(map[string]string)
	currentTime := time.Now()

	if m.textBody != nil {
		text = &textContainer{
			TextBody:  m.textBody,
			container: newContainer(),
		}
		text.MIMEType = mime.TextPlain

		main = text
		mainContentType = text.MIMEType
	}

	if m.htmlBody != nil {
		html = &htmlContainer{
			ContentIds: contentIds,
			HTMLBody:   m.htmlBody,
			Images:     m.images,
		}
		html.MIMEType = mime.TextHTML

		main = html
		mainContentType = html.MIMEType
	}

	if m.textBody != nil && m.htmlBody != nil {
		alternative = &alternativeContainer{}
		alternative.MIMEType = mime.MultipartAlternative
		alternative.Children = append(alternative.Children, text, html)

		main = alternative
		mainContentType = alternative.MIMEType
	}

	if len(m.images) > 0 {
		related = &relatedContainer{
			Images:     m.images,
			ContentIds: contentIds,
			Time:       currentTime,
		}
		related.MIMEType = mime.MultipartRelated

		if m.textBody != nil {
			related.Children = append(related.Children, alternative)
		} else {
			related.Children = append(related.Children, html)
		}

		main = related
		mainContentType = related.MIMEType
	}

	if len(m.attachments) > 0 {
		mixed = &mixedContainer{
			Time:        currentTime,
			Attachments: m.attachments,
		}
		mixed.MIMEType = mime.MultipartMixed

		if related != nil {
			mixed.Children = append(mixed.Children, related)
		} else if alternative != nil {
			mixed.Children = append(mixed.Children, alternative)
		} else if html != nil {
			mixed.Children = append(mixed.Children, html)
		} else {
			mixed.Children = append(mixed.Children, text)
		}

		main = mixed
		mainContentType = mixed.MIMEType
	}

	m.AddMailHeader(headers.ContentType(mainContentType, params.StringValue("boundary", w.Boundary())))

	err := serialiseHeaders(buf, m.headers)

	if err != nil {
		return nil, err
	}

	main.WriteContent(buf, w)

	if err != nil {
		return nil, errors.Wrap(err, "There was an issue serialising the headers.")
	}

	return buf, nil
}
