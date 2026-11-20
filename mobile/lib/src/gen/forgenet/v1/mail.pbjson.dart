// This is a generated file - do not edit.
//
// Generated from forgenet/v1/mail.proto.

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

@$core.Deprecated('Use mailTLSDescriptor instead')
const MailTLS$json = {
  '1': 'MailTLS',
  '2': [
    {'1': 'MAIL_TLS_UNSPECIFIED', '2': 0},
    {'1': 'MAIL_TLS_STARTTLS', '2': 1},
    {'1': 'MAIL_TLS_IMPLICIT', '2': 2},
    {'1': 'MAIL_TLS_NONE', '2': 3},
  ],
};

/// Descriptor for `MailTLS`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mailTLSDescriptor = $convert.base64Decode(
    'CgdNYWlsVExTEhgKFE1BSUxfVExTX1VOU1BFQ0lGSUVEEAASFQoRTUFJTF9UTFNfU1RBUlRUTF'
    'MQARIVChFNQUlMX1RMU19JTVBMSUNJVBACEhEKDU1BSUxfVExTX05PTkUQAw==');

@$core.Deprecated('Use mailSourceDescriptor instead')
const MailSource$json = {
  '1': 'MailSource',
  '2': [
    {'1': 'MAIL_SOURCE_UNSPECIFIED', '2': 0},
    {'1': 'MAIL_SOURCE_SERVER', '2': 1},
    {'1': 'MAIL_SOURCE_DASHBOARD', '2': 2},
  ],
};

/// Descriptor for `MailSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mailSourceDescriptor = $convert.base64Decode(
    'CgpNYWlsU291cmNlEhsKF01BSUxfU09VUkNFX1VOU1BFQ0lGSUVEEAASFgoSTUFJTF9TT1VSQ0'
    'VfU0VSVkVSEAESGQoVTUFJTF9TT1VSQ0VfREFTSEJPQVJEEAI=');

@$core.Deprecated('Use getMailStatusRequestDescriptor instead')
const GetMailStatusRequest$json = {
  '1': 'GetMailStatusRequest',
};

/// Descriptor for `GetMailStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMailStatusRequestDescriptor =
    $convert.base64Decode('ChRHZXRNYWlsU3RhdHVzUmVxdWVzdA==');

@$core.Deprecated('Use getMailStatusResponseDescriptor instead')
const GetMailStatusResponse$json = {
  '1': 'GetMailStatusResponse',
  '2': [
    {'1': 'enabled', '3': 1, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'host', '3': 2, '4': 1, '5': 9, '10': 'host'},
    {'1': 'port', '3': 3, '4': 1, '5': 13, '10': 'port'},
    {
      '1': 'tls',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MailTLS',
      '10': 'tls'
    },
    {'1': 'from', '3': 5, '4': 1, '5': 9, '10': 'from'},
    {'1': 'authenticated', '3': 6, '4': 1, '5': 8, '10': 'authenticated'},
    {
      '1': 'source',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MailSource',
      '10': 'source'
    },
    {'1': 'editable', '3': 8, '4': 1, '5': 8, '10': 'editable'},
    {'1': 'username', '3': 9, '4': 1, '5': 9, '10': 'username'},
    {'1': 'disabled', '3': 10, '4': 1, '5': 8, '10': 'disabled'},
    {
      '1': 'password_unreadable',
      '3': 11,
      '4': 1,
      '5': 8,
      '10': 'passwordUnreadable'
    },
    {
      '1': 'updated_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    {'1': 'updated_by', '3': 13, '4': 1, '5': 9, '10': 'updatedBy'},
  ],
};

/// Descriptor for `GetMailStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMailStatusResponseDescriptor = $convert.base64Decode(
    'ChVHZXRNYWlsU3RhdHVzUmVzcG9uc2USGAoHZW5hYmxlZBgBIAEoCFIHZW5hYmxlZBISCgRob3'
    'N0GAIgASgJUgRob3N0EhIKBHBvcnQYAyABKA1SBHBvcnQSJgoDdGxzGAQgASgOMhQuZm9yZ2Vu'
    'ZXQudjEuTWFpbFRMU1IDdGxzEhIKBGZyb20YBSABKAlSBGZyb20SJAoNYXV0aGVudGljYXRlZB'
    'gGIAEoCFINYXV0aGVudGljYXRlZBIvCgZzb3VyY2UYByABKA4yFy5mb3JnZW5ldC52MS5NYWls'
    'U291cmNlUgZzb3VyY2USGgoIZWRpdGFibGUYCCABKAhSCGVkaXRhYmxlEhoKCHVzZXJuYW1lGA'
    'kgASgJUgh1c2VybmFtZRIaCghkaXNhYmxlZBgKIAEoCFIIZGlzYWJsZWQSLwoTcGFzc3dvcmRf'
    'dW5yZWFkYWJsZRgLIAEoCFIScGFzc3dvcmRVbnJlYWRhYmxlEjkKCnVwZGF0ZWRfYXQYDCABKA'
    'syGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgl1cGRhdGVkQXQSHQoKdXBkYXRlZF9ieRgN'
    'IAEoCVIJdXBkYXRlZEJ5');

