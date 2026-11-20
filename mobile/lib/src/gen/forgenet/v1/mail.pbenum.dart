// This is a generated file - do not edit.
//
// Generated from forgenet/v1/mail.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// MailTLS est la protection de la liaison avec le relais (ADR-0014).
class MailTLS extends $pb.ProtobufEnum {
  static const MailTLS MAIL_TLS_UNSPECIFIED =
      MailTLS._(0, _omitEnumNames ? '' : 'MAIL_TLS_UNSPECIFIED');
  static const MailTLS MAIL_TLS_STARTTLS =
      MailTLS._(1, _omitEnumNames ? '' : 'MAIL_TLS_STARTTLS');
  static const MailTLS MAIL_TLS_IMPLICIT =
      MailTLS._(2, _omitEnumNames ? '' : 'MAIL_TLS_IMPLICIT');
  static const MailTLS MAIL_TLS_NONE =
      MailTLS._(3, _omitEnumNames ? '' : 'MAIL_TLS_NONE');

  static const $core.List<MailTLS> values = <MailTLS>[
    MAIL_TLS_UNSPECIFIED,
    MAIL_TLS_STARTTLS,
    MAIL_TLS_IMPLICIT,
    MAIL_TLS_NONE,
  ];

  static final $core.List<MailTLS?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MailTLS? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MailTLS._(super.value, super.name);
}

/// MailSource dit d'où viennent les réglages en vigueur (ADR-0015).
class MailSource extends $pb.ProtobufEnum {
  static const MailSource MAIL_SOURCE_UNSPECIFIED =
      MailSource._(0, _omitEnumNames ? '' : 'MAIL_SOURCE_UNSPECIFIED');
  static const MailSource MAIL_SOURCE_SERVER =
      MailSource._(1, _omitEnumNames ? '' : 'MAIL_SOURCE_SERVER');
  static const MailSource MAIL_SOURCE_DASHBOARD =
      MailSource._(2, _omitEnumNames ? '' : 'MAIL_SOURCE_DASHBOARD');

  static const $core.List<MailSource> values = <MailSource>[
    MAIL_SOURCE_UNSPECIFIED,
    MAIL_SOURCE_SERVER,
    MAIL_SOURCE_DASHBOARD,
  ];

  static final $core.List<MailSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MailSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MailSource._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
