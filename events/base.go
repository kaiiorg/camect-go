package events

import "encoding/json"

type Base struct {
	Type Type `json:"type"`

	// Description is a Type dependent description
	// Type == ModeEvent: what mode the hub has been switched to
	Description string `json:"desc"`

	// CamId camera ID
	CamId string `json:"cam_id"`
	// CamName camera name
	CamName string `json:"cam_name"`

	// raw is the original JSON
	raw []byte
}

func New(data []byte) (*Base, error) {
	b := &Base{
		raw: data,
	}
	err := json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Base) Raw() []byte {
	return b.raw
}
