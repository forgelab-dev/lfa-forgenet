package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor journalise chaque appel gRPC avec sa durée et son code de retour.
func LoggingUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		logEvent := logger.Debug()
		if st.Code() != codes.OK {
			logEvent = logger.Warn()
		}

		logEvent.
			Str("grpc_method", info.FullMethod).
			Str("grpc_code", st.Code().String()).
			Dur("duration", duration).
			Msg("appel gRPC unitaire traité")

		return resp, err
	}
}

// RecoveryUnaryInterceptor intercepte les paniques au sein des handlers gRPC pour éviter le crash du serveur.
func RecoveryUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error().
					Str("grpc_method", info.FullMethod).
					Interface("panic", r).
					Msg("panique interceptée dans un handler gRPC")
				// Le contenu de la panique reste dans le journal : il peut révéler
				// des détails internes et n'aide pas le client.
				err = status.Error(codes.Internal, "erreur interne du serveur")
			}
		}()

		return handler(ctx, req)
	}
}
