package domain

import "time"

type DenialStat struct {
	IP         string    `json:"ip"          example:"1.2.3.4"`
	Denials    int64     `json:"denials"     example:"42"`
	LastDenied time.Time `json:"last_denied" example:"2026-05-06T13:00:00+04:00"`
	Reason     string    `json:"reason"      example:"blacklist"  enums:"blacklist,graylist,not_in_whitelist"`
}

type DenialsRepo interface {
	Record(ip, reason string)
	GetAll() []DenialStat
}
