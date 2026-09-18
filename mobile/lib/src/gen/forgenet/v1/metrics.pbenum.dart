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

import 'package:protobuf/protobuf.dart' as $pb;

/// MetricUnit définit la liste fermée et exhaustive des unités de mesure.
class MetricUnit extends $pb.ProtobufEnum {
  static const MetricUnit METRIC_UNIT_UNSPECIFIED =
      MetricUnit._(0, _omitEnumNames ? '' : 'METRIC_UNIT_UNSPECIFIED');
  static const MetricUnit METRIC_UNIT_PERCENT =
      MetricUnit._(1, _omitEnumNames ? '' : 'METRIC_UNIT_PERCENT');
  static const MetricUnit METRIC_UNIT_BYTES =
      MetricUnit._(2, _omitEnumNames ? '' : 'METRIC_UNIT_BYTES');
  static const MetricUnit METRIC_UNIT_NANOSECONDS =
      MetricUnit._(3, _omitEnumNames ? '' : 'METRIC_UNIT_NANOSECONDS');
  static const MetricUnit METRIC_UNIT_PACKETS =
      MetricUnit._(4, _omitEnumNames ? '' : 'METRIC_UNIT_PACKETS');
  static const MetricUnit METRIC_UNIT_OPERATIONS =
      MetricUnit._(5, _omitEnumNames ? '' : 'METRIC_UNIT_OPERATIONS');
  static const MetricUnit METRIC_UNIT_SECONDS =
      MetricUnit._(6, _omitEnumNames ? '' : 'METRIC_UNIT_SECONDS');
  static const MetricUnit METRIC_UNIT_COUNT =
      MetricUnit._(7, _omitEnumNames ? '' : 'METRIC_UNIT_COUNT');

  static const $core.List<MetricUnit> values = <MetricUnit>[
    METRIC_UNIT_UNSPECIFIED,
    METRIC_UNIT_PERCENT,
    METRIC_UNIT_BYTES,
    METRIC_UNIT_NANOSECONDS,
    METRIC_UNIT_PACKETS,
    METRIC_UNIT_OPERATIONS,
    METRIC_UNIT_SECONDS,
    METRIC_UNIT_COUNT,
  ];

  static final $core.List<MetricUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static MetricUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricUnit._(super.value, super.name);
}

