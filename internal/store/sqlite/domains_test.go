package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func newDomain(id, host, project string) *store.Domain {
	return &store.Domain{ID: id, Host: host, ProjectID: project, Service: "web", Port: 8080, CreatedBy: "user:usr_1"}
}

func TestDomains_CreateGetDelete(t *testing.T) {
	ctx := context.Background()
	domains := newTestStore(t).Domains()

	d := newDomain("dom_1", "app.example.org", "boutique")
	require.NoError(t, domains.CreateDomain(ctx, d))
	assert.False(t, d.CreatedAt.IsZero(), "la date de création est fixée à l'insertion")

	got, err := domains.GetDomain(ctx, "dom_1")
	require.NoError(t, err)
	assert.Equal(t, "app.example.org", got.Host)
	assert.Equal(t, "boutique", got.ProjectID)
	assert.Equal(t, "web", got.Service)
	assert.Equal(t, 8080, got.Port)
	assert.Equal(t, "user:usr_1", got.CreatedBy)
	assert.WithinDuration(t, d.CreatedAt, got.CreatedAt, time.Microsecond)

	require.NoError(t, domains.DeleteDomain(ctx, "dom_1"))
	_, err = domains.GetDomain(ctx, "dom_1")
	assert.ErrorIs(t, err, store.ErrNotFound)
	assert.ErrorIs(t, domains.DeleteDomain(ctx, "dom_1"), store.ErrNotFound)
}

func TestDomains_HostIsUnique(t *testing.T) {
	ctx := context.Background()
	domains := newTestStore(t).Domains()

	require.NoError(t, domains.CreateDomain(ctx, newDomain("dom_1", "app.example.org", "boutique")))
	err := domains.CreateDomain(ctx, newDomain("dom_2", "app.example.org", "vitrine"))
	assert.ErrorIs(t, err, store.ErrConflict)
}

func TestDomains_ListByProjectSortedByHost(t *testing.T) {
	ctx := context.Background()
	domains := newTestStore(t).Domains()

	for _, d := range []*store.Domain{
		newDomain("dom_1", "www.example.org", "vitrine"),
		newDomain("dom_2", "b.example.org", "boutique"),
		newDomain("dom_3", "a.example.org", "boutique"),
	} {
		require.NoError(t, domains.CreateDomain(ctx, d))
	}

	all, err := domains.ListDomains(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, []string{"a.example.org", "b.example.org", "www.example.org"}, hosts(all))

	boutique, err := domains.ListDomains(ctx, "boutique")
	require.NoError(t, err)
	assert.Equal(t, []string{"a.example.org", "b.example.org"}, hosts(boutique))

	none, err := domains.ListDomains(ctx, "absent")
	require.NoError(t, err)
	assert.Empty(t, none)
}

func TestDomains_RejectsInvalidPort(t *testing.T) {
	d := newDomain("dom_1", "app.example.org", "boutique")
	d.Port = 70000
	assert.Error(t, newTestStore(t).Domains().CreateDomain(context.Background(), d))
}

func hosts(list []*store.Domain) []string {
	out := make([]string, 0, len(list))
	for _, d := range list {
		out = append(out, d.Host)
	}
	return out
}
