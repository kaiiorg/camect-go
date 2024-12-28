package camect_go

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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

func (h *Hub) Events(ctx context.Context, buffer int) (*CamectEvents, error) {
	if h.ctx == nil {
		h.ctx = ctx
	} else {
		return nil, ErrAlreadyListeningForEvents
	}

	if buffer <= 0 {
		h.logger.Warn("buffer must be > 0. Set to 1")
		buffer = 1
	}

	eventsChan := newCamectEvent(buffer)
	err := h.connectAndStartListener(eventsChan)
	if err != nil {
		return nil, err
	}

	return eventsChan, nil
}

func (h *Hub) Info() (*Info, error) {
	_, rawInfo, err := h.request(http.MethodGet, homeInfoPath, nil)
	if err != nil {
		return nil, err
	}

	return infoFromJson(rawInfo)
}

func (h *Hub) SetMode(newMode HubMode, reason string) error {
	params := url.Values{}
	params.Set("Mode", newMode.String())

	if reason != "" {
		params.Set("Reason", reason)
	}

	respCode, rawResp, err := h.request(http.MethodGet, setModePath, params)
	if err != nil {
		return err
	}

	if respCode == 200 {
		return nil
	}

	resp := map[string]interface{}{}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return errors.Join(ErrFailedToSetMode, err)
	}

	errMsg, ok := resp["err_msg"]
	if !ok {
		return errors.Join(ErrFailedToSetMode, ErrReasonNotProvided)
	}

	return errors.Join(ErrFailedToSetMode, fmt.Errorf("%#v", errMsg))
}

func (h *Hub) Cameras() ([]Camera, error) {
	respCode, rawResp, err := h.request(http.MethodGet, listCamerasPath, nil)
	if err != nil {
		return nil, err
	}

	if respCode != 200 {
		resp := map[string]interface{}{}
		err = json.Unmarshal(rawResp, &resp)
		if err != nil {
			return nil, errors.Join(ErrFailedToGetCamera, err)
		}

		errMsg, ok := resp["err_msg"]
		if !ok {
			return nil, errors.Join(ErrFailedToGetCamera, ErrReasonNotProvided)
		}

		return nil, errors.Join(ErrFailedToGetCamera, fmt.Errorf("%#v", errMsg))
	}

	cameraMap := map[string][]Camera{}
	err = json.Unmarshal(rawResp, &cameraMap)
	if err != nil {
		return nil, errors.Join(ErrFailedToGetCamera, err)
	}

	cameras, ok := cameraMap["camera"]
	if !ok {
		return nil, errors.Join(ErrFailedToGetCamera, ErrUnexpectedDataShape)
	}

	return cameras, nil
}
