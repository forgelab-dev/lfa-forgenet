// This is a generated file - do not edit.
//
// Generated from forgenet/v1/domain.proto.

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

import 'domain.pb.dart' as $0;

export 'domain.pb.dart';

/// DomainService gère les noms de domaine servis par le reverse proxy. Tant
/// que le proxy est désactivé dans la configuration du serveur, chaque méthode
/// répond FAILED_PRECONDITION.
@$pb.GrpcServiceName('forgenet.v1.DomainService')
class DomainServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DomainServiceClient(super.channel, {super.options, super.interceptors});

  /// ListDomains liste les domaines, d'un projet ou de tous.
  $grpc.ResponseFuture<$0.ListDomainsResponse> listDomains(
    $0.ListDomainsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDomains, request, options: options);
  }

  /// AddDomain relie un domaine au service déployé d'un projet. Un nom déjà
  /// pris répond ALREADY_EXISTS ; un service absent de la stack, NOT_FOUND.
  $grpc.ResponseFuture<$0.AddDomainResponse> addDomain(
    $0.AddDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addDomain, request, options: options);
  }

  /// RemoveDomain retire un domaine. Caddy cesse de le servir avant que le
  /// domaine ne soit effacé.
  $grpc.ResponseFuture<$0.RemoveDomainResponse> removeDomain(
    $0.RemoveDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeDomain, request, options: options);
  }

  // method descriptors

  static final _$listDomains =
      $grpc.ClientMethod<$0.ListDomainsRequest, $0.ListDomainsResponse>(
          '/forgenet.v1.DomainService/ListDomains',
          ($0.ListDomainsRequest value) => value.writeToBuffer(),
          $0.ListDomainsResponse.fromBuffer);
  static final _$addDomain =
      $grpc.ClientMethod<$0.AddDomainRequest, $0.AddDomainResponse>(
          '/forgenet.v1.DomainService/AddDomain',
          ($0.AddDomainRequest value) => value.writeToBuffer(),
          $0.AddDomainResponse.fromBuffer);
  static final _$removeDomain =
      $grpc.ClientMethod<$0.RemoveDomainRequest, $0.RemoveDomainResponse>(
          '/forgenet.v1.DomainService/RemoveDomain',
          ($0.RemoveDomainRequest value) => value.writeToBuffer(),
          $0.RemoveDomainResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.DomainService')
abstract class DomainServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.DomainService';

  DomainServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListDomainsRequest, $0.ListDomainsResponse>(
            'ListDomains',
            listDomains_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDomainsRequest.fromBuffer(value),
            ($0.ListDomainsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddDomainRequest, $0.AddDomainResponse>(
        'AddDomain',
        addDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AddDomainRequest.fromBuffer(value),
        ($0.AddDomainResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RemoveDomainRequest, $0.RemoveDomainResponse>(
            'RemoveDomain',
            removeDomain_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RemoveDomainRequest.fromBuffer(value),
            ($0.RemoveDomainResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListDomainsResponse> listDomains_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListDomainsRequest> $request) async {
    return listDomains($call, await $request);
  }

  $async.Future<$0.ListDomainsResponse> listDomains(
      $grpc.ServiceCall call, $0.ListDomainsRequest request);

  $async.Future<$0.AddDomainResponse> addDomain_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddDomainRequest> $request) async {
    return addDomain($call, await $request);
  }

  $async.Future<$0.AddDomainResponse> addDomain(
      $grpc.ServiceCall call, $0.AddDomainRequest request);

  $async.Future<$0.RemoveDomainResponse> removeDomain_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveDomainRequest> $request) async {
    return removeDomain($call, await $request);
  }

  $async.Future<$0.RemoveDomainResponse> removeDomain(
      $grpc.ServiceCall call, $0.RemoveDomainRequest request);
}
