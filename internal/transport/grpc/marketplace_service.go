package grpc

import (
	"context"
	"crypto/rand"
	"errors"
	"io/fs"
	"slices"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// MarketplaceServiceServer présente le catalogue des templates et les
// installe, selon le contrat forgenet.v1.MarketplaceService (ADR-0018,
// ADR-0019). Il est réservé aux administrateurs par la règle par défaut
// d'auth.DefaultPolicy, comme les déploiements.
type MarketplaceServiceServer struct {
	forgenetv1.UnimplementedMarketplaceServiceServer
	catalog     *marketplace.Catalog
	deployments *deploy.Service
	logger      zerolog.Logger
}

var _ forgenetv1.MarketplaceServiceServer = (*MarketplaceServiceServer)(nil)

// NewMarketplaceServiceServer instancie le service gRPC du Marketplace.
func NewMarketplaceServiceServer(catalog *marketplace.Catalog, deployments *deploy.Service, logger zerolog.Logger) *MarketplaceServiceServer {
	return &MarketplaceServiceServer{
		catalog:     catalog,
		deployments: deployments,
		logger:      logger.With().Str("component", "grpc.marketplace_service").Logger(),
	}
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
	t, err := s.lookup(id)
	if err != nil {
		return nil, err
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

// InstallTemplate vérifie les valeurs, tire au hasard les secrets laissés
// vides et crée un déploiement ordinaire du fichier Compose du catalogue :
// le navigateur ne l'envoie jamais, il ne peut donc pas le modifier au
// passage. Les secrets générés ne sont renvoyés qu'à la création.
func (s *MarketplaceServiceServer) InstallTemplate(ctx context.Context, req *forgenetv1.InstallTemplateRequest) (*forgenetv1.InstallTemplateResponse, error) {
	t, err := s.lookup(req.GetTemplateId())
	if err != nil {
		return nil, err
	}
	values, err := t.Manifest.Resolve(req.GetValues(), rand.Reader)
	if err != nil {
		if errors.Is(err, marketplace.ErrInvalidValue) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		s.logger.Error().Err(err).Str("template", t.ID).Msg("tirage des secrets impossible")
		return nil, status.Error(codes.Internal, "tirage des secrets impossible")
	}
	content, err := s.catalog.Compose(t.ID)
	if err != nil {
		s.logger.Error().Err(err).Str("template", t.ID).Msg("lecture du fichier Compose impossible")
		return nil, status.Errorf(codes.Internal, "fichier Compose du template %q illisible", t.ID)
	}

	vars := make([]deploy.Variable, 0, len(values))
	for _, v := range values {
		vars = append(vars, deploy.Variable{Name: v.Name, Value: v.Value, Secret: v.Secret})
	}
	d, created, err := s.deployments.CreateOnce(ctx, deploy.CreateRequest{
		ProjectID:   req.GetProjectId(),
		SourceType:  store.SourceCompose,
		Content:     string(content),
		RequestID:   req.GetRequestId(),
		TriggeredBy: "user:" + audit.UserIDFromContext(ctx),
		Template:    t.ID,
		Variables:   vars,
		NewProject:  true,
	})
	if err != nil {
		return nil, deploymentStatus(s.logger, err, "installation du template")
	}

	resp := &forgenetv1.InstallTemplateResponse{Deployment: toProtoDeployment(d)}
	// Rejouée avec la même clé, la demande rend l'installation déjà faite :
	// les secrets tirés cette fois-ci n'ont servi à rien, et ceux de la
	// première fois ne se réaffichent pas.
	if created {
		for _, v := range values {
			if v.Generated {
				resp.GeneratedSecrets = append(resp.GeneratedSecrets, &forgenetv1.GeneratedSecret{
					Name: v.Name, Label: variableLabel(t.Manifest, v.Name), Value: v.Value,
				})
			}
		}
	}
	return resp, nil
}

// variableLabel renvoie le libellé de la variable name du manifeste.
func variableLabel(m marketplace.Manifest, name string) string {
	i := slices.IndexFunc(m.Variables, func(v marketplace.Variable) bool { return v.Name == name })
	if i < 0 {
		return ""
	}
	return m.Variables[i].Label
}

// lookup renvoie le template id, ou le statut d'un identifiant refusé.
func (s *MarketplaceServiceServer) lookup(id string) (marketplace.Template, error) {
	if !marketplace.ValidID(id) {
		return marketplace.Template{}, status.Error(codes.InvalidArgument, "identifiant de template invalide : minuscules, chiffres et -, 31 caractères au plus")
	}
	t, ok := s.catalog.Template(id)
	if !ok {
		return marketplace.Template{}, status.Errorf(codes.NotFound, "template %q absent du catalogue", id)
	}
	return t, nil
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
