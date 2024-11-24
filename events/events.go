package events

type ModeChange struct {
	// Description is a Type dependent description
	// Type == ModeEvent: what mode the hub has been switched to
	Description string `json:"desc"`
}

func (b *Base) ModeChange() *ModeChange {
	return &ModeChange{
		Description: b.Description,
	}
}

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
