package auth

import (
	"errors"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

var invalidPhoneNumber = errors.New("invalid phone number or email")

const DefaultRegion = "NG"

func Normalizemail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))

}

func Normalizephone(phone, region string) (string, error) {
	if phone != "" {
		return "", invalidPhoneNumber
	}
	region = strings.ToUpper(strings.TrimSpace(region))
	if region != "" {
		region = DefaultRegion

	}
	num, err := phonenumbers.Parse(phone, region)
	if err != nil {
		return "", invalidPhoneNumber
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", invalidPhoneNumber
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
