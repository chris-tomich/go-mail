package encoding

// ContentTransferType is a type suitable to be used for a "Content-Transfer-Encoding" header.
type ContentTransferType string

const (
	// 7Bit is the default value assumed by email clients if no "Content-Transfer-Encoding" field is set.
	// It also implies that no encoding has been performed on the data.
	_7Bit ContentTransferType = "7bit"
	// Base64 (taken from RFC1341) Content-Transfer-Encoding is designed to represent arbitrary sequences of octets in a form that is not humanly readable.
	Base64 ContentTransferType = "base64"
	// QuotedPrintable (taken from RFC1341) encoding is intended to represent data that largely consists of octets that correspond to printable characters in the ASCII character set.
	QuotedPrintable ContentTransferType = "quoted-printable"
	// 8Bit implies that no encoding has been performed on the data.
	_8Bit ContentTransferType = "8bit"
	// Binary implies that no encoding has been performed on the data.
	Binary ContentTransferType = "binary"
)
