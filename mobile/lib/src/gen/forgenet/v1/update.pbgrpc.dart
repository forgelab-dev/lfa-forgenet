// This is a generated file - do not edit.
//
// Generated from forgenet/v1/update.proto.

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

import 'update.pb.dart' as $0;

export 'update.pb.dart';

/// UpdateService annonce les nouvelles versions et met l'installation à jour
/// (ADR-0011). Réservé aux administrateurs.
@$pb.GrpcServiceName('forgenet.v1.UpdateService')
class UpdateServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  UpdateServiceClient(super.channel, {super.options, super.interceptors});

  /// GetUpdateStatus donne la version en service, la version disponible et la
  /// dernière mise à jour lancée.
  $grpc.ResponseFuture<$0.GetUpdateStatusResponse> getUpdateStatus(
    $0.GetUpdateStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUpdateStatus, request, options: options);
  }

  /// CheckForUpdates interroge les releases tout de suite, sans attendre la
  /// vérification quotidienne. Répété trop vite, il renvoie le dernier résultat.
  $grpc.ResponseFuture<$0.CheckForUpdatesResponse> checkForUpdates(
    $0.CheckForUpdatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkForUpdates, request, options: options);
  }

  /// StartUpdate lance la mise à jour vers la version disponible, la seule que
  /// le serveur accepte d'installer : le client ne choisit ni la version ni
  /// l'image. Rien à installer répond FAILED_PRECONDITION ; une mise à jour
  /// déjà en cours, ABORTED.
  $grpc.ResponseFuture<$0.StartUpdateResponse> startUpdate(
    $0.StartUpdateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startUpdate, request, options: options);
  }

  // method descriptors

  static final _$getUpdateStatus =
      $grpc.ClientMethod<$0.GetUpdateStatusRequest, $0.GetUpdateStatusResponse>(
          '/forgenet.v1.UpdateService/GetUpdateStatus',
          ($0.GetUpdateStatusRequest value) => value.writeToBuffer(),
          $0.GetUpdateStatusResponse.fromBuffer);
  static final _$checkForUpdates =
      $grpc.ClientMethod<$0.CheckForUpdatesRequest, $0.CheckForUpdatesResponse>(
          '/forgenet.v1.UpdateService/CheckForUpdates',
          ($0.CheckForUpdatesRequest value) => value.writeToBuffer(),
          $0.CheckForUpdatesResponse.fromBuffer);
  static final _$startUpdate =
      $grpc.ClientMethod<$0.StartUpdateRequest, $0.StartUpdateResponse>(
          '/forgenet.v1.UpdateService/StartUpdate',
          ($0.StartUpdateRequest value) => value.writeToBuffer(),
          $0.StartUpdateResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.UpdateService')
abstract class UpdateServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.UpdateService';

  UpdateServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetUpdateStatusRequest,
            $0.GetUpdateStatusResponse>(
        'GetUpdateStatus',
        getUpdateStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUpdateStatusRequest.fromBuffer(value),
        ($0.GetUpdateStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckForUpdatesRequest,
            $0.CheckForUpdatesResponse>(
        'CheckForUpdates',
        checkForUpdates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckForUpdatesRequest.fromBuffer(value),
        ($0.CheckForUpdatesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartUpdateRequest, $0.StartUpdateResponse>(
            'StartUpdate',
            startUpdate_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartUpdateRequest.fromBuffer(value),
            ($0.StartUpdateResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.GetUpdateStatusResponse> getUpdateStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUpdateStatusRequest> $request) async {
    return getUpdateStatus($call, await $request);
  }

  $async.Future<$0.GetUpdateStatusResponse> getUpdateStatus(
      $grpc.ServiceCall call, $0.GetUpdateStatusRequest request);

  $async.Future<$0.CheckForUpdatesResponse> checkForUpdates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckForUpdatesRequest> $request) async {
    return checkForUpdates($call, await $request);
  }

  $async.Future<$0.CheckForUpdatesResponse> checkForUpdates(
      $grpc.ServiceCall call, $0.CheckForUpdatesRequest request);

  $async.Future<$0.StartUpdateResponse> startUpdate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartUpdateRequest> $request) async {
    return startUpdate($call, await $request);
  }

  $async.Future<$0.StartUpdateResponse> startUpdate(
      $grpc.ServiceCall call, $0.StartUpdateRequest request);
}
