package version

import "time"

type VersionSpec struct {
	Version       string    `json:"version"`
	Channel       string    `json:"channel"`
	ReleasedAt    time.Time `json:"released_at"`
	Description   string    `json:"description"`
	CanAutoUpdate bool      `json:"can_auto_update"`
	URL           string    `json:"url"`
}
