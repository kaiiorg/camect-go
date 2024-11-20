package camect_go

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/coder/websocket"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

type Hub struct {
	ip       string
	username string
	password string

	ctx context.Context
	wg  sync.WaitGroup

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

func (h *Hub) Events(ctx context.Context, buffer int) (<-chan string, error) {
	if h.ctx == nil {
		h.ctx = ctx
	} else {
		return nil, AlreadyListeningForEvents
	}

	if buffer <= 0 {
		h.logger.Warn("buffer must be > 0. Set to 1")
		buffer = 1
	}

	conn, err := h.connectWs()
	if err != nil {
		return nil, err
	}

	eventsChan := make(chan string, buffer)
	go h.eventListener(eventsChan, conn)

	return eventsChan, nil
}

func (h *Hub) connectWs() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   h.ip,
		Path:   eventWsPath,
	}
	conn, _, err := websocket.Dial(
		h.ctx,
		u.String(),
		&websocket.DialOptions{
			HTTPClient: &h.httpClient,
			HTTPHeader: http.Header{
				"Authorization": []string{fmt.Sprintf("Basic %s", h.auth())},
			},
		},
	)

	return conn, err
}

func (h *Hub) eventListener(eventsChan chan string, conn *websocket.Conn) {
	type webSocketMessage struct {
		Type websocket.MessageType
		Data []byte
		Err  error
	}
	messageChan := make(chan *webSocketMessage, 1)
	go func() {
		for {
			m := &webSocketMessage{}
			m.Type, m.Data, m.Err = conn.Read(h.ctx)
			messageChan <- m
			if m.Err != nil {
				break
			}
		}
	}()

	for {
		select {
		case <-h.ctx.Done():
			conn.Close(websocket.StatusGoingAway, "asked to exit")
			return
		case m := <-messageChan:
			if m.Err != nil {
				h.logger.Error(m.Err.Error())
				return
			}
			h.logger.Info(
				"Got message",
				"type", m.Type,
				"data", string(m.Data),
			)
		}
	}

}

func (h *Hub) Info() (*Info, error) {
	rawInfo, err := h.request(http.MethodGet, homeInfoPath, nil)
	if err != nil {
		return nil, err
	}

	return infoFromJson(rawInfo)
}