/// MetricName définit la liste fermée des métriques observables dans le système.
class MetricName extends $pb.ProtobufEnum {
  static const MetricName METRIC_NAME_UNSPECIFIED =
      MetricName._(0, _omitEnumNames ? '' : 'METRIC_NAME_UNSPECIFIED');
  static const MetricName METRIC_NAME_CPU_USAGE_PERCENT =
      MetricName._(1, _omitEnumNames ? '' : 'METRIC_NAME_CPU_USAGE_PERCENT');
  static const MetricName METRIC_NAME_CPU_SYSTEM_NANOSECONDS = MetricName._(
      2, _omitEnumNames ? '' : 'METRIC_NAME_CPU_SYSTEM_NANOSECONDS');
  static const MetricName METRIC_NAME_MEMORY_USAGE_BYTES =
      MetricName._(3, _omitEnumNames ? '' : 'METRIC_NAME_MEMORY_USAGE_BYTES');
  static const MetricName METRIC_NAME_MEMORY_LIMIT_BYTES =
      MetricName._(4, _omitEnumNames ? '' : 'METRIC_NAME_MEMORY_LIMIT_BYTES');
  static const MetricName METRIC_NAME_MEMORY_CACHE_BYTES =
      MetricName._(5, _omitEnumNames ? '' : 'METRIC_NAME_MEMORY_CACHE_BYTES');
  static const MetricName METRIC_NAME_MEMORY_USAGE_PERCENT =
      MetricName._(6, _omitEnumNames ? '' : 'METRIC_NAME_MEMORY_USAGE_PERCENT');
  static const MetricName METRIC_NAME_NETWORK_RX_BYTES_TOTAL = MetricName._(
      7, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_RX_BYTES_TOTAL');
  static const MetricName METRIC_NAME_NETWORK_TX_BYTES_TOTAL = MetricName._(
      8, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_TX_BYTES_TOTAL');
  static const MetricName METRIC_NAME_NETWORK_RX_PACKETS_TOTAL = MetricName._(
      9, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_RX_PACKETS_TOTAL');
  static const MetricName METRIC_NAME_NETWORK_TX_PACKETS_TOTAL = MetricName._(
      10, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_TX_PACKETS_TOTAL');
  static const MetricName METRIC_NAME_NETWORK_RX_ERRORS_TOTAL = MetricName._(
      11, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_RX_ERRORS_TOTAL');
  static const MetricName METRIC_NAME_NETWORK_TX_ERRORS_TOTAL = MetricName._(
      12, _omitEnumNames ? '' : 'METRIC_NAME_NETWORK_TX_ERRORS_TOTAL');
  static const MetricName METRIC_NAME_DISK_READ_BYTES_TOTAL = MetricName._(
      13, _omitEnumNames ? '' : 'METRIC_NAME_DISK_READ_BYTES_TOTAL');
  static const MetricName METRIC_NAME_DISK_WRITE_BYTES_TOTAL = MetricName._(
      14, _omitEnumNames ? '' : 'METRIC_NAME_DISK_WRITE_BYTES_TOTAL');
  static const MetricName METRIC_NAME_DISK_READ_OPS_TOTAL =
      MetricName._(15, _omitEnumNames ? '' : 'METRIC_NAME_DISK_READ_OPS_TOTAL');
  static const MetricName METRIC_NAME_DISK_WRITE_OPS_TOTAL = MetricName._(
      16, _omitEnumNames ? '' : 'METRIC_NAME_DISK_WRITE_OPS_TOTAL');
  static const MetricName METRIC_NAME_STORAGE_TOTAL_BYTES =
      MetricName._(17, _omitEnumNames ? '' : 'METRIC_NAME_STORAGE_TOTAL_BYTES');
  static const MetricName METRIC_NAME_STORAGE_USED_BYTES =
      MetricName._(18, _omitEnumNames ? '' : 'METRIC_NAME_STORAGE_USED_BYTES');
  static const MetricName METRIC_NAME_STORAGE_FREE_BYTES =
      MetricName._(19, _omitEnumNames ? '' : 'METRIC_NAME_STORAGE_FREE_BYTES');
  static const MetricName METRIC_NAME_STORAGE_USAGE_PERCENT = MetricName._(
      20, _omitEnumNames ? '' : 'METRIC_NAME_STORAGE_USAGE_PERCENT');
  static const MetricName METRIC_NAME_SYSTEM_UPTIME_SECONDS = MetricName._(
      21, _omitEnumNames ? '' : 'METRIC_NAME_SYSTEM_UPTIME_SECONDS');
  static const MetricName METRIC_NAME_SYSTEM_ACTIVE_CONTAINERS_COUNT =
      MetricName._(22,
          _omitEnumNames ? '' : 'METRIC_NAME_SYSTEM_ACTIVE_CONTAINERS_COUNT');

  static const $core.List<MetricName> values = <MetricName>[
    METRIC_NAME_UNSPECIFIED,
    METRIC_NAME_CPU_USAGE_PERCENT,
    METRIC_NAME_CPU_SYSTEM_NANOSECONDS,
    METRIC_NAME_MEMORY_USAGE_BYTES,
    METRIC_NAME_MEMORY_LIMIT_BYTES,
    METRIC_NAME_MEMORY_CACHE_BYTES,
    METRIC_NAME_MEMORY_USAGE_PERCENT,
    METRIC_NAME_NETWORK_RX_BYTES_TOTAL,
    METRIC_NAME_NETWORK_TX_BYTES_TOTAL,
    METRIC_NAME_NETWORK_RX_PACKETS_TOTAL,
    METRIC_NAME_NETWORK_TX_PACKETS_TOTAL,
    METRIC_NAME_NETWORK_RX_ERRORS_TOTAL,
    METRIC_NAME_NETWORK_TX_ERRORS_TOTAL,
    METRIC_NAME_DISK_READ_BYTES_TOTAL,
    METRIC_NAME_DISK_WRITE_BYTES_TOTAL,
    METRIC_NAME_DISK_READ_OPS_TOTAL,
    METRIC_NAME_DISK_WRITE_OPS_TOTAL,
    METRIC_NAME_STORAGE_TOTAL_BYTES,
    METRIC_NAME_STORAGE_USED_BYTES,
    METRIC_NAME_STORAGE_FREE_BYTES,
    METRIC_NAME_STORAGE_USAGE_PERCENT,
    METRIC_NAME_SYSTEM_UPTIME_SECONDS,
    METRIC_NAME_SYSTEM_ACTIVE_CONTAINERS_COUNT,
  ];

  static final $core.List<MetricName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 22);
  static MetricName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricName._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
