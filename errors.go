package camect_go

import "errors"

var (
	ErrAlreadyListeningForEvents = errors.New("already listening for events")
	ErrFailedToSetMode           = errors.New("failed to set mode")
	ErrReasonNotProvided         = errors.New("reason was not provided")
	ErrFailedToGetCamera         = errors.New("failed to get cameras")
	ErrUnexpectedDataShape       = errors.New("the hub returned JSON in an unexpected shape")
)
