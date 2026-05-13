package usecase

import (
	"net"
	"net/http"
	"proxy/internal/domain"
	"strings"
	"sync"
	"time"
)

type UpstreamService struct {
	mu         sync.RWMutex
	upstreams  []domain.UpstreamStatus
	httpClient *http.Client
}

func NewUpstreamService(targets []domain.UpstreamTarget) *UpstreamService {
	svc := &UpstreamService{
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
	svc.initUpstreams(targets)
	go svc.startHealthChecker()
	return svc
}

func (s *UpstreamService) initUpstreams(targets []domain.UpstreamTarget) {
	s.upstreams = make([]domain.UpstreamStatus, len(targets))
	for i, t := range targets {
		s.upstreams[i] = domain.UpstreamStatus{
			ID:             t.ID,
			Name:           t.Name,
			URL:            t.URL,
			LatencyHistory: make([]int, 20),
		}
	}
}

func (s *UpstreamService) startHealthChecker() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		s.checkAllUpstreams()
	}
}

func (s *UpstreamService) checkAllUpstreams() {
	var wg sync.WaitGroup
	for idx := range s.upstreams {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.checkUpstream(i)
		}(idx)
	}
	wg.Wait()
}

func (s *UpstreamService) checkUpstream(idx int) {
	s.mu.RLock()
	up := s.upstreams[idx]
	s.mu.RUnlock()

	start := time.Now()
	var checkErr error
	var statusCode int

	addr := strings.TrimPrefix(up.URL, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	if idx := strings.Index(addr, "/"); idx >= 0 {
		addr = addr[:idx]
	}

	if !strings.HasPrefix(up.URL, "http") {
		conn, err := net.DialTimeout("tcp", up.URL, 2*time.Second)
		checkErr = err
		if err == nil {
			conn.Close()
			statusCode = 200
		} else {
			statusCode = 500
		}
	} else {
		resp, err := s.httpClient.Get(up.URL)
		checkErr = err
		if err == nil {
			resp.Body.Close()
			statusCode = resp.StatusCode
		}
	}

	latencyMs := int(time.Since(start).Milliseconds())

	s.mu.Lock()
	defer s.mu.Unlock()

	ptr := &s.upstreams[idx]
	ptr.LastCheck = time.Now().Format("15:04:05")
	ptr.Latency = latencyMs
	if len(ptr.LatencyHistory) == 0 {
		ptr.LatencyHistory = make([]int, 20)
	}
	ptr.LatencyHistory = append(ptr.LatencyHistory[1:], latencyMs)
	ptr.RequestsTotal++

	if checkErr != nil {
		ptr.Status = "down"
		ptr.ErrorRate = 100.0
		return
	}
	switch {
	case statusCode >= 500:
		ptr.Status = "degraded"
		ptr.ErrorRate = (ptr.ErrorRate*float64(ptr.RequestsTotal-1) + 100) / float64(ptr.RequestsTotal)
	case statusCode >= 400:
		ptr.Status = "degraded"
		ptr.ErrorRate = (ptr.ErrorRate*float64(ptr.RequestsTotal-1) + 50) / float64(ptr.RequestsTotal)
	default:
		ptr.Status = "healthy"
		ptr.ErrorRate = (ptr.ErrorRate * float64(ptr.RequestsTotal-1)) / float64(ptr.RequestsTotal)
	}
}

func (s *UpstreamService) GetStatuses() []domain.UpstreamStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.UpstreamStatus, len(s.upstreams))
	copy(result, s.upstreams)
	return result
}
