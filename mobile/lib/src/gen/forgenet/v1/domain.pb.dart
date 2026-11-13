// This is a generated file - do not edit.
//
// Generated from forgenet/v1/domain.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// Domain relie un nom de domaine au service d'une application déployée :
/// Caddy obtient son certificat et lui transmet le trafic (ADR-0007, ADR-0008).
class Domain extends $pb.GeneratedMessage {
  factory Domain({
    $core.String? id,
    $core.String? host,
    $core.String? projectId,
    $core.String? service,
    $core.int? port,
    $core.String? createdBy,
    $1.Timestamp? createdAt,
  }) {
    final result = Domain._();
    if (id != null) result.id = id;
    if (host != null) result.host = host;
    if (projectId != null) result.projectId = projectId;
    if (service != null) result.service = service;
    if (port != null) result.port = port;
    if (createdBy != null) result.createdBy = createdBy;
    if (createdAt != null) result.createdAt = createdAt;
    return result;
  }

  Domain._();

  factory Domain.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Domain()..mergeFromBuffer(data, registry);
  factory Domain.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Domain()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Domain',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: Domain.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'host')
    ..aOS(3, _omitFieldNames ? '' : 'projectId')
    ..aOS(4, _omitFieldNames ? '' : 'service')
    ..aI(5, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aOS(6, _omitFieldNames ? '' : 'createdBy')
    ..aOM<$1.Timestamp>(7, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Domain clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Domain copyWith(void Function(Domain) updates) =>
      super.copyWith((message) => updates(message as Domain)) as Domain;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Domain() / Domain.new instead')
  static Domain create() => Domain._();
  static $pb.GeneratedMessage $_createMessage() => Domain._();
  @$core.override
  Domain createEmptyInstance() => Domain._();
  @$core.pragma('dart2js:noInline')
  static Domain getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Domain>(Domain.$_createMessage);
  static Domain? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  /// Nom normalisé : minuscules, sans point final, nom internationalisé en
  /// punycode (« xn-- »).
  @$pb.TagNumber(2)
  $core.String get host => $_getSZ(1);
  @$pb.TagNumber(2)
  set host($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasHost() => $_has(1);
  @$pb.TagNumber(2)
  void clearHost() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get projectId => $_getSZ(2);
  @$pb.TagNumber(3)
  set projectId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProjectId() => $_has(2);
  @$pb.TagNumber(3)
  void clearProjectId() => $_clearField(3);

  /// Service du fichier Compose du projet.
  @$pb.TagNumber(4)
  $core.String get service => $_getSZ(3);
  @$pb.TagNumber(4)
  set service($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasService() => $_has(3);
  @$pb.TagNumber(4)
  void clearService() => $_clearField(4);

  /// Port où le service écoute dans son conteneur.
  @$pb.TagNumber(5)
  $core.int get port => $_getIZ(4);
  @$pb.TagNumber(5)
  set port($core.int value) => $_setUnsignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPort() => $_has(4);
  @$pb.TagNumber(5)
  void clearPort() => $_clearField(5);

  /// Auteur de l'ajout, ex: "user:usr_…".
  @$pb.TagNumber(6)
  $core.String get createdBy => $_getSZ(5);
  @$pb.TagNumber(6)
  set createdBy($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCreatedBy() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedBy() => $_clearField(6);

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
}

class ListDomainsRequest extends $pb.GeneratedMessage {
  factory ListDomainsRequest({
    $core.String? projectId,
  }) {
    final result = ListDomainsRequest._();
    if (projectId != null) result.projectId = projectId;
    return result;
  }

  ListDomainsRequest._();

  factory ListDomainsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDomainsRequest()..mergeFromBuffer(data, registry);
  factory ListDomainsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDomainsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDomainsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListDomainsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'projectId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDomainsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDomainsRequest copyWith(void Function(ListDomainsRequest) updates) =>
      super.copyWith((message) => updates(message as ListDomainsRequest))
          as ListDomainsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ListDomainsRequest() / ListDomainsRequest.new instead')
  static ListDomainsRequest create() => ListDomainsRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListDomainsRequest._();
  @$core.override
  ListDomainsRequest createEmptyInstance() => ListDomainsRequest._();
  @$core.pragma('dart2js:noInline')
  static ListDomainsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDomainsRequest>(
          ListDomainsRequest.$_createMessage);
  static ListDomainsRequest? _defaultInstance;

  /// Vide : les domaines de tous les projets.
  @$pb.TagNumber(1)
  $core.String get projectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set projectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasProjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProjectId() => $_clearField(1);
}

class ListDomainsResponse extends $pb.GeneratedMessage {
  factory ListDomainsResponse({
    $core.Iterable<Domain>? domains,
  }) {
    final result = ListDomainsResponse._();
    if (domains != null) result.domains.addAll(domains);
    return result;
  }

  ListDomainsResponse._();

  factory ListDomainsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDomainsResponse()..mergeFromBuffer(data, registry);
  factory ListDomainsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDomainsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDomainsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListDomainsResponse.$_createMessage)
    ..pPM<Domain>(1, _omitFieldNames ? '' : 'domains',
        subBuilder: Domain.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDomainsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDomainsResponse copyWith(void Function(ListDomainsResponse) updates) =>
      super.copyWith((message) => updates(message as ListDomainsResponse))
          as ListDomainsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use ListDomainsResponse() / ListDomainsResponse.new instead')
  static ListDomainsResponse create() => ListDomainsResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListDomainsResponse._();
  @$core.override
  ListDomainsResponse createEmptyInstance() => ListDomainsResponse._();
  @$core.pragma('dart2js:noInline')
  static ListDomainsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDomainsResponse>(
          ListDomainsResponse.$_createMessage);
  static ListDomainsResponse? _defaultInstance;

  /// Triés par nom de domaine.
  @$pb.TagNumber(1)
  $pb.PbList<Domain> get domains => $_getList(0);
}

class AddDomainRequest extends $pb.GeneratedMessage {
  factory AddDomainRequest({
    $core.String? host,
    $core.String? projectId,
    $core.String? service,
    $core.int? port,
  }) {
    final result = AddDomainRequest._();
    if (host != null) result.host = host;
    if (projectId != null) result.projectId = projectId;
    if (service != null) result.service = service;
    if (port != null) result.port = port;
    return result;
  }

  AddDomainRequest._();

  factory AddDomainRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddDomainRequest()..mergeFromBuffer(data, registry);
  factory AddDomainRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddDomainRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddDomainRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: AddDomainRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'host')
    ..aOS(2, _omitFieldNames ? '' : 'projectId')
    ..aOS(3, _omitFieldNames ? '' : 'service')
    ..aI(4, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddDomainRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddDomainRequest copyWith(void Function(AddDomainRequest) updates) =>
      super.copyWith((message) => updates(message as AddDomainRequest))
          as AddDomainRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddDomainRequest() / AddDomainRequest.new instead')
  static AddDomainRequest create() => AddDomainRequest._();
  static $pb.GeneratedMessage $_createMessage() => AddDomainRequest._();
  @$core.override
  AddDomainRequest createEmptyInstance() => AddDomainRequest._();
  @$core.pragma('dart2js:noInline')
  static AddDomainRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddDomainRequest>(
          AddDomainRequest.$_createMessage);
  static AddDomainRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get host => $_getSZ(0);
  @$pb.TagNumber(1)
  set host($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHost() => $_has(0);
  @$pb.TagNumber(1)
  void clearHost() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get projectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set projectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasProjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearProjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get service => $_getSZ(2);
  @$pb.TagNumber(3)
  set service($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasService() => $_has(2);
  @$pb.TagNumber(3)
  void clearService() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get port => $_getIZ(3);
  @$pb.TagNumber(4)
  set port($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPort() => $_has(3);
  @$pb.TagNumber(4)
  void clearPort() => $_clearField(4);
}

class AddDomainResponse extends $pb.GeneratedMessage {
  factory AddDomainResponse({
    Domain? domain,
  }) {
    final result = AddDomainResponse._();
    if (domain != null) result.domain = domain;
    return result;
  }

  AddDomainResponse._();

  factory AddDomainResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddDomainResponse()..mergeFromBuffer(data, registry);
  factory AddDomainResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AddDomainResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddDomainResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: AddDomainResponse.$_createMessage)
    ..aOM<Domain>(1, _omitFieldNames ? '' : 'domain',
        subBuilder: Domain.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddDomainResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddDomainResponse copyWith(void Function(AddDomainResponse) updates) =>
      super.copyWith((message) => updates(message as AddDomainResponse))
          as AddDomainResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AddDomainResponse() / AddDomainResponse.new instead')
  static AddDomainResponse create() => AddDomainResponse._();
  static $pb.GeneratedMessage $_createMessage() => AddDomainResponse._();
  @$core.override
  AddDomainResponse createEmptyInstance() => AddDomainResponse._();
  @$core.pragma('dart2js:noInline')
  static AddDomainResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddDomainResponse>(
          AddDomainResponse.$_createMessage);
  static AddDomainResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Domain get domain => $_getN(0);
  @$pb.TagNumber(1)
  set domain(Domain value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDomain() => $_has(0);
  @$pb.TagNumber(1)
  void clearDomain() => $_clearField(1);
  @$pb.TagNumber(1)
  Domain ensureDomain() => $_ensure(0);
}

class RemoveDomainRequest extends $pb.GeneratedMessage {
  factory RemoveDomainRequest({
    $core.String? domainId,
  }) {
    final result = RemoveDomainRequest._();
    if (domainId != null) result.domainId = domainId;
    return result;
  }

  RemoveDomainRequest._();

  factory RemoveDomainRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveDomainRequest()..mergeFromBuffer(data, registry);
  factory RemoveDomainRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveDomainRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveDomainRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RemoveDomainRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'domainId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveDomainRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveDomainRequest copyWith(void Function(RemoveDomainRequest) updates) =>
      super.copyWith((message) => updates(message as RemoveDomainRequest))
          as RemoveDomainRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use RemoveDomainRequest() / RemoveDomainRequest.new instead')
  static RemoveDomainRequest create() => RemoveDomainRequest._();
  static $pb.GeneratedMessage $_createMessage() => RemoveDomainRequest._();
  @$core.override
  RemoveDomainRequest createEmptyInstance() => RemoveDomainRequest._();
  @$core.pragma('dart2js:noInline')
  static RemoveDomainRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveDomainRequest>(
          RemoveDomainRequest.$_createMessage);
  static RemoveDomainRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get domainId => $_getSZ(0);
  @$pb.TagNumber(1)
  set domainId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDomainId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDomainId() => $_clearField(1);
}

class RemoveDomainResponse extends $pb.GeneratedMessage {
  factory RemoveDomainResponse() => RemoveDomainResponse._();

  RemoveDomainResponse._();

  factory RemoveDomainResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveDomainResponse()..mergeFromBuffer(data, registry);
  factory RemoveDomainResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RemoveDomainResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveDomainResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RemoveDomainResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveDomainResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveDomainResponse copyWith(void Function(RemoveDomainResponse) updates) =>
      super.copyWith((message) => updates(message as RemoveDomainResponse))
          as RemoveDomainResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RemoveDomainResponse() / RemoveDomainResponse.new instead')
  static RemoveDomainResponse create() => RemoveDomainResponse._();
  static $pb.GeneratedMessage $_createMessage() => RemoveDomainResponse._();
  @$core.override
  RemoveDomainResponse createEmptyInstance() => RemoveDomainResponse._();
  @$core.pragma('dart2js:noInline')
  static RemoveDomainResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveDomainResponse>(
          RemoveDomainResponse.$_createMessage);
  static RemoveDomainResponse? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
