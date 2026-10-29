package handlers

import (
	"regexp"
	"strings"
	"unicode"
)

// Field length limits
const (
	maxNicknameLen = 30
	minNicknameLen = 3
	maxNameLen     = 50
	maxEmailLen    = 254
	maxPasswordLen = 72 // bcrypt hard limit
	minPasswordLen = 8
	minAge         = 13
	maxAge         = 120
	maxPostTitle   = 200
	minPostTitle   = 3
	maxPostContent = 10000
	maxCategories  = 5
	maxCategoryLen = 50
	maxCommentLen  = 2000
	minCommentLen  = 1
	maxMessageLen  = 2000
	maxIdentifier  = 254
)

var (
	// Matches the frontend regex: /^[^\s@]+@[^\s@]+\.[^\s@]+$/
	emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	// nickname: letters, digits, underscores, hyphens only
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	// digit check — mirrors the frontend /\d/ test used on names
	digitRegex = regexp.MustCompile(`\d`)
	// UUID v4 format
	uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// validateEmail checks basic email format and length.
func validateEmail(email string) string {
	if len(email) > maxEmailLen {
		return "email must not exceed 254 characters"
	}
	if !emailRegex.MatchString(email) {
		return "email address is not valid"
	}
	return ""
}

// validateNickname checks length and allowed characters.
func validateNickname(nickname string) string {
	if len(nickname) < minNicknameLen {
		return "nickname must be at least 3 characters"
	}
	if len(nickname) > maxNicknameLen {
		return "nickname must not exceed 30 characters"
	}
	if !nicknameRegex.MatchString(nickname) {
		return "nickname may only contain letters, digits, underscores and hyphens"
	}
	return ""
}

// validateName checks a first/last name field — mirrors frontend /\d/ check.
func validateName(field, value string) string {
	if len(value) < 1 {
		return field + " is required"
	}
	if len(value) > maxNameLen {
		return field + " must not exceed 50 characters"
	}
	if digitRegex.MatchString(value) {
		return field + " must not contain numbers"
	}
	return ""
}

// validatePassword checks length and character requirements — mirrors frontend rules.
func validatePassword(password string) string {
	if len(password) < minPasswordLen {
		return "password must be at least 8 characters"
	}
	if len(password) > maxPasswordLen {
		return "password must not exceed 72 characters"
	}
	hasLower, hasUpper := false, false
	for _, r := range password {
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsUpper(r) {
			hasUpper = true
		}
	}
	if !hasLower {
		return "password must contain at least one lowercase letter"
	}
	if !hasUpper {
		return "password must contain at least one uppercase letter"
	}
	return ""
}

// validateGender ensures the value is one of the three options the frontend sends.
func validateGender(gender string) string {
	switch strings.ToLower(gender) {
	case "male", "female", "other":
		return ""
	}
	return "gender must be one of: male, female, other"
}

// isValidUUID returns true if s is a valid UUID v4.
func isValidUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

// validatePostTitle checks post title length.
func validatePostTitle(title string) string {
	if len(title) < minPostTitle {
		return "title must be at least 3 characters"
	}
	if len(title) > maxPostTitle {
		return "title must not exceed 200 characters"
	}
	return ""
}

// validatePostContent checks post body length.
func validatePostContent(content string) string {
	if len(content) > maxPostContent {
		return "post content must not exceed 10,000 characters"
	}
	return ""
}

// validateCategories checks the list of categories.
func validateCategories(cats []string) string {
	if len(cats) == 0 {
		return "at least one category is required"
	}
	if len(cats) > maxCategories {
		return "you may select at most 5 categories"
	}
	for _, c := range cats {
		if len(c) < 1 || len(c) > maxCategoryLen {
			return "each category must be between 1 and 50 characters"
		}
	}
	return ""
}

// validateCommentContent checks comment length.
func validateCommentContent(content string) string {
	if len(content) < minCommentLen {
		return "comment must not be empty"
	}
	if len(content) > maxCommentLen {
		return "comment must not exceed 2,000 characters"
	}
	return ""
}

// validateMessageContent checks chat message length.
func validateMessageContent(content string) string {
	if len(content) > maxMessageLen {
		return "message must not exceed 2,000 characters"
	}
	return ""
}
