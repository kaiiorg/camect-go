package events

import "encoding/json"

type Base struct {
	Type Type `json:"type"`

	// Description is a Type dependent description
	// Type == ModeEvent: what mode the hub has been switched to
	Description string `json:"desc"`
}

func New(data []byte) (*Base, error) {
	b := &Base{}
	err := json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
