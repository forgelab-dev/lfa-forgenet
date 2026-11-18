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

import 'mail.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'mail.pbenum.dart';

class GetMailStatusRequest extends $pb.GeneratedMessage {
  factory GetMailStatusRequest() => GetMailStatusRequest._();

  GetMailStatusRequest._();

  factory GetMailStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetMailStatusRequest()..mergeFromBuffer(data, registry);
  factory GetMailStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetMailStatusRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMailStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetMailStatusRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMailStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMailStatusRequest copyWith(void Function(GetMailStatusRequest) updates) =>
      super.copyWith((message) => updates(message as GetMailStatusRequest))
          as GetMailStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetMailStatusRequest() / GetMailStatusRequest.new instead')
  static GetMailStatusRequest create() => GetMailStatusRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetMailStatusRequest._();
  @$core.override
  GetMailStatusRequest createEmptyInstance() => GetMailStatusRequest._();
  @$core.pragma('dart2js:noInline')
  static GetMailStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMailStatusRequest>(
          GetMailStatusRequest.$_createMessage);
  static GetMailStatusRequest? _defaultInstance;
}

/// GetMailStatusResponse décrit l'envoi d'e-mails. Les identifiants du relais
/// ne quittent jamais le serveur : seule leur présence est dite.
class GetMailStatusResponse extends $pb.GeneratedMessage {
  factory GetMailStatusResponse({
    $core.bool? enabled,
    $core.String? host,
    $core.int? port,
    MailTLS? tls,
    $core.String? from,
    $core.bool? authenticated,
  }) {
    final result = GetMailStatusResponse._();
    if (enabled != null) result.enabled = enabled;
    if (host != null) result.host = host;
    if (port != null) result.port = port;
    if (tls != null) result.tls = tls;
    if (from != null) result.from = from;
    if (authenticated != null) result.authenticated = authenticated;
    return result;
  }

  GetMailStatusResponse._();

  factory GetMailStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetMailStatusResponse()..mergeFromBuffer(data, registry);
  factory GetMailStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetMailStatusResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMailStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetMailStatusResponse.$_createMessage)
    ..aOB(1, _omitFieldNames ? '' : 'enabled')
    ..aOS(2, _omitFieldNames ? '' : 'host')
    ..aI(3, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aE<MailTLS>(4, _omitFieldNames ? '' : 'tls', enumValues: MailTLS.values)
    ..aOS(5, _omitFieldNames ? '' : 'from')
    ..aOB(6, _omitFieldNames ? '' : 'authenticated')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMailStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMailStatusResponse copyWith(
          void Function(GetMailStatusResponse) updates) =>
      super.copyWith((message) => updates(message as GetMailStatusResponse))
          as GetMailStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetMailStatusResponse() / GetMailStatusResponse.new instead')
  static GetMailStatusResponse create() => GetMailStatusResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetMailStatusResponse._();
  @$core.override
  GetMailStatusResponse createEmptyInstance() => GetMailStatusResponse._();
  @$core.pragma('dart2js:noInline')
  static GetMailStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMailStatusResponse>(
          GetMailStatusResponse.$_createMessage);
  static GetMailStatusResponse? _defaultInstance;

  /// Faux : aucun relais n'est configuré, ForgeNet n'envoie rien.
  @$pb.TagNumber(1)
  $core.bool get enabled => $_getBF(0);
  @$pb.TagNumber(1)
  set enabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasEnabled() => $_has(0);
  @$pb.TagNumber(1)
  void clearEnabled() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get host => $_getSZ(1);
  @$pb.TagNumber(2)
  set host($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasHost() => $_has(1);
  @$pb.TagNumber(2)
  void clearHost() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get port => $_getIZ(2);
  @$pb.TagNumber(3)
  set port($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPort() => $_has(2);
  @$pb.TagNumber(3)
  void clearPort() => $_clearField(3);

  @$pb.TagNumber(4)
  MailTLS get tls => $_getN(3);
  @$pb.TagNumber(4)
  set tls(MailTLS value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasTls() => $_has(3);
  @$pb.TagNumber(4)
  void clearTls() => $_clearField(4);

  /// Expéditeur des messages, tel qu'il apparaît au destinataire.
  @$pb.TagNumber(5)
  $core.String get from => $_getSZ(4);
  @$pb.TagNumber(5)
  set from($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasFrom() => $_has(4);
  @$pb.TagNumber(5)
  void clearFrom() => $_clearField(5);

  /// Vrai : le serveur s'authentifie auprès du relais.
  @$pb.TagNumber(6)
  $core.bool get authenticated => $_getBF(5);
  @$pb.TagNumber(6)
  set authenticated($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasAuthenticated() => $_has(5);
  @$pb.TagNumber(6)
  void clearAuthenticated() => $_clearField(6);
}

class SendTestMailRequest extends $pb.GeneratedMessage {
  factory SendTestMailRequest() => SendTestMailRequest._();

  SendTestMailRequest._();

  factory SendTestMailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SendTestMailRequest()..mergeFromBuffer(data, registry);
  factory SendTestMailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SendTestMailRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTestMailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SendTestMailRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTestMailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTestMailRequest copyWith(void Function(SendTestMailRequest) updates) =>
      super.copyWith((message) => updates(message as SendTestMailRequest))
          as SendTestMailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use SendTestMailRequest() / SendTestMailRequest.new instead')
  static SendTestMailRequest create() => SendTestMailRequest._();
  static $pb.GeneratedMessage $_createMessage() => SendTestMailRequest._();
  @$core.override
  SendTestMailRequest createEmptyInstance() => SendTestMailRequest._();
  @$core.pragma('dart2js:noInline')
  static SendTestMailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTestMailRequest>(
          SendTestMailRequest.$_createMessage);
  static SendTestMailRequest? _defaultInstance;
}

class SendTestMailResponse extends $pb.GeneratedMessage {
  factory SendTestMailResponse({
    $core.String? to,
  }) {
    final result = SendTestMailResponse._();
    if (to != null) result.to = to;
    return result;
  }

  SendTestMailResponse._();

  factory SendTestMailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SendTestMailResponse()..mergeFromBuffer(data, registry);
  factory SendTestMailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SendTestMailResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTestMailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SendTestMailResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'to')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTestMailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTestMailResponse copyWith(void Function(SendTestMailResponse) updates) =>
      super.copyWith((message) => updates(message as SendTestMailResponse))
          as SendTestMailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SendTestMailResponse() / SendTestMailResponse.new instead')
  static SendTestMailResponse create() => SendTestMailResponse._();
  static $pb.GeneratedMessage $_createMessage() => SendTestMailResponse._();
  @$core.override
  SendTestMailResponse createEmptyInstance() => SendTestMailResponse._();
  @$core.pragma('dart2js:noInline')
  static SendTestMailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTestMailResponse>(
          SendTestMailResponse.$_createMessage);
  static SendTestMailResponse? _defaultInstance;

  /// Adresse du compte connecté, à qui l'e-mail d'essai est parti.
  @$pb.TagNumber(1)
  $core.String get to => $_getSZ(0);
  @$pb.TagNumber(1)
  set to($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTo() => $_has(0);
  @$pb.TagNumber(1)
  void clearTo() => $_clearField(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
