package camect_go

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/kaiiorg/camect-go/events"

	"github.com/coder/websocket"
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

	eventsChan := make(chan string, buffer)
	err := h.connectAndStartListener(eventsChan)
	if err != nil {
		return nil, err
	}

	return eventsChan, nil
}

func (h *Hub) connectAndStartListener(eventsChan chan string) error {
	conn, err := h.connectWs()
	if err != nil {
		return err
	}
	go h.eventListener(eventsChan, conn)
	return nil
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
	messageChan := make(chan *webSocketMessage, 1)
	// Anonymous goroutine to continuously read from the web socket connection
	// Exits when a read returns an error. Force this goroutine to exit by closing the connection
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-h.ctx.Done():
			conn.Close(websocket.StatusGoingAway, "asked to exit")
			return
		case <-ticker.C:
			err := conn.Ping(h.ctx)
			if err != nil {
				h.logger.Info("ping failed", "error", err)
			}
		case m := <-messageChan:
			if m.Err != nil {
				// TODO send error via channel
				h.logger.Error(m.Err.Error())

				if errors.Is(h.ctx.Err(), context.Canceled) {
					return
				}

				for {
					for i := 0; i < 5; i++ {
						if errors.Is(h.ctx.Err(), context.Canceled) {
							h.logger.Warn("asked to exit; stopping attempts to reconnect websocket")
							return
						}
						time.Sleep(time.Second)
					}

					err := h.connectAndStartListener(eventsChan)
					if err != nil {
						h.logger.Info("failed to reconnect, will retry in 5 seconds", "error", err)
						continue
					}
					h.logger.Info("reconnected websocket")
					return
				}
			}

			baseEvent, err := events.New(m.Data)
			if err != nil {
				h.logger.Error("failed to unmarshal event json", "error", err)
				continue
			}

			switch baseEvent.Type {
			case events.ModeEvent:
				// TODO send ModeChange events via channel
				h.logger.Info("got mode changed event", "data", fmt.Sprintf("%#v", baseEvent.ModeChange()))
			case events.AlertDisabledEvent:
				// TODO send AlertDisabledEvent events via channel
				h.logger.Info("got alert disabled event", "data", fmt.Sprintf("%#v", baseEvent.AlertDisabled()))
			case events.AlertEnabledEvent:
				// TODO send AlertEnabledEvent events via channel
				h.logger.Info("got alert enabled event", "data", fmt.Sprintf("%#v", baseEvent.AlertEnabled()))
			default:
				// TODO send unknown events via channel
				h.logger.Warn("got unknown event", "type", baseEvent.Type, "raw", string(baseEvent.Raw()))
			}
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
