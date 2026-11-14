// This is a generated file - do not edit.
//
// Generated from forgenet/v1/update.proto.

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

@$core.Deprecated('Use updateStepDescriptor instead')
const UpdateStep$json = {
  '1': 'UpdateStep',
  '2': [
    {'1': 'UPDATE_STEP_UNSPECIFIED', '2': 0},
    {'1': 'UPDATE_STEP_CHECK', '2': 1},
    {'1': 'UPDATE_STEP_PULL', '2': 2},
    {'1': 'UPDATE_STEP_BACKUP', '2': 3},
    {'1': 'UPDATE_STEP_SWITCH', '2': 4},
    {'1': 'UPDATE_STEP_HEALTH', '2': 5},
    {'1': 'UPDATE_STEP_ROLLBACK', '2': 6},
    {'1': 'UPDATE_STEP_DONE', '2': 7},
  ],
};

/// Descriptor for `UpdateStep`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List updateStepDescriptor = $convert.base64Decode(
    'CgpVcGRhdGVTdGVwEhsKF1VQREFURV9TVEVQX1VOU1BFQ0lGSUVEEAASFQoRVVBEQVRFX1NURV'
    'BfQ0hFQ0sQARIUChBVUERBVEVfU1RFUF9QVUxMEAISFgoSVVBEQVRFX1NURVBfQkFDS1VQEAMS'
    'FgoSVVBEQVRFX1NURVBfU1dJVENIEAQSFgoSVVBEQVRFX1NURVBfSEVBTFRIEAUSGAoUVVBEQV'
    'RFX1NURVBfUk9MTEJBQ0sQBhIUChBVUERBVEVfU1RFUF9ET05FEAc=');

@$core.Deprecated('Use updateResultDescriptor instead')
const UpdateResult$json = {
  '1': 'UpdateResult',
  '2': [
    {'1': 'UPDATE_RESULT_UNSPECIFIED', '2': 0},
    {'1': 'UPDATE_RESULT_RUNNING', '2': 1},
    {'1': 'UPDATE_RESULT_SUCCEEDED', '2': 2},
    {'1': 'UPDATE_RESULT_ROLLED_BACK', '2': 3},
    {'1': 'UPDATE_RESULT_FAILED', '2': 4},
  ],
};

/// Descriptor for `UpdateResult`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List updateResultDescriptor = $convert.base64Decode(
    'CgxVcGRhdGVSZXN1bHQSHQoZVVBEQVRFX1JFU1VMVF9VTlNQRUNJRklFRBAAEhkKFVVQREFURV'
    '9SRVNVTFRfUlVOTklORxABEhsKF1VQREFURV9SRVNVTFRfU1VDQ0VFREVEEAISHQoZVVBEQVRF'
    'X1JFU1VMVF9ST0xMRURfQkFDSxADEhgKFFVQREFURV9SRVNVTFRfRkFJTEVEEAQ=');

@$core.Deprecated('Use availableReleaseDescriptor instead')
const AvailableRelease$json = {
  '1': 'AvailableRelease',
  '2': [
    {'1': 'version', '3': 1, '4': 1, '5': 9, '10': 'version'},
    {'1': 'notes_url', '3': 2, '4': 1, '5': 9, '10': 'notesUrl'},
  ],
};

/// Descriptor for `AvailableRelease`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List availableReleaseDescriptor = $convert.base64Decode(
    'ChBBdmFpbGFibGVSZWxlYXNlEhgKB3ZlcnNpb24YASABKAlSB3ZlcnNpb24SGwoJbm90ZXNfdX'
    'JsGAIgASgJUghub3Rlc1VybA==');

@$core.Deprecated('Use updateRunDescriptor instead')
const UpdateRun$json = {
  '1': 'UpdateRun',
  '2': [
    {'1': 'from_version', '3': 1, '4': 1, '5': 9, '10': 'fromVersion'},
    {'1': 'to_version', '3': 2, '4': 1, '5': 9, '10': 'toVersion'},
    {
      '1': 'step',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.UpdateStep',
      '10': 'step'
    },
    {
      '1': 'result',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.UpdateResult',
      '10': 'result'
    },
    {'1': 'reason', '3': 5, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'backup', '3': 6, '4': 1, '5': 9, '10': 'backup'},
    {
      '1': 'started_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'finished_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'finishedAt'
    },
  ],
};

/// Descriptor for `UpdateRun`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRunDescriptor = $convert.base64Decode(
    'CglVcGRhdGVSdW4SIQoMZnJvbV92ZXJzaW9uGAEgASgJUgtmcm9tVmVyc2lvbhIdCgp0b192ZX'
    'JzaW9uGAIgASgJUgl0b1ZlcnNpb24SKwoEc3RlcBgDIAEoDjIXLmZvcmdlbmV0LnYxLlVwZGF0'
    'ZVN0ZXBSBHN0ZXASMQoGcmVzdWx0GAQgASgOMhkuZm9yZ2VuZXQudjEuVXBkYXRlUmVzdWx0Ug'
    'ZyZXN1bHQSFgoGcmVhc29uGAUgASgJUgZyZWFzb24SFgoGYmFja3VwGAYgASgJUgZiYWNrdXAS'
    'OQoKc3RhcnRlZF9hdBgHIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0ZW'
    'RBdBI7CgtmaW5pc2hlZF9hdBgIIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCmZp'
    'bmlzaGVkQXQ=');

