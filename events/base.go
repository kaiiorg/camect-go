package events

import "encoding/json"

type Base struct {
	Type Type `json:"type"`

	Description      string           `json:"desc"`
	CamId            string           `json:"cam_id"`
	CamName          string           `json:"cam_name"`
	Url              string           `json:"url"`
	DetectedObjects  []string         `json:"detected_obj"`
	RegionOfInterest RegionOfInterest `json:"roi"`

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
