package domain

type ClientStat struct {
	IP           string  `json:"ip"            example:"1.2.3.4"`
	Requests     int64   `json:"requests"      example:"4200"`
	TrafficMB    float64 `json:"traffic_mb"    example:"12.5"`
	AvgLatencyMs float64 `json:"avg_latency_ms" example:"18.3"`
	ErrorRate    float64 `json:"error_rate"    example:"1.2"`
	LastSeen     string  `json:"last_seen"     example:"13:05:00"`
}

type ClientStatsRepo interface {
	Update(stat ClientStat)
	GetTop(limit int, sortBy string) []ClientStat
}
