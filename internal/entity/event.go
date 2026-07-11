package entity

import "time"

type SourceType string
type EventType string

const (
	SourceGitHub  SourceType = "github"
	SourceWeather SourceType = "weather"
	SourceReddit  SourceType = "reddit"
)

type Event struct {
	ID          string     `json:"id"`
	ExternalID  string     `json:"external_id"`
	Title       string     `json:"title"`
	Source      SourceType `json:"source"`
	Type        EventType `json:"type"`
	Payload     []byte     `json:"payload"`
	CollectedAt time.Time  `json:"collected_at"`
	OccuredAt   time.Time  `json:"occured_time"`
}
