// This is a generated file - do not edit.
//
// Generated from forgenet/v1/metrics.proto.

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

import 'metrics.pb.dart' as $0;

export 'metrics.pb.dart';

/// MetricsService fournit la collecte de télémesures en temps réel, en streaming et historiques.
@$pb.GrpcServiceName('forgenet.v1.MetricsService')
class MetricsServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  MetricsServiceClient(super.channel, {super.options, super.interceptors});

  /// GetContainerMetrics récupère les métriques instantanées d'un conteneur.
  $grpc.ResponseFuture<$0.GetContainerMetricsResponse> getContainerMetrics(
    $0.GetContainerMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContainerMetrics, request, options: options);
  }

  /// StreamContainerMetrics transmet en temps réel le flux de métriques d'un conteneur.
  $grpc.ResponseStream<$0.StreamContainerMetricsResponse>
      streamContainerMetrics(
    $0.StreamContainerMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamContainerMetrics, $async.Stream.fromIterable([request]),
        options: options);
  }

  /// QueryContainerMetrics interroge les séries temporelles historiques d'un conteneur avec résolution ajustable et garde-fous de bande passante.
  $grpc.ResponseFuture<$0.QueryContainerMetricsResponse> queryContainerMetrics(
    $0.QueryContainerMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$queryContainerMetrics, request, options: options);
  }

  /// GetSystemMetrics récupère l'état d'utilisation instantané du système hôte.
  $grpc.ResponseFuture<$0.GetSystemMetricsResponse> getSystemMetrics(
    $0.GetSystemMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSystemMetrics, request, options: options);
  }

  /// StreamSystemMetrics transmet en continu l'état des ressources globales du nœud hôte.
  $grpc.ResponseStream<$0.StreamSystemMetricsResponse> streamSystemMetrics(
    $0.StreamSystemMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamSystemMetrics, $async.Stream.fromIterable([request]),
        options: options);
  }

  /// QuerySystemMetrics interroge les séries temporelles historiques du système hôte.
  $grpc.ResponseFuture<$0.QuerySystemMetricsResponse> querySystemMetrics(
    $0.QuerySystemMetricsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$querySystemMetrics, request, options: options);
  }

  // method descriptors

  static final _$getContainerMetrics = $grpc.ClientMethod<
          $0.GetContainerMetricsRequest, $0.GetContainerMetricsResponse>(
      '/forgenet.v1.MetricsService/GetContainerMetrics',
      ($0.GetContainerMetricsRequest value) => value.writeToBuffer(),
      $0.GetContainerMetricsResponse.fromBuffer);
  static final _$streamContainerMetrics = $grpc.ClientMethod<
          $0.StreamContainerMetricsRequest, $0.StreamContainerMetricsResponse>(
      '/forgenet.v1.MetricsService/StreamContainerMetrics',
      ($0.StreamContainerMetricsRequest value) => value.writeToBuffer(),
      $0.StreamContainerMetricsResponse.fromBuffer);
  static final _$queryContainerMetrics = $grpc.ClientMethod<
          $0.QueryContainerMetricsRequest, $0.QueryContainerMetricsResponse>(
      '/forgenet.v1.MetricsService/QueryContainerMetrics',
      ($0.QueryContainerMetricsRequest value) => value.writeToBuffer(),
      $0.QueryContainerMetricsResponse.fromBuffer);
  static final _$getSystemMetrics = $grpc.ClientMethod<
          $0.GetSystemMetricsRequest, $0.GetSystemMetricsResponse>(
      '/forgenet.v1.MetricsService/GetSystemMetrics',
      ($0.GetSystemMetricsRequest value) => value.writeToBuffer(),
      $0.GetSystemMetricsResponse.fromBuffer);
  static final _$streamSystemMetrics = $grpc.ClientMethod<
          $0.StreamSystemMetricsRequest, $0.StreamSystemMetricsResponse>(
      '/forgenet.v1.MetricsService/StreamSystemMetrics',
      ($0.StreamSystemMetricsRequest value) => value.writeToBuffer(),
      $0.StreamSystemMetricsResponse.fromBuffer);
  static final _$querySystemMetrics = $grpc.ClientMethod<
          $0.QuerySystemMetricsRequest, $0.QuerySystemMetricsResponse>(
      '/forgenet.v1.MetricsService/QuerySystemMetrics',
      ($0.QuerySystemMetricsRequest value) => value.writeToBuffer(),
      $0.QuerySystemMetricsResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.MetricsService')
abstract class MetricsServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.MetricsService';

  MetricsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetContainerMetricsRequest,
            $0.GetContainerMetricsResponse>(
        'GetContainerMetrics',
        getContainerMetrics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContainerMetricsRequest.fromBuffer(value),
        ($0.GetContainerMetricsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamContainerMetricsRequest,
            $0.StreamContainerMetricsResponse>(
        'StreamContainerMetrics',
        streamContainerMetrics_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.StreamContainerMetricsRequest.fromBuffer(value),
        ($0.StreamContainerMetricsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.QueryContainerMetricsRequest,
            $0.QueryContainerMetricsResponse>(
        'QueryContainerMetrics',
        queryContainerMetrics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.QueryContainerMetricsRequest.fromBuffer(value),
        ($0.QueryContainerMetricsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSystemMetricsRequest,
            $0.GetSystemMetricsResponse>(
        'GetSystemMetrics',
        getSystemMetrics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSystemMetricsRequest.fromBuffer(value),
        ($0.GetSystemMetricsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamSystemMetricsRequest,
            $0.StreamSystemMetricsResponse>(
        'StreamSystemMetrics',
        streamSystemMetrics_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.StreamSystemMetricsRequest.fromBuffer(value),
        ($0.StreamSystemMetricsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.QuerySystemMetricsRequest,
            $0.QuerySystemMetricsResponse>(
        'QuerySystemMetrics',
        querySystemMetrics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.QuerySystemMetricsRequest.fromBuffer(value),
        ($0.QuerySystemMetricsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.GetContainerMetricsResponse> getContainerMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetContainerMetricsRequest> $request) async {
    return getContainerMetrics($call, await $request);
  }

  $async.Future<$0.GetContainerMetricsResponse> getContainerMetrics(
      $grpc.ServiceCall call, $0.GetContainerMetricsRequest request);

  $async.Stream<$0.StreamContainerMetricsResponse> streamContainerMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StreamContainerMetricsRequest> $request) async* {
    yield* streamContainerMetrics($call, await $request);
  }

  $async.Stream<$0.StreamContainerMetricsResponse> streamContainerMetrics(
      $grpc.ServiceCall call, $0.StreamContainerMetricsRequest request);

  $async.Future<$0.QueryContainerMetricsResponse> queryContainerMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.QueryContainerMetricsRequest> $request) async {
    return queryContainerMetrics($call, await $request);
  }

  $async.Future<$0.QueryContainerMetricsResponse> queryContainerMetrics(
      $grpc.ServiceCall call, $0.QueryContainerMetricsRequest request);

  $async.Future<$0.GetSystemMetricsResponse> getSystemMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSystemMetricsRequest> $request) async {
    return getSystemMetrics($call, await $request);
  }

  $async.Future<$0.GetSystemMetricsResponse> getSystemMetrics(
      $grpc.ServiceCall call, $0.GetSystemMetricsRequest request);

  $async.Stream<$0.StreamSystemMetricsResponse> streamSystemMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StreamSystemMetricsRequest> $request) async* {
    yield* streamSystemMetrics($call, await $request);
  }

  $async.Stream<$0.StreamSystemMetricsResponse> streamSystemMetrics(
      $grpc.ServiceCall call, $0.StreamSystemMetricsRequest request);

  $async.Future<$0.QuerySystemMetricsResponse> querySystemMetrics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.QuerySystemMetricsRequest> $request) async {
    return querySystemMetrics($call, await $request);
  }

  $async.Future<$0.QuerySystemMetricsResponse> querySystemMetrics(
      $grpc.ServiceCall call, $0.QuerySystemMetricsRequest request);
}
