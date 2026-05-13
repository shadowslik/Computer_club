package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockClientStatsRepo struct {
	stats []domain.ClientStat
}

func (m *mockClientStatsRepo) Update(stat domain.ClientStat) {
	m.stats = append(m.stats, stat)
}

func (m *mockClientStatsRepo) GetTop(limit int, sortBy string) []domain.ClientStat {
	if limit > 0 && limit < len(m.stats) {
		return m.stats[:limit]
	}
	return m.stats
}

func TestClientStatsUseCase_Update_DelegatesToRepo(t *testing.T) {
	repo := &mockClientStatsRepo{}
	uc := NewClientStatsUseCase(repo)

	uc.Update(domain.ClientStat{IP: "1.1.1.1", Requests: 10})
	require.Len(t, repo.stats, 1)
	assert.Equal(t, "1.1.1.1", repo.stats[0].IP)
}

func TestClientStatsUseCase_GetTop_DelegatesToRepo(t *testing.T) {
	repo := &mockClientStatsRepo{
		stats: []domain.ClientStat{
			{IP: "a", Requests: 100},
			{IP: "b", Requests: 50},
		},
	}
	uc := NewClientStatsUseCase(repo)

	top := uc.GetTop(10, "requests")
	assert.Len(t, top, 2)
}
