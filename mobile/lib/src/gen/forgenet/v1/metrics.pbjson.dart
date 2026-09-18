// This is a generated file - do not edit.
//
// Generated from forgenet/v1/metrics.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use metricUnitDescriptor instead')
const MetricUnit$json = {
  '1': 'MetricUnit',
  '2': [
    {'1': 'METRIC_UNIT_UNSPECIFIED', '2': 0},
    {'1': 'METRIC_UNIT_PERCENT', '2': 1},
    {'1': 'METRIC_UNIT_BYTES', '2': 2},
    {'1': 'METRIC_UNIT_NANOSECONDS', '2': 3},
    {'1': 'METRIC_UNIT_PACKETS', '2': 4},
    {'1': 'METRIC_UNIT_OPERATIONS', '2': 5},
    {'1': 'METRIC_UNIT_SECONDS', '2': 6},
    {'1': 'METRIC_UNIT_COUNT', '2': 7},
  ],
};

/// Descriptor for `MetricUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricUnitDescriptor = $convert.base64Decode(
    'CgpNZXRyaWNVbml0EhsKF01FVFJJQ19VTklUX1VOU1BFQ0lGSUVEEAASFwoTTUVUUklDX1VOSV'
    'RfUEVSQ0VOVBABEhUKEU1FVFJJQ19VTklUX0JZVEVTEAISGwoXTUVUUklDX1VOSVRfTkFOT1NF'
    'Q09ORFMQAxIXChNNRVRSSUNfVU5JVF9QQUNLRVRTEAQSGgoWTUVUUklDX1VOSVRfT1BFUkFUSU'
    '9OUxAFEhcKE01FVFJJQ19VTklUX1NFQ09ORFMQBhIVChFNRVRSSUNfVU5JVF9DT1VOVBAH');

@$core.Deprecated('Use metricNameDescriptor instead')
const MetricName$json = {
  '1': 'MetricName',
  '2': [
    {'1': 'METRIC_NAME_UNSPECIFIED', '2': 0},
    {'1': 'METRIC_NAME_CPU_USAGE_PERCENT', '2': 1},
    {'1': 'METRIC_NAME_CPU_SYSTEM_NANOSECONDS', '2': 2},
    {'1': 'METRIC_NAME_MEMORY_USAGE_BYTES', '2': 3},
    {'1': 'METRIC_NAME_MEMORY_LIMIT_BYTES', '2': 4},
    {'1': 'METRIC_NAME_MEMORY_CACHE_BYTES', '2': 5},
    {'1': 'METRIC_NAME_MEMORY_USAGE_PERCENT', '2': 6},
    {'1': 'METRIC_NAME_NETWORK_RX_BYTES_TOTAL', '2': 7},
    {'1': 'METRIC_NAME_NETWORK_TX_BYTES_TOTAL', '2': 8},
    {'1': 'METRIC_NAME_NETWORK_RX_PACKETS_TOTAL', '2': 9},
    {'1': 'METRIC_NAME_NETWORK_TX_PACKETS_TOTAL', '2': 10},
    {'1': 'METRIC_NAME_NETWORK_RX_ERRORS_TOTAL', '2': 11},
    {'1': 'METRIC_NAME_NETWORK_TX_ERRORS_TOTAL', '2': 12},
    {'1': 'METRIC_NAME_DISK_READ_BYTES_TOTAL', '2': 13},
    {'1': 'METRIC_NAME_DISK_WRITE_BYTES_TOTAL', '2': 14},
    {'1': 'METRIC_NAME_DISK_READ_OPS_TOTAL', '2': 15},
    {'1': 'METRIC_NAME_DISK_WRITE_OPS_TOTAL', '2': 16},
    {'1': 'METRIC_NAME_STORAGE_TOTAL_BYTES', '2': 17},
    {'1': 'METRIC_NAME_STORAGE_USED_BYTES', '2': 18},
    {'1': 'METRIC_NAME_STORAGE_FREE_BYTES', '2': 19},
    {'1': 'METRIC_NAME_STORAGE_USAGE_PERCENT', '2': 20},
    {'1': 'METRIC_NAME_SYSTEM_UPTIME_SECONDS', '2': 21},
    {'1': 'METRIC_NAME_SYSTEM_ACTIVE_CONTAINERS_COUNT', '2': 22},
  ],
};

