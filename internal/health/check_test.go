package health

import (
	"net/netip"
	"testing"
	"time"

	"github.com/qdm12/ddns-updater/internal/constants"
	"github.com/qdm12/ddns-updater/internal/models"
	"github.com/qdm12/ddns-updater/internal/provider"
	"github.com/qdm12/ddns-updater/internal/records"
	"github.com/stretchr/testify/assert"
)

type memoryDatabase []records.Record

func (m memoryDatabase) SelectAll() []records.Record {
	return m
}

type hostnameProvider struct {
	provider.Provider
	hostname string
}

func (p hostnameProvider) String() string {
	return p.hostname
}

func Test_isHealthy_noIP(t *testing.T) {
	t.Parallel()

	const hostname = "record0.example.com"
	db := memoryDatabase{{
		Provider: hostnameProvider{hostname: hostname},
		History: []models.HistoryEvent{
			{IP: netip.MustParseAddr("1.2.3.4"), Time: time.Date(2025, time.May, 1, 12, 0, 0, 0, time.UTC)},
		},
		Status: constants.NOIP,
	}}

	// The resolver is left nil since a failing record needs no lookup.
	err := isHealthy(t.Context(), db, nil)

	assert.ErrorIs(t, err, ErrRecordUpdateFailed)
	assert.ErrorContains(t, err, hostname)
}
