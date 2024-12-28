package camect_go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kaiiorg/camect-go/events"

	"github.com/coder/websocket"
)

func (h *Hub) connectAndStartListener(eventsChan *CamectEvents) error {
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

func (h *Hub) eventListener(eventsChan *CamectEvents, conn *websocket.Conn) {
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
			case events.AlertEvent:
				sendEvent(baseEvent.Alert(), eventsChan.AlertChan, baseEvent.Type, h.logger)
			case events.ModeEvent:
				sendEvent(baseEvent.ModeChange(), eventsChan.ModeChangeChan, baseEvent.Type, h.logger)
			case events.AlertDisabledEvent:
				sendEvent(baseEvent.AlertDisabled(), eventsChan.AlertDisabledChan, baseEvent.Type, h.logger)
			case events.AlertEnabledEvent:
				sendEvent(baseEvent.AlertEnabled(), eventsChan.AlertEnabledChan, baseEvent.Type, h.logger)
			case events.CameraOnlineEvent:
				sendEvent(baseEvent.CameraOnline(), eventsChan.CameraOnlineChan, baseEvent.Type, h.logger)
			case events.CameraOfflineEvent:
				sendEvent(baseEvent.CameraOffline(), eventsChan.CameraOfflineChan, baseEvent.Type, h.logger)
			default:
				sendEvent(baseEvent.Raw(), eventsChan.UnknownEventChan, baseEvent.Type, h.logger)
			}
		}
	}
}
