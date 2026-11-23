package grpc

import (
	"context"
	"errors"
	"io/fs"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// MarketplaceServiceServer présente le catalogue des templates selon le
// contrat forgenet.v1.MarketplaceService (ADR-0018). Il est réservé aux
// administrateurs par la règle par défaut d'auth.DefaultPolicy.
type MarketplaceServiceServer struct {
	forgenetv1.UnimplementedMarketplaceServiceServer
	catalog *marketplace.Catalog
	logger  zerolog.Logger
}

var _ forgenetv1.MarketplaceServiceServer = (*MarketplaceServiceServer)(nil)

// NewMarketplaceServiceServer instancie le service gRPC du Marketplace.
func NewMarketplaceServiceServer(catalog *marketplace.Catalog, logger zerolog.Logger) *MarketplaceServiceServer {
	return &MarketplaceServiceServer{catalog: catalog, logger: logger.With().Str("component", "grpc.marketplace_service").Logger()}
}

// ListTemplates donne tout le catalogue.
func (s *MarketplaceServiceServer) ListTemplates(_ context.Context, _ *forgenetv1.ListTemplatesRequest) (*forgenetv1.ListTemplatesResponse, error) {
	list := s.catalog.Templates()
	out := make([]*forgenetv1.Template, 0, len(list))
	for _, t := range list {
		out = append(out, templateToProto(t))
	}
	return &forgenetv1.ListTemplatesResponse{Templates: out}, nil
}

// GetTemplate donne un template et son README.
func (s *MarketplaceServiceServer) GetTemplate(_ context.Context, req *forgenetv1.GetTemplateRequest) (*forgenetv1.GetTemplateResponse, error) {
	id := req.GetId()
	if !marketplace.ValidID(id) {
		return nil, status.Error(codes.InvalidArgument, "identifiant de template invalide : minuscules, chiffres et -, 31 caractères au plus")
	}
	t, ok := s.catalog.Template(id)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "template %q absent du catalogue", id)
	}
	readme, err := s.catalog.Readme(id)
	if err != nil {
		// Le catalogue est embarqué : un README introuvable ici est une
		// erreur de construction, pas une erreur de l'appelant.
		if !errors.Is(err, fs.ErrNotExist) {
			s.logger.Error().Err(err).Str("template", id).Msg("lecture du README impossible")
		}
		return nil, status.Errorf(codes.Internal, "README du template %q illisible", id)
	}
	return &forgenetv1.GetTemplateResponse{Template: templateToProto(t), Readme: string(readme)}, nil
}

func templateToProto(t marketplace.Template) *forgenetv1.Template {
	m := t.Manifest
	out := &forgenetv1.Template{
		Id:          t.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Version:     m.Version,
		License:     m.License,
		Homepage:    m.Homepage,
		IconSvg:     string(t.Icon),
	}
	for _, v := range m.Variables {
		out.Variables = append(out.Variables, &forgenetv1.TemplateVariable{
			Name:         v.Name,
			Label:        v.Label,
			Help:         v.Help,
			DefaultValue: v.Default,
			Pattern:      v.Pattern,
			Optional:     v.Optional,
			Secret:       v.Secret,
			Generate:     uint32(v.Generate), //nolint:gosec // G115 : 0, ou de 16 à 128 (Manifest.Validate)
		})
	}
	for _, p := range m.Ports {
		out.Ports = append(out.Ports, &forgenetv1.TemplatePort{
			Service:     p.Service,
			Port:        uint32(p.Port), //nolint:gosec // G115 : port vérifié de 1 à 65535 (Manifest.Validate)
			Protocol:    p.Protocol,
			Description: p.Description,
		})
	}
	for _, v := range m.Volumes {
		out.Volumes = append(out.Volumes, &forgenetv1.TemplateVolume{Name: v.Name, Description: v.Description})
	}
	return out
}
