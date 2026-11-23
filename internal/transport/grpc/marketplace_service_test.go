package grpc_test

import (
	"io/fs"
	"path"
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
	"github.com/forgelab-dev/lfa-forgenet/templates"
)

// demoCatalog met le template d'exemple, qui n'est pas au catalogue réel, sous
// le nom « demo ».
func demoCatalog(t *testing.T) fs.FS {
	t.Helper()
	fsys := fstest.MapFS{}
	for _, name := range []string{marketplace.ManifestFile, marketplace.ComposeFile, marketplace.IconFile, marketplace.ReadmeFile} {
		data, err := fs.ReadFile(templates.FS, path.Join("_exemple", name))
		require.NoError(t, err)
		fsys[path.Join("demo", name)] = &fstest.MapFile{Data: data}
	}
	return fsys
}

func startMarketplaceServer(t *testing.T, fsys fs.FS) (*testServer, forgenetv1.MarketplaceServiceClient) {
	t.Helper()
	catalog, err := marketplace.NewCatalog(fsys)
	require.NoError(t, err)
	ts := startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, _ store.Store, _ audit.Auditor) {
		forgenetv1.RegisterMarketplaceServiceServer(srv, transportGRPC.NewMarketplaceServiceServer(catalog, zerolog.Nop()))
	})
	return ts, forgenetv1.NewMarketplaceServiceClient(ts.conn)
}

func TestMarketplaceService_AdminsOnly(t *testing.T) {
	ts, client := startMarketplaceServer(t, demoCatalog(t))
	_, err := client.ListTemplates(as(ts.memberToken), &forgenetv1.ListTemplatesRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err), "installer une application est réservé aux admins, comme déployer")
	_, err = client.GetTemplate(as(ts.memberToken), &forgenetv1.GetTemplateRequest{Id: "demo"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestMarketplaceService_ListTemplates(t *testing.T) {
	ts, client := startMarketplaceServer(t, demoCatalog(t))
	resp, err := client.ListTemplates(as(ts.adminToken), &forgenetv1.ListTemplatesRequest{})
	require.NoError(t, err)
	require.Len(t, resp.GetTemplates(), 1)

	tpl := resp.GetTemplates()[0]
	assert.Equal(t, "demo", tpl.GetId())
	assert.Equal(t, "Page d'exemple", tpl.GetName())
	assert.Equal(t, "other", tpl.GetCategory())
	assert.Equal(t, "1.31.8", tpl.GetVersion())
	assert.Equal(t, "BSD-2-Clause", tpl.GetLicense())
	assert.Equal(t, "https://nginx.org", tpl.GetHomepage())
	assert.Contains(t, tpl.GetIconSvg(), "<svg")

	require.Len(t, tpl.GetVariables(), 1)
	v := tpl.GetVariables()[0]
	assert.Equal(t, "SITE_TITRE", v.GetName())
	assert.Equal(t, "Titre de la page", v.GetLabel())
	assert.Equal(t, "Bonjour depuis ForgeNet", v.GetDefaultValue())
	assert.NotEmpty(t, v.GetPattern())
	assert.NotEmpty(t, v.GetHelp())
	assert.False(t, v.GetSecret())

	require.Len(t, tpl.GetPorts(), 1)
	p := tpl.GetPorts()[0]
	assert.Equal(t, []any{"web", uint32(80), "http", "La page, à relier à un domaine"},
		[]any{p.GetService(), p.GetPort(), p.GetProtocol(), p.GetDescription()})
	require.Len(t, tpl.GetVolumes(), 1)
	assert.Equal(t, "site", tpl.GetVolumes()[0].GetName())
}

func TestMarketplaceService_EmptyCatalog(t *testing.T) {
	ts, client := startMarketplaceServer(t, fstest.MapFS{})
	resp, err := client.ListTemplates(as(ts.adminToken), &forgenetv1.ListTemplatesRequest{})
	require.NoError(t, err)
	assert.Empty(t, resp.GetTemplates(), "un catalogue vide n'est pas une erreur : le dashboard affiche son état vide")
}

func TestMarketplaceService_GetTemplate(t *testing.T) {
	ts, client := startMarketplaceServer(t, demoCatalog(t))
	resp, err := client.GetTemplate(as(ts.adminToken), &forgenetv1.GetTemplateRequest{Id: "demo"})
	require.NoError(t, err)
	assert.Equal(t, "demo", resp.GetTemplate().GetId())
	assert.Contains(t, resp.GetReadme(), "#", "le README est rendu tel quel, en Markdown")

	_, err = client.GetTemplate(as(ts.adminToken), &forgenetv1.GetTemplateRequest{Id: "../demo"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = client.GetTemplate(as(ts.adminToken), &forgenetv1.GetTemplateRequest{Id: "postgres"})
	assert.Equal(t, codes.NotFound, status.Code(err))
}