/// Descriptor for `MetricName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricNameDescriptor = $convert.base64Decode(
    'CgpNZXRyaWNOYW1lEhsKF01FVFJJQ19OQU1FX1VOU1BFQ0lGSUVEEAASIQodTUVUUklDX05BTU'
    'VfQ1BVX1VTQUdFX1BFUkNFTlQQARImCiJNRVRSSUNfTkFNRV9DUFVfU1lTVEVNX05BTk9TRUNP'
    'TkRTEAISIgoeTUVUUklDX05BTUVfTUVNT1JZX1VTQUdFX0JZVEVTEAMSIgoeTUVUUklDX05BTU'
    'VfTUVNT1JZX0xJTUlUX0JZVEVTEAQSIgoeTUVUUklDX05BTUVfTUVNT1JZX0NBQ0hFX0JZVEVT'
    'EAUSJAogTUVUUklDX05BTUVfTUVNT1JZX1VTQUdFX1BFUkNFTlQQBhImCiJNRVRSSUNfTkFNRV'
    '9ORVRXT1JLX1JYX0JZVEVTX1RPVEFMEAcSJgoiTUVUUklDX05BTUVfTkVUV09SS19UWF9CWVRF'
    'U19UT1RBTBAIEigKJE1FVFJJQ19OQU1FX05FVFdPUktfUlhfUEFDS0VUU19UT1RBTBAJEigKJE'
    '1FVFJJQ19OQU1FX05FVFdPUktfVFhfUEFDS0VUU19UT1RBTBAKEicKI01FVFJJQ19OQU1FX05F'
    'VFdPUktfUlhfRVJST1JTX1RPVEFMEAsSJwojTUVUUklDX05BTUVfTkVUV09SS19UWF9FUlJPUl'
    'NfVE9UQUwQDBIlCiFNRVRSSUNfTkFNRV9ESVNLX1JFQURfQllURVNfVE9UQUwQDRImCiJNRVRS'
    'SUNfTkFNRV9ESVNLX1dSSVRFX0JZVEVTX1RPVEFMEA4SIwofTUVUUklDX05BTUVfRElTS19SRU'
    'FEX09QU19UT1RBTBAPEiQKIE1FVFJJQ19OQU1FX0RJU0tfV1JJVEVfT1BTX1RPVEFMEBASIwof'
    'TUVUUklDX05BTUVfU1RPUkFHRV9UT1RBTF9CWVRFUxAREiIKHk1FVFJJQ19OQU1FX1NUT1JBR0'
    'VfVVNFRF9CWVRFUxASEiIKHk1FVFJJQ19OQU1FX1NUT1JBR0VfRlJFRV9CWVRFUxATEiUKIU1F'
    'VFJJQ19OQU1FX1NUT1JBR0VfVVNBR0VfUEVSQ0VOVBAUEiUKIU1FVFJJQ19OQU1FX1NZU1RFTV'
    '9VUFRJTUVfU0VDT05EUxAVEi4KKk1FVFJJQ19OQU1FX1NZU1RFTV9BQ1RJVkVfQ09OVEFJTkVS'
    'U19DT1VOVBAW');

@$core.Deprecated('Use metricKeyDescriptor instead')
const MetricKey$json = {
  '1': 'MetricKey',
  '2': [
    {
      '1': 'known',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MetricName',
      '9': 0,
      '10': 'known'
    },
    {'1': 'custom', '3': 2, '4': 1, '5': 9, '9': 0, '10': 'custom'},
  ],
  '8': [
    {'1': 'key'},
  ],
};

/// Descriptor for `MetricKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricKeyDescriptor = $convert.base64Decode(
    'CglNZXRyaWNLZXkSLwoFa25vd24YASABKA4yFy5mb3JnZW5ldC52MS5NZXRyaWNOYW1lSABSBW'
    'tub3duEhgKBmN1c3RvbRgCIAEoCUgAUgZjdXN0b21CBQoDa2V5');

@$core.Deprecated('Use timeRangeDescriptor instead')
const TimeRange$json = {
  '1': 'TimeRange',
  '2': [
    {
      '1': 'start_time',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startTime'
    },
    {
      '1': 'end_time',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endTime'
    },
  ],
};

/// Descriptor for `TimeRange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeRangeDescriptor = $convert.base64Decode(
    'CglUaW1lUmFuZ2USOQoKc3RhcnRfdGltZRgBIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3'
    'RhbXBSCXN0YXJ0VGltZRI1CghlbmRfdGltZRgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1l'
    'c3RhbXBSB2VuZFRpbWU=');

