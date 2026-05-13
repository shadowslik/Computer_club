package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/config"
	"proxy/internal/domain"
)

type mockRateLimitRepo struct {
	rules    []domain.RateRule
	counters map[string]*domain.Counter
	connKeys []string
	decrKeys []string

	addRuleErr     error
	removeRuleErr  error
	getRulesErr    error
	findRulesErr   error
	getCounterErr  error
	saveCounterErr error
	incrConnErr    error
	decrConnErr    error
}

func newMockRLRepo() *mockRateLimitRepo {
	return &mockRateLimitRepo{
		counters: make(map[string]*domain.Counter),
	}
}

func (m *mockRateLimitRepo) AddRule(_ context.Context, rule domain.RateRule) error {
	if m.addRuleErr != nil {
		return m.addRuleErr
	}
	m.rules = append(m.rules, rule)
	return nil
}

func (m *mockRateLimitRepo) RemoveRule(_ context.Context, id string) error {
	if m.removeRuleErr != nil {
		return m.removeRuleErr
	}
	filtered := m.rules[:0]
	for _, r := range m.rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	m.rules = filtered
	return nil
}

func (m *mockRateLimitRepo) GetRules(_ context.Context) ([]domain.RateRule, error) {
	return m.rules, m.getRulesErr
}

func (m *mockRateLimitRepo) FindRules(_ context.Context, _ string) ([]domain.RateRule, error) {
	return nil, m.findRulesErr
}

func (m *mockRateLimitRepo) GetCounter(_ context.Context, key string) (*domain.Counter, error) {
	if m.getCounterErr != nil {
		return nil, m.getCounterErr
	}
	if c, ok := m.counters[key]; ok {
		return c, nil
	}
	c := &domain.Counter{
		LastSecond: time.Now(),
		LastMinute: time.Now(),
		LastHour:   time.Now(),
		LastDay:    time.Now(),
	}
	m.counters[key] = c
	return c, nil
}

func (m *mockRateLimitRepo) SaveCounter(_ context.Context, key string, c *domain.Counter) error {
	if m.saveCounterErr != nil {
		return m.saveCounterErr
	}
	m.counters[key] = c
	return nil
}

func (m *mockRateLimitRepo) IncrementConn(_ context.Context, key string) error {
	if m.incrConnErr != nil {
		return m.incrConnErr
	}
	m.connKeys = append(m.connKeys, key)
	if c, ok := m.counters[key]; ok {
		c.ActiveConns++
	}
	return nil
}

func (m *mockRateLimitRepo) DecrementConn(_ context.Context, key string) error {
	if m.decrConnErr != nil {
		return m.decrConnErr
	}
	m.decrKeys = append(m.decrKeys, key)
	return nil
}

func minimalConfig() *config.Config {
	return &config.Config{
		RateLimit: config.RateLimitConfig{
			DefaultRPS:            10,
			DefaultRPM:            100,
			DefaultRPH:            1000,
			DefaultRPD:            10000,
			DefaultMaxUploadBps:   -1,
			DefaultMaxDownloadBps: -1,
			DefaultMaxTrafficDay:  -1,
			DefaultMaxConcurrent:  50,
			DefaultMaxNewPerSec:   20,
		},
	}
}

func TestRateLimitService_Check_WithinLimits(t *testing.T) {
	repo := newMockRLRepo()
	svc := NewRateLimitService(repo, minimalConfig())

	result, err := svc.Check(context.Background(), "1.1.1.1", 0, nil)
	require.NoError(t, err)
	assert.True(t, result.Allowed)
}

func TestRateLimitService_Check_ExceedsMaxRPS(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultRPS = 2
	svc := NewRateLimitService(repo, cfg)

	ctx := context.Background()
	ip := "2.2.2.2"

	counter := &domain.Counter{
		RequestsSecond: 2,
		LastSecond:     time.Now(),
		LastMinute:     time.Now(),
		LastHour:       time.Now(),
		LastDay:        time.Now(),
	}
	repo.counters[ip] = counter

	result, err := svc.Check(ctx, ip, 0, nil)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Contains(t, result.Reason, "RPS")
}

func TestRateLimitService_Check_ExceedsMaxRPM(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultRPS = 1000
	cfg.RateLimit.DefaultRPM = 3
	svc := NewRateLimitService(repo, cfg)

	ip := "3.3.3.3"
	counter := &domain.Counter{
		RequestsSecond: 0,
		RequestsMinute: 3,
		LastSecond:     time.Now(),
		LastMinute:     time.Now(),
		LastHour:       time.Now(),
		LastDay:        time.Now(),
	}
	repo.counters[ip] = counter

	result, err := svc.Check(context.Background(), ip, 0, nil)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Contains(t, result.Reason, "RPM")
}

