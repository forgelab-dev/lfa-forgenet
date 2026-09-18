// This is a generated file - do not edit.
//
// Generated from forgenet/v1/metrics.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;
import 'package:protobuf/well_known_types/google/protobuf/timestamp.pb.dart'
    as $1;

import 'metrics.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'metrics.pbenum.dart';

/// TimeRange définit une fenêtre temporelle pour les requêtes d'historique.
class TimeRange extends $pb.GeneratedMessage {
  factory TimeRange({
    $1.Timestamp? startTime,
    $1.Timestamp? endTime,
  }) {
    final result = TimeRange._();
    if (startTime != null) result.startTime = startTime;
    if (endTime != null) result.endTime = endTime;
    return result;
  }

  TimeRange._();

  factory TimeRange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TimeRange()..mergeFromBuffer(data, registry);
  factory TimeRange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      TimeRange()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimeRange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: TimeRange.$_createMessage)
    ..aOM<$1.Timestamp>(1, _omitFieldNames ? '' : 'startTime',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<$1.Timestamp>(2, _omitFieldNames ? '' : 'endTime',
        subBuilder: $1.Timestamp.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeRange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeRange copyWith(void Function(TimeRange) updates) =>
      super.copyWith((message) => updates(message as TimeRange)) as TimeRange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use TimeRange() / TimeRange.new instead')
  static TimeRange create() => TimeRange._();
  static $pb.GeneratedMessage $_createMessage() => TimeRange._();
  @$core.override
  TimeRange createEmptyInstance() => TimeRange._();
  @$core.pragma('dart2js:noInline')
  static TimeRange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimeRange>(TimeRange.$_createMessage);
  static TimeRange? _defaultInstance;

  @$pb.TagNumber(1)
  $1.Timestamp get startTime => $_getN(0);
  @$pb.TagNumber(1)
  set startTime($1.Timestamp value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasStartTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearStartTime() => $_clearField(1);
  @$pb.TagNumber(1)
  $1.Timestamp ensureStartTime() => $_ensure(0);

  @$pb.TagNumber(2)
  $1.Timestamp get endTime => $_getN(1);
  @$pb.TagNumber(2)
  set endTime($1.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasEndTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearEndTime() => $_clearField(2);
  @$pb.TagNumber(2)
  $1.Timestamp ensureEndTime() => $_ensure(1);
}

/// MetricSample représente un point de mesure unique dans une série temporelle.
class MetricSample extends $pb.GeneratedMessage {
  factory MetricSample({
    $1.Timestamp? timestamp,
    $core.double? value,
  }) {
    final result = MetricSample._();
    if (timestamp != null) result.timestamp = timestamp;
    if (value != null) result.value = value;
    return result;
  }

  MetricSample._();

  factory MetricSample.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MetricSample()..mergeFromBuffer(data, registry);
  factory MetricSample.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MetricSample()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricSample',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: MetricSample.$_createMessage)
    ..aOM<$1.Timestamp>(1, _omitFieldNames ? '' : 'timestamp',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aD(2, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricSample clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricSample copyWith(void Function(MetricSample) updates) =>
      super.copyWith((message) => updates(message as MetricSample))
          as MetricSample;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use MetricSample() / MetricSample.new instead')
  static MetricSample create() => MetricSample._();
  static $pb.GeneratedMessage $_createMessage() => MetricSample._();
  @$core.override
  MetricSample createEmptyInstance() => MetricSample._();
  @$core.pragma('dart2js:noInline')
  static MetricSample getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MetricSample>(
          MetricSample.$_createMessage);
  static MetricSample? _defaultInstance;

  @$pb.TagNumber(1)
  $1.Timestamp get timestamp => $_getN(0);
  @$pb.TagNumber(1)
  set timestamp($1.Timestamp value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasTimestamp() => $_has(0);
  @$pb.TagNumber(1)
  void clearTimestamp() => $_clearField(1);
  @$pb.TagNumber(1)
  $1.Timestamp ensureTimestamp() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.double get value => $_getN(1);
  @$pb.TagNumber(2)
  set value($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(2)
  void clearValue() => $_clearField(2);
}

/// MetricSeries regroupe une série de mesures chronologiques pour une métrique donnée.
class MetricSeries extends $pb.GeneratedMessage {
  factory MetricSeries({
    MetricName? metricName,
    MetricUnit? unit,
    $core.Iterable<MetricSample>? samples,
  }) {
    final result = MetricSeries._();
    if (metricName != null) result.metricName = metricName;
    if (unit != null) result.unit = unit;
    if (samples != null) result.samples.addAll(samples);
    return result;
  }

  MetricSeries._();

  factory MetricSeries.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MetricSeries()..mergeFromBuffer(data, registry);
  factory MetricSeries.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MetricSeries()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricSeries',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: MetricSeries.$_createMessage)
    ..aE<MetricName>(1, _omitFieldNames ? '' : 'metricName',
        enumValues: MetricName.values)
    ..aE<MetricUnit>(2, _omitFieldNames ? '' : 'unit',
        enumValues: MetricUnit.values)
    ..pPM<MetricSample>(3, _omitFieldNames ? '' : 'samples',
        subBuilder: MetricSample.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricSeries clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricSeries copyWith(void Function(MetricSeries) updates) =>
      super.copyWith((message) => updates(message as MetricSeries))
          as MetricSeries;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use MetricSeries() / MetricSeries.new instead')
  static MetricSeries create() => MetricSeries._();
  static $pb.GeneratedMessage $_createMessage() => MetricSeries._();
  @$core.override
  MetricSeries createEmptyInstance() => MetricSeries._();
  @$core.pragma('dart2js:noInline')
  static MetricSeries getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MetricSeries>(
          MetricSeries.$_createMessage);
  static MetricSeries? _defaultInstance;

  @$pb.TagNumber(1)
  MetricName get metricName => $_getN(0);
  @$pb.TagNumber(1)
  set metricName(MetricName value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMetricName() => $_has(0);
  @$pb.TagNumber(1)
  void clearMetricName() => $_clearField(1);

  @$pb.TagNumber(2)
  MetricUnit get unit => $_getN(1);
  @$pb.TagNumber(2)
  set unit(MetricUnit value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasUnit() => $_has(1);
  @$pb.TagNumber(2)
  void clearUnit() => $_clearField(2);

  @$pb.TagNumber(3)
  $pb.PbList<MetricSample> get samples => $_getList(2);
}

/// CpuMetrics définit la consommation de ressources processeur.
class CpuMetrics extends $pb.GeneratedMessage {
  factory CpuMetrics({
    $core.double? cpuUsagePercent,
    $core.int? onlineCpus,
    $fixnum.Int64? systemCpuNanoseconds,
  }) {
    final result = CpuMetrics._();
    if (cpuUsagePercent != null) result.cpuUsagePercent = cpuUsagePercent;
    if (onlineCpus != null) result.onlineCpus = onlineCpus;
    if (systemCpuNanoseconds != null)
      result.systemCpuNanoseconds = systemCpuNanoseconds;
    return result;
  }

  CpuMetrics._();

  factory CpuMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CpuMetrics()..mergeFromBuffer(data, registry);
  factory CpuMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      CpuMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CpuMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: CpuMetrics.$_createMessage)
    ..aD(1, _omitFieldNames ? '' : 'cpuUsagePercent')
    ..aI(2, _omitFieldNames ? '' : 'onlineCpus', fieldType: $pb.PbFieldType.OU3)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'systemCpuNanoseconds', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CpuMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CpuMetrics copyWith(void Function(CpuMetrics) updates) =>
      super.copyWith((message) => updates(message as CpuMetrics)) as CpuMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use CpuMetrics() / CpuMetrics.new instead')
  static CpuMetrics create() => CpuMetrics._();
  static $pb.GeneratedMessage $_createMessage() => CpuMetrics._();
  @$core.override
  CpuMetrics createEmptyInstance() => CpuMetrics._();
  @$core.pragma('dart2js:noInline')
  static CpuMetrics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CpuMetrics>(CpuMetrics.$_createMessage);
  static CpuMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $core.double get cpuUsagePercent => $_getN(0);
  @$pb.TagNumber(1)
  set cpuUsagePercent($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCpuUsagePercent() => $_has(0);
  @$pb.TagNumber(1)
  void clearCpuUsagePercent() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get onlineCpus => $_getIZ(1);
  @$pb.TagNumber(2)
  set onlineCpus($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOnlineCpus() => $_has(1);
  @$pb.TagNumber(2)
  void clearOnlineCpus() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get systemCpuNanoseconds => $_getI64(2);
  @$pb.TagNumber(3)
  set systemCpuNanoseconds($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasSystemCpuNanoseconds() => $_has(2);
  @$pb.TagNumber(3)
  void clearSystemCpuNanoseconds() => $_clearField(3);
}

/// MemoryMetrics définit l'utilisation de la mémoire RAM et du cache.
class MemoryMetrics extends $pb.GeneratedMessage {
  factory MemoryMetrics({
    $fixnum.Int64? usageBytes,
    $fixnum.Int64? limitBytes,
    $fixnum.Int64? cacheBytes,
    $core.double? usagePercent,
  }) {
    final result = MemoryMetrics._();
    if (usageBytes != null) result.usageBytes = usageBytes;
    if (limitBytes != null) result.limitBytes = limitBytes;
    if (cacheBytes != null) result.cacheBytes = cacheBytes;
    if (usagePercent != null) result.usagePercent = usagePercent;
    return result;
  }

  MemoryMetrics._();

  factory MemoryMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MemoryMetrics()..mergeFromBuffer(data, registry);
  factory MemoryMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      MemoryMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MemoryMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: MemoryMetrics.$_createMessage)
    ..a<$fixnum.Int64>(
        1, _omitFieldNames ? '' : 'usageBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'limitBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'cacheBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aD(4, _omitFieldNames ? '' : 'usagePercent')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MemoryMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MemoryMetrics copyWith(void Function(MemoryMetrics) updates) =>
      super.copyWith((message) => updates(message as MemoryMetrics))
          as MemoryMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use MemoryMetrics() / MemoryMetrics.new instead')
  static MemoryMetrics create() => MemoryMetrics._();
  static $pb.GeneratedMessage $_createMessage() => MemoryMetrics._();
  @$core.override
  MemoryMetrics createEmptyInstance() => MemoryMetrics._();
  @$core.pragma('dart2js:noInline')
  static MemoryMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MemoryMetrics>(
          MemoryMetrics.$_createMessage);
  static MemoryMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get usageBytes => $_getI64(0);
  @$pb.TagNumber(1)
  set usageBytes($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUsageBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearUsageBytes() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get limitBytes => $_getI64(1);
  @$pb.TagNumber(2)
  set limitBytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimitBytes() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimitBytes() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get cacheBytes => $_getI64(2);
  @$pb.TagNumber(3)
  set cacheBytes($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCacheBytes() => $_has(2);
  @$pb.TagNumber(3)
  void clearCacheBytes() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get usagePercent => $_getN(3);
  @$pb.TagNumber(4)
  set usagePercent($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasUsagePercent() => $_has(3);
  @$pb.TagNumber(4)
  void clearUsagePercent() => $_clearField(4);
}

/// NetworkMetrics détaille le trafic réseau cumulé (compteurs Prometheus).
class NetworkMetrics extends $pb.GeneratedMessage {
  factory NetworkMetrics({
    $fixnum.Int64? rxBytesTotal,
    $fixnum.Int64? txBytesTotal,
    $fixnum.Int64? rxPacketsTotal,
    $fixnum.Int64? txPacketsTotal,
    $fixnum.Int64? rxErrorsTotal,
    $fixnum.Int64? txErrorsTotal,
  }) {
    final result = NetworkMetrics._();
    if (rxBytesTotal != null) result.rxBytesTotal = rxBytesTotal;
    if (txBytesTotal != null) result.txBytesTotal = txBytesTotal;
    if (rxPacketsTotal != null) result.rxPacketsTotal = rxPacketsTotal;
    if (txPacketsTotal != null) result.txPacketsTotal = txPacketsTotal;
    if (rxErrorsTotal != null) result.rxErrorsTotal = rxErrorsTotal;
    if (txErrorsTotal != null) result.txErrorsTotal = txErrorsTotal;
    return result;
  }

  NetworkMetrics._();

  factory NetworkMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkMetrics()..mergeFromBuffer(data, registry);
  factory NetworkMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      NetworkMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: NetworkMetrics.$_createMessage)
    ..a<$fixnum.Int64>(
        1, _omitFieldNames ? '' : 'rxBytesTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'txBytesTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'rxPacketsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'txPacketsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        5, _omitFieldNames ? '' : 'rxErrorsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        6, _omitFieldNames ? '' : 'txErrorsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkMetrics copyWith(void Function(NetworkMetrics) updates) =>
      super.copyWith((message) => updates(message as NetworkMetrics))
          as NetworkMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use NetworkMetrics() / NetworkMetrics.new instead')
  static NetworkMetrics create() => NetworkMetrics._();
  static $pb.GeneratedMessage $_createMessage() => NetworkMetrics._();
  @$core.override
  NetworkMetrics createEmptyInstance() => NetworkMetrics._();
  @$core.pragma('dart2js:noInline')
  static NetworkMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<NetworkMetrics>(
          NetworkMetrics.$_createMessage);
  static NetworkMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get rxBytesTotal => $_getI64(0);
  @$pb.TagNumber(1)
  set rxBytesTotal($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRxBytesTotal() => $_has(0);
  @$pb.TagNumber(1)
  void clearRxBytesTotal() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get txBytesTotal => $_getI64(1);
  @$pb.TagNumber(2)
  set txBytesTotal($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTxBytesTotal() => $_has(1);
  @$pb.TagNumber(2)
  void clearTxBytesTotal() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get rxPacketsTotal => $_getI64(2);
  @$pb.TagNumber(3)
  set rxPacketsTotal($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRxPacketsTotal() => $_has(2);
  @$pb.TagNumber(3)
  void clearRxPacketsTotal() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get txPacketsTotal => $_getI64(3);
  @$pb.TagNumber(4)
  set txPacketsTotal($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasTxPacketsTotal() => $_has(3);
  @$pb.TagNumber(4)
  void clearTxPacketsTotal() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get rxErrorsTotal => $_getI64(4);
  @$pb.TagNumber(5)
  set rxErrorsTotal($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasRxErrorsTotal() => $_has(4);
  @$pb.TagNumber(5)
  void clearRxErrorsTotal() => $_clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get txErrorsTotal => $_getI64(5);
  @$pb.TagNumber(6)
  set txErrorsTotal($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasTxErrorsTotal() => $_has(5);
  @$pb.TagNumber(6)
  void clearTxErrorsTotal() => $_clearField(6);
}

/// DiskIoMetrics détaille les opérations d'E/S sur le disque (compteurs Prometheus).
class DiskIoMetrics extends $pb.GeneratedMessage {
  factory DiskIoMetrics({
    $fixnum.Int64? readBytesTotal,
    $fixnum.Int64? writeBytesTotal,
    $fixnum.Int64? readOpsTotal,
    $fixnum.Int64? writeOpsTotal,
  }) {
    final result = DiskIoMetrics._();
    if (readBytesTotal != null) result.readBytesTotal = readBytesTotal;
    if (writeBytesTotal != null) result.writeBytesTotal = writeBytesTotal;
    if (readOpsTotal != null) result.readOpsTotal = readOpsTotal;
    if (writeOpsTotal != null) result.writeOpsTotal = writeOpsTotal;
    return result;
  }

  DiskIoMetrics._();

  factory DiskIoMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DiskIoMetrics()..mergeFromBuffer(data, registry);
  factory DiskIoMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      DiskIoMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DiskIoMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: DiskIoMetrics.$_createMessage)
    ..a<$fixnum.Int64>(
        1, _omitFieldNames ? '' : 'readBytesTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'writeBytesTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'readOpsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'writeOpsTotal', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiskIoMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DiskIoMetrics copyWith(void Function(DiskIoMetrics) updates) =>
      super.copyWith((message) => updates(message as DiskIoMetrics))
          as DiskIoMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use DiskIoMetrics() / DiskIoMetrics.new instead')
  static DiskIoMetrics create() => DiskIoMetrics._();
  static $pb.GeneratedMessage $_createMessage() => DiskIoMetrics._();
  @$core.override
  DiskIoMetrics createEmptyInstance() => DiskIoMetrics._();
  @$core.pragma('dart2js:noInline')
  static DiskIoMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DiskIoMetrics>(
          DiskIoMetrics.$_createMessage);
  static DiskIoMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get readBytesTotal => $_getI64(0);
  @$pb.TagNumber(1)
  set readBytesTotal($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasReadBytesTotal() => $_has(0);
  @$pb.TagNumber(1)
  void clearReadBytesTotal() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get writeBytesTotal => $_getI64(1);
  @$pb.TagNumber(2)
  set writeBytesTotal($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasWriteBytesTotal() => $_has(1);
  @$pb.TagNumber(2)
  void clearWriteBytesTotal() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get readOpsTotal => $_getI64(2);
  @$pb.TagNumber(3)
  set readOpsTotal($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasReadOpsTotal() => $_has(2);
  @$pb.TagNumber(3)
  void clearReadOpsTotal() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get writeOpsTotal => $_getI64(3);
  @$pb.TagNumber(4)
  set writeOpsTotal($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasWriteOpsTotal() => $_has(3);
  @$pb.TagNumber(4)
  void clearWriteOpsTotal() => $_clearField(4);
}

/// StorageMetrics détaille l'occupation du stockage disque.
class StorageMetrics extends $pb.GeneratedMessage {
  factory StorageMetrics({
    $fixnum.Int64? totalBytes,
    $fixnum.Int64? usedBytes,
    $fixnum.Int64? freeBytes,
    $core.double? usagePercent,
  }) {
    final result = StorageMetrics._();
    if (totalBytes != null) result.totalBytes = totalBytes;
    if (usedBytes != null) result.usedBytes = usedBytes;
    if (freeBytes != null) result.freeBytes = freeBytes;
    if (usagePercent != null) result.usagePercent = usagePercent;
    return result;
  }

  StorageMetrics._();

  factory StorageMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StorageMetrics()..mergeFromBuffer(data, registry);
  factory StorageMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StorageMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StorageMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StorageMetrics.$_createMessage)
    ..a<$fixnum.Int64>(
        1, _omitFieldNames ? '' : 'totalBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'usedBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'freeBytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aD(4, _omitFieldNames ? '' : 'usagePercent')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StorageMetrics copyWith(void Function(StorageMetrics) updates) =>
      super.copyWith((message) => updates(message as StorageMetrics))
          as StorageMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use StorageMetrics() / StorageMetrics.new instead')
  static StorageMetrics create() => StorageMetrics._();
  static $pb.GeneratedMessage $_createMessage() => StorageMetrics._();
  @$core.override
  StorageMetrics createEmptyInstance() => StorageMetrics._();
  @$core.pragma('dart2js:noInline')
  static StorageMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<StorageMetrics>(
          StorageMetrics.$_createMessage);
  static StorageMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get totalBytes => $_getI64(0);
  @$pb.TagNumber(1)
  set totalBytes($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTotalBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalBytes() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get usedBytes => $_getI64(1);
  @$pb.TagNumber(2)
  set usedBytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUsedBytes() => $_has(1);
  @$pb.TagNumber(2)
  void clearUsedBytes() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get freeBytes => $_getI64(2);
  @$pb.TagNumber(3)
  set freeBytes($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasFreeBytes() => $_has(2);
  @$pb.TagNumber(3)
  void clearFreeBytes() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get usagePercent => $_getN(3);
  @$pb.TagNumber(4)
  set usagePercent($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasUsagePercent() => $_has(3);
  @$pb.TagNumber(4)
  void clearUsagePercent() => $_clearField(4);
}

/// ContainerMetrics regroupe l'ensemble des métriques d'un conteneur à un instant T.
class ContainerMetrics extends $pb.GeneratedMessage {
  factory ContainerMetrics({
    $core.String? containerId,
    $1.Timestamp? collectedAt,
    CpuMetrics? cpu,
    MemoryMetrics? memory,
    NetworkMetrics? network,
    DiskIoMetrics? diskIo,
  }) {
    final result = ContainerMetrics._();
    if (containerId != null) result.containerId = containerId;
    if (collectedAt != null) result.collectedAt = collectedAt;
    if (cpu != null) result.cpu = cpu;
    if (memory != null) result.memory = memory;
    if (network != null) result.network = network;
    if (diskIo != null) result.diskIo = diskIo;
    return result;
  }

  ContainerMetrics._();

  factory ContainerMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ContainerMetrics()..mergeFromBuffer(data, registry);
  factory ContainerMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      ContainerMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContainerMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: ContainerMetrics.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..aOM<$1.Timestamp>(2, _omitFieldNames ? '' : 'collectedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..aOM<CpuMetrics>(3, _omitFieldNames ? '' : 'cpu',
        subBuilder: CpuMetrics.$_createMessage)
    ..aOM<MemoryMetrics>(4, _omitFieldNames ? '' : 'memory',
        subBuilder: MemoryMetrics.$_createMessage)
    ..aOM<NetworkMetrics>(5, _omitFieldNames ? '' : 'network',
        subBuilder: NetworkMetrics.$_createMessage)
    ..aOM<DiskIoMetrics>(6, _omitFieldNames ? '' : 'diskIo',
        subBuilder: DiskIoMetrics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContainerMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContainerMetrics copyWith(void Function(ContainerMetrics) updates) =>
      super.copyWith((message) => updates(message as ContainerMetrics))
          as ContainerMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use ContainerMetrics() / ContainerMetrics.new instead')
  static ContainerMetrics create() => ContainerMetrics._();
  static $pb.GeneratedMessage $_createMessage() => ContainerMetrics._();
  @$core.override
  ContainerMetrics createEmptyInstance() => ContainerMetrics._();
  @$core.pragma('dart2js:noInline')
  static ContainerMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ContainerMetrics>(
          ContainerMetrics.$_createMessage);
  static ContainerMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $1.Timestamp get collectedAt => $_getN(1);
  @$pb.TagNumber(2)
  set collectedAt($1.Timestamp value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasCollectedAt() => $_has(1);
  @$pb.TagNumber(2)
  void clearCollectedAt() => $_clearField(2);
  @$pb.TagNumber(2)
  $1.Timestamp ensureCollectedAt() => $_ensure(1);

  @$pb.TagNumber(3)
  CpuMetrics get cpu => $_getN(2);
  @$pb.TagNumber(3)
  set cpu(CpuMetrics value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasCpu() => $_has(2);
  @$pb.TagNumber(3)
  void clearCpu() => $_clearField(3);
  @$pb.TagNumber(3)
  CpuMetrics ensureCpu() => $_ensure(2);

  @$pb.TagNumber(4)
  MemoryMetrics get memory => $_getN(3);
  @$pb.TagNumber(4)
  set memory(MemoryMetrics value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasMemory() => $_has(3);
  @$pb.TagNumber(4)
  void clearMemory() => $_clearField(4);
  @$pb.TagNumber(4)
  MemoryMetrics ensureMemory() => $_ensure(3);

  @$pb.TagNumber(5)
  NetworkMetrics get network => $_getN(4);
  @$pb.TagNumber(5)
  set network(NetworkMetrics value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasNetwork() => $_has(4);
  @$pb.TagNumber(5)
  void clearNetwork() => $_clearField(5);
  @$pb.TagNumber(5)
  NetworkMetrics ensureNetwork() => $_ensure(4);

  @$pb.TagNumber(6)
  DiskIoMetrics get diskIo => $_getN(5);
  @$pb.TagNumber(6)
  set diskIo(DiskIoMetrics value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasDiskIo() => $_has(5);
  @$pb.TagNumber(6)
  void clearDiskIo() => $_clearField(6);
  @$pb.TagNumber(6)
  DiskIoMetrics ensureDiskIo() => $_ensure(5);
}

/// SystemMetrics regroupe les métriques globales du nœud hôte.
class SystemMetrics extends $pb.GeneratedMessage {
  factory SystemMetrics({
    $1.Timestamp? collectedAt,
    $fixnum.Int64? uptimeSeconds,
    CpuMetrics? cpu,
    MemoryMetrics? memory,
    StorageMetrics? storage,
    NetworkMetrics? network,
    DiskIoMetrics? diskIo,
    $core.int? activeContainersCount,
  }) {
    final result = SystemMetrics._();
    if (collectedAt != null) result.collectedAt = collectedAt;
    if (uptimeSeconds != null) result.uptimeSeconds = uptimeSeconds;
    if (cpu != null) result.cpu = cpu;
    if (memory != null) result.memory = memory;
    if (storage != null) result.storage = storage;
    if (network != null) result.network = network;
    if (diskIo != null) result.diskIo = diskIo;
    if (activeContainersCount != null)
      result.activeContainersCount = activeContainersCount;
    return result;
  }

  SystemMetrics._();

  factory SystemMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SystemMetrics()..mergeFromBuffer(data, registry);
  factory SystemMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      SystemMetrics()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SystemMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: SystemMetrics.$_createMessage)
    ..aOM<$1.Timestamp>(1, _omitFieldNames ? '' : 'collectedAt',
        subBuilder: $1.Timestamp.$_createMessage)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'uptimeSeconds', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOM<CpuMetrics>(3, _omitFieldNames ? '' : 'cpu',
        subBuilder: CpuMetrics.$_createMessage)
    ..aOM<MemoryMetrics>(4, _omitFieldNames ? '' : 'memory',
        subBuilder: MemoryMetrics.$_createMessage)
    ..aOM<StorageMetrics>(5, _omitFieldNames ? '' : 'storage',
        subBuilder: StorageMetrics.$_createMessage)
    ..aOM<NetworkMetrics>(6, _omitFieldNames ? '' : 'network',
        subBuilder: NetworkMetrics.$_createMessage)
    ..aOM<DiskIoMetrics>(7, _omitFieldNames ? '' : 'diskIo',
        subBuilder: DiskIoMetrics.$_createMessage)
    ..aI(8, _omitFieldNames ? '' : 'activeContainersCount',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SystemMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SystemMetrics copyWith(void Function(SystemMetrics) updates) =>
      super.copyWith((message) => updates(message as SystemMetrics))
          as SystemMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated('Use SystemMetrics() / SystemMetrics.new instead')
  static SystemMetrics create() => SystemMetrics._();
  static $pb.GeneratedMessage $_createMessage() => SystemMetrics._();
  @$core.override
  SystemMetrics createEmptyInstance() => SystemMetrics._();
  @$core.pragma('dart2js:noInline')
  static SystemMetrics getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SystemMetrics>(
          SystemMetrics.$_createMessage);
  static SystemMetrics? _defaultInstance;

  @$pb.TagNumber(1)
  $1.Timestamp get collectedAt => $_getN(0);
  @$pb.TagNumber(1)
  set collectedAt($1.Timestamp value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasCollectedAt() => $_has(0);
  @$pb.TagNumber(1)
  void clearCollectedAt() => $_clearField(1);
  @$pb.TagNumber(1)
  $1.Timestamp ensureCollectedAt() => $_ensure(0);

  @$pb.TagNumber(2)
  $fixnum.Int64 get uptimeSeconds => $_getI64(1);
  @$pb.TagNumber(2)
  set uptimeSeconds($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUptimeSeconds() => $_has(1);
  @$pb.TagNumber(2)
  void clearUptimeSeconds() => $_clearField(2);

  @$pb.TagNumber(3)
  CpuMetrics get cpu => $_getN(2);
  @$pb.TagNumber(3)
  set cpu(CpuMetrics value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasCpu() => $_has(2);
  @$pb.TagNumber(3)
  void clearCpu() => $_clearField(3);
  @$pb.TagNumber(3)
  CpuMetrics ensureCpu() => $_ensure(2);

  @$pb.TagNumber(4)
  MemoryMetrics get memory => $_getN(3);
  @$pb.TagNumber(4)
  set memory(MemoryMetrics value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasMemory() => $_has(3);
  @$pb.TagNumber(4)
  void clearMemory() => $_clearField(4);
  @$pb.TagNumber(4)
  MemoryMetrics ensureMemory() => $_ensure(3);

  @$pb.TagNumber(5)
  StorageMetrics get storage => $_getN(4);
  @$pb.TagNumber(5)
  set storage(StorageMetrics value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasStorage() => $_has(4);
  @$pb.TagNumber(5)
  void clearStorage() => $_clearField(5);
  @$pb.TagNumber(5)
  StorageMetrics ensureStorage() => $_ensure(4);

  @$pb.TagNumber(6)
  NetworkMetrics get network => $_getN(5);
  @$pb.TagNumber(6)
  set network(NetworkMetrics value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasNetwork() => $_has(5);
  @$pb.TagNumber(6)
  void clearNetwork() => $_clearField(6);
  @$pb.TagNumber(6)
  NetworkMetrics ensureNetwork() => $_ensure(5);

  @$pb.TagNumber(7)
  DiskIoMetrics get diskIo => $_getN(6);
  @$pb.TagNumber(7)
  set diskIo(DiskIoMetrics value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasDiskIo() => $_has(6);
  @$pb.TagNumber(7)
  void clearDiskIo() => $_clearField(7);
  @$pb.TagNumber(7)
  DiskIoMetrics ensureDiskIo() => $_ensure(6);

  @$pb.TagNumber(8)
  $core.int get activeContainersCount => $_getIZ(7);
  @$pb.TagNumber(8)
  set activeContainersCount($core.int value) => $_setUnsignedInt32(7, value);
  @$pb.TagNumber(8)
  $core.bool hasActiveContainersCount() => $_has(7);
  @$pb.TagNumber(8)
  void clearActiveContainersCount() => $_clearField(8);
}

class GetContainerMetricsRequest extends $pb.GeneratedMessage {
  factory GetContainerMetricsRequest({
    $core.String? containerId,
  }) {
    final result = GetContainerMetricsRequest._();
    if (containerId != null) result.containerId = containerId;
    return result;
  }

  GetContainerMetricsRequest._();

  factory GetContainerMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerMetricsRequest()..mergeFromBuffer(data, registry);
  factory GetContainerMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContainerMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetContainerMetricsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerMetricsRequest copyWith(
          void Function(GetContainerMetricsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetContainerMetricsRequest))
          as GetContainerMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetContainerMetricsRequest() / GetContainerMetricsRequest.new instead')
  static GetContainerMetricsRequest create() => GetContainerMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetContainerMetricsRequest._();
  @$core.override
  GetContainerMetricsRequest createEmptyInstance() =>
      GetContainerMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetContainerMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContainerMetricsRequest>(
          GetContainerMetricsRequest.$_createMessage);
  static GetContainerMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);
}

class GetContainerMetricsResponse extends $pb.GeneratedMessage {
  factory GetContainerMetricsResponse({
    ContainerMetrics? metrics,
  }) {
    final result = GetContainerMetricsResponse._();
    if (metrics != null) result.metrics = metrics;
    return result;
  }

  GetContainerMetricsResponse._();

  factory GetContainerMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerMetricsResponse()..mergeFromBuffer(data, registry);
  factory GetContainerMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetContainerMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetContainerMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetContainerMetricsResponse.$_createMessage)
    ..aOM<ContainerMetrics>(1, _omitFieldNames ? '' : 'metrics',
        subBuilder: ContainerMetrics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetContainerMetricsResponse copyWith(
          void Function(GetContainerMetricsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetContainerMetricsResponse))
          as GetContainerMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetContainerMetricsResponse() / GetContainerMetricsResponse.new instead')
  static GetContainerMetricsResponse create() =>
      GetContainerMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      GetContainerMetricsResponse._();
  @$core.override
  GetContainerMetricsResponse createEmptyInstance() =>
      GetContainerMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetContainerMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetContainerMetricsResponse>(
          GetContainerMetricsResponse.$_createMessage);
  static GetContainerMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ContainerMetrics get metrics => $_getN(0);
  @$pb.TagNumber(1)
  set metrics(ContainerMetrics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMetrics() => $_has(0);
  @$pb.TagNumber(1)
  void clearMetrics() => $_clearField(1);
  @$pb.TagNumber(1)
  ContainerMetrics ensureMetrics() => $_ensure(0);
}

class StreamContainerMetricsRequest extends $pb.GeneratedMessage {
  factory StreamContainerMetricsRequest({
    $core.String? containerId,
    $core.int? intervalSeconds,
  }) {
    final result = StreamContainerMetricsRequest._();
    if (containerId != null) result.containerId = containerId;
    if (intervalSeconds != null) result.intervalSeconds = intervalSeconds;
    return result;
  }

  StreamContainerMetricsRequest._();

  factory StreamContainerMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerMetricsRequest()..mergeFromBuffer(data, registry);
  factory StreamContainerMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamContainerMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamContainerMetricsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..aI(2, _omitFieldNames ? '' : 'intervalSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerMetricsRequest copyWith(
          void Function(StreamContainerMetricsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StreamContainerMetricsRequest))
          as StreamContainerMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamContainerMetricsRequest() / StreamContainerMetricsRequest.new instead')
  static StreamContainerMetricsRequest create() =>
      StreamContainerMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamContainerMetricsRequest._();
  @$core.override
  StreamContainerMetricsRequest createEmptyInstance() =>
      StreamContainerMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static StreamContainerMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamContainerMetricsRequest>(
          StreamContainerMetricsRequest.$_createMessage);
  static StreamContainerMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  /// Intervalle d'échantillonnage en secondes.
  /// Si non spécifié ou égal à 0, le serveur utilise son intervalle par défaut (2s).
  @$pb.TagNumber(2)
  $core.int get intervalSeconds => $_getIZ(1);
  @$pb.TagNumber(2)
  set intervalSeconds($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIntervalSeconds() => $_has(1);
  @$pb.TagNumber(2)
  void clearIntervalSeconds() => $_clearField(2);
}

class StreamContainerMetricsResponse extends $pb.GeneratedMessage {
  factory StreamContainerMetricsResponse({
    ContainerMetrics? metrics,
  }) {
    final result = StreamContainerMetricsResponse._();
    if (metrics != null) result.metrics = metrics;
    return result;
  }

  StreamContainerMetricsResponse._();

  factory StreamContainerMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerMetricsResponse()..mergeFromBuffer(data, registry);
  factory StreamContainerMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamContainerMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamContainerMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamContainerMetricsResponse.$_createMessage)
    ..aOM<ContainerMetrics>(1, _omitFieldNames ? '' : 'metrics',
        subBuilder: ContainerMetrics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamContainerMetricsResponse copyWith(
          void Function(StreamContainerMetricsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StreamContainerMetricsResponse))
          as StreamContainerMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamContainerMetricsResponse() / StreamContainerMetricsResponse.new instead')
  static StreamContainerMetricsResponse create() =>
      StreamContainerMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamContainerMetricsResponse._();
  @$core.override
  StreamContainerMetricsResponse createEmptyInstance() =>
      StreamContainerMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static StreamContainerMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamContainerMetricsResponse>(
          StreamContainerMetricsResponse.$_createMessage);
  static StreamContainerMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ContainerMetrics get metrics => $_getN(0);
  @$pb.TagNumber(1)
  set metrics(ContainerMetrics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMetrics() => $_has(0);
  @$pb.TagNumber(1)
  void clearMetrics() => $_clearField(1);
  @$pb.TagNumber(1)
  ContainerMetrics ensureMetrics() => $_ensure(0);
}

class QueryContainerMetricsRequest extends $pb.GeneratedMessage {
  factory QueryContainerMetricsRequest({
    $core.String? containerId,
    TimeRange? range,
    $core.int? stepSeconds,
    $core.int? maxPoints,
    $core.Iterable<MetricName>? metricNames,
  }) {
    final result = QueryContainerMetricsRequest._();
    if (containerId != null) result.containerId = containerId;
    if (range != null) result.range = range;
    if (stepSeconds != null) result.stepSeconds = stepSeconds;
    if (maxPoints != null) result.maxPoints = maxPoints;
    if (metricNames != null) result.metricNames.addAll(metricNames);
    return result;
  }

  QueryContainerMetricsRequest._();

  factory QueryContainerMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QueryContainerMetricsRequest()..mergeFromBuffer(data, registry);
  factory QueryContainerMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QueryContainerMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryContainerMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: QueryContainerMetricsRequest.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..aOM<TimeRange>(2, _omitFieldNames ? '' : 'range',
        subBuilder: TimeRange.$_createMessage)
    ..aI(3, _omitFieldNames ? '' : 'stepSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(4, _omitFieldNames ? '' : 'maxPoints', fieldType: $pb.PbFieldType.OU3)
    ..pc<MetricName>(
        5, _omitFieldNames ? '' : 'metricNames', $pb.PbFieldType.KE,
        valueOf: MetricName.valueOf,
        enumValues: MetricName.values,
        defaultEnumValue: MetricName.METRIC_NAME_UNSPECIFIED)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryContainerMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryContainerMetricsRequest copyWith(
          void Function(QueryContainerMetricsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as QueryContainerMetricsRequest))
          as QueryContainerMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use QueryContainerMetricsRequest() / QueryContainerMetricsRequest.new instead')
  static QueryContainerMetricsRequest create() =>
      QueryContainerMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      QueryContainerMetricsRequest._();
  @$core.override
  QueryContainerMetricsRequest createEmptyInstance() =>
      QueryContainerMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static QueryContainerMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryContainerMetricsRequest>(
          QueryContainerMetricsRequest.$_createMessage);
  static QueryContainerMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  @$pb.TagNumber(2)
  TimeRange get range => $_getN(1);
  @$pb.TagNumber(2)
  set range(TimeRange value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasRange() => $_has(1);
  @$pb.TagNumber(2)
  void clearRange() => $_clearField(2);
  @$pb.TagNumber(2)
  TimeRange ensureRange() => $_ensure(1);

  /// Résolution d'agrégation souhaitée en secondes (ex: 60s pour un point par minute).
  /// Si non renseigné ou égal à 0, le serveur applique un pas d'agrégation automatique adapté à la plage demandée.
  @$pb.TagNumber(3)
  $core.int get stepSeconds => $_getIZ(2);
  @$pb.TagNumber(3)
  set stepSeconds($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStepSeconds() => $_has(2);
  @$pb.TagNumber(3)
  void clearStepSeconds() => $_clearField(3);

  /// Nombre maximum de points retournés par série pour économiser la bande passante.
  /// Si 0, le serveur applique son plafond maximal (ex: 1000 points).
  @$pb.TagNumber(4)
  $core.int get maxPoints => $_getIZ(3);
  @$pb.TagNumber(4)
  set maxPoints($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMaxPoints() => $_has(3);
  @$pb.TagNumber(4)
  void clearMaxPoints() => $_clearField(4);

  /// Liste optionnelle des métriques spécifiques à extraire. Si vide, toutes les métriques sont retournées.
  @$pb.TagNumber(5)
  $pb.PbList<MetricName> get metricNames => $_getList(4);
}

class QueryContainerMetricsResponse extends $pb.GeneratedMessage {
  factory QueryContainerMetricsResponse({
    $core.String? containerId,
    $core.Iterable<MetricSeries>? series,
  }) {
    final result = QueryContainerMetricsResponse._();
    if (containerId != null) result.containerId = containerId;
    if (series != null) result.series.addAll(series);
    return result;
  }

  QueryContainerMetricsResponse._();

  factory QueryContainerMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QueryContainerMetricsResponse()..mergeFromBuffer(data, registry);
  factory QueryContainerMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QueryContainerMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryContainerMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: QueryContainerMetricsResponse.$_createMessage)
    ..aOS(1, _omitFieldNames ? '' : 'containerId')
    ..pPM<MetricSeries>(2, _omitFieldNames ? '' : 'series',
        subBuilder: MetricSeries.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryContainerMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryContainerMetricsResponse copyWith(
          void Function(QueryContainerMetricsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as QueryContainerMetricsResponse))
          as QueryContainerMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use QueryContainerMetricsResponse() / QueryContainerMetricsResponse.new instead')
  static QueryContainerMetricsResponse create() =>
      QueryContainerMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      QueryContainerMetricsResponse._();
  @$core.override
  QueryContainerMetricsResponse createEmptyInstance() =>
      QueryContainerMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static QueryContainerMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryContainerMetricsResponse>(
          QueryContainerMetricsResponse.$_createMessage);
  static QueryContainerMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get containerId => $_getSZ(0);
  @$pb.TagNumber(1)
  set containerId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContainerId() => $_has(0);
  @$pb.TagNumber(1)
  void clearContainerId() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<MetricSeries> get series => $_getList(1);
}

class GetSystemMetricsRequest extends $pb.GeneratedMessage {
  factory GetSystemMetricsRequest() => GetSystemMetricsRequest._();

  GetSystemMetricsRequest._();

  factory GetSystemMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSystemMetricsRequest()..mergeFromBuffer(data, registry);
  factory GetSystemMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSystemMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSystemMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetSystemMetricsRequest.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSystemMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSystemMetricsRequest copyWith(
          void Function(GetSystemMetricsRequest) updates) =>
      super.copyWith((message) => updates(message as GetSystemMetricsRequest))
          as GetSystemMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSystemMetricsRequest() / GetSystemMetricsRequest.new instead')
  static GetSystemMetricsRequest create() => GetSystemMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() => GetSystemMetricsRequest._();
  @$core.override
  GetSystemMetricsRequest createEmptyInstance() => GetSystemMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static GetSystemMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSystemMetricsRequest>(
          GetSystemMetricsRequest.$_createMessage);
  static GetSystemMetricsRequest? _defaultInstance;
}

class GetSystemMetricsResponse extends $pb.GeneratedMessage {
  factory GetSystemMetricsResponse({
    SystemMetrics? metrics,
  }) {
    final result = GetSystemMetricsResponse._();
    if (metrics != null) result.metrics = metrics;
    return result;
  }

  GetSystemMetricsResponse._();

  factory GetSystemMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSystemMetricsResponse()..mergeFromBuffer(data, registry);
  factory GetSystemMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      GetSystemMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSystemMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: GetSystemMetricsResponse.$_createMessage)
    ..aOM<SystemMetrics>(1, _omitFieldNames ? '' : 'metrics',
        subBuilder: SystemMetrics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSystemMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSystemMetricsResponse copyWith(
          void Function(GetSystemMetricsResponse) updates) =>
      super.copyWith((message) => updates(message as GetSystemMetricsResponse))
          as GetSystemMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use GetSystemMetricsResponse() / GetSystemMetricsResponse.new instead')
  static GetSystemMetricsResponse create() => GetSystemMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() => GetSystemMetricsResponse._();
  @$core.override
  GetSystemMetricsResponse createEmptyInstance() =>
      GetSystemMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static GetSystemMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSystemMetricsResponse>(
          GetSystemMetricsResponse.$_createMessage);
  static GetSystemMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  SystemMetrics get metrics => $_getN(0);
  @$pb.TagNumber(1)
  set metrics(SystemMetrics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMetrics() => $_has(0);
  @$pb.TagNumber(1)
  void clearMetrics() => $_clearField(1);
  @$pb.TagNumber(1)
  SystemMetrics ensureMetrics() => $_ensure(0);
}

class StreamSystemMetricsRequest extends $pb.GeneratedMessage {
  factory StreamSystemMetricsRequest({
    $core.int? intervalSeconds,
  }) {
    final result = StreamSystemMetricsRequest._();
    if (intervalSeconds != null) result.intervalSeconds = intervalSeconds;
    return result;
  }

  StreamSystemMetricsRequest._();

  factory StreamSystemMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamSystemMetricsRequest()..mergeFromBuffer(data, registry);
  factory StreamSystemMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamSystemMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamSystemMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamSystemMetricsRequest.$_createMessage)
    ..aI(1, _omitFieldNames ? '' : 'intervalSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSystemMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSystemMetricsRequest copyWith(
          void Function(StreamSystemMetricsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StreamSystemMetricsRequest))
          as StreamSystemMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamSystemMetricsRequest() / StreamSystemMetricsRequest.new instead')
  static StreamSystemMetricsRequest create() => StreamSystemMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamSystemMetricsRequest._();
  @$core.override
  StreamSystemMetricsRequest createEmptyInstance() =>
      StreamSystemMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static StreamSystemMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamSystemMetricsRequest>(
          StreamSystemMetricsRequest.$_createMessage);
  static StreamSystemMetricsRequest? _defaultInstance;

  /// Intervalle d'échantillonnage en secondes.
  /// Si non spécifié ou égal à 0, le serveur utilise son intervalle par défaut (5s).
  @$pb.TagNumber(1)
  $core.int get intervalSeconds => $_getIZ(0);
  @$pb.TagNumber(1)
  set intervalSeconds($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIntervalSeconds() => $_has(0);
  @$pb.TagNumber(1)
  void clearIntervalSeconds() => $_clearField(1);
}

class StreamSystemMetricsResponse extends $pb.GeneratedMessage {
  factory StreamSystemMetricsResponse({
    SystemMetrics? metrics,
  }) {
    final result = StreamSystemMetricsResponse._();
    if (metrics != null) result.metrics = metrics;
    return result;
  }

  StreamSystemMetricsResponse._();

  factory StreamSystemMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamSystemMetricsResponse()..mergeFromBuffer(data, registry);
  factory StreamSystemMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      StreamSystemMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamSystemMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: StreamSystemMetricsResponse.$_createMessage)
    ..aOM<SystemMetrics>(1, _omitFieldNames ? '' : 'metrics',
        subBuilder: SystemMetrics.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSystemMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSystemMetricsResponse copyWith(
          void Function(StreamSystemMetricsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StreamSystemMetricsResponse))
          as StreamSystemMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use StreamSystemMetricsResponse() / StreamSystemMetricsResponse.new instead')
  static StreamSystemMetricsResponse create() =>
      StreamSystemMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      StreamSystemMetricsResponse._();
  @$core.override
  StreamSystemMetricsResponse createEmptyInstance() =>
      StreamSystemMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static StreamSystemMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamSystemMetricsResponse>(
          StreamSystemMetricsResponse.$_createMessage);
  static StreamSystemMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  SystemMetrics get metrics => $_getN(0);
  @$pb.TagNumber(1)
  set metrics(SystemMetrics value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMetrics() => $_has(0);
  @$pb.TagNumber(1)
  void clearMetrics() => $_clearField(1);
  @$pb.TagNumber(1)
  SystemMetrics ensureMetrics() => $_ensure(0);
}

class QuerySystemMetricsRequest extends $pb.GeneratedMessage {
  factory QuerySystemMetricsRequest({
    TimeRange? range,
    $core.int? stepSeconds,
    $core.int? maxPoints,
    $core.Iterable<MetricName>? metricNames,
  }) {
    final result = QuerySystemMetricsRequest._();
    if (range != null) result.range = range;
    if (stepSeconds != null) result.stepSeconds = stepSeconds;
    if (maxPoints != null) result.maxPoints = maxPoints;
    if (metricNames != null) result.metricNames.addAll(metricNames);
    return result;
  }

  QuerySystemMetricsRequest._();

  factory QuerySystemMetricsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QuerySystemMetricsRequest()..mergeFromBuffer(data, registry);
  factory QuerySystemMetricsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QuerySystemMetricsRequest()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuerySystemMetricsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: QuerySystemMetricsRequest.$_createMessage)
    ..aOM<TimeRange>(1, _omitFieldNames ? '' : 'range',
        subBuilder: TimeRange.$_createMessage)
    ..aI(2, _omitFieldNames ? '' : 'stepSeconds',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(3, _omitFieldNames ? '' : 'maxPoints', fieldType: $pb.PbFieldType.OU3)
    ..pc<MetricName>(
        4, _omitFieldNames ? '' : 'metricNames', $pb.PbFieldType.KE,
        valueOf: MetricName.valueOf,
        enumValues: MetricName.values,
        defaultEnumValue: MetricName.METRIC_NAME_UNSPECIFIED)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySystemMetricsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySystemMetricsRequest copyWith(
          void Function(QuerySystemMetricsRequest) updates) =>
      super.copyWith((message) => updates(message as QuerySystemMetricsRequest))
          as QuerySystemMetricsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use QuerySystemMetricsRequest() / QuerySystemMetricsRequest.new instead')
  static QuerySystemMetricsRequest create() => QuerySystemMetricsRequest._();
  static $pb.GeneratedMessage $_createMessage() =>
      QuerySystemMetricsRequest._();
  @$core.override
  QuerySystemMetricsRequest createEmptyInstance() =>
      QuerySystemMetricsRequest._();
  @$core.pragma('dart2js:noInline')
  static QuerySystemMetricsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuerySystemMetricsRequest>(
          QuerySystemMetricsRequest.$_createMessage);
  static QuerySystemMetricsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  TimeRange get range => $_getN(0);
  @$pb.TagNumber(1)
  set range(TimeRange value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasRange() => $_has(0);
  @$pb.TagNumber(1)
  void clearRange() => $_clearField(1);
  @$pb.TagNumber(1)
  TimeRange ensureRange() => $_ensure(0);

  /// Résolution d'agrégation souhaitée en secondes.
  /// Si 0, le serveur applique son pas d'agrégation automatique.
  @$pb.TagNumber(2)
  $core.int get stepSeconds => $_getIZ(1);
  @$pb.TagNumber(2)
  set stepSeconds($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasStepSeconds() => $_has(1);
  @$pb.TagNumber(2)
  void clearStepSeconds() => $_clearField(2);

  /// Nombre maximum de points retournés par série.
  /// Si 0, le serveur applique son plafond par défaut (ex: 1000 points).
  @$pb.TagNumber(3)
  $core.int get maxPoints => $_getIZ(2);
  @$pb.TagNumber(3)
  set maxPoints($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMaxPoints() => $_has(2);
  @$pb.TagNumber(3)
  void clearMaxPoints() => $_clearField(3);

  /// Liste optionnelle des métriques spécifiques à extraire. Si vide, toutes les métriques sont retournées.
  @$pb.TagNumber(4)
  $pb.PbList<MetricName> get metricNames => $_getList(3);
}

class QuerySystemMetricsResponse extends $pb.GeneratedMessage {
  factory QuerySystemMetricsResponse({
    $core.Iterable<MetricSeries>? series,
  }) {
    final result = QuerySystemMetricsResponse._();
    if (series != null) result.series.addAll(series);
    return result;
  }

  QuerySystemMetricsResponse._();

  factory QuerySystemMetricsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QuerySystemMetricsResponse()..mergeFromBuffer(data, registry);
  factory QuerySystemMetricsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      QuerySystemMetricsResponse()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuerySystemMetricsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'forgenet.v1'),
      createEmptyInstance: QuerySystemMetricsResponse.$_createMessage)
    ..pPM<MetricSeries>(1, _omitFieldNames ? '' : 'series',
        subBuilder: MetricSeries.$_createMessage)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySystemMetricsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySystemMetricsResponse copyWith(
          void Function(QuerySystemMetricsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as QuerySystemMetricsResponse))
          as QuerySystemMetricsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  @$core.Deprecated(
      'Use QuerySystemMetricsResponse() / QuerySystemMetricsResponse.new instead')
  static QuerySystemMetricsResponse create() => QuerySystemMetricsResponse._();
  static $pb.GeneratedMessage $_createMessage() =>
      QuerySystemMetricsResponse._();
  @$core.override
  QuerySystemMetricsResponse createEmptyInstance() =>
      QuerySystemMetricsResponse._();
  @$core.pragma('dart2js:noInline')
  static QuerySystemMetricsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuerySystemMetricsResponse>(
          QuerySystemMetricsResponse.$_createMessage);
  static QuerySystemMetricsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<MetricSeries> get series => $_getList(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
