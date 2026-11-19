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

  /// SuspendMember empêche le compte de se connecter et ferme tout de suite ses
  /// sessions, jetons d'API et liens à usage unique. Réversible : le compte et
  /// son historique restent. Son propre compte, ou le dernier administrateur
  /// actif : FAILED_PRECONDITION. Compte inconnu : NOT_FOUND.
  $grpc.ResponseFuture<$0.SuspendMemberResponse> suspendMember(
    $0.SuspendMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$suspendMember, request, options: options);
  }

  /// ResumeMember réactive un compte suspendu ; ses sessions fermées le
  /// restent. Compte inconnu : NOT_FOUND.
  $grpc.ResponseFuture<$0.ResumeMemberResponse> resumeMember(
    $0.ResumeMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resumeMember, request, options: options);
  }

  /// CloseMemberSessions ferme les sessions et jetons d'API du compte sans le
  /// suspendre (appareil perdu). Sur son propre compte, la session qui fait
  /// l'appel reste ouverte. Compte inconnu : NOT_FOUND.
  $grpc.ResponseFuture<$0.CloseMemberSessionsResponse> closeMemberSessions(
    $0.CloseMemberSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$closeMemberSessions, request, options: options);
  }

  /// DeleteMember supprime le compte, définitivement ; une invitation en
  /// attente est ainsi annulée. Le journal d'audit garde ses entrées. Nom de
  /// confirmation différent : INVALID_ARGUMENT. Son propre compte, ou le
  /// dernier administrateur actif : FAILED_PRECONDITION. Compte inconnu :
  /// NOT_FOUND.
  $grpc.ResponseFuture<$0.DeleteMemberResponse> deleteMember(
    $0.DeleteMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMember, request, options: options);
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
  static final _$suspendMember =
      $grpc.ClientMethod<$0.SuspendMemberRequest, $0.SuspendMemberResponse>(
          '/forgenet.v1.TeamService/SuspendMember',
          ($0.SuspendMemberRequest value) => value.writeToBuffer(),
          $0.SuspendMemberResponse.fromBuffer);
  static final _$resumeMember =
      $grpc.ClientMethod<$0.ResumeMemberRequest, $0.ResumeMemberResponse>(
          '/forgenet.v1.TeamService/ResumeMember',
          ($0.ResumeMemberRequest value) => value.writeToBuffer(),
          $0.ResumeMemberResponse.fromBuffer);
  static final _$closeMemberSessions = $grpc.ClientMethod<
          $0.CloseMemberSessionsRequest, $0.CloseMemberSessionsResponse>(
      '/forgenet.v1.TeamService/CloseMemberSessions',
      ($0.CloseMemberSessionsRequest value) => value.writeToBuffer(),
      $0.CloseMemberSessionsResponse.fromBuffer);
  static final _$deleteMember =
      $grpc.ClientMethod<$0.DeleteMemberRequest, $0.DeleteMemberResponse>(
          '/forgenet.v1.TeamService/DeleteMember',
          ($0.DeleteMemberRequest value) => value.writeToBuffer(),
          $0.DeleteMemberResponse.fromBuffer);
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
    $addMethod(
        $grpc.ServiceMethod<$0.SuspendMemberRequest, $0.SuspendMemberResponse>(
            'SuspendMember',
            suspendMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SuspendMemberRequest.fromBuffer(value),
            ($0.SuspendMemberResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ResumeMemberRequest, $0.ResumeMemberResponse>(
            'ResumeMember',
            resumeMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ResumeMemberRequest.fromBuffer(value),
            ($0.ResumeMemberResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CloseMemberSessionsRequest,
            $0.CloseMemberSessionsResponse>(
        'CloseMemberSessions',
        closeMemberSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CloseMemberSessionsRequest.fromBuffer(value),
        ($0.CloseMemberSessionsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteMemberRequest, $0.DeleteMemberResponse>(
            'DeleteMember',
            deleteMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteMemberRequest.fromBuffer(value),
            ($0.DeleteMemberResponse value) => value.writeToBuffer()));
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

  $async.Future<$0.SuspendMemberResponse> suspendMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SuspendMemberRequest> $request) async {
    return suspendMember($call, await $request);
  }

  $async.Future<$0.SuspendMemberResponse> suspendMember(
      $grpc.ServiceCall call, $0.SuspendMemberRequest request);

  $async.Future<$0.ResumeMemberResponse> resumeMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResumeMemberRequest> $request) async {
    return resumeMember($call, await $request);
  }

  $async.Future<$0.ResumeMemberResponse> resumeMember(
      $grpc.ServiceCall call, $0.ResumeMemberRequest request);

  $async.Future<$0.CloseMemberSessionsResponse> closeMemberSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CloseMemberSessionsRequest> $request) async {
    return closeMemberSessions($call, await $request);
  }

  $async.Future<$0.CloseMemberSessionsResponse> closeMemberSessions(
      $grpc.ServiceCall call, $0.CloseMemberSessionsRequest request);

  $async.Future<$0.DeleteMemberResponse> deleteMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteMemberRequest> $request) async {
    return deleteMember($call, await $request);
  }

  $async.Future<$0.DeleteMemberResponse> deleteMember(
      $grpc.ServiceCall call, $0.DeleteMemberRequest request);
}
