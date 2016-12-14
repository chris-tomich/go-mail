package attachments

import "io"

// FileAttachment contains the data necessary to send a file as an attachment.
type FileAttachment struct{}

// InlineImage contains the data necessary to embed an image into an email with a HTML body.
type InlineImage struct{}

// FileAttachmentFromFile will return a FileAttachment based upon the file path provided.
func FileAttachmentFromFile(filePath string) *FileAttachment {
	return &FileAttachment{}
}

// FileAttachmentFromReader will return a FileAttachment based upon the data provided and provided filename.
func FileAttachmentFromReader(filename string, reader io.Reader) *FileAttachment {
	return &FileAttachment{}
}

// InlineImageFromFile will return an InlineImage based upon the file path provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromFile(imageName string, imagePath string) *InlineImage {
	return &InlineImage{}
}

// InlineImageFromReader will return an InlineImage based upon the data provided. The imageName is
// the value used in the corresponding HTML body.
func InlineImageFromReader(imageName string, reader io.Reader) *InlineImage {
	return &InlineImage{}
}
