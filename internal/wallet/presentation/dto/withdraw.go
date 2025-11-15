package dto

import "time"

type Transaction struct {
	Amount      int64
	Idempotency string
	ReleaseTime *time.Time `json:"release_time,omitempty"`
}
