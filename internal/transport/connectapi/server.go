// Package connectapi sert les services gRPC du serveur en Connect et en
// gRPC-Web sur le serveur HTTP, pour le dashboard (ADR-0006).
//
// Les implémentations restent celles de grpc-go : elles s'enregistrent ici
// par vanguardgrpc. Les intercepteurs grpc-go ne sont pas rejoués par ce
// pont : la reprise après panique, la journalisation et surtout le contrôle
// d'accès (auth.Guard) sont donc rebranchés en intercepteurs Connect.
package connectapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/vanguard/vanguardgrpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
)

// MaxMessageBytes borne la taille d'un message reçu, comme la limite par
// défaut de grpc-go (4 Mio) sur la porte gRPC.
const MaxMessageBytes = 4 << 20

// Server est la porte Connect du serveur.
type Server struct {
	connect *connect.Server
}

// New crée la porte. Le garde est obligatoire : il applique la même
// politique d'accès que sur la porte gRPC.
func New(guard *auth.Guard, logger zerolog.Logger) *Server {
	componentLogger := logger.With().Str("component", "transport.connect").Logger()
	return &Server{connect: connect.NewServer(
		recovery(componentLogger),
		logging(componentLogger),
		authorize(guard),
	)}
}

// Registrar reçoit les services comme un *grpc.Server :
// forgenetv1.RegisterXServiceServer(s.Registrar(), impl).
func (s *Server) Registrar() grpc.ServiceRegistrar {
	return vanguardgrpc.NewServiceRegistrar(s.connect)
}

// Mount installe une route par méthode (« /forgenet.v1.Service/Méthode »)
// sur mux. À appeler après l'enregistrement des services.
func (s *Server) Mount(mux *http.ServeMux) {
	streaming := make(map[string]bool)
	for spec := range s.connect.Specs() {
		if spec.StreamType != connect.StreamTypeUnary {
			streaming[spec.Procedure] = true
		}
	}
	connecthttp.Mount(&streamDeadlineMux{mux: mux, streaming: streaming}, s.connect,
		connecthttp.WithReadMaxBytes(MaxMessageBytes))
}

// streamDeadlineMux lève le délai d'écriture du serveur HTTP sur les flux :
// sans cela, un flux de journaux serait coupé au bout de 15 secondes
// (server.write_timeout). Les appels simples gardent ce délai.
type streamDeadlineMux struct {
	mux       *http.ServeMux
	streaming map[string]bool
}

func (m *streamDeadlineMux) Handle(pattern string, h http.Handler) {
	if !m.streaming[pattern] {
		m.mux.Handle(pattern, h)
		return
	}
	m.mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
		h.ServeHTTP(w, r)
	}))
}

// authorize applique auth.Guard avant tout décodage du message : un appel
// refusé ne coûte ni décompression ni désérialisation.
func authorize(guard *auth.Guard) connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			var token, clientAddr string
			if info, ok := connect.CallInfoForServerContext(ctx); ok {
				token = auth.BearerToken(info.RequestHeader().Get("Authorization"))
				clientAddr = info.PeerAddr
			}
			ctx, err := guard.Authorize(ctx, spec.Procedure, token, clientAddr)
			if err != nil {
				return toConnectError(err)
			}
			return next(ctx, spec, stream)
		}
	}
}

// toConnectError convertit une erreur gRPC du garde : les codes Connect
// reprennent ceux de gRPC, valeur pour valeur.
func toConnectError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return connect.NewError(connect.CodeInternal, "erreur interne du serveur")
	}
	return connect.NewError(connect.Code(st.Code()), st.Message()) //nolint:gosec // codes gRPC de 0 à 16
}

func logging(logger zerolog.Logger) connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			start := time.Now()
			err := next(ctx, spec, stream)
			code := connect.CodeOf(err)
			event := logger.Debug()
			if err != nil {
				event = logger.Warn()
			}
			event.Str("procedure", spec.Procedure).Str("code", codeName(err, code)).
				Dur("duration", time.Since(start)).Msg("appel Connect traité")
			return err
		}
	}
}

func codeName(err error, code connect.Code) string {
	if err == nil {
		return "ok"
	}
	return strings.ToLower(code.String())
}

// recovery transforme une panique en erreur interne ; son contenu reste dans
// le journal, il peut révéler des détails internes.
func recovery(logger zerolog.Logger) connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error().Str("procedure", spec.Procedure).Interface("panic", r).
						Msg("panique interceptée dans un handler Connect")
					err = connect.NewError(connect.CodeInternal, "erreur interne du serveur")
				}
			}()
			return next(ctx, spec, stream)
		}
	}
}
