package memory

import "time"

type Memory struct {
	ID           string     `json:"id"`
	ContainerTag string     `json:"containerTag"`
	Content      string     `json:"content"`
	Source       string     `json:"source,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
}

type SaveRequest struct {
	ContainerTag string
	Content      string
	Source       string
	ExpiresAt    *time.Time
}

type RecallRequest struct {
	ContainerTag string
	Query        string
	Limit        int
}

type RecallResult struct {
	Memory Memory `json:"memory"`
	Score  int64  `json:"score"`
	Reason string `json:"reason"`
}
