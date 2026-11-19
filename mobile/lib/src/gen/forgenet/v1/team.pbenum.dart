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

import 'package:protobuf/protobuf.dart' as $pb;

/// MemberStatus dit où en est un compte (ADR-0014).
class MemberStatus extends $pb.ProtobufEnum {
  static const MemberStatus MEMBER_STATUS_UNSPECIFIED =
      MemberStatus._(0, _omitEnumNames ? '' : 'MEMBER_STATUS_UNSPECIFIED');

  /// Le compte a son mot de passe et peut se connecter.
  static const MemberStatus MEMBER_STATUS_ACTIVE =
      MemberStatus._(1, _omitEnumNames ? '' : 'MEMBER_STATUS_ACTIVE');

  /// Le compte attend son mot de passe, et une invitation valable court.
  static const MemberStatus MEMBER_STATUS_INVITED =
      MemberStatus._(2, _omitEnumNames ? '' : 'MEMBER_STATUS_INVITED');

  /// Le compte attend son mot de passe, mais aucune invitation valable ne
  /// court : elle a expiré, ou le compte a été créé sans elle.
  static const MemberStatus MEMBER_STATUS_INVITATION_EXPIRED = MemberStatus._(
      3, _omitEnumNames ? '' : 'MEMBER_STATUS_INVITATION_EXPIRED');

  /// Le compte est suspendu : il ne peut plus se connecter, et aucun de ses
  /// jetons n'est accepté. Un administrateur peut le réactiver.
  static const MemberStatus MEMBER_STATUS_SUSPENDED =
      MemberStatus._(4, _omitEnumNames ? '' : 'MEMBER_STATUS_SUSPENDED');

  static const $core.List<MemberStatus> values = <MemberStatus>[
    MEMBER_STATUS_UNSPECIFIED,
    MEMBER_STATUS_ACTIVE,
    MEMBER_STATUS_INVITED,
    MEMBER_STATUS_INVITATION_EXPIRED,
    MEMBER_STATUS_SUSPENDED,
  ];

  static final $core.List<MemberStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static MemberStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MemberStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
