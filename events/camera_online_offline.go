package events

type CameraOnline struct {
	CamId   string `json:"cam_id"`
	CamName string `json:"cam_name"`
}

func (b *Base) CameraOnline() *CameraOnline {
	return &CameraOnline{
		CamId:   b.CamId,
		CamName: b.CamName,
	}
}

type CameraOffline struct {
	CamId   string `json:"cam_id"`
	CamName string `json:"cam_name"`
}

func (b *Base) CameraOffline() *CameraOffline {
	return &CameraOffline{
		CamId:   b.CamId,
		CamName: b.CamName,
	}
}
