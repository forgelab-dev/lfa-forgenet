package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
)

// UnaryServerInterceptor applique l'authentification obligatoire sur tous les appels gRPC,
// à l'exception des méthodes explicitement exemptées (ex: health checks).
func UnaryServerInterceptor(authSvc Service, exemptMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if exemptMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		newCtx, err := authenticate(ctx, authSvc)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// StreamServerInterceptor applique l'authentification obligatoire sur tous les flux de streaming gRPC.
func StreamServerInterceptor(authSvc Service, exemptMethods map[string]bool) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if exemptMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		newCtx, err := authenticate(ss.Context(), authSvc)
		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
		return handler(srv, wrapped)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func authenticate(ctx context.Context, authSvc Service) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "métadonnées de requête gRPC manquantes")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 || strings.TrimSpace(authHeaders[0]) == "" {
		return nil, status.Error(codes.Unauthenticated, "en-tête 'authorization' requis")
	}

	tokenStr := strings.TrimSpace(authHeaders[0])
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}

	claims, err := authSvc.ValidateToken(ctx, tokenStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentification invalide : %v", err)
	}

	clientIP := ""
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		clientIP = p.Addr.String()
	}

	ctx = WithUserContext(ctx, claims)
	ctx = audit.WithAuditContext(ctx, claims.UserID, clientIP)

	return ctx, nil
}
