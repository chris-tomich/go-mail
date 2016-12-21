package headers

import (
	"strconv"
	"strings"

	"github.com/chris-tomich/go-mail/message/headers/encoding"
	"github.com/chris-tomich/go-mail/message/headers/mime"
	"fmt"
	"github.com/chris-tomich/go-mail/message/headers/params"
)

func aggregateAddresses(addresses ...string) string {
	var fromAddresses string

	for _, address := range addresses {
		address = strings.Trim(address, " ")

		fromAddresses += "; " + address
	}

	return fromAddresses[2:]
}

// From returns a properly formatted header for the given addresses.
// The address can be of the following formats.
// Simple Email - chris.tomich@email.com
// Email with Name - Chris Tomich <chris.tomich@email.com>
func From(address string) (string, string) {
	return "From", address
}

// To returns a properly formatted header for the given addresses.
// The addresses can be of the following formats.
// Simple Email - chris.tomich@email.com
// Email with Name - Chris Tomich <chris.tomich@email.com>
func To(addresses ...string) (string, string) {
	return "To", aggregateAddresses(addresses...)
}

// Subject returns a properly formatted header for the given subject line.
func Subject(subject string) (string, string) {
	return "Subject", subject
}

// MIMEVersion returns a "MIME-Version" header with the given version.
func MIMEVersion(version float64) (string, string) {
	return "MIME-Version", strconv.FormatFloat(version, 'f', 1, 64)
}

// ContentType returns a "Content-Type" header with the given MIME type.
func ContentType(mimeType mime.Type, additionalParams ...params.Header) (string, string) {
	if len(additionalParams) > 0 {
		return "Content-Type", fmt.Sprintf("%v%v", mimeType, params.Aggregate(additionalParams))
	} else {
		return "Content-Type", string(mimeType)
	}
}

func ContentDescription(description string) (string, string) {
	return "Content-Description", description
}

func ContentDisposition(disposition params.ContentDisposition, additionalParams ...params.Header) (string, string) {
	return "Content-Disposition", fmt.Sprintf("%v%v", disposition, params.Aggregate(additionalParams))
}

func ContentId(contentId string) (string, string) {
	return "Content-ID", fmt.Sprintf("<%v>", contentId)
}

// ContentTransferEncoding returns a "Content-Transfer-Encoding" header with the provided encoding type.
// If no "Content-Transfer-Encoding" header is provided, according to RFC1341, 7bit will be assumed by clients.
func ContentTransferEncoding(encoding encoding.ContentTransferType) (string, string) {
	return "Content-Transfer-Encoding", string(encoding)
}