@$core.Deprecated('Use getUpdateStatusRequestDescriptor instead')
const GetUpdateStatusRequest$json = {
  '1': 'GetUpdateStatusRequest',
};

/// Descriptor for `GetUpdateStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUpdateStatusRequestDescriptor =
    $convert.base64Decode('ChZHZXRVcGRhdGVTdGF0dXNSZXF1ZXN0');

@$core.Deprecated('Use getUpdateStatusResponseDescriptor instead')
const GetUpdateStatusResponse$json = {
  '1': 'GetUpdateStatusResponse',
  '2': [
    {'1': 'current_version', '3': 1, '4': 1, '5': 9, '10': 'currentVersion'},
    {
      '1': 'available',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.AvailableRelease',
      '10': 'available'
    },
    {
      '1': 'checked_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'checkedAt'
    },
    {'1': 'check_error', '3': 4, '4': 1, '5': 9, '10': 'checkError'},
    {'1': 'check_enabled', '3': 5, '4': 1, '5': 8, '10': 'checkEnabled'},
    {'1': 'supported', '3': 6, '4': 1, '5': 8, '10': 'supported'},
    {
      '1': 'last_run',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.UpdateRun',
      '10': 'lastRun'
    },
  ],
};

/// Descriptor for `GetUpdateStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUpdateStatusResponseDescriptor = $convert.base64Decode(
    'ChdHZXRVcGRhdGVTdGF0dXNSZXNwb25zZRInCg9jdXJyZW50X3ZlcnNpb24YASABKAlSDmN1cn'
    'JlbnRWZXJzaW9uEjsKCWF2YWlsYWJsZRgCIAEoCzIdLmZvcmdlbmV0LnYxLkF2YWlsYWJsZVJl'
    'bGVhc2VSCWF2YWlsYWJsZRI5CgpjaGVja2VkX2F0GAMgASgLMhouZ29vZ2xlLnByb3RvYnVmLl'
    'RpbWVzdGFtcFIJY2hlY2tlZEF0Eh8KC2NoZWNrX2Vycm9yGAQgASgJUgpjaGVja0Vycm9yEiMK'
    'DWNoZWNrX2VuYWJsZWQYBSABKAhSDGNoZWNrRW5hYmxlZBIcCglzdXBwb3J0ZWQYBiABKAhSCX'
    'N1cHBvcnRlZBIxCghsYXN0X3J1bhgHIAEoCzIWLmZvcmdlbmV0LnYxLlVwZGF0ZVJ1blIHbGFz'
    'dFJ1bg==');

@$core.Deprecated('Use checkForUpdatesRequestDescriptor instead')
const CheckForUpdatesRequest$json = {
  '1': 'CheckForUpdatesRequest',
};

/// Descriptor for `CheckForUpdatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkForUpdatesRequestDescriptor =
    $convert.base64Decode('ChZDaGVja0ZvclVwZGF0ZXNSZXF1ZXN0');

@$core.Deprecated('Use checkForUpdatesResponseDescriptor instead')
const CheckForUpdatesResponse$json = {
  '1': 'CheckForUpdatesResponse',
  '2': [
    {
      '1': 'status',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.GetUpdateStatusResponse',
      '10': 'status'
    },
  ],
};

/// Descriptor for `CheckForUpdatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkForUpdatesResponseDescriptor =
    $convert.base64Decode(
        'ChdDaGVja0ZvclVwZGF0ZXNSZXNwb25zZRI8CgZzdGF0dXMYASABKAsyJC5mb3JnZW5ldC52MS'
        '5HZXRVcGRhdGVTdGF0dXNSZXNwb25zZVIGc3RhdHVz');

@$core.Deprecated('Use startUpdateRequestDescriptor instead')
const StartUpdateRequest$json = {
  '1': 'StartUpdateRequest',
};

/// Descriptor for `StartUpdateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startUpdateRequestDescriptor =
    $convert.base64Decode('ChJTdGFydFVwZGF0ZVJlcXVlc3Q=');

@$core.Deprecated('Use startUpdateResponseDescriptor instead')
const StartUpdateResponse$json = {
  '1': 'StartUpdateResponse',
  '2': [
    {'1': 'to_version', '3': 1, '4': 1, '5': 9, '10': 'toVersion'},
  ],
};

/// Descriptor for `StartUpdateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startUpdateResponseDescriptor = $convert.base64Decode(
    'ChNTdGFydFVwZGF0ZVJlc3BvbnNlEh0KCnRvX3ZlcnNpb24YASABKAlSCXRvVmVyc2lvbg==');
