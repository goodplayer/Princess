package validation

import "regexp"

func ValidateUserName(username string) bool {
	// [a-zA-Z0-9] and length >= 6 && length <= 20
	match, err := regexp.MatchString("[a-zA-Z0-9]+", username)
	if err != nil {
		panic(err)
	} else if !match {
		return false
	}
	return len(username) >= 6 && len(username) <= 20
}
