package camect_go

import (
	"crypto/tls"
	"net/http"
)

type Hub struct {
	ip       string
	username string
	password string

	httpClient http.Client
}

func New(ip, username, password string) *Hub {
	return &Hub{
		ip:       ip,
		username: username,
		password: password,
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (h *Hub) Info() (*Info, error) {
	rawInfo, err := h.request(http.MethodGet, homeInfoPath, nil)
	if err != nil {
		return nil, err
	}

	return infoFromJson(rawInfo)
}