@$core.Deprecated('Use metricSampleDescriptor instead')
const MetricSample$json = {
  '1': 'MetricSample',
  '2': [
    {
      '1': 'timestamp',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'timestamp'
    },
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
  ],
};

/// Descriptor for `MetricSample`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricSampleDescriptor = $convert.base64Decode(
    'CgxNZXRyaWNTYW1wbGUSOAoJdGltZXN0YW1wGAEgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbW'
    'VzdGFtcFIJdGltZXN0YW1wEhQKBXZhbHVlGAIgASgBUgV2YWx1ZQ==');

@$core.Deprecated('Use metricSeriesDescriptor instead')
const MetricSeries$json = {
  '1': 'MetricSeries',
  '2': [
    {
      '1': 'key',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.MetricKey',
      '10': 'key'
    },
    {
      '1': 'unit',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MetricUnit',
      '10': 'unit'
    },
    {
      '1': 'samples',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.MetricSample',
      '10': 'samples'
    },
  ],
};

/// Descriptor for `MetricSeries`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricSeriesDescriptor = $convert.base64Decode(
    'CgxNZXRyaWNTZXJpZXMSKAoDa2V5GAEgASgLMhYuZm9yZ2VuZXQudjEuTWV0cmljS2V5UgNrZX'
    'kSKwoEdW5pdBgCIAEoDjIXLmZvcmdlbmV0LnYxLk1ldHJpY1VuaXRSBHVuaXQSMwoHc2FtcGxl'
    'cxgDIAMoCzIZLmZvcmdlbmV0LnYxLk1ldHJpY1NhbXBsZVIHc2FtcGxlcw==');

@$core.Deprecated('Use cpuMetricsDescriptor instead')
const CpuMetrics$json = {
  '1': 'CpuMetrics',
  '2': [
    {'1': 'cpu_usage_percent', '3': 1, '4': 1, '5': 1, '10': 'cpuUsagePercent'},
    {'1': 'online_cpus', '3': 2, '4': 1, '5': 13, '10': 'onlineCpus'},
    {
      '1': 'system_cpu_nanoseconds',
      '3': 3,
      '4': 1,
      '5': 4,
      '10': 'systemCpuNanoseconds'
    },
  ],
};

/// Descriptor for `CpuMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cpuMetricsDescriptor = $convert.base64Decode(
    'CgpDcHVNZXRyaWNzEioKEWNwdV91c2FnZV9wZXJjZW50GAEgASgBUg9jcHVVc2FnZVBlcmNlbn'
    'QSHwoLb25saW5lX2NwdXMYAiABKA1SCm9ubGluZUNwdXMSNAoWc3lzdGVtX2NwdV9uYW5vc2Vj'
    'b25kcxgDIAEoBFIUc3lzdGVtQ3B1TmFub3NlY29uZHM=');

@$core.Deprecated('Use memoryMetricsDescriptor instead')
const MemoryMetrics$json = {
  '1': 'MemoryMetrics',
  '2': [
    {'1': 'usage_bytes', '3': 1, '4': 1, '5': 4, '10': 'usageBytes'},
    {'1': 'limit_bytes', '3': 2, '4': 1, '5': 4, '10': 'limitBytes'},
    {'1': 'cache_bytes', '3': 3, '4': 1, '5': 4, '10': 'cacheBytes'},
    {'1': 'usage_percent', '3': 4, '4': 1, '5': 1, '10': 'usagePercent'},
  ],
};

/// Descriptor for `MemoryMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List memoryMetricsDescriptor = $convert.base64Decode(
    'Cg1NZW1vcnlNZXRyaWNzEh8KC3VzYWdlX2J5dGVzGAEgASgEUgp1c2FnZUJ5dGVzEh8KC2xpbW'
    'l0X2J5dGVzGAIgASgEUgpsaW1pdEJ5dGVzEh8KC2NhY2hlX2J5dGVzGAMgASgEUgpjYWNoZUJ5'
    'dGVzEiMKDXVzYWdlX3BlcmNlbnQYBCABKAFSDHVzYWdlUGVyY2VudA==');