@$core.Deprecated('Use sendTestMailRequestDescriptor instead')
const SendTestMailRequest$json = {
  '1': 'SendTestMailRequest',
};

/// Descriptor for `SendTestMailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTestMailRequestDescriptor =
    $convert.base64Decode('ChNTZW5kVGVzdE1haWxSZXF1ZXN0');

@$core.Deprecated('Use sendTestMailResponseDescriptor instead')
const SendTestMailResponse$json = {
  '1': 'SendTestMailResponse',
  '2': [
    {'1': 'to', '3': 1, '4': 1, '5': 9, '10': 'to'},
  ],
};

/// Descriptor for `SendTestMailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTestMailResponseDescriptor = $convert
    .base64Decode('ChRTZW5kVGVzdE1haWxSZXNwb25zZRIOCgJ0bxgBIAEoCVICdG8=');

@$core.Deprecated('Use saveMailSettingsRequestDescriptor instead')
const SaveMailSettingsRequest$json = {
  '1': 'SaveMailSettingsRequest',
  '2': [
    {'1': 'host', '3': 1, '4': 1, '5': 9, '10': 'host'},
    {'1': 'port', '3': 2, '4': 1, '5': 13, '10': 'port'},
    {
      '1': 'tls',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.forgenet.v1.MailTLS',
      '10': 'tls'
    },
    {'1': 'username', '3': 4, '4': 1, '5': 9, '10': 'username'},
    {'1': 'from', '3': 5, '4': 1, '5': 9, '10': 'from'},
    {'1': 'keep_password', '3': 6, '4': 1, '5': 8, '10': 'keepPassword'},
    {'1': 'password', '3': 7, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `SaveMailSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List saveMailSettingsRequestDescriptor = $convert.base64Decode(
    'ChdTYXZlTWFpbFNldHRpbmdzUmVxdWVzdBISCgRob3N0GAEgASgJUgRob3N0EhIKBHBvcnQYAi'
    'ABKA1SBHBvcnQSJgoDdGxzGAMgASgOMhQuZm9yZ2VuZXQudjEuTWFpbFRMU1IDdGxzEhoKCHVz'
    'ZXJuYW1lGAQgASgJUgh1c2VybmFtZRISCgRmcm9tGAUgASgJUgRmcm9tEiMKDWtlZXBfcGFzc3'
    'dvcmQYBiABKAhSDGtlZXBQYXNzd29yZBIaCghwYXNzd29yZBgHIAEoCVIIcGFzc3dvcmQ=');

@$core.Deprecated('Use saveMailSettingsResponseDescriptor instead')
const SaveMailSettingsResponse$json = {
  '1': 'SaveMailSettingsResponse',
  '2': [
    {
      '1': 'status',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.GetMailStatusResponse',
      '10': 'status'
    },
    {'1': 'to', '3': 2, '4': 1, '5': 9, '10': 'to'},
  ],
};

/// Descriptor for `SaveMailSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List saveMailSettingsResponseDescriptor =
    $convert.base64Decode(
        'ChhTYXZlTWFpbFNldHRpbmdzUmVzcG9uc2USOgoGc3RhdHVzGAEgASgLMiIuZm9yZ2VuZXQudj'
        'EuR2V0TWFpbFN0YXR1c1Jlc3BvbnNlUgZzdGF0dXMSDgoCdG8YAiABKAlSAnRv');

@$core.Deprecated('Use disableMailRequestDescriptor instead')
const DisableMailRequest$json = {
  '1': 'DisableMailRequest',
};

/// Descriptor for `DisableMailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableMailRequestDescriptor =
    $convert.base64Decode('ChJEaXNhYmxlTWFpbFJlcXVlc3Q=');

@$core.Deprecated('Use disableMailResponseDescriptor instead')
const DisableMailResponse$json = {
  '1': 'DisableMailResponse',
  '2': [
    {
      '1': 'status',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.GetMailStatusResponse',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DisableMailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableMailResponseDescriptor = $convert.base64Decode(
    'ChNEaXNhYmxlTWFpbFJlc3BvbnNlEjoKBnN0YXR1cxgBIAEoCzIiLmZvcmdlbmV0LnYxLkdldE'
    '1haWxTdGF0dXNSZXNwb25zZVIGc3RhdHVz');

@$core.Deprecated('Use resetMailSettingsRequestDescriptor instead')
const ResetMailSettingsRequest$json = {
  '1': 'ResetMailSettingsRequest',
};

/// Descriptor for `ResetMailSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetMailSettingsRequestDescriptor =
    $convert.base64Decode('ChhSZXNldE1haWxTZXR0aW5nc1JlcXVlc3Q=');

@$core.Deprecated('Use resetMailSettingsResponseDescriptor instead')
const ResetMailSettingsResponse$json = {
  '1': 'ResetMailSettingsResponse',
  '2': [
    {
      '1': 'status',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.forgenet.v1.GetMailStatusResponse',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ResetMailSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetMailSettingsResponseDescriptor =
    $convert.base64Decode(
        'ChlSZXNldE1haWxTZXR0aW5nc1Jlc3BvbnNlEjoKBnN0YXR1cxgBIAEoCzIiLmZvcmdlbmV0Ln'
        'YxLkdldE1haWxTdGF0dXNSZXNwb25zZVIGc3RhdHVz');
