package CheckProvider

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func GetEmailProvider(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid Email Format")
	}

	domain := parts[1]
	var provider string
	var err error

	dnsServers := []string{"8.8.8.8:53", "8.8.4.4:53", "1.1.1.1:53", "1.0.0.1:53"}

	for retries := 0; retries < len(dnsServers); retries++ {
		provider, err = getMXRecordProvider(domain, dnsServers[retries])
		if err == nil {
			break
		}

		log.Printf("Error getting email provider for %s using DNS server %s (attempt %d): %v", email, dnsServers[retries], retries+1, err)
		time.Sleep(1 * time.Second) // Wait for 1 second before retrying

	}

	if err != nil {
		return "", err
	}

	return provider, nil
}

func getMXRecordProvider(domain, dnsServer string) (string, error) {
	mxRecords, err := getMXRecords(domain, dnsServer)
	if err != nil {
		return "", err
	}

	if len(mxRecords) == 0 {
		return "", fmt.Errorf("no MX records found for %s", domain)
	}

	provider := mxRecords[0].Mx
	return provider, nil
}

func getMXRecords(domain, dnsServer string) ([]*dns.MX, error) {
	c := dns.Client{
		ReadTimeout: 5 * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)

	r, _, err := c.Exchange(&m, dnsServer)
	if err != nil {
		return nil, err
	}

	if len(r.Answer) == 0 {
		return nil, fmt.Errorf("no Records found for given DNS query")
	}

	mxRecords := make([]*dns.MX, 0, len(r.Answer))
	for _, ans := range r.Answer {
		if mx, ok := ans.(*dns.MX); ok {
			mxRecords = append(mxRecords, mx)
		}
	}

	sort.Slice(mxRecords, func(i, j int) bool {
		return mxRecords[i].Preference < mxRecords[j].Preference
	})

	return mxRecords, nil

}
