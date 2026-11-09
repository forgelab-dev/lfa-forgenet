// This is a generated file - do not edit.
//
// Generated from forgenet/v1/auth.proto.

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

import 'auth.pb.dart' as $0;

export 'auth.pb.dart';

/// AuthService ouvre et ferme les sessions du dashboard et du CLI.
@$pb.GrpcServiceName('forgenet.v1.AuthService')
class AuthServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AuthServiceClient(super.channel, {super.options, super.interceptors});

  /// Login ouvre une session. Seule méthode appelable sans jeton ; un échec
  /// répond UNAUTHENTICATED sans dire si le compte existe, et des échecs
  /// répétés répondent RESOURCE_EXHAUSTED.
  $grpc.ResponseFuture<$0.LoginResponse> login(
    $0.LoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$login, request, options: options);
  }

  /// Logout révoque la session du jeton présenté.
  $grpc.ResponseFuture<$0.LogoutResponse> logout(
    $0.LogoutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$logout, request, options: options);
  }

  /// GetCurrentUser renvoie le compte du jeton présenté, avec son rôle actuel.
  $grpc.ResponseFuture<$0.GetCurrentUserResponse> getCurrentUser(
    $0.GetCurrentUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCurrentUser, request, options: options);
  }

  // method descriptors

  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/forgenet.v1.AuthService/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$logout =
      $grpc.ClientMethod<$0.LogoutRequest, $0.LogoutResponse>(
          '/forgenet.v1.AuthService/Logout',
          ($0.LogoutRequest value) => value.writeToBuffer(),
          $0.LogoutResponse.fromBuffer);
  static final _$getCurrentUser =
      $grpc.ClientMethod<$0.GetCurrentUserRequest, $0.GetCurrentUserResponse>(
          '/forgenet.v1.AuthService/GetCurrentUser',
          ($0.GetCurrentUserRequest value) => value.writeToBuffer(),
          $0.GetCurrentUserResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.AuthService')
abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.AuthService';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogoutRequest, $0.LogoutResponse>(
        'Logout',
        logout_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogoutRequest.fromBuffer(value),
        ($0.LogoutResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCurrentUserRequest,
            $0.GetCurrentUserResponse>(
        'GetCurrentUser',
        getCurrentUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCurrentUserRequest.fromBuffer(value),
        ($0.GetCurrentUserResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.LoginResponse> login_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login(
      $grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$0.LogoutResponse> logout_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LogoutRequest> $request) async {
    return logout($call, await $request);
  }

  $async.Future<$0.LogoutResponse> logout(
      $grpc.ServiceCall call, $0.LogoutRequest request);

  $async.Future<$0.GetCurrentUserResponse> getCurrentUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCurrentUserRequest> $request) async {
    return getCurrentUser($call, await $request);
  }

  $async.Future<$0.GetCurrentUserResponse> getCurrentUser(
      $grpc.ServiceCall call, $0.GetCurrentUserRequest request);
}
