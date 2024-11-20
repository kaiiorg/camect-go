package camect_go

import (
	"io"
	"net/http"
	"net/url"
)

func (h *Hub) request(method, path string, params url.Values) ([]byte, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     h.ip,
		Path:     path,
		RawQuery: params.Encode(),
	}
	url := u.String()
	h.logger.Debug(
		"building request",
		"method", method,
		"url", url,
		"username", h.username,
		"password", "[REDACTED]",
	)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(h.username, h.password)

	h.logger.Debug("sending request")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	h.logger.Debug("reading response body")
	return io.ReadAll(resp.Body)
}
