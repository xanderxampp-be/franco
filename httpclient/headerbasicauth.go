package httpclient

import (
	"encoding/base64"
	"net/http"
)

func GenerateHeaderBasicAuthGeneric(timestamp, usr, pwd string) http.Header {
	contentType := "application/json"
	headers := http.Header{
		"Content-Type":  []string{contentType},
		"Request-Time":  []string{timestamp},
		"Authorization": []string{basicAuth(usr, pwd)},
	}

	return headers
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