@$core.Deprecated('Use networkMetricsDescriptor instead')
const NetworkMetrics$json = {
  '1': 'NetworkMetrics',
  '2': [
    {'1': 'rx_bytes_total', '3': 1, '4': 1, '5': 4, '10': 'rxBytesTotal'},
    {'1': 'tx_bytes_total', '3': 2, '4': 1, '5': 4, '10': 'txBytesTotal'},
    {'1': 'rx_packets_total', '3': 3, '4': 1, '5': 4, '10': 'rxPacketsTotal'},
    {'1': 'tx_packets_total', '3': 4, '4': 1, '5': 4, '10': 'txPacketsTotal'},
    {'1': 'rx_errors_total', '3': 5, '4': 1, '5': 4, '10': 'rxErrorsTotal'},
    {'1': 'tx_errors_total', '3': 6, '4': 1, '5': 4, '10': 'txErrorsTotal'},
  ],
};

/// Descriptor for `NetworkMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkMetricsDescriptor = $convert.base64Decode(
    'Cg5OZXR3b3JrTWV0cmljcxIkCg5yeF9ieXRlc190b3RhbBgBIAEoBFIMcnhCeXRlc1RvdGFsEi'
    'QKDnR4X2J5dGVzX3RvdGFsGAIgASgEUgx0eEJ5dGVzVG90YWwSKAoQcnhfcGFja2V0c190b3Rh'
    'bBgDIAEoBFIOcnhQYWNrZXRzVG90YWwSKAoQdHhfcGFja2V0c190b3RhbBgEIAEoBFIOdHhQYW'
    'NrZXRzVG90YWwSJgoPcnhfZXJyb3JzX3RvdGFsGAUgASgEUg1yeEVycm9yc1RvdGFsEiYKD3R4'
    'X2Vycm9yc190b3RhbBgGIAEoBFINdHhFcnJvcnNUb3RhbA==');

@$core.Deprecated('Use diskIoMetricsDescriptor instead')
const DiskIoMetrics$json = {
  '1': 'DiskIoMetrics',
  '2': [
    {'1': 'read_bytes_total', '3': 1, '4': 1, '5': 4, '10': 'readBytesTotal'},
    {'1': 'write_bytes_total', '3': 2, '4': 1, '5': 4, '10': 'writeBytesTotal'},
    {'1': 'read_ops_total', '3': 3, '4': 1, '5': 4, '10': 'readOpsTotal'},
    {'1': 'write_ops_total', '3': 4, '4': 1, '5': 4, '10': 'writeOpsTotal'},
  ],
};

/// Descriptor for `DiskIoMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List diskIoMetricsDescriptor = $convert.base64Decode(
    'Cg1EaXNrSW9NZXRyaWNzEigKEHJlYWRfYnl0ZXNfdG90YWwYASABKARSDnJlYWRCeXRlc1RvdG'
    'FsEioKEXdyaXRlX2J5dGVzX3RvdGFsGAIgASgEUg93cml0ZUJ5dGVzVG90YWwSJAoOcmVhZF9v'
    'cHNfdG90YWwYAyABKARSDHJlYWRPcHNUb3RhbBImCg93cml0ZV9vcHNfdG90YWwYBCABKARSDX'
    'dyaXRlT3BzVG90YWw=');

@$core.Deprecated('Use storageMetricsDescriptor instead')
const StorageMetrics$json = {
  '1': 'StorageMetrics',
  '2': [
    {'1': 'total_bytes', '3': 1, '4': 1, '5': 4, '10': 'totalBytes'},
    {'1': 'used_bytes', '3': 2, '4': 1, '5': 4, '10': 'usedBytes'},
    {'1': 'free_bytes', '3': 3, '4': 1, '5': 4, '10': 'freeBytes'},
    {'1': 'usage_percent', '3': 4, '4': 1, '5': 1, '10': 'usagePercent'},
  ],
};

/// Descriptor for `StorageMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storageMetricsDescriptor = $convert.base64Decode(
    'Cg5TdG9yYWdlTWV0cmljcxIfCgt0b3RhbF9ieXRlcxgBIAEoBFIKdG90YWxCeXRlcxIdCgp1c2'
    'VkX2J5dGVzGAIgASgEUgl1c2VkQnl0ZXMSHQoKZnJlZV9ieXRlcxgDIAEoBFIJZnJlZUJ5dGVz'
    'EiMKDXVzYWdlX3BlcmNlbnQYBCABKAFSDHVzYWdlUGVyY2VudA==');

