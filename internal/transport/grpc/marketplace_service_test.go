package grpc_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"io/fs"
	"path"
	"regexp"
	"testing"
	"testing/fstest"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
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

// baseManifest décrit un template dont le mot de passe est tiré au hasard.
const baseManifest = `schema: 1
name: Base de test
description: Une base qui demande un mot de passe.
category: database
version: "1"
service: db
license: MIT
homepage: https://example.org
variables:
  - name: DB_PASSWORD
    label: Mot de passe de la base
    secret: true
    generate: 24
`

const baseCompose = `services:
  db:
    image: alpine:3.21
    command: ["sleep", "3600"]
    environment:
      DB_PASSWORD: ${DB_PASSWORD:?}
`

// installCatalog contient « demo » (une variable, sans secret) et « base »
// (un secret généré).
func installCatalog(t *testing.T) fs.FS {
	t.Helper()
	fsys := demoCatalog(t).(fstest.MapFS)
	fsys["base/manifest.yaml"] = &fstest.MapFile{Data: []byte(baseManifest)}
	fsys["base/docker-compose.yml"] = &fstest.MapFile{Data: []byte(baseCompose)}
	fsys["base/icon.svg"] = fsys["demo/icon.svg"]
	fsys["base/README.md"] = &fstest.MapFile{Data: []byte("# Base de test\n")}
	return fsys
}

func testBox(t *testing.T) *sealed.Box {
	t.Helper()
	key := make([]byte, sealed.KeySize)
	_, err := rand.Read(key)
	require.NoError(t, err)
	box, err := sealed.NewBox(key)
	require.NoError(t, err)
	return box
}

func startMarketplaceServer(t *testing.T, fsys fs.FS, opts ...deploy.Option) (*testServer, forgenetv1.MarketplaceServiceClient) {
	t.Helper()
	catalog, err := marketplace.NewCatalog(fsys)
	require.NoError(t, err)
	ts := startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		// L'exécuteur n'est pas lancé : une installation reste en attente.
		deployments := deploy.New(st.Deployments(), &fakeStacks{}, auditor, zerolog.Nop(), opts...)
		forgenetv1.RegisterMarketplaceServiceServer(srv, transportGRPC.NewMarketplaceServiceServer(catalog, deployments, zerolog.Nop()))
	})
	return ts, forgenetv1.NewMarketplaceServiceClient(ts.conn)
}

