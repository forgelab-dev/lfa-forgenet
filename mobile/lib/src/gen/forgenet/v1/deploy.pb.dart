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
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

import 'container.pbenum.dart' as $2;
import 'deploy.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'deploy.pbenum.dart';

class RecreateStrategy extends $pb.GeneratedMessage {
  factory RecreateStrategy() => RecreateStrategy._();

  RecreateStrategy._();

  factory RecreateStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RecreateStrategy()..mergeFromBuffer(data, registry);
  factory RecreateStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RecreateStrategy()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecreateStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RecreateStrategy.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecreateStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecreateStrategy copyWith(void Function(RecreateStrategy) updates) =>
      super.copyWith((message) => updates(message as RecreateStrategy))
          as RecreateStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RecreateStrategy() / RecreateStrategy.new instead')
  static RecreateStrategy create() => RecreateStrategy._();
  static $pb.GeneratedMessage $_createMessage() => RecreateStrategy._();
  @$core.override
  RecreateStrategy createEmptyInstance() => RecreateStrategy._();
  @$core.pragma('dart2js:noInline')
  static RecreateStrategy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RecreateStrategy>(
          RecreateStrategy.$_createMessage);
  static RecreateStrategy? _defaultInstance;
}

class RollingStrategy extends $pb.GeneratedMessage {
  factory RollingStrategy({
    $core.int? maxSurge,
    $core.int? maxUnavailable,
  }) {
    final result = RollingStrategy._();
    if (maxSurge != null) result.maxSurge = maxSurge;
    if (maxUnavailable != null) result.maxUnavailable = maxUnavailable;
    return result;
  }

  RollingStrategy._();

