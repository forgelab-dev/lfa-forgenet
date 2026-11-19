// This is a generated file - do not edit.
//
// Generated from forgenet/v1/team.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

import 'auth.pbenum.dart' as $2;
import 'team.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'team.pbenum.dart';

/// Member est un compte du serveur, vu par un administrateur. Le mot de passe
/// et ses empreintes ne quittent jamais le serveur.
class Member extends $pb.GeneratedMessage {
  factory Member({
    $core.String? id,
    $core.String? username,
    $core.String? email,
    $2.UserRole? role,
    MemberStatus? status,
    $1.Timestamp? invitationExpiresAt,
    $1.Timestamp? createdAt,
    $1.Timestamp? suspendedAt,
  }) {
    final result = Member._();
    if (id != null) result.id = id;
    if (username != null) result.username = username;
    if (email != null) result.email = email;
    if (role != null) result.role = role;
    if (status != null) result.status = status;
    if (invitationExpiresAt != null)
      result.invitationExpiresAt = invitationExpiresAt;
    if (createdAt != null) result.createdAt = createdAt;
    if (suspendedAt != null) result.suspendedAt = suspendedAt;
    return result;
  }

  Member._();

  factory Member.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Member()..mergeFromBuffer(data, registry);
  factory Member.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Member()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Member',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: Member.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'username')
    ..aOS(3, _omitFieldNames ? '' : 'email')
    ..aE<$2.UserRole>(4, _omitFieldNames ? '' : 'role',
        enumValues: $2.UserRole.values)
    ..aE<MemberStatus>(5, _omitFieldNames ? '' : 'status',
        enumValues: MemberStatus.values)
    ..aOM<$1.Timestamp>(6, _omitFieldNames ? '' : 'invitationExpiresAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(7, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(8, _omitFieldNames ? '' : 'suspendedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Member clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Member copyWith(void Function(Member) updates) =>
      super.copyWith((message) => updates(message as Member)) as Member;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Member() / Member.new instead')
  static Member create() => Member._();
  static $pb.GeneratedMessage $_createMessage() => Member._();
  @$core.override
  Member createEmptyInstance() => Member._();
  @$core.pragma('dart2js:noInline')
  static Member getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Member>(Member.$_createMessage);
  static Member? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get username => $_getSZ(1);
  @$pb.TagNumber(2)
  set username($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUsername() => $_has(1);
  @$pb.TagNumber(2)
  void clearUsername() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get email => $_getSZ(2);
  @$pb.TagNumber(3)
  set email($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasEmail() => $_has(2);
  @$pb.TagNumber(3)
  void clearEmail() => $_clearField(3);

  @$pb.TagNumber(4)
  $2.UserRole get role => $_getN(3);
  @$pb.TagNumber(4)
  set role($2.UserRole value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasRole() => $_has(3);
  @$pb.TagNumber(4)
  void clearRole() => $_clearField(4);

  @$pb.TagNumber(5)
  MemberStatus get status => $_getN(4);
  @$pb.TagNumber(5)
  set status(MemberStatus value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasStatus() => $_has(4);
  @$pb.TagNumber(5)
  void clearStatus() => $_clearField(5);

  /// Fin de validité de l'invitation en cours ; absente hors MEMBER_STATUS_INVITED.
  @$pb.TagNumber(6)
  $1.Timestamp get invitationExpiresAt => $_getN(5);
  @$pb.TagNumber(6)
  set invitationExpiresAt($1.Timestamp value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasInvitationExpiresAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearInvitationExpiresAt() => $_clearField(6);
  @$pb.TagNumber(6)
  $1.Timestamp ensureInvitationExpiresAt() => $_ensure(5);

  @$pb.TagNumber(7)
  $1.Timestamp get createdAt => $_getN(6);
  @$pb.TagNumber(7)
  set createdAt($1.Timestamp value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasCreatedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearCreatedAt() => $_clearField(7);
  @$pb.TagNumber(7)
  $1.Timestamp ensureCreatedAt() => $_ensure(6);

  /// Date de la suspension ; absente hors MEMBER_STATUS_SUSPENDED.
  @$pb.TagNumber(8)
  $1.Timestamp get suspendedAt => $_getN(7);
  @$pb.TagNumber(8)
  set suspendedAt($1.Timestamp value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasSuspendedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearSuspendedAt() => $_clearField(8);
  @$pb.TagNumber(8)
  $1.Timestamp ensureSuspendedAt() => $_ensure(7);
}

class ListMembersRequest extends $pb.GeneratedMessage {
  factory ListMembersRequest() => ListMembersRequest._();

  ListMembersRequest._();

  factory ListMembersRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListMembersRequest()..mergeFromBuffer(data, registry);
  factory ListMembersRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListMembersRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMembersRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListMembersRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMembersRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMembersRequest copyWith(void Function(ListMembersRequest) updates) =>
      super.copyWith((message) => updates(message as ListMembersRequest))
          as ListMembersRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ListMembersRequest() / ListMembersRequest.new instead')
  static ListMembersRequest create() => ListMembersRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListMembersRequest._();
  @$core.override
  ListMembersRequest createEmptyInstance() => ListMembersRequest._();
  @$core.pragma('dart2js:noInline')
  static ListMembersRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMembersRequest>(
          ListMembersRequest.$_createMessage);
  static ListMembersRequest? _defaultInstance;
}

class ListMembersResponse extends $pb.GeneratedMessage {
  factory ListMembersResponse({
    $core.Iterable<Member>? members,
    $core.bool? mailEnabled,
  }) {
    final result = ListMembersResponse._();
    if (members != null) result.members.addAll(members);
    if (mailEnabled != null) result.mailEnabled = mailEnabled;
    return result;
  }

  ListMembersResponse._();

  factory ListMembersResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListMembersResponse()..mergeFromBuffer(data, registry);
  factory ListMembersResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListMembersResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMembersResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListMembersResponse.$_createMessage)
    ..pPM<Member>(1, _omitFieldNames ? '' : 'members',
        subBuilder: Member.$_createMessage)
    ..aOB(2, _omitFieldNames ? '' : 'mailEnabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMembersResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMembersResponse copyWith(void Function(ListMembersResponse) updates) =>
      super.copyWith((message) => updates(message as ListMembersResponse))
          as ListMembersResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use ListMembersResponse() / ListMembersResponse.new instead')
  static ListMembersResponse create() => ListMembersResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListMembersResponse._();
  @$core.override
  ListMembersResponse createEmptyInstance() => ListMembersResponse._();
  @$core.pragma('dart2js:noInline')
  static ListMembersResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMembersResponse>(
          ListMembersResponse.$_createMessage);
  static ListMembersResponse? _defaultInstance;

  /// Par ancienneté, le plus ancien en tête.
  @$pb.TagNumber(1)
  $pb.PbList<Member> get members => $_getList(0);

  /// Vrai si un relais SMTP est configuré : les invitations partent alors par
  /// e-mail ; sinon, le lien est rendu une fois à l'administrateur.
  @$pb.TagNumber(2)
  $core.bool get mailEnabled => $_getBF(1);
  @$pb.TagNumber(2)
  set mailEnabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMailEnabled() => $_has(1);
  @$pb.TagNumber(2)
  void clearMailEnabled() => $_clearField(2);
}

class InviteMemberRequest extends $pb.GeneratedMessage {
  factory InviteMemberRequest({
    $core.String? username,
    $core.String? email,
    $2.UserRole? role,
  }) {
    final result = InviteMemberRequest._();
    if (username != null) result.username = username;
    if (email != null) result.email = email;
    if (role != null) result.role = role;
    return result;
  }

  InviteMemberRequest._();

  factory InviteMemberRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InviteMemberRequest()..mergeFromBuffer(data, registry);
  factory InviteMemberRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InviteMemberRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InviteMemberRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: InviteMemberRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'username')
    ..aOS(2, _omitFieldNames ? '' : 'email')
    ..aE<$2.UserRole>(3, _omitFieldNames ? '' : 'role',
        enumValues: $2.UserRole.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InviteMemberRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InviteMemberRequest copyWith(void Function(InviteMemberRequest) updates) =>
      super.copyWith((message) => updates(message as InviteMemberRequest))
          as InviteMemberRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use InviteMemberRequest() / InviteMemberRequest.new instead')
  static InviteMemberRequest create() => InviteMemberRequest._();
  static $pb.GeneratedMessage $_createMessage() => InviteMemberRequest._();
  @$core.override
  InviteMemberRequest createEmptyInstance() => InviteMemberRequest._();
  @$core.pragma('dart2js:noInline')
  static InviteMemberRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InviteMemberRequest>(
          InviteMemberRequest.$_createMessage);
  static InviteMemberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get username => $_getSZ(0);
  @$pb.TagNumber(1)
  set username($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUsername() => $_has(0);
  @$pb.TagNumber(1)
  void clearUsername() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get email => $_getSZ(1);
  @$pb.TagNumber(2)
  set email($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasEmail() => $_has(1);
  @$pb.TagNumber(2)
  void clearEmail() => $_clearField(2);

  /// USER_ROLE_ADMIN ou USER_ROLE_MEMBER.
  @$pb.TagNumber(3)
  $2.UserRole get role => $_getN(2);
  @$pb.TagNumber(3)
  set role($2.UserRole value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasRole() => $_has(2);
  @$pb.TagNumber(3)
  void clearRole() => $_clearField(3);
}

/// InvitationResult est la réponse commune à une invitation et à son renvoi.
class InvitationResult extends $pb.GeneratedMessage {
  factory InvitationResult({
    Member? member,
    $core.bool? emailed,
    $core.String? invitationUrl,
    $core.bool? localDashboard,
  }) {
    final result = InvitationResult._();
    if (member != null) result.member = member;
    if (emailed != null) result.emailed = emailed;
    if (invitationUrl != null) result.invitationUrl = invitationUrl;
    if (localDashboard != null) result.localDashboard = localDashboard;
    return result;
  }

  InvitationResult._();

  factory InvitationResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InvitationResult()..mergeFromBuffer(data, registry);
  factory InvitationResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InvitationResult()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvitationResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: InvitationResult.$_createMessage)
    ..aOM<Member>(1, _omitFieldNames ? '' : 'member',
        subBuilder: Member.$_createMessage)
    ..aOB(2, _omitFieldNames ? '' : 'emailed')
    ..aOS(3, _omitFieldNames ? '' : 'invitationUrl')
    ..aOB(4, _omitFieldNames ? '' : 'localDashboard')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvitationResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvitationResult copyWith(void Function(InvitationResult) updates) =>
      super.copyWith((message) => updates(message as InvitationResult))
          as InvitationResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use InvitationResult() / InvitationResult.new instead')
  static InvitationResult create() => InvitationResult._();
  static $pb.GeneratedMessage $_createMessage() => InvitationResult._();
  @$core.override
  InvitationResult createEmptyInstance() => InvitationResult._();
  @$core.pragma('dart2js:noInline')
  static InvitationResult getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<InvitationResult>(
          InvitationResult.$_createMessage);
  static InvitationResult? _defaultInstance;

  @$pb.TagNumber(1)
  Member get member => $_getN(0);
  @$pb.TagNumber(1)
  set member(Member value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMember() => $_has(0);
  @$pb.TagNumber(1)
  void clearMember() => $_clearField(1);
  @$pb.TagNumber(1)
  Member ensureMember() => $_ensure(0);

  /// Vrai si le lien est parti par e-mail (en tâche de fond).
  @$pb.TagNumber(2)
  $core.bool get emailed => $_getBF(1);
  @$pb.TagNumber(2)
  set emailed($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasEmailed() => $_has(1);
  @$pb.TagNumber(2)
  void clearEmailed() => $_clearField(2);

  /// Le lien d'invitation, seulement quand il n'est pas parti par e-mail. Il
  /// n'est jamais rendu une seconde fois : un nouveau lien annule celui-ci.
  @$pb.TagNumber(3)
  $core.String get invitationUrl => $_getSZ(2);
  @$pb.TagNumber(3)
  set invitationUrl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasInvitationUrl() => $_has(2);
  @$pb.TagNumber(3)
  void clearInvitationUrl() => $_clearField(3);

  /// Vrai si le dashboard n'a pas de domaine : le lien s'ouvre par un tunnel
  /// SSH (ADR-0009), et l'invité doit le savoir.
  @$pb.TagNumber(4)
  $core.bool get localDashboard => $_getBF(3);
  @$pb.TagNumber(4)
  set localDashboard($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLocalDashboard() => $_has(3);
  @$pb.TagNumber(4)
  void clearLocalDashboard() => $_clearField(4);
}

class InviteMemberResponse extends $pb.GeneratedMessage {
  factory InviteMemberResponse({
    InvitationResult? result,
  }) {
    final result$ = InviteMemberResponse._();
    if (result != null) result$.result = result;
    return result$;
  }

  InviteMemberResponse._();

  factory InviteMemberResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InviteMemberResponse()..mergeFromBuffer(data, registry);
  factory InviteMemberResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      InviteMemberResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InviteMemberResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: InviteMemberResponse.$_createMessage)
    ..aOM<InvitationResult>(1, _omitFieldNames ? '' : 'result',
        subBuilder: InvitationResult.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InviteMemberResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InviteMemberResponse copyWith(void Function(InviteMemberResponse) updates) =>
      super.copyWith((message) => updates(message as InviteMemberResponse))
          as InviteMemberResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use InviteMemberResponse() / InviteMemberResponse.new instead')
  static InviteMemberResponse create() => InviteMemberResponse._();
  static $pb.GeneratedMessage $_createMessage() => InviteMemberResponse._();
  @$core.override
  InviteMemberResponse createEmptyInstance() => InviteMemberResponse._();
  @$core.pragma('dart2js:noInline')
  static InviteMemberResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InviteMemberResponse>(
          InviteMemberResponse.$_createMessage);
  static InviteMemberResponse? _defaultInstance;

  @$pb.TagNumber(1)
  InvitationResult get result => $_getN(0);
  @$pb.TagNumber(1)
  set result(InvitationResult value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasResult() => $_has(0);
  @$pb.TagNumber(1)
  void clearResult() => $_clearField(1);
  @$pb.TagNumber(1)
  InvitationResult ensureResult() => $_ensure(0);
}

class RenewInvitationRequest extends $pb.GeneratedMessage {
  factory RenewInvitationRequest({
    $core.String? userId,
  }) {
    final result = RenewInvitationRequest._();
    if (userId != null) result.userId = userId;
    return result;
  }

  RenewInvitationRequest._();

  factory RenewInvitationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RenewInvitationRequest()..mergeFromBuffer(data, registry);
  factory RenewInvitationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RenewInvitationRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RenewInvitationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RenewInvitationRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewInvitationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewInvitationRequest copyWith(
          void Function(RenewInvitationRequest) updates) =>
      super.copyWith((message) => updates(message as RenewInvitationRequest))
          as RenewInvitationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RenewInvitationRequest() / RenewInvitationRequest.new instead')
  static RenewInvitationRequest create() => RenewInvitationRequest._();
  static $pb.GeneratedMessage $_createMessage() => RenewInvitationRequest._();
  @$core.override
  RenewInvitationRequest createEmptyInstance() => RenewInvitationRequest._();
  @$core.pragma('dart2js:noInline')
  static RenewInvitationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RenewInvitationRequest>(
          RenewInvitationRequest.$_createMessage);
  static RenewInvitationRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);
}

class RenewInvitationResponse extends $pb.GeneratedMessage {
  factory RenewInvitationResponse({
    InvitationResult? result,
  }) {
    final result$ = RenewInvitationResponse._();
    if (result != null) result$.result = result;
    return result$;
  }

  RenewInvitationResponse._();

  factory RenewInvitationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RenewInvitationResponse()..mergeFromBuffer(data, registry);
  factory RenewInvitationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RenewInvitationResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RenewInvitationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RenewInvitationResponse.$_createMessage)
    ..aOM<InvitationResult>(1, _omitFieldNames ? '' : 'result',
        subBuilder: InvitationResult.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewInvitationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewInvitationResponse copyWith(
          void Function(RenewInvitationResponse) updates) =>
      super.copyWith((message) => updates(message as RenewInvitationResponse))
          as RenewInvitationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RenewInvitationResponse() / RenewInvitationResponse.new instead')
  static RenewInvitationResponse create() => RenewInvitationResponse._();
  static $pb.GeneratedMessage $_createMessage() => RenewInvitationResponse._();
  @$core.override
  RenewInvitationResponse createEmptyInstance() => RenewInvitationResponse._();
  @$core.pragma('dart2js:noInline')
  static RenewInvitationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RenewInvitationResponse>(
          RenewInvitationResponse.$_createMessage);
  static RenewInvitationResponse? _defaultInstance;

  @$pb.TagNumber(1)
  InvitationResult get result => $_getN(0);
  @$pb.TagNumber(1)
  set result(InvitationResult value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasResult() => $_has(0);
  @$pb.TagNumber(1)
  void clearResult() => $_clearField(1);
  @$pb.TagNumber(1)
  InvitationResult ensureResult() => $_ensure(0);
}

class SuspendMemberRequest extends $pb.GeneratedMessage {
  factory SuspendMemberRequest({
    $core.String? userId,
  }) {
    final result = SuspendMemberRequest._();
    if (userId != null) result.userId = userId;
    return result;
  }

  SuspendMemberRequest._();

  factory SuspendMemberRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SuspendMemberRequest()..mergeFromBuffer(data, registry);
  factory SuspendMemberRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SuspendMemberRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SuspendMemberRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SuspendMemberRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SuspendMemberRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SuspendMemberRequest copyWith(void Function(SuspendMemberRequest) updates) =>
      super.copyWith((message) => updates(message as SuspendMemberRequest))
          as SuspendMemberRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SuspendMemberRequest() / SuspendMemberRequest.new instead')
  static SuspendMemberRequest create() => SuspendMemberRequest._();
  static $pb.GeneratedMessage $_createMessage() => SuspendMemberRequest._();
  @$core.override
  SuspendMemberRequest createEmptyInstance() => SuspendMemberRequest._();
  @$core.pragma('dart2js:noInline')
  static SuspendMemberRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SuspendMemberRequest>(
          SuspendMemberRequest.$_createMessage);
  static SuspendMemberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);
}

class SuspendMemberResponse extends $pb.GeneratedMessage {
  factory SuspendMemberResponse({
    Member? member,
    $fixnum.Int64? closedSessions,
  }) {
    final result = SuspendMemberResponse._();
    if (member != null) result.member = member;
    if (closedSessions != null) result.closedSessions = closedSessions;
    return result;
  }

  SuspendMemberResponse._();

  factory SuspendMemberResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SuspendMemberResponse()..mergeFromBuffer(data, registry);
  factory SuspendMemberResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SuspendMemberResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SuspendMemberResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SuspendMemberResponse.$_createMessage)
    ..aOM<Member>(1, _omitFieldNames ? '' : 'member',
        subBuilder: Member.$_createMessage)
    ..aInt64(2, _omitFieldNames ? '' : 'closedSessions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SuspendMemberResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SuspendMemberResponse copyWith(
          void Function(SuspendMemberResponse) updates) =>
      super.copyWith((message) => updates(message as SuspendMemberResponse))
          as SuspendMemberResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SuspendMemberResponse() / SuspendMemberResponse.new instead')
  static SuspendMemberResponse create() => SuspendMemberResponse._();
  static $pb.GeneratedMessage $_createMessage() => SuspendMemberResponse._();
  @$core.override
  SuspendMemberResponse createEmptyInstance() => SuspendMemberResponse._();
  @$core.pragma('dart2js:noInline')
  static SuspendMemberResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SuspendMemberResponse>(
          SuspendMemberResponse.$_createMessage);
  static SuspendMemberResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Member get member => $_getN(0);
  @$pb.TagNumber(1)
  set member(Member value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMember() => $_has(0);
  @$pb.TagNumber(1)
  void clearMember() => $_clearField(1);
  @$pb.TagNumber(1)
  Member ensureMember() => $_ensure(0);

  /// Nombre de sessions et de jetons d'API fermés.
  @$pb.TagNumber(2)
  $fixnum.Int64 get closedSessions => $_getI64(1);
  @$pb.TagNumber(2)
  set closedSessions($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasClosedSessions() => $_has(1);
  @$pb.TagNumber(2)
  void clearClosedSessions() => $_clearField(2);
}

class ResumeMemberRequest extends $pb.GeneratedMessage {
  factory ResumeMemberRequest({
    $core.String? userId,
  }) {
    final result = ResumeMemberRequest._();
    if (userId != null) result.userId = userId;
    return result;
  }

  ResumeMemberRequest._();

  factory ResumeMemberRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResumeMemberRequest()..mergeFromBuffer(data, registry);
  factory ResumeMemberRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResumeMemberRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResumeMemberRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ResumeMemberRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeMemberRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeMemberRequest copyWith(void Function(ResumeMemberRequest) updates) =>
      super.copyWith((message) => updates(message as ResumeMemberRequest))
          as ResumeMemberRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use ResumeMemberRequest() / ResumeMemberRequest.new instead')
  static ResumeMemberRequest create() => ResumeMemberRequest._();
  static $pb.GeneratedMessage $_createMessage() => ResumeMemberRequest._();
  @$core.override
  ResumeMemberRequest createEmptyInstance() => ResumeMemberRequest._();
  @$core.pragma('dart2js:noInline')
  static ResumeMemberRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResumeMemberRequest>(
          ResumeMemberRequest.$_createMessage);
  static ResumeMemberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);
}

class ResumeMemberResponse extends $pb.GeneratedMessage {
  factory ResumeMemberResponse({
    Member? member,
  }) {
    final result = ResumeMemberResponse._();
    if (member != null) result.member = member;
    return result;
  }

  ResumeMemberResponse._();

  factory ResumeMemberResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResumeMemberResponse()..mergeFromBuffer(data, registry);
  factory ResumeMemberResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResumeMemberResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResumeMemberResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ResumeMemberResponse.$_createMessage)
    ..aOM<Member>(1, _omitFieldNames ? '' : 'member',
        subBuilder: Member.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeMemberResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeMemberResponse copyWith(void Function(ResumeMemberResponse) updates) =>
      super.copyWith((message) => updates(message as ResumeMemberResponse))
          as ResumeMemberResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ResumeMemberResponse() / ResumeMemberResponse.new instead')
  static ResumeMemberResponse create() => ResumeMemberResponse._();
  static $pb.GeneratedMessage $_createMessage() => ResumeMemberResponse._();
  @$core.override
  ResumeMemberResponse createEmptyInstance() => ResumeMemberResponse._();
  @$core.pragma('dart2js:noInline')
  static ResumeMemberResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResumeMemberResponse>(
          ResumeMemberResponse.$_createMessage);
  static ResumeMemberResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Member get member => $_getN(0);
  @$pb.TagNumber(1)
  set member(Member value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMember() => $_has(0);
  @$pb.TagNumber(1)
  void clearMember() => $_clearField(1);
  @$pb.TagNumber(1)
  Member ensureMember() => $_ensure(0);
}

class CloseMemberSessionsRequest extends $pb.GeneratedMessage {
  factory CloseMemberSessionsRequest({
    $core.String? userId,
  }) {
    final result = CloseMemberSessionsRequest._();
    if (userId != null) result.userId = userId;
    return result;
  }

  CloseMemberSessionsRequest._();

  factory CloseMemberSessionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CloseMemberSessionsRequest()..mergeFromBuffer(data, registry);
  factory CloseMemberSessionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CloseMemberSessionsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloseMemberSessionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CloseMemberSessionsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloseMemberSessionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloseMemberSessionsRequest copyWith(
          void Function(CloseMemberSessionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CloseMemberSessionsRequest))
          as CloseMemberSessionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CloseMemberSessionsRequest() / CloseMemberSessionsRequest.new instead')
  static CloseMemberSessionsRequest create() => CloseMemberSessionsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      CloseMemberSessionsRequest._();
  @$core.override
  CloseMemberSessionsRequest createEmptyInstance() =>
      CloseMemberSessionsRequest._();
  @$core.pragma('dart2js:noInline')
  static CloseMemberSessionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloseMemberSessionsRequest>(
          CloseMemberSessionsRequest.$_createMessage);
  static CloseMemberSessionsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);
}

class CloseMemberSessionsResponse extends $pb.GeneratedMessage {
  factory CloseMemberSessionsResponse({
    $fixnum.Int64? closedSessions,
  }) {
    final result = CloseMemberSessionsResponse._();
    if (closedSessions != null) result.closedSessions = closedSessions;
    return result;
  }

  CloseMemberSessionsResponse._();

  factory CloseMemberSessionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CloseMemberSessionsResponse()..mergeFromBuffer(data, registry);
  factory CloseMemberSessionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CloseMemberSessionsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloseMemberSessionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CloseMemberSessionsResponse.$_createMessage)
    ..aInt64(1, _omitFieldNames ? '' : 'closedSessions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloseMemberSessionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloseMemberSessionsResponse copyWith(
          void Function(CloseMemberSessionsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CloseMemberSessionsResponse))
          as CloseMemberSessionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CloseMemberSessionsResponse() / CloseMemberSessionsResponse.new instead')
  static CloseMemberSessionsResponse create() =>
      CloseMemberSessionsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      CloseMemberSessionsResponse._();
  @$core.override
  CloseMemberSessionsResponse createEmptyInstance() =>
      CloseMemberSessionsResponse._();
  @$core.pragma('dart2js:noInline')
  static CloseMemberSessionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloseMemberSessionsResponse>(
          CloseMemberSessionsResponse.$_createMessage);
  static CloseMemberSessionsResponse? _defaultInstance;

  /// Nombre de sessions et de jetons d'API fermés.
  @$pb.TagNumber(1)
  $fixnum.Int64 get closedSessions => $_getI64(0);
  @$pb.TagNumber(1)
  set closedSessions($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasClosedSessions() => $_has(0);
  @$pb.TagNumber(1)
  void clearClosedSessions() => $_clearField(1);
}

class DeleteMemberRequest extends $pb.GeneratedMessage {
  factory DeleteMemberRequest({
    $core.String? userId,
    $core.String? confirmUsername,
  }) {
    final result = DeleteMemberRequest._();
    if (userId != null) result.userId = userId;
    if (confirmUsername != null) result.confirmUsername = confirmUsername;
    return result;
  }

  DeleteMemberRequest._();

  factory DeleteMemberRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteMemberRequest()..mergeFromBuffer(data, registry);
  factory DeleteMemberRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteMemberRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMemberRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DeleteMemberRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..aOS(2, _omitFieldNames ? '' : 'confirmUsername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMemberRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMemberRequest copyWith(void Function(DeleteMemberRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteMemberRequest))
          as DeleteMemberRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DeleteMemberRequest() / DeleteMemberRequest.new instead')
  static DeleteMemberRequest create() => DeleteMemberRequest._();
  static $pb.GeneratedMessage $_createMessage() => DeleteMemberRequest._();
  @$core.override
  DeleteMemberRequest createEmptyInstance() => DeleteMemberRequest._();
  @$core.pragma('dart2js:noInline')
  static DeleteMemberRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMemberRequest>(
          DeleteMemberRequest.$_createMessage);
  static DeleteMemberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);

  /// Le nom du compte, retapé pour confirmer : une suppression ne part pas
  /// d'un clic ou d'un identifiant mal recopié.
  @$pb.TagNumber(2)
  $core.String get confirmUsername => $_getSZ(1);
  @$pb.TagNumber(2)
  set confirmUsername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasConfirmUsername() => $_has(1);
  @$pb.TagNumber(2)
  void clearConfirmUsername() => $_clearField(2);
}

class DeleteMemberResponse extends $pb.GeneratedMessage {
  factory DeleteMemberResponse() => DeleteMemberResponse._();

  DeleteMemberResponse._();

  factory DeleteMemberResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteMemberResponse()..mergeFromBuffer(data, registry);
  factory DeleteMemberResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeleteMemberResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMemberResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DeleteMemberResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMemberResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMemberResponse copyWith(void Function(DeleteMemberResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteMemberResponse))
          as DeleteMemberResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use DeleteMemberResponse() / DeleteMemberResponse.new instead')
  static DeleteMemberResponse create() => DeleteMemberResponse._();
  static $pb.GeneratedMessage $_createMessage() => DeleteMemberResponse._();
  @$core.override
  DeleteMemberResponse createEmptyInstance() => DeleteMemberResponse._();
  @$core.pragma('dart2js:noInline')
  static DeleteMemberResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMemberResponse>(
          DeleteMemberResponse.$_createMessage);
  static DeleteMemberResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
