// This is a generated file - do not edit.
//
// Generated from forgenet/v1/deploy.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// DeploymentState représente la machine à états explicite du cycle de vie d'un déploiement.
/// Transitions autorisées :
/// PENDING -> BUILDING, CANCELLED
/// BUILDING -> DEPLOYING, FAILED, CANCELLED
/// DEPLOYING -> HEALTH_CHECKING, FAILED, CANCELLED
/// HEALTH_CHECKING -> SUCCESS, FAILED, ROLLING_BACK
/// ROLLING_BACK -> ROLLED_BACK, FAILED
class DeploymentState extends $pb.ProtobufEnum {
  static const DeploymentState DEPLOYMENT_STATE_UNSPECIFIED = DeploymentState._(
      0, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_UNSPECIFIED');
  static const DeploymentState DEPLOYMENT_STATE_PENDING =
      DeploymentState._(1, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_PENDING');
  static const DeploymentState DEPLOYMENT_STATE_BUILDING =
      DeploymentState._(2, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_BUILDING');
  static const DeploymentState DEPLOYMENT_STATE_DEPLOYING =
      DeploymentState._(3, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_DEPLOYING');
  static const DeploymentState DEPLOYMENT_STATE_HEALTH_CHECKING =
      DeploymentState._(
          4, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_HEALTH_CHECKING');
  static const DeploymentState DEPLOYMENT_STATE_SUCCESS =
      DeploymentState._(5, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_SUCCESS');
  static const DeploymentState DEPLOYMENT_STATE_FAILED =
      DeploymentState._(6, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_FAILED');
  static const DeploymentState DEPLOYMENT_STATE_CANCELLED =
      DeploymentState._(7, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_CANCELLED');
  static const DeploymentState DEPLOYMENT_STATE_ROLLING_BACK =
      DeploymentState._(
          8, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_ROLLING_BACK');
  static const DeploymentState DEPLOYMENT_STATE_ROLLED_BACK = DeploymentState._(
      9, _omitEnumNames ? '' : 'DEPLOYMENT_STATE_ROLLED_BACK');

  static const $core.List<DeploymentState> values = <DeploymentState>[
    DEPLOYMENT_STATE_UNSPECIFIED,
    DEPLOYMENT_STATE_PENDING,
    DEPLOYMENT_STATE_BUILDING,
    DEPLOYMENT_STATE_DEPLOYING,
    DEPLOYMENT_STATE_HEALTH_CHECKING,
    DEPLOYMENT_STATE_SUCCESS,
    DEPLOYMENT_STATE_FAILED,
    DEPLOYMENT_STATE_CANCELLED,
    DEPLOYMENT_STATE_ROLLING_BACK,
    DEPLOYMENT_STATE_ROLLED_BACK,
  ];

  static final $core.List<DeploymentState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static DeploymentState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeploymentState._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
