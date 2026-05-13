package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy/internal/domain"
)

func TestViolatorsRepo_AddAndGetAll(t *testing.T) {
	r := NewViolatorsRepo()
	v := domain.Violator{
		IP:          "10.0.0.1",
		Endpoint:    "/api",
		Requests:    200,
		Limit:       100,
		BlockedAt:   time.Now(),
		BannedUntil: time.Now().Add(time.Minute),
	}
	r.Add(v)

	all := r.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, "10.0.0.1", all[0].IP)
}

func TestViolatorsRepo_Remove(t *testing.T) {
	r := NewViolatorsRepo()
	r.Add(domain.Violator{IP: "10.0.0.2", BannedUntil: time.Now().Add(time.Minute)})
	r.Add(domain.Violator{IP: "10.0.0.3", BannedUntil: time.Now().Add(time.Minute)})

	r.Remove("10.0.0.2")

	all := r.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, "10.0.0.3", all[0].IP)
}

func TestViolatorsRepo_IsBanned_Active(t *testing.T) {
	r := NewViolatorsRepo()
	r.Add(domain.Violator{IP: "10.0.0.4", BannedUntil: time.Now().Add(time.Minute)})

	assert.True(t, r.IsBanned("10.0.0.4"))
}

func TestViolatorsRepo_IsBanned_Expired(t *testing.T) {
	r := NewViolatorsRepo()
	r.Add(domain.Violator{IP: "10.0.0.5", BannedUntil: time.Now().Add(-time.Second)})

	assert.False(t, r.IsBanned("10.0.0.5"))
}

func TestViolatorsRepo_IsBanned_NotFound(t *testing.T) {
	r := NewViolatorsRepo()
	assert.False(t, r.IsBanned("9.9.9.9"))
}

func TestViolatorsRepo_GetAll_Empty(t *testing.T) {
	r := NewViolatorsRepo()
	assert.Empty(t, r.GetAll())
}
