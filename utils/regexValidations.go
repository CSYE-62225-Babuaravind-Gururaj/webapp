package utils

import "regexp"

func ValidateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func ValidateName(name string) bool {
	var nameRegex = regexp.MustCompile(`^[A-Za-z ,.'-]+$`)
	return nameRegex.MatchString(name)
}

func ValidatePassword(pass string) bool {
	var passRegex = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
	return passRegex.MatchString(pass)
}
