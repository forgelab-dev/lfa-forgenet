package grpc_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/domains"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// domainProxy simule Caddy ; err est renvoyée à chaque rechargement.
type domainProxy struct {
	mu  sync.Mutex
	err error
}

func (p *domainProxy) Apply(context.Context, []proxy.Route) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.err
}

func (p *domainProxy) fail(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.err = err
}

type domainNetwork struct{}

func (domainNetwork) EnsureNetwork(context.Context, string, map[string]string) error { return nil }
func (domainNetwork) ConnectToNetwork(context.Context, string, string, []string) error {
	return nil
}
func (domainNetwork) DisconnectFromNetwork(context.Context, string, string) error { return nil }

// domainStacks : la stack « boutique » a un service « web ».
type domainStacks struct{}

func (domainStacks) Services(_ context.Context, project string) ([]compose.Service, error) {
	if project == "forgenet-boutique" {
		return []compose.Service{{ID: "c1", Service: "web"}}, nil
	}
	return nil, nil
}

type domainTestServer struct {
	*testServer
	domains forgenetv1.DomainServiceClient
	proxy   *domainProxy
}

// startDomainServer sert DomainService derrière le vrai garde ; enabled faux
// simule un reverse proxy désactivé.
func startDomainServer(t *testing.T, enabled bool) *domainTestServer {
	t.Helper()
	dts := &domainTestServer{proxy: &domainProxy{}}
	dts.testServer = startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		var svc *domains.Service
		if enabled {
			svc = domains.New(st.Domains(), dts.proxy, domainNetwork{}, domainStacks{}, auditor,
				domains.Config{Network: "forgenet-edge"}, zerolog.Nop())
		}
		forgenetv1.RegisterDomainServiceServer(srv, transportGRPC.NewDomainServiceServer(svc, zerolog.Nop()))
	})
	dts.domains = forgenetv1.NewDomainServiceClient(dts.conn)
	return dts
}

func addRequest(host string) *forgenetv1.AddDomainRequest {
	return &forgenetv1.AddDomainRequest{Host: host, ProjectId: "boutique", Service: "web", Port: 8080}
}

func TestDomainService_RolePolicy(t *testing.T) {
	ts := startDomainServer(t, true)

	calls := map[string]func(ctx context.Context) error{
		"List": func(ctx context.Context) error {
			_, err := ts.domains.ListDomains(ctx, &forgenetv1.ListDomainsRequest{})
			return err
		},
		"Add": func(ctx context.Context) error {
			_, err := ts.domains.AddDomain(ctx, addRequest("role.example.org"))
			return err
		},
		"Remove": func(ctx context.Context) error {
			_, err := ts.domains.RemoveDomain(ctx, &forgenetv1.RemoveDomainRequest{DomainId: "dom_absent"})
			return err
		},
	}
	for name, call := range calls {
		assert.Equal(t, codes.Unauthenticated, status.Code(call(context.Background())), "%s sans jeton", name)
		memberCode := status.Code(call(as(ts.memberToken)))
		if name == "List" {
			assert.Equal(t, codes.OK, memberCode, "un membre liste les domaines")
		} else {
			assert.Equal(t, codes.PermissionDenied, memberCode, "%s par un membre", name)
		}
		code := status.Code(call(as(ts.adminToken)))
		assert.NotContains(t, []codes.Code{codes.Unauthenticated, codes.PermissionDenied}, code, "%s par un admin", name)
	}
}

func TestDomainService_AddListRemove(t *testing.T) {
	ts := startDomainServer(t, true)
	ctx := as(ts.adminToken)

	added, err := ts.domains.AddDomain(ctx, addRequest("Shop.Example.org"))
	require.NoError(t, err)
	d := added.GetDomain()
	assert.Equal(t, "shop.example.org", d.GetHost())
	assert.Equal(t, "boutique", d.GetProjectId())
	assert.Equal(t, "web", d.GetService())
	assert.Equal(t, uint32(8080), d.GetPort())
	assert.True(t, strings.HasPrefix(d.GetCreatedBy(), "user:usr_"), d.GetCreatedBy())
	assert.NotNil(t, d.GetCreatedAt())

	listed, err := ts.domains.ListDomains(as(ts.memberToken), &forgenetv1.ListDomainsRequest{ProjectId: "boutique"})
	require.NoError(t, err)
	require.Len(t, listed.GetDomains(), 1)
	assert.Equal(t, d.GetId(), listed.GetDomains()[0].GetId())

	_, err = ts.domains.RemoveDomain(ctx, &forgenetv1.RemoveDomainRequest{DomainId: d.GetId()})
	require.NoError(t, err)
	listed, err = ts.domains.ListDomains(ctx, &forgenetv1.ListDomainsRequest{})
	require.NoError(t, err)
	assert.Empty(t, listed.GetDomains())

	events, err := ts.store.Audit().List(context.Background(), 10, 0)
	require.NoError(t, err)
	var actions []string
	for _, e := range events {
		actions = append(actions, e.Action)
	}
	assert.Contains(t, actions, "domain.add")
	assert.Contains(t, actions, "domain.remove")
}

func TestDomainService_ErrorCodes(t *testing.T) {
	ts := startDomainServer(t, true)
	ctx := as(ts.adminToken)
	_, err := ts.domains.AddDomain(ctx, addRequest("pris.example.org"))
	require.NoError(t, err)

	cases := []struct {
		name string
		req  *forgenetv1.AddDomainRequest
		want codes.Code
	}{
		{"domaine invalide", addRequest("pas un domaine"), codes.InvalidArgument},
		{"port hors bornes", &forgenetv1.AddDomainRequest{Host: "a.example.org", ProjectId: "boutique", Service: "web", Port: 70000}, codes.InvalidArgument},
		{"service non déployé", &forgenetv1.AddDomainRequest{Host: "a.example.org", ProjectId: "boutique", Service: "api", Port: 80}, codes.NotFound},
		{"domaine déjà pris", addRequest("PRIS.example.org"), codes.AlreadyExists},
	}
	for _, c := range cases {
		_, err := ts.domains.AddDomain(ctx, c.req)
		assert.Equal(t, c.want, status.Code(err), c.name)
	}

	_, err = ts.domains.RemoveDomain(ctx, &forgenetv1.RemoveDomainRequest{DomainId: "dom_absent"})
	assert.Equal(t, codes.NotFound, status.Code(err))

	ts.proxy.fail(&proxy.LoadError{Message: "loading config: refusé"})
	_, err = ts.domains.AddDomain(ctx, addRequest("b.example.org"))
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	ts.proxy.fail(fmt.Errorf("%w : connexion refusée", proxy.ErrUnavailable))
	_, err = ts.domains.AddDomain(ctx, addRequest("c.example.org"))
	assert.Equal(t, codes.Unavailable, status.Code(err))
	assert.NotContains(t, status.Convert(err).Message(), "connexion refusée", "le détail reste dans le journal")
}

func TestDomainService_Disabled(t *testing.T) {
	ts := startDomainServer(t, false)
	ctx := as(ts.adminToken)

	_, err := ts.domains.ListDomains(ctx, &forgenetv1.ListDomainsRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "proxy.enabled")
	_, err = ts.domains.AddDomain(ctx, addRequest("a.example.org"))
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = ts.domains.RemoveDomain(ctx, &forgenetv1.RemoveDomainRequest{DomainId: "dom_1"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}
