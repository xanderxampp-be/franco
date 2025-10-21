package httpclient

import (
	"net/http"
)

func GenerateHeaderBasic() http.Header {
	contentType := "application/json"
	headers := http.Header{
		"Content-Type": []string{contentType},
	}

	return headers
}
