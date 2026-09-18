// This is a generated file - do not edit.
//
// Generated from forgenet/v1/container.proto.

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

import 'container.pb.dart' as $0;

export 'container.pb.dart';

/// ContainerService fournit les méthodes gRPC pour gérer le cycle de vie des conteneurs.
@$pb.GrpcServiceName('forgenet.v1.ContainerService')
class ContainerServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ContainerServiceClient(super.channel, {super.options, super.interceptors});

  /// ListContainers retourne la liste paginée des conteneurs existants.
  $grpc.ResponseFuture<$0.ListContainersResponse> listContainers(
    $0.ListContainersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listContainers, request, options: options);
  }

  /// GetContainer récupère les informations détaillées d'un conteneur spécifique par son ID.
  $grpc.ResponseFuture<$0.GetContainerResponse> getContainer(
    $0.GetContainerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContainer, request, options: options);
  }

  /// StartContainer démarre un conteneur arrêté.
  $grpc.ResponseFuture<$0.StartContainerResponse> startContainer(
    $0.StartContainerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startContainer, request, options: options);
  }

  /// StopContainer arrête proprement un conteneur en cours d'exécution avec un délai d'attente (grace period).
  $grpc.ResponseFuture<$0.StopContainerResponse> stopContainer(
    $0.StopContainerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopContainer, request, options: options);
  }

  /// RestartContainer redémarre un conteneur.
  $grpc.ResponseFuture<$0.RestartContainerResponse> restartContainer(
    $0.RestartContainerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restartContainer, request, options: options);
  }

  /// StreamContainerLogs ouvre un flux gRPC pour lire les logs d'un conteneur en temps réel (streaming).
  $grpc.ResponseStream<$0.StreamContainerLogsResponse> streamContainerLogs(
    $0.StreamContainerLogsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamContainerLogs, $async.Stream.fromIterable([request]),
        options: options);
  }

  // method descriptors

  static final _$listContainers =
      $grpc.ClientMethod<$0.ListContainersRequest, $0.ListContainersResponse>(
          '/forgenet.v1.ContainerService/ListContainers',
          ($0.ListContainersRequest value) => value.writeToBuffer(),
          $0.ListContainersResponse.fromBuffer);
  static final _$getContainer =
      $grpc.ClientMethod<$0.GetContainerRequest, $0.GetContainerResponse>(
          '/forgenet.v1.ContainerService/GetContainer',
          ($0.GetContainerRequest value) => value.writeToBuffer(),
          $0.GetContainerResponse.fromBuffer);
  static final _$startContainer =
      $grpc.ClientMethod<$0.StartContainerRequest, $0.StartContainerResponse>(
          '/forgenet.v1.ContainerService/StartContainer',
          ($0.StartContainerRequest value) => value.writeToBuffer(),
          $0.StartContainerResponse.fromBuffer);
  static final _$stopContainer =
      $grpc.ClientMethod<$0.StopContainerRequest, $0.StopContainerResponse>(
          '/forgenet.v1.ContainerService/StopContainer',
          ($0.StopContainerRequest value) => value.writeToBuffer(),
          $0.StopContainerResponse.fromBuffer);
  static final _$restartContainer = $grpc.ClientMethod<
          $0.RestartContainerRequest, $0.RestartContainerResponse>(
      '/forgenet.v1.ContainerService/RestartContainer',
      ($0.RestartContainerRequest value) => value.writeToBuffer(),
      $0.RestartContainerResponse.fromBuffer);
  static final _$streamContainerLogs = $grpc.ClientMethod<
          $0.StreamContainerLogsRequest, $0.StreamContainerLogsResponse>(
      '/forgenet.v1.ContainerService/StreamContainerLogs',
      ($0.StreamContainerLogsRequest value) => value.writeToBuffer(),
      $0.StreamContainerLogsResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.ContainerService')
abstract class ContainerServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.ContainerService';

  ContainerServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ListContainersRequest,
            $0.ListContainersResponse>(
        'ListContainers',
        listContainers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListContainersRequest.fromBuffer(value),
        ($0.ListContainersResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetContainerRequest, $0.GetContainerResponse>(
            'GetContainer',
            getContainer_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetContainerRequest.fromBuffer(value),
            ($0.GetContainerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartContainerRequest,
            $0.StartContainerResponse>(
        'StartContainer',
        startContainer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartContainerRequest.fromBuffer(value),
        ($0.StartContainerResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StopContainerRequest, $0.StopContainerResponse>(
            'StopContainer',
            stopContainer_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StopContainerRequest.fromBuffer(value),
            ($0.StopContainerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RestartContainerRequest,
            $0.RestartContainerResponse>(
        'RestartContainer',
        restartContainer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RestartContainerRequest.fromBuffer(value),
        ($0.RestartContainerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamContainerLogsRequest,
            $0.StreamContainerLogsResponse>(
        'StreamContainerLogs',
        streamContainerLogs_Pre,
        false,
        true,
        ($core.List<$core.int> value) =>
            $0.StreamContainerLogsRequest.fromBuffer(value),
        ($0.StreamContainerLogsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListContainersResponse> listContainers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListContainersRequest> $request) async {
    return listContainers($call, await $request);
  }

  $async.Future<$0.ListContainersResponse> listContainers(
      $grpc.ServiceCall call, $0.ListContainersRequest request);

  $async.Future<$0.GetContainerResponse> getContainer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetContainerRequest> $request) async {
    return getContainer($call, await $request);
  }

  $async.Future<$0.GetContainerResponse> getContainer(
      $grpc.ServiceCall call, $0.GetContainerRequest request);

  $async.Future<$0.StartContainerResponse> startContainer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartContainerRequest> $request) async {
    return startContainer($call, await $request);
  }

  $async.Future<$0.StartContainerResponse> startContainer(
      $grpc.ServiceCall call, $0.StartContainerRequest request);

  $async.Future<$0.StopContainerResponse> stopContainer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopContainerRequest> $request) async {
    return stopContainer($call, await $request);
  }

  $async.Future<$0.StopContainerResponse> stopContainer(
      $grpc.ServiceCall call, $0.StopContainerRequest request);

  $async.Future<$0.RestartContainerResponse> restartContainer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RestartContainerRequest> $request) async {
    return restartContainer($call, await $request);
  }

  $async.Future<$0.RestartContainerResponse> restartContainer(
      $grpc.ServiceCall call, $0.RestartContainerRequest request);

  $async.Stream<$0.StreamContainerLogsResponse> streamContainerLogs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StreamContainerLogsRequest> $request) async* {
    yield* streamContainerLogs($call, await $request);
  }

  $async.Stream<$0.StreamContainerLogsResponse> streamContainerLogs(
      $grpc.ServiceCall call, $0.StreamContainerLogsRequest request);
}
