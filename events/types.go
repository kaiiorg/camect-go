package events

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	ModeEvent          Type = "mode"
	AlertDisabledEvent Type = "alert_disabled"
	AlertEnabledEvent  Type = "alert_enabled"
	AlertEvent         Type = "alert"
	CameraOnlineEvent  Type = "camera_online"
	CameraOfflineEvent Type = "camera_offline"
)
