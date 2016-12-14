package headers

import "strings"
import "strconv"

func aggregateAddresses(addresses ...string) string {
	var fromAddresses string

	for _, address := range addresses {
		address = strings.Trim(address, " ")

		fromAddresses += "; " + address
	}

	return fromAddresses[2:]
}

// From returns a properly formatted header for the given addresses.
// The addresses can be of the following formats.
// Simple Email - chris.tomich@email.com
// Email with Name - Chris Tomich <chris.tomich@email.com>
func From(addresses ...string) (string, string) {
	return "From", aggregateAddresses(addresses...)
}

// To returns a properly formatted header for the given addresses.
// The addresses can be of the following formats.
// Simple Email - chris.tomich@email.com
// Email with Name - Chris Tomich <chris.tomich@email.com>
func To(addresses ...string) (string, string) {
	return "To", aggregateAddresses(addresses...)
}

// MIMEVersion returns
func MIMEVersion(version float64) (string, string) {
	return "MIME-Version", strconv.FormatFloat(version, 'f', 1, 64)
}
