// This is a generated file - do not edit.
//
// Generated from forgenet/v1/update.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// UpdateStep est l'étape où se trouve, ou s'est arrêtée, une mise à jour
/// (ADR-0011).
class UpdateStep extends $pb.ProtobufEnum {
  static const UpdateStep UPDATE_STEP_UNSPECIFIED =
      UpdateStep._(0, _omitEnumNames ? '' : 'UPDATE_STEP_UNSPECIFIED');
  static const UpdateStep UPDATE_STEP_CHECK =
      UpdateStep._(1, _omitEnumNames ? '' : 'UPDATE_STEP_CHECK');
  static const UpdateStep UPDATE_STEP_PULL =
      UpdateStep._(2, _omitEnumNames ? '' : 'UPDATE_STEP_PULL');
  static const UpdateStep UPDATE_STEP_BACKUP =
      UpdateStep._(3, _omitEnumNames ? '' : 'UPDATE_STEP_BACKUP');
  static const UpdateStep UPDATE_STEP_SWITCH =
      UpdateStep._(4, _omitEnumNames ? '' : 'UPDATE_STEP_SWITCH');
  static const UpdateStep UPDATE_STEP_HEALTH =
      UpdateStep._(5, _omitEnumNames ? '' : 'UPDATE_STEP_HEALTH');
  static const UpdateStep UPDATE_STEP_ROLLBACK =
      UpdateStep._(6, _omitEnumNames ? '' : 'UPDATE_STEP_ROLLBACK');
  static const UpdateStep UPDATE_STEP_DONE =
      UpdateStep._(7, _omitEnumNames ? '' : 'UPDATE_STEP_DONE');

  static const $core.List<UpdateStep> values = <UpdateStep>[
    UPDATE_STEP_UNSPECIFIED,
    UPDATE_STEP_CHECK,
    UPDATE_STEP_PULL,
    UPDATE_STEP_BACKUP,
    UPDATE_STEP_SWITCH,
    UPDATE_STEP_HEALTH,
    UPDATE_STEP_ROLLBACK,
    UPDATE_STEP_DONE,
  ];

  static final $core.List<UpdateStep?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static UpdateStep? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UpdateStep._(super.value, super.name);
}

/// UpdateResult est l'issue d'une mise à jour.
class UpdateResult extends $pb.ProtobufEnum {
  static const UpdateResult UPDATE_RESULT_UNSPECIFIED =
      UpdateResult._(0, _omitEnumNames ? '' : 'UPDATE_RESULT_UNSPECIFIED');
  static const UpdateResult UPDATE_RESULT_RUNNING =
      UpdateResult._(1, _omitEnumNames ? '' : 'UPDATE_RESULT_RUNNING');
  static const UpdateResult UPDATE_RESULT_SUCCEEDED =
      UpdateResult._(2, _omitEnumNames ? '' : 'UPDATE_RESULT_SUCCEEDED');
  static const UpdateResult UPDATE_RESULT_ROLLED_BACK =
      UpdateResult._(3, _omitEnumNames ? '' : 'UPDATE_RESULT_ROLLED_BACK');
  static const UpdateResult UPDATE_RESULT_FAILED =
      UpdateResult._(4, _omitEnumNames ? '' : 'UPDATE_RESULT_FAILED');

  static const $core.List<UpdateResult> values = <UpdateResult>[
    UPDATE_RESULT_UNSPECIFIED,
    UPDATE_RESULT_RUNNING,
    UPDATE_RESULT_SUCCEEDED,
    UPDATE_RESULT_ROLLED_BACK,
    UPDATE_RESULT_FAILED,
  ];

  static final $core.List<UpdateResult?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static UpdateResult? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UpdateResult._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
