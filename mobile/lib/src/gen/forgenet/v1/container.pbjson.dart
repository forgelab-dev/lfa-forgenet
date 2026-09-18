// This is a generated file - do not edit.
//
// Generated from forgenet/v1/container.proto.

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

@$core.Deprecated('Use containerStateDescriptor instead')
const ContainerState$json = {
  '1': 'ContainerState',
  '2': [
    {'1': 'CONTAINER_STATE_UNSPECIFIED', '2': 0},
    {'1': 'CONTAINER_STATE_CREATED', '2': 1},
    {'1': 'CONTAINER_STATE_RUNNING', '2': 2},
    {'1': 'CONTAINER_STATE_PAUSED', '2': 3},
    {'1': 'CONTAINER_STATE_RESTARTING', '2': 4},
    {'1': 'CONTAINER_STATE_EXITED', '2': 5},
    {'1': 'CONTAINER_STATE_DEAD', '2': 6},
  ],
};

/// Descriptor for `ContainerState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List containerStateDescriptor = $convert.base64Decode(
    'Cg5Db250YWluZXJTdGF0ZRIfChtDT05UQUlORVJfU1RBVEVfVU5TUEVDSUZJRUQQABIbChdDT0'
    '5UQUlORVJfU1RBVEVfQ1JFQVRFRBABEhsKF0NPTlRBSU5FUl9TVEFURV9SVU5OSU5HEAISGgoW'
    'Q09OVEFJTkVSX1NUQVRFX1BBVVNFRBADEh4KGkNPTlRBSU5FUl9TVEFURV9SRVNUQVJUSU5HEA'
    'QSGgoWQ09OVEFJTkVSX1NUQVRFX0VYSVRFRBAFEhgKFENPTlRBSU5FUl9TVEFURV9ERUFEEAY=');

@$core.Deprecated('Use logStreamDescriptor instead')
const LogStream$json = {
  '1': 'LogStream',
  '2': [
    {'1': 'LOG_STREAM_UNSPECIFIED', '2': 0},
    {'1': 'LOG_STREAM_STDERR', '2': 1},
    {'1': 'LOG_STREAM_STDOUT', '2': 2},
  ],
};

/// Descriptor for `LogStream`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logStreamDescriptor = $convert.base64Decode(
    'CglMb2dTdHJlYW0SGgoWTE9HX1NUUkVBTV9VTlNQRUNJRklFRBAAEhUKEUxPR19TVFJFQU1fU1'
    'RERVJSEAESFQoRTE9HX1NUUkVBTV9TVERPVVQQAg==');

@$core.Deprecated('Use portMappingDescriptor instead')
const PortMapping$json = {
  '1': 'PortMapping',
  '2': [
    {'1': 'host_port', '3': 1, '4': 1, '5': 13, '10': 'hostPort'},
    {'1': 'container_port', '3': 2, '4': 1, '5': 13, '10': 'containerPort'},
    {'1': 'protocol', '3': 3, '4': 1, '5': 9, '10': 'protocol'},
  ],
};

/// Descriptor for `PortMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List portMappingDescriptor = $convert.base64Decode(
    'CgtQb3J0TWFwcGluZxIbCglob3N0X3BvcnQYASABKA1SCGhvc3RQb3J0EiUKDmNvbnRhaW5lcl'
    '9wb3J0GAIgASgNUg1jb250YWluZXJQb3J0EhoKCHByb3RvY29sGAMgASgJUghwcm90b2NvbA==');

@$core.Deprecated('Use containerDescriptor instead')
const Container$json = {
  '1': 'Container',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'image', '3': 3, '4': 1, '5': 9, '10': 'image'},
    {
      '1': 'state',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.ContainerState',
      '10': 'state'
    },
    {
      '1': 'ports',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.PortMapping',
      '10': 'ports'
    },
    {
      '1': 'created_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Container`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List containerDescriptor = $convert.base64Decode(
    'CglDb250YWluZXISDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWUSFAoFaW1hZ2'
    'UYAyABKAlSBWltYWdlEjEKBXN0YXRlGAQgASgOMhsuZm9yZ2VuZXQudjEuQ29udGFpbmVyU3Rh'
    'dGVSBXN0YXRlEi4KBXBvcnRzGAUgAygLMhguZm9yZ2VuZXQudjEuUG9ydE1hcHBpbmdSBXBvcn'
    'RzEjkKCmNyZWF0ZWRfYXQYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVh'
    'dGVkQXQ=');

