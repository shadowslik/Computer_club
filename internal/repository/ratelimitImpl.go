package repository

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"proxy/internal/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RateLimitStore struct {
	mu       sync.Mutex
	rules    []domain.RateRule
	counters map[string]*domain.Counter // ключ: IP-строка
	fileName string
}

func NewRateLimitStore(file string) (*RateLimitStore, error) {
	s := &RateLimitStore{
		counters: make(map[string]*domain.Counter),
		fileName: file,
	}
	if err := s.loadRules(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *RateLimitStore) loadRules() error {
	data, err := os.ReadFile(s.fileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.rules)
}

func (s *RateLimitStore) saveRules() error {
	data, err := json.MarshalIndent(s.rules, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.fileName + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.fileName)
}

func (s *RateLimitStore) AddRule(_ context.Context, rule domain.RateRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rule.ID = uuid.NewString()
	s.rules = append(s.rules, rule)
	return s.saveRules()
}

func (s *RateLimitStore) RemoveRule(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := s.rules[:0]
	for _, r := range s.rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	s.rules = filtered
	return s.saveRules()
}

func (s *RateLimitStore) GetRules(_ context.Context) ([]domain.RateRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]domain.RateRule, len(s.rules))
	copy(result, s.rules)
	return result, nil
}

func (s *RateLimitStore) FindRules(_ context.Context, ipStr string) ([]domain.RateRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ip := net.ParseIP(ipStr)
	var matched []domain.RateRule

	for _, rule := range s.rules {
		switch rule.Target {
		case domain.TargetIP:
			if rule.Value == ipStr {
				matched = append(matched, rule)
			}
		case domain.TargetSubnet:
			_, network, err := net.ParseCIDR(rule.Value)
			if err == nil && network.Contains(ip) {
				matched = append(matched, rule)
			}
		}
	}
	return matched, nil
}

func (s *RateLimitStore) getOrCreate(key string) *domain.Counter {
	c, ok := s.counters[key]
	if !ok {
		now := time.Now()
		c = &domain.Counter{
			LastSecond: now,
			LastMinute: now,
			LastHour:   now,
			LastDay:    now,
		}
		s.counters[key] = c
	}
	return c
}

func (s *RateLimitStore) GetCounter(_ context.Context, key string) (*domain.Counter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.getOrCreate(key)
	copy := *c
	return &copy, nil
}

func (s *RateLimitStore) SaveCounter(_ context.Context, key string, c *domain.Counter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *c
	s.counters[key] = &copy
	return nil
}

func (s *RateLimitStore) IncrementConn(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.getOrCreate(key).ActiveConns++
	return nil
}

func (s *RateLimitStore) DecrementConn(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.getOrCreate(key)
	if c.ActiveConns > 0 {
		c.ActiveConns--
	}
	return nil
}
