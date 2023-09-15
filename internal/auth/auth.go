package auth

import (
	"errors"
	"net/http"
	"strings"
)

// APIKey extracts API key from
// the headers of HTTP request
// Example:
// Authorization: APIKey {insert api_key here}
func APIKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "APIKey" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}