  factory RollingStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollingStrategy()..mergeFromBuffer(data, registry);
  factory RollingStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollingStrategy()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RollingStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RollingStrategy.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'maxSurge', fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'maxUnavailable',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollingStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollingStrategy copyWith(void Function(RollingStrategy) updates) =>
      super.copyWith((message) => updates(message as RollingStrategy))
          as RollingStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use RollingStrategy() / RollingStrategy.new instead')
  static RollingStrategy create() => RollingStrategy._();
  static $pb.GeneratedMessage $_createMessage() => RollingStrategy._();
  @$core.override
  RollingStrategy createEmptyInstance() => RollingStrategy._();
  @$core.pragma('dart2js:noInline')
  static RollingStrategy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RollingStrategy>(
          RollingStrategy.$_createMessage);
  static RollingStrategy? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get maxSurge => $_getIZ(0);
  @$pb.TagNumber(1)
  set maxSurge($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMaxSurge() => $_has(0);
  @$pb.TagNumber(1)
  void clearMaxSurge() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get maxUnavailable => $_getIZ(1);
  @$pb.TagNumber(2)
  set maxUnavailable($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMaxUnavailable() => $_has(1);
  @$pb.TagNumber(2)
  void clearMaxUnavailable() => $_clearField(2);
}

class BlueGreenStrategy extends $pb.GeneratedMessage {
  factory BlueGreenStrategy({
    $core.int? autoPromoteSeconds,
  }) {
    final result = BlueGreenStrategy._();
    if (autoPromoteSeconds != null)
      result.autoPromoteSeconds = autoPromoteSeconds;
    return result;
  }

  BlueGreenStrategy._();

  factory BlueGreenStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BlueGreenStrategy()..mergeFromBuffer(data, registry);
  factory BlueGreenStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      BlueGreenStrategy()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BlueGreenStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: BlueGreenStrategy.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'autoPromoteSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlueGreenStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlueGreenStrategy copyWith(void Function(BlueGreenStrategy) updates) =>
      super.copyWith((message) => updates(message as BlueGreenStrategy))
          as BlueGreenStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use BlueGreenStrategy() / BlueGreenStrategy.new instead')
  static BlueGreenStrategy create() => BlueGreenStrategy._();
  static $pb.GeneratedMessage $_createMessage() => BlueGreenStrategy._();
  @$core.override
  BlueGreenStrategy createEmptyInstance() => BlueGreenStrategy._();
  @$core.pragma('dart2js:noInline')
  static BlueGreenStrategy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlueGreenStrategy>(
          BlueGreenStrategy.$_createMessage);
  static BlueGreenStrategy? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get autoPromoteSeconds => $_getIZ(0);
  @$pb.TagNumber(1)
  set autoPromoteSeconds($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAutoPromoteSeconds() => $_has(0);
  @$pb.TagNumber(1)
  void clearAutoPromoteSeconds() => $_clearField(1);
}

class CanaryStrategy extends $pb.GeneratedMessage {
  factory CanaryStrategy({
    $core.int? initialWeightPercent,
    $core.int? stepWeightPercent,
    $core.int? stepIntervalSeconds,
  }) {
    final result = CanaryStrategy._();
    if (initialWeightPercent != null)
      result.initialWeightPercent = initialWeightPercent;
    if (stepWeightPercent != null) result.stepWeightPercent = stepWeightPercent;
    if (stepIntervalSeconds != null)
      result.stepIntervalSeconds = stepIntervalSeconds;
    return result;
  }

  CanaryStrategy._();

  factory CanaryStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CanaryStrategy()..mergeFromBuffer(data, registry);
  factory CanaryStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CanaryStrategy()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CanaryStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CanaryStrategy.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'initialWeightPercent',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'stepWeightPercent',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(3, _omitFieldNames ? '' : 'stepIntervalSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CanaryStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CanaryStrategy copyWith(void Function(CanaryStrategy) updates) =>
      super.copyWith((message) => updates(message as CanaryStrategy))
          as CanaryStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CanaryStrategy() / CanaryStrategy.new instead')
  static CanaryStrategy create() => CanaryStrategy._();
  static $pb.GeneratedMessage $_createMessage() => CanaryStrategy._();
  @$core.override
  CanaryStrategy createEmptyInstance() => CanaryStrategy._();
  @$core.pragma('dart2js:noInline')
  static CanaryStrategy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CanaryStrategy>(
          CanaryStrategy.$_createMessage);
  static CanaryStrategy? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get initialWeightPercent => $_getIZ(0);
  @$pb.TagNumber(1)
  set initialWeightPercent($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasInitialWeightPercent() => $_has(0);
  @$pb.TagNumber(1)
  void clearInitialWeightPercent() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get stepWeightPercent => $_getIZ(1);
  @$pb.TagNumber(2)
  set stepWeightPercent($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasStepWeightPercent() => $_has(1);
  @$pb.TagNumber(2)
  void clearStepWeightPercent() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get stepIntervalSeconds => $_getIZ(2);
  @$pb.TagNumber(3)
  set stepIntervalSeconds($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStepIntervalSeconds() => $_has(2);
  @$pb.TagNumber(3)
  void clearStepIntervalSeconds() => $_clearField(3);
}

enum StrategyConfig_Strategy { recreate, rolling, blueGreen, canary, notSet }

/// StrategyConfig définit la méthode d'exécution et sa configuration associée.
class StrategyConfig extends $pb.GeneratedMessage {
  factory StrategyConfig({
    RecreateStrategy? recreate,
    RollingStrategy? rolling,
    BlueGreenStrategy? blueGreen,
    CanaryStrategy? canary,
  }) {
    final result = StrategyConfig._();
    if (recreate != null) result.recreate = recreate;
    if (rolling != null) result.rolling = rolling;
    if (blueGreen != null) result.blueGreen = blueGreen;
    if (canary != null) result.canary = canary;
    return result;
  }

  StrategyConfig._();

  factory StrategyConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StrategyConfig()..mergeFromBuffer(data, registry);
  factory StrategyConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StrategyConfig()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, StrategyConfig_Strategy>
      _StrategyConfig_StrategyByTag = {
    1: StrategyConfig_Strategy.recreate,
    2: StrategyConfig_Strategy.rolling,
    3: StrategyConfig_Strategy.blueGreen,
    4: StrategyConfig_Strategy.canary,
    0: StrategyConfig_Strategy.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StrategyConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StrategyConfig.$_createMessage)
    ..oo(0, [1, 2, 3, 4])
    ..aOM<RecreateStrategy>(1, _omitFieldNames ? '' : 'recreate',
        subBuilder: RecreateStrategy.$_createMessage)
    ..aOM<RollingStrategy>(2, _omitFieldNames ? '' : 'rolling',
        subBuilder: RollingStrategy.$_createMessage)
    ..aOM<BlueGreenStrategy>(3, _omitFieldNames ? '' : 'blueGreen',
        subBuilder: BlueGreenStrategy.$_createMessage)
    ..aOM<CanaryStrategy>(4, _omitFieldNames ? '' : 'canary',
        subBuilder: CanaryStrategy.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StrategyConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StrategyConfig copyWith(void Function(StrategyConfig) updates) =>
      super.copyWith((message) => updates(message as StrategyConfig))
          as StrategyConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use StrategyConfig() / StrategyConfig.new instead')
  static StrategyConfig create() => StrategyConfig._();
  static $pb.GeneratedMessage $_createMessage() => StrategyConfig._();
  @$core.override
  StrategyConfig createEmptyInstance() => StrategyConfig._();
  @$core.pragma('dart2js:noInline')
  static StrategyConfig getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<StrategyConfig>(
          StrategyConfig.$_createMessage);
  static StrategyConfig? _defaultInstance;

  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  @$pb.TagNumber(4)
  StrategyConfig_Strategy whichStrategy() =>
      _StrategyConfig_StrategyByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  @$pb.TagNumber(4)
  void clearStrategy() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  RecreateStrategy get recreate => $_getN(0);
  @$pb.TagNumber(1)
  set recreate(RecreateStrategy value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRecreate() => $_has(0);
  @$pb.TagNumber(1)
  void clearRecreate() => $_clearField(1);
  @$pb.TagNumber(1)
  RecreateStrategy ensureRecreate() => $_ensure(0);

  @$pb.TagNumber(2)
  RollingStrategy get rolling => $_getN(1);
  @$pb.TagNumber(2)
  set rolling(RollingStrategy value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasRolling() => $_has(1);
  @$pb.TagNumber(2)
  void clearRolling() => $_clearField(2);
  @$pb.TagNumber(2)
  RollingStrategy ensureRolling() => $_ensure(1);

  @$pb.TagNumber(3)
  BlueGreenStrategy get blueGreen => $_getN(2);
  @$pb.TagNumber(3)
  set blueGreen(BlueGreenStrategy value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasBlueGreen() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlueGreen() => $_clearField(3);
  @$pb.TagNumber(3)
  BlueGreenStrategy ensureBlueGreen() => $_ensure(2);

  @$pb.TagNumber(4)
  CanaryStrategy get canary => $_getN(3);
  @$pb.TagNumber(4)
  set canary(CanaryStrategy value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasCanary() => $_has(3);
  @$pb.TagNumber(4)
  void clearCanary() => $_clearField(4);
  @$pb.TagNumber(4)
  CanaryStrategy ensureCanary() => $_ensure(3);
}

/// GitSource définit un déploiement basé sur un dépôt Git à builder.
class GitSource extends $pb.GeneratedMessage {
  factory GitSource({
    $core.String? repositoryUrl,
    $core.String? branch,
    $core.String? commitSha,
    $core.String? dockerfilePath,
  }) {
    final result = GitSource._();
    if (repositoryUrl != null) result.repositoryUrl = repositoryUrl;
    if (branch != null) result.branch = branch;
    if (commitSha != null) result.commitSha = commitSha;
    if (dockerfilePath != null) result.dockerfilePath = dockerfilePath;
    return result;
  }

  GitSource._();

  factory GitSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GitSource()..mergeFromBuffer(data, registry);
  factory GitSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GitSource()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GitSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GitSource.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'repositoryUrl')
    ..aOS(2, _omitFieldNames ? '' : 'branch')
    ..aOS(3, _omitFieldNames ? '' : 'commitSha')
    ..aOS(4, _omitFieldNames ? '' : 'dockerfilePath')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GitSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GitSource copyWith(void Function(GitSource) updates) =>
      super.copyWith((message) => updates(message as GitSource)) as GitSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use GitSource() / GitSource.new instead')
  static GitSource create() => GitSource._();
  static $pb.GeneratedMessage $_createMessage() => GitSource._();
  @$core.override
  GitSource createEmptyInstance() => GitSource._();
  @$core.pragma('dart2js:noInline')
  static GitSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GitSource>(GitSource.$_createMessage);
  static GitSource? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get repositoryUrl => $_getSZ(0);
  @$pb.TagNumber(1)
  set repositoryUrl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRepositoryUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearRepositoryUrl() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get branch => $_getSZ(1);
  @$pb.TagNumber(2)
  set branch($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBranch() => $_has(1);
  @$pb.TagNumber(2)
  void clearBranch() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get commitSha => $_getSZ(2);
  @$pb.TagNumber(3)
  set commitSha($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCommitSha() => $_has(2);
  @$pb.TagNumber(3)
  void clearCommitSha() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get dockerfilePath => $_getSZ(3);
  @$pb.TagNumber(4)
  set dockerfilePath($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDockerfilePath() => $_has(3);
  @$pb.TagNumber(4)
  void clearDockerfilePath() => $_clearField(4);
}

/// ImageSource définit un déploiement depuis une image Docker pré-construite.
class ImageSource extends $pb.GeneratedMessage {
  factory ImageSource({
    $core.String? imageRef,
  }) {
    final result = ImageSource._();
    if (imageRef != null) result.imageRef = imageRef;
    return result;
  }

  ImageSource._();

  factory ImageSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ImageSource()..mergeFromBuffer(data, registry);
  factory ImageSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ImageSource()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImageSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ImageSource.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'imageRef')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImageSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImageSource copyWith(void Function(ImageSource) updates) =>
      super.copyWith((message) => updates(message as ImageSource))
          as ImageSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ImageSource() / ImageSource.new instead')
  static ImageSource create() => ImageSource._();
  static $pb.GeneratedMessage $_createMessage() => ImageSource._();
  @$core.override
  ImageSource createEmptyInstance() => ImageSource._();
  @$core.pragma('dart2js:noInline')
  static ImageSource getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ImageSource>(
          ImageSource.$_createMessage);
  static ImageSource? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get imageRef => $_getSZ(0);
  @$pb.TagNumber(1)
  set imageRef($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasImageRef() => $_has(0);
  @$pb.TagNumber(1)
  void clearImageRef() => $_clearField(1);
}

/// ComposeSource définit un déploiement basé sur un fichier Docker Compose.
class ComposeSource extends $pb.GeneratedMessage {
  factory ComposeSource({
    $core.String? composeContent,
  }) {
    final result = ComposeSource._();
    if (composeContent != null) result.composeContent = composeContent;
    return result;
  }

  ComposeSource._();

  factory ComposeSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ComposeSource()..mergeFromBuffer(data, registry);
  factory ComposeSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ComposeSource()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ComposeSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ComposeSource.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'composeContent')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ComposeSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ComposeSource copyWith(void Function(ComposeSource) updates) =>
      super.copyWith((message) => updates(message as ComposeSource))
          as ComposeSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ComposeSource() / ComposeSource.new instead')
  static ComposeSource create() => ComposeSource._();
  static $pb.GeneratedMessage $_createMessage() => ComposeSource._();
  @$core.override
  ComposeSource createEmptyInstance() => ComposeSource._();
  @$core.pragma('dart2js:noInline')
  static ComposeSource getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ComposeSource>(
          ComposeSource.$_createMessage);
  static ComposeSource? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get composeContent => $_getSZ(0);
  @$pb.TagNumber(1)
  set composeContent($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasComposeContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearComposeContent() => $_clearField(1);
}

enum DeploymentSource_Source { git, image, compose, notSet }

/// DeploymentSource garantit qu'une seule source est spécifiée à la fois.
class DeploymentSource extends $pb.GeneratedMessage {
  factory DeploymentSource({
    GitSource? git,
    ImageSource? image,
    ComposeSource? compose,
  }) {
    final result = DeploymentSource._();
    if (git != null) result.git = git;
    if (image != null) result.image = image;
    if (compose != null) result.compose = compose;
    return result;
  }

  DeploymentSource._();

  factory DeploymentSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeploymentSource()..mergeFromBuffer(data, registry);
  factory DeploymentSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeploymentSource()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, DeploymentSource_Source>
      _DeploymentSource_SourceByTag = {
    1: DeploymentSource_Source.git,
    2: DeploymentSource_Source.image,
    3: DeploymentSource_Source.compose,
    0: DeploymentSource_Source.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeploymentSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DeploymentSource.$_createMessage)
    ..oo(0, [1, 2, 3])
    ..aOM<GitSource>(1, _omitFieldNames ? '' : 'git',
        subBuilder: GitSource.$_createMessage)
    ..aOM<ImageSource>(2, _omitFieldNames ? '' : 'image',
        subBuilder: ImageSource.$_createMessage)
    ..aOM<ComposeSource>(3, _omitFieldNames ? '' : 'compose',
        subBuilder: ComposeSource.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentSource copyWith(void Function(DeploymentSource) updates) =>
      super.copyWith((message) => updates(message as DeploymentSource))
          as DeploymentSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DeploymentSource() / DeploymentSource.new instead')
  static DeploymentSource create() => DeploymentSource._();
  static $pb.GeneratedMessage $_createMessage() => DeploymentSource._();
  @$core.override
  DeploymentSource createEmptyInstance() => DeploymentSource._();
  @$core.pragma('dart2js:noInline')
  static DeploymentSource getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeploymentSource>(
          DeploymentSource.$_createMessage);
  static DeploymentSource? _defaultInstance;

  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  DeploymentSource_Source whichSource() =>
      _DeploymentSource_SourceByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(1)
  @$pb.TagNumber(2)
  @$pb.TagNumber(3)
  void clearSource() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  GitSource get git => $_getN(0);
  @$pb.TagNumber(1)
  set git(GitSource value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasGit() => $_has(0);
  @$pb.TagNumber(1)
  void clearGit() => $_clearField(1);
  @$pb.TagNumber(1)
  GitSource ensureGit() => $_ensure(0);

  @$pb.TagNumber(2)
  ImageSource get image => $_getN(1);
  @$pb.TagNumber(2)
  set image(ImageSource value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasImage() => $_has(1);
  @$pb.TagNumber(2)
  void clearImage() => $_clearField(2);
  @$pb.TagNumber(2)
  ImageSource ensureImage() => $_ensure(1);

  @$pb.TagNumber(3)
  ComposeSource get compose => $_getN(2);
  @$pb.TagNumber(3)
  set compose(ComposeSource value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasCompose() => $_has(2);
  @$pb.TagNumber(3)
  void clearCompose() => $_clearField(3);
  @$pb.TagNumber(3)
  ComposeSource ensureCompose() => $_ensure(2);
}

/// Deployment représente l'entité globale d'une opération de déploiement.
class Deployment extends $pb.GeneratedMessage {
  factory Deployment({
    $core.String? id,
    $core.String? projectId,
    DeploymentSource? source,
    DeploymentState? state,
    StrategyConfig? strategy,
    $core.String? errorMessage,
    $core.String? targetEnv,
    $core.String? previousDeploymentId,
    $core.String? triggeredBy,
    $1.Timestamp? createdAt,
    $1.Timestamp? updatedAt,
    $1.Timestamp? finishedAt,
  }) {
    final result = Deployment._();
    if (id != null) result.id = id;
    if (projectId != null) result.projectId = projectId;
    if (source != null) result.source = source;
    if (state != null) result.state = state;
    if (strategy != null) result.strategy = strategy;
    if (errorMessage != null) result.errorMessage = errorMessage;
    if (targetEnv != null) result.targetEnv = targetEnv;
    if (previousDeploymentId != null)
      result.previousDeploymentId = previousDeploymentId;
    if (triggeredBy != null) result.triggeredBy = triggeredBy;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (finishedAt != null) result.finishedAt = finishedAt;
    return result;
  }

  Deployment._();

  factory Deployment.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Deployment()..mergeFromBuffer(data, registry);
  factory Deployment.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      Deployment()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Deployment',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: Deployment.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'projectId')
    ..aOM<DeploymentSource>(3, _omitFieldNames ? '' : 'source',
        subBuilder: DeploymentSource.$_createMessage)
    ..aE<DeploymentState>(4, _omitFieldNames ? '' : 'state',
        enumValues: DeploymentState.values)
    ..aOM<StrategyConfig>(5, _omitFieldNames ? '' : 'strategy',
        subBuilder: StrategyConfig.$_createMessage)
    ..aOS(6, _omitFieldNames ? '' : 'errorMessage')
    ..aOS(7, _omitFieldNames ? '' : 'targetEnv')
    ..aOS(8, _omitFieldNames ? '' : 'previousDeploymentId')
    ..aOS(9, _omitFieldNames ? '' : 'triggeredBy')
    ..aOM<$1.Timestamp>(10, _omitFieldNames ? '' : 'createdAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(11, _omitFieldNames ? '' : 'updatedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(12, _omitFieldNames ? '' : 'finishedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployment clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployment copyWith(void Function(Deployment) updates) =>
      super.copyWith((message) => updates(message as Deployment)) as Deployment;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use Deployment() / Deployment.new instead')
  static Deployment create() => Deployment._();
  static $pb.GeneratedMessage $_createMessage() => Deployment._();
  @$core.override
  Deployment createEmptyInstance() => Deployment._();
  @$core.pragma('dart2js:noInline')
  static Deployment getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Deployment>(Deployment.$_createMessage);
  static Deployment? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get projectId => $_getSZ(1);
  @$pb.TagNumber(2)
  set projectId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasProjectId() => $_has(1);
  @$pb.TagNumber(2)
  void clearProjectId() => $_clearField(2);

  @$pb.TagNumber(3)
  DeploymentSource get source => $_getN(2);
  @$pb.TagNumber(3)
  set source(DeploymentSource value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasSource() => $_has(2);
  @$pb.TagNumber(3)
  void clearSource() => $_clearField(3);
  @$pb.TagNumber(3)
  DeploymentSource ensureSource() => $_ensure(2);

  @$pb.TagNumber(4)
  DeploymentState get state => $_getN(3);
  @$pb.TagNumber(4)
  set state(DeploymentState value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(4)
  void clearState() => $_clearField(4);

  @$pb.TagNumber(5)
  StrategyConfig get strategy => $_getN(4);
  @$pb.TagNumber(5)
  set strategy(StrategyConfig value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasStrategy() => $_has(4);
  @$pb.TagNumber(5)
  void clearStrategy() => $_clearField(5);
  @$pb.TagNumber(5)
  StrategyConfig ensureStrategy() => $_ensure(4);

  @$pb.TagNumber(6)
  $core.String get errorMessage => $_getSZ(5);
  @$pb.TagNumber(6)
  set errorMessage($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasErrorMessage() => $_has(5);
  @$pb.TagNumber(6)
  void clearErrorMessage() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get targetEnv => $_getSZ(6);
  @$pb.TagNumber(7)
  set targetEnv($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasTargetEnv() => $_has(6);
  @$pb.TagNumber(7)
  void clearTargetEnv() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get previousDeploymentId => $_getSZ(7);
  @$pb.TagNumber(8)
  set previousDeploymentId($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasPreviousDeploymentId() => $_has(7);
  @$pb.TagNumber(8)
  void clearPreviousDeploymentId() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get triggeredBy => $_getSZ(8);
  @$pb.TagNumber(9)
  set triggeredBy($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasTriggeredBy() => $_has(8);
  @$pb.TagNumber(9)
  void clearTriggeredBy() => $_clearField(9);

  @$pb.TagNumber(10)
  $1.Timestamp get createdAt => $_getN(9);
  @$pb.TagNumber(10)
  set createdAt($1.Timestamp value) => $_setField(10, value);
  @$pb.TagNumber(10)
  $core.bool hasCreatedAt() => $_has(9);
  @$pb.TagNumber(10)
  void clearCreatedAt() => $_clearField(10);
  @$pb.TagNumber(10)
  $1.Timestamp ensureCreatedAt() => $_ensure(9);

  @$pb.TagNumber(11)
  $1.Timestamp get updatedAt => $_getN(10);
  @$pb.TagNumber(11)
  set updatedAt($1.Timestamp value) => $_setField(11, value);
  @$pb.TagNumber(11)
  $core.bool hasUpdatedAt() => $_has(10);
  @$pb.TagNumber(11)
  void clearUpdatedAt() => $_clearField(11);
  @$pb.TagNumber(11)
  $1.Timestamp ensureUpdatedAt() => $_ensure(10);

  @$pb.TagNumber(12)
  $1.Timestamp get finishedAt => $_getN(11);
  @$pb.TagNumber(12)
  set finishedAt($1.Timestamp value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasFinishedAt() => $_has(11);
  @$pb.TagNumber(12)
  void clearFinishedAt() => $_clearField(12);
  @$pb.TagNumber(12)
  $1.Timestamp ensureFinishedAt() => $_ensure(11);
}

/// DeploymentLogEntry représente une ligne de log d'exécution du build ou déploiement.
class DeploymentLogEntry extends $pb.GeneratedMessage {
  factory DeploymentLogEntry({
    $1.Timestamp? emittedAt,
    $core.String? message,
    $core.String? step,
    $2.LogStream? stream,
  }) {
    final result = DeploymentLogEntry._();
    if (emittedAt != null) result.emittedAt = emittedAt;
    if (message != null) result.message = message;
    if (step != null) result.step = step;
    if (stream != null) result.stream = stream;
    return result;
  }

  DeploymentLogEntry._();

  factory DeploymentLogEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeploymentLogEntry()..mergeFromBuffer(data, registry);
  factory DeploymentLogEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DeploymentLogEntry()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeploymentLogEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DeploymentLogEntry.$_createMessage)
    ..aOM<$1.Timestamp>(1, _omitFieldNames ? '' : 'emittedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'message')
    ..aOS(3, _omitFieldNames ? '' : 'step')
    ..aE<$2.LogStream>(4, _omitFieldNames ? '' : 'stream',
        enumValues: $2.LogStream.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentLogEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentLogEntry copyWith(void Function(DeploymentLogEntry) updates) =>
      super.copyWith((message) => updates(message as DeploymentLogEntry))
          as DeploymentLogEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DeploymentLogEntry() / DeploymentLogEntry.new instead')
  static DeploymentLogEntry create() => DeploymentLogEntry._();
  static $pb.GeneratedMessage $_createMessage() => DeploymentLogEntry._();
  @$core.override
  DeploymentLogEntry createEmptyInstance() => DeploymentLogEntry._();
  @$core.pragma('dart2js:noInline')
  static DeploymentLogEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeploymentLogEntry>(
          DeploymentLogEntry.$_createMessage);
  static DeploymentLogEntry? _defaultInstance;

  @$pb.TagNumber(1)
  $1.Timestamp get emittedAt => $_getN(0);
  @$pb.TagNumber(1)
  set emittedAt($1.Timestamp value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasEmittedAt() => $_has(0);
  @$pb.TagNumber(1)
  void clearEmittedAt() => $_clearField(1);
  @$pb.TagNumber(1)
  $1.Timestamp ensureEmittedAt() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(2)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get step => $_getSZ(2);
  @$pb.TagNumber(3)
  set step($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStep() => $_has(2);
  @$pb.TagNumber(3)
  void clearStep() => $_clearField(3);

  @$pb.TagNumber(4)
  $2.LogStream get stream => $_getN(3);
  @$pb.TagNumber(4)
  set stream($2.LogStream value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasStream() => $_has(3);
  @$pb.TagNumber(4)
  void clearStream() => $_clearField(4);
}

class CreateDeploymentRequest extends $pb.GeneratedMessage {
  factory CreateDeploymentRequest({
    $core.String? projectId,
    DeploymentSource? source,
    StrategyConfig? strategy,
    $core.String? targetEnv,
    $core.String? requestId,
  }) {
    final result = CreateDeploymentRequest._();
    if (projectId != null) result.projectId = projectId;
    if (source != null) result.source = source;
    if (strategy != null) result.strategy = strategy;
    if (targetEnv != null) result.targetEnv = targetEnv;
    if (requestId != null) result.requestId = requestId;
    return result;
  }

  CreateDeploymentRequest._();

  factory CreateDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateDeploymentRequest()..mergeFromBuffer(data, registry);
  factory CreateDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateDeploymentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CreateDeploymentRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'projectId')
    ..aOM<DeploymentSource>(2, _omitFieldNames ? '' : 'source',
        subBuilder: DeploymentSource.$_createMessage)
    ..aOM<StrategyConfig>(3, _omitFieldNames ? '' : 'strategy',
        subBuilder: StrategyConfig.$_createMessage)
    ..aOS(4, _omitFieldNames ? '' : 'targetEnv')
    ..aOS(5, _omitFieldNames ? '' : 'requestId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentRequest copyWith(
          void Function(CreateDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as CreateDeploymentRequest))
          as CreateDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateDeploymentRequest() / CreateDeploymentRequest.new instead')
  static CreateDeploymentRequest create() => CreateDeploymentRequest._();
  static $pb.GeneratedMessage $_createMessage() => CreateDeploymentRequest._();
  @$core.override
  CreateDeploymentRequest createEmptyInstance() => CreateDeploymentRequest._();
  @$core.pragma('dart2js:noInline')
  static CreateDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDeploymentRequest>(
          CreateDeploymentRequest.$_createMessage);
  static CreateDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get projectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set projectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasProjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProjectId() => $_clearField(1);

  @$pb.TagNumber(2)
  DeploymentSource get source => $_getN(1);
  @$pb.TagNumber(2)
  set source(DeploymentSource value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasSource() => $_has(1);
  @$pb.TagNumber(2)
  void clearSource() => $_clearField(2);
  @$pb.TagNumber(2)
  DeploymentSource ensureSource() => $_ensure(1);

  @$pb.TagNumber(3)
  StrategyConfig get strategy => $_getN(2);
  @$pb.TagNumber(3)
  set strategy(StrategyConfig value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasStrategy() => $_has(2);
  @$pb.TagNumber(3)
  void clearStrategy() => $_clearField(3);
  @$pb.TagNumber(3)
  StrategyConfig ensureStrategy() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.String get targetEnv => $_getSZ(3);
  @$pb.TagNumber(4)
  set targetEnv($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasTargetEnv() => $_has(3);
  @$pb.TagNumber(4)
  void clearTargetEnv() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get requestId => $_getSZ(4);
  @$pb.TagNumber(5)
  set requestId($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasRequestId() => $_has(4);
  @$pb.TagNumber(5)
  void clearRequestId() => $_clearField(5);
}

class CreateDeploymentResponse extends $pb.GeneratedMessage {
  factory CreateDeploymentResponse({
    Deployment? deployment,
  }) {
    final result = CreateDeploymentResponse._();
    if (deployment != null) result.deployment = deployment;
    return result;
  }

  CreateDeploymentResponse._();

  factory CreateDeploymentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateDeploymentResponse()..mergeFromBuffer(data, registry);
  factory CreateDeploymentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CreateDeploymentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDeploymentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CreateDeploymentResponse.$_createMessage)
    ..aOM<Deployment>(1, _omitFieldNames ? '' : 'deployment',
        subBuilder: Deployment.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentResponse copyWith(
          void Function(CreateDeploymentResponse) updates) =>
      super.copyWith((message) => updates(message as CreateDeploymentResponse))
          as CreateDeploymentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CreateDeploymentResponse() / CreateDeploymentResponse.new instead')
  static CreateDeploymentResponse create() => CreateDeploymentResponse._();
  static $pb.GeneratedMessage $_createMessage() => CreateDeploymentResponse._();
  @$core.override
  CreateDeploymentResponse createEmptyInstance() =>
      CreateDeploymentResponse._();
  @$core.pragma('dart2js:noInline')
  static CreateDeploymentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDeploymentResponse>(
          CreateDeploymentResponse.$_createMessage);
  static CreateDeploymentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Deployment get deployment => $_getN(0);
  @$pb.TagNumber(1)
  set deployment(Deployment value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDeployment() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeployment() => $_clearField(1);
  @$pb.TagNumber(1)
  Deployment ensureDeployment() => $_ensure(0);
}

class GetDeploymentRequest extends $pb.GeneratedMessage {
  factory GetDeploymentRequest({
    $core.String? deploymentId,
  }) {
    final result = GetDeploymentRequest._();
    if (deploymentId != null) result.deploymentId = deploymentId;
    return result;
  }

  GetDeploymentRequest._();

  factory GetDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetDeploymentRequest()..mergeFromBuffer(data, registry);
  factory GetDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetDeploymentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetDeploymentRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deploymentId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentRequest copyWith(void Function(GetDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as GetDeploymentRequest))
          as GetDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetDeploymentRequest() / GetDeploymentRequest.new instead')
  static GetDeploymentRequest create() => GetDeploymentRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetDeploymentRequest._();
  @$core.override
  GetDeploymentRequest createEmptyInstance() => GetDeploymentRequest._();
  @$core.pragma('dart2js:noInline')
  static GetDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDeploymentRequest>(
          GetDeploymentRequest.$_createMessage);
  static GetDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deploymentId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deploymentId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeploymentId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeploymentId() => $_clearField(1);
}

class GetDeploymentResponse extends $pb.GeneratedMessage {
  factory GetDeploymentResponse({
    Deployment? deployment,
  }) {
    final result = GetDeploymentResponse._();
    if (deployment != null) result.deployment = deployment;
    return result;
  }

  GetDeploymentResponse._();

  factory GetDeploymentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetDeploymentResponse()..mergeFromBuffer(data, registry);
  factory GetDeploymentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetDeploymentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDeploymentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetDeploymentResponse.$_createMessage)
    ..aOM<Deployment>(1, _omitFieldNames ? '' : 'deployment',
        subBuilder: Deployment.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentResponse copyWith(
          void Function(GetDeploymentResponse) updates) =>
      super.copyWith((message) => updates(message as GetDeploymentResponse))
          as GetDeploymentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetDeploymentResponse() / GetDeploymentResponse.new instead')
  static GetDeploymentResponse create() => GetDeploymentResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetDeploymentResponse._();
  @$core.override
  GetDeploymentResponse createEmptyInstance() => GetDeploymentResponse._();
  @$core.pragma('dart2js:noInline')
  static GetDeploymentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDeploymentResponse>(
          GetDeploymentResponse.$_createMessage);
  static GetDeploymentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Deployment get deployment => $_getN(0);
  @$pb.TagNumber(1)
  set deployment(Deployment value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDeployment() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeployment() => $_clearField(1);
  @$pb.TagNumber(1)
  Deployment ensureDeployment() => $_ensure(0);
}

class ListDeploymentsRequest extends $pb.GeneratedMessage {
  factory ListDeploymentsRequest({
    $core.String? projectId,
    $core.int? pageSize,
    $core.String? pageToken,
  }) {
    final result = ListDeploymentsRequest._();
    if (projectId != null) result.projectId = projectId;
    if (pageSize != null) result.pageSize = pageSize;
    if (pageToken != null) result.pageToken = pageToken;
    return result;
  }

  ListDeploymentsRequest._();

  factory ListDeploymentsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDeploymentsRequest()..mergeFromBuffer(data, registry);
  factory ListDeploymentsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDeploymentsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDeploymentsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListDeploymentsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'projectId')
    ..aI(2, _omitFieldNames ? '' : 'pageSize')
    ..aOS(3, _omitFieldNames ? '' : 'pageToken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeploymentsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeploymentsRequest copyWith(
          void Function(ListDeploymentsRequest) updates) =>
      super.copyWith((message) => updates(message as ListDeploymentsRequest))
          as ListDeploymentsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListDeploymentsRequest() / ListDeploymentsRequest.new instead')
  static ListDeploymentsRequest create() => ListDeploymentsRequest._();
  static $pb.GeneratedMessage $_createMessage() => ListDeploymentsRequest._();
  @$core.override
  ListDeploymentsRequest createEmptyInstance() => ListDeploymentsRequest._();
  @$core.pragma('dart2js:noInline')
  static ListDeploymentsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDeploymentsRequest>(
          ListDeploymentsRequest.$_createMessage);
  static ListDeploymentsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get projectId => $_getSZ(0);
  @$pb.TagNumber(1)
  set projectId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasProjectId() => $_has(0);
  @$pb.TagNumber(1)
  void clearProjectId() => $_clearField(1);

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

class ListDeploymentsResponse extends $pb.GeneratedMessage {
  factory ListDeploymentsResponse({
    $core.Iterable<Deployment>? deployments,
    $core.String? nextPageToken,
  }) {
    final result = ListDeploymentsResponse._();
    if (deployments != null) result.deployments.addAll(deployments);
    if (nextPageToken != null) result.nextPageToken = nextPageToken;
    return result;
  }

  ListDeploymentsResponse._();

  factory ListDeploymentsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDeploymentsResponse()..mergeFromBuffer(data, registry);
  factory ListDeploymentsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ListDeploymentsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDeploymentsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ListDeploymentsResponse.$_createMessage)
    ..pPM<Deployment>(1, _omitFieldNames ? '' : 'deployments',
        subBuilder: Deployment.$_createMessage)
    ..aOS(2, _omitFieldNames ? '' : 'nextPageToken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeploymentsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeploymentsResponse copyWith(
          void Function(ListDeploymentsResponse) updates) =>
      super.copyWith((message) => updates(message as ListDeploymentsResponse))
          as ListDeploymentsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use ListDeploymentsResponse() / ListDeploymentsResponse.new instead')
  static ListDeploymentsResponse create() => ListDeploymentsResponse._();
  static $pb.GeneratedMessage $_createMessage() => ListDeploymentsResponse._();
  @$core.override
  ListDeploymentsResponse createEmptyInstance() => ListDeploymentsResponse._();
  @$core.pragma('dart2js:noInline')
  static ListDeploymentsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDeploymentsResponse>(
          ListDeploymentsResponse.$_createMessage);
  static ListDeploymentsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Deployment> get deployments => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get nextPageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set nextPageToken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasNextPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearNextPageToken() => $_clearField(2);
}

class CancelDeploymentRequest extends $pb.GeneratedMessage {
  factory CancelDeploymentRequest({
    $core.String? deploymentId,
    $core.String? reason,
  }) {
    final result = CancelDeploymentRequest._();
    if (deploymentId != null) result.deploymentId = deploymentId;
    if (reason != null) result.reason = reason;
    return result;
  }

  CancelDeploymentRequest._();

  factory CancelDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CancelDeploymentRequest()..mergeFromBuffer(data, registry);
  factory CancelDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CancelDeploymentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CancelDeploymentRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deploymentId')
    ..aOS(2, _omitFieldNames ? '' : 'reason')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelDeploymentRequest copyWith(
          void Function(CancelDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as CancelDeploymentRequest))
          as CancelDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CancelDeploymentRequest() / CancelDeploymentRequest.new instead')
  static CancelDeploymentRequest create() => CancelDeploymentRequest._();
  static $pb.GeneratedMessage $_createMessage() => CancelDeploymentRequest._();
  @$core.override
  CancelDeploymentRequest createEmptyInstance() => CancelDeploymentRequest._();
  @$core.pragma('dart2js:noInline')
  static CancelDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelDeploymentRequest>(
          CancelDeploymentRequest.$_createMessage);
  static CancelDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deploymentId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deploymentId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeploymentId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeploymentId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get reason => $_getSZ(1);
  @$pb.TagNumber(2)
  set reason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReason() => $_has(1);
  @$pb.TagNumber(2)
  void clearReason() => $_clearField(2);
}

class CancelDeploymentResponse extends $pb.GeneratedMessage {
  factory CancelDeploymentResponse({
    Deployment? deployment,
  }) {
    final result = CancelDeploymentResponse._();
    if (deployment != null) result.deployment = deployment;
    return result;
  }

  CancelDeploymentResponse._();

  factory CancelDeploymentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CancelDeploymentResponse()..mergeFromBuffer(data, registry);
  factory CancelDeploymentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CancelDeploymentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelDeploymentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CancelDeploymentResponse.$_createMessage)
    ..aOM<Deployment>(1, _omitFieldNames ? '' : 'deployment',
        subBuilder: Deployment.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelDeploymentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelDeploymentResponse copyWith(
          void Function(CancelDeploymentResponse) updates) =>
      super.copyWith((message) => updates(message as CancelDeploymentResponse))
          as CancelDeploymentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use CancelDeploymentResponse() / CancelDeploymentResponse.new instead')
  static CancelDeploymentResponse create() => CancelDeploymentResponse._();
  static $pb.GeneratedMessage $_createMessage() => CancelDeploymentResponse._();
  @$core.override
  CancelDeploymentResponse createEmptyInstance() =>
      CancelDeploymentResponse._();
  @$core.pragma('dart2js:noInline')
  static CancelDeploymentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelDeploymentResponse>(
          CancelDeploymentResponse.$_createMessage);
  static CancelDeploymentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Deployment get deployment => $_getN(0);
  @$pb.TagNumber(1)
  set deployment(Deployment value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDeployment() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeployment() => $_clearField(1);
  @$pb.TagNumber(1)
  Deployment ensureDeployment() => $_ensure(0);
}

class RollbackDeploymentRequest extends $pb.GeneratedMessage {
  factory RollbackDeploymentRequest({
    $core.String? targetDeploymentId,
    $core.String? requestId,
  }) {
    final result = RollbackDeploymentRequest._();
    if (targetDeploymentId != null)
      result.targetDeploymentId = targetDeploymentId;
    if (requestId != null) result.requestId = requestId;
    return result;
  }

  RollbackDeploymentRequest._();

  factory RollbackDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollbackDeploymentRequest()..mergeFromBuffer(data, registry);
  factory RollbackDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollbackDeploymentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RollbackDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RollbackDeploymentRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'targetDeploymentId')
    ..aOS(2, _omitFieldNames ? '' : 'requestId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollbackDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollbackDeploymentRequest copyWith(
          void Function(RollbackDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as RollbackDeploymentRequest))
          as RollbackDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RollbackDeploymentRequest() / RollbackDeploymentRequest.new instead')
  static RollbackDeploymentRequest create() => RollbackDeploymentRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      RollbackDeploymentRequest._();
  @$core.override
  RollbackDeploymentRequest createEmptyInstance() =>
      RollbackDeploymentRequest._();
  @$core.pragma('dart2js:noInline')
  static RollbackDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RollbackDeploymentRequest>(
          RollbackDeploymentRequest.$_createMessage);
  static RollbackDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get targetDeploymentId => $_getSZ(0);
  @$pb.TagNumber(1)
  set targetDeploymentId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTargetDeploymentId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTargetDeploymentId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get requestId => $_getSZ(1);
  @$pb.TagNumber(2)
  set requestId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRequestId() => $_has(1);
  @$pb.TagNumber(2)
  void clearRequestId() => $_clearField(2);
}

class RollbackDeploymentResponse extends $pb.GeneratedMessage {
  factory RollbackDeploymentResponse({
    Deployment? deployment,
  }) {
    final result = RollbackDeploymentResponse._();
    if (deployment != null) result.deployment = deployment;
    return result;
  }

  RollbackDeploymentResponse._();

  factory RollbackDeploymentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollbackDeploymentResponse()..mergeFromBuffer(data, registry);
  factory RollbackDeploymentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      RollbackDeploymentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RollbackDeploymentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: RollbackDeploymentResponse.$_createMessage)
    ..aOM<Deployment>(1, _omitFieldNames ? '' : 'deployment',
        subBuilder: Deployment.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollbackDeploymentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RollbackDeploymentResponse copyWith(
          void Function(RollbackDeploymentResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RollbackDeploymentResponse))
          as RollbackDeploymentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use RollbackDeploymentResponse() / RollbackDeploymentResponse.new instead')
  static RollbackDeploymentResponse create() => RollbackDeploymentResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      RollbackDeploymentResponse._();
  @$core.override
  RollbackDeploymentResponse createEmptyInstance() =>
      RollbackDeploymentResponse._();
  @$core.pragma('dart2js:noInline')
  static RollbackDeploymentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RollbackDeploymentResponse>(
          RollbackDeploymentResponse.$_createMessage);
  static RollbackDeploymentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Deployment get deployment => $_getN(0);
  @$pb.TagNumber(1)
  set deployment(Deployment value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDeployment() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeployment() => $_clearField(1);
  @$pb.TagNumber(1)
  Deployment ensureDeployment() => $_ensure(0);
}

class StreamDeploymentLogsRequest extends $pb.GeneratedMessage {
  factory StreamDeploymentLogsRequest({
    $core.String? deploymentId,
    $core.bool? follow,
    $1.Timestamp? since,
    $2.LogStream? stream,
  }) {
    final result = StreamDeploymentLogsRequest._();
    if (deploymentId != null) result.deploymentId = deploymentId;
    if (follow != null) result.follow = follow;
    if (since != null) result.since = since;
    if (stream != null) result.stream = stream;
    return result;
  }

  StreamDeploymentLogsRequest._();

  factory StreamDeploymentLogsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentLogsRequest()..mergeFromBuffer(data, registry);
  factory StreamDeploymentLogsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentLogsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDeploymentLogsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamDeploymentLogsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deploymentId')
    ..aOB(2, _omitFieldNames ? '' : 'follow')
    ..aOM<$1.Timestamp>(3, _omitFieldNames ? '' : 'since',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aE<$2.LogStream>(4, _omitFieldNames ? '' : 'stream',
        enumValues: $2.LogStream.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentLogsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentLogsRequest copyWith(
          void Function(StreamDeploymentLogsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StreamDeploymentLogsRequest))
          as StreamDeploymentLogsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamDeploymentLogsRequest() / StreamDeploymentLogsRequest.new instead')
  static StreamDeploymentLogsRequest create() =>
      StreamDeploymentLogsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamDeploymentLogsRequest._();
  @$core.override
  StreamDeploymentLogsRequest createEmptyInstance() =>
      StreamDeploymentLogsRequest._();
  @$core.pragma('dart2js:noInline')
  static StreamDeploymentLogsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDeploymentLogsRequest>(
          StreamDeploymentLogsRequest.$_createMessage);
  static StreamDeploymentLogsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deploymentId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deploymentId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeploymentId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeploymentId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get follow => $_getBF(1);
  @$pb.TagNumber(2)
  set follow($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasFollow() => $_has(1);
  @$pb.TagNumber(2)
  void clearFollow() => $_clearField(2);

  @$pb.TagNumber(3)
  $1.Timestamp get since => $_getN(2);
  @$pb.TagNumber(3)
  set since($1.Timestamp value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasSince() => $_has(2);
  @$pb.TagNumber(3)
  void clearSince() => $_clearField(3);
  @$pb.TagNumber(3)
  $1.Timestamp ensureSince() => $_ensure(2);

  /// Filtrage par flux de sortie.
  /// Si la valeur est LOG_STREAM_UNSPECIFIED, les deux flux (STDOUT et STDERR) sont fusionnés et renvoyés.
  @$pb.TagNumber(4)
  $2.LogStream get stream => $_getN(3);
  @$pb.TagNumber(4)
  set stream($2.LogStream value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasStream() => $_has(3);
  @$pb.TagNumber(4)
  void clearStream() => $_clearField(4);
}

class StreamDeploymentLogsResponse extends $pb.GeneratedMessage {
  factory StreamDeploymentLogsResponse({
    DeploymentLogEntry? logEntry,
  }) {
    final result = StreamDeploymentLogsResponse._();
    if (logEntry != null) result.logEntry = logEntry;
    return result;
  }

  StreamDeploymentLogsResponse._();

  factory StreamDeploymentLogsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentLogsResponse()..mergeFromBuffer(data, registry);
  factory StreamDeploymentLogsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentLogsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDeploymentLogsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamDeploymentLogsResponse.$_createMessage)
    ..aOM<DeploymentLogEntry>(1, _omitFieldNames ? '' : 'logEntry',
        subBuilder: DeploymentLogEntry.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentLogsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentLogsResponse copyWith(
          void Function(StreamDeploymentLogsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StreamDeploymentLogsResponse))
          as StreamDeploymentLogsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamDeploymentLogsResponse() / StreamDeploymentLogsResponse.new instead')
  static StreamDeploymentLogsResponse create() =>
      StreamDeploymentLogsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamDeploymentLogsResponse._();
  @$core.override
  StreamDeploymentLogsResponse createEmptyInstance() =>
      StreamDeploymentLogsResponse._();
  @$core.pragma('dart2js:noInline')
  static StreamDeploymentLogsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDeploymentLogsResponse>(
          StreamDeploymentLogsResponse.$_createMessage);
  static StreamDeploymentLogsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  DeploymentLogEntry get logEntry => $_getN(0);
  @$pb.TagNumber(1)
  set logEntry(DeploymentLogEntry value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasLogEntry() => $_has(0);
  @$pb.TagNumber(1)
  void clearLogEntry() => $_clearField(1);
  @$pb.TagNumber(1)
  DeploymentLogEntry ensureLogEntry() => $_ensure(0);
}

class StreamDeploymentRequest extends $pb.GeneratedMessage {
  factory StreamDeploymentRequest({
    $core.String? deploymentId,
  }) {
    final result = StreamDeploymentRequest._();
    if (deploymentId != null) result.deploymentId = deploymentId;
    return result;
  }

  StreamDeploymentRequest._();

  factory StreamDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentRequest()..mergeFromBuffer(data, registry);
  factory StreamDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamDeploymentRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'deploymentId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentRequest copyWith(
          void Function(StreamDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as StreamDeploymentRequest))
          as StreamDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamDeploymentRequest() / StreamDeploymentRequest.new instead')
  static StreamDeploymentRequest create() => StreamDeploymentRequest._();
  static $pb.GeneratedMessage $_createMessage() => StreamDeploymentRequest._();
  @$core.override
  StreamDeploymentRequest createEmptyInstance() => StreamDeploymentRequest._();
  @$core.pragma('dart2js:noInline')
  static StreamDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDeploymentRequest>(
          StreamDeploymentRequest.$_createMessage);
  static StreamDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deploymentId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deploymentId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeploymentId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeploymentId() => $_clearField(1);
}

class StreamDeploymentResponse extends $pb.GeneratedMessage {
  factory StreamDeploymentResponse({
    Deployment? deployment,
  }) {
    final result = StreamDeploymentResponse._();
    if (deployment != null) result.deployment = deployment;
    return result;
  }

  StreamDeploymentResponse._();

  factory StreamDeploymentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentResponse()..mergeFromBuffer(data, registry);
  factory StreamDeploymentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamDeploymentResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDeploymentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamDeploymentResponse.$_createMessage)
    ..aOM<Deployment>(1, _omitFieldNames ? '' : 'deployment',
        subBuilder: Deployment.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDeploymentResponse copyWith(
          void Function(StreamDeploymentResponse) updates) =>
      super.copyWith((message) => updates(message as StreamDeploymentResponse))
          as StreamDeploymentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamDeploymentResponse() / StreamDeploymentResponse.new instead')
  static StreamDeploymentResponse create() => StreamDeploymentResponse._();
  static $pb.GeneratedMessage $_createMessage() => StreamDeploymentResponse._();
  @$core.override
  StreamDeploymentResponse createEmptyInstance() =>
      StreamDeploymentResponse._();
  @$core.pragma('dart2js:noInline')
  static StreamDeploymentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDeploymentResponse>(
          StreamDeploymentResponse.$_createMessage);
  static StreamDeploymentResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Deployment get deployment => $_getN(0);
  @$pb.TagNumber(1)
  set deployment(Deployment value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasDeployment() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeployment() => $_clearField(1);
  @$pb.TagNumber(1)
  Deployment ensureDeployment() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