@$core.Deprecated('Use containerMetricsDescriptor instead')
const ContainerMetrics$json = {
  '1': 'ContainerMetrics',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {
      '1': 'collected_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'collectedAt'
    },
    {
      '1': 'cpu',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.CpuMetrics',
      '10': 'cpu'
    },
    {
      '1': 'memory',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.MemoryMetrics',
      '10': 'memory'
    },
    {
      '1': 'network',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.NetworkMetrics',
      '10': 'network'
    },
    {
      '1': 'disk_io',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.DiskIoMetrics',
      '10': 'diskIo'
    },
  ],
};

/// Descriptor for `ContainerMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List containerMetricsDescriptor = $convert.base64Decode(
    'ChBDb250YWluZXJNZXRyaWNzEiEKDGNvbnRhaW5lcl9pZBgBIAEoCVILY29udGFpbmVySWQSPQ'
    'oMY29sbGVjdGVkX2F0GAIgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFILY29sbGVj'
    'dGVkQXQSKQoDY3B1GAMgASgLMhcuZm9yZ2VuZXQudjEuQ3B1TWV0cmljc1IDY3B1EjIKBm1lbW'
    '9yeRgEIAEoCzIaLmZvcmdlbmV0LnYxLk1lbW9yeU1ldHJpY3NSBm1lbW9yeRI1CgduZXR3b3Jr'
    'GAUgASgLMhsuZm9yZ2VuZXQudjEuTmV0d29ya01ldHJpY3NSB25ldHdvcmsSMwoHZGlza19pbx'
    'gGIAEoCzIaLmZvcmdlbmV0LnYxLkRpc2tJb01ldHJpY3NSBmRpc2tJbw==');

@$core.Deprecated('Use systemMetricsDescriptor instead')
const SystemMetrics$json = {
  '1': 'SystemMetrics',
  '2': [
    {
      '1': 'collected_at',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'collectedAt'
    },
    {'1': 'uptime_seconds', '3': 2, '4': 1, '5': 4, '10': 'uptimeSeconds'},
    {
      '1': 'cpu',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.CpuMetrics',
      '10': 'cpu'
    },
    {
      '1': 'memory',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.MemoryMetrics',
      '10': 'memory'
    },
    {
      '1': 'storage',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.StorageMetrics',
      '10': 'storage'
    },
    {
      '1': 'network',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.NetworkMetrics',
      '10': 'network'
    },
    {
      '1': 'disk_io',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.DiskIoMetrics',
      '10': 'diskIo'
    },
    {
      '1': 'active_containers_count',
      '3': 8,
      '4': 1,
      '5': 13,
      '10': 'activeContainersCount'
    },
  ],
};

/// Descriptor for `SystemMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List systemMetricsDescriptor = $convert.base64Decode(
    'Cg1TeXN0ZW1NZXRyaWNzEj0KDGNvbGxlY3RlZF9hdBgBIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi'
    '5UaW1lc3RhbXBSC2NvbGxlY3RlZEF0EiUKDnVwdGltZV9zZWNvbmRzGAIgASgEUg11cHRpbWVT'
    'ZWNvbmRzEikKA2NwdRgDIAEoCzIXLmZvcmdlbmV0LnYxLkNwdU1ldHJpY3NSA2NwdRIyCgZtZW'
    '1vcnkYBCABKAsyGi5mb3JnZW5ldC52MS5NZW1vcnlNZXRyaWNzUgZtZW1vcnkSNQoHc3RvcmFn'
    'ZRgFIAEoCzIbLmZvcmdlbmV0LnYxLlN0b3JhZ2VNZXRyaWNzUgdzdG9yYWdlEjUKB25ldHdvcm'
    'sYBiABKAsyGy5mb3JnZW5ldC52MS5OZXR3b3JrTWV0cmljc1IHbmV0d29yaxIzCgdkaXNrX2lv'
    'GAcgASgLMhouZm9yZ2VuZXQudjEuRGlza0lvTWV0cmljc1IGZGlza0lvEjYKF2FjdGl2ZV9jb2'
    '50YWluZXJzX2NvdW50GAggASgNUhVhY3RpdmVDb250YWluZXJzQ291bnQ=');

@$core.Deprecated('Use getContainerMetricsRequestDescriptor instead')
const GetContainerMetricsRequest$json = {
  '1': 'GetContainerMetricsRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
  ],
};

/// Descriptor for `GetContainerMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContainerMetricsRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDb250YWluZXJNZXRyaWNzUmVxdWVzdBIhCgxjb250YWluZXJfaWQYASABKAlSC2Nvbn'
        'RhaW5lcklk');

