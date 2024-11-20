package camect_go

import (
	"crypto/tls"
	"log/slog"
	"net/http"
)

type Hub struct {
	ip       string
	username string
	password string

	logger     *slog.Logger
	httpClient http.Client
}

func New(ip, username, password string, logger *slog.Logger) *Hub {
	if logger == nil {
		logger = slog.Default()
	}

	return &Hub{
		ip:       ip,
		username: username,
		password: password,

		logger: logger,
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
