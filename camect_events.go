package camect_go

import (
	"log/slog"

	"github.com/kaiiorg/camect-go/events"
)

type CamectEvents struct {
	// AlertChan receives detection alerts from the hub
	AlertChan chan *events.Alert
	// ModeChangeChan receives hub mode (DEFAULT or HOME) changes from the hub
	ModeChangeChan chan *events.ModeChange
	// AlertDisabledChan receives camera alert disabled messages from the hub
	AlertDisabledChan chan *events.AlertDisabled
	// AlertEnabledChan receives camera alert disabled messages from the hub
	AlertEnabledChan chan *events.AlertEnabled
	// CameraOnlineChan receives camera alert disabled messages from the hub
	CameraOnlineChan chan *events.CameraOnline
	// CameraOfflineChan receives camera alert disabled messages from the hub
	CameraOfflineChan chan *events.CameraOffline
	// UnknownEventChan receives unmarshalled JSON when an unknown event type is sent from the hub.
	UnknownEventChan chan []byte
}

func newCamectEvent(bufferSize int) *CamectEvents {
	return &CamectEvents{
		AlertChan:         make(chan *events.Alert, bufferSize),
		ModeChangeChan:    make(chan *events.ModeChange, bufferSize),
		AlertDisabledChan: make(chan *events.AlertDisabled, bufferSize),
		AlertEnabledChan:  make(chan *events.AlertEnabled, bufferSize),
		CameraOnlineChan:  make(chan *events.CameraOnline, bufferSize),
		CameraOfflineChan: make(chan *events.CameraOffline, bufferSize),
		UnknownEventChan:  make(chan []byte, bufferSize),
	}
}

func (ce *CamectEvents) close() {
	close(ce.AlertChan)
	close(ce.ModeChangeChan)
	close(ce.AlertDisabledChan)
	close(ce.AlertEnabledChan)
	close(ce.CameraOnlineChan)
	close(ce.CameraOfflineChan)
	close(ce.UnknownEventChan)
}

func sendEvent[T any](event T, eventChan chan T, eventType events.Type, log *slog.Logger) {
	// Don't push the event into the channel if there's no space (so we don't block)
	if cap(eventChan) > 0 {
		eventChan <- event
	} else {
		log.Warn("channel backed up; dropping event", "type", eventType.String())
	}
}
