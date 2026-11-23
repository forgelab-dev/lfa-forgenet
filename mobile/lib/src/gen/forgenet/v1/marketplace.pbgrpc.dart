// This is a generated file - do not edit.
//
// Generated from forgenet/v1/marketplace.proto.

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

import 'marketplace.pb.dart' as $0;

export 'marketplace.pb.dart';

/// MarketplaceService présente le catalogue des templates embarqué dans le
/// serveur (ADR-0018) : chaque version de ForgeNet propose les templates
/// vérifiés par sa CI, sans aller les chercher sur le réseau. Réservé aux
/// administrateurs, comme les déploiements.
@$pb.GrpcServiceName('forgenet.v1.MarketplaceService')
class MarketplaceServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  MarketplaceServiceClient(super.channel, {super.options, super.interceptors});

  /// ListTemplates donne tout le catalogue.
  $grpc.ResponseFuture<$0.ListTemplatesResponse> listTemplates(
    $0.ListTemplatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTemplates, request, options: options);
  }

  /// GetTemplate donne un template et son README. Un identifiant mal formé
  /// répond INVALID_ARGUMENT ; un template absent du catalogue, NOT_FOUND.
  $grpc.ResponseFuture<$0.GetTemplateResponse> getTemplate(
    $0.GetTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTemplate, request, options: options);
  }

  // method descriptors

  static final _$listTemplates =
      $grpc.ClientMethod<$0.ListTemplatesRequest, $0.ListTemplatesResponse>(
          '/forgenet.v1.MarketplaceService/ListTemplates',
          ($0.ListTemplatesRequest value) => value.writeToBuffer(),
          $0.ListTemplatesResponse.fromBuffer);
  static final _$getTemplate =
      $grpc.ClientMethod<$0.GetTemplateRequest, $0.GetTemplateResponse>(
          '/forgenet.v1.MarketplaceService/GetTemplate',
          ($0.GetTemplateRequest value) => value.writeToBuffer(),
          $0.GetTemplateResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.MarketplaceService')
abstract class MarketplaceServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.MarketplaceService';

  MarketplaceServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListTemplatesRequest, $0.ListTemplatesResponse>(
            'ListTemplates',
            listTemplates_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListTemplatesRequest.fromBuffer(value),
            ($0.ListTemplatesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetTemplateRequest, $0.GetTemplateResponse>(
            'GetTemplate',
            getTemplate_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetTemplateRequest.fromBuffer(value),
            ($0.GetTemplateResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListTemplatesResponse> listTemplates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTemplatesRequest> $request) async {
    return listTemplates($call, await $request);
  }

  $async.Future<$0.ListTemplatesResponse> listTemplates(
      $grpc.ServiceCall call, $0.ListTemplatesRequest request);

  $async.Future<$0.GetTemplateResponse> getTemplate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetTemplateRequest> $request) async {
    return getTemplate($call, await $request);
  }

  $async.Future<$0.GetTemplateResponse> getTemplate(
      $grpc.ServiceCall call, $0.GetTemplateRequest request);
}