@$core.Deprecated('Use getContainerMetricsResponseDescriptor instead')
const GetContainerMetricsResponse$json = {
  '1': 'GetContainerMetricsResponse',
  '2': [
    {
      '1': 'metrics',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.ContainerMetrics',
      '10': 'metrics'
    },
  ],
};

/// Descriptor for `GetContainerMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContainerMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRDb250YWluZXJNZXRyaWNzUmVzcG9uc2USNwoHbWV0cmljcxgBIAEoCzIdLmZvcmdlbm'
        'V0LnYxLkNvbnRhaW5lck1ldHJpY3NSB21ldHJpY3M=');

@$core.Deprecated('Use streamContainerMetricsRequestDescriptor instead')
const StreamContainerMetricsRequest$json = {
  '1': 'StreamContainerMetricsRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {'1': 'interval_seconds', '3': 2, '4': 1, '5': 13, '10': 'intervalSeconds'},
  ],
};

/// Descriptor for `StreamContainerMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamContainerMetricsRequestDescriptor =
    $convert.base64Decode(
        'Ch1TdHJlYW1Db250YWluZXJNZXRyaWNzUmVxdWVzdBIhCgxjb250YWluZXJfaWQYASABKAlSC2'
        'NvbnRhaW5lcklkEikKEGludGVydmFsX3NlY29uZHMYAiABKA1SD2ludGVydmFsU2Vjb25kcw==');

@$core.Deprecated('Use streamContainerMetricsResponseDescriptor instead')
const StreamContainerMetricsResponse$json = {
  '1': 'StreamContainerMetricsResponse',
  '2': [
    {
      '1': 'metrics',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.ContainerMetrics',
      '10': 'metrics'
    },
  ],
};

/// Descriptor for `StreamContainerMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamContainerMetricsResponseDescriptor =
    $convert.base64Decode(
        'Ch5TdHJlYW1Db250YWluZXJNZXRyaWNzUmVzcG9uc2USNwoHbWV0cmljcxgBIAEoCzIdLmZvcm'
        'dlbmV0LnYxLkNvbnRhaW5lck1ldHJpY3NSB21ldHJpY3M=');

@$core.Deprecated('Use queryContainerMetricsRequestDescriptor instead')
const QueryContainerMetricsRequest$json = {
  '1': 'QueryContainerMetricsRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {
      '1': 'range',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.TimeRange',
      '10': 'range'
    },
    {'1': 'step_seconds', '3': 3, '4': 1, '5': 13, '10': 'stepSeconds'},
    {'1': 'max_points', '3': 4, '4': 1, '5': 13, '10': 'maxPoints'},
    {
      '1': 'keys',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.MetricKey',
      '10': 'keys'
    },
  ],
};

/// Descriptor for `QueryContainerMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryContainerMetricsRequestDescriptor = $convert.base64Decode(
    'ChxRdWVyeUNvbnRhaW5lck1ldHJpY3NSZXF1ZXN0EiEKDGNvbnRhaW5lcl9pZBgBIAEoCVILY2'
    '9udGFpbmVySWQSLAoFcmFuZ2UYAiABKAsyFi5mb3JnZW5ldC52MS5UaW1lUmFuZ2VSBXJhbmdl'
    'EiEKDHN0ZXBfc2Vjb25kcxgDIAEoDVILc3RlcFNlY29uZHMSHQoKbWF4X3BvaW50cxgEIAEoDV'
    'IJbWF4UG9pbnRzEioKBGtleXMYBSADKAsyFi5mb3JnZW5ldC52MS5NZXRyaWNLZXlSBGtleXM=');

@$core.Deprecated('Use queryContainerMetricsResponseDescriptor instead')
const QueryContainerMetricsResponse$json = {
  '1': 'QueryContainerMetricsResponse',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {
      '1': 'series',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.MetricSeries',
      '10': 'series'
    },
  ],
};

/// Descriptor for `QueryContainerMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryContainerMetricsResponseDescriptor =
    $convert.base64Decode(
        'Ch1RdWVyeUNvbnRhaW5lck1ldHJpY3NSZXNwb25zZRIhCgxjb250YWluZXJfaWQYASABKAlSC2'
        'NvbnRhaW5lcklkEjEKBnNlcmllcxgCIAMoCzIZLmZvcmdlbmV0LnYxLk1ldHJpY1Nlcmllc1IG'
        'c2VyaWVz');

