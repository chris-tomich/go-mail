package mime

import "fmt"

// Multipart is a MIME type specifically for 'multipart/alternative' and 'multipart/mixed' that has additional settings
type Multipart Type

const (
	// MultipartAlternative is the MIME type for 'multipart/alternative'
	MultipartAlternative Multipart = "multipart/alternative"
	// MultipartMixed is the MIME type for 'multipart/mixed'
	MultipartMixed Multipart = "multipart/mixed"
	// MultipartRelated is the MIME type for 'multipart/related'
	MultipartRelated Multipart = "multipart/related"
)

// Type returns this Multipart as a Type.
func (m Multipart) Type() Type {
	return Type(m)
}

// SetBoundary set's the boundary setting for a 'multipart/alternative' or 'multipart/mixed' MIME type.
func (m Multipart) SetBoundary(boundary string) Type {
	return Type(fmt.Sprintf("%v; boundary=%v", m, boundary))
}
