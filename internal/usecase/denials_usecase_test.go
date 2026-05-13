package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

type mockDenialsRepo struct {
	records []domain.DenialStat
}

func (m *mockDenialsRepo) Record(ip, reason string) {
	for i, s := range m.records {
		if s.IP == ip {
			m.records[i].Denials++
			m.records[i].LastDenied = time.Now()
			m.records[i].Reason = reason
			return
		}
	}
	m.records = append(m.records, domain.DenialStat{
		IP:         ip,
		Denials:    1,
		LastDenied: time.Now(),
		Reason:     reason,
	})
}

func (m *mockDenialsRepo) GetAll() []domain.DenialStat {
	return m.records
}

func TestDenialsUseCase_GetAll_ReturnsRepoData(t *testing.T) {
	repo := &mockDenialsRepo{}
	repo.records = []domain.DenialStat{
		{IP: "1.1.1.1", Denials: 3, Reason: "blacklist"},
		{IP: "2.2.2.2", Denials: 1, Reason: "graylist"},
	}
	uc := NewDenialsUseCase(repo)

	all := uc.GetAll()
	require.Len(t, all, 2)
	assert.Equal(t, "1.1.1.1", all[0].IP)
	assert.Equal(t, int64(3), all[0].Denials)
	assert.Equal(t, "2.2.2.2", all[1].IP)
}

func TestDenialsUseCase_GetAll_EmptyRepo(t *testing.T) {
	repo := &mockDenialsRepo{}
	uc := NewDenialsUseCase(repo)

	all := uc.GetAll()
	assert.Empty(t, all)
}

func TestDenialsUseCase_Record_IncrementsDenials(t *testing.T) {
	repo := &mockDenialsRepo{}
	uc := NewDenialsUseCase(repo)

	uc.Record("3.3.3.3", "blacklist")
	uc.Record("3.3.3.3", "blacklist")

	all := uc.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, int64(2), all[0].Denials)
}

func TestDenialsUseCase_Record_MultipleIPsStackCorrectly(t *testing.T) {
	repo := &mockDenialsRepo{}
	uc := NewDenialsUseCase(repo)

	uc.Record("10.0.0.1", "blacklist")
	uc.Record("10.0.0.2", "graylist")
	uc.Record("10.0.0.1", "blacklist")

	all := uc.GetAll()
	require.Len(t, all, 2)

	ipMap := make(map[string]domain.DenialStat)
	for _, s := range all {
		ipMap[s.IP] = s
	}

	assert.Equal(t, int64(2), ipMap["10.0.0.1"].Denials)
	assert.Equal(t, int64(1), ipMap["10.0.0.2"].Denials)
}
