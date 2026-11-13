// This is a generated file - do not edit.
//
// Generated from forgenet/v1/domain.proto.

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

@$core.Deprecated('Use domainDescriptor instead')
const Domain$json = {
  '1': 'Domain',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'host', '3': 2, '4': 1, '5': 9, '10': 'host'},
    {'1': 'project_id', '3': 3, '4': 1, '5': 9, '10': 'projectId'},
    {'1': 'service', '3': 4, '4': 1, '5': 9, '10': 'service'},
    {'1': 'port', '3': 5, '4': 1, '5': 13, '10': 'port'},
    {'1': 'created_by', '3': 6, '4': 1, '5': 9, '10': 'createdBy'},
    {
      '1': 'created_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Domain`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainDescriptor = $convert.base64Decode(
    'CgZEb21haW4SDgoCaWQYASABKAlSAmlkEhIKBGhvc3QYAiABKAlSBGhvc3QSHQoKcHJvamVjdF'
    '9pZBgDIAEoCVIJcHJvamVjdElkEhgKB3NlcnZpY2UYBCABKAlSB3NlcnZpY2USEgoEcG9ydBgF'
    'IAEoDVIEcG9ydBIdCgpjcmVhdGVkX2J5GAYgASgJUgljcmVhdGVkQnkSOQoKY3JlYXRlZF9hdB'
    'gHIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use listDomainsRequestDescriptor instead')
const ListDomainsRequest$json = {
  '1': 'ListDomainsRequest',
  '2': [
    {'1': 'project_id', '3': 1, '4': 1, '5': 9, '10': 'projectId'},
  ],
};

/// Descriptor for `ListDomainsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDomainsRequestDescriptor =
    $convert.base64Decode(
        'ChJMaXN0RG9tYWluc1JlcXVlc3QSHQoKcHJvamVjdF9pZBgBIAEoCVIJcHJvamVjdElk');

@$core.Deprecated('Use listDomainsResponseDescriptor instead')
const ListDomainsResponse$json = {
  '1': 'ListDomainsResponse',
  '2': [
    {
      '1': 'domains',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.forgenet.v1.Domain',
      '10': 'domains'
    },
  ],
};

/// Descriptor for `ListDomainsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDomainsResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0RG9tYWluc1Jlc3BvbnNlEi0KB2RvbWFpbnMYASADKAsyEy5mb3JnZW5ldC52MS5Eb2'
    '1haW5SB2RvbWFpbnM=');

@$core.Deprecated('Use addDomainRequestDescriptor instead')
const AddDomainRequest$json = {
  '1': 'AddDomainRequest',
  '2': [
    {'1': 'host', '3': 1, '4': 1, '5': 9, '10': 'host'},
    {'1': 'project_id', '3': 2, '4': 1, '5': 9, '10': 'projectId'},
    {'1': 'service', '3': 3, '4': 1, '5': 9, '10': 'service'},
    {'1': 'port', '3': 4, '4': 1, '5': 13, '10': 'port'},
  ],
};

/// Descriptor for `AddDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addDomainRequestDescriptor = $convert.base64Decode(
    'ChBBZGREb21haW5SZXF1ZXN0EhIKBGhvc3QYASABKAlSBGhvc3QSHQoKcHJvamVjdF9pZBgCIA'
    'EoCVIJcHJvamVjdElkEhgKB3NlcnZpY2UYAyABKAlSB3NlcnZpY2USEgoEcG9ydBgEIAEoDVIE'
    'cG9ydA==');

@$core.Deprecated('Use addDomainResponseDescriptor instead')
const AddDomainResponse$json = {
  '1': 'AddDomainResponse',
  '2': [
    {
      '1': 'domain',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.Domain',
      '10': 'domain'
    },
  ],
};

/// Descriptor for `AddDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addDomainResponseDescriptor = $convert.base64Decode(
    'ChFBZGREb21haW5SZXNwb25zZRIrCgZkb21haW4YASABKAsyEy5mb3JnZW5ldC52MS5Eb21haW'
    '5SBmRvbWFpbg==');

@$core.Deprecated('Use removeDomainRequestDescriptor instead')
const RemoveDomainRequest$json = {
  '1': 'RemoveDomainRequest',
  '2': [
    {'1': 'domain_id', '3': 1, '4': 1, '5': 9, '10': 'domainId'},
  ],
};

/// Descriptor for `RemoveDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeDomainRequestDescriptor =
    $convert.base64Decode(
        'ChNSZW1vdmVEb21haW5SZXF1ZXN0EhsKCWRvbWFpbl9pZBgBIAEoCVIIZG9tYWluSWQ=');

@$core.Deprecated('Use removeDomainResponseDescriptor instead')
const RemoveDomainResponse$json = {
  '1': 'RemoveDomainResponse',
};

/// Descriptor for `RemoveDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeDomainResponseDescriptor =
    $convert.base64Decode('ChRSZW1vdmVEb21haW5SZXNwb25zZQ==');
