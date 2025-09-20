package types

import "time"

type QueueItem struct {
	Message     string    `json:"message"`
	Id          string    `json:"id"`
	AvailableAt time.Time `json:"available_at"`
}
