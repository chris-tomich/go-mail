package attachments

import "io"
import "bytes"

// EmbeddedBinaryObject contains the data necessary to attach a file or embed an image into an email.
type EmbeddedBinaryObject struct {
	// FileName is the name of the attachment or image.
	FileName string
	// Data is the base64 encoded data for this attachment or image.
	Data bytes.Buffer
}

// FileAttachmentFromFile will return a EmbeddedBinaryObject based upon the file path provided.
func FileAttachmentFromFile(filePath string) *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{}
}

// FileAttachmentFromReader will return a EmbeddedBinaryObject based upon the data provided and provided filename.
func FileAttachmentFromReader(filename string, reader io.Reader) *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{}
}

// InlineImageFromPath will return an EmbeddedBinaryObject based upon the file path provided. The name of the image
// will be inferred from it's filename.
func InlineImageFromPath(imagePath string) *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{}
}

// InlineImageFromFile will return an EmbeddedBinaryObject based upon the file path provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromFile(imageName string, imagePath string) *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{}
}

// InlineImageFromReader will return an EmbeddedBinaryObject based upon the data provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromReader(imageName string, reader io.Reader) *EmbeddedBinaryObject {
	return &EmbeddedBinaryObject{}
}
