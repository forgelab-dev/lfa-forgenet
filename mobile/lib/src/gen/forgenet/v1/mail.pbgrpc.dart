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

/// MailService montre l'envoi d'e-mails et l'essaie. Réservé aux
/// administrateurs. Les réglages se changent sur le serveur, par install.sh
/// (ADR-0014).
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
}
