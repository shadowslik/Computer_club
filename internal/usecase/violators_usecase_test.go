package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockViolatorsRepo struct {
	violators []domain.Violator
}

func (m *mockViolatorsRepo) Add(v domain.Violator) {
	m.violators = append(m.violators, v)
}

func (m *mockViolatorsRepo) GetAll() []domain.Violator {
	return m.violators
}

func (m *mockViolatorsRepo) Remove(ip string) {
	filtered := m.violators[:0]
	for _, v := range m.violators {
		if v.IP != ip {
			filtered = append(filtered, v)
		}
	}
	m.violators = filtered
}

func (m *mockViolatorsRepo) IsBanned(ip string) bool {
	for _, v := range m.violators {
		if v.IP == ip {
			return true
		}
	}
	return false
}

func TestViolatorsUseCase_GetAll_ReturnsRepoList(t *testing.T) {
	repo := &mockViolatorsRepo{
		violators: []domain.Violator{
			{IP: "1.1.1.1", Endpoint: "/api", Requests: 200, Limit: 100, BlockedAt: time.Now()},
			{IP: "2.2.2.2", Endpoint: "/api", Requests: 150, Limit: 100, BlockedAt: time.Now()},
		},
	}
	uc := NewViolatorsUseCase(repo)

	all := uc.GetAll()
	require.Len(t, all, 2)
	assert.Equal(t, "1.1.1.1", all[0].IP)
	assert.Equal(t, "2.2.2.2", all[1].IP)
}

func TestViolatorsUseCase_GetAll_Empty(t *testing.T) {
	repo := &mockViolatorsRepo{}
	uc := NewViolatorsUseCase(repo)

	all := uc.GetAll()
	assert.Empty(t, all)
}

func TestViolatorsUseCase_Unban_CallsRemove(t *testing.T) {
	repo := &mockViolatorsRepo{
		violators: []domain.Violator{
			{IP: "3.3.3.3"},
			{IP: "4.4.4.4"},
		},
	}
	uc := NewViolatorsUseCase(repo)

	uc.Unban("3.3.3.3")
	all := uc.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, "4.4.4.4", all[0].IP)
}

func TestViolatorsUseCase_Unban_NonExistingIP_NoError(t *testing.T) {
	repo := &mockViolatorsRepo{
		violators: []domain.Violator{{IP: "5.5.5.5"}},
	}
	uc := NewViolatorsUseCase(repo)

	uc.Unban("9.9.9.9")
	assert.Len(t, uc.GetAll(), 1)
}

func TestViolatorsUseCase_Add_DelegatesToRepo(t *testing.T) {
	repo := &mockViolatorsRepo{}
	uc := NewViolatorsUseCase(repo)

	v := domain.Violator{IP: "6.6.6.6", Endpoint: "/test", Requests: 300, Limit: 100}
	uc.Add(v)

	all := uc.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, "6.6.6.6", all[0].IP)
}

func TestViolatorsUseCase_IsBanned_DelegatesToRepo(t *testing.T) {
	repo := &mockViolatorsRepo{
		violators: []domain.Violator{{IP: "7.7.7.7"}},
	}
	uc := NewViolatorsUseCase(repo)

	assert.True(t, uc.IsBanned("7.7.7.7"))
	assert.False(t, uc.IsBanned("8.8.8.8"))
}
