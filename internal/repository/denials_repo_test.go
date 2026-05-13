package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDenialsRepo_RecordAndGetAll(t *testing.T) {
	r := NewDenialsRepo()

	r.Record("1.1.1.1", "blacklist")
	r.Record("2.2.2.2", "graylist")

	all := r.GetAll()
	require.Len(t, all, 2)

	ipMap := make(map[string]string)
	for _, s := range all {
		ipMap[s.IP] = s.Reason
	}
	assert.Equal(t, "blacklist", ipMap["1.1.1.1"])
	assert.Equal(t, "graylist", ipMap["2.2.2.2"])
}

func TestDenialsRepo_Record_IncrementsExistingIP(t *testing.T) {
	r := NewDenialsRepo()

	r.Record("3.3.3.3", "blacklist")
	r.Record("3.3.3.3", "blacklist")
	r.Record("3.3.3.3", "graylist")

	all := r.GetAll()
	require.Len(t, all, 1)
	assert.Equal(t, int64(3), all[0].Denials)
	assert.Equal(t, "graylist", all[0].Reason) // last reason wins
}

func TestDenialsRepo_GetAll_Empty(t *testing.T) {
	r := NewDenialsRepo()
	assert.Empty(t, r.GetAll())
}