@$core.Deprecated('Use getSystemMetricsRequestDescriptor instead')
const GetSystemMetricsRequest$json = {
  '1': 'GetSystemMetricsRequest',
};

/// Descriptor for `GetSystemMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSystemMetricsRequestDescriptor =
    $convert.base64Decode('ChdHZXRTeXN0ZW1NZXRyaWNzUmVxdWVzdA==');

@$core.Deprecated('Use getSystemMetricsResponseDescriptor instead')
const GetSystemMetricsResponse$json = {
  '1': 'GetSystemMetricsResponse',
  '2': [
    {
      '1': 'metrics',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.SystemMetrics',
      '10': 'metrics'
    },
  ],
};

/// Descriptor for `GetSystemMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSystemMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRTeXN0ZW1NZXRyaWNzUmVzcG9uc2USNAoHbWV0cmljcxgBIAEoCzIaLmZvcmdlbmV0Ln'
        'YxLlN5c3RlbU1ldHJpY3NSB21ldHJpY3M=');

@$core.Deprecated('Use streamSystemMetricsRequestDescriptor instead')
const StreamSystemMetricsRequest$json = {
  '1': 'StreamSystemMetricsRequest',
  '2': [
    {'1': 'interval_seconds', '3': 1, '4': 1, '5': 13, '10': 'intervalSeconds'},
  ],
};

/// Descriptor for `StreamSystemMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamSystemMetricsRequestDescriptor =
    $convert.base64Decode(
        'ChpTdHJlYW1TeXN0ZW1NZXRyaWNzUmVxdWVzdBIpChBpbnRlcnZhbF9zZWNvbmRzGAEgASgNUg'
        '9pbnRlcnZhbFNlY29uZHM=');

@$core.Deprecated('Use streamSystemMetricsResponseDescriptor instead')
const StreamSystemMetricsResponse$json = {
  '1': 'StreamSystemMetricsResponse',
  '2': [
    {
      '1': 'metrics',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.SystemMetrics',
      '10': 'metrics'
    },
  ],
};

/// Descriptor for `StreamSystemMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamSystemMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChtTdHJlYW1TeXN0ZW1NZXRyaWNzUmVzcG9uc2USNAoHbWV0cmljcxgBIAEoCzIaLmZvcmdlbm'
        'V0LnYxLlN5c3RlbU1ldHJpY3NSB21ldHJpY3M=');

@$core.Deprecated('Use querySystemMetricsRequestDescriptor instead')
const QuerySystemMetricsRequest$json = {
  '1': 'QuerySystemMetricsRequest',
  '2': [
    {
      '1': 'range',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.TimeRange',
      '10': 'range'
    },
    {'1': 'step_seconds', '3': 2, '4': 1, '5': 13, '10': 'stepSeconds'},
    {'1': 'max_points', '3': 3, '4': 1, '5': 13, '10': 'maxPoints'},
    {
      '1': 'keys',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.MetricKey',
      '10': 'keys'
    },
  ],
};

/// Descriptor for `QuerySystemMetricsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List querySystemMetricsRequestDescriptor = $convert.base64Decode(
    'ChlRdWVyeVN5c3RlbU1ldHJpY3NSZXF1ZXN0EiwKBXJhbmdlGAEgASgLMhYuZm9yZ2VuZXQudj'
    'EuVGltZVJhbmdlUgVyYW5nZRIhCgxzdGVwX3NlY29uZHMYAiABKA1SC3N0ZXBTZWNvbmRzEh0K'
    'Cm1heF9wb2ludHMYAyABKA1SCW1heFBvaW50cxIqCgRrZXlzGAQgAygLMhYuZm9yZ2VuZXQudj'
    'EuTWV0cmljS2V5UgRrZXlz');

@$core.Deprecated('Use querySystemMetricsResponseDescriptor instead')
const QuerySystemMetricsResponse$json = {
  '1': 'QuerySystemMetricsResponse',
  '2': [
    {
      '1': 'series',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.MetricSeries',
      '10': 'series'
    },
  ],
};

/// Descriptor for `QuerySystemMetricsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List querySystemMetricsResponseDescriptor =
    $convert.base64Decode(
        'ChpRdWVyeVN5c3RlbU1ldHJpY3NSZXNwb25zZRIxCgZzZXJpZXMYASADKAsyGS5mb3JnZW5ldC'
        '52MS5NZXRyaWNTZXJpZXNSBnNlcmllcw==');
