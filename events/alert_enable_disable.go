package events

type AlertDisabled struct {
	CamId   string `json:"cam_id"`
	CamName string `json:"cam_name"`
}

func (b *Base) AlertDisabled() *AlertDisabled {
	return &AlertDisabled{
		CamId:   b.CamId,
		CamName: b.CamName,
	}
}

type AlertEnabled struct {
	CamId   string `json:"cam_id"`
	CamName string `json:"cam_name"`
}

func (b *Base) AlertEnabled() *AlertEnabled {
	return &AlertEnabled{
		CamId:   b.CamId,
		CamName: b.CamName,
	}
}
