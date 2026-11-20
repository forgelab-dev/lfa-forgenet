// This is a generated file - do not edit.
//
// Generated from forgenet/v1/mail.proto.

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

import 'mail.pb.dart' as $0;

export 'mail.pb.dart';

/// MailService montre l'envoi d'e-mails, l'essaie et le règle. Réservé aux
/// administrateurs. Les réglages du dashboard l'emportent sur ceux
/// d'install.sh tant qu'ils existent (ADR-0015).
@$pb.GrpcServiceName('forgenet.v1.MailService')
class MailServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  MailServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.GetMailStatusResponse> getMailStatus(
    $0.GetMailStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMailStatus, request, options: options);
  }

  /// SendTestMail envoie un e-mail d'essai à l'administrateur connecté et
  /// attend la réponse du relais. Sans relais : FAILED_PRECONDITION. Refus du
  /// relais : FAILED_PRECONDITION avec sa réponse. Envois répétés :
  /// RESOURCE_EXHAUSTED.
  $grpc.ResponseFuture<$0.SendTestMailResponse> sendTestMail(
    $0.SendTestMailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendTestMail, request, options: options);
  }

  /// SaveMailSettings essaie les nouveaux réglages par un e-mail à
  /// l'administrateur connecté, puis les enregistre et les applique sans
  /// redémarrage. Réglages invalides : INVALID_ARGUMENT. Refus du relais :
  /// FAILED_PRECONDITION avec sa réponse, et les anciens réglages restent.
  /// Serveur sans clé : FAILED_PRECONDITION. Essais répétés :
  /// RESOURCE_EXHAUSTED.
  $grpc.ResponseFuture<$0.SaveMailSettingsResponse> saveMailSettings(
    $0.SaveMailSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$saveMailSettings, request, options: options);
  }

  /// DisableMail coupe l'envoi ; les valeurs restent pour le rouvrir.
  /// Serveur sans clé : FAILED_PRECONDITION.
  $grpc.ResponseFuture<$0.DisableMailResponse> disableMail(
    $0.DisableMailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableMail, request, options: options);
  }

  /// ResetMailSettings efface les réglages du dashboard : ceux d'install.sh
  /// (.env) s'appliquent de nouveau.
  $grpc.ResponseFuture<$0.ResetMailSettingsResponse> resetMailSettings(
    $0.ResetMailSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resetMailSettings, request, options: options);
  }

  // method descriptors

  static final _$getMailStatus =
      $grpc.ClientMethod<$0.GetMailStatusRequest, $0.GetMailStatusResponse>(
          '/forgenet.v1.MailService/GetMailStatus',
          ($0.GetMailStatusRequest value) => value.writeToBuffer(),
          $0.GetMailStatusResponse.fromBuffer);
  static final _$sendTestMail =
      $grpc.ClientMethod<$0.SendTestMailRequest, $0.SendTestMailResponse>(
          '/forgenet.v1.MailService/SendTestMail',
          ($0.SendTestMailRequest value) => value.writeToBuffer(),
          $0.SendTestMailResponse.fromBuffer);
  static final _$saveMailSettings = $grpc.ClientMethod<
          $0.SaveMailSettingsRequest, $0.SaveMailSettingsResponse>(
      '/forgenet.v1.MailService/SaveMailSettings',
      ($0.SaveMailSettingsRequest value) => value.writeToBuffer(),
      $0.SaveMailSettingsResponse.fromBuffer);
  static final _$disableMail =
      $grpc.ClientMethod<$0.DisableMailRequest, $0.DisableMailResponse>(
          '/forgenet.v1.MailService/DisableMail',
          ($0.DisableMailRequest value) => value.writeToBuffer(),
          $0.DisableMailResponse.fromBuffer);
  static final _$resetMailSettings = $grpc.ClientMethod<
          $0.ResetMailSettingsRequest, $0.ResetMailSettingsResponse>(
      '/forgenet.v1.MailService/ResetMailSettings',
      ($0.ResetMailSettingsRequest value) => value.writeToBuffer(),
      $0.ResetMailSettingsResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.MailService')
abstract class MailServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.MailService';

  MailServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.GetMailStatusRequest, $0.GetMailStatusResponse>(
            'GetMailStatus',
            getMailStatus_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetMailStatusRequest.fromBuffer(value),
            ($0.GetMailStatusResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SendTestMailRequest, $0.SendTestMailResponse>(
            'SendTestMail',
            sendTestMail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SendTestMailRequest.fromBuffer(value),
            ($0.SendTestMailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SaveMailSettingsRequest,
            $0.SaveMailSettingsResponse>(
        'SaveMailSettings',
        saveMailSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SaveMailSettingsRequest.fromBuffer(value),
        ($0.SaveMailSettingsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DisableMailRequest, $0.DisableMailResponse>(
            'DisableMail',
            disableMail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DisableMailRequest.fromBuffer(value),
            ($0.DisableMailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResetMailSettingsRequest,
            $0.ResetMailSettingsResponse>(
        'ResetMailSettings',
        resetMailSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResetMailSettingsRequest.fromBuffer(value),
        ($0.ResetMailSettingsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.GetMailStatusResponse> getMailStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMailStatusRequest> $request) async {
    return getMailStatus($call, await $request);
  }

  $async.Future<$0.GetMailStatusResponse> getMailStatus(
      $grpc.ServiceCall call, $0.GetMailStatusRequest request);

  $async.Future<$0.SendTestMailResponse> sendTestMail_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendTestMailRequest> $request) async {
    return sendTestMail($call, await $request);
  }

  $async.Future<$0.SendTestMailResponse> sendTestMail(
      $grpc.ServiceCall call, $0.SendTestMailRequest request);

  $async.Future<$0.SaveMailSettingsResponse> saveMailSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SaveMailSettingsRequest> $request) async {
    return saveMailSettings($call, await $request);
  }

  $async.Future<$0.SaveMailSettingsResponse> saveMailSettings(
      $grpc.ServiceCall call, $0.SaveMailSettingsRequest request);

  $async.Future<$0.DisableMailResponse> disableMail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisableMailRequest> $request) async {
    return disableMail($call, await $request);
  }

  $async.Future<$0.DisableMailResponse> disableMail(
      $grpc.ServiceCall call, $0.DisableMailRequest request);

  $async.Future<$0.ResetMailSettingsResponse> resetMailSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResetMailSettingsRequest> $request) async {
    return resetMailSettings($call, await $request);
  }

  $async.Future<$0.ResetMailSettingsResponse> resetMailSettings(
      $grpc.ServiceCall call, $0.ResetMailSettingsRequest request);
}
