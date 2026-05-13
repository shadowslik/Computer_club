package domain

import (
	"context"
	"time"
)

type RuleTarget string

const (
	TargetIP     RuleTarget = "ip"
	TargetSubnet RuleTarget = "subnet"
)

type RateRule struct {
	ID     string     `json:"id"     example:"550e8400-e29b-41d4-a716-446655440000"`
	Target RuleTarget `json:"target" example:"ip"         enums:"ip,subnet"`
	Value  string     `json:"value"  example:"192.168.1.1"`

	MaxRPS int `json:"max_rps" example:"100"`
	MaxRPM int `json:"max_rpm" example:"1000"`
	MaxRPH int `json:"max_rph" example:"10000"`
	MaxRPD int `json:"max_rpd" example:"100000"`

	MaxUploadBps   int64 `json:"max_upload_bps"   example:"1048576"`
	MaxDownloadBps int64 `json:"max_download_bps" example:"1048576"`
	MaxTrafficDay  int64 `json:"max_traffic_day" example:"104857600"`

	MaxConcurrent   int `json:"max_concurrent"    example:"100"`
	MaxNewPerSecond int `json:"max_new_per_second" example:"20"`
}

type Counter struct {
	RequestsSecond int
	RequestsMinute int
	RequestsHour   int
	RequestsDay    int

	UploadDay   int64
	DownloadDay int64

	ActiveConns     int
	NewConnsLastSec int

	LastSecond time.Time
	LastMinute time.Time
	LastHour   time.Time
	LastDay    time.Time
}

type LimitResult struct {
	Allowed bool
	Reason  string
}

type RateLimitRepo interface {
	AddRule(ctx context.Context, rule RateRule) error
	RemoveRule(ctx context.Context, id string) error
	GetRules(ctx context.Context) ([]RateRule, error)
	FindRules(ctx context.Context, ip string) ([]RateRule, error)
	GetCounter(ctx context.Context, key string) (*Counter, error)
	SaveCounter(ctx context.Context, key string, c *Counter) error
	IncrementConn(ctx context.Context, key string) error
	DecrementConn(ctx context.Context, key string) error
}

type RateLimitUseCase interface {
	Check(ctx context.Context, ip string, uploadBytes int64, rules []RateRule) (LimitResult, error)
	TrackDownload(ctx context.Context, ip string, bytes int64) error
	OnConnect(ctx context.Context, ip string) (LimitResult, error)
	OnDisconnect(ctx context.Context, ip string) error
	AddRule(ctx context.Context, rule RateRule) error
	RemoveRule(ctx context.Context, id string) error
	GetRules(ctx context.Context) ([]RateRule, error)
	GetRulesForIP(ctx context.Context, ip string) ([]RateRule, error)
}
