package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockIpRepo struct {
	data         map[string]bool
	getAllErr    error
	addErr       error
	removeErr    error
	isAllowedErr error
}

func newMockIpRepo() *mockIpRepo {
	return &mockIpRepo{data: make(map[string]bool)}
}

func (m *mockIpRepo) GetAll() (map[string]bool, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.data, nil
}

func (m *mockIpRepo) Add(ctx context.Context, ip string) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.data[ip] = true
	return nil
}

func (m *mockIpRepo) Remove(ctx context.Context, ip string) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	delete(m.data, ip)
	return nil
}

func (m *mockIpRepo) IsAllowed(ctx context.Context, ip string) (bool, error) {
	if m.isAllowedErr != nil {
		return false, m.isAllowedErr
	}
	return m.data[ip], nil
}

func TestHTTPUseCase_Add_ValidNewIP(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	added, err := uc.Add(context.Background(), "192.168.1.1")
	require.NoError(t, err)
	assert.True(t, added)
	assert.True(t, repo.data["192.168.1.1"])
}

func TestHTTPUseCase_Add_DuplicateIP(t *testing.T) {
	repo := newMockIpRepo()
	repo.data["10.0.0.1"] = true
	uc := NewHTTPUseCase(repo)

	added, err := uc.Add(context.Background(), "10.0.0.1")
	require.NoError(t, err)
	assert.False(t, added)
}

func TestHTTPUseCase_Add_InvalidIP(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	added, err := uc.Add(context.Background(), "not-an-ip")
	assert.Error(t, err)
	assert.False(t, added)
}

func TestHTTPUseCase_Add_CIDR(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	added, err := uc.Add(context.Background(), "10.0.0.0/8")
	require.NoError(t, err)
	assert.True(t, added)
}

func TestHTTPUseCase_Add_RepoError(t *testing.T) {
	repo := newMockIpRepo()
	repo.addErr = errors.New("storage failure")
	uc := NewHTTPUseCase(repo)

	added, err := uc.Add(context.Background(), "1.2.3.4")
	assert.Error(t, err)
	assert.False(t, added)
}

func TestHTTPUseCase_Remove_ExistingIP(t *testing.T) {
	repo := newMockIpRepo()
	repo.data["5.5.5.5"] = true
	uc := NewHTTPUseCase(repo)

	removed, err := uc.Remove(context.Background(), "5.5.5.5")
	require.NoError(t, err)
	assert.True(t, removed)
	assert.False(t, repo.data["5.5.5.5"])
}

func TestHTTPUseCase_Remove_NonExistingIP(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	removed, err := uc.Remove(context.Background(), "9.9.9.9")
	require.NoError(t, err)
	assert.False(t, removed)
}

func TestHTTPUseCase_Remove_InvalidIP(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	removed, err := uc.Remove(context.Background(), "bad##ip")
	assert.Error(t, err)
	assert.False(t, removed)
}

func TestHTTPUseCase_Remove_RepoError(t *testing.T) {
	repo := newMockIpRepo()
	repo.data["7.7.7.7"] = true
	repo.removeErr = errors.New("remove failure")
	uc := NewHTTPUseCase(repo)

	removed, err := uc.Remove(context.Background(), "7.7.7.7")
	assert.Error(t, err)
	assert.False(t, removed)
}

func TestHTTPUseCase_IsAllowed_Valid(t *testing.T) {
	repo := newMockIpRepo()
	repo.data["8.8.8.8"] = true
	uc := NewHTTPUseCase(repo)

	ok, err := uc.IsAllowed(context.Background(), "8.8.8.8")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestHTTPUseCase_IsAllowed_InvalidIP(t *testing.T) {
	repo := newMockIpRepo()
	uc := NewHTTPUseCase(repo)

	ok, err := uc.IsAllowed(context.Background(), "not-valid")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestHTTPUseCase_GetAll(t *testing.T) {
	repo := newMockIpRepo()
	repo.data["1.1.1.1"] = true
	repo.data["2.2.2.2"] = true
	uc := NewHTTPUseCase(repo)

	all, err := uc.GetAll()
	require.NoError(t, err)
	assert.Len(t, all, 2)
	assert.True(t, all["1.1.1.1"])
	assert.True(t, all["2.2.2.2"])
}

func TestHTTPUseCase_GetAll_RepoError(t *testing.T) {
	repo := newMockIpRepo()
	repo.getAllErr = errors.New("db error")
	uc := NewHTTPUseCase(repo)

	all, err := uc.GetAll()
	assert.Error(t, err)
	assert.Nil(t, all)
}