func TestMarketplaceService_AdminsOnly(t *testing.T) {
	ts, client := startMarketplaceServer(t, demoCatalog(t))
	_, err := client.ListTemplates(as(ts.memberToken), &forgenetv1.ListTemplatesRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err), "installer une application est réservé aux admins, comme déployer")
	_, err = client.GetTemplate(as(ts.memberToken), &forgenetv1.GetTemplateRequest{Id: "demo"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = client.InstallTemplate(as(ts.memberToken), &forgenetv1.InstallTemplateRequest{TemplateId: "demo", ProjectId: "page"})
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

func TestMarketplaceService_InstallTemplate(t *testing.T) {
	ts, client := startMarketplaceServer(t, installCatalog(t), deploy.WithBox(testBox(t)))
	ctx := as(ts.adminToken)

	resp, err := client.InstallTemplate(ctx, &forgenetv1.InstallTemplateRequest{TemplateId: "base", ProjectId: "ma-base"})
	require.NoError(t, err)
	d := resp.GetDeployment()
	assert.Equal(t, "ma-base", d.GetProjectId())
	assert.Equal(t, "base", d.GetTemplate())
	assert.Equal(t, forgenetv1.DeploymentState_DEPLOYMENT_STATE_PENDING, d.GetState(), "un déploiement ordinaire, mis en file")
	assert.Equal(t, baseCompose, d.GetSource().GetCompose().GetComposeContent(), "le fichier du catalogue, jamais un fichier reçu")

	require.Len(t, resp.GetGeneratedSecrets(), 1)
	secret := resp.GetGeneratedSecrets()[0]
	assert.Equal(t, "DB_PASSWORD", secret.GetName())
	assert.Equal(t, "Mot de passe de la base", secret.GetLabel())
	assert.Regexp(t, regexp.MustCompile(`^[A-Za-z0-9]{24}$`), secret.GetValue())

	stored, err := ts.store.Deployments().GetDeployment(context.Background(), d.GetId())
	require.NoError(t, err)
	require.NotNil(t, stored.VariablesSealed)
	assert.False(t, bytes.Contains(stored.VariablesSealed, []byte(secret.GetValue())), "chiffré en base")

	// Relire le déploiement ne redonne pas le secret.
	assert.NotContains(t, d.String(), secret.GetValue())

	_, err = client.InstallTemplate(ctx, &forgenetv1.InstallTemplateRequest{TemplateId: "demo", ProjectId: "ma-base"})
	assert.Equal(t, codes.AlreadyExists, status.Code(err), "une installation ne remplace jamais une application")
}

func TestMarketplaceService_InstallTemplateReplayGivesNoSecret(t *testing.T) {
	ts, client := startMarketplaceServer(t, installCatalog(t), deploy.WithBox(testBox(t)))
	ctx := as(ts.adminToken)
	req := &forgenetv1.InstallTemplateRequest{TemplateId: "base", ProjectId: "ma-base", RequestId: "req-1"}

	first, err := client.InstallTemplate(ctx, req)
	require.NoError(t, err)
	require.Len(t, first.GetGeneratedSecrets(), 1)

	again, err := client.InstallTemplate(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, first.GetDeployment().GetId(), again.GetDeployment().GetId())
	assert.Empty(t, again.GetGeneratedSecrets(), "un secret n'est montré qu'une fois")
}

func TestMarketplaceService_InstallTemplateRefusals(t *testing.T) {
	ts, client := startMarketplaceServer(t, installCatalog(t), deploy.WithBox(testBox(t)))
	ctx := as(ts.adminToken)

	cases := map[string]struct {
		req  *forgenetv1.InstallTemplateRequest
		code codes.Code
		want string
	}{
		"identifiant mal formé": {&forgenetv1.InstallTemplateRequest{TemplateId: "../demo", ProjectId: "p"}, codes.InvalidArgument, "identifiant de template"},
		"template absent":       {&forgenetv1.InstallTemplateRequest{TemplateId: "postgres", ProjectId: "p"}, codes.NotFound, "postgres"},
		"nom d'application":     {&forgenetv1.InstallTemplateRequest{TemplateId: "demo", ProjectId: "Pas Bon"}, codes.InvalidArgument, "identifiant de projet"},
		"valeur hors motif": {&forgenetv1.InstallTemplateRequest{
			TemplateId: "demo", ProjectId: "p", Values: map[string]string{"SITE_TITRE": "<script>"},
		}, codes.InvalidArgument, "SITE_TITRE"},
		"variable inconnue": {&forgenetv1.InstallTemplateRequest{
			TemplateId: "demo", ProjectId: "p", Values: map[string]string{"AUTRE": "x"},
		}, codes.InvalidArgument, "AUTRE"},
	}
	for name, c := range cases {
		_, err := client.InstallTemplate(ctx, c.req)
		assert.Equal(t, c.code, status.Code(err), name)
		assert.Contains(t, status.Convert(err).Message(), c.want, name)
	}

	list, err := ts.store.Deployments().ListDeployments(context.Background(), "", 10, 0)
	require.NoError(t, err)
	assert.Empty(t, list, "rien n'est enregistré")
}

func TestMarketplaceService_InstallTemplateNeedsSecretKey(t *testing.T) {
	ts, client := startMarketplaceServer(t, installCatalog(t))
	_, err := client.InstallTemplate(as(ts.adminToken), &forgenetv1.InstallTemplateRequest{TemplateId: "base", ProjectId: "ma-base"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "FORGENET_SECRET_KEY")
}
