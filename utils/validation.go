package utils

import "regexp"

func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return regex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6 // Add more rules if needed
}
