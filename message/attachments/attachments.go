package attachments

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Represents a buffer that holds base64 encoded information.
type Base64Buffer struct {
	*bytes.Buffer
}

// EmbeddedBinaryObject contains the data necessary to attach a file or embed an image into an email.
type EmbeddedBinaryObject struct {
	// FileName is the name of the attachment or image.
	FileName string
	// Data is the base64 encoded data for this attachment or image.
	Data Base64Buffer
}

func NewEmbeddedBinaryObject() *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{
		Data: Base64Buffer{
			Buffer: &bytes.Buffer{},
		},
	}
}

func loadEmbeddedBinaryObject(filename string, b []byte) (*EmbeddedBinaryObject, error) {
	o := NewEmbeddedBinaryObject()
	o.FileName = filename

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

	return loadEmbeddedBinaryObject(filename, b)
}

// FileAttachmentFromReader will return a EmbeddedBinaryObject based upon the data provided and provided filename.
// Filename is the name of the file to be attached, not the path to a file.
func FileAttachmentFromReader(filename string, reader io.Reader) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return loadEmbeddedBinaryObject(filename, b)
}

// InlineImageFromPath will return an EmbeddedBinaryObject based upon the file path provided. The name of the image
// will be inferred from it's filename.
func InlineImageFromPath(imagePath string) (*EmbeddedBinaryObject, error) {
	_, imageName := filepath.Split(imagePath)

	return InlineImageFromFile(imageName, imagePath)
}

// InlineImageFromFile will return an EmbeddedBinaryObject based upon the file path provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromFile(imageName string, imagePath string) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadFile(imagePath)

	if err != nil {
		return nil, err
	}

	return loadEmbeddedBinaryObject(imageName, b)
}

// InlineImageFromReader will return an EmbeddedBinaryObject based upon the data provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromReader(imageName string, reader io.Reader) (*EmbeddedBinaryObject, error) {
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return loadEmbeddedBinaryObject(imageName, b)
}
