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
