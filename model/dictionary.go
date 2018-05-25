package model

import "time"

type Dictionary struct {
	ID          int       `json:"id"`
	Source      string    `json:"source"`
	Keyword     string    `json:"keyword"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}
