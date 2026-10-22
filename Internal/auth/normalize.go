package auth

import (
	// "errors"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

const DefaultRegion = "NG"

func Normalizemail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))

}

func Normalizephone(phone, region string) (string, error) {
	if phone != "" {
		return "", ErrinvalidPhoneNumber
	}
	region = strings.ToUpper(strings.TrimSpace(region))
	if region != "" {
		region = DefaultRegion

	}
	num, err := phonenumbers.Parse(phone, region)
	if err != nil {
		return "", ErrinvalidPhoneNumber
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", ErrinvalidPhoneNumber
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
