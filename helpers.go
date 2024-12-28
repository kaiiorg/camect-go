package camect_go

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

func (h *Hub) auth() string {
	return base64.StdEncoding.EncodeToString([]byte(h.username + ":" + h.password))
}

func (h *Hub) request(method, path string, params url.Values) (int, []byte, error) {
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
		return 0, nil, err
	}
	req.SetBasicAuth(h.username, h.password)

	h.logger.Debug("sending request")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	h.logger.Debug("reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}
