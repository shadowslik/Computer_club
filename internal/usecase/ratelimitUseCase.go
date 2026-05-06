package usecase

import (
	"context"
	"fmt"
	"proxy/internal/config"
	"proxy/internal/domain"
	"time"
)

type RateLimitService struct {
	repo domain.RateLimitRepo
	cfg  *config.Config
}

func NewRateLimitService(repo domain.RateLimitRepo, cfg *config.Config) *RateLimitService {
	return &RateLimitService{repo: repo, cfg: cfg}
}

func resetIfNeeded(c *domain.Counter) {
	now := time.Now()
	if now.Sub(c.LastSecond) >= time.Second {
		c.RequestsSecond = 0
		c.NewConnsLastSec = 0
		c.LastSecond = now
	}
	if now.Sub(c.LastMinute) >= time.Minute {
		c.RequestsMinute = 0
		c.LastMinute = now
	}
	if now.Sub(c.LastHour) >= time.Hour {
		c.RequestsHour = 0
		c.LastHour = now
	}
	if now.Sub(c.LastDay) >= 24*time.Hour {
		c.RequestsDay = 0
		c.UploadDay = 0
		c.DownloadDay = 0
		c.LastDay = now
	}
}

func isEnabled(v int) bool     { return v >= 0 }
func isEnabled64(v int64) bool { return v >= 0 }

func (s *RateLimitService) defaultRule() domain.RateRule {
	rl := s.cfg.RateLimit
	return domain.RateRule{
		ID:              "default",
		Target:          domain.TargetIP,
		MaxRPS:          rl.DefaultRPS,
		MaxRPM:          rl.DefaultRPM,
		MaxRPH:          rl.DefaultRPH,
		MaxRPD:          rl.DefaultRPD,
		MaxUploadBps:    rl.DefaultMaxUploadBps,
		MaxDownloadBps:  rl.DefaultMaxDownloadBps,
		MaxTrafficDay:   rl.DefaultMaxTrafficDay,
		MaxConcurrent:   rl.DefaultMaxConcurrent,
		MaxNewPerSecond: rl.DefaultMaxNewPerSec,
	}
}

func (s *RateLimitService) GetRulesForIP(ctx context.Context, ip string) ([]domain.RateRule, error) {
	rules, err := s.repo.FindRules(ctx, ip)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return []domain.RateRule{s.defaultRule()}, nil
	}
	return rules, nil
}

func (s *RateLimitService) Check(ctx context.Context, ip string, uploadBytes int64, rules []domain.RateRule) (domain.LimitResult, error) {
	if len(rules) == 0 {
		rules = []domain.RateRule{s.defaultRule()}
	}

	counter, err := s.repo.GetCounter(ctx, ip)
	if err != nil {
		return domain.LimitResult{Allowed: true}, err
	}

	resetIfNeeded(counter)

	counter.RequestsSecond++
	counter.RequestsMinute++
	counter.RequestsHour++
	counter.RequestsDay++
	counter.UploadDay += uploadBytes

	for _, rule := range rules {
		if isEnabled(rule.MaxRPS) && counter.RequestsSecond > rule.MaxRPS {
			return domain.LimitResult{false, fmt.Sprintf("превышен RPS: %d req/s", rule.MaxRPS)}, nil
		}
		if isEnabled(rule.MaxRPM) && counter.RequestsMinute > rule.MaxRPM {
			return domain.LimitResult{false, fmt.Sprintf("превышен RPM: %d req/min", rule.MaxRPM)}, nil
		}
		if isEnabled(rule.MaxRPH) && counter.RequestsHour > rule.MaxRPH {
			return domain.LimitResult{false, fmt.Sprintf("превышен RPH: %d req/h", rule.MaxRPH)}, nil
		}
		if isEnabled(rule.MaxRPD) && counter.RequestsDay > rule.MaxRPD {
			return domain.LimitResult{false, fmt.Sprintf("превышен RPD: %d req/day", rule.MaxRPD)}, nil
		}
		if isEnabled64(rule.MaxTrafficDay) && (counter.UploadDay+counter.DownloadDay) > rule.MaxTrafficDay {
			return domain.LimitResult{false, "превышен дневной лимит трафика"}, nil
		}
		if isEnabled(rule.MaxConcurrent) && counter.ActiveConns > rule.MaxConcurrent {
			return domain.LimitResult{false, fmt.Sprintf("превышен лимит соединений: %d", rule.MaxConcurrent)}, nil
		}
	}

	if err := s.repo.SaveCounter(ctx, ip, counter); err != nil {
		return domain.LimitResult{Allowed: true}, err
	}
	return domain.LimitResult{Allowed: true}, nil
}

func (s *RateLimitService) TrackDownload(ctx context.Context, ip string, bytes int64) error {
	counter, err := s.repo.GetCounter(ctx, ip)
	if err != nil {
		return err
	}
	resetIfNeeded(counter)
	counter.DownloadDay += bytes
	return s.repo.SaveCounter(ctx, ip, counter)
}

func (s *RateLimitService) OnConnect(ctx context.Context, ip string) (domain.LimitResult, error) {
	rules, err := s.repo.FindRules(ctx, ip)
	if err != nil {
		return domain.LimitResult{Allowed: true}, err
	}
	if len(rules) == 0 {
		rules = []domain.RateRule{s.defaultRule()}
	}

	counter, err := s.repo.GetCounter(ctx, ip)
	if err != nil {
		return domain.LimitResult{Allowed: true}, err
	}

	resetIfNeeded(counter)
	counter.NewConnsLastSec++

	for _, rule := range rules {
		if isEnabled(rule.MaxNewPerSecond) && counter.NewConnsLastSec > rule.MaxNewPerSecond {
			return domain.LimitResult{
				Allowed: false,
				Reason:  fmt.Sprintf("превышен лимит новых соединений: %d/s", rule.MaxNewPerSecond),
			}, nil
		}
		if isEnabled(rule.MaxConcurrent) && counter.ActiveConns >= rule.MaxConcurrent {
			return domain.LimitResult{
				Allowed: false,
				Reason:  fmt.Sprintf("превышен лимит соединений: %d", rule.MaxConcurrent),
			}, nil
		}
	}

	if err := s.repo.IncrementConn(ctx, ip); err != nil {
		return domain.LimitResult{Allowed: true}, err
	}
	if err := s.repo.SaveCounter(ctx, ip, counter); err != nil {
		return domain.LimitResult{Allowed: true}, err
	}
	return domain.LimitResult{Allowed: true}, nil
}

func (s *RateLimitService) OnDisconnect(ctx context.Context, ip string) error {
	return s.repo.DecrementConn(ctx, ip)
}

func (s *RateLimitService) AddRule(ctx context.Context, rule domain.RateRule) error {
	return s.repo.AddRule(ctx, rule)
}

func (s *RateLimitService) RemoveRule(ctx context.Context, id string) error {
	return s.repo.RemoveRule(ctx, id)
}

func (s *RateLimitService) GetRules(ctx context.Context) ([]domain.RateRule, error) {
	return s.repo.GetRules(ctx)
}
