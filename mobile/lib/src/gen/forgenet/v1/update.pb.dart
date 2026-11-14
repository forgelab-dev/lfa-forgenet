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
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

import 'update.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'update.pbenum.dart';

/// AvailableRelease est une version publiée plus récente que celle en service.
class AvailableRelease extends $pb.GeneratedMessage {
  factory AvailableRelease({
    $core.String? version,
    $core.String? notesUrl,
  }) {
    final result = AvailableRelease._();
    if (version != null) result.version = version;
    if (notesUrl != null) result.notesUrl = notesUrl;
    return result;
  }

  AvailableRelease._();

  factory AvailableRelease.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AvailableRelease()..mergeFromBuffer(data, registry);
  factory AvailableRelease.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      AvailableRelease()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AvailableRelease',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: AvailableRelease.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'version')
    ..aOS(2, _omitFieldNames ? '' : 'notesUrl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AvailableRelease clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AvailableRelease copyWith(void Function(AvailableRelease) updates) =>
      super.copyWith((message) => updates(message as AvailableRelease))
          as AvailableRelease;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use AvailableRelease() / AvailableRelease.new instead')
  static AvailableRelease create() => AvailableRelease._();
  static $pb.GeneratedMessage $_createMessage() => AvailableRelease._();
  @$core.override
  AvailableRelease createEmptyInstance() => AvailableRelease._();
  @$core.pragma('dart2js:noInline')
  static AvailableRelease getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AvailableRelease>(
          AvailableRelease.$_createMessage);
  static AvailableRelease? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get version => $_getSZ(0);
  @$pb.TagNumber(1)
  set version($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearVersion() => $_clearField(1);

  /// Notes de version sur GitHub ; vide si la fiche n'en donne pas.
  @$pb.TagNumber(2)
  $core.String get notesUrl => $_getSZ(1);
  @$pb.TagNumber(2)
  set notesUrl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasNotesUrl() => $_has(1);
  @$pb.TagNumber(2)
  void clearNotesUrl() => $_clearField(2);
}

/// UpdateRun décrit la dernière mise à jour lancée sur l'installation.
class UpdateRun extends $pb.GeneratedMessage {
  factory UpdateRun({
    $core.String? fromVersion,
    $core.String? toVersion,
    UpdateStep? step,
    UpdateResult? result,
    $core.String? reason,
    $core.String? backup,
    $1.Timestamp? startedAt,
    $1.Timestamp? finishedAt,
  }) {
    final result$ = UpdateRun._();
    if (fromVersion != null) result$.fromVersion = fromVersion;
    if (toVersion != null) result$.toVersion = toVersion;
    if (step != null) result$.step = step;
    if (result != null) result$.result = result;
    if (reason != null) result$.reason = reason;
    if (backup != null) result$.backup = backup;
    if (startedAt != null) result$.startedAt = startedAt;
    if (finishedAt != null) result$.finishedAt = finishedAt;
    return result$;
  }

  UpdateRun._();

  factory UpdateRun.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRun()..mergeFromBuffer(data, registry);
  factory UpdateRun.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      UpdateRun()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRun',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: UpdateRun.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'fromVersion')
    ..aOS(2, _omitFieldNames ? '' : 'toVersion')
    ..aE<UpdateStep>(3, _omitFieldNames ? '' : 'step',
        enumValues: UpdateStep.values)
    ..aE<UpdateResult>(4, _omitFieldNames ? '' : 'result',
        enumValues: UpdateResult.values)
    ..aOS(5, _omitFieldNames ? '' : 'reason')
    ..aOS(6, _omitFieldNames ? '' : 'backup')
    ..aOM<$1.Timestamp>(7, _omitFieldNames ? '' : 'startedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(8, _omitFieldNames ? '' : 'finishedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRun clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRun copyWith(void Function(UpdateRun) updates) =>
      super.copyWith((message) => updates(message as UpdateRun)) as UpdateRun;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use UpdateRun() / UpdateRun.new instead')
  static UpdateRun create() => UpdateRun._();
  static $pb.GeneratedMessage $_createMessage() => UpdateRun._();
  @$core.override
  UpdateRun createEmptyInstance() => UpdateRun._();
  @$core.pragma('dart2js:noInline')
  static UpdateRun getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRun>(UpdateRun.$_createMessage);
  static UpdateRun? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get fromVersion => $_getSZ(0);
  @$pb.TagNumber(1)
  set fromVersion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasFromVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearFromVersion() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get toVersion => $_getSZ(1);
  @$pb.TagNumber(2)
  set toVersion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasToVersion() => $_has(1);
  @$pb.TagNumber(2)
  void clearToVersion() => $_clearField(2);

  @$pb.TagNumber(3)
  UpdateStep get step => $_getN(2);
  @$pb.TagNumber(3)
  set step(UpdateStep value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasStep() => $_has(2);
  @$pb.TagNumber(3)
  void clearStep() => $_clearField(3);

  @$pb.TagNumber(4)
  UpdateResult get result => $_getN(3);
  @$pb.TagNumber(4)
  set result(UpdateResult value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasResult() => $_has(3);
  @$pb.TagNumber(4)
  void clearResult() => $_clearField(4);

  /// Raison d'un échec ou d'un retour arrière.
  @$pb.TagNumber(5)
  $core.String get reason => $_getSZ(4);
  @$pb.TagNumber(5)
  set reason($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasReason() => $_has(4);
  @$pb.TagNumber(5)
  void clearReason() => $_clearField(5);

  /// Dossier de la sauvegarde prise avant la bascule, sur l'hôte.
  @$pb.TagNumber(6)
  $core.String get backup => $_getSZ(5);
  @$pb.TagNumber(6)
  set backup($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasBackup() => $_has(5);
  @$pb.TagNumber(6)
  void clearBackup() => $_clearField(6);

  @$pb.TagNumber(7)
  $1.Timestamp get startedAt => $_getN(6);
  @$pb.TagNumber(7)
  set startedAt($1.Timestamp value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasStartedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearStartedAt() => $_clearField(7);
  @$pb.TagNumber(7)
  $1.Timestamp ensureStartedAt() => $_ensure(6);

  /// Absent tant que la mise à jour est en cours.
  @$pb.TagNumber(8)
  $1.Timestamp get finishedAt => $_getN(7);
  @$pb.TagNumber(8)
  set finishedAt($1.Timestamp value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasFinishedAt() => $_has(7);
  @$pb.TagNumber(8)
  void clearFinishedAt() => $_clearField(8);
  @$pb.TagNumber(8)
  $1.Timestamp ensureFinishedAt() => $_ensure(7);
}

class GetUpdateStatusRequest extends $pb.GeneratedMessage {
  factory GetUpdateStatusRequest() => GetUpdateStatusRequest._();

  GetUpdateStatusRequest._();

  factory GetUpdateStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetUpdateStatusRequest()..mergeFromBuffer(data, registry);
  factory GetUpdateStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetUpdateStatusRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUpdateStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetUpdateStatusRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUpdateStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUpdateStatusRequest copyWith(
          void Function(GetUpdateStatusRequest) updates) =>
      super.copyWith((message) => updates(message as GetUpdateStatusRequest))
          as GetUpdateStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetUpdateStatusRequest() / GetUpdateStatusRequest.new instead')
  static GetUpdateStatusRequest create() => GetUpdateStatusRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetUpdateStatusRequest._();
  @$core.override
  GetUpdateStatusRequest createEmptyInstance() => GetUpdateStatusRequest._();
  @$core.pragma('dart2js:noInline')
  static GetUpdateStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUpdateStatusRequest>(
          GetUpdateStatusRequest.$_createMessage);
  static GetUpdateStatusRequest? _defaultInstance;
}

class GetUpdateStatusResponse extends $pb.GeneratedMessage {
  factory GetUpdateStatusResponse({
    $core.String? currentVersion,
    AvailableRelease? available,
    $1.Timestamp? checkedAt,
    $core.String? checkError,
    $core.bool? checkEnabled,
    $core.bool? supported,
    UpdateRun? lastRun,
  }) {
    final result = GetUpdateStatusResponse._();
    if (currentVersion != null) result.currentVersion = currentVersion;
    if (available != null) result.available = available;
    if (checkedAt != null) result.checkedAt = checkedAt;
    if (checkError != null) result.checkError = checkError;
    if (checkEnabled != null) result.checkEnabled = checkEnabled;
    if (supported != null) result.supported = supported;
    if (lastRun != null) result.lastRun = lastRun;
    return result;
  }

  GetUpdateStatusResponse._();

  factory GetUpdateStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetUpdateStatusResponse()..mergeFromBuffer(data, registry);
  factory GetUpdateStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetUpdateStatusResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUpdateStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetUpdateStatusResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'currentVersion')
    ..aOM<AvailableRelease>(2, _omitFieldNames ? '' : 'available',
        subBuilder: AvailableRelease.$_createMessage)
    ..aOM<$1.Timestamp>(3, _omitFieldNames ? '' : 'checkedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOS(4, _omitFieldNames ? '' : 'checkError')
    ..aOB(5, _omitFieldNames ? '' : 'checkEnabled')
    ..aOB(6, _omitFieldNames ? '' : 'supported')
    ..aOM<UpdateRun>(7, _omitFieldNames ? '' : 'lastRun',
        subBuilder: UpdateRun.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUpdateStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUpdateStatusResponse copyWith(
          void Function(GetUpdateStatusResponse) updates) =>
      super.copyWith((message) => updates(message as GetUpdateStatusResponse))
          as GetUpdateStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetUpdateStatusResponse() / GetUpdateStatusResponse.new instead')
  static GetUpdateStatusResponse create() => GetUpdateStatusResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetUpdateStatusResponse._();
  @$core.override
  GetUpdateStatusResponse createEmptyInstance() => GetUpdateStatusResponse._();
  @$core.pragma('dart2js:noInline')
  static GetUpdateStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUpdateStatusResponse>(
          GetUpdateStatusResponse.$_createMessage);
  static GetUpdateStatusResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get currentVersion => $_getSZ(0);
  @$pb.TagNumber(1)
  set currentVersion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCurrentVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearCurrentVersion() => $_clearField(1);

  /// Absent : à jour, ou aucune vérification réussie.
  @$pb.TagNumber(2)
  AvailableRelease get available => $_getN(1);
  @$pb.TagNumber(2)
  set available(AvailableRelease value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasAvailable() => $_has(1);
  @$pb.TagNumber(2)
  void clearAvailable() => $_clearField(2);
  @$pb.TagNumber(2)
  AvailableRelease ensureAvailable() => $_ensure(1);

  /// Dernière vérification réussie ; absente si aucune.
  @$pb.TagNumber(3)
  $1.Timestamp get checkedAt => $_getN(2);
  @$pb.TagNumber(3)
  set checkedAt($1.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasCheckedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearCheckedAt() => $_clearField(3);
  @$pb.TagNumber(3)
  $1.Timestamp ensureCheckedAt() => $_ensure(2);

  /// Erreur de la dernière vérification, vide si elle a réussi.
  @$pb.TagNumber(4)
  $core.String get checkError => $_getSZ(3);
  @$pb.TagNumber(4)
  set checkError($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCheckError() => $_has(3);
  @$pb.TagNumber(4)
  void clearCheckError() => $_clearField(4);

  /// Faux si la vérification quotidienne est coupée (updates.check: false).
  @$pb.TagNumber(5)
  $core.bool get checkEnabled => $_getBF(4);
  @$pb.TagNumber(5)
  set checkEnabled($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(5)
  $core.bool hasCheckEnabled() => $_has(4);
  @$pb.TagNumber(5)
  void clearCheckEnabled() => $_clearField(5);

  /// Vrai si cette installation sait se mettre à jour depuis le dashboard
  /// (installation standard d'install.sh) ; sinon, la mise à jour se fait à
  /// la main.
  @$pb.TagNumber(6)
  $core.bool get supported => $_getBF(5);
  @$pb.TagNumber(6)
  set supported($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasSupported() => $_has(5);
  @$pb.TagNumber(6)
  void clearSupported() => $_clearField(6);

  /// Dernière mise à jour lancée ; absente si aucune.
  @$pb.TagNumber(7)
  UpdateRun get lastRun => $_getN(6);
  @$pb.TagNumber(7)
  set lastRun(UpdateRun value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasLastRun() => $_has(6);
  @$pb.TagNumber(7)
  void clearLastRun() => $_clearField(7);
  @$pb.TagNumber(7)
  UpdateRun ensureLastRun() => $_ensure(6);
}

class CheckForUpdatesRequest extends $pb.GeneratedMessage {
  factory CheckForUpdatesRequest() => CheckForUpdatesRequest._();

  CheckForUpdatesRequest._();

  factory CheckForUpdatesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckForUpdatesRequest()..mergeFromBuffer(data, registry);
  factory CheckForUpdatesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckForUpdatesRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckForUpdatesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CheckForUpdatesRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckForUpdatesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckForUpdatesRequest copyWith(
          void Function(CheckForUpdatesRequest) updates) =>
      super.copyWith((message) => updates(message as CheckForUpdatesRequest))
          as CheckForUpdatesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CheckForUpdatesRequest() / CheckForUpdatesRequest.new instead')
  static CheckForUpdatesRequest create() => CheckForUpdatesRequest._();
  static $pb.GeneratedMessage $_createMessage() => CheckForUpdatesRequest._();
  @$core.override
  CheckForUpdatesRequest createEmptyInstance() => CheckForUpdatesRequest._();
  @$core.pragma('dart2js:noInline')
  static CheckForUpdatesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckForUpdatesRequest>(
          CheckForUpdatesRequest.$_createMessage);
  static CheckForUpdatesRequest? _defaultInstance;
}

class CheckForUpdatesResponse extends $pb.GeneratedMessage {
  factory CheckForUpdatesResponse({
    GetUpdateStatusResponse? status,
  }) {
    final result = CheckForUpdatesResponse._();
    if (status != null) result.status = status;
    return result;
  }

  CheckForUpdatesResponse._();

  factory CheckForUpdatesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckForUpdatesResponse()..mergeFromBuffer(data, registry);
  factory CheckForUpdatesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CheckForUpdatesResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckForUpdatesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CheckForUpdatesResponse.$_createMessage)
    ..aOM<GetUpdateStatusResponse>(1, _omitFieldNames ? '' : 'status',
        subBuilder: GetUpdateStatusResponse.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckForUpdatesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckForUpdatesResponse copyWith(
          void Function(CheckForUpdatesResponse) updates) =>
      super.copyWith((message) => updates(message as CheckForUpdatesResponse))
          as CheckForUpdatesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CheckForUpdatesResponse() / CheckForUpdatesResponse.new instead')
  static CheckForUpdatesResponse create() => CheckForUpdatesResponse._();
  static $pb.GeneratedMessage $_createMessage() => CheckForUpdatesResponse._();
  @$core.override
  CheckForUpdatesResponse createEmptyInstance() => CheckForUpdatesResponse._();
  @$core.pragma('dart2js:noInline')
  static CheckForUpdatesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckForUpdatesResponse>(
          CheckForUpdatesResponse.$_createMessage);
  static CheckForUpdatesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  GetUpdateStatusResponse get status => $_getN(0);
  @$pb.TagNumber(1)
  set status(GetUpdateStatusResponse value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);
  @$pb.TagNumber(1)
  GetUpdateStatusResponse ensureStatus() => $_ensure(0);
}

class StartUpdateRequest extends $pb.GeneratedMessage {
  factory StartUpdateRequest() => StartUpdateRequest._();

  StartUpdateRequest._();

  factory StartUpdateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartUpdateRequest()..mergeFromBuffer(data, registry);
  factory StartUpdateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartUpdateRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartUpdateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StartUpdateRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartUpdateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartUpdateRequest copyWith(void Function(StartUpdateRequest) updates) =>
      super.copyWith((message) => updates(message as StartUpdateRequest))
          as StartUpdateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use StartUpdateRequest() / StartUpdateRequest.new instead')
  static StartUpdateRequest create() => StartUpdateRequest._();
  static $pb.GeneratedMessage $_createMessage() => StartUpdateRequest._();
  @$core.override
  StartUpdateRequest createEmptyInstance() => StartUpdateRequest._();
  @$core.pragma('dart2js:noInline')
  static StartUpdateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartUpdateRequest>(
          StartUpdateRequest.$_createMessage);
  static StartUpdateRequest? _defaultInstance;
}

class StartUpdateResponse extends $pb.GeneratedMessage {
  factory StartUpdateResponse({
    $core.String? toVersion,
  }) {
    final result = StartUpdateResponse._();
    if (toVersion != null) result.toVersion = toVersion;
    return result;
  }

  StartUpdateResponse._();

  factory StartUpdateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartUpdateResponse()..mergeFromBuffer(data, registry);
  factory StartUpdateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StartUpdateResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartUpdateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StartUpdateResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'toVersion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartUpdateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartUpdateResponse copyWith(void Function(StartUpdateResponse) updates) =>
      super.copyWith((message) => updates(message as StartUpdateResponse))
          as StartUpdateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core
      .Deprecated('Use StartUpdateResponse() / StartUpdateResponse.new instead')
  static StartUpdateResponse create() => StartUpdateResponse._();
  static $pb.GeneratedMessage $_createMessage() => StartUpdateResponse._();
  @$core.override
  StartUpdateResponse createEmptyInstance() => StartUpdateResponse._();
  @$core.pragma('dart2js:noInline')
  static StartUpdateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartUpdateResponse>(
          StartUpdateResponse.$_createMessage);
  static StartUpdateResponse? _defaultInstance;

  /// Version qui va être installée.
  @$pb.TagNumber(1)
  $core.String get toVersion => $_getSZ(0);
  @$pb.TagNumber(1)
  set toVersion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearToVersion() => $_clearField(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
