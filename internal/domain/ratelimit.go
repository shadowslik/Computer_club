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
	ID     string     `json:"id"`
	Target RuleTarget `json:"target"`
	Value  string     `json:"value"`

	// 2.2.1 — запросы
	MaxRPS int `json:"max_rps"` // -1 = отключён, 0 = запрещено всё
	MaxRPM int `json:"max_rpm"`
	MaxRPH int `json:"max_rph"`
	MaxRPD int `json:"max_rpd"`

	// 2.2.2 — пропускная способность (байт)
	MaxUploadBps   int64 `json:"max_upload_bps"`   // байт за один запрос
	MaxDownloadBps int64 `json:"max_download_bps"` // байт за один запрос
	MaxTrafficDay  int64 `json:"max_traffic_day"`  // суммарно за день

	// 2.2.3 — соединения
	MaxConcurrent   int `json:"max_concurrent"`     // одновременных
	MaxNewPerSecond int `json:"max_new_per_second"` // новых в секунду
}

type Counter struct {
	// запросы
	RequestsSecond int
	RequestsMinute int
	RequestsHour   int
	RequestsDay    int

	// трафик за день
	UploadDay   int64
	DownloadDay int64

	// соединения
	ActiveConns     int
	NewConnsLastSec int

	// временные метки сброса
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
	Check(ctx context.Context, ip string, uploadBytes int64) (LimitResult, error)
	TrackDownload(ctx context.Context, ip string, bytes int64) error
	OnConnect(ctx context.Context, ip string) (LimitResult, error)
	OnDisconnect(ctx context.Context, ip string) error
	AddRule(ctx context.Context, rule RateRule) error
	RemoveRule(ctx context.Context, id string) error
	GetRules(ctx context.Context) ([]RateRule, error)
	GetRulesForIP(ctx context.Context, ip string) ([]RateRule, error) // ← новый
}
