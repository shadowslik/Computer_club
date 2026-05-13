package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

func tmpRateLimitStore(t *testing.T) (*RateLimitStore, string) {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "rules.json")
	store, err := NewRateLimitStore(file)
	require.NoError(t, err)
	return store, file
}

func TestRateLimitStore_NewStore_NonExistentFile(t *testing.T) {
	store, _ := tmpRateLimitStore(t)
	rules, err := store.GetRules(context.Background())
	require.NoError(t, err)
	assert.Empty(t, rules)
}

func TestRateLimitStore_AddAndGetRules(t *testing.T) {
	store, _ := tmpRateLimitStore(t)

	rule := domain.RateRule{Target: domain.TargetIP, Value: "1.2.3.4", MaxRPS: 10}
	err := store.AddRule(context.Background(), rule)
	require.NoError(t, err)

	rules, err := store.GetRules(context.Background())
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "1.2.3.4", rules[0].Value)
	assert.NotEmpty(t, rules[0].ID)
}

func TestRateLimitStore_RemoveRule(t *testing.T) {
	store, _ := tmpRateLimitStore(t)

	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "2.2.2.2"}))
	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "3.3.3.3"}))

	rules, _ := store.GetRules(context.Background())
	require.Len(t, rules, 2)
	idToRemove := rules[0].ID

	require.NoError(t, store.RemoveRule(context.Background(), idToRemove))

	rules, _ = store.GetRules(context.Background())
	assert.Len(t, rules, 1)
	assert.NotEqual(t, idToRemove, rules[0].ID)
}

func TestRateLimitStore_FindRules_IPMatch(t *testing.T) {
	store, _ := tmpRateLimitStore(t)

	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "4.4.4.4", MaxRPS: 5}))
	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "5.5.5.5", MaxRPS: 20}))

	matched, err := store.FindRules(context.Background(), "4.4.4.4")
	require.NoError(t, err)
	require.Len(t, matched, 1)
	assert.Equal(t, "4.4.4.4", matched[0].Value)
}

func TestRateLimitStore_FindRules_SubnetMatch(t *testing.T) {
	store, _ := tmpRateLimitStore(t)

	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetSubnet, Value: "10.0.0.0/8", MaxRPS: 3}))

	matched, err := store.FindRules(context.Background(), "10.1.2.3")
	require.NoError(t, err)
	require.Len(t, matched, 1)
	assert.Equal(t, "10.0.0.0/8", matched[0].Value)
}

func TestRateLimitStore_FindRules_NoMatch(t *testing.T) {
	store, _ := tmpRateLimitStore(t)

	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "6.6.6.6"}))

	matched, err := store.FindRules(context.Background(), "7.7.7.7")
	require.NoError(t, err)
	assert.Empty(t, matched)
}

func TestRateLimitStore_GetAndSaveCounter(t *testing.T) {
	store, _ := tmpRateLimitStore(t)
	ctx := context.Background()

	c, err := store.GetCounter(ctx, "1.1.1.1")
	require.NoError(t, err)
	c.RequestsSecond = 42

	require.NoError(t, store.SaveCounter(ctx, "1.1.1.1", c))

	c2, err := store.GetCounter(ctx, "1.1.1.1")
	require.NoError(t, err)
	assert.Equal(t, 42, c2.RequestsSecond)
}

func TestRateLimitStore_IncrementAndDecrementConn(t *testing.T) {
	store, _ := tmpRateLimitStore(t)
	ctx := context.Background()
	ip := "8.8.8.8"

	require.NoError(t, store.IncrementConn(ctx, ip))
	require.NoError(t, store.IncrementConn(ctx, ip))

	c, _ := store.GetCounter(ctx, ip)
	assert.Equal(t, 2, c.ActiveConns)

	require.NoError(t, store.DecrementConn(ctx, ip))
	c, _ = store.GetCounter(ctx, ip)
	assert.Equal(t, 1, c.ActiveConns)
}

func TestRateLimitStore_DecrementConn_NeverNegative(t *testing.T) {
	store, _ := tmpRateLimitStore(t)
	ctx := context.Background()

	require.NoError(t, store.DecrementConn(ctx, "9.9.9.9"))
	c, _ := store.GetCounter(ctx, "9.9.9.9")
	assert.Equal(t, 0, c.ActiveConns)
}

func TestRateLimitStore_PersistsRulesToDisk(t *testing.T) {
	store, file := tmpRateLimitStore(t)

	require.NoError(t, store.AddRule(context.Background(), domain.RateRule{Target: domain.TargetIP, Value: "11.11.11.11"}))

	store2, err := NewRateLimitStore(file)
	require.NoError(t, err)

	rules, err := store2.GetRules(context.Background())
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "11.11.11.11", rules[0].Value)
}

func TestRateLimitStore_Load_InvalidJSON_Error(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "rules.json")
	require.NoError(t, os.WriteFile(file, []byte("{invalid json"), 0644))

	_, err := NewRateLimitStore(file)
	assert.Error(t, err)
}
