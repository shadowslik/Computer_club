package repository

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTmpIPRepo(t *testing.T, initialIPs []string) (*IpRepoImpl, string) {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "ips.json")

	if initialIPs != nil {
		data, err := json.MarshalIndent(ipFile{IPs: initialIPs}, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(file, data, 0644))
	}

	repo, err := NewIpRepoImpl(file, 100, 30*time.Second, zap.NewNop())
	require.NoError(t, err)
	return repo, file
}

func TestIpRepoImpl_GetAll_Empty(t *testing.T) {
	repo, _ := newTmpIPRepo(t, nil)
	all, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, all)
}

func TestIpRepoImpl_GetAll_WithEntries(t *testing.T) {
	repo, _ := newTmpIPRepo(t, []string{"1.2.3.4", "5.6.7.8"})
	all, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, all, 2)
	assert.True(t, all["1.2.3.4"])
	assert.True(t, all["5.6.7.8"])
}

func TestIpRepoImpl_IsAllowed_Present(t *testing.T) {
	repo, _ := newTmpIPRepo(t, []string{"10.0.0.1"})
	ok, err := repo.IsAllowed(context.Background(), "10.0.0.1")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestIpRepoImpl_IsAllowed_Absent(t *testing.T) {
	repo, _ := newTmpIPRepo(t, nil)
	ok, err := repo.IsAllowed(context.Background(), "10.0.0.2")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestIpRepoImpl_Add_PersistsAndMatches(t *testing.T) {
	repo, _ := newTmpIPRepo(t, nil)

	err := repo.Add(context.Background(), "192.168.1.1")
	require.NoError(t, err)

	ok, err := repo.IsAllowed(context.Background(), "192.168.1.1")
	require.NoError(t, err)
	assert.True(t, ok)

	all, err := repo.GetAll()
	require.NoError(t, err)
	assert.True(t, all["192.168.1.1"])
}

func TestIpRepoImpl_Remove_RemovesAndPersists(t *testing.T) {
	repo, _ := newTmpIPRepo(t, []string{"172.16.0.1", "172.16.0.2"})

	err := repo.Remove(context.Background(), "172.16.0.1")
	require.NoError(t, err)

	ok, err := repo.IsAllowed(context.Background(), "172.16.0.1")
	require.NoError(t, err)
	assert.False(t, ok)

	all, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, all, 1)
	assert.True(t, all["172.16.0.2"])
}

func TestIpRepoImpl_Add_CIDR(t *testing.T) {
	repo, _ := newTmpIPRepo(t, nil)

	err := repo.Add(context.Background(), "10.0.0.0/8")
	require.NoError(t, err)

	ok, err := repo.IsAllowed(context.Background(), "10.1.2.3")
	require.NoError(t, err)
	assert.True(t, ok)
}
