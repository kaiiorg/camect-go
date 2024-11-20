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

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(h.username, h.password)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
