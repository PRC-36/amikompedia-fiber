package util

import (
	"fmt"
	"strings"
)

type UsernameAndUUID struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

func FromStringToUsernameAndUUID(usernameAndId string) (*UsernameAndUUID, error) {
	parts := strings.Split(usernameAndId, ":")

	// Check if we have at least two parts
	if len(parts) >= 2 {
		// Assign the values to variables a and b
		a := parts[0]
		b := parts[1]

		return &UsernameAndUUID{Username: a, UserID: b}, nil
	} else {
		return nil, fmt.Errorf("invalid input string format")
	}
}

func FromUsernameAndUUIDToString(username, uuid string) string {
	return fmt.Sprintf("%s:%s", username, uuid)
}
