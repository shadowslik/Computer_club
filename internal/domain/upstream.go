package domain

type UpstreamTarget struct {
	ID   string
	Name string
	URL  string
	Type string
}

type UpstreamStatus struct {
	ID             string  `json:"id"             example:"java-backend"`
	Name           string  `json:"name"           example:"Java Backend"`
	URL            string  `json:"url"            example:"http://localhost:8080/health"`
	Status         string  `json:"status"         example:"healthy"  enums:"healthy,degraded,down"`
	Latency        int     `json:"latency"        example:"12"`
	LatencyHistory []int   `json:"latencyHistory"`
	ErrorRate      float64 `json:"errorRate"      example:"0.5"`
	RequestsTotal  int64   `json:"requestsTotal"  example:"1234"`
	LastCheck      string  `json:"lastCheck"      example:"13:05:00"`
}
