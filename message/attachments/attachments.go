package attachments

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"path/filepath"
	"github.com/chris-tomich/go-mail/message/headers/mime"
)

// Represents a buffer that holds base64 encoded information.
type Base64Buffer struct {
	*bytes.Buffer
}

// EmbeddedBinaryObject contains the data necessary to attach a file or embed an image into an email.
type EmbeddedBinaryObject struct {
	// Filename is the name of the attachment or image.
	Filename string
	// Data is the base64 encoded data for this attachment or image.
	Data     Base64Buffer
	// MimeType stores the MIME type detected for this file based upon it's filename extension.
	MIMEType mime.Type
}

func NewEmbeddedBinaryObject() *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{
		Data: Base64Buffer{
			Buffer: &bytes.Buffer{},
		},
	}
}

func loadEmbeddedBinaryObject(filename string, extension string, b []byte) (*EmbeddedBinaryObject, error) {
	o := NewEmbeddedBinaryObject()
	o.Filename = filename

	matchingMime := mime.Types[extension]

	if matchingMime == "" {
		o.MIMEType = mime.ApplicationOctetStream
	} else {
		o.MIMEType = mime.Type(matchingMime)
	}

	wc := base64.NewEncoder(base64.StdEncoding, o.Data)
	_, err := wc.Write(b)

	if err != nil {
		return nil, err
	}

	err = wc.Close()

	if err != nil {
		return nil, err
	}

	return o, nil
}

// FileAttachmentFromFile will return a EmbeddedBinaryObject based upon the file path provided.
func FileAttachmentFromFile(filePath string) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	_, filename := filepath.Split(filePath)
	ext := filepath.Ext(filePath)[1:]

	return loadEmbeddedBinaryObject(filename, ext, b)
}

// FileAttachmentFromReader will return a EmbeddedBinaryObject based upon the data provided and provided filename.
// Filename is the name of the file to be attached, not the path to a file.
func FileAttachmentFromReader(filename string, reader io.Reader) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)[1:]

	return loadEmbeddedBinaryObject(filename, ext, b)
}

// InlineImageFromPath will return an EmbeddedBinaryObject based upon the file path provided. The name of the image
// will be inferred from it's filename.
func InlineImageFromPath(imagePath string) (*EmbeddedBinaryObject, error) {
	_, imageName := filepath.Split(imagePath)

	return InlineImageFromFile(imageName, imagePath)
}

// InlineImageFromFile will return an EmbeddedBinaryObject based upon the file path provided. The imageName is
// the value used in the <img src='' /> tags in the corresponding HTML body. It is recommended this is the actual
// file's name so that the MIME type can be properly detected.
func InlineImageFromFile(imageName string, imagePath string) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadFile(imagePath)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(imageName)[1:]

	return loadEmbeddedBinaryObject(imageName, ext, b)
}

// InlineImageFromReader will return an EmbeddedBinaryObject based upon the data provided. The imageName is
// the value used in the <img src='' /> tags in the corresponding HTML body. It is recommended this is the actual
// file's name so that the MIME type can be properly detected.
func InlineImageFromReader(imageName string, reader io.Reader) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(imageName)[1:]

	return loadEmbeddedBinaryObject(imageName, ext, b)
}
