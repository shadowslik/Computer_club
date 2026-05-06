package domain

import "time"

// Violator represents a client that exceeded rate limits and is currently banned.
type Violator struct {
	IP          string    `json:"ip"           example:"1.2.3.4"`
	Endpoint    string    `json:"endpoint"     example:"/api/resource"`
	Requests    int       `json:"requests"     example:"150"`
	Limit       int       `json:"limit"        example:"100"`
	BlockedAt   time.Time `json:"blocked_at"   example:"2026-05-06T13:00:00+04:00"`
	BannedUntil time.Time `json:"banned_until" example:"2026-05-06T13:05:00+04:00"`
}

type ViolatorsRepo interface {
	Add(v Violator)
	GetAll() []Violator
	Remove(ip string)
}