@$core.Deprecated('Use listContainersRequestDescriptor instead')
const ListContainersRequest$json = {
  '1': 'ListContainersRequest',
  '2': [
    {'1': 'include_stopped', '3': 1, '4': 1, '5': 8, '10': 'includeStopped'},
    {'1': 'page_size', '3': 2, '4': 1, '5': 5, '10': 'pageSize'},
    {'1': 'page_token', '3': 3, '4': 1, '5': 9, '10': 'pageToken'},
  ],
};

/// Descriptor for `ListContainersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContainersRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0Q29udGFpbmVyc1JlcXVlc3QSJwoPaW5jbHVkZV9zdG9wcGVkGAEgASgIUg5pbmNsdW'
    'RlU3RvcHBlZBIbCglwYWdlX3NpemUYAiABKAVSCHBhZ2VTaXplEh0KCnBhZ2VfdG9rZW4YAyAB'
    'KAlSCXBhZ2VUb2tlbg==');

@$core.Deprecated('Use listContainersResponseDescriptor instead')
const ListContainersResponse$json = {
  '1': 'ListContainersResponse',
  '2': [
    {
      '1': 'containers',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.Container',
      '10': 'containers'
    },
    {'1': 'next_page_token', '3': 2, '4': 1, '5': 9, '10': 'nextPageToken'},
  ],
};

/// Descriptor for `ListContainersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContainersResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0Q29udGFpbmVyc1Jlc3BvbnNlEjYKCmNvbnRhaW5lcnMYASADKAsyFi5mb3JnZW5ldC'
    '52MS5Db250YWluZXJSCmNvbnRhaW5lcnMSJgoPbmV4dF9wYWdlX3Rva2VuGAIgASgJUg1uZXh0'
    'UGFnZVRva2Vu');

@$core.Deprecated('Use getContainerRequestDescriptor instead')
const GetContainerRequest$json = {
  '1': 'GetContainerRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
  ],
};

/// Descriptor for `GetContainerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContainerRequestDescriptor = $convert.base64Decode(
    'ChNHZXRDb250YWluZXJSZXF1ZXN0EiEKDGNvbnRhaW5lcl9pZBgBIAEoCVILY29udGFpbmVySW'
    'Q=');

@$core.Deprecated('Use getContainerResponseDescriptor instead')
const GetContainerResponse$json = {
  '1': 'GetContainerResponse',
  '2': [
    {
      '1': 'container',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Container',
      '10': 'container'
    },
  ],
};

/// Descriptor for `GetContainerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContainerResponseDescriptor = $convert.base64Decode(
    'ChRHZXRDb250YWluZXJSZXNwb25zZRI0Cgljb250YWluZXIYASABKAsyFi5mb3JnZW5ldC52MS'
    '5Db250YWluZXJSCWNvbnRhaW5lcg==');

@$core.Deprecated('Use startContainerRequestDescriptor instead')
const StartContainerRequest$json = {
  '1': 'StartContainerRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
  ],
};

/// Descriptor for `StartContainerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startContainerRequestDescriptor = $convert.base64Decode(
    'ChVTdGFydENvbnRhaW5lclJlcXVlc3QSIQoMY29udGFpbmVyX2lkGAEgASgJUgtjb250YWluZX'
    'JJZA==');

@$core.Deprecated('Use startContainerResponseDescriptor instead')
const StartContainerResponse$json = {
  '1': 'StartContainerResponse',
  '2': [
    {
      '1': 'container',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Container',
      '10': 'container'
    },
  ],
};

/// Descriptor for `StartContainerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startContainerResponseDescriptor =
    $convert.base64Decode(
        'ChZTdGFydENvbnRhaW5lclJlc3BvbnNlEjQKCWNvbnRhaW5lchgBIAEoCzIWLmZvcmdlbmV0Ln'
        'YxLkNvbnRhaW5lclIJY29udGFpbmVy');

@$core.Deprecated('Use stopContainerRequestDescriptor instead')
const StopContainerRequest$json = {
  '1': 'StopContainerRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {'1': 'timeout_seconds', '3': 2, '4': 1, '5': 13, '10': 'timeoutSeconds'},
  ],
};