func TestRateLimitService_Check_ExceedsMaxRPD(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultRPS = 1000
	cfg.RateLimit.DefaultRPM = 100000
	cfg.RateLimit.DefaultRPH = 1000000
	cfg.RateLimit.DefaultRPD = 5
	svc := NewRateLimitService(repo, cfg)

	ip := "4.4.4.4"
	counter := &domain.Counter{
		RequestsDay: 5,
		LastSecond:  time.Now(),
		LastMinute:  time.Now(),
		LastHour:    time.Now(),
		LastDay:     time.Now(),
	}
	repo.counters[ip] = counter

	result, err := svc.Check(context.Background(), ip, 0, nil)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Contains(t, result.Reason, "RPD")
}

func TestRateLimitService_Check_ExceedsMaxTrafficDay(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultRPS = 10000
	cfg.RateLimit.DefaultRPM = 100000
	cfg.RateLimit.DefaultRPH = 1000000
	cfg.RateLimit.DefaultRPD = 10000000
	cfg.RateLimit.DefaultMaxTrafficDay = 100
	svc := NewRateLimitService(repo, cfg)

	ip := "5.5.5.5"
	counter := &domain.Counter{
		DownloadDay: 90,
		LastSecond:  time.Now(),
		LastMinute:  time.Now(),
		LastHour:    time.Now(),
		LastDay:     time.Now(),
	}
	repo.counters[ip] = counter

	result, err := svc.Check(context.Background(), ip, 20, nil)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
}

func TestRateLimitService_TrackDownload_IncrementsDownloadDay(t *testing.T) {
	repo := newMockRLRepo()
	svc := NewRateLimitService(repo, minimalConfig())
	ip := "6.6.6.6"

	err := svc.TrackDownload(context.Background(), ip, 512)
	require.NoError(t, err)

	assert.Equal(t, int64(512), repo.counters[ip].DownloadDay)
}

func TestRateLimitService_OnConnect_WithinMaxConcurrent(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultMaxConcurrent = 5
	svc := NewRateLimitService(repo, cfg)

	result, err := svc.OnConnect(context.Background(), "7.7.7.7")
	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Contains(t, repo.connKeys, "7.7.7.7")
}

func TestRateLimitService_OnConnect_ExceedsMaxConcurrent(t *testing.T) {
	repo := newMockRLRepo()
	cfg := minimalConfig()
	cfg.RateLimit.DefaultMaxConcurrent = 2
	cfg.RateLimit.DefaultMaxNewPerSec = 100
	svc := NewRateLimitService(repo, cfg)

	ip := "8.8.8.8"
	counter := &domain.Counter{
		ActiveConns: 2,
		LastSecond:  time.Now(),
		LastMinute:  time.Now(),
		LastHour:    time.Now(),
		LastDay:     time.Now(),
	}
	repo.counters[ip] = counter

	result, err := svc.OnConnect(context.Background(), ip)
	require.NoError(t, err)
	assert.False(t, result.Allowed)
}

func TestRateLimitService_OnDisconnect_CallsDecrement(t *testing.T) {
	repo := newMockRLRepo()
	svc := NewRateLimitService(repo, minimalConfig())

	err := svc.OnDisconnect(context.Background(), "9.9.9.9")
	require.NoError(t, err)
	assert.Contains(t, repo.decrKeys, "9.9.9.9")
}

func TestRateLimitService_AddRule_DelegatesToRepo(t *testing.T) {
	repo := newMockRLRepo()
	svc := NewRateLimitService(repo, minimalConfig())

	rule := domain.RateRule{ID: "rule1", MaxRPS: 50}
	err := svc.AddRule(context.Background(), rule)
	require.NoError(t, err)
	assert.Len(t, repo.rules, 1)
	assert.Equal(t, "rule1", repo.rules[0].ID)
}

func TestRateLimitService_RemoveRule_DelegatesToRepo(t *testing.T) {
	repo := newMockRLRepo()
	repo.rules = []domain.RateRule{{ID: "r1"}, {ID: "r2"}}
	svc := NewRateLimitService(repo, minimalConfig())

	err := svc.RemoveRule(context.Background(), "r1")
	require.NoError(t, err)
	assert.Len(t, repo.rules, 1)
	assert.Equal(t, "r2", repo.rules[0].ID)
}

func TestRateLimitService_GetRules_DelegatesToRepo(t *testing.T) {
	repo := newMockRLRepo()
	repo.rules = []domain.RateRule{{ID: "x"}, {ID: "y"}}
	svc := NewRateLimitService(repo, minimalConfig())

	rules, err := svc.GetRules(context.Background())
	require.NoError(t, err)
	assert.Len(t, rules, 2)
}

func TestRateLimitService_GetRulesForIP_NoMatch_ReturnsDefault(t *testing.T) {
	repo := newMockRLRepo()
	svc := NewRateLimitService(repo, minimalConfig())

	rules, err := svc.GetRulesForIP(context.Background(), "1.2.3.4")
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "default", rules[0].ID)
}

func TestRateLimitService_GetRulesForIP_WithMatch(t *testing.T) {
	repo := newMockRLRepo()
	repo.rules = []domain.RateRule{{ID: "specific", MaxRPS: 5}}
	svc := NewRateLimitService(repo, minimalConfig())

	rules, err := svc.GetRulesForIP(context.Background(), "9.9.9.9")
	require.NoError(t, err)
	assert.NotEmpty(t, rules)
}
