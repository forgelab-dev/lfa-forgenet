// This is a generated file - do not edit.
//
// Generated from forgenet/v1/team.proto.

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

import 'team.pb.dart' as $0;

export 'team.pb.dart';

/// TeamService gère les membres du serveur et leurs invitations (ADR-0014).
/// Réservé aux administrateurs.
@$pb.GrpcServiceName('forgenet.v1.TeamService')
class TeamServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TeamServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ListMembersResponse> listMembers(
    $0.ListMembersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMembers, request, options: options);
  }

  /// InviteMember crée le compte sans mot de passe utilisable et fait un lien
  /// d'invitation valable 7 jours. Nom ou adresse déjà pris : ALREADY_EXISTS.
  /// Champ invalide : INVALID_ARGUMENT.
  $grpc.ResponseFuture<$0.InviteMemberResponse> inviteMember(
    $0.InviteMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$inviteMember, request, options: options);
  }

  /// RenewInvitation refait le lien d'un compte qui n'a pas encore choisi son
  /// mot de passe ; le nouveau lien annule l'ancien. Compte déjà actif :
  /// FAILED_PRECONDITION. Compte inconnu : NOT_FOUND.
  $grpc.ResponseFuture<$0.RenewInvitationResponse> renewInvitation(
    $0.RenewInvitationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$renewInvitation, request, options: options);
  }

  // method descriptors

  static final _$listMembers =
      $grpc.ClientMethod<$0.ListMembersRequest, $0.ListMembersResponse>(
          '/forgenet.v1.TeamService/ListMembers',
          ($0.ListMembersRequest value) => value.writeToBuffer(),
          $0.ListMembersResponse.fromBuffer);
  static final _$inviteMember =
      $grpc.ClientMethod<$0.InviteMemberRequest, $0.InviteMemberResponse>(
          '/forgenet.v1.TeamService/InviteMember',
          ($0.InviteMemberRequest value) => value.writeToBuffer(),
          $0.InviteMemberResponse.fromBuffer);
  static final _$renewInvitation =
      $grpc.ClientMethod<$0.RenewInvitationRequest, $0.RenewInvitationResponse>(
          '/forgenet.v1.TeamService/RenewInvitation',
          ($0.RenewInvitationRequest value) => value.writeToBuffer(),
          $0.RenewInvitationResponse.fromBuffer);
}

@$pb.GrpcServiceName('forgenet.v1.TeamService')
abstract class TeamServiceBase extends $grpc.Service {
  $core.String get $name => 'forgenet.v1.TeamService';

  TeamServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListMembersRequest, $0.ListMembersResponse>(
            'ListMembers',
            listMembers_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListMembersRequest.fromBuffer(value),
            ($0.ListMembersResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.InviteMemberRequest, $0.InviteMemberResponse>(
            'InviteMember',
            inviteMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.InviteMemberRequest.fromBuffer(value),
            ($0.InviteMemberResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RenewInvitationRequest,
            $0.RenewInvitationResponse>(
        'RenewInvitation',
        renewInvitation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RenewInvitationRequest.fromBuffer(value),
        ($0.RenewInvitationResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListMembersResponse> listMembers_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListMembersRequest> $request) async {
    return listMembers($call, await $request);
  }

  $async.Future<$0.ListMembersResponse> listMembers(
      $grpc.ServiceCall call, $0.ListMembersRequest request);

  $async.Future<$0.InviteMemberResponse> inviteMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.InviteMemberRequest> $request) async {
    return inviteMember($call, await $request);
  }

  $async.Future<$0.InviteMemberResponse> inviteMember(
      $grpc.ServiceCall call, $0.InviteMemberRequest request);

  $async.Future<$0.RenewInvitationResponse> renewInvitation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RenewInvitationRequest> $request) async {
    return renewInvitation($call, await $request);
  }

  $async.Future<$0.RenewInvitationResponse> renewInvitation(
      $grpc.ServiceCall call, $0.RenewInvitationRequest request);
}
