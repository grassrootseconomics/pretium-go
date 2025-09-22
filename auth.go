package pretium

import (
	"net/http"
)

func (fc *PretiumClient) setAuthHeaders(req *http.Request) error {
	req.Header.Set("x-api-key", fc.apiKey)
	return nil
}
