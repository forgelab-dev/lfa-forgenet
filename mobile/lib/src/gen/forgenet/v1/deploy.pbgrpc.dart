// This is a generated file - do not edit.
//
// Generated from forgenet/v1/deploy.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'deploy.pb.dart' as $0;

export 'deploy.pb.dart';

/// DeploymentService gère le cycle de vie asynchrone des déploiements des applications.
@$pb.GrpcServiceName('forgenet.v1.DeploymentService')
class DeploymentServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DeploymentServiceClient(super.channel, {super.options, super.interceptors});

  /// CreateDeployment initie un déploiement de manière asynchrone (retourne immédiatement l'état PENDING).
  $grpc.ResponseFuture<$0.CreateDeploymentResponse> createDeployment(
    $0.CreateDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDeployment, request, options: options);
  }

  /// GetDeployment récupère les métadonnées d'un déploiement donné.
  $grpc.ResponseFuture<$0.GetDeploymentResponse> getDeployment(
    $0.GetDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeployment, request, options: options);
  }

  /// ListDeployments liste l'historique des déploiements d'un projet.
  $grpc.ResponseFuture<$0.ListDeploymentsResponse> listDeployments(
    $0.ListDeploymentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDeployments, request, options: options);
  }

  /// CancelDeployment interrompt un déploiement en cours s'il est dans un état annulable (PENDING, BUILDING, DEPLOYING).
  $grpc.ResponseFuture<$0.CancelDeploymentResponse> cancelDeployment(
    $0.CancelDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelDeployment, request, options: options);
  }

  /// RollbackDeployment initie un nouveau déploiement basé sur un état stable antérieur.
  $grpc.ResponseFuture<$0.RollbackDeploymentResponse> rollbackDeployment(
    $0.RollbackDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rollbackDeployment, request, options: options);
  }

  /// StreamDeploymentLogs fournit en temps réel les logs de build et de déploiement (STDOUT/STDERR).
  $grpc.ResponseStream<$0.StreamDeploymentLogsResponse> streamDeploymentLogs(
    $0.StreamDeploymentLogsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamDeploymentLogs, $async.Stream.fromIterable([request]),
        options: options);
  }

  /// StreamDeployment transmet en streaming les mises à jour complètes de l'état du déploiement.
  $grpc.ResponseStream<$0.StreamDeploymentResponse> streamDeployment(
    $0.StreamDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamDeployment, $async.Stream.fromIterable([request]),
        options: options);
  }

  // method descriptors

  static final _$createDeployment = $grpc.ClientMethod<
          $0.CreateDeploymentRequest, $0.CreateDeploymentResponse>(
      '/forgenet.v1.DeploymentService/CreateDeployment',
      ($0.CreateDeploymentRequest value) => value.writeToBuffer(),
      $0.CreateDeploymentResponse.fromBuffer);
  static final _$getDeployment =
      $grpc.ClientMethod<$0.GetDeploymentRequest, $0.GetDeploymentResponse>(
          '/forgenet.v1.DeploymentService/GetDeployment',
          ($0.GetDeploymentRequest value) => value.writeToBuffer(),
          $0.GetDeploymentResponse.fromBuffer);
  static final _$listDeployments =
      $grpc.ClientMethod<$0.ListDeploymentsRequest, $0.ListDeploymentsResponse>(
          '/forgenet.v1.DeploymentService/ListDeployments',
          ($0.ListDeploymentsRequest value) => value.writeToBuffer(),
          $0.ListDeploymentsResponse.fromBuffer);
  static final _$cancelDeployment = $grpc.ClientMethod<
          $0.CancelDeploymentRequest, $0.CancelDeploymentResponse>(
      '/forgenet.v1.DeploymentService/CancelDeployment',
      ($0.CancelDeploymentRequest value) => value.writeToBuffer(),
      $0.CancelDeploymentResponse.fromBuffer);
  static final _$rollbackDeployment = $grpc.ClientMethod<
          $0.RollbackDeploymentRequest, $0.RollbackDeploymentResponse>(
      '/forgenet.v1.DeploymentService/RollbackDeployment',
      ($0.RollbackDeploymentRequest value) => value.writeToBuffer(),
      $0.RollbackDeploymentResponse.fromBuffer);
  static final _$streamDeploymentLogs = $grpc.ClientMethod<
          $0.StreamDeploymentLogsRequest, $0.StreamDeploymentLogsResponse>(
      '/forgenet.v1.DeploymentService/StreamDeploymentLogs',
      ($0.StreamDeploymentLogsRequest value) => value.writeToBuffer(),
      $0.StreamDeploymentLogsResponse.fromBuffer);
  static final _$streamDeployment = $grpc.ClientMethod<
          $0.StreamDeploymentRequest, $0.StreamDeploymentResponse>(
      '/forgenet.v1.DeploymentService/StreamDeployment',
      ($0.StreamDeploymentRequest value) => value.writeToBuffer(),
      $0.StreamDeploymentResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.DeploymentService')
abstract class DeploymentServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.DeploymentService';

  DeploymentServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateDeploymentRequest,
            $0.CreateDeploymentResponse>(
        'CreateDeployment',
        createDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDeploymentRequest.fromBuffer(value),
        ($0.CreateDeploymentResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetDeploymentRequest, $0.GetDeploymentResponse>(
            'GetDeployment',
            getDeployment_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetDeploymentRequest.fromBuffer(value),
            ($0.GetDeploymentResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDeploymentsRequest,
            $0.ListDeploymentsResponse>(
        'ListDeployments',
        listDeployments_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDeploymentsRequest.fromBuffer(value),
        ($0.ListDeploymentsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelDeploymentRequest,
            $0.CancelDeploymentResponse>(
        'CancelDeployment',
        cancelDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelDeploymentRequest.fromBuffer(value),
        ($0.CancelDeploymentResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RollbackDeploymentRequest,
            $0.RollbackDeploymentResponse>(
        'RollbackDeployment',
        rollbackDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RollbackDeploymentRequest.fromBuffer(value),
        ($0.RollbackDeploymentResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamDeploymentLogsRequest,
            $0.StreamDeploymentLogsResponse>(
        'StreamDeploymentLogs',
        streamDeploymentLogs_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.StreamDeploymentLogsRequest.fromBuffer(value),
        ($0.StreamDeploymentLogsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamDeploymentRequest,
            $0.StreamDeploymentResponse>(
        'StreamDeployment',
        streamDeployment_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.StreamDeploymentRequest.fromBuffer(value),
        ($0.StreamDeploymentResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CreateDeploymentResponse> createDeployment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDeploymentRequest> $request) async {
    return createDeployment($call, await $request);
  }

  $async.Future<$0.CreateDeploymentResponse> createDeployment(
      $grpc.ServiceCall call, $0.CreateDeploymentRequest request);

  $async.Future<$0.GetDeploymentResponse> getDeployment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDeploymentRequest> $request) async {
    return getDeployment($call, await $request);
  }

  $async.Future<$0.GetDeploymentResponse> getDeployment(
      $grpc.ServiceCall call, $0.GetDeploymentRequest request);

  $async.Future<$0.ListDeploymentsResponse> listDeployments_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDeploymentsRequest> $request) async {
    return listDeployments($call, await $request);
  }

  $async.Future<$0.ListDeploymentsResponse> listDeployments(
      $grpc.ServiceCall call, $0.ListDeploymentsRequest request);

  $async.Future<$0.CancelDeploymentResponse> cancelDeployment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelDeploymentRequest> $request) async {
    return cancelDeployment($call, await $request);
  }

  $async.Future<$0.CancelDeploymentResponse> cancelDeployment(
      $grpc.ServiceCall call, $0.CancelDeploymentRequest request);

  $async.Future<$0.RollbackDeploymentResponse> rollbackDeployment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RollbackDeploymentRequest> $request) async {
    return rollbackDeployment($call, await $request);
  }

  $async.Future<$0.RollbackDeploymentResponse> rollbackDeployment(
      $grpc.ServiceCall call, $0.RollbackDeploymentRequest request);

  $async.Stream<$0.StreamDeploymentLogsResponse> streamDeploymentLogs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StreamDeploymentLogsRequest> $request) async* {
    yield* streamDeploymentLogs($call, await $request);
  }

  $async.Stream<$0.StreamDeploymentLogsResponse> streamDeploymentLogs(
      $grpc.ServiceCall call, $0.StreamDeploymentLogsRequest request);

  $async.Stream<$0.StreamDeploymentResponse> streamDeployment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StreamDeploymentRequest> $request) async* {
    yield* streamDeployment($call, await $request);
  }

  $async.Stream<$0.StreamDeploymentResponse> streamDeployment(
      $grpc.ServiceCall call, $0.StreamDeploymentRequest request);
}
