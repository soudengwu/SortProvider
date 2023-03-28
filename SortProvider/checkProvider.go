package CheckProvider

import (
	"fmt"
	"net"
	"strings"
)

func GetEmailProvider(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid Email Format")
	}

	domain := parts[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return "", err
	}

	if len(mxRecords) == 0 {
		return "", fmt.Errorf("no MX records found for %s", domain)
	}

	provider := mxRecords[0].Host
	return provider, nil
}
