// This is a generated file - do not edit.
//
// Generated from forgenet/v1/container.proto.

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

import 'container.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'container.pbenum.dart';

/// PortMapping définit la redirection d'un port hôte vers un port conteneur.
class PortMapping extends $pb.GeneratedMessage {
  factory PortMapping({
    $core.int? hostPort,
    $core.int? containerPort,
    $core.String? protocol,
  }) {
    final result = PortMapping._();
    if (hostPort != null) result.hostPort = hostPort;
    if (containerPort != null) result.containerPort = containerPort;
    if (protocol != null) result.protocol = protocol;
    return result;
  }

  PortMapping._();

  factory PortMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PortMapping()..mergeFromBuffer(data, registry);
  factory PortMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      PortMapping()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PortMapping',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: PortMapping.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'hostPort', fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'containerPort',
        fieldType: $pb.PbFieldType.OU3)
    ..aOS(3, _omitFieldNames ? '' : 'protocol')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PortMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PortMapping copyWith(void Function(PortMapping) updates) =>
      super.copyWith((message) => updates(message as PortMapping))
          as PortMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use PortMapping() / PortMapping.new instead')
  static PortMapping create() => PortMapping._();
  static $pb.GeneratedMessage $_createMessage() => PortMapping._();
  @$core.override
  PortMapping createEmptyInstance() => PortMapping._();
  @$core.pragma('dart2js:noInline')
  static PortMapping getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PortMapping>(
          PortMapping.$_createMessage);
  static PortMapping? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get hostPort => $_getIZ(0);
  @$pb.TagNumber(1)
  set hostPort($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHostPort() => $_has(0);
  @$pb.TagNumber(1)
  void clearHostPort() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get containerPort => $_getIZ(1);
  @$pb.TagNumber(2)
  set containerPort($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContainerPort() => $_has(1);
  @$pb.TagNumber(2)
  void clearContainerPort() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get protocol => $_getSZ(2);
  @$pb.TagNumber(3)
  set protocol($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasProtocol() => $_has(2);
  @$pb.TagNumber(3)
  void clearProtocol() => $_clearField(3);
}

/// Container représente l'état complet et les métadonnées d'un conteneur.
class Container extends $pb.GeneratedMessage {
  factory Container({
    $core.String? id,
    $core.String? name,
    $core.String? image,
    ContainerState? state,
    $core.Iterable<PortMapping>? ports,
    $1.Timestamp? createdAt,
  }) {
    final result = Container._();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (image != null) result.image = image;
    if (state != null) result.state = state;
    if (ports != null) result.ports.addAll(ports);
    if (createdAt != null) result.createdAt = createdAt;
    return result;
  }

  Container._();

  factory Container.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Container()..mergeFromBuffer(data, registry);
  factory Container.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Container()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Container',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: Container.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'image')
    ..aE<ContainerState>(4, _omitFieldNames ? '' : 'state',
        enumValues: ContainerState.values)
    ..pPM<PortMapping>(5, _omitFieldNames ? '' : 'ports',
        subBuilder: PortMapping.$_createMessage)
    ..aOM<$1.Timestamp>(6, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Container clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Container copyWith(void Function(Container) updates) =>
      super.copyWith((message) => updates(message as Container)) as Container;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Container() / Container.new instead')
  static Container create() => Container._();
  static $pb.GeneratedMessage $_createMessage() => Container._();
  @$core.override
  Container createEmptyInstance() => Container._();
  @$core.pragma('dart2js:noInline')
  static Container getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Container>(Container.$_createMessage);
  static Container? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get image => $_getSZ(2);
  @$pb.TagNumber(3)
  set image($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasImage() => $_has(2);
  @$pb.TagNumber(3)
  void clearImage() => $_clearField(3);

  @$pb.TagNumber(4)
  ContainerState get state => $_getN(3);
  @$pb.TagNumber(4)
  set state(ContainerState value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(4)
  void clearState() => $_clearField(4);

  @$pb.TagNumber(5)
  $pb.PbList<PortMapping> get ports => $_getList(4);

  @$pb.TagNumber(6)
  $1.Timestamp get createdAt => $_getN(5);
  @$pb.TagNumber(6)
  set createdAt($1.Timestamp value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => $_clearField(6);
  @$pb.TagNumber(6)
  $1.Timestamp ensureCreatedAt() => $_ensure(5);
}

class ListContainersRequest extends $pb.GeneratedMessage {
  factory ListContainersRequest({
    $core.bool? includeStopped,
    $core.int? pageSize,
    $core.String? pageToken,
  }) {
    final result = ListContainersRequest._();
    if (includeStopped != null) result.includeStopped = includeStopped;
    if (pageSize != null) result.pageSize = pageSize;
    if (pageToken != null) result.pageToken = pageToken;
    return result;
  }

  ListContainersRequest._();

  factory ListContainersRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListContainersRequest()..mergeFromBuffer(data, registry);
  factory ListContainersRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListContainersRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListContainersRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListContainersRequest.$_createMessage)
    ..aOB(1, _omitFieldNames ? '' : 'includeStopped')
    ..aI(2, _omitFieldNames ? '' : 'pageSize')
    ..aOS(3, _omitFieldNames ? '' : 'pageToken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContainersRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContainersRequest copyWith(
          void Function(ListContainersRequest) updates) =>
      super.copyWith((message) => updates(message as ListContainersRequest))
          as ListContainersRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListContainersRequest() / ListContainersRequest.new instead')
  static ListContainersRequest create() => ListContainersRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListContainersRequest._();
  @$core.override
  ListContainersRequest createEmptyInstance() => ListContainersRequest._();
  @$core.pragma('dart2js:noInline')
  static ListContainersRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListContainersRequest>(
          ListContainersRequest.$_createMessage);
  static ListContainersRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get includeStopped => $_getBF(0);
  @$pb.TagNumber(1)
  set includeStopped($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIncludeStopped() => $_has(0);
  @$pb.TagNumber(1)
  void clearIncludeStopped() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get pageSize => $_getIZ(1);
  @$pb.TagNumber(2)
  set pageSize($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPageSize() => $_has(1);
  @$pb.TagNumber(2)
  void clearPageSize() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get pageToken => $_getSZ(2);
  @$pb.TagNumber(3)
  set pageToken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPageToken() => $_has(2);
  @$pb.TagNumber(3)
  void clearPageToken() => $_clearField(3);
}

class ListContainersResponse extends $pb.GeneratedMessage {
  factory ListContainersResponse({
    $core.Iterable<Container>? containers,
    $core.String? nextPageToken,
  }) {
    final result = ListContainersResponse._();
    if (containers != null) result.containers.addAll(containers);
    if (nextPageToken != null) result.nextPageToken = nextPageToken;
    return result;
  }

  ListContainersResponse._();

  factory ListContainersResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListContainersResponse()..mergeFromBuffer(data, registry);
  factory ListContainersResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListContainersResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListContainersResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListContainersResponse.$_createMessage)
    ..pPM<Container>(1, _omitFieldNames ? '' : 'containers',
        subBuilder: Container.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'nextPageToken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContainersResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContainersResponse copyWith(
          void Function(ListContainersResponse) updates) =>
      super.copyWith((message) => updates(message as ListContainersResponse))
          as ListContainersResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListContainersResponse() / ListContainersResponse.new instead')
  static ListContainersResponse create() => ListContainersResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListContainersResponse._();
  @$core.override
  ListContainersResponse createEmptyInstance() => ListContainersResponse._();
  @$core.pragma('dart2js:noInline')
  static ListContainersResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListContainersResponse>(
          ListContainersResponse.$_createMessage);
  static ListContainersResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Container> get containers => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get nextPageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set nextPageToken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasNextPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearNextPageToken() => $_clearField(2);
}

class GetContainerRequest extends $pb.GeneratedMessage {
  factory GetContainerRequest({
    $core.String? containerId,
  }) {
    final result = GetContainerRequest._();
    if (containerId != null) result.containerId = containerId;
    return result;
  }

  GetContainerRequest._();

  factory GetContainerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerRequest()..mergeFromBuffer(data, registry);
  factory GetContainerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContainerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetContainerRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerRequest copyWith(void Function(GetContainerRequest) updates) =>
      super.copyWith((message) => updates(message as GetContainerRequest))
          as GetContainerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use GetContainerRequest() / GetContainerRequest.new instead')
  static GetContainerRequest create() => GetContainerRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetContainerRequest._();
  @$core.override
  GetContainerRequest createEmptyInstance() => GetContainerRequest._();
  @$core.pragma('dart2js:noInline')
  static GetContainerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContainerRequest>(
          GetContainerRequest.$_createMessage);
  static GetContainerRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);
}

class GetContainerResponse extends $pb.GeneratedMessage {
  factory GetContainerResponse({
    Container? container,
  }) {
    final result = GetContainerResponse._();
    if (container != null) result.container = container;
    return result;
  }

  GetContainerResponse._();

  factory GetContainerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerResponse()..mergeFromBuffer(data, registry);
  factory GetContainerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContainerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetContainerResponse.$_createMessage)
    ..aOM<Container>(1, _omitFieldNames ? '' : 'container',
        subBuilder: Container.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerResponse copyWith(void Function(GetContainerResponse) updates) =>
      super.copyWith((message) => updates(message as GetContainerResponse))
          as GetContainerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetContainerResponse() / GetContainerResponse.new instead')
  static GetContainerResponse create() => GetContainerResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetContainerResponse._();
  @$core.override
  GetContainerResponse createEmptyInstance() => GetContainerResponse._();
  @$core.pragma('dart2js:noInline')
  static GetContainerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContainerResponse>(
          GetContainerResponse.$_createMessage);
  static GetContainerResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Container get container => $_getN(0);
  @$pb.TagNumber(1)
  set container(Container value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasContainer() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainer() => $_clearField(1);
  @$pb.TagNumber(1)
  Container ensureContainer() => $_ensure(0);
}

class StartContainerRequest extends $pb.GeneratedMessage {
  factory StartContainerRequest({
    $core.String? containerId,
  }) {
    final result = StartContainerRequest._();
    if (containerId != null) result.containerId = containerId;
    return result;
  }

  StartContainerRequest._();

  factory StartContainerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartContainerRequest()..mergeFromBuffer(data, registry);
  factory StartContainerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartContainerRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartContainerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StartContainerRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartContainerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartContainerRequest copyWith(
          void Function(StartContainerRequest) updates) =>
      super.copyWith((message) => updates(message as StartContainerRequest))
          as StartContainerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StartContainerRequest() / StartContainerRequest.new instead')
  static StartContainerRequest create() => StartContainerRequest._();
  static $pb.GeneratedMessage $_createMessage() => StartContainerRequest._();
  @$core.override
  StartContainerRequest createEmptyInstance() => StartContainerRequest._();
  @$core.pragma('dart2js:noInline')
  static StartContainerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartContainerRequest>(
          StartContainerRequest.$_createMessage);
  static StartContainerRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);
}

class StartContainerResponse extends $pb.GeneratedMessage {
  factory StartContainerResponse({
    Container? container,
  }) {
    final result = StartContainerResponse._();
    if (container != null) result.container = container;
    return result;
  }

  StartContainerResponse._();

  factory StartContainerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartContainerResponse()..mergeFromBuffer(data, registry);
  factory StartContainerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartContainerResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartContainerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StartContainerResponse.$_createMessage)
    ..aOM<Container>(1, _omitFieldNames ? '' : 'container',
        subBuilder: Container.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartContainerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartContainerResponse copyWith(
          void Function(StartContainerResponse) updates) =>
      super.copyWith((message) => updates(message as StartContainerResponse))
          as StartContainerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StartContainerResponse() / StartContainerResponse.new instead')
  static StartContainerResponse create() => StartContainerResponse._();
  static $pb.GeneratedMessage $_createMessage() => StartContainerResponse._();
  @$core.override
  StartContainerResponse createEmptyInstance() => StartContainerResponse._();
  @$core.pragma('dart2js:noInline')
  static StartContainerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartContainerResponse>(
          StartContainerResponse.$_createMessage);
  static StartContainerResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Container get container => $_getN(0);
  @$pb.TagNumber(1)
  set container(Container value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasContainer() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainer() => $_clearField(1);
  @$pb.TagNumber(1)
  Container ensureContainer() => $_ensure(0);
}

class StopContainerRequest extends $pb.GeneratedMessage {
  factory StopContainerRequest({
    $core.String? containerId,
    $core.int? timeoutSeconds,
  }) {
    final result = StopContainerRequest._();
    if (containerId != null) result.containerId = containerId;
    if (timeoutSeconds != null) result.timeoutSeconds = timeoutSeconds;
    return result;
  }

  StopContainerRequest._();

  factory StopContainerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StopContainerRequest()..mergeFromBuffer(data, registry);
  factory StopContainerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StopContainerRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopContainerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StopContainerRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..aI(2, _omitFieldNames ? '' : 'timeoutSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopContainerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopContainerRequest copyWith(void Function(StopContainerRequest) updates) =>
      super.copyWith((message) => updates(message as StopContainerRequest))
          as StopContainerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StopContainerRequest() / StopContainerRequest.new instead')
  static StopContainerRequest create() => StopContainerRequest._();
  static $pb.GeneratedMessage $_createMessage() => StopContainerRequest._();
  @$core.override
  StopContainerRequest createEmptyInstance() => StopContainerRequest._();
  @$core.pragma('dart2js:noInline')
  static StopContainerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopContainerRequest>(
          StopContainerRequest.$_createMessage);
  static StopContainerRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get timeoutSeconds => $_getIZ(1);
  @$pb.TagNumber(2)
  set timeoutSeconds($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTimeoutSeconds() => $_has(1);
  @$pb.TagNumber(2)
  void clearTimeoutSeconds() => $_clearField(2);
}

class StopContainerResponse extends $pb.GeneratedMessage {
  factory StopContainerResponse({
    Container? container,
  }) {
    final result = StopContainerResponse._();
    if (container != null) result.container = container;
    return result;
  }

  StopContainerResponse._();

  factory StopContainerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StopContainerResponse()..mergeFromBuffer(data, registry);
  factory StopContainerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StopContainerResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopContainerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StopContainerResponse.$_createMessage)
    ..aOM<Container>(1, _omitFieldNames ? '' : 'container',
        subBuilder: Container.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopContainerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopContainerResponse copyWith(
          void Function(StopContainerResponse) updates) =>
      super.copyWith((message) => updates(message as StopContainerResponse))
          as StopContainerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StopContainerResponse() / StopContainerResponse.new instead')
  static StopContainerResponse create() => StopContainerResponse._();
  static $pb.GeneratedMessage $_createMessage() => StopContainerResponse._();
  @$core.override
  StopContainerResponse createEmptyInstance() => StopContainerResponse._();
  @$core.pragma('dart2js:noInline')
  static StopContainerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopContainerResponse>(
          StopContainerResponse.$_createMessage);
  static StopContainerResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Container get container => $_getN(0);
  @$pb.TagNumber(1)
  set container(Container value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasContainer() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainer() => $_clearField(1);
  @$pb.TagNumber(1)
  Container ensureContainer() => $_ensure(0);
}

class RestartContainerRequest extends $pb.GeneratedMessage {
  factory RestartContainerRequest({
    $core.String? containerId,
  }) {
    final result = RestartContainerRequest._();
    if (containerId != null) result.containerId = containerId;
    return result;
  }

  RestartContainerRequest._();

  factory RestartContainerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RestartContainerRequest()..mergeFromBuffer(data, registry);
  factory RestartContainerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RestartContainerRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestartContainerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RestartContainerRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestartContainerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestartContainerRequest copyWith(
          void Function(RestartContainerRequest) updates) =>
      super.copyWith((message) => updates(message as RestartContainerRequest))
          as RestartContainerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RestartContainerRequest() / RestartContainerRequest.new instead')
  static RestartContainerRequest create() => RestartContainerRequest._();
  static $pb.GeneratedMessage $_createMessage() => RestartContainerRequest._();
  @$core.override
  RestartContainerRequest createEmptyInstance() => RestartContainerRequest._();
  @$core.pragma('dart2js:noInline')
  static RestartContainerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestartContainerRequest>(
          RestartContainerRequest.$_createMessage);
  static RestartContainerRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);
}

class RestartContainerResponse extends $pb.GeneratedMessage {
  factory RestartContainerResponse({
    Container? container,
  }) {
    final result = RestartContainerResponse._();
    if (container != null) result.container = container;
    return result;
  }

  RestartContainerResponse._();

  factory RestartContainerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RestartContainerResponse()..mergeFromBuffer(data, registry);
  factory RestartContainerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RestartContainerResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestartContainerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RestartContainerResponse.$_createMessage)
    ..aOM<Container>(1, _omitFieldNames ? '' : 'container',
        subBuilder: Container.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestartContainerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestartContainerResponse copyWith(
          void Function(RestartContainerResponse) updates) =>
      super.copyWith((message) => updates(message as RestartContainerResponse))
          as RestartContainerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RestartContainerResponse() / RestartContainerResponse.new instead')
  static RestartContainerResponse create() => RestartContainerResponse._();
  static $pb.GeneratedMessage $_createMessage() => RestartContainerResponse._();
  @$core.override
  RestartContainerResponse createEmptyInstance() =>
      RestartContainerResponse._();
  @$core.pragma('dart2js:noInline')
  static RestartContainerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestartContainerResponse>(
          RestartContainerResponse.$_createMessage);
  static RestartContainerResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Container get container => $_getN(0);
  @$pb.TagNumber(1)
  set container(Container value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasContainer() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainer() => $_clearField(1);
  @$pb.TagNumber(1)
  Container ensureContainer() => $_ensure(0);
}

class StreamContainerLogsRequest extends $pb.GeneratedMessage {
  factory StreamContainerLogsRequest({
    $core.String? containerId,
    $core.bool? follow,
    $core.int? tailLines,
    $1.Timestamp? since,
    LogStream? stream,
  }) {
    final result = StreamContainerLogsRequest._();
    if (containerId != null) result.containerId = containerId;
    if (follow != null) result.follow = follow;
    if (tailLines != null) result.tailLines = tailLines;
    if (since != null) result.since = since;
    if (stream != null) result.stream = stream;
    return result;
  }

  StreamContainerLogsRequest._();

  factory StreamContainerLogsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerLogsRequest()..mergeFromBuffer(data, registry);
  factory StreamContainerLogsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerLogsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamContainerLogsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamContainerLogsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..aOB(2, _omitFieldNames ? '' : 'follow')
    ..aI(3, _omitFieldNames ? '' : 'tailLines')
    ..aOM<$1.Timestamp>(4, _omitFieldNames ? '' : 'since',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aE<LogStream>(5, _omitFieldNames ? '' : 'stream',
        enumValues: LogStream.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerLogsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerLogsRequest copyWith(
          void Function(StreamContainerLogsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StreamContainerLogsRequest))
          as StreamContainerLogsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamContainerLogsRequest() / StreamContainerLogsRequest.new instead')
  static StreamContainerLogsRequest create() => StreamContainerLogsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamContainerLogsRequest._();
  @$core.override
  StreamContainerLogsRequest createEmptyInstance() =>
      StreamContainerLogsRequest._();
  @$core.pragma('dart2js:noInline')
  static StreamContainerLogsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamContainerLogsRequest>(
          StreamContainerLogsRequest.$_createMessage);
  static StreamContainerLogsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get follow => $_getBF(1);
  @$pb.TagNumber(2)
  set follow($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasFollow() => $_has(1);
  @$pb.TagNumber(2)
  void clearFollow() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get tailLines => $_getIZ(2);
  @$pb.TagNumber(3)
  set tailLines($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTailLines() => $_has(2);
  @$pb.TagNumber(3)
  void clearTailLines() => $_clearField(3);

  @$pb.TagNumber(4)
  $1.Timestamp get since => $_getN(3);
  @$pb.TagNumber(4)
  set since($1.Timestamp value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasSince() => $_has(3);
  @$pb.TagNumber(4)
  void clearSince() => $_clearField(4);
  @$pb.TagNumber(4)
  $1.Timestamp ensureSince() => $_ensure(3);

  @$pb.TagNumber(5)
  LogStream get stream => $_getN(4);
  @$pb.TagNumber(5)
  set stream(LogStream value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasStream() => $_has(4);
  @$pb.TagNumber(5)
  void clearStream() => $_clearField(5);
}

class StreamContainerLogsResponse extends $pb.GeneratedMessage {
  factory StreamContainerLogsResponse({
    $core.List<$core.int>? chunk,
    $1.Timestamp? emitted,
    LogStream? source,
  }) {
    final result = StreamContainerLogsResponse._();
    if (chunk != null) result.chunk = chunk;
    if (emitted != null) result.emitted = emitted;
    if (source != null) result.source = source;
    return result;
  }

  StreamContainerLogsResponse._();

  factory StreamContainerLogsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerLogsResponse()..mergeFromBuffer(data, registry);
  factory StreamContainerLogsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerLogsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamContainerLogsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamContainerLogsResponse.$_createMessage)
    ..a<$core.List<$core.int>>(
        1, _omitFieldNames ? '' : 'chunk', $pb.PbFieldType.OY)
    ..aOM<$1.Timestamp>(2, _omitFieldNames ? '' : 'emitted',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aE<LogStream>(3, _omitFieldNames ? '' : 'source',
        enumValues: LogStream.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerLogsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerLogsResponse copyWith(
          void Function(StreamContainerLogsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StreamContainerLogsResponse))
          as StreamContainerLogsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamContainerLogsResponse() / StreamContainerLogsResponse.new instead')
  static StreamContainerLogsResponse create() =>
      StreamContainerLogsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamContainerLogsResponse._();
  @$core.override
  StreamContainerLogsResponse createEmptyInstance() =>
      StreamContainerLogsResponse._();
  @$core.pragma('dart2js:noInline')
  static StreamContainerLogsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamContainerLogsResponse>(
          StreamContainerLogsResponse.$_createMessage);
  static StreamContainerLogsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get chunk => $_getN(0);
  @$pb.TagNumber(1)
  set chunk($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChunk() => $_has(0);
  @$pb.TagNumber(1)
  void clearChunk() => $_clearField(1);

  @$pb.TagNumber(2)
  $1.Timestamp get emitted => $_getN(1);
  @$pb.TagNumber(2)
  set emitted($1.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasEmitted() => $_has(1);
  @$pb.TagNumber(2)
  void clearEmitted() => $_clearField(2);
  @$pb.TagNumber(2)
  $1.Timestamp ensureEmitted() => $_ensure(1);

  @$pb.TagNumber(3)
  LogStream get source => $_getN(2);
  @$pb.TagNumber(3)
  set source(LogStream value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasSource() => $_has(2);
  @$pb.TagNumber(3)
  void clearSource() => $_clearField(3);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
