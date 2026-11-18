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

  /// ChangePassword change le mot de passe du compte connecté, puis ferme ses
  /// autres sessions et jetons. Un mot de passe actuel incorrect ou un nouveau
  /// mot de passe refusé répond INVALID_ARGUMENT (jamais UNAUTHENTICATED, qui
  /// fermerait la session du demandeur) ; des échecs répétés répondent
  /// RESOURCE_EXHAUSTED. Le changement est inscrit au journal d'audit.
  $grpc.ResponseFuture<$0.ChangePasswordResponse> changePassword(
    $0.ChangePasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changePassword, request, options: options);
  }

  /// InspectAccessLink décrit un lien à usage unique encore valable, sans
  /// l'utiliser : la page dit pour quel compte on choisit un mot de passe.
  /// Appelable sans jeton. Un lien inconnu, expiré, déjà utilisé ou remplacé
  /// répond NOT_FOUND, sans dire lequel.
  $grpc.ResponseFuture<$0.InspectAccessLinkResponse> inspectAccessLink(
    $0.InspectAccessLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$inspectAccessLink, request, options: options);
  }

  /// RedeemAccessLink utilise le lien : il fixe le mot de passe choisi, ferme
  /// les sessions et jetons du compte, puis ouvre une session. Appelable sans
  /// jeton. Un lien qui ne vaut plus répond NOT_FOUND ; un mot de passe refusé
  /// répond INVALID_ARGUMENT et laisse le lien valable. L'usage est inscrit au
  /// journal d'audit.
  $grpc.ResponseFuture<$0.RedeemAccessLinkResponse> redeemAccessLink(
    $0.RedeemAccessLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$redeemAccessLink, request, options: options);
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
  static final _$changePassword =
      $grpc.ClientMethod<$0.ChangePasswordRequest, $0.ChangePasswordResponse>(
          '/forgenet.v1.AuthService/ChangePassword',
          ($0.ChangePasswordRequest value) => value.writeToBuffer(),
          $0.ChangePasswordResponse.fromBuffer);
  static final _$inspectAccessLink = $grpc.ClientMethod<
          $0.InspectAccessLinkRequest, $0.InspectAccessLinkResponse>(
      '/forgenet.v1.AuthService/InspectAccessLink',
      ($0.InspectAccessLinkRequest value) => value.writeToBuffer(),
      $0.InspectAccessLinkResponse.fromBuffer);
  static final _$redeemAccessLink = $grpc.ClientMethod<
          $0.RedeemAccessLinkRequest, $0.RedeemAccessLinkResponse>(
      '/forgenet.v1.AuthService/RedeemAccessLink',
      ($0.RedeemAccessLinkRequest value) => value.writeToBuffer(),
      $0.RedeemAccessLinkResponse.fromBuffer);
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
    $addMethod($grpc.ServiceMethod<$0.ChangePasswordRequest,
            $0.ChangePasswordResponse>(
        'ChangePassword',
        changePassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangePasswordRequest.fromBuffer(value),
        ($0.ChangePasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.InspectAccessLinkRequest,
            $0.InspectAccessLinkResponse>(
        'InspectAccessLink',
        inspectAccessLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.InspectAccessLinkRequest.fromBuffer(value),
        ($0.InspectAccessLinkResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RedeemAccessLinkRequest,
            $0.RedeemAccessLinkResponse>(
        'RedeemAccessLink',
        redeemAccessLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RedeemAccessLinkRequest.fromBuffer(value),
        ($0.RedeemAccessLinkResponse value) => value.writeToBuffer()));
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

  $async.Future<$0.ChangePasswordResponse> changePassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ChangePasswordRequest> $request) async {
    return changePassword($call, await $request);
  }

  $async.Future<$0.ChangePasswordResponse> changePassword(
      $grpc.ServiceCall call, $0.ChangePasswordRequest request);

  $async.Future<$0.InspectAccessLinkResponse> inspectAccessLink_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.InspectAccessLinkRequest> $request) async {
    return inspectAccessLink($call, await $request);
  }

  $async.Future<$0.InspectAccessLinkResponse> inspectAccessLink(
      $grpc.ServiceCall call, $0.InspectAccessLinkRequest request);

  $async.Future<$0.RedeemAccessLinkResponse> redeemAccessLink_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RedeemAccessLinkRequest> $request) async {
    return redeemAccessLink($call, await $request);
  }

  $async.Future<$0.RedeemAccessLinkResponse> redeemAccessLink(
      $grpc.ServiceCall call, $0.RedeemAccessLinkRequest request);
}
