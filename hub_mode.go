package camect_go

type HubMode string

func (hm HubMode) String() string {
	return string(hm)
}

const (
	ModeDefault = "DEFAULT"
	ModeHome    = "HOME"

	// ModeDisarmed is what the Camect V2 UI shows ModeDefault as
	ModeDisarmed = ModeDefault
	// ModeNormal is what the Camect V2 UI shows ModeHome as
	ModeNormal = ModeHome
)
