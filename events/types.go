package events

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	ModeEvent Type = "mode"
)
