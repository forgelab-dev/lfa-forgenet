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

/// ContainerState représente les différents états possibles du cycle de vie d'un conteneur.
class ContainerState extends $pb.ProtobufEnum {
  static const ContainerState CONTAINER_STATE_UNSPECIFIED =
      ContainerState._(0, _omitEnumNames ? '' : 'CONTAINER_STATE_UNSPECIFIED');
  static const ContainerState CONTAINER_STATE_CREATED =
      ContainerState._(1, _omitEnumNames ? '' : 'CONTAINER_STATE_CREATED');
  static const ContainerState CONTAINER_STATE_RUNNING =
      ContainerState._(2, _omitEnumNames ? '' : 'CONTAINER_STATE_RUNNING');
  static const ContainerState CONTAINER_STATE_PAUSED =
      ContainerState._(3, _omitEnumNames ? '' : 'CONTAINER_STATE_PAUSED');
  static const ContainerState CONTAINER_STATE_RESTARTING =
      ContainerState._(4, _omitEnumNames ? '' : 'CONTAINER_STATE_RESTARTING');
  static const ContainerState CONTAINER_STATE_EXITED =
      ContainerState._(5, _omitEnumNames ? '' : 'CONTAINER_STATE_EXITED');
  static const ContainerState CONTAINER_STATE_DEAD =
      ContainerState._(6, _omitEnumNames ? '' : 'CONTAINER_STATE_DEAD');

  static const $core.List<ContainerState> values = <ContainerState>[
    CONTAINER_STATE_UNSPECIFIED,
    CONTAINER_STATE_CREATED,
    CONTAINER_STATE_RUNNING,
    CONTAINER_STATE_PAUSED,
    CONTAINER_STATE_RESTARTING,
    CONTAINER_STATE_EXITED,
    CONTAINER_STATE_DEAD,
  ];

  static final $core.List<ContainerState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ContainerState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContainerState._(super.value, super.name);
}

/// LogStream identifie le flux de sortie du conteneur (stderr ou stdout)
class LogStream extends $pb.ProtobufEnum {
  static const LogStream LOG_STREAM_UNSPECIFIED =
      LogStream._(0, _omitEnumNames ? '' : 'LOG_STREAM_UNSPECIFIED');
  static const LogStream LOG_STREAM_STDERR =
      LogStream._(1, _omitEnumNames ? '' : 'LOG_STREAM_STDERR');
  static const LogStream LOG_STREAM_STDOUT =
      LogStream._(2, _omitEnumNames ? '' : 'LOG_STREAM_STDOUT');

  static const $core.List<LogStream> values = <LogStream>[
    LOG_STREAM_UNSPECIFIED,
    LOG_STREAM_STDERR,
    LOG_STREAM_STDOUT,
  ];

  static final $core.List<LogStream?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LogStream? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogStream._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
