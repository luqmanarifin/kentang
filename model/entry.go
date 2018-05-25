package model

import "time"

type Entry struct {
	ID        int       `json:"id"`
	Source    string    `json:"source"`
	Keyword   string    `json:"keyword"`
	Timestamp time.Time `json:"timestamp"`
}
