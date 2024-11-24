package events

type Alert struct {
	Description      string           `json:"desc"`
	Url              string           `json:"url"`
	CamId            string           `json:"cam_id"`
	CamName          string           `json:"cam_name"`
	DetectedObjects  []string         `json:"detected_obj"`
	RegionOfInterest RegionOfInterest `json:"roi"`
}

func (b *Base) Alert() *Alert {
	return &Alert{
		Description:      b.Description,
		Url:              b.Url,
		CamId:            b.CamId,
		CamName:          b.CamName,
		DetectedObjects:  b.DetectedObjects,
		RegionOfInterest: b.RegionOfInterest,
	}
}

type RegionOfInterest struct {
	Contour []Contour `json:"contour"`
	Object  []Object  `json:"object"`
}

type Contour struct {
	Point []Point `json:"point"`
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Object struct {
	Name          string    `json:"name"`
	MinSize       float32   `json:"min_size"`
	MaxSize       float32   `json:"max_size"`
	Contour       []Contour `json:"contour"`
	MovementTrace Contour   `json:"movement_trace"`
}
