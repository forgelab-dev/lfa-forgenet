package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor journalise chaque appel simple avec sa durée et son code.
func LoggingUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		logCall(logger, info.FullMethod, start, err)
		return resp, err
	}
}

// LoggingStreamInterceptor journalise chaque flux à sa fermeture.
func LoggingStreamInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		logCall(logger, info.FullMethod, start, err)
		return err
	}
}

// RecoveryUnaryInterceptor transforme une panique d'un handler en erreur
// Internal, pour qu'un appel défaillant ne fasse pas tomber le serveur.
func RecoveryUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer recoverPanic(logger, info.FullMethod, &err)
		return handler(ctx, req)
	}
}

// RecoveryStreamInterceptor fait de même pour les flux.
func RecoveryStreamInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer recoverPanic(logger, info.FullMethod, &err)
		return handler(srv, ss)
	}
}

func logCall(logger zerolog.Logger, method string, start time.Time, err error) {
	code := status.Code(err)
	event := logger.Debug()
	if code != codes.OK {
		event = logger.Warn()
	}
	event.Str("grpc_method", method).Str("grpc_code", code.String()).Dur("duration", time.Since(start)).Msg("appel gRPC traité")
}

// recoverPanic doit être appelée en defer. Le contenu de la panique reste dans
// le journal : il peut révéler des détails internes et n'aide pas le client.
func recoverPanic(logger zerolog.Logger, method string, err *error) {
	if r := recover(); r != nil {
		logger.Error().Str("grpc_method", method).Interface("panic", r).Msg("panique interceptée dans un handler gRPC")
		*err = status.Error(codes.Internal, "erreur interne du serveur")
	}
}
