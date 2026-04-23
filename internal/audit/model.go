package audit

import "time"

type Event struct {
	TraceID   string    `json:"trace_id"`
	Subject   string    `json:"subject"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Decision  string    `json:"decision"`
	Rule      string    `json:"rule"`
	CreatedAt time.Time `json:"created_at"`
}
