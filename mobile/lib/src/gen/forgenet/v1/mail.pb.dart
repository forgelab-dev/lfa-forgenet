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
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

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

/// GetMailStatusResponse décrit l'envoi d'e-mails. Le mot de passe du relais
/// ne quitte jamais le serveur : seule sa présence est dite (authenticated).
class GetMailStatusResponse extends $pb.GeneratedMessage {
  factory GetMailStatusResponse({
    $core.bool? enabled,
    $core.String? host,
    $core.int? port,
    MailTLS? tls,
    $core.String? from,
    $core.bool? authenticated,
    MailSource? source,
    $core.bool? editable,
    $core.String? username,
    $core.bool? disabled,
    $core.bool? passwordUnreadable,
    $1.Timestamp? updatedAt,
    $core.String? updatedBy,
  }) {
    final result = GetMailStatusResponse._();
    if (enabled != null) result.enabled = enabled;
    if (host != null) result.host = host;
    if (port != null) result.port = port;
    if (tls != null) result.tls = tls;
    if (from != null) result.from = from;
    if (authenticated != null) result.authenticated = authenticated;
    if (source != null) result.source = source;
    if (editable != null) result.editable = editable;
    if (username != null) result.username = username;
    if (disabled != null) result.disabled = disabled;
    if (passwordUnreadable != null)
      result.passwordUnreadable = passwordUnreadable;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (updatedBy != null) result.updatedBy = updatedBy;
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
    ..aE<MailSource>(7, _omitFieldNames ? '' : 'source',
        enumValues: MailSource.values)
    ..aOB(8, _omitFieldNames ? '' : 'editable')
    ..aOS(9, _omitFieldNames ? '' : 'username')
    ..aOB(10, _omitFieldNames ? '' : 'disabled')
    ..aOB(11, _omitFieldNames ? '' : 'passwordUnreadable')
    ..aOM<$1.Timestamp>(12, _omitFieldNames ? '' : 'updatedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOS(13, _omitFieldNames ? '' : 'updatedBy')
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

  /// Faux : aucun relais n'est réglé, ou l'envoi est coupé ; ForgeNet
  /// n'envoie rien.
  @$pb.TagNumber(1)
  $core.bool get enabled => $_getBF(0);
  @$pb.TagNumber(1)
  set enabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasEnabled() => $_has(0);
  @$pb.TagNumber(1)
  void clearEnabled() => $_clearField(1);

  /// Les valeurs restent données quand l'envoi est coupé, pour le rouvrir.
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

  /// Vrai : le serveur s'authentifie auprès du relais (un mot de passe est
  /// enregistré).
  @$pb.TagNumber(6)
  $core.bool get authenticated => $_getBF(5);
  @$pb.TagNumber(6)
  set authenticated($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasAuthenticated() => $_has(5);
  @$pb.TagNumber(6)
  void clearAuthenticated() => $_clearField(6);

  @$pb.TagNumber(7)
  MailSource get source => $_getN(6);
  @$pb.TagNumber(7)
  set source(MailSource value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasSource() => $_has(6);
  @$pb.TagNumber(7)
  void clearSource() => $_clearField(7);

  /// Vrai : les réglages se changent depuis le dashboard (le serveur a sa clé
  /// FORGENET_SECRET_KEY).
  @$pb.TagNumber(8)
  $core.bool get editable => $_getBF(7);
  @$pb.TagNumber(8)
  set editable($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(8)
  $core.bool hasEditable() => $_has(7);
  @$pb.TagNumber(8)
  void clearEditable() => $_clearField(8);

  /// Identifiant auprès du relais, pour le formulaire.
  @$pb.TagNumber(9)
  $core.String get username => $_getSZ(8);
  @$pb.TagNumber(9)
  set username($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasUsername() => $_has(8);
  @$pb.TagNumber(9)
  void clearUsername() => $_clearField(9);

  /// Vrai : l'envoi a été coupé depuis le dashboard.
  @$pb.TagNumber(10)
  $core.bool get disabled => $_getBF(9);
  @$pb.TagNumber(10)
  set disabled($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(10)
  $core.bool hasDisabled() => $_has(9);
  @$pb.TagNumber(10)
  void clearDisabled() => $_clearField(10);

  /// Vrai : des réglages du dashboard existent, mais leur mot de passe ne
  /// s'ouvre plus avec la clé du serveur ; ceux de .env s'appliquent.
  @$pb.TagNumber(11)
  $core.bool get passwordUnreadable => $_getBF(10);
  @$pb.TagNumber(11)
  set passwordUnreadable($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(11)
  $core.bool hasPasswordUnreadable() => $_has(10);
  @$pb.TagNumber(11)
  void clearPasswordUnreadable() => $_clearField(11);

  /// Dernier enregistrement depuis le dashboard : quand et par quel compte.
  @$pb.TagNumber(12)
  $1.Timestamp get updatedAt => $_getN(11);
  @$pb.TagNumber(12)
  set updatedAt($1.Timestamp value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasUpdatedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearUpdatedAt() => $_clearField(12);
  @$pb.TagNumber(12)
  $1.Timestamp ensureUpdatedAt() => $_ensure(11);

  @$pb.TagNumber(13)
  $core.String get updatedBy => $_getSZ(12);
  @$pb.TagNumber(13)
  set updatedBy($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasUpdatedBy() => $_has(12);
  @$pb.TagNumber(13)
  void clearUpdatedBy() => $_clearField(13);
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

/// SaveMailSettingsRequest porte les nouveaux réglages du relais.
class SaveMailSettingsRequest extends $pb.GeneratedMessage {
  factory SaveMailSettingsRequest({
    $core.String? host,
    $core.int? port,
    MailTLS? tls,
    $core.String? username,
    $core.String? from,
    $core.bool? keepPassword,
    $core.String? password,
  }) {
    final result = SaveMailSettingsRequest._();
    if (host != null) result.host = host;
    if (port != null) result.port = port;
    if (tls != null) result.tls = tls;
    if (username != null) result.username = username;
    if (from != null) result.from = from;
    if (keepPassword != null) result.keepPassword = keepPassword;
    if (password != null) result.password = password;
    return result;
  }

  SaveMailSettingsRequest._();

  factory SaveMailSettingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SaveMailSettingsRequest()..mergeFromBuffer(data, registry);
  factory SaveMailSettingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SaveMailSettingsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SaveMailSettingsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SaveMailSettingsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'host')
    ..aI(2, _omitFieldNames ? '' : 'port', fieldType: $pb.PbFieldType.OU3)
    ..aE<MailTLS>(3, _omitFieldNames ? '' : 'tls', enumValues: MailTLS.values)
    ..aOS(4, _omitFieldNames ? '' : 'username')
    ..aOS(5, _omitFieldNames ? '' : 'from')
    ..aOB(6, _omitFieldNames ? '' : 'keepPassword')
    ..aOS(7, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SaveMailSettingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SaveMailSettingsRequest copyWith(
          void Function(SaveMailSettingsRequest) updates) =>
      super.copyWith((message) => updates(message as SaveMailSettingsRequest))
          as SaveMailSettingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SaveMailSettingsRequest() / SaveMailSettingsRequest.new instead')
  static SaveMailSettingsRequest create() => SaveMailSettingsRequest._();
  static $pb.GeneratedMessage $_createMessage() => SaveMailSettingsRequest._();
  @$core.override
  SaveMailSettingsRequest createEmptyInstance() => SaveMailSettingsRequest._();
  @$core.pragma('dart2js:noInline')
  static SaveMailSettingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SaveMailSettingsRequest>(
          SaveMailSettingsRequest.$_createMessage);
  static SaveMailSettingsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get host => $_getSZ(0);
  @$pb.TagNumber(1)
  set host($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHost() => $_has(0);
  @$pb.TagNumber(1)
  void clearHost() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get port => $_getIZ(1);
  @$pb.TagNumber(2)
  set port($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPort() => $_has(1);
  @$pb.TagNumber(2)
  void clearPort() => $_clearField(2);

  /// MAIL_TLS_UNSPECIFIED vaut STARTTLS.
  @$pb.TagNumber(3)
  MailTLS get tls => $_getN(2);
  @$pb.TagNumber(3)
  set tls(MailTLS value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasTls() => $_has(2);
  @$pb.TagNumber(3)
  void clearTls() => $_clearField(3);

  /// Vide : le relais est utilisé sans authentification.
  @$pb.TagNumber(4)
  $core.String get username => $_getSZ(3);
  @$pb.TagNumber(4)
  set username($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasUsername() => $_has(3);
  @$pb.TagNumber(4)
  void clearUsername() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get from => $_getSZ(4);
  @$pb.TagNumber(5)
  set from($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasFrom() => $_has(4);
  @$pb.TagNumber(5)
  void clearFrom() => $_clearField(5);

  /// Vrai : le mot de passe enregistré est gardé, password est ignoré.
  @$pb.TagNumber(6)
  $core.bool get keepPassword => $_getBF(5);
  @$pb.TagNumber(6)
  set keepPassword($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasKeepPassword() => $_has(5);
  @$pb.TagNumber(6)
  void clearKeepPassword() => $_clearField(6);

  /// Nouveau mot de passe, quand keep_password est faux. Jamais renvoyé.
  @$pb.TagNumber(7)
  $core.String get password => $_getSZ(6);
  @$pb.TagNumber(7)
  set password($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasPassword() => $_has(6);
  @$pb.TagNumber(7)
  void clearPassword() => $_clearField(7);
}

class SaveMailSettingsResponse extends $pb.GeneratedMessage {
  factory SaveMailSettingsResponse({
    GetMailStatusResponse? status,
    $core.String? to,
  }) {
    final result = SaveMailSettingsResponse._();
    if (status != null) result.status = status;
    if (to != null) result.to = to;
    return result;
  }

  SaveMailSettingsResponse._();

  factory SaveMailSettingsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SaveMailSettingsResponse()..mergeFromBuffer(data, registry);
  factory SaveMailSettingsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SaveMailSettingsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SaveMailSettingsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SaveMailSettingsResponse.$_createMessage)
    ..aOM<GetMailStatusResponse>(1, _omitFieldNames ? '' : 'status',
        subBuilder: GetMailStatusResponse.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'to')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SaveMailSettingsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SaveMailSettingsResponse copyWith(
          void Function(SaveMailSettingsResponse) updates) =>
      super.copyWith((message) => updates(message as SaveMailSettingsResponse))
          as SaveMailSettingsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use SaveMailSettingsResponse() / SaveMailSettingsResponse.new instead')
  static SaveMailSettingsResponse create() => SaveMailSettingsResponse._();
  static $pb.GeneratedMessage $_createMessage() => SaveMailSettingsResponse._();
  @$core.override
  SaveMailSettingsResponse createEmptyInstance() =>
      SaveMailSettingsResponse._();
  @$core.pragma('dart2js:noInline')
  static SaveMailSettingsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SaveMailSettingsResponse>(
          SaveMailSettingsResponse.$_createMessage);
  static SaveMailSettingsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  GetMailStatusResponse get status => $_getN(0);
  @$pb.TagNumber(1)
  set status(GetMailStatusResponse value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);
  @$pb.TagNumber(1)
  GetMailStatusResponse ensureStatus() => $_ensure(0);

  /// Adresse du compte connecté, à qui l'e-mail d'essai est parti.
  @$pb.TagNumber(2)
  $core.String get to => $_getSZ(1);
  @$pb.TagNumber(2)
  set to($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTo() => $_has(1);
  @$pb.TagNumber(2)
  void clearTo() => $_clearField(2);
}

class DisableMailRequest extends $pb.GeneratedMessage {
  factory DisableMailRequest() => DisableMailRequest._();

  DisableMailRequest._();

  factory DisableMailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DisableMailRequest()..mergeFromBuffer(data, registry);
  factory DisableMailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DisableMailRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableMailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DisableMailRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableMailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableMailRequest copyWith(void Function(DisableMailRequest) updates) =>
      super.copyWith((message) => updates(message as DisableMailRequest))
          as DisableMailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DisableMailRequest() / DisableMailRequest.new instead')
  static DisableMailRequest create() => DisableMailRequest._();
  static $pb.GeneratedMessage $_createMessage() => DisableMailRequest._();
  @$core.override
  DisableMailRequest createEmptyInstance() => DisableMailRequest._();
  @$core.pragma('dart2js:noInline')
  static DisableMailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableMailRequest>(
          DisableMailRequest.$_createMessage);
  static DisableMailRequest? _defaultInstance;
}

class DisableMailResponse extends $pb.GeneratedMessage {
  factory DisableMailResponse({
    GetMailStatusResponse? status,
  }) {
    final result = DisableMailResponse._();
    if (status != null) result.status = status;
    return result;
  }

  DisableMailResponse._();

  factory DisableMailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DisableMailResponse()..mergeFromBuffer(data, registry);
  factory DisableMailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DisableMailResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableMailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DisableMailResponse.$_createMessage)
    ..aOM<GetMailStatusResponse>(1, _omitFieldNames ? '' : 'status',
        subBuilder: GetMailStatusResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableMailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableMailResponse copyWith(void Function(DisableMailResponse) updates) =>
      super.copyWith((message) => updates(message as DisableMailResponse))
          as DisableMailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use DisableMailResponse() / DisableMailResponse.new instead')
  static DisableMailResponse create() => DisableMailResponse._();
  static $pb.GeneratedMessage $_createMessage() => DisableMailResponse._();
  @$core.override
  DisableMailResponse createEmptyInstance() => DisableMailResponse._();
  @$core.pragma('dart2js:noInline')
  static DisableMailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableMailResponse>(
          DisableMailResponse.$_createMessage);
  static DisableMailResponse? _defaultInstance;

  @$pb.TagNumber(1)
  GetMailStatusResponse get status => $_getN(0);
  @$pb.TagNumber(1)
  set status(GetMailStatusResponse value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);
  @$pb.TagNumber(1)
  GetMailStatusResponse ensureStatus() => $_ensure(0);
}

class ResetMailSettingsRequest extends $pb.GeneratedMessage {
  factory ResetMailSettingsRequest() => ResetMailSettingsRequest._();

  ResetMailSettingsRequest._();

  factory ResetMailSettingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResetMailSettingsRequest()..mergeFromBuffer(data, registry);
  factory ResetMailSettingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResetMailSettingsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResetMailSettingsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ResetMailSettingsRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetMailSettingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetMailSettingsRequest copyWith(
          void Function(ResetMailSettingsRequest) updates) =>
      super.copyWith((message) => updates(message as ResetMailSettingsRequest))
          as ResetMailSettingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ResetMailSettingsRequest() / ResetMailSettingsRequest.new instead')
  static ResetMailSettingsRequest create() => ResetMailSettingsRequest._();
  static $pb.GeneratedMessage $_createMessage() => ResetMailSettingsRequest._();
  @$core.override
  ResetMailSettingsRequest createEmptyInstance() =>
      ResetMailSettingsRequest._();
  @$core.pragma('dart2js:noInline')
  static ResetMailSettingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResetMailSettingsRequest>(
          ResetMailSettingsRequest.$_createMessage);
  static ResetMailSettingsRequest? _defaultInstance;
}

class ResetMailSettingsResponse extends $pb.GeneratedMessage {
  factory ResetMailSettingsResponse({
    GetMailStatusResponse? status,
  }) {
    final result = ResetMailSettingsResponse._();
    if (status != null) result.status = status;
    return result;
  }

  ResetMailSettingsResponse._();

  factory ResetMailSettingsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResetMailSettingsResponse()..mergeFromBuffer(data, registry);
  factory ResetMailSettingsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ResetMailSettingsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResetMailSettingsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ResetMailSettingsResponse.$_createMessage)
    ..aOM<GetMailStatusResponse>(1, _omitFieldNames ? '' : 'status',
        subBuilder: GetMailStatusResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetMailSettingsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetMailSettingsResponse copyWith(
          void Function(ResetMailSettingsResponse) updates) =>
      super.copyWith((message) => updates(message as ResetMailSettingsResponse))
          as ResetMailSettingsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ResetMailSettingsResponse() / ResetMailSettingsResponse.new instead')
  static ResetMailSettingsResponse create() => ResetMailSettingsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      ResetMailSettingsResponse._();
  @$core.override
  ResetMailSettingsResponse createEmptyInstance() =>
      ResetMailSettingsResponse._();
  @$core.pragma('dart2js:noInline')
  static ResetMailSettingsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResetMailSettingsResponse>(
          ResetMailSettingsResponse.$_createMessage);
  static ResetMailSettingsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  GetMailStatusResponse get status => $_getN(0);
  @$pb.TagNumber(1)
  set status(GetMailStatusResponse value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);
  @$pb.TagNumber(1)
  GetMailStatusResponse ensureStatus() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
