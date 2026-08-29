package model

import "time"

type RocketMsg struct {
	Body           string    `json:"body"`
	MsgId          string    `json:"msgId"`
	Tags           string    `json:"tags"`
	Keys           string    `json:"keys"`
	StoreTimestamp time.Time `json:"storeTimestamp"`
}