/// Descriptor for `StopContainerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopContainerRequestDescriptor = $convert.base64Decode(
    'ChRTdG9wQ29udGFpbmVyUmVxdWVzdBIhCgxjb250YWluZXJfaWQYASABKAlSC2NvbnRhaW5lck'
    'lkEicKD3RpbWVvdXRfc2Vjb25kcxgCIAEoDVIOdGltZW91dFNlY29uZHM=');

@$core.Deprecated('Use stopContainerResponseDescriptor instead')
const StopContainerResponse$json = {
  '1': 'StopContainerResponse',
  '2': [
    {
      '1': 'container',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Container',
      '10': 'container'
    },
  ],
};

/// Descriptor for `StopContainerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopContainerResponseDescriptor = $convert.base64Decode(
    'ChVTdG9wQ29udGFpbmVyUmVzcG9uc2USNAoJY29udGFpbmVyGAEgASgLMhYuZm9yZ2VuZXQudj'
    'EuQ29udGFpbmVyUgljb250YWluZXI=');

@$core.Deprecated('Use restartContainerRequestDescriptor instead')
const RestartContainerRequest$json = {
  '1': 'RestartContainerRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
  ],
};

/// Descriptor for `RestartContainerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restartContainerRequestDescriptor =
    $convert.base64Decode(
        'ChdSZXN0YXJ0Q29udGFpbmVyUmVxdWVzdBIhCgxjb250YWluZXJfaWQYASABKAlSC2NvbnRhaW'
        '5lcklk');

@$core.Deprecated('Use restartContainerResponseDescriptor instead')
const RestartContainerResponse$json = {
  '1': 'RestartContainerResponse',
  '2': [
    {
      '1': 'container',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Container',
      '10': 'container'
    },
  ],
};

/// Descriptor for `RestartContainerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restartContainerResponseDescriptor =
    $convert.base64Decode(
        'ChhSZXN0YXJ0Q29udGFpbmVyUmVzcG9uc2USNAoJY29udGFpbmVyGAEgASgLMhYuZm9yZ2VuZX'
        'QudjEuQ29udGFpbmVyUgljb250YWluZXI=');

@$core.Deprecated('Use streamContainerLogsRequestDescriptor instead')
const StreamContainerLogsRequest$json = {
  '1': 'StreamContainerLogsRequest',
  '2': [
    {'1': 'container_id', '3': 1, '4': 1, '5': 9, '10': 'containerId'},
    {'1': 'follow', '3': 2, '4': 1, '5': 8, '10': 'follow'},
    {'1': 'tail_lines', '3': 3, '4': 1, '5': 5, '10': 'tailLines'},
    {
      '1': 'since',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'since'
    },
    {
      '1': 'stream',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.LogStream',
      '10': 'stream'
    },
  ],
};

/// Descriptor for `StreamContainerLogsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamContainerLogsRequestDescriptor = $convert.base64Decode(
    'ChpTdHJlYW1Db250YWluZXJMb2dzUmVxdWVzdBIhCgxjb250YWluZXJfaWQYASABKAlSC2Nvbn'
    'RhaW5lcklkEhYKBmZvbGxvdxgCIAEoCFIGZm9sbG93Eh0KCnRhaWxfbGluZXMYAyABKAVSCXRh'
    'aWxMaW5lcxIwCgVzaW5jZRgEIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSBXNpbm'
    'NlEi4KBnN0cmVhbRgFIAEoDjIWLmZvcmdlbmV0LnYxLkxvZ1N0cmVhbVIGc3RyZWFt');

@$core.Deprecated('Use streamContainerLogsResponseDescriptor instead')
const StreamContainerLogsResponse$json = {
  '1': 'StreamContainerLogsResponse',
  '2': [
    {'1': 'chunk', '3': 1, '4': 1, '5': 12, '10': 'chunk'},
    {
      '1': 'emitted',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'emitted'
    },
    {
      '1': 'source',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.LogStream',
      '10': 'source'
    },
  ],
};

/// Descriptor for `StreamContainerLogsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamContainerLogsResponseDescriptor =
    $convert.base64Decode(
        'ChtTdHJlYW1Db250YWluZXJMb2dzUmVzcG9uc2USFAoFY2h1bmsYASABKAxSBWNodW5rEjQKB2'
        'VtaXR0ZWQYAiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgdlbWl0dGVkEi4KBnNv'
        'dXJjZRgDIAEoDjIWLmZvcmdlbmV0LnYxLkxvZ1N0cmVhbVIGc291cmNl');